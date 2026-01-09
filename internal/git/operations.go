package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/yigitozgumus/grip/internal/config"
)

// IsGitRepo checks if a directory is a git repository
func IsGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetLocalBranches returns all local branch names
func GetLocalBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	return filterEmpty(branches), nil
}

// GetRemoteBranches returns all remote branch names
func GetRemoteBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "-r", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	return filterEmpty(branches), nil
}

// SwitchBranch switches to the specified branch
func SwitchBranch(repoPath, branch string) error {
	cmd := exec.Command("git", "-C", repoPath, "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RefreshRepoInfo gets complete repository information
func RefreshRepoInfo(repoPath string) (*config.Repository, error) {
	currentBranch, err := GetCurrentBranch(repoPath)
	if err != nil {
		return nil, err
	}

	localBranches, err := GetLocalBranches(repoPath)
	if err != nil {
		return nil, err
	}

	remoteBranches, err := GetRemoteBranches(repoPath)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &config.Repository{
		Path:           repoPath,
		CurrentBranch:  currentBranch,
		LocalBranches:  localBranches,
		RemoteBranches: remoteBranches,
		LastBranchScan: now,
		LastUpdated:    now,
	}, nil
}

// FetchRemote fetches updates from remote
func FetchRemote(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "fetch", "--prune")
	return cmd.Run()
}

func filterEmpty(strs []string) []string {
	result := make([]string, 0, len(strs))
	for _, s := range strs {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
