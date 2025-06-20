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
	"xarc.dev/taskvanguard/pkg/utils"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze tasks for categorization and priority insights",
	Long: `Analyze your TaskWarrior tasks to get AI-powered insights about
categorization, priority adjustments, and potential task relationships.`,
	Run: func(cmd *cobra.Command, args []string) {

		s := spinner.New(spinner.CharSets[40], 100*time.Millisecond) 
		s.Prefix = "→ Analyzing your task list... "
		s.Start()
	
		env, err := taskwarrior.Bootstrap(cmd)
		if err != nil {
			s.Stop()
			fmt.Println(theme.Error("Bootstrap failed: " + err.Error()))
			return
		}

		var taskList []types.Task
		// If analyze is used without arguments, ask user about analyzing all tasks
		if len(args) == 0 {
			s.Stop() // Stop spinner before user interaction
			
			// Ask if user wants to analyze all tasks
			if !promptAnalyzeAllTasks(env.Config.Settings.TaskImportLimit) {
				// User declined, show examples and exit
				showFilterExamples()
				return
			}
			
			// User accepted, restart spinner and get all tasks
			s.Start()
			taskList, err = env.Client.GetPendingTasks()
			if err != nil {
				s.Stop()
				fmt.Println(theme.Error("Failed to get tasks: " + err.Error()))
				return
			}
		} else {
			// Use arguments as filter
			taskList, err = env.Client.GetPendingTasksWithArgsFiltered(env.Config, args)
			if err != nil {
				s.Stop()
				fmt.Println(theme.Error("Failed to get filtered tasks: " + err.Error()))
				return
			}
		}

		// Apply task limit
		limit := env.Config.Settings.TaskImportLimit
		if len(taskList) > limit {
			taskList = taskList[:limit]
			fmt.Printf("Limited to first %d tasks for analysis\n", limit)
		}

		if len(taskList) == 0 {
			s.Stop()
			fmt.Println(theme.Warn("No tasks found matching criteria"))
			return
		}

		// Count and display task count
		fmt.Printf("\n→ Found %d tasks for analysis!\n", len(taskList))

		// Convert tasks to task args format for batch processing
		taskArgs := make([]string, len(taskList))
		for i, task := range taskList {
			var sb strings.Builder
			sb.WriteString(task.Description)
			if task.Project != "" {
				sb.WriteString(" project:")
				sb.WriteString(task.Project)
			}

			for _, tag := range task.Tags {
				sb.WriteString(" +")
				sb.WriteString(tag)
			}

			if task.Priority != "" {
				sb.WriteString(" priority:")
				sb.WriteString(task.Priority)
			}

			taskArgs[i] = sb.String()
		}

		// Analyze batch
		suggestions, err := analyzer.AnalyzeBatchTasksWithLLM(
			env.Config, 
			taskArgs, 
			env.UserGoals, 
			env.UserProjects,
		)
		s.Stop()
		if err != nil {
			fmt.Println(theme.Error("Analysis failed: " + err.Error()))
			return
		}

		if env.Config.Settings.EnableLowercase {
			lowercaseTaskBatchSuggestion(suggestions)
		}

		// === User Prompt: Edit Mode Selection ===
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("How do you want to proceed? [o]ne-by-one / [e]dit all: ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		var oneByOneMode, massEditMode bool
		switch input {
		case "o", "":
			oneByOneMode = true
		case "e":
			massEditMode = true
		default:
			fmt.Println("Invalid mode selected.")
			return
		}

		if oneByOneMode {
			err := oneByOneInteractiveApply(*env.Client, reader, taskList, suggestions)
			if err != nil {
				fmt.Println(theme.Error("Failed to apply suggestions: " + err.Error()))
			}
			return
		}

		if massEditMode {
			err := massEditViaEditor(*env.Client, taskList, suggestions)
			if err != nil {
				fmt.Println(theme.Error("Mass edit failed: " + err.Error()))
			}
			return
		}

	},
}



func oneByOneInteractiveApply(
	client taskwarrior.Client,
	reader *bufio.Reader,
	taskList []types.Task,
	suggestions *types.BatchTaskSuggestion,
) error {
	var acceptAll, denyAll bool

	fmt.Println(theme.Title("\n=== Task Analysis Results ==="))
	userChoices := make([]bool, len(suggestions.TaskAnalyses))

	for i, suggestion := range suggestions.TaskAnalyses {
		if acceptAll {
			userChoices[i] = true
			continue
		}
		if denyAll {
			userChoices[i] = false
			continue
		}

		orig := taskList[i]
		fmt.Printf("\n%s Task %d: %s%s\n", theme.Info("["), suggestion.TaskIndex, orig.Description, theme.Info("]"))

		if suggestion.RefinedTask != orig.Description {
			fmt.Printf("  %s %s\n", theme.Title("Refined:"), suggestion.RefinedTask)
		}
		if suggestion.Project != "" {
			fmt.Printf("  %s %s\n", theme.Title("Project:"), suggestion.Project)
		}
		if len(suggestion.SuggestedTags) > 0 {
			fmt.Printf("  %s %v\n", theme.Title("Tags:"), suggestion.SuggestedTags)
		}
		if suggestion.GoalAlignment != "" {
			fmt.Printf("  %s %s\n", theme.Title("Goal Alignment:"), suggestion.GoalAlignment)
		}
		if len(suggestion.Subtasks) > 0 {
			fmt.Println(theme.Title("Subtasks:"))
			for _, subtask := range suggestion.Subtasks {
				fmt.Printf("    - %s\n", subtask)
			}
		}

		fmt.Print("Accept task suggestion? [y]es/[N]o/[a]ll/[q]uit: ")
		inputRaw, _ := reader.ReadString('\n')
		input := strings.ToLower(strings.TrimSpace(inputRaw))

		switch input {
		case "y", "yes":
			userChoices[i] = true
		case "a", "all":
			userChoices[i] = true
			acceptAll = true
		case "q", "quit":
			denyAll = true
		default:
			userChoices[i] = false
		}
	}

	// === Apply Suggestions OneByOne ===
	fmt.Println(theme.Title("\n=== Applying Accepted Suggestions ==="))
	for i, apply := range userChoices {
		if !apply {
			continue
		}

		orig := taskList[i]
		suggestion := suggestions.TaskAnalyses[i]
		args := utils.TaskSuggestionToArgs(suggestion)
		_, err := client.ModifyTaskInTaskWarrior(orig.ID, args)

		if err != nil {
			fmt.Println(theme.Error(fmt.Sprintf("Task %d: %s", i, err.Error())))
		} else {
			fmt.Println(theme.Success(fmt.Sprintf("Task %d applied: %s", i, orig.Description)))
		}
	}

	return nil
}

func lowercaseTaskBatchSuggestion(batch *types.BatchTaskSuggestion) {
	for i := range batch.TaskAnalyses {
		analysis := &batch.TaskAnalyses[i]
		analysis.RefinedTask = strings.ToLower(analysis.RefinedTask)
		analysis.GoalAlignment = strings.ToLower(analysis.GoalAlignment)
		analysis.Project = strings.ToLower(analysis.Project)
		
		for j, s := range analysis.Subtasks {
			analysis.Subtasks[j] = strings.ToLower(s)
		}
		
		for k, v := range analysis.AdditionalInfo {
			analysis.AdditionalInfo[k] = strings.ToLower(v)
		}
		
		for j, tag := range analysis.SuggestedTags {
			analysis.SuggestedTags[j] = strings.ToLower(tag)
		}
	}
}

func massEditViaEditor(client taskwarrior.Client, taskList []types.Task, suggestions *types.BatchTaskSuggestion) error {
	var commands []string
	
	// Generate task modify commands
	for i, suggestion := range suggestions.TaskAnalyses {
		orig := taskList[i]
		args := utils.TaskSuggestionToArgs(suggestion)
		
		if len(args) > 0 {
			cmdParts := []string{"task", "modify", strconv.Itoa(orig.ID)}
			cmdParts = append(cmdParts, args...)
			commands = append(commands, strings.Join(cmdParts, " "))
		}
	}
	
	if len(commands) == 0 {
		fmt.Println(theme.Warn("No modifications to apply"))
		return nil
	}
	
	// Create temporary file with commands
	tempFile, err := os.CreateTemp("", "taskvanguard-edit-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write commands to temp file with explanatory header
	header := "#!/bin/bash\n# TaskVanguard Analysis Results - Mass Edit Mode\n#\n# The commands below will modify your TaskWarrior tasks based on AI analysis.\n# You can:\n# - Delete any lines you don't want to execute\n# - Modify the task modify commands as needed\n# - Add comments with # prefix\n# - All remaining commands will be executed when you save and exit\n#\n# Format: task modify <task_id> <modifications>\n\n"
	content := header + strings.Join(commands, "\n") + "\n"
	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tempFile.Close()
	
	// Get editor from environment or use default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // fallback to nano
	}
	
	// Open editor
	fmt.Printf("Opening editor (%s) to edit commands...\n", editor)
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}
	
	// Read modified commands
	modifiedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read modified file: %w", err)
	}
	
	modifiedCommands := strings.Split(strings.TrimSpace(string(modifiedContent)), "\n")
	
	// Execute commands
	fmt.Println(theme.Title("\n=== Executing Commands ==="))
	for i, command := range modifiedCommands {
		command = strings.TrimSpace(command)
		if command == "" || strings.HasPrefix(command, "#") {
			continue // Skip empty lines and comments
		}
		
		fmt.Printf("Executing: %s\n", command)
		
		// Parse and execute the command
		parts := strings.Fields(command)
		if len(parts) < 3 || parts[0] != "task" || parts[1] != "modify" {
			fmt.Printf(theme.Warn("Skipping invalid command: %s\n"), command)
			continue
		}
		
		// Extract task ID and arguments
		taskID, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Printf(theme.Error("Invalid task ID in command: %s\n"), command)
			continue
		}
		
		args := parts[3:]
		_, err = client.ModifyTaskInTaskWarrior(taskID, args)
		if err != nil {
			fmt.Printf(theme.Error("Command %d failed: %s\n"), i+1, err.Error())
		} else {
			fmt.Printf(theme.Success("Command %d executed successfully\n"), i+1)
		}
	}
	
	return nil
}

func promptAnalyzeAllTasks(limit int) bool {
	fmt.Printf("%s %s %d %s %s: ", 
		theme.Title("→"), 
		theme.Info("Analyze all pending tasks (up to"), 
		limit, 
		theme.Info("tasks)?"), 
		"[Y]es/[n]o")
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	
	return input == "y" || input == "yes" || input == ""
}

func showFilterExamples() {
	fmt.Println(theme.Title("\n🧬 Filter Examples:"))
	fmt.Printf("  %s %s\n", theme.Info("Project:"), "vanguard analyze project:work")
	fmt.Printf("  %s %s\n", theme.Info("Tag:"), "vanguard analyze +urgent")
	fmt.Printf("  %s %s\n", theme.Info("Exclude tag:"), "vanguard analyze -someday")
	fmt.Printf("  %s %s\n", theme.Info("Priority:"), "vanguard analyze priority:H")
	fmt.Printf("  %s %s\n", theme.Info("Combination:"), "vanguard analyze project:work +urgent priority:H")
	fmt.Printf("  %s %s\n", theme.Info("Description:"), "vanguard analyze /meeting/")
	fmt.Println()
}
