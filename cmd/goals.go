package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/taskvanguard/taskvanguard/internal/goals"
	"github.com/taskvanguard/taskvanguard/internal/taskwarrior"
	"github.com/taskvanguard/taskvanguard/pkg/theme"
)

// getGoalsManager creates a goals manager with the current config
func getGoalsManager(cmd *cobra.Command) (*goals.Manager, error) {
	env, err := taskwarrior.Bootstrap(cmd)
	if err != nil {
		return nil, err
	}
	return goals.NewManager(env.Config), nil
}

var goalsCmd = &cobra.Command{
	Use:   "goals",
	Short: "Manage your goals and track task alignment",
	Long: `Manage your long-term goals and see how your current tasks
align with achieving those goals.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	goalsCmd.AddCommand(goalsListCmd)
	goalsCmd.AddCommand(goalsAddCmd)
	goalsCmd.AddCommand(goalsShowCmd)
	goalsCmd.AddCommand(goalsModifyCmd)
	goalsCmd.AddCommand(goalsDeleteCmd)
	goalsCmd.AddCommand(goalsLinkCmd)
	goalsCmd.AddCommand(goalsUnlinkCmd)
	goalsCmd.AddCommand(goalsLinksCmd)
	// goalsCmd.AddCommand(goalsAlignCmd)
}

var goalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all goals",
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			color.Red("Error initializing goals manager: %v", err)
			return
		}

		goals, err := goalsManager.ListGoals()
		if err != nil {
			color.Red("Error listing goals: %v", err)
			return
		}

		if len(goals) == 0 {
			color.Yellow("No goals found.")
			return
		}

		color.Green("Goals:")
		for _, goal := range goals {
			fmt.Printf("  %d: %s\n", goal.ID, goal.Description)
		}
	},
}

var goalsAddCmd = &cobra.Command{
	Use:   "add [goal description]",
	Short: "Add a new goal",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			color.Red("Error initializing goals manager: %v", err)
			return
		}

		output, taskID, err := goalsManager.AddGoal(args)
		if err != nil {
			color.Red("Error adding goal: %v", err)
			return
		}

		color.Green("Goal added successfully:")
		fmt.Printf("  ID: %d\n", taskID)
		fmt.Printf("  Output: %s\n", strings.TrimSpace(output))
	},
}

var goalsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show goal or task details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			fmt.Println(theme.Error(fmt.Sprintf("Error initializing goals manager: %v", err)))
			return
		}

		task, err := goalsManager.ShowGoal(args[0])
		if err != nil {
			fmt.Println(theme.Error(fmt.Sprintf("Error showing task/goal: %v", err)))
			return
		}

		if task == nil {
			fmt.Println(theme.Warn("Task/goal not found."))
			return
		}

		// Determine if this is a goal or task
		isGoal, err := goalsManager.IsGoal(args[0])
		if err != nil {
			fmt.Println(theme.Error(fmt.Sprintf("Error determining if ID is a goal: %v", err)))
			return
		}

		// Display appropriate header
		fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
		if isGoal {
			fmt.Println(theme.Title("          🏔️ GOAL DETAILS:"))
		} else {
			fmt.Println(theme.Title("          ✅ TASK DETAILS:"))
		}
		fmt.Println(theme.Title("───────────────────────────────────────────────"))

		// Display details with colors
		fmt.Printf("  %s %d\n", theme.Info("ID:"), task.ID)
		fmt.Printf("  %s %s\n", theme.Info("UUID:"), task.UUID)
		fmt.Printf("  %s %s\n", theme.Info("Description:"), task.Description)
		fmt.Printf("  %s %s\n", theme.Info("Status:"), task.Status)
		
		if task.Project != "" {
			fmt.Printf("  %s %s\n", theme.Info("Project:"), task.Project)
		}
		
		if len(task.Tags) > 0 {
			fmt.Printf("  %s %s\n", theme.Info("Tags:"), strings.Join(task.Tags, ", "))
		}
		
		if task.Priority != "" {
			fmt.Printf("  %s %s\n", theme.Info("Priority:"), task.Priority)
		}
		
		fmt.Printf("  %s %.2f\n", theme.Info("Urgency:"), task.Urgency)
	},
}

var goalsModifyCmd = &cobra.Command{
	Use:   "modify <goal_id> [arguments...]",
	Short: "Modify an existing goal",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			color.Red("Error initializing goals manager: %v", err)
			return
		}

		goalID, err := strconv.Atoi(args[0])
		if err != nil {
			color.Red("Invalid goal ID: %v", err)
			return
		}

		output, err := goalsManager.ModifyGoal(goalID, args[1:])
		if err != nil {
			color.Red("Error modifying goal: %v", err)
			return
		}

		color.Green("Goal modified successfully:")
		fmt.Printf("  Output: %s\n", strings.TrimSpace(output))
	},
}

var goalsDeleteCmd = &cobra.Command{
	Use:   "delete <goal_id>",
	Short: "Delete a goal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			color.Red("Error initializing goals manager: %v", err)
			return
		}

		err = goalsManager.DeleteGoal(args[0])
		if err != nil {
			color.Red("Error deleting goal: %v", err)
			return
		}

		color.Green("Goal deleted successfully.")
	},
}

var goalsLinkCmd = &cobra.Command{
	Use:   "link <id1> <id2>",
	Short: "Link a task and goal together (order-agnostic)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			color.Red("Error initializing goals manager: %v", err)
			return
		}

		err = goalsManager.Link(args[0], args[1])
		if err != nil {
			color.Red("Error linking task and goal: %v", err)
			return
		}

		color.Green("Task and goal linked successfully.")
	},
}

var goalsUnlinkCmd = &cobra.Command{
	Use:   "unlink <id1> <id2>",
	Short: "Remove link between task and goal (order-agnostic)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			color.Red("Error initializing goals manager: %v", err)
			return
		}

		err = goalsManager.Unlink(args[0], args[1])
		if err != nil {
			color.Red("Error unlinking task and goal: %v", err)
			return
		}

		color.Green("Task and goal unlinked successfully.")
	},
}

var goalsLinksCmd = &cobra.Command{
	Use:   "links <id>",
	Short: "Show links for a task or goal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		goalsManager, err := getGoalsManager(cmd)
		if err != nil {
			color.Red("Error initializing goals manager: %v", err)
			return
		}

		links, err := goalsManager.ShowLinks(args[0])
		if err != nil {
			color.Red("Error showing links: %v", err)
			return
		}

		if len(links) == 0 {
			color.Yellow("No links found for this task/goal.")
			return
		}

		isGoal, err := goalsManager.IsGoal(args[0])
		if err != nil {
			color.Red("Error determining if ID is a goal: %v", err)
			return
		}

		if isGoal {
			color.Green("Tasks linked to this goal:")
		} else {
			color.Green("Goals linked to this task:")
		}

		for _, link := range links {
			fmt.Printf("  %d: %s\n", link.ID, link.Description)
		}
	},
}

// var goalsAlignCmd = &cobra.Command{
// 	Use:   "align",
// 	Short: "Show how current tasks align with goals",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// TODO: Implement alignment analysis
// 		color.Yellow("Goal alignment analysis not yet implemented.")
// 	},
// }