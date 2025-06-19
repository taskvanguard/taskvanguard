package utils

import (
	"strings"
)

type ParsedTask struct {
	Title    string
	Tags     []string
	Project  string
	Priority string
}

// prefixMatches checks if input matches one of the expected keywords (e.g., "project", "priority")
func prefixMatches(input string, full string) bool {
	return strings.HasPrefix(full, input)
}

func ParseTaskArgs(taskArgs string) ParsedTask {
	var parsed ParsedTask
	var titleParts []string

	args := strings.Fields(taskArgs)

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "+"):
			parsed.Tags = append(parsed.Tags, strings.TrimPrefix(arg, "+"))

		case strings.Contains(arg, ":"):
			parts := strings.SplitN(arg, ":", 2)
			key := strings.ToLower(parts[0])
			value := parts[1]

			switch {
			case prefixMatches(key, "priority"):
				parsed.Priority = strings.ToUpper(value)

			case prefixMatches(key, "project"):
				parsed.Project = value
			}

		default:
			titleParts = append(titleParts, arg)
		}
	}

	parsed.Title = strings.Join(titleParts, " ")
	return parsed
}