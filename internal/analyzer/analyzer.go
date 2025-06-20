package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/taskvanguard/taskvanguard/internal/llm"
	"github.com/taskvanguard/taskvanguard/internal/prompts"
	"github.com/taskvanguard/taskvanguard/pkg/filter"
	"github.com/taskvanguard/taskvanguard/pkg/theme"
	"github.com/taskvanguard/taskvanguard/pkg/types"
	"github.com/taskvanguard/taskvanguard/pkg/utils"
)

// cleanMarkdownCodeFences removes markdown code fence markers from LLM responses
func cleanMarkdownCodeFences(response string) string {
	response = strings.TrimSpace(response)
	
	// Remove ```json at the beginning
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSpace(response)
	}
	
	// Remove ``` at the end
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}
	
	return response
}

func AnalyzeSingleTaskWithLLM(cfg *types.Config, taskArgs string, userGoals []types.Task, projects []string) (*types.TaskSuggestion, error) {
	args := utils.ParseTaskArgs(taskArgs)

	if !filter.ShouldIncludeByTags(args.Tags, cfg.Filters) {
		return nil, fmt.Errorf("at least one of the tags is blacklisted for LLM processing: tags=%v", args.Tags)
	}

	if !filter.ShouldIncludeByProject(args.Project, cfg.Filters) {
		return nil, fmt.Errorf("the project this task is assigned to is blacklisted for LLM processing: project=%q", args.Project)
	}

	task := prompts.Task{
		Description: args.Title,
		Tags:        args.Tags,
		Project:     args.Project,
		Priority:    args.Priority,
		// DueDate:     "2025-06-15",
	}

	data := buildTemplateData(cfg, []prompts.Task{task}, userGoals, projects)
	data.Task = task

	response, err := sendLLMRequest(cfg, "task_analysis_single.md", data)
	if err != nil {
		return nil, err
	}

	// Clean markdown code fences from response
	cleanedResponse := cleanMarkdownCodeFences(response)

	var suggestion types.TaskSuggestion
	if err := json.Unmarshal([]byte(cleanedResponse), &suggestion); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	return &suggestion, nil
}

func AnalyzeBatchTasksWithLLM(cfg *types.Config, taskArgsList []string, userGoals []types.Task, projects []string) (*types.BatchTaskSuggestion, error) {
	const batchSize = 20
	var allAnalyses []types.TaskAnalysisResult

	// Process tasks in chunks of 15
	for i := 0; i < len(taskArgsList); i += batchSize {
		end := i + batchSize
		if end > len(taskArgsList) {
			end = len(taskArgsList)
		}

		batchTaskArgs := taskArgsList[i:end]
		tasks := make([]prompts.Task, 0, len(batchTaskArgs))
		taskIndexOffset := i

		for _, taskArgs := range batchTaskArgs {
			args := utils.ParseTaskArgs(taskArgs)

			if !filter.ShouldIncludeByTags(args.Tags, cfg.Filters) {
				return nil, fmt.Errorf("at least one of the tags is blacklisted for LLM processing: tags=%v", args.Tags)
			}

			if !filter.ShouldIncludeByProject(args.Project, cfg.Filters) {
				return nil, fmt.Errorf("the project this task is assigned to is blacklisted for LLM processing: project=%q", args.Project)
			}

			tasks = append(tasks, prompts.Task{
				Description: args.Title,
				Tags:        args.Tags,
				Project:     args.Project,
				Priority:    args.Priority,
				// DueDate:     "2025-06-15",
			})
		}

		data := buildTemplateData(cfg, tasks, userGoals, projects)
		data.Tasks = tasks

		response, err := sendLLMRequest(cfg, "task_analysis_batch.md", data)
		if err != nil {
			return nil, fmt.Errorf("failed to process batch %d-%d: %v", i+1, end, err)
		}

		// Clean markdown code fences from response
		cleanedResponse := cleanMarkdownCodeFences(response)

		var batchSuggestion types.BatchTaskSuggestion
		if err := json.Unmarshal([]byte(cleanedResponse), &batchSuggestion); err != nil {
			return nil, fmt.Errorf("failed to parse LLM response for batch %d-%d: %v", i+1, end, err)
		}

		// Adjust task indices to reflect global position
		for j := range batchSuggestion.TaskAnalyses {
			batchSuggestion.TaskAnalyses[j].TaskIndex = taskIndexOffset + j + 1
		}

		allAnalyses = append(allAnalyses, batchSuggestion.TaskAnalyses...)
	}

	return &types.BatchTaskSuggestion{
		TaskAnalyses: allAnalyses,
	}, nil
}

func BuildExampleJSON(userAnnotations []prompts.Annotation) string {
	buf := &bytes.Buffer{}
	buf.WriteString("{\n")
	for i, ann := range userAnnotations {
		// Always output the value as the annotation description (escaped)
		descJSON, _ := json.Marshal(ann.Description) // This will handle quotes, escapes, etc.
		fmt.Fprintf(buf, `  "%s": %s`, ann.Name, descJSON)
		if i != len(userAnnotations)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}
	buf.WriteString("}")
	return buf.String()
}

func buildTemplateData(cfg *types.Config, tasks []prompts.Task, userGoals []types.Task, projects []string) prompts.TemplateData {
	data := prompts.TemplateData{
		UserContext: prompts.UserContext{
			UserTags:        []prompts.Tag{},
			UserAnnotations: []prompts.Annotation{},
			UserProjects:    projects,
			UserGoals:       prompts.ToPromptGoals(userGoals),
		},
	}

	// Add Tags and Annotations from config
	for name, meta := range cfg.Tags {
		if !filterTags(name, cfg.Filters.TagFilterMode, cfg.Filters.TagFilterTags) {
			continue
		}

		data.UserContext.UserTags = append(data.UserContext.UserTags, prompts.Tag{
			Name:        name,
			Description: meta.Desc,
		})
	}

	for _, meta := range cfg.Annotations {
		data.UserContext.UserAnnotations = append(data.UserContext.UserAnnotations, prompts.Annotation{
			Name:        meta.Label,
			Description: meta.Desc,
		})
	}

	data.ExampleOutput = BuildExampleJSON(data.UserContext.UserAnnotations)

	return data
}

func sendLLMRequest(cfg *types.Config, templateName string, data prompts.TemplateData) (string, error) {
	llmClient, err := llm.NewClient(&cfg.LLM)
	if err != nil {
		return "", err
	}

	rendered, err := prompts.RenderTemplate(templateName, data)
	if err != nil {
		return "", err
	}

	messages := []llm.Message{
		{Role: "user", Content: rendered},
	}

	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Request"))
		fmt.Println(theme.Info(messages))
	}

	if !cfg.Settings.EnableLLM {
		return "", fmt.Errorf("sending API Request to LLM is disabled via config")
	}

	response, err := llmClient.Chat(messages)
	if cfg.Settings.Debug {
		fmt.Println(theme.Title("LLM Response"))
		fmt.Println(theme.Info(response))
	}
	if err != nil {
		return "", err
	}

	return response, nil
}

func filterTags(tagName string, mode string, list []string) bool {
	switch mode {
	case "blacklist":
		for _, v := range list {
			if tagName == v {
				return false // skip if blacklisted
			}
		}
		return true // include if not blacklisted
	case "whitelist":
		for _, v := range list {
			if tagName == v {
				return true // include if whitelisted
			}
		}
		return false // skip if not whitelisted
	default:
		return true // include everything if mode is unknown
	}
}

