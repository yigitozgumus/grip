package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yigitozgumus/grip/internal/config"
	"github.com/yigitozgumus/grip/internal/discovery"
	"github.com/yigitozgumus/grip/internal/git"
	"github.com/yigitozgumus/grip/internal/ui"
	"github.com/yigitozgumus/grip/internal/watcher"
)

//go:embed shell/zsh.sh
var zshScript string

//go:embed shell/bash.sh
var bashScript string

//go:embed shell/fish.fish
var fishScript string

var cfgManager *config.Manager

func main() {
	var err error
	cfgManager, err = config.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "grip",
	Short: "Git Workspace Manager - Fast project & branch switching",
	Long: `grip helps you navigate between git repositories and branches quickly.

Add workspaces, scan for repos, then use fuzzy search to jump between projects.`,
}

var addCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Add a workspace directory to track",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		depth, _ := cmd.Flags().GetInt("depth")

		if name == "" {
			name = filepath.Base(path)
		}

		if err := cfgManager.AddWorkspace(name, path, depth); err != nil {
			return err
		}

		if err := cfgManager.Save(); err != nil {
			return err
		}

		color.Green("Added workspace: %s (%s)", name, path)
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scan all workspaces and index repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := cfgManager.GetConfig()

		if len(conf.Workspaces) == 0 {
			return fmt.Errorf("no workspaces configured. Use 'grip add <path>' first")
		}

		cfgManager.ClearRepos()

		for _, ws := range conf.Workspaces {
			fmt.Printf("Scanning %s (%s)...\n", ws.Name, ws.Path)

			repos, err := discovery.ScanWorkspace(ws)
			if err != nil {
				color.Yellow("Warning: %v", err)
				continue
			}

			for _, repo := range repos {
				cfgManager.UpdateRepo(repo.Path, repo)
				fmt.Printf("  %s (%s)\n", filepath.Base(repo.Path), repo.CurrentBranch)
			}
		}

		if err := cfgManager.Save(); err != nil {
			return err
		}

		conf = cfgManager.GetConfig()
		color.Green("Indexed %d repositories", len(conf.Repos))
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show all tracked repositories",
	Run: func(cmd *cobra.Command, args []string) {
		conf := cfgManager.GetConfig()

		if len(conf.Repos) == 0 {
			fmt.Println("No repositories tracked. Run 'grip init' first.")
			return
		}

		green := color.New(color.FgGreen).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		for path, repo := range conf.Repos {
			name := filepath.Base(path)
			fmt.Printf("%-30s %s  %s\n",
				green(name),
				cyan(repo.CurrentBranch),
				path,
			)
		}
	},
}

var cdCmd = &cobra.Command{
	Use:   "cd [project]",
	Short: "Select a project and print its path (use with shell wrapper)",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := cfgManager.GetConfig()

		if len(conf.Repos) == 0 {
			return fmt.Errorf("no repositories tracked. Run 'grip init' first")
		}

		var selected *ui.ProjectItem

		if len(args) > 0 {
			// Direct project name match
			projectName := args[0]
			for path, repo := range conf.Repos {
				if filepath.Base(path) == projectName {
					selected = &ui.ProjectItem{Path: path, Repo: repo}
					break
				}
			}
			if selected == nil {
				return fmt.Errorf("project not found: %s", projectName)
			}
		} else {
			// Fuzzy select
			var err error
			selected, err = ui.SelectProject(conf.Repos)
			if err != nil {
				return err
			}
		}

		// Output path only (for shell wrapper to cd)
		fmt.Println(selected.Path)
		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch [project]",
	Short: "Select a project and branch, then switch",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := cfgManager.GetConfig()

		if len(conf.Repos) == 0 {
			return fmt.Errorf("no repositories tracked. Run 'grip init' first")
		}

		var selected *ui.ProjectItem

		if len(args) > 0 {
			projectName := args[0]
			for path, repo := range conf.Repos {
				if filepath.Base(path) == projectName {
					selected = &ui.ProjectItem{Path: path, Repo: repo}
					break
				}
			}
			if selected == nil {
				return fmt.Errorf("project not found: %s", projectName)
			}
		} else {
			var err error
			selected, err = ui.SelectProject(conf.Repos)
			if err != nil {
				return err
			}
		}

		// Select branch
		branch, err := ui.SelectBranch(selected.Repo)
		if err != nil {
			return err
		}

		// Switch branch
		if err := git.SwitchBranch(selected.Path, branch); err != nil {
			color.Red("Failed to switch branch: %v", err)
			return err
		}

		color.Green("Switched to %s", branch)

		// Update config
		selected.Repo.CurrentBranch = branch
		selected.Repo.LastUpdated = time.Now()
		cfgManager.UpdateRepo(selected.Path, selected.Repo)
		cfgManager.Save()

		// Output path for shell wrapper
		fmt.Println(selected.Path)
		return nil
	},
}

var branchesCmd = &cobra.Command{
	Use:   "branches [project]",
	Short: "Switch branch in current or specified project",
	RunE: func(cmd *cobra.Command, args []string) error {
		refresh, _ := cmd.Flags().GetBool("refresh")
		conf := cfgManager.GetConfig()

		var repo *config.Repository
		var repoPath string

		if len(args) > 0 {
			// Find by project name
			projectName := args[0]
			for path, r := range conf.Repos {
				if filepath.Base(path) == projectName {
					repo = r
					repoPath = path
					break
				}
			}
			if repo == nil {
				return fmt.Errorf("project not found: %s", projectName)
			}
		} else {
			// Use current directory
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			repo = conf.Repos[cwd]
			repoPath = cwd

			if repo == nil {
				// Try to find repo containing current dir
				for path, r := range conf.Repos {
					if cwd == path || isSubPath(path, cwd) {
						repo = r
						repoPath = path
						break
					}
				}
			}

			if repo == nil {
				return fmt.Errorf("current directory is not a tracked repository")
			}
		}

		// Refresh branch cache if requested
		if refresh {
			fmt.Println("Refreshing branch list...")
			r, err := git.RefreshRepoInfo(repoPath)
			if err != nil {
				return err
			}
			repo = r
			cfgManager.UpdateRepo(repoPath, repo)
			cfgManager.Save()
		}

		// Select branch
		branch, err := ui.SelectBranch(repo)
		if err != nil {
			return err
		}

		// Switch
		if err := git.SwitchBranch(repoPath, branch); err != nil {
			color.Red("Failed to switch branch: %v", err)
			return err
		}

		color.Green("Switched to %s", branch)

		// Update config
		repo.CurrentBranch = branch
		repo.LastUpdated = time.Now()
		cfgManager.UpdateRepo(repoPath, repo)
		cfgManager.Save()

		return nil
	},
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start daemon to watch repositories for branch changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := watcher.New(cfgManager)
		if err != nil {
			return err
		}

		if err := w.Start(); err != nil {
			return err
		}

		color.Green("Watching repositories for changes... (Ctrl+C to stop)")

		// Wait for interrupt
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		fmt.Println("\nStopping watcher...")
		return w.Stop()
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Output shell integration snippet for your rc file",
	Long: `Output the shell integration snippet to add to your shell configuration.

Examples:
  grip shell              # Auto-detect shell
  grip shell --zsh        # Force zsh
  grip shell >> ~/.zshrc  # Append to config`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check explicit flags first
		useZsh, _ := cmd.Flags().GetBool("zsh")
		useBash, _ := cmd.Flags().GetBool("bash")
		useFish, _ := cmd.Flags().GetBool("fish")

		var shell string
		switch {
		case useZsh:
			shell = "zsh"
		case useBash:
			shell = "bash"
		case useFish:
			shell = "fish"
		default:
			// Auto-detect from $SHELL
			shell = detectShell()
		}

		// Output the embedded script content
		switch shell {
		case "fish":
			fmt.Print(fishScript)
		case "bash":
			fmt.Print(bashScript)
		default: // zsh
			fmt.Print(zshScript)
		}

		return nil
	},
}

func detectShell() string {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return "zsh" // Default
	}

	shellName := filepath.Base(shellPath)
	switch shellName {
	case "fish":
		return "fish"
	case "bash":
		return "bash"
	default:
		return "zsh"
	}
}

func isSubPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return !filepath.IsAbs(rel) && rel != ".." && rel[:2] != ".."
}

func init() {
	addCmd.Flags().StringP("name", "n", "", "Workspace name (defaults to directory name)")
	addCmd.Flags().IntP("depth", "d", 2, "Maximum directory depth to scan")

	branchesCmd.Flags().Bool("refresh", false, "Refresh branch cache from git")

	shellCmd.Flags().Bool("zsh", false, "Output zsh integration")
	shellCmd.Flags().Bool("bash", false, "Output bash integration")
	shellCmd.Flags().Bool("fish", false, "Output fish integration")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(branchesCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(shellCmd)
}
