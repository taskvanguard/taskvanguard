package filter

import (
	"strings"

	"github.com/taskvanguard/taskvanguard/pkg/types"
)

// FilterTasks filters tasks based on project and tag blacklist/whitelist configuration
func FilterTasks(tasks []types.Task, cfg *types.Config) []types.Task {
	if cfg == nil {
		return tasks
	}

	var filtered []types.Task
	for _, task := range tasks {
		if ShouldIncludeTask(task, cfg) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// ShouldIncludeTask determines if a task should be included based on filter configuration
func ShouldIncludeTask(task types.Task, cfg *types.Config) bool {
	// Check project filtering
	if !ShouldIncludeByProject(task.Project, cfg.Filters) {
		return false
	}

	// Check tag filtering
	if !ShouldIncludeByTags(task.Tags, cfg.Filters) {
		return false
	}

	return true
}

// ShouldIncludeByProject checks if task should be included based on project filters
func ShouldIncludeByProject(project string, filters types.FiltersConfig) bool {
	if len(filters.ProjectFilterProjects) == 0 {
		return true // No project filter configured
	}

	projectLower := strings.ToLower(project)
	containsProject := false

	for _, filterProject := range filters.ProjectFilterProjects {
		if strings.ToLower(filterProject) == projectLower {
			containsProject = true
			break
		}
	}

	switch strings.ToLower(filters.ProjectFilterMode) {
	case "blacklist":
		return !containsProject // Exclude if project is in blacklist
	case "whitelist":
		return containsProject // Include only if project is in whitelist
	default:
		return true // No filtering if mode is not recognized
	}
}

// ShouldIncludeByTags checks if task should be included based on tag filters
func ShouldIncludeByTags(tags []string, filters types.FiltersConfig) bool {
	if len(filters.TagFilterTags) == 0 {
		return true // No tag filter configured
	}

	// Convert all tags to lowercase for comparison
	tagsLower := make([]string, len(tags))
	for i, tag := range tags {
		tagsLower[i] = strings.ToLower(tag)
	}

	filterTagsLower := make([]string, len(filters.TagFilterTags))
	for i, tag := range filters.TagFilterTags {
		filterTagsLower[i] = strings.ToLower(tag)
	}

	switch strings.ToLower(filters.TagFilterMode) {
	case "blacklist":
		// Exclude if task has any blacklisted tag
		for _, taskTag := range tagsLower {
			for _, filterTag := range filterTagsLower {
				if taskTag == filterTag {
					return false
				}
			}
		}
		return true
	case "whitelist":
		// Include only if task has at least one whitelisted tag
		for _, taskTag := range tagsLower {
			for _, filterTag := range filterTagsLower {
				if taskTag == filterTag {
					return true
				}
			}
		}
		return false
	default:
		return true // No filtering if mode is not recognized
	}
}

// FilterProjects filters a list of project names based on project filter configuration
func FilterProjects(projects []string, cfg *types.Config) []string {
	if cfg == nil || len(cfg.Filters.ProjectFilterProjects) == 0 {
		return projects
	}

	var filtered []string
	for _, project := range projects {
		if ShouldIncludeByProject(project, cfg.Filters) {
			filtered = append(filtered, project)
		}
	}
	return filtered
}

// FilterTags filters a map of tag names and counts based on tag filter configuration
func FilterTags(tags map[string]int, cfg *types.Config) map[string]int {
	if cfg == nil || len(cfg.Filters.TagFilterTags) == 0 {
		return tags
	}

	filtered := make(map[string]int)
	for tag, count := range tags {
		if ShouldIncludeByTag(tag, cfg.Filters) {
			filtered[tag] = count
		}
	}
	return filtered
}

// ShouldIncludeByTag checks if a single tag should be included based on tag filters
func ShouldIncludeByTag(tag string, filters types.FiltersConfig) bool {
	if len(filters.TagFilterTags) == 0 {
		return true // No tag filter configured
	}

	tagLower := strings.ToLower(tag)
	containsTag := false

	for _, filterTag := range filters.TagFilterTags {
		if strings.ToLower(filterTag) == tagLower {
			containsTag = true
			break
		}
	}

	switch strings.ToLower(filters.TagFilterMode) {
	case "blacklist":
		return !containsTag // Exclude if tag is in blacklist
	case "whitelist":
		return containsTag // Include only if tag is in whitelist
	default:
		return true // No filtering if mode is not recognized
	}
}

