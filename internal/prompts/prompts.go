package prompts

import (
	"bytes"

	"text/template"

	"github.com/taskvanguard/taskvanguard/assets"
	"github.com/taskvanguard/taskvanguard/pkg/types"
)

// Structs shared for templating
type Task struct {
	Description  string
	Tags         []string
	Project      string
	Priority     string
	DueDate      string
	Annotations  []string
}

type Tag struct {
	Name        string
	Description string
}

type Goal struct {
	Description string
	Priority 	string
}

type UserContext struct {
	UserTags        	[]Tag
	UserAnnotations 	[]Annotation
	UserProjects    	[]string
	UserGoals           []Goal
}

type Annotation struct {
	Name        string
	Description string
}

type TemplateData struct {
	Task        	Task
	Tasks       	[]Task
	UserContext 	UserContext
	ExampleOutput 	string 
}

// RenderTemplate renders a Markdown template from path with the given data
func RenderTemplate(filename string, data TemplateData) (string, error) {
	tmplBytes, err := LoadPrompt(filename)
	if err != nil {
		return "", err
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	tmpl := template.New("prompt").Funcs(funcMap)

	// Load user_context.md template if it exists
	userContextBytes, err := LoadPrompt("user_context.md")
	if err == nil {
		_, err = tmpl.New("user_context.md").Parse(string(userContextBytes))
		if err != nil {
			return "", err
		}
	}

	_, err = tmpl.Parse(string(tmplBytes))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func LoadPrompt(filename string) (string, error) {
	prompt, err := assets.Load(filename)
	if err != nil {
		return "", err
	}
	return prompt, nil
}

func ToPromptGoals(tasks []types.Task) []Goal {

	goals := make([]Goal, 0, len(tasks))

	for _, t := range tasks {
		goals = append(goals, Goal{
			Description: t.Description,
			Priority:    t.Priority,
		})
	}
	return goals
}