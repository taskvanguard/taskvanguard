package cmd

import (
	"github.com/spf13/cobra"
)

var goalsCmd = &cobra.Command{
	Use:   "goals",
	Short: "Manage your goals and track task alignment",
	Long: `Manage your long-term goals and see how your current tasks
align with achieving those goals.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement goals management
		// First run taskvanguard init
		// analyze manifest.md
		// Socratic goal refinement: The LLM challenges vague goals.
		//“I want to be healthier” → “What does that mean in your current environment? Define success in 1, 3, and 12 months.”

	},
}

func init() {
	goalsCmd.AddCommand(goalsListCmd)
	goalsCmd.AddCommand(goalsAddCmd)
	goalsCmd.AddCommand(goalsAlignCmd)
}

var goalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all goals",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: List goals
	},
}

var goalsAddCmd = &cobra.Command{
	Use:   "add [goal description]",
	Short: "Add a new goal",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Add goal
	},
}

var goalsAlignCmd = &cobra.Command{
	Use:   "align",
	Short: "Show how current tasks align with goals",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Show task-goal alignment
	},
}