package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

const (
	configDir  = ".config/grip"
	configFile = "config.json"
)

// Manager handles configuration persistence and thread-safe access
type Manager struct {
	mu     sync.RWMutex
	config *Config
	path   string
}

// NewManager creates a new config manager and loads existing config
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, configDir, configFile)

	m := &Manager{
		config: NewConfig(),
		path:   configPath,
	}

	// Create config dir if not exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, err
	}

	// Load existing config (ignore if doesn't exist)
	if err := m.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return m, nil
}

// Load reads the config from disk
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, m.config)
}

// Save writes the config to disk
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

// AddWorkspace adds a new workspace to track
func (m *Manager) AddWorkspace(name, path string, maxDepth int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate
	for _, ws := range m.config.Workspaces {
		if ws.Path == path {
			return nil // Already exists
		}
	}

	m.config.Workspaces = append(m.config.Workspaces, Workspace{
		Name:     name,
		Path:     path,
		MaxDepth: maxDepth,
	})

	return nil
}

// RemoveWorkspace removes a workspace by path
func (m *Manager) RemoveWorkspace(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	workspaces := make([]Workspace, 0, len(m.config.Workspaces))
	for _, ws := range m.config.Workspaces {
		if ws.Path != path {
			workspaces = append(workspaces, ws)
		}
	}
	m.config.Workspaces = workspaces
}

// GetConfig returns a copy of the current config
func (m *Manager) GetConfig() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// UpdateRepo updates or adds a repository
func (m *Manager) UpdateRepo(path string, repo *Repository) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Repos[path] = repo
}

// GetRepo retrieves a repository by path
func (m *Manager) GetRepo(path string) *Repository {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config.Repos[path]
}

// ClearRepos removes all tracked repositories
func (m *Manager) ClearRepos() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Repos = make(map[string]*Repository)
}
