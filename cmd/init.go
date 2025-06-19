package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"xarc.dev/taskvanguard/internal/config"
	"xarc.dev/taskvanguard/pkg/theme"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup TaskVanguard and backup your Tasks",
	Long: `Initialize TaskVanguard with comprehensive setup including configuration,
shell integration, tag management, goal tracking, and task backup.`,
	Run: func(cmd *cobra.Command, args []string) {
		
		var twConfigPath string
		if twEnvPath := os.Getenv("TASKRC"); twEnvPath != "" {
			twConfigPath = twEnvPath
		} else {
			twConfigPath = filepath.Join(os.Getenv("HOME"), ".taskrc")
		}
		
		var configPath string
		if envPath := os.Getenv("TASKVANGUARD_CONFIG"); envPath != "" {
			configPath = envPath
		} else {
			configDir, err := os.UserConfigDir()
			if err != nil {
				fmt.Printf(theme.Error("config path could not be set: %v"), err)
				return
			}
			configPath = filepath.Join(configDir, "taskvanguard", "vanguardrc.yaml")
		}

		printIntro(configPath)

		var response string
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			// Step 1: Shell Integration
			setupShellAlias()

			// Step 2: Basic Initialization
			setupConfiguration(configPath)

			// Step 3: Task Backup
			setupTaskBackup()

			// Step 4: Tag Management & Analysis
			setupTagManagement()

			// Step 5: Goal Tracking Setup
			setupGoalTracking(twConfigPath)

			fmt.Println(theme.Success("🎉 TaskVanguard initialization complete!"))
			fmt.Println("")
			fmt.Println(theme.Info("You can now use 'vanguard' or 'tvg' (if alias was added) to access TaskVanguard commands. Try tvg add <task desc>."))
		} else {
			fmt.Println(theme.Info("Initialization cancelled."))
		}
	},
}

func printIntro(configPath string) {
	// Header
	header := []string{
		" _____         _   __     __                                    _ ",
		"|_   _|_ _ ___| | _\\ \\   / /_ _ _ __   __ _ _   _  __ _ _ __ __| |",
		"  | |/ _` / __| |/ /\\ \\ / / _` | '_ \\ / _` | | | |/ _` | '__/ _` |",
		"  | | (_| \\__ \\   <  \\ V / (_| | | | | (_| | |_| | (_| | | | (_| |",
		"  |_|\\__,_|___/_|\\_\\  \\_/ \\__,_|_| |_|\\__, |\\__,_|\\__,_|_|  \\__,_|",
		"                                      |___/                       ",
	}

	// Print header with Title color
	for _, line := range header {
		fmt.Println(theme.Title(line))
	}

	// Message
	// fmt.Println(theme.Title("It's a good idea to make sure your Tasks have:"))
	// fmt.Println(theme.Success("✔") + " " + theme.Title("Good descriptive Titles"))
	// fmt.Println(theme.Success("✔") + " " + theme.Title("Annotations with additional infos"))
	// fmt.Println(theme.Success("✔") + " " + theme.Title("Tags that are self-explanatory"))
	fmt.Println("")
	fmt.Printf("Configuration will be stored at: \n%s\n", theme.Info(configPath))
	fmt.Println("")
	fmt.Print(theme.Title("🚀 Do you want to initialize TaskVanguard? ") + theme.Info("y/n") + ": ")
}

func setupShellAlias() {
	fmt.Println("")
	fmt.Printf(theme.Title("Would you like to add a shortcut alias '%s' for TaskVanguard to your shell config? ") + theme.Info("y/n") + ": ", theme.Info("tvg"))
	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
		fmt.Println(theme.Title("        1) ALIAS"))
		fmt.Println(theme.Title("───────────────────────────────────────────────"))
		fmt.Println(theme.Title("You might want to add this line to your shell config file:"))
		fmt.Println(theme.Success("alias tvg=\"vanguard\""))
		fmt.Println("Then run: " + theme.Info("source ~/.bashrc") + " (or your shell config file)")
		fmt.Println("")
		fmt.Println("Press Enter to continue...")
		fmt.Scanln(&response)
		return
	}
}

func setupConfiguration(configPath string) {
	fmt.Println(theme.Title("───────────────────────────────────────────────"))
	fmt.Println(theme.Title("        2) BASIC CONFIGURATION"))
	fmt.Println(theme.Title("───────────────────────────────────────────────"))

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Print(theme.Warn("Configuration file already exists. Overwrite? ") + theme.Info("y/n") + ": ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println(theme.Info("Keeping existing configuration."))
			return
		}
	}

	// Create config directory
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf(theme.Error("Failed to create config directory: %v"), err)
		return
	}

	fmt.Println("")

	// // Prompt for LLM provider
	// fmt.Print(theme.Title("Enter Name of LLM provider (openai/deepseek): ") + ": ")
	// var provider string
	// fmt.Scanln(&provider)
	// if provider == "" {
	// 	provider = "openai"
	// }

	// // Prompt for LLM model
	// fmt.Print(theme.Title("Enter model: ") + ": ")
	// var model string
	// fmt.Scanln(&model)
	// if model == "" {
	// 	model = "https://openrouter.ai/api/v1"
	// }

	// // Prompt for LLM Baseurl
	// fmt.Print(theme.Title("Enter baseurl (openai/deepseek): ") + ": ")
	// var baseUrl string
	// fmt.Scanln(&baseUrl)
	// if baseUrl == "" {
	// 	baseUrl = "https://openrouter.ai/api/v1"
	// }

	// // Prompt for API key
	// fmt.Print(theme.Title("Enter your API key: "))
	// reader := bufio.NewReader(os.Stdin)
	// apiKey, _ := reader.ReadString('\n')
	// apiKey = strings.TrimSpace(apiKey)

	// // Create basic config
	// cfg := &types.Config{
	// 	LLM: types.LLMConfig{
	// 		Provider: provider,
	// 		APIKey:   apiKey,
	// 		Model:    getDefaultModel(provider),
	// 	},
	// 	Tags: make(map[string]types.TagsMeta),
	// }

	fmt.Println(theme.Success("✔ Generated Config file successfully!"))
	fmt.Println("")

	provider := "openai"
	apiKey := "YOUR_API_KEY_HERE"
	model := "gpt-3.5-turbo"
	fmt.Println("")

	// cfg := &types.Config{
	// 	LLM: types.LLMConfig{
	// 		Provider: provider,
	// 		APIKey:   apiKey,
	// 		Model:    model,
	// 	},
	// 	Tags: make(map[string]types.TagsMeta),
	// }

	fmt.Printf(theme.Title("LLM:\n Provider: %s\n APIKey: \"%s\"\n Model: %s\n"), provider, apiKey, model)

	// Save config
	_, err := config.CreateDefaultConfig(configPath)
	// data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf(theme.Error("Failed to create configuration file: %v"), err)
		return
	}

	// if err := os.WriteFile(configPath, data, 0644); err != nil {
	// 	fmt.Printf(theme.Error("Failed to save configuration: %v"), err)
	// 	return
	// }

	fmt.Println(theme.Warn("Go here and put your API Credentials in:"))
	fmt.Println(theme.Info(configPath))

	fmt.Println("")
	fmt.Println("Press Enter to continue...")
	var response string
	fmt.Scanln(&response)

}

func setupTaskBackup() {
	fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
	fmt.Println(theme.Title("        3) TASK BACKUP"))
	fmt.Println(theme.Title("───────────────────────────────────────────────"))
	fmt.Print(theme.Title("Would you like to back up your existing tasks before proceeding? ") + theme.Info("y/n") + ": ")
	fmt.Println("")

	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		backupPath := filepath.Join(os.Getenv("HOME"), ".config", "taskvanguard", "task_backup.json")
		cmd := exec.Command("task", "export")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf(theme.Error("Failed to export tasks: %v"), err)
			return
		}

		if err := os.WriteFile(backupPath, output, 0644); err != nil {
			fmt.Printf(theme.Error("Failed to save backup: %v"), err)
			return
		}

		fmt.Printf(theme.Success("✔ Tasks backed up to: %s"), backupPath)
		fmt.Println("")
	}
	fmt.Println("")
	fmt.Println("Press Enter to continue...")
	fmt.Scanln(&response)
}

func setupTagManagement() {
	fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
	fmt.Println(theme.Title("        4) FINAL STEP"))
	fmt.Println(theme.Title("───────────────────────────────────────────────"))

	// var response string
	// fmt.Print(theme.Title("Would you like to analyze existing TaskWarrior tasks for tags? ") + theme.Info("y/n") + ": ")
	// fmt.Scanln(&response)
	// if response == "y" || response == "Y" {
	// 	// Get existing tags from TaskWarrior
	// 	tags, err := getExistingTags()
	// 	if err != nil {
	// 		fmt.Printf(theme.Error("Failed to get existing tags: %v"), err)
	// 	} else if len(tags) > 0 {
	// 		fmt.Println(theme.Info("Found existing tags:"))
	// 		for _, tag := range tags {
	// 			fmt.Printf("  %s\n", theme.Success(tag))
	// 		}
	// 		fmt.Println("")
	// 		fmt.Print(theme.Title("Import all these tags or use default set? ") + theme.Info("import/default") + ": ")
	// 		fmt.Scanln(&response)
	// 		if response == "import" {
	// 			fmt.Println(theme.Success("✔ Existing tags will be imported during analysis."))
	// 		}
	// 	}
	// }
	fmt.Println(theme.Title("These are the tags and annotations that will be assigned via LLM."))
	fmt.Println("You can change these in your config.")
	fmt.Println("")
	// Suggest common useful tags
	fmt.Println(theme.Info("Suggested tags:"))
	commonTags := []string{
		theme.Info("+sb: ") + "snowball", 
		theme.Info("+cut: ") + "saves time/money", 
		theme.Info("+fast: ") + "quick task", 
		theme.Info("+key: ") + "high impact",
	}
	for _, tag := range commonTags {
		fmt.Printf("  %s\n", theme.Success(tag))
	}
	fmt.Println("")

	fmt.Println(theme.Info("Suggested annotations:"))
	annotations := []string{
		theme.Info("short_reward: ") + "immediate benefit",
		theme.Info("long_reward: ") + "strategic benefit",
		theme.Info("risk: ") + "if not done",
		theme.Info("tip: ") + "practical, actionable, insightful",
	}

	for _, annotation := range annotations {
		fmt.Printf("  %s\n", theme.Success(annotation))
	}
	fmt.Println("")
}

func setupGoalTracking(twConfigPath string) {
	fmt.Println(theme.Title("\n───────────────────────────────────────────────"))
	fmt.Println(theme.Title("        5) GOAL TRACKING"))
	fmt.Println(theme.Title("───────────────────────────────────────────────"))

	fmt.Print(theme.Title("Would you like to enable goal tracking features? (recommended) ") + theme.Info("y/n") + ": ")
	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		fmt.Println("")
		fmt.Println(theme.Info("Goal tracking enables:"))
		fmt.Println("  • Create goals as tasks in project:goals")
		fmt.Println("  • Link other tasks to goals via goal:\"name\" in task commands")
		fmt.Println("  • Auto-assign new tasks to goals")
		fmt.Println("")

		fmt.Println(theme.Warn("Required: Add these UDAs (User-Defined Attributes) to your .taskrc:"))
		fmt.Println(theme.Info(twConfigPath))		
		fmt.Println("")
		fmt.Println(theme.Success("uda.goal.label=Goal"))
		fmt.Println(theme.Success("uda.goal.type=string"))
		fmt.Println(theme.Success("uda.goal.values="))
		fmt.Println(theme.Success("uda.skipped.label=Skipped"))
		fmt.Println(theme.Success("uda.skipped.type=numeric"))
		fmt.Println(theme.Success("uda.skipped.values=0"))
		fmt.Println("")

		fmt.Print(theme.Title("Would you like me to add these to your .taskrc automatically? ") + theme.Info("y/n") + ": ")
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			if err := addUDAsToTaskrc(twConfigPath); err != nil {
				fmt.Printf(theme.Error("Failed to update .taskrc: %v"), err)
				fmt.Println(theme.Warn("Please add the UDAs manually to your .taskrc"))
			} else {
				fmt.Println(theme.Success("✔ UDAs added to .taskrc successfully!"))
			}
		}
	}
	fmt.Println("")
}

// func getDefaultModel(provider string) string {
// 	switch provider {
// 	case "deepseek":
// 		return "deepseek-chat"
// 	default:
// 		return "gpt-3.5-turbo"
// 	}
// }

// func getExistingTags() ([]string, error) {
// 	cmd := exec.Command("task", "_tags")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}

// 	tagsStr := strings.TrimSpace(string(output))
// 	if tagsStr == "" {
// 		return []string{}, nil
// 	}

// 	return strings.Fields(tagsStr), nil
// }

func addUDAsToTaskrc(twConfigPath string) error {
	
	// Read existing .taskrc
	content := ""
	if data, err := os.ReadFile(twConfigPath); err == nil {
		content = string(data)
	}

	// UDAs to add
	udas := []string{
		"uda.goal.label=Goal",
		"uda.goal.type=string",
		"uda.goal.values=",
		"uda.skipped.label=Skipped",
		"uda.skipped.type=numeric",
		"uda.skipped.values=0",
	}

	// Check if UDAs already exist
	for _, uda := range udas {
		if !strings.Contains(content, uda) {
			content += "\n" + uda
		}
	}

	// Write back to .taskrc
	return os.WriteFile(twConfigPath, []byte(content), 0644)
}

