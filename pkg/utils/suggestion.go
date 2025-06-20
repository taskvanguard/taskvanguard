package utils

import (
	"github.com/taskvanguard/taskvanguard/pkg/types"
)

func TaskSuggestionToArgs(s types.TaskAnalysisResult) []string {
	var args []string

	if s.RefinedTask != "" {
		args = append(args, s.RefinedTask)
	}

	if len(s.SuggestedTags) > 0 {
		args = append(args, s.SuggestedTags...) 
		// for _, tag := range s.SuggestedTags {
		// 	args = append(args, tag)
		// }
	}

	if s.Project != "" {
		args = append(args, "project:"+s.Project)
	}

	if prio, ok := s.AdditionalInfo["priority"]; ok && prio != "" {
		args = append(args, "priority:"+prio)
	}

	return args
}
