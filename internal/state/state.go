package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type State struct {
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	LastUpdate  time.Time              `json:"last_update"`
	Resources   map[string]interface{} `json:"resources"`
	Outputs     map[string]interface{} `json:"outputs"`
}

type StateManager struct {
	statePath string
	lockPath  string
	mu        sync.Mutex
}

func NewStateManager(baseDir string) *StateManager {
	return &StateManager{
		statePath: filepath.Join(baseDir, "states"),
		lockPath:  filepath.Join(baseDir, "locks"),
	}
}

func (sm *StateManager) Lock(env string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	lockFile := filepath.Join(sm.lockPath, fmt.Sprintf("%s.lock", env))
	if err := os.MkdirAll(sm.lockPath, 0755); err != nil {
		return fmt.Errorf("failed to create lock directory: %w", err)
	}

	if _, err := os.Stat(lockFile); err == nil {
		return fmt.Errorf("environment %s is locked", env)
	}

	lockInfo := struct {
		Time    time.Time `json:"time"`
		Process int       `json:"process"`
	}{
		Time:    time.Now(),
		Process: os.Getpid(),
	}

	data, err := json.Marshal(lockInfo)
	if err != nil {
		return fmt.Errorf("failed to create lock info: %w", err)
	}

	return os.WriteFile(lockFile, data, 0644)
}

func (sm *StateManager) Unlock(env string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	lockFile := filepath.Join(sm.lockPath, fmt.Sprintf("%s.lock", env))
	return os.Remove(lockFile)
}

func (sm *StateManager) SaveState(env string, state *State) error {
	if err := os.MkdirAll(sm.statePath, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	state.LastUpdate = time.Now()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	stateFile := filepath.Join(sm.statePath, fmt.Sprintf("%s.json", env))
	return os.WriteFile(stateFile, data, 0644)
}

func (sm *StateManager) LoadState(env string) (*State, error) {
	stateFile := filepath.Join(sm.statePath, fmt.Sprintf("%s.json", env))
	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{
				Environment: env,
				Resources:   make(map[string]interface{}),
				Outputs:     make(map[string]interface{}),
			}, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}
