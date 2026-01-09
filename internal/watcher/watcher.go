package watcher

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yigitozgumus/grip/internal/config"
	"github.com/yigitozgumus/grip/internal/git"
)

// Watcher monitors git repositories for branch changes
type Watcher struct {
	watcher *fsnotify.Watcher
	cfg     *config.Manager
	mu      sync.Mutex
	repos   map[string]bool // Track watched repos
	done    chan struct{}
}

// New creates a new file watcher
func New(cfg *config.Manager) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: w,
		cfg:     cfg,
		repos:   make(map[string]bool),
		done:    make(chan struct{}),
	}, nil
}

// Start begins watching all tracked repositories
func (w *Watcher) Start() error {
	conf := w.cfg.GetConfig()

	for repoPath := range conf.Repos {
		if err := w.watchRepo(repoPath); err != nil {
			log.Printf("Failed to watch %s: %v", repoPath, err)
		}
	}

	go w.eventLoop()

	return nil
}

func (w *Watcher) watchRepo(repoPath string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.repos[repoPath] {
		return nil // Already watching
	}

	// Watch the .git/HEAD file for branch changes
	headPath := filepath.Join(repoPath, ".git", "HEAD")
	if err := w.watcher.Add(headPath); err != nil {
		return err
	}

	w.repos[repoPath] = true
	log.Printf("Watching: %s", filepath.Base(repoPath))
	return nil
}

func (w *Watcher) eventLoop() {
	for {
		select {
		case <-w.done:
			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				w.handleHeadChange(event.Name)
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (w *Watcher) handleHeadChange(headPath string) {
	// Extract repo path from .git/HEAD path
	// headPath = /path/to/repo/.git/HEAD
	// repoPath = /path/to/repo
	repoPath := filepath.Dir(filepath.Dir(headPath))

	// Get current branch
	branch, err := git.GetCurrentBranch(repoPath)
	if err != nil {
		log.Printf("Failed to get current branch for %s: %v", repoPath, err)
		return
	}

	// Update config
	conf := w.cfg.GetConfig()
	if repo, exists := conf.Repos[repoPath]; exists {
		if repo.CurrentBranch != branch {
			repo.CurrentBranch = branch
			repo.LastUpdated = time.Now()
			w.cfg.UpdateRepo(repoPath, repo)

			if err := w.cfg.Save(); err != nil {
				log.Printf("Failed to save config: %v", err)
			} else {
				log.Printf("Updated %s: %s", filepath.Base(repoPath), branch)
			}
		}
	}
}

// Stop stops the watcher
func (w *Watcher) Stop() error {
	close(w.done)
	return w.watcher.Close()
}

// AddRepo adds a new repository to watch
func (w *Watcher) AddRepo(repoPath string) error {
	return w.watchRepo(repoPath)
}
