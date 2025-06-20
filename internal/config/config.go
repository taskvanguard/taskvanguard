package config

import (
	"os"
	"path/filepath"

	"github.com/taskvanguard/taskvanguard/pkg/types"
	"gopkg.in/yaml.v3"
)

func Load() (*types.Config, error) {
	var configPath string
	
	if envPath := os.Getenv("TASKVANGUARD_CONFIG"); envPath != "" {
		configPath = envPath
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(configDir, "taskvanguard", "vanguardrc.yaml")
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return CreateDefaultConfig(configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func CreateDefaultConfig(configPath string) (*types.Config, error) {
	config := &types.Config{
		Settings: types.Settings{
			Debug: false,
			SplitTasks: true,
			EnableLLM: true,
			EnableGoals: true,
			GoalProjectName: "goals",
			TaskImportLimit: 500,
			TaskProcessingBatchSize: 15,
			GuidingQuestionAmount: 6,
		},
		LLM: types.LLMConfig{
			Provider: "openai",
			Model:    "gpt-3.5-turbo",
			BaseURL:  "https://openrouter.ai/api/v1",
			APIKey:   "<api key>",
		},
		Filters: types.FiltersConfig{
			TagFilterMode: "blacklist",
			TagFilterTags: []string{"private", "confidential"},
			ProjectFilterMode: "blacklist",
			ProjectFilterProjects: []string{"pers.secret", "work.secret"},
		},
		Tags: map[string]types.TagsMeta{
			"cut": {
				Desc:          "Task has the potential to save time or cost in the future",
				UrgencyFactor: 1.2,
			},
			"key": {
				Desc:          "Task is impacting goals",
				UrgencyFactor: 1.2,
			},
			"fast": {
				Desc:          "Task is probably done in very short time (10 mins or less)",
				UrgencyFactor: 1.2,
			},
			"sb": {
				Desc:          "Task is potentially snowballing positively or negatively and offers high roi",
				UrgencyFactor: 1.3,
			},
		},
		Annotations: map[string]types.AnnotationsMeta{
			"short_reward": {
				Label:       "Short Reward",
				Symbol:      "●",
				Desc: 	   	 "Immediate benefit",
			},
			"long_reward": {
				Label:       "Long Reward",
				Symbol:      "●",
				Desc: 	   	 "Strategic benefit",
			},
			"risk": {
				Label:       "Risk",
				Symbol:      "●",
				Desc: 	   	 "If not done",
			},
			"tip": {
				Label:       "Tip",
				Symbol:      "●",
				Desc: 	   	 "Practical, actionable",
			},
		},
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, err
	}


	return config, nil
}
