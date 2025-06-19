package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"xarc.dev/taskvanguard/internal/llm"
	"xarc.dev/taskvanguard/internal/prompts"
	"xarc.dev/taskvanguard/internal/taskwarrior"
	"xarc.dev/taskvanguard/pkg/theme"
	"xarc.dev/taskvanguard/pkg/types"
)

// GuideResponse is returned by the LLM guide session.
type GuideResponse struct {
	Question       string `json:"question"`
	AnswersSummary string `json:"answers-summary"`
	GoalSummary    string `json:"goal-summary"`
	GoalAction     string `json:"goal-action"`
	GoalName       string `json:"goal-name"`
}

// QuestionAnswer represents a single question-answer pair in the guide session.
type QuestionAnswer struct {
	Question string
	Answer   string
}

// RoadmapTask represents a task in the generated roadmap.
type RoadmapTask struct {
	ID            int      `json:"id"`
	Description   string   `json:"description"`
	Project       string   `json:"project"`
	Tags          []string `json:"tags"`
	Depends       []int    `json:"depends"`
	Priority      string   `json:"priority"`
	Estimate      string   `json:"estimate"`
	Resources     []string `json:"resources"`
	Risks         string   `json:"risks"`
	Metrics       string   `json:"metrics"`
	DecisionPoint bool     `json:"decision_point"`
}

// TaskWarriorTask represents a task in TaskWarrior JSON format.
type TaskWarriorTask struct {
	UUID        string                   `json:"uuid"`
	Status      string                   `json:"status"`
	Entry       string                   `json:"entry"`
	Description string                   `json:"description"`
	Modified    string                   `json:"modified"`
	Project     string                   `json:"project,omitempty"`
	Tags        []string                 `json:"tags,omitempty"`
	Priority    string                   `json:"priority,omitempty"`
	Depends     string                   `json:"depends,omitempty"`
	Annotations []TaskWarriorAnnotation  `json:"annotations,omitempty"`
}

// TaskWarriorAnnotation represents an annotation in TaskWarrior format.
type TaskWarriorAnnotation struct {
	Entry       string `json:"entry"`
	Description string `json:"description"`
}

type GuideQuestionData struct {
	QAHistory         string
	QuestionThreshold int
	QuestionCount     int
}

type GuideSummaryData struct {
	QAHistory string
}

type GuideRoadmapData struct {
	GoalSummary    string
	AnswersSummary string
	UserTags       string
}

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Asks questions about a specific goals and creates action plan",
	Long: `Asks one question after another to better understand a specific goal. Then provides a roadmap that contains actionable tasks and subtasks.`,
	Run: runGuide,
}

func init() {
	guideCmd.Flags().Int("questions-count", 6, "Set amount of questions asked before presenting the roadmap")
}

func runGuide(cmd *cobra.Command, args []string) {
	env, err := taskwarrior.Bootstrap(cmd)
	if err != nil {
		fmt.Println(theme.Error(err.Error()))
		return
	}

	questionsCount, _ := cmd.Flags().GetInt("questions-count")
	if questionsCount <= 0 {
		questionsCount = env.Config.Settings.GuidingQuestionAmount
		if questionsCount <= 0 {
			questionsCount = 6
		}
	}
	questionsCount += 2

	goal := promptForGoal(questionsCount)
	if goal == "" {
		fmt.Println(theme.Warn("No goal provided. Exiting."))
		return
	}

	timeframe := promptForTimeframe(questionsCount)
	if timeframe == "" {
		fmt.Println(theme.Warn("No timeframe provided. Exiting."))
		return
	}

	qaHistory := []QuestionAnswer{{
		Question: "What is a specific goal you want to achieve?",
		Answer:   goal,
	}, {
		Question: "What is a realistic timeframe for achieving this goal?",
		Answer:   timeframe,
	}}

	guideResult, err := conductQuestioningSession(env.Config, qaHistory, questionsCount)
	if err != nil {
		fmt.Println(theme.Error(err.Error()))
		return
	}

	if !confirmGoal(guideResult) {
		fmt.Println(theme.Warn("Guide session cancelled."))
		return
	}

	generateRoadmap(env.Config, guideResult)
}

func promptForGoal(totalQuestions int) string {
	fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
	fmt.Println(theme.Title("          🏔️ GOAL GUIDANCE:"))
	fmt.Println(theme.Title("───────────────────────────────────────────────"))
	fmt.Printf("%s %s %s: ", theme.Title("→"), theme.Info(fmt.Sprintf("[1/%d] Question: ", totalQuestions)), theme.Info("What is a specific goal you want to achieve?"))
	
	reader := bufio.NewReader(os.Stdin)
	goal, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(goal)
}

func promptForTimeframe(totalQuestions int) string {
	fmt.Printf("%s %s %s: ", theme.Title("→"), theme.Info(fmt.Sprintf("[2/%d] Question: ", totalQuestions)), theme.Info("What is a realistic timeframe for achieving this goal?"))
	
	reader := bufio.NewReader(os.Stdin)
	timeframe, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(timeframe)
}

func conductQuestioningSession(cfg *types.Config, qaHistory []QuestionAnswer, maxQuestions int) (*GuideResponse, error) {
	llmClient, err := llm.NewClient(&cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("init llm client: %w", err)
	}

	questionCount := 1
	
	for questionCount < maxQuestions {
		currentQuestionCount := len(qaHistory)
		prompt, err := createQuestionPrompt(qaHistory, maxQuestions, currentQuestionCount)
		if err != nil {
			return nil, fmt.Errorf("create prompt: %w", err)
		}

		if cfg.Settings.Debug {
			fmt.Println(theme.Title("LLM Question Prompt"))
			fmt.Println(theme.Info(prompt))
		}

		s := spinner.New(spinner.CharSets[40], 100*time.Millisecond)
		s.Prefix = "Thinking... "
		s.Start()

		messages := []llm.Message{{
			Role:    "user",
			Content: prompt,
		}}

		if !cfg.Settings.EnableLLM {
			s.Stop()
			return nil, fmt.Errorf("sending API Request to LLM is disabled via config")
		}

		response, err := llmClient.Chat(messages)
		s.Stop()
		
		if err != nil {
			return nil, fmt.Errorf("llm chat error: %w", err)
		}

		if cfg.Settings.Debug {
			fmt.Println(theme.Title("LLM Question Response"))
			fmt.Println(theme.Info(response))
		}

		var questionResp struct {
			Question string `json:"question"`
		}
		if err := json.Unmarshal([]byte(response), &questionResp); err != nil {
			return nil, fmt.Errorf("unmarshall llm response: %w", err)
		}

		if questionResp.Question == "" {
			break
		}

		questionCount++
		fmt.Printf("\n%s %s %s: ", theme.Title("→"), theme.Info(fmt.Sprintf("[%d/%d] Question: ", questionCount+1, maxQuestions)), questionResp.Question)
		
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read answer: %w", err)
		}
		answer = strings.TrimSpace(answer)

		qaHistory = append(qaHistory, QuestionAnswer{
			Question: questionResp.Question,
			Answer:   answer,
		})
	}

	// Now use the summary template
	prompt, err := createSummaryPrompt(qaHistory)
	if err != nil {
		return nil, fmt.Errorf("create summary prompt: %w", err)
	}

	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Summary Prompt"))
		fmt.Println(theme.Info(prompt))
	}

	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond)
	s.Prefix = "Summarizing... "
	s.Start()

	messages := []llm.Message{{
		Role:    "user",
		Content: prompt,
	}}

	response, err := llmClient.Chat(messages)
	s.Stop()
	
	if err != nil {
		return nil, fmt.Errorf("summary llm chat error: %w", err)
	}

	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Summary Response"))
		fmt.Println(theme.Info(response))
	}

	var finalResp GuideResponse
	if err := json.Unmarshal([]byte(response), &finalResp); err != nil {
		return nil, fmt.Errorf("unmarshall summary llm response: %w", err)
	}

	return &finalResp, nil
}

func createQuestionPrompt(qaHistory []QuestionAnswer, maxQuestions int, currentQuestionCount int) (string, error) {
	var builder strings.Builder
	for i, qa := range qaHistory {
		builder.WriteString(fmt.Sprintf("Q%d: %s\nA%d: %s\n\n", i+1, qa.Question, i+1, qa.Answer))
	}
	qaHistoryStr := builder.String()

	questionData := GuideQuestionData{
		QAHistory:         qaHistoryStr,
		QuestionThreshold: maxQuestions,
		QuestionCount:     currentQuestionCount,
	}

	template, err := prompts.LoadPrompt("guide_questions.md")
	if err != nil {
		return "", err
	}

	prompt := strings.ReplaceAll(template, "{{ .QAHistory }}", questionData.QAHistory)
	prompt = strings.ReplaceAll(prompt, "{{ .QuestionThreshold }}", fmt.Sprintf("%d", questionData.QuestionThreshold))
	prompt = strings.ReplaceAll(prompt, "{{ .QuestionCount }}", fmt.Sprintf("%d", questionData.QuestionCount))

	return prompt, nil
}

func createSummaryPrompt(qaHistory []QuestionAnswer) (string, error) {
	var builder strings.Builder
	for i, qa := range qaHistory {
		builder.WriteString(fmt.Sprintf("Q%d: %s\nA%d: %s\n\n", i+1, qa.Question, i+1, qa.Answer))
	}
	qaHistoryStr := builder.String()

	summaryData := GuideSummaryData{
		QAHistory: qaHistoryStr,
	}

	template, err := prompts.LoadPrompt("guide_summary.md")
	if err != nil {
		return "", err
	}

	prompt := strings.ReplaceAll(template, "{{ .QAHistory }}", summaryData.QAHistory)

	return prompt, nil
}

func confirmGoal(guideResult *GuideResponse) bool {
	fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
	fmt.Println(theme.Title("          📑 GOAL SUMMARY:"))
	fmt.Println(theme.Title("───────────────────────────────────────────────"))
	
	fmt.Printf("%s %s\n\n", theme.Info("Goal:"), guideResult.GoalSummary)
	fmt.Printf("%s\n%s\n\n", theme.Info("Key Details:"), guideResult.AnswersSummary)
	
	fmt.Printf("%s %s %s: ", theme.Title("→"), theme.Info("Is this accurate?"), "[Y]es/[n]o/[a]dd more info")
	
	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))
	
	switch response {
	case "y", "yes", "":
		return true
	case "a", "add":
		fmt.Printf("%s %s: ", theme.Info("Additional info"), "")
		reader := bufio.NewReader(os.Stdin)
		additionalInfo, _ := reader.ReadString('\n')
		additionalInfo = strings.TrimSpace(additionalInfo)
		if additionalInfo != "" {
			guideResult.AnswersSummary += "\n- " + additionalInfo
		}
		return true
	default:
		return false
	}
}

func createRoadmapPrompt(cfg *types.Config, guideResult *GuideResponse) (string, error) {
	// get user tags from config
	var userTags []string
	for tagName := range cfg.Tags {
		userTags = append(userTags, tagName)
	}
	userTagsStr := strings.Join(userTags, ", ")
	if userTagsStr == "" {
		userTagsStr = "key, sb, fast, cut (use appropriate tags based on task characteristics)"
	}

	template, err := prompts.LoadPrompt("guide_roadmap.md")
	if err != nil {
		return "", err
	}

	prompt := strings.ReplaceAll(template, "{{ .GoalSummary }}", guideResult.GoalSummary)
	prompt = strings.ReplaceAll(prompt, "{{ .AnswersSummary }}", guideResult.AnswersSummary)
	prompt = strings.ReplaceAll(prompt, "{{ .UserTags }}", userTagsStr)

	return prompt, nil
}

// generateRoadmap creates a roadmap from the guide result and displays it.
func generateRoadmap(cfg *types.Config, guideResult *GuideResponse) {
	fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
	fmt.Println(theme.Title("          🗺️  ROADMAP GENERATION:"))
	fmt.Println(theme.Title("───────────────────────────────────────────────"))
	
	llmClient, err := llm.NewClient(&cfg.LLM)
	if err != nil {
		fmt.Printf("%s %s\n", theme.Error("❌ Failed to initialize LLM client:"), err.Error())
		return
	}

	prompt, err := createRoadmapPrompt(cfg, guideResult)
	if err != nil {
		fmt.Printf("%s %s\n", theme.Error("❌ Failed to create roadmap prompt:"), err.Error())
		return
	}

	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond)
	s.Prefix = "Creating roadmap... "
	s.Start()

	messages := []llm.Message{{
		Role:    "user",
		Content: prompt,
	}}

	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Roadmap Prompt"))
		fmt.Println(theme.Info(prompt))
	}

	if !cfg.Settings.EnableLLM {
		s.Stop()
		fmt.Printf("%s %s\n", theme.Error("❌ LLM disabled:"), "Enable LLM in config to generate roadmap")
		return
	}

	response, err := llmClient.Chat(messages)
	s.Stop()
	
	if err != nil {
		fmt.Printf("%s %s\n", theme.Error("❌ LLM error:"), err.Error())
		return
	}

	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Roadmap Response"))
		fmt.Println(theme.Info(response))
	}

	var roadmapTasks []RoadmapTask
	if err := json.Unmarshal([]byte(response), &roadmapTasks); err != nil {
		fmt.Printf("%s %s\n", theme.Error("❌ Failed to parse roadmap:"), err.Error())
		fmt.Printf("%s\n%s\n", theme.Warn("Raw response:"), response)
		return
	}

	displayRoadmap(roadmapTasks, guideResult)
	
	if err := generateRoadmapMarkdown(roadmapTasks, guideResult); err != nil {
		fmt.Printf("%s %s\n", theme.Warn("⚠️  Failed to generate markdown file:"), err.Error())
	} else {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		var filename string
		if guideResult.GoalName != "" {
			filename = fmt.Sprintf("roadmap_%s_%s.md", guideResult.GoalName, timestamp)
		} else {
			filename = fmt.Sprintf("roadmap_%s.md", timestamp)
		}
		fmt.Printf("%s %s\n", theme.Success("📄 Roadmap saved to:"), filename)
	}
	
	if promptForTaskImport() {
		if err := importTasksToTaskWarrior(roadmapTasks); err != nil {
			fmt.Printf("%s %s\n", theme.Error("❌ Failed to import tasks:"), err.Error())
		} else {
			fmt.Printf("%s %s\n", theme.Success("✅ Tasks imported successfully!"), "")
		}
	}
}

func displayRoadmap(tasks []RoadmapTask, guideResult *GuideResponse) {
	fmt.Printf("%s %s\n\n", theme.Success("✅ Roadmap Generated!"), "")
	fmt.Printf("%s %s\n", theme.Info("🎯 Goal:"), guideResult.GoalSummary)
	fmt.Printf("%s %d tasks\n\n", theme.Info("📋 Plan:"), len(tasks))

	for i, task := range tasks {
		fmt.Printf("%s %s\n", theme.Title(fmt.Sprintf("%d.", i+1)), theme.Success(task.Description))
		
		if task.Project != "" {
			fmt.Printf("   %s %s\n", theme.Info("📁 Project:"), task.Project)
		}
		
		if len(task.Tags) > 0 {
			fmt.Printf("   %s %s\n", theme.Info("🏷️  Tags:"), strings.Join(task.Tags, ", "))
		}
		
		if task.Priority != "" {
			priorityColor := theme.Info
			if task.Priority == "High" {
				priorityColor = theme.Error
			} else if task.Priority == "Medium" {
				priorityColor = theme.Warn
			}
			fmt.Printf("   %s %s\n", theme.Info("⚡ Priority:"), priorityColor(task.Priority))
		}
		
		if task.Estimate != "" {
			fmt.Printf("   %s %s\n", theme.Info("⏱️  Estimate:"), task.Estimate)
		}
		
		if len(task.Depends) > 0 {
			var dependsStr []string
			for _, dep := range task.Depends {
				dependsStr = append(dependsStr, fmt.Sprintf("#%d", dep))
			}
			fmt.Printf("   %s %s\n", theme.Info("🔗 Depends:"), strings.Join(dependsStr, ", "))
		}
		
		if len(task.Resources) > 0 {
			fmt.Printf("   %s %s\n", theme.Info("🛠️  Resources:"), strings.Join(task.Resources, ", "))
		}
		
		if task.Risks != "" {
			fmt.Printf("   %s %s\n", theme.Warn("⚠️  Risks:"), task.Risks)
		}
		
		if task.Metrics != "" {
			fmt.Printf("   %s %s\n", theme.Info("📊 Success:"), task.Metrics)
		}
		
		if task.DecisionPoint {
			fmt.Printf("   %s %s\n", theme.Error("🔄 Decision Point:"), "Review and adapt here")
		}
		
		fmt.Println()
	}

	fmt.Printf("%s %s\n", theme.Warn("💡 Next Steps:"), "Tasks are ready for import into TaskWarrior")
	fmt.Printf("%s %s\n", theme.Info("🚀 Pro Tip:"), "Use `vanguard analyze` ")
}

func generateRoadmapMarkdown(tasks []RoadmapTask, guideResult *GuideResponse) error {
	markdown := formatRoadmapMarkdown(tasks, guideResult)
	
	// Create filename with goal-name and timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	var filename string
	if guideResult.GoalName != "" {
		filename = fmt.Sprintf("roadmap_%s_%s.md", guideResult.GoalName, timestamp)
	} else {
		filename = fmt.Sprintf("roadmap_%s.md", timestamp)
	}
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create markdown file: %w", err)
	}
	defer file.Close()
	
	_, err = file.WriteString(markdown)
	if err != nil {
		return fmt.Errorf("failed to write markdown content: %w", err)
	}
	
	return nil
}

func formatRoadmapMarkdown(tasks []RoadmapTask, guideResult *GuideResponse) string {
	var md strings.Builder
	
	// Header
	md.WriteString("# Goal Roadmap\n\n")
	
	// Goal summary
	md.WriteString("## 🏔️ Goal\n\n")
	md.WriteString(guideResult.GoalSummary + "\n\n")
	
	// Key details
	if guideResult.AnswersSummary != "" {
		md.WriteString("## 📋 Key Details\n\n")
		md.WriteString(guideResult.AnswersSummary + "\n\n")
	}
	
	// Tasks
	md.WriteString(fmt.Sprintf("## 🗺️ Action Plan (%d tasks)\n\n", len(tasks)))
	
	for i, task := range tasks {
		md.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, task.Description))
		
		// Task details in a table format
		md.WriteString("| Attribute | Value |\n")
		md.WriteString("|-----------|-------|\n")
		
		if task.Project != "" {
			md.WriteString(fmt.Sprintf("| 📁 Project | %s |\n", task.Project))
		}
		
		if len(task.Tags) > 0 {
			md.WriteString(fmt.Sprintf("| 🏷️ Tags | %s |\n", strings.Join(task.Tags, ", ")))
		}
		
		if task.Priority != "" {
			md.WriteString(fmt.Sprintf("| ⚡ Priority | %s |\n", task.Priority))
		}
		
		if task.Estimate != "" {
			md.WriteString(fmt.Sprintf("| ⏱️ Estimate | %s |\n", task.Estimate))
		}
		
		if len(task.Depends) > 0 {
			var dependsStr []string
			for _, dep := range task.Depends {
				dependsStr = append(dependsStr, fmt.Sprintf("#%d", dep))
			}
			md.WriteString(fmt.Sprintf("| 🔗 Depends | %s |\n", strings.Join(dependsStr, ", ")))
		}
		
		if len(task.Resources) > 0 {
			md.WriteString(fmt.Sprintf("| 🛠️ Resources | %s |\n", strings.Join(task.Resources, ", ")))
		}
		
		if task.Risks != "" {
			md.WriteString(fmt.Sprintf("| 🚧 Risks | %s |\n", task.Risks))
		}
		
		if task.Metrics != "" {
			md.WriteString(fmt.Sprintf("| 📈 Success Metrics | %s |\n", task.Metrics))
		}
		
		if task.DecisionPoint {
			md.WriteString("| 🔄 Decision Point | Review and adapt here |\n")
		}
		
		md.WriteString("\n")
	}
	
	// Footer
	md.WriteString("---\n\n")
	md.WriteString("**💡 Next Steps:** Tasks are ready for import into TaskWarrior\n\n")
	md.WriteString("**🚀 Pro Tip:** Start with the highest priority tasks first\n\n")
	md.WriteString(fmt.Sprintf("*Generated on %s*\n", time.Now().Format("2006-01-02 15:04:05")))
	
	return md.String()
}

func promptForTaskImport() bool {
	fmt.Printf("\n%s %s %s: ", theme.Title("→"), theme.Info("Import tasks into TaskWarrior?"), "[Y]es/[n]o")
	
	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))
	
	return response == "y" || response == "yes" || response == ""
}

func convertToTaskWarriorFormat(roadmapTasks []RoadmapTask) ([]TaskWarriorTask, map[string]string, error) {
	var twTasks []TaskWarriorTask
	idToUUID := make(map[string]string)
	now := time.Now().UTC().Format("20060102T150405Z")
	
	// First pass: create UUIDs and basic tasks
	for _, task := range roadmapTasks {
		taskUUID := uuid.New().String()
		idToUUID[fmt.Sprintf("%d", task.ID)] = taskUUID
		
		// Convert priority format according to TaskWarrior spec (H/M/L or empty)
		priority := ""
		switch strings.ToLower(task.Priority) {
		case "high":
			priority = "H"
		case "medium":
			priority = "M"
		case "low":
			priority = "L"
		}
		
		// Create annotations for custom fields
		var annotations []TaskWarriorAnnotation
		
		if task.Estimate != "" {
			annotations = append(annotations, TaskWarriorAnnotation{
				Entry:       now,
				Description: "Estimate: " + task.Estimate,
			})
		}
		
		if len(task.Resources) > 0 {
			annotations = append(annotations, TaskWarriorAnnotation{
				Entry:       now,
				Description: "Resources: " + strings.Join(task.Resources, ", "),
			})
		}
		
		if task.Risks != "" {
			annotations = append(annotations, TaskWarriorAnnotation{
				Entry:       now,
				Description: "Risks: " + task.Risks,
			})
		}
		
		if task.Metrics != "" {
			annotations = append(annotations, TaskWarriorAnnotation{
				Entry:       now,
				Description: "Success: " + task.Metrics,
			})
		}
		
		if task.DecisionPoint {
			annotations = append(annotations, TaskWarriorAnnotation{
				Entry:       now,
				Description: "Decision Point: Review and adapt here",
			})
		}
		
		twTask := TaskWarriorTask{
			UUID:        taskUUID,
			Status:      "pending",
			Entry:       now,
			Modified:    now,
			Description: task.Description,
			Project:     task.Project,
			Tags:        task.Tags,
			Priority:    priority,
			Annotations: annotations,
		}
		
		twTasks = append(twTasks, twTask)
	}
	
	// Second pass: handle dependencies
	for i, task := range roadmapTasks {
		if len(task.Depends) > 0 {
			var dependUUIDs []string
			for _, depID := range task.Depends {
				if depUUID, exists := idToUUID[fmt.Sprintf("%d", depID)]; exists {
					dependUUIDs = append(dependUUIDs, depUUID)
				}
			}
			if len(dependUUIDs) > 0 {
				// TaskWarrior spec: comma-separated UUIDs without spaces
				twTasks[i].Depends = strings.Join(dependUUIDs, ",")
			}
		}
	}
	
	return twTasks, idToUUID, nil
}

func importTasksToTaskWarrior(roadmapTasks []RoadmapTask) error {
	twTasks, _, err := convertToTaskWarriorFormat(roadmapTasks)
	if err != nil {
		return fmt.Errorf("failed to convert tasks: %w", err)
	}
	
	// Create temporary file
	tempDir, err := os.MkdirTemp("", "taskvanguard-import-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	tempFile := filepath.Join(tempDir, "tasks.json")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()
	
	// Write each task as a separate JSON line according to TaskWarrior spec
	for _, task := range twTasks {
		taskJSON, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal task: %w", err)
		}
		
		// TaskWarrior import format: one JSON object per line, no pretty printing
		if _, err := file.WriteString(string(taskJSON) + "\n"); err != nil {
			return fmt.Errorf("failed to write task to file: %w", err)
		}
	}
	
	file.Close()
	
	// Import using task import command with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	s := spinner.New(spinner.CharSets[40], 100*time.Millisecond)
	s.Prefix = "Importing tasks... "
	s.Start()
	defer s.Stop()
	
	cmd := exec.CommandContext(ctx, "task", "import", tempFile)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("task import failed: %w\nOutput: %s", err, string(output))
	}
	
	return nil
}
