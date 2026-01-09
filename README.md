# grip

Fast project and branch switching for git workspaces.

![Demo](https://raw.githubusercontent.com/yigitozgumus/grip/main/demo.gif)

## Features

- **Fuzzy project search** - Quickly find and jump to any git repository
- **Branch switching** - Switch branches with fuzzy search across local and remote
- **Workspace scanning** - Automatically discover repos in your project directories
- **Background watcher** - Tracks branch changes from any tool (IDE, CLI, etc.)
- **Shell integration** - Native functions for bash, zsh, and fish

## Installation

### Using Go

```bash
go install github.com/yigitozgumus/grip/cmd/grip@latest
```

### From Source

```bash
git clone https://github.com/yigitozgumus/grip.git
cd grip
make build
make install  # Copies to ~/.local/bin/
```

### Homebrew (coming soon)

```bash
brew install yigitozgumus/tap/grip
```

## Setup

### 1. Add a workspace

```bash
grip add ~/projects --name work --depth 2
```

### 2. Scan for repositories

```bash
grip init
```

### 3. Add shell integration

```bash
# For zsh
grip shell >> ~/.zshrc
source ~/.zshrc

# For bash
grip shell --bash >> ~/.bashrc

# For fish
grip shell --fish >> ~/.config/fish/config.fish
```

## Usage

### Shell Commands

| Command | Description |
|---------|-------------|
| `gpcd` | Fuzzy select a project and cd into it |
| `gpcd myproject` | Jump directly to "myproject" |
| `gpsw` | Select project + branch, switch and cd |
| `gpb` | Switch branch in current repository |
| `gps` | Show all tracked repositories |
| `gpi` | Re-scan all workspaces |

### CLI Commands

```bash
grip add <path>           # Add a workspace to track
grip init                 # Scan and index all repositories
grip status               # List all tracked repos
grip cd [project]         # Select project, print path
grip switch [project]     # Select project + branch
grip branches [--refresh] # Switch branch in current repo
grip watch                # Start background watcher daemon
grip shell                # Output shell integration
```

## Configuration

Config is stored at `~/.config/grip/config.json`

```json
{
  "workspaces": [
    {
      "name": "work",
      "path": "/Users/you/projects",
      "max_depth": 2
    }
  ],
  "repos": {
    "/Users/you/projects/myapp": {
      "current_branch": "main",
      "local_branches": ["main", "feature/x"],
      "remote_branches": ["origin/main"]
    }
  }
}
```

## How It Works

1. **Workspaces**: Directories containing git repositories (scanned up to `max_depth`)
2. **Indexing**: `grip init` discovers repos and caches branch information
3. **Fuzzy Search**: Interactive selection powered by [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder)
4. **Shell Integration**: Functions wrap `grip` to enable actual directory changes
5. **Watcher**: Optional daemon monitors `.git/HEAD` files for real-time branch updates

## License

MIT License - see [LICENSE](LICENSE)
