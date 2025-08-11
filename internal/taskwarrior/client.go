package taskwarrior

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/taskvanguard/taskvanguard/pkg/filter"
	"github.com/taskvanguard/taskvanguard/pkg/types"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) AddTaskToTaskWarrior(args []string) (string, int, error) {
	cmdArgs := append([]string{"add"}, args...)
	cmd := exec.Command("task", cmdArgs...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "TaskWarrior not found or failed. Task creation failed: " +
			err.Error() + "\nOutput: " + string(output)
		return string(output), 0, errors.New(msg)
	}

	// Extract task ID from output using regex
	// TaskWarrior typically outputs: "Created task 123."
	re := regexp.MustCompile(`Created task (\d+)\.`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		msg := "could not extract task ID from output: " + string(output)
		return string(output), 0, errors.New(msg)
	}
	
	taskID, err := strconv.Atoi(matches[1])
	if err != nil {
		msg := fmt.Sprintf("could not parse task ID: %v\nOutput: %s", err, string(output))
		return string(output), 0, errors.New(msg)
	}
	return string(output), taskID, nil
}

func (c *Client) ModifyTaskInTaskWarrior(taskId int, args []string) (string, error) {
	cmdArgs := append([]string{"modify", strconv.Itoa(taskId)}, args...)

	cmd := exec.Command("task", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "TaskWarrior not found or failed. Task " +
			strconv.Itoa(taskId) +
			" modification failed: " + err.Error() +
			"\nOutput: " + string(output)
		return "", errors.New(msg)
	}

	return string(output), nil
}

func (c *Client) AddSingleAnnotation(taskId string, value string) error {
	cmd := exec.Command("task", taskId, "annotate", value)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "failed to add annotation to task " + taskId + ": " + err.Error() + "Output: " + string(output)
		return errors.New(msg)
	}
	
	return nil
}

func (c *Client) GetTasks() ([]types.Task, error) {
	cmd := exec.Command("task", "export")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []types.Task
	if err := json.Unmarshal(output, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *Client) GetPendingTasks() ([]types.Task, error) {
	cmd := exec.Command("task", "status:pending", "export")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []types.Task
	if err := json.Unmarshal(output, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTasksFiltered returns all tasks with filtering applied
func (c *Client) GetTasksFiltered(cfg *types.Config) ([]types.Task, error) {
	tasks, err := c.GetTasks()
	if err != nil {
		return nil, err
	}
	return filter.FilterTasks(tasks, cfg), nil
}

// GetPendingTasksFiltered returns pending tasks with filtering applied
func (c *Client) GetPendingTasksFiltered(cfg *types.Config) ([]types.Task, error) {
	tasks, err := c.GetPendingTasks()
	if err != nil {
		return nil, err
	}
	return filter.FilterTasks(tasks, cfg), nil
}

func (c *Client) GetTaskByID(id string) (*types.Task, error) {
	cmd := exec.Command("task", id, "export")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []types.Task
	if err := json.Unmarshal(output, &tasks); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	return &tasks[0], nil
}

func (c *Client) IsAvailable() bool {
	_, err := exec.LookPath("task")
	return err == nil
}

func (c *Client) GetVersion() (string, error) {
	cmd := exec.Command("task", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(output))
	return version, nil
}

func (c *Client) GetTags() (map[string]int, error) {
	cmd := exec.Command("task",
		"rc.verbose=nothing",
		"rc.report.tagscounter.columns=tag,count",
		"tags",
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	tagCounts := make(map[string]int)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			tag := fields[0]
			count, err := strconv.Atoi(fields[1])
			if err != nil {
				continue // skip malformed line
			}
			tagCounts[tag] = count
		}
	}

	return tagCounts, nil
}

func (c *Client) GetGoals() ([]types.Task, error) {
	cmd := exec.Command("task", "project:goals", "export")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []types.Task
	if err := json.Unmarshal(output, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetGoalsFiltered returns goal tasks with filtering applied
func (c *Client) GetGoalsFiltered(cfg *types.Config) ([]types.Task, error) {
	tasks, err := c.GetGoals()
	if err != nil {
		return nil, err
	}
	return filter.FilterTasks(tasks, cfg), nil
}

func (c *Client) GetProjects() ([]string, error) {
	cmd := exec.Command("task", "projects")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	projects := []string{}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip header lines and empty lines
		if line == "" || strings.Contains(line, "Project") || strings.Contains(line, "----") {
			continue
		}
		
		// Extract project name from the line (first column)
		fields := strings.Fields(line)
		if len(fields) > 0 {
			project := fields[0]
			projects = append(projects, project)
		}
	}

	return projects, nil
}

// GetProjectsFiltered returns projects with filtering applied
func (c *Client) GetProjectsFiltered(cfg *types.Config) ([]string, error) {
	projects, err := c.GetProjects()
	if err != nil {
		return nil, err
	}
	return filter.FilterProjects(projects, cfg), nil
}

// GetTagsFiltered returns tags with filtering applied
func (c *Client) GetTagsFiltered(cfg *types.Config) (map[string]int, error) {
	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}
	return filter.FilterTags(tags, cfg), nil
}

func (c *Client) GetPendingTasksWithArgs(filterArgs []string) ([]types.Task, error) {
	args := append([]string{"task"}, "status:pending")
	args = append(args, filterArgs...)
	args = append(args, "export")

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []types.Task
	if err := json.Unmarshal(output, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTasksWithFilter returns tasks with custom filter arguments
func (c *Client) GetTasksWithFilter(filterArgs []string) ([]types.Task, error) {
	args := append([]string{"task"}, filterArgs...)
	args = append(args, "export")
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []types.Task
	if err := json.Unmarshal(output, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTasksFiltered returns all tasks with filtering applied
func (c *Client) GetTasksWithFilterFiltered(cfg *types.Config, filterArgs []string) ([]types.Task, error) {
	tasks, err := c.GetTasksWithFilter(filterArgs)
	if err != nil {
		return nil, err
	}
	return filter.FilterTasks(tasks, cfg), nil
}

// GetTasksFiltered returns all tasks with filtering applied
func (c *Client) GetPendingTasksWithArgsFiltered(cfg *types.Config, filterArgs []string) ([]types.Task, error) {
	tasks, err := c.GetPendingTasksWithArgs(filterArgs)
	if err != nil {
		return nil, err
	}
	return filter.FilterTasks(tasks, cfg), nil
}

func (c *Client) StartTask(taskId string) error {
	cmd := exec.Command("task", taskId, "start")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "failed to start task " + taskId + ": " + err.Error() + "Output: " + string(output)
		return errors.New(msg)
	}
	
	return nil
}