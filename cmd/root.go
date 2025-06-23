package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskvanguard",
	Short: "AI-powered task guide and motivator for TaskWarrior",
	Long: `TaskVanguard is a CLI tool that works alongside TaskWarrior to provide
AI-powered guidance, motivation, and task categorization using LLM APIs.

🚀 CORE FEATURES:
• Smart adding - AI-powered tagging (+sb, +cut, +fast, etc.) and annotations
• Goal management - Link tasks to goals
• Guidance - Helps figuring out the next best task to do and roadmaps 

📋 AVAILABLE COMMANDS:
• init     - Complete setup wizard (config, backup, tags, goals)
• add      - Enhanced task creation with AI assistance
• analyze  - Analyze task and provides recommendations
• spot     - Picks one high impact, high urgency task to do right now
• guide    - Asks a series of questions -> generates roadmap to achieve goal
• goals    - Manage strategic goals and link tasks to them

🔧 CONFIGURATION:
Config stored at: ~/.config/taskvanguard/vanguardrc.yaml
Supports OpenAI and DeepSeek LLM providers

⚔️ QUICK START:
1. Run 'taskvanguard init' to set up configuration
2. Try 'taskvanguard analyze' for detailed analysis
3. Use 'taskvanguard spot' to start completing tasks

For any unrecognized commands, TaskVanguard forwards them directly to TaskWarrior.`,
	Run: forwardToTaskWarrior,
}

// Setup and wire up subcommands here
func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(spotCmd)
	rootCmd.AddCommand(goalsCmd)
	rootCmd.AddCommand(guideCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

// Forwards unrecognized commands directly to Taskwarrior.
func forwardToTaskWarrior(cmd *cobra.Command, args []string) {
	if len(os.Args) <= 1 {
		fmt.Println("Please provide task arguments (same as 'task' command).")
		os.Exit(1)
	}

	// Forward everything after "taskvanguard" to "task"
	c := exec.Command("task", os.Args[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running task: %v\n", err)
		os.Exit(1)
	}
}
