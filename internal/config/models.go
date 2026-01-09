package config

import "time"

// Workspace represents a directory containing git repositories
type Workspace struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	MaxDepth int    `json:"max_depth"`
}

// Repository represents a tracked git repository
type Repository struct {
	Path           string    `json:"path"`
	CurrentBranch  string    `json:"current_branch"`
	LocalBranches  []string  `json:"local_branches"`
	RemoteBranches []string  `json:"remote_branches"`
	LastBranchScan time.Time `json:"last_branch_scan"`
	LastUpdated    time.Time `json:"last_updated"`
}

// Config holds the complete application configuration
type Config struct {
	Workspaces []Workspace            `json:"workspaces"`
	Repos      map[string]*Repository `json:"repos"`
}

// NewConfig creates a new empty configuration
func NewConfig() *Config {
	return &Config{
		Workspaces: make([]Workspace, 0),
		Repos:      make(map[string]*Repository),
	}
}
