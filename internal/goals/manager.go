package goals

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/taskvanguard/taskvanguard/internal/taskwarrior"
	"github.com/taskvanguard/taskvanguard/pkg/types"
)

type Manager struct {
	client *taskwarrior.Client
	config *types.Config
}

func NewManager(config *types.Config) *Manager {
	return &Manager{
		client: taskwarrior.NewClient(),
		config: config,
	}
}

// ListGoals returns all goals (tasks in project:goals)
func (m *Manager) ListGoals() ([]types.Task, error) {
	return m.client.GetGoals()
}

// AddGoal creates a new goal with the given arguments
func (m *Manager) AddGoal(args []string) (string, int, error) {
	// Prepend project:<goal_project_name> to the arguments
	goalProjectArg := fmt.Sprintf("project:%s", m.config.Settings.GoalProjectName)
	goalsArgs := append([]string{goalProjectArg}, args...)
	return m.client.AddTaskToTaskWarrior(goalsArgs)
}

// ModifyGoal modifies an existing goal
func (m *Manager) ModifyGoal(goalID int, args []string) (string, error) {
	return m.client.ModifyTaskInTaskWarrior(goalID, args)
}

// DeleteGoal deletes a goal by ID
func (m *Manager) DeleteGoal(goalID string) error {
	cmd := exec.Command("task", goalID, "delete")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete goal %s: %v\nOutput: %s", goalID, err, string(output))
	}
	return nil
}

// ShowGoal shows details of a goal or task by ID
func (m *Manager) ShowGoal(id string) (*types.Task, error) {
	return m.client.GetTaskByID(id)
}

// LinkTaskToGoal links a task to a goal using the goal UDA
func (m *Manager) LinkTaskToGoal(taskID, goalUUID string) error {
	cmd := exec.Command("task", taskID, "modify", fmt.Sprintf("goal:%s", goalUUID))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to link task %s to goal %s: %v\nOutput: %s", taskID, goalUUID, err, string(output))
	}
	return nil
}

// UnlinkTaskFromGoal removes the goal link from a task
func (m *Manager) UnlinkTaskFromGoal(taskID string) error {
	cmd := exec.Command("task", taskID, "modify", "goal:")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unlink task %s from goal: %v\nOutput: %s", taskID, err, string(output))
	}
	return nil
}

// Link determines which ID is a task and which is a goal, then links them
func (m *Manager) Link(id1, id2 string) error {
	task1, err := m.client.GetTaskByID(id1)
	if err != nil {
		return fmt.Errorf("failed to get task %s: %v", id1, err)
	}
	
	task2, err := m.client.GetTaskByID(id2)
	if err != nil {
		return fmt.Errorf("failed to get task %s: %v", id2, err)
	}

	if task1 == nil || task2 == nil {
		return errors.New("one or both IDs not found")
	}

	// Determine which is the goal and which is the task
	var taskID, goalUUID string

	if task1.Project == m.config.Settings.GoalProjectName && task2.Project != m.config.Settings.GoalProjectName {
		taskID = id2
		goalUUID = task1.UUID
	} else if task2.Project == m.config.Settings.GoalProjectName && task1.Project != m.config.Settings.GoalProjectName {
		taskID = id1
		goalUUID = task2.UUID
	} else if task1.Project == m.config.Settings.GoalProjectName && task2.Project == m.config.Settings.GoalProjectName {
		return errors.New("both items are goals - cannot link two goals together")
	} else {
		return errors.New("both items are tasks - cannot link two tasks together")
	}

	return m.LinkTaskToGoal(taskID, goalUUID)
}

// Unlink removes the link between a task and goal (order-agnostic)
func (m *Manager) Unlink(id1, id2 string) error {
	task1, err := m.client.GetTaskByID(id1)
	if err != nil {
		return fmt.Errorf("failed to get task %s: %v", id1, err)
	}
	
	task2, err := m.client.GetTaskByID(id2)
	if err != nil {
		return fmt.Errorf("failed to get task %s: %v", id2, err)
	}

	if task1 == nil || task2 == nil {
		return errors.New("one or both IDs not found")
	}

	// Determine which is the task (non-goal)
	var taskID string

	if task1.Project == m.config.Settings.GoalProjectName && task2.Project != m.config.Settings.GoalProjectName {
		taskID = id2
	} else if task2.Project == m.config.Settings.GoalProjectName && task1.Project != m.config.Settings.GoalProjectName {
		taskID = id1
	} else {
		return errors.New("cannot determine which item is the task to unlink")
	}

	return m.UnlinkTaskFromGoal(taskID)
}

// GetLinkedTasks returns all tasks linked to a specific goal
func (m *Manager) GetLinkedTasks(goalUUID string) ([]types.Task, error) {
	return m.client.GetTasksWithFilter([]string{fmt.Sprintf("goal:%s", goalUUID)})
}

// GetLinkedGoal returns the goal linked to a specific task
func (m *Manager) GetLinkedGoal(taskID string) (*types.Task, error) {
	task, err := m.client.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	
	if task == nil {
		return nil, errors.New("task not found")
	}

	// Extract goal UUID from task (would need to check UDA fields)
	// This would require extending the Task struct to include UDA fields
	// For now, we'll implement basic functionality
	
	return nil, errors.New("getting linked goal not yet implemented - requires UDA field support")
}

// ShowLinks shows all links for a given ID (goal or task)
func (m *Manager) ShowLinks(id string) ([]types.Task, error) {
	task, err := m.client.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	
	if task == nil {
		return nil, errors.New("task/goal not found")
	}

	if task.Project == m.config.Settings.GoalProjectName {
		// It's a goal, show linked tasks
		return m.GetLinkedTasks(task.UUID)
	} else {
		// It's a task, show linked goal (return as slice for consistency)
		linkedGoal, err := m.GetLinkedGoal(id)
		if err != nil {
			return nil, err
		}
		if linkedGoal == nil {
			return []types.Task{}, nil
		}
		return []types.Task{*linkedGoal}, nil
	}
}

// ValidateGoalID checks if the given ID refers to a goal
func (m *Manager) ValidateGoalID(id string) error {
	task, err := m.client.GetTaskByID(id)
	if err != nil {
		return err
	}
	
	if task == nil {
		return errors.New("goal not found")
	}
	
	if task.Project != m.config.Settings.GoalProjectName {
		return errors.New("ID does not refer to a goal")
	}
	
	return nil
}

// IsGoal checks if a given ID refers to a goal
func (m *Manager) IsGoal(id string) (bool, error) {
	task, err := m.client.GetTaskByID(id)
	if err != nil {
		return false, err
	}
	
	if task == nil {
		return false, nil
	}
	
	return task.Project == m.config.Settings.GoalProjectName, nil
}