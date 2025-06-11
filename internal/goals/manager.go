package goals

// import (
// 	"encoding/json"
// 	"os"
// 	"path/filepath"
// 	"time"

// 	"xarc.dev/taskvanguard/pkg/types"
// )

// type Manager struct {
// 	goalsFile string
// }

// func NewManager() (*Manager, error) {
// 	configDir, err := os.UserConfigDir()
// 	if err != nil {
// 		return nil, err
// 	}

// 	goalsDir := filepath.Join(configDir, "taskvanguard")
// 	if err := os.MkdirAll(goalsDir, 0755); err != nil {
// 		return nil, err
// 	}

// 	return &Manager{
// 		goalsFile: filepath.Join(goalsDir, "goals.json"),
// 	}, nil
// }

// func (m *Manager) LoadGoals() ([]types.Goal, error) {
// 	if _, err := os.Stat(m.goalsFile); os.IsNotExist(err) {
// 		return []types.Goal{}, nil
// 	}

// 	data, err := os.ReadFile(m.goalsFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var goals []types.Goal
// 	if err := json.Unmarshal(data, &goals); err != nil {
// 		return nil, err
// 	}

// 	return goals, nil
// }

// func (m *Manager) SaveGoals(goals []types.Goal) error {
// 	data, err := json.MarshalIndent(goals, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile(m.goalsFile, data, 0644)
// }

// func (m *Manager) AddGoal(title, description string) (*types.Goal, error) {
// 	goals, err := m.LoadGoals()
// 	if err != nil {
// 		return nil, err
// 	}

// 	goal := types.Goal{
// 		ID:          generateID(),
// 		Title:       title,
// 		Description: description,
// 		Priority:    "medium",
// 		Created:     types.TWTime(time.Now()),
// 		Modified:    types.TWTime(time.Now()),
// 		Status:      "active",
// 	}

// 	goals = append(goals, goal)
// 	if err := m.SaveGoals(goals); err != nil {
// 		return nil, err
// 	}

// 	return &goal, nil
// }

// func generateID() string {
// 	return time.Now().Format("20060102150405")
// }