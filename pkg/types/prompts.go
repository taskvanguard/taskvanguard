package types

type TaskSuggestion struct {
	SuggestedTags  []string            `json:"suggested_tags"`
	GoalAlignment  string              `json:"goal_alignment"`
	Project        string              `json:"project"`
	RefinedTask    string              `json:"refined_task"`
	AdditionalInfo map[string]string   `json:"additional_infos"`
	Subtasks       []string            `json:"subtasks"`
}

type TaskAnalysisResult struct {
	TaskIndex      int                 `json:"task_index"`
	SuggestedTags  []string            `json:"suggested_tags"`
	GoalAlignment  string              `json:"goal_alignment"`
	Project        string              `json:"project"`
	RefinedTask    string              `json:"refined_task"`
	AdditionalInfo map[string]string   `json:"additional_infos"`
	Subtasks       []string            `json:"subtasks"`
}

type BatchTaskSuggestion struct {
	TaskAnalyses []TaskAnalysisResult `json:"task_analyses"`
}

// type Annotation struct {
// 	Name        string
// 	Description string
// }

// type AdditionalInfos struct {
// 	ShortReward string `json:"short_reward"`
// 	LongReward  string `json:"long_reward"`
// 	Risk        string `json:"risk"`
// 	Tip         string `json:"tip"`
// }


// // Structs shared for templating
// type Task struct {
// 	Description  string
// 	Tags         []string
// 	Project      string
// 	Priority     string
// 	DueDate      string
// 	Annotations  []string
// }

// type Tag struct {
// 	Name        string
// 	Description string
// }

// type Goal struct {
// 	Description string
// 	Priority 	string
// }

// type UserContext struct {
// 	UserTags        	[]Tag
// 	UserAnnotations 	[]Annotation
// 	UserProjects    	[]string
// 	UserGoals           []Goal
// }

// type Annotation struct {
// 	Name        string
// 	Description string
// }

// type TemplateData struct {
// 	Task        Task
// 	UserContext UserContext
// }
