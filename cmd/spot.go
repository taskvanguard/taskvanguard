package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/taskvanguard/taskvanguard/internal/goals"
	"github.com/taskvanguard/taskvanguard/internal/llm"
	"github.com/taskvanguard/taskvanguard/internal/taskwarrior"
	"github.com/taskvanguard/taskvanguard/pkg/theme"
	"github.com/taskvanguard/taskvanguard/pkg/types"
)

type SpotlightResult struct {
	TaskID     int    `json:"task_id"`
	Title      string `json:"title"`
	Reason     string `json:"reason"`
	History    string `json:"history"`
	Goal	   string `json:"goal"`
	Estimated  string `json:"estimated"`
	ContextTag string `json:"context_tag"`
	Next 	   string `json:"next"`
}

type TaskContext struct {
	Mood      string    `json:"mood"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}

var spotCmd = &cobra.Command{
	Use:   "spot",
	Short: "Analyze urgent tasks and suggests the best one",
	Long: `Analyzes the tasklist, presents a single, high impact task and helps with tackling it. 
	Reframes it, breaks it down into Microtasks and gives some beneficial infos`,
	Run: runSpot,
}

func init() {
	spotCmd.Flags().String("mood", "", "Set mood context (tired, energetic, focused, etc.). Overrides TASKVANGUARD_MOOD")
	spotCmd.Flags().String("context", "", "Set location context (home, office, travel, etc.). Overrides TASKVANGUARD_LOCATION")
	spotCmd.Flags().Bool("refresh", false, "Ignore environment variables and ask fresh")
	spotCmd.Flags().Bool("no-prompt", false, "Run in passive mode without prompts")
}

func runSpot(cmd *cobra.Command, args []string) {
	noPrompt, _ := cmd.Flags().GetBool("no-prompt")
	moodFlag, _ := cmd.Flags().GetString("mood")
	contextFlag, _ := cmd.Flags().GetString("context")
	refresh, _ := cmd.Flags().GetBool("refresh")

	env, err := taskwarrior.Bootstrap(cmd)
	if err != nil {
		fmt.Println(theme.Error(err.Error()))
		return
	}

	if noPrompt {
		runPassiveSpotlight(env.Client, env.Config, moodFlag, contextFlag, refresh, args)
	} else {
		runInteractiveSpotlight(env.Client, env.Config, moodFlag, contextFlag, refresh, args)
	}
}

func runPassiveSpotlight(client *taskwarrior.Client, cfg *types.Config, moodFlag string, contextFlag string, refresh bool, filterArgs []string) {
	taskContext := loadContextFromEnv(moodFlag, contextFlag, refresh)
	
	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond) 
	s.Prefix = "Working... "
	s.Start()

	task, err := pickSpotlightTask(client, cfg, taskContext, filterArgs)
	if err != nil {
    	fmt.Println("❌", theme.Error(err.Error()))
    	return
	}

	s.Stop()

	displaySpotlight(task, true)
}

func runInteractiveSpotlight(client *taskwarrior.Client, cfg *types.Config, moodFlag string, contextFlag string, refresh bool, filterArgs []string) {
	taskContext := askOrLoadContextFromEnv(moodFlag, contextFlag, refresh)

	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond) 
	s.Prefix = "Working... "
	s.Start()

	task, err := pickSpotlightTask(client, cfg, taskContext, filterArgs)
	if err != nil {
    	fmt.Println("❌", theme.Error(err.Error()))
    	return
	}

	s.Stop()

	displaySpotlight(task, false)
	promptUserAction(client, task)
}

func pickSpotlightTask(client *taskwarrior.Client, cfg *types.Config, taskContext TaskContext, filterArgs []string) (SpotlightResult, error) {

	var tasks []types.Task
	var err error

	if len(filterArgs) > 0 {
		tasks, err = client.GetPendingTasksWithArgsFiltered(cfg, filterArgs)
	} else {
		tasks, err = client.GetPendingTasksFiltered(cfg)
	}

	if err != nil {
		return SpotlightResult{}, fmt.Errorf("fetch tasks: %w", err)
	}

	if len(tasks) == 0 {
		return SpotlightResult{}, fmt.Errorf("no pending tasks: %w", err)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Urgency > tasks[j].Urgency
	})

	if len(tasks) > 20 {
		tasks = tasks[:20]
	}

	llmClient, err := llm.NewClient(&cfg.LLM)
	if err != nil {
		return SpotlightResult{}, fmt.Errorf("init llm client: %w", err)
	}

	prompt := createSpotlightPrompt(taskContext, tasks, cfg)
	messages := []llm.Message{
		{Role: "user", Content: prompt},
	}

	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Request"))
		fmt.Println(theme.Info(messages))
	}

	if !cfg.Settings.EnableLLM {
		return SpotlightResult{}, fmt.Errorf("sending API Request to LLM is disabled via config")
	}

	response, err := llmClient.Chat(messages)
	
	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Response"))
		fmt.Println(theme.Info(response))
	}
	if err != nil {
		return SpotlightResult{}, fmt.Errorf("llm chat error: %w", err)
	}

	var result SpotlightResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return SpotlightResult{}, fmt.Errorf("unmarshall llm response: %w", err)
	}

	return result, nil
}

func displaySpotlight(t SpotlightResult, silent bool) {
	if silent {
		fmt.Printf("🎯 Spotlight: %s (%s) — %s\n", t.Title, t.Estimated, t.Reason)
	} else {
		fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
		fmt.Println(theme.Title("          🎯 SPOTLIGHT TASK:"))
		fmt.Println(theme.Title("───────────────────────────────────────────────"))

		// fmt.Println("\n🎯 Spotlight Task:\n")
		fmt.Printf("%s %s %s\n", "⚡", theme.Info("Task:"), theme.Success(t.Title))
		fmt.Printf("%s %s %s   ", "🕒", theme.Info("Estimated:"), t.Estimated)
		fmt.Printf("%s %s %s\n", "🏷️", theme.Info("Context:"), t.ContextTag)
		fmt.Printf("%s %s %s\n", "🔥", theme.Info("Why now:"), t.Reason)
		if (len(t.History) > 0) {
			fmt.Printf("%s %s %s\n", "📖", theme.Info("History:"), t.History)
		}
		if (len(t.Goal) > 0) {
			fmt.Printf("%s %s %s\n", "⛰️", theme.Info("Goal:"), t.Goal)
		}
		fmt.Printf("%s %s %s\n", "☑️", theme.Info("Next:"), t.Next)
		fmt.Printf("%s %s\n", theme.Warn("🚀 Ready?"), "You can do this now. Hit enter to start.")
	}
}

func askOrLoadContextFromEnv(moodFlag string, contextFlag string, refresh bool) TaskContext {
	if !refresh {
		if taskContext := loadContextFromEnv("", "", false); taskContext.Mood != "neutral" || taskContext.Location != "unknown" {
			return taskContext
		}
	}

	// fmt.Print("Where are you right now? ([h]ome/[o]ffice/[t]ravel/[]other): ")
	fmt.Printf("%s %s %s: ", theme.Title("→"), theme.Info("Current Location?"), "[h]ome [o]ffice [t]ravel other ")
	var location string
	fmt.Scanln(&location)
	switch location {
	case "h", "home":
		location = "home"
	case "o", "office":
		location = "office"
	case "t", "tired":
		location = "travel"
	default:
		if moodFlag != "" {
			location = moodFlag			
			fmt.Printf(theme.Warn("Context set via --context : `%s`\n"), location)
		} else {
			location = "unknown"
			fmt.Printf(theme.Warn("No input. Defaulting to `%s`\n"), location)
		}
	}

	// fmt.Print("How are you feeling? ([e]nergetic/[f]ocused/[t]ired/[s]tressed/[n]eutral): ")
	fmt.Printf("%s %s %s: ", theme.Title("→"), theme.Info("Current energy?"), "[e]nergetic [f]ocused [t]ired [s]tressed [n]eutral ")
	var mood string
	fmt.Scanln(&mood)
	switch mood {
	case "e", "energetic":
		mood = "energetic"
	case "f", "focused":
		mood = "focused"
	case "t", "tired":
		mood = "tired"
	case "s", "stressed":
		mood = "stressed"
	default:
		if contextFlag != "" {
			mood = contextFlag			
			fmt.Printf(theme.Warn("Mood set via --mood `%s`\n"), mood)
		} else {
			mood = "neutral"
			fmt.Printf(theme.Warn("No mood is still a mood. Set to `%s`\n"), mood)
		}

	}

	taskContext := TaskContext{
		Mood:      mood,
		Location:  location,
		Timestamp: time.Now(),
	}

	saveContextToEnv(taskContext)
	return taskContext
}

func saveContextToEnv(taskContext TaskContext) {
	os.Setenv("TASKVANGUARD_MOOD", taskContext.Mood)
	os.Setenv("TASKVANGUARD_LOCATION", taskContext.Location)
	os.Setenv("TASKVANGUARD_TIME", taskContext.Timestamp.Format(time.RFC3339))
}

func loadContextFromEnv(moodFlag, contextFlag string, refresh bool) TaskContext {
	var mood, location string
	var timestamp time.Time

	if moodFlag != "" {
		mood = moodFlag
	} else if !refresh {
		mood = os.Getenv("TASKVANGUARD_MOOD")
	}
	if mood == "" {
		mood = "neutral"
	}

	if contextFlag != "" {
		location = contextFlag
	} else if !refresh {
		location = os.Getenv("TASKVANGUARD_LOCATION")
	}
	if location == "" {
		location = "unknown"
	}

	if !refresh {
		if timeStr := os.Getenv("TASKVANGUARD_TIME"); timeStr != "" {
			if parsedTime, err := time.Parse(time.RFC3339, timeStr); err == nil {
				if time.Since(parsedTime) < 4*time.Hour {
					timestamp = parsedTime
				}
			}
		}
	}
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	return TaskContext{
		Mood:      mood,
		Location:  location,
		Timestamp: timestamp,
	}
}

func createSpotlightPrompt(taskContext TaskContext, tasks []types.Task, cfg *types.Config) string {
	// Create enhanced tasks with goal descriptions for LLM
	type TaskForLLM struct {
		types.Task
		GoalDescription string `json:"goal_description,omitempty"`
	}

	var enhancedTasks []TaskForLLM
	goalsManager := goals.NewManager(cfg)

	for _, task := range tasks {
		enhanced := TaskForLLM{Task: task}
		
		// If task has a goal, get its description
		if task.Goal != "" {
			linkedGoal, err := goalsManager.GetLinkedGoal(strconv.Itoa(task.ID))
			if err == nil && linkedGoal != nil {
				enhanced.GoalDescription = linkedGoal.Description
			}
		}
		
		enhancedTasks = append(enhancedTasks, enhanced)
	}

	tasksJSON, err := json.MarshalIndent(enhancedTasks, "  ", "  ")
	if err != nil {
		tasksJSON = []byte("[]")
	}

	prompt := fmt.Sprintf(`You are a productivity expert helping someone choose the best task to work on right now.

Context:
- Current mood: %s
- Current location: %s  
- Time: %s

Here are the top pending tasks sorted by urgency (JSON format):

%s`, taskContext.Mood, taskContext.Location, time.Now().Format("3:04 PM"), string(tasksJSON))

	prompt += `
Analyze the provided tasks and select the ONE best task for right now based on current context.

Rules:
- Address the user as "you"
- Return only valid JSON
- Use empty strings ("") for optional fields when no data exists
- For "history" field: if task was previously skipped, explain why; otherwise use empty string
- Be direct and concise

Required JSON format:
{
  "task_id": number,
  "title": "string - simplified task title",
  "reason": "string - why this task fits now",
  "estimated": "string - time estimate format: 'X min' or 'X hour'",
  "history": "string - reasons for skipping if any, otherwise empty",
  "goal": "string - associated goal if relevant, otherwise empty", 
  "context_tag": "string - single word context descriptor",
  "next": "string - immediate first action to take"
}`

	return prompt
}

func promptUserAction(client *taskwarrior.Client, spotlightTask SpotlightResult) error {

	fmt.Printf("%s %s %s: ", theme.Title("→"), theme.Info("Do this task now?"), "[Y]es/[n]o/[s]nooze")

	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))

	var isTaskSkipped bool

	switch strings.ToLower(response) {
	case "y", "yes", "":
		fmt.Println(theme.Success("Momentum: Task started!"))
		client.StartTask(strconv.Itoa(spotlightTask.TaskID))
		return nil

	case "s", "snooze":
		isTaskSkipped = true
		fmt.Println(theme.Warn("Task snoozed. It'll be back."))
	case "n", "no":
		isTaskSkipped = true
		fmt.Printf("%s %s %s", theme.Error("Blocked."), theme.Info("What's stopping you?"), "[quick note]: ")
		reader := bufio.NewReader(os.Stdin)
		reason, _ := reader.ReadString('\n')
		reason = strings.TrimSpace(reason)
		if reason != "" {
			client.AddSingleAnnotation(strconv.Itoa(spotlightTask.TaskID), fmt.Sprintf("Skipped: %s", reason))
			fmt.Println(theme.Info("Noted."))
		}
	}

	if isTaskSkipped {
		task, err := client.GetTaskByID(strconv.Itoa(spotlightTask.TaskID))
		if err == nil {
			task.Skipped ++
			arg := "skipped:" + strconv.Itoa(task.Skipped)
			client.ModifyTaskInTaskWarrior(spotlightTask.TaskID, []string{arg})
		} else {
			return err
		}
	}
	
	return nil
}