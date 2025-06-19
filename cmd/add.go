package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"xarc.dev/taskvanguard/internal/analyzer"
	"xarc.dev/taskvanguard/internal/taskwarrior"
	"xarc.dev/taskvanguard/pkg/theme"
	"xarc.dev/taskvanguard/pkg/types"
)

var addCmd = &cobra.Command{
	Use:   "add [task arguments...]",
	Short: "Analyze and add a task with AI suggestions",
	Long: `Analyzes the task arguments using AI to suggest tags and categorization,
then offers to add the task to TaskWarrior with the suggested enhancements.`,
	Run: runSmartAdd,
}

func init() {
	addCmd.Flags().Bool("no-subtask", false, "Disable subtask splitting for this command")
	addCmd.Flags().Bool("no-tags", false, "Disable adding tags intelligently")
	addCmd.Flags().Bool("no-annotations", false, "Disable adding annotations")
}

type Option struct {
	key    string
	prompt string
}

func runSmartAdd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide task arguments (same ones as with 'task add' command)")
		return
	}
	
	env, err := taskwarrior.Bootstrap(cmd)
	if err != nil {
		fmt.Println(theme.Error(err.Error()))
		return
	}

	output, newTaskId, err := env.Client.AddTaskToTaskWarrior(args)
	fmt.Print(string(output))
	if err != nil {
		fmt.Printf("Could not create the original Task you provided: %v\n", err)
		return
	}

	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond) 
	s.Prefix = "Working... "
	s.Start()

	taskArgs := strings.Join(args, " ")
	suggestion, err := analyzer.AnalyzeSingleTaskWithLLM(env.Config, taskArgs, env.UserGoals, env.UserProjects)
	if err != nil {
		fmt.Printf("Error analyzing task: %v\n", err)
		return
	}

	s.Stop()

	if env.Config.Settings.EnableLowercase {
		lowercaseTaskSuggestion(suggestion)
	}

	displaySuggestions(env.Config, taskArgs, suggestion)
	userConfirmations := askUserConfirmation(env.Config, suggestion)

	if !anyAccepted(userConfirmations) {
		fmt.Println(theme.Success("\nAdded only provided Task without modifications."))
		return
	}

	if userConfirmations["subtasks"] && len(suggestion.Subtasks) > 0 {
		if err := addSubtasks(*env.Client, suggestion.Subtasks); err != nil {
			fmt.Printf("Error adding subtasks: %v\n", err)
		}
		return
	}

	enhancedArgs := buildEnhancedTaskArgs(suggestion, userConfirmations)
	if err := addAnnotationsInTaskWarrior(env.Config, newTaskId, suggestion.AdditionalInfo); err != nil {
		fmt.Printf("Error adding annotations to task id: %d \nError: %v\n", newTaskId, err)
		return
	}

	if _, err := env.Client.ModifyTaskInTaskWarrior(newTaskId, enhancedArgs); err != nil {
		fmt.Printf("Error adding task: %v\n", err)
		return
	}

	fmt.Printf("\n%s Task %d modified successfully!\n", theme.Success("✓"), newTaskId)
}

func displaySuggestions(cfg *types.Config, taskArgs string, suggestion *types.TaskSuggestion) {
	fmt.Println(theme.Title("\nOriginal Task:"))
	fmt.Println(theme.Info(taskArgs))

	// 1. Show suggestion in colorful, readable way
	fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
	fmt.Println(theme.Title("        REFINED TASK SUGGESTION"))
	fmt.Println(theme.Title("───────────────────────────────────────────────\n"))

	fmt.Printf("%s %s\n", theme.Info("Task:"), theme.Success(suggestion.RefinedTask))
	
	
	if cfg.Settings.SplitTasks {
		if len(suggestion.Subtasks) > 0 {
			fmt.Println(theme.Info("\nSubtasks:"))
			for _, st := range suggestion.Subtasks {
				fmt.Printf("  %s %s\n", theme.Success("▸"), theme.Success(st))
			}
		}
	}
	
	if len(suggestion.SuggestedTags) > 0 {
		fmt.Printf("\n%s %s\n", theme.Info("Suggested Tags:"), theme.Warn(strings.Join(suggestion.SuggestedTags, " ")))
	}
	if suggestion.Project != "" {
		fmt.Printf("%s %s\n", theme.Info("Suggested Project:"), theme.Success(suggestion.Project))
	}

	if len(suggestion.AdditionalInfo) > 0 {
		fmt.Println(theme.Info("\nAnnotations:"))
		for key, val := range suggestion.AdditionalInfo {
			symbol, label := "▸", key
			if meta, ok := cfg.Annotations[key]; ok {
				if meta.Symbol != "" {
					symbol = meta.Symbol
				}
				if meta.Label != "" {
					label = meta.Label
				}
			}
			fmt.Printf("%s %s: %s\n", symbol, label, val)
		}
	}

	fmt.Println(theme.Title("\n────────────────────────────────────────────"))
	
}

func askUserConfirmation(cfg *types.Config, suggestion *types.TaskSuggestion) map[string]bool {

	// 2. Prompt loop
	reader := bufio.NewReader(os.Stdin)
	applyAll := false
	denyAll := false

	options := []Option{}
	promptSuffix := " [y]es/[N]o/[a]ll/[q]uit: ";
	splitPrompt := "split into subtasks?"

	if cfg.Settings.SplitTasks && len(suggestion.Subtasks) > 0 {
		options = append(options, Option{
			key:    "subtasks",
			prompt: theme.Info(splitPrompt),
		})
	}

	options = append(options,
		Option{"title",       theme.Info("apply refined title?")},
		Option{"tags",        theme.Info("apply suggested tags?")},
		Option{"project",     theme.Info("apply suggested project?")},
		Option{"annotations", theme.Info("apply annotations?")},
	)

	applied := make(map[string]bool)

	for _, opt := range options {
		if applyAll {
			applied[opt.key] = true
			continue
		}

		if denyAll {
			applied[opt.key] = false
			continue
		}

		// TODO: Show what is done with each prompt
		fullPrompt := fmt.Sprintf("%s%s", opt.prompt, promptSuffix)
		fmt.Print(fullPrompt)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "y", "yes":
			applied[opt.key] = true
		case "a", "all":
			applied[opt.key] = true
			applyAll = true
		case "q", "quit":
			denyAll = true
		default:
			applied[opt.key] = false
		}
	}

	fmt.Println(theme.Title("\nYour choices:"))
	for _, opt := range options {
		status := theme.Error("No")
		if applied[opt.key] {
			status = theme.Success("Yes")
		}
		fmt.Printf("  %s: %s\n", theme.Info(opt.key), status)
	}
	return applied

}

// func buildEnhancedTaskArgs(originalArgs []string, suggestion *types.TaskSuggestion, userConfirmations map[string]bool) []string {
func buildEnhancedTaskArgs(suggestion *types.TaskSuggestion, userConfirmations map[string]bool) []string {
	var args []string
	
	// Use refined title if accepted, otherwise dont change title
	if userConfirmations["title"] && suggestion.RefinedTask != "" {
		args = append(args, suggestion.RefinedTask)
	}

	// Add suggested tags if accepted
	if userConfirmations["tags"] && len(suggestion.SuggestedTags) > 0 {
		args = append(args, suggestion.SuggestedTags...)
	}
	
	// Add project if accepted
	if userConfirmations["project"] && suggestion.Project != "" {
		args = append(args, "project:"+suggestion.Project)
	}
	
	return args
}

func addAnnotationsInTaskWarrior(cfg *types.Config, taskId int, additionalInfo map[string]string) error {
	taskIdStr := strconv.Itoa(taskId)

	for key, info := range additionalInfo {
		if info == "" {
			continue
		}
		// Lookup symbol and label
		symbol := ""
		label := key
		if meta, ok := cfg.Annotations[key]; ok {
			if meta.Symbol != "" {
				symbol = meta.Symbol + " "
			}
			if meta.Label != "" {
				label = meta.Label
			}
		}
		annotationText := fmt.Sprintf("%s%s: %s", symbol, label, info)
		if err := addSingleAnnotation(taskIdStr, annotationText); err != nil {
			return err
		}
	}

	return nil
}

func addSingleAnnotation(taskId string, value string) error {
	cmd := exec.Command("task", taskId, "annotate", value)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add annotation to task %s: %v\nOutput: %s", taskId, err, string(output))
	}
	
	return nil
}

func anyAccepted(confirmations map[string]bool) bool {
	for _, v := range confirmations {
		if v {
			return true
		}
	}
	return false
}

func addSubtasks(client taskwarrior.Client, subtasks []string) error {
	for _, subtask := range subtasks {
		if subtask == "" {
			continue
		}
		output, _, err := client.AddTaskToTaskWarrior([]string{subtask})
		if err != nil {
			return err
		}
		fmt.Print(output)
	}
	return nil
}

func lowercaseTaskSuggestion(t *types.TaskSuggestion) {
    t.RefinedTask = strings.ToLower(t.RefinedTask)
    t.GoalAlignment = strings.ToLower(t.GoalAlignment)
    t.Project = strings.ToLower(t.Project)
    for i, s := range t.Subtasks {
        t.Subtasks[i] = strings.ToLower(s)
    }
    for k, v := range t.AdditionalInfo {
        t.AdditionalInfo[k] = strings.ToLower(v)
    }
    for i, tag := range t.SuggestedTags {
        t.SuggestedTags[i] = strings.ToLower(tag)
    }
}
