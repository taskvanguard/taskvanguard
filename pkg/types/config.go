package types

type Config struct {
	LLM      	LLMConfig      				`yaml:"llm"`
	Tags     	map[string]TagsMeta    		`yaml:"tags"`     // <== flattened here
	Settings 	Settings	    			`yaml:"settings"`
	Annotations map[string]AnnotationsMeta  `yaml:"annotations"`
	Filters 	FiltersConfig			    `yaml:"filters"`
}

type FiltersConfig struct {
	TagFilterMode 			string 		`yaml:"tag_filter_mode"`
	TagFilterTags 			[]string  	`yaml:"tag_filter_tags"`
	ProjectFilterMode		string 		`yaml:"project_filter_mode"`
	ProjectFilterProjects	[]string 	`yaml:"project_filter_projects"`
}

type LLMConfig struct {
	Provider string `yaml:"provider"` // "openai" or "deepseek"
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
	BaseURL  string `yaml:"base_url"`
}

type TagsMeta struct {
	Desc          string  `yaml:"desc"`
	UrgencyFactor float64 `yaml:"urgency_factor"`
}

type Settings struct {
	Debug 			   		bool   `yaml:"debug"`
	EnableLLM 		   		bool   `yaml:"enable_llm"`
    SplitTasks         		bool   `yaml:"split_tasks"`
    AutoImportTags     		bool   `yaml:"auto_import_tags"`
    AutoImportProjects 		bool   `yaml:"auto_import_projects"`
	EnableLowercase	   		bool   `yaml:"enable_lowercase"`
	EnableTagging 	   		bool   `yaml:"enable_tagging"`
    EnableAnnotations  		bool   `yaml:"enable_annotations"`
	EnableGoals 	   		bool   `yaml:"enable_goals"`
    GoalProjectName    		string `yaml:"goal_project_name"`
	TaskImportLimit 		int	   `yaml:"task_import_limit"`
    TaskProcessingBatchSize int	   `yaml:"task_processing_batch_size"`
    GuidingQuestionAmount   int    `yaml:"guiding_question_amount"`
}

type AnnotationsMeta struct {
	Label          string  `yaml:"label"`
	Symbol string `yaml:"symbol"`
	Desc string `yaml:"description"`	
}

// var ErrMissingLLMKey = errors.New("LLM API key not configured. Run 'taskvanguard init' first")

// func (cfg *Config) Validate() error {
// 	if cfg.LLM.APIKey == "" {
// 		return ErrMissingLLMKey
// 	}
// 	return nil
// }