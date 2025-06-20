package taskwarrior

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/taskvanguard/taskvanguard/internal/config"

	"github.com/taskvanguard/taskvanguard/pkg/types"
)

type RuntimeContext struct {
	Config       *types.Config
	Client       *Client
	UserGoals    []types.Task
	UserProjects []string
	UserAllTasks []types.Task
}

func Bootstrap(cmd *cobra.Command) (*RuntimeContext, error) {
// func Bootstrap(cmd *cobra.Command, withTasks bool) (*RuntimeContext, error) {
	client := NewClient()
	if !client.IsAvailable() {
		return nil, errors.New("TaskWarrior not found. Please install TaskWarrior first")
	}

	cfg, err := loadAndApplyFlagOverrides(cmd)
	if err != nil {
		return nil, err
	}

	// 🔄 Enrich config after loading + applying CLI flags
	if err := EnrichConfigWithTW(cfg, client); err != nil {
		return nil, fmt.Errorf("failed to enrich config with Taskwarrior tags: %v", err)
	}

	if cfg.LLM.APIKey == "" {
		return nil, errors.New("LLM API key not configured. Run 'taskvanguard init' first")
	}

	goals, err := client.GetGoalsFiltered(cfg)
	if err != nil {
		return nil, err
	}

	projects, err := client.GetProjectsFiltered(cfg)
	if err != nil {
		return nil, err
	}

	// var tasks []types.Task
	// if withTasks {
	// 	tasks, err = client.GetTasksFiltered(cfg)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return &RuntimeContext{
		Config:       cfg,
		Client:       client,
		UserGoals:    goals,
		UserProjects: projects,
		// UserAllTasks: tasks,
	}, nil
}

func loadAndApplyFlagOverrides(cmd *cobra.Command) (*types.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}
	if noSubtask, _ := cmd.Flags().GetBool("no-subtasks"); noSubtask {
		cfg.Settings.SplitTasks = false
	}
	if noAnnotations, _ := cmd.Flags().GetBool("no-annotations"); noAnnotations {
		cfg.Settings.EnableAnnotations = false
	}
	if noTags, _ := cmd.Flags().GetBool("no-tags"); noTags {
		cfg.Settings.EnableTagging = false
	}
	return cfg, nil
}

func EnrichConfigWithTW(cfg *types.Config, client *Client) error {
	if !cfg.Settings.AutoImportTags {
		return nil
	}

	if cfg.Tags == nil {
		cfg.Tags = make(map[string]types.TagsMeta)
	}

	tags, err := client.GetTagsFiltered(cfg)
	if err != nil {
		return err
	}

	for tag := range tags {
		if _, exists := cfg.Tags[tag]; !exists {
			cfg.Tags[tag] = types.TagsMeta{
				Desc:          "",      // auto import
				UrgencyFactor: 1.0,
			}
		}
	}

	return nil
}
