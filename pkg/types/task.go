package types

type Task struct {
	ID          int       `json:"id,omitempty"`
	UUID        string    `json:"uuid"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority,omitempty"`
	Project     string    `json:"project,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Due         *string   `json:"due,omitempty"`
	Entry       TWTime    `json:"entry"`
	Modified    TWTime    `json:"modified"`
	Urgency     float64   `json:"urgency,omitempty"`
	Skipped		int		  `json:"skipped"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

type Annotation struct {
	Entry       string    `json:"entry"`
	Description string    `json:"description"`
}

// type Goal struct {
// 	ID          string    `json:"id"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	Deadline    *string   `json:"deadline,omitempty"`
// 	Priority    string    `json:"priority"`
// 	Created     TWTime    `json:"created"`
// 	Modified    TWTime    `json:"modified"`
// 	Status      string    `json:"status"`
// }

// type TaskAnalysis struct {
// 	TaskID           string   `json:"task_id"`
// 	SuggestedProject string   `json:"suggested_project,omitempty"`
// 	SuggestedTags    []string `json:"suggested_tags,omitempty"`
// 	SuggestedPriority string  `json:"suggested_priority,omitempty"`
// 	EstimatedEffort  string   `json:"estimated_effort,omitempty"`
// 	ImpactScore      int      `json:"impact_score"`
// 	SnowballRisk     string   `json:"snowball_risk"`
// 	Reasoning        string   `json:"reasoning"`
// }

// type ImpactAnalysis struct {
// 	TaskID            string   `json:"task_id"`
// 	ImpactScore       int      `json:"impact_score"`
// 	SnowballPotential string   `json:"snowball_potential"`
// 	ConsequenceDelay  string   `json:"consequence_delay"`
// 	RelatedTasks      []string `json:"related_tasks,omitempty"`
// 	Recommendations   []string `json:"recommendations"`
// }
