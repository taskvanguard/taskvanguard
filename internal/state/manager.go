package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type TaskContext struct {
	Mood      string    `json:"mood"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}

type StateManager struct {
	statePath string
}

func NewStateManager() (*StateManager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	
	statePath := filepath.Join(configDir, "taskvanguard", "state.json")
	
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return nil, err
	}
	
	return &StateManager{
		statePath: statePath,
	}, nil
}

func (sm *StateManager) SaveContext(context TaskContext) error {
	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(sm.statePath, data, 0644)
}

func (sm *StateManager) LoadContext() (TaskContext, error) {
	context := TaskContext{
		Mood:      "neutral",
		Location:  "unknown",
		Timestamp: time.Now(),
	}
	
	if _, err := os.Stat(sm.statePath); os.IsNotExist(err) {
		return context, nil
	}
	
	data, err := os.ReadFile(sm.statePath)
	if err != nil {
		return context, err
	}
	
	var savedContext TaskContext
	if err := json.Unmarshal(data, &savedContext); err != nil {
		return context, err
	}
	
	if time.Since(savedContext.Timestamp) > 4*time.Hour {
		return TaskContext{
			Mood:      "neutral",
			Location:  "unknown", 
			Timestamp: time.Now(),
		}, nil
	}
	
	return savedContext, nil
}

func (sm *StateManager) IsContextFresh() (bool, error) {
	context, err := sm.LoadContext()
	if err != nil {
		return false, err
	}
	
	return context.Mood != "neutral" || context.Location != "unknown", nil
}