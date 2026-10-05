package scope

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/salttis/tasx/internal/config"
)

var invalidProjectName = regexp.MustCompile(`[^a-z0-9_-]+`)

type Scope struct {
	Name        string
	Root        string
	TasksPath   string
	RoadmapPath string
	Missing     bool
}

func List(cfg config.Config) ([]Scope, error) {
	scopes := []Scope{personal(cfg)}
	registry := filepath.Join(cfg.PersonalDir, "ai", "projects")
	data, err := os.ReadFile(registry)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read project registry %q: %w", registry, err)
	}
	seen := map[string]bool{filepath.Clean(scopes[0].TasksPath): true}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, root, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(name) == "" || strings.TrimSpace(root) == "" {
			continue
		}
		root = strings.TrimSpace(root)
		tasksPath := filepath.Join(root, ".ai", "tasks")
		key := filepath.Clean(tasksPath)
		if seen[key] {
			continue
		}
		seen[key] = true
		scopes = append(scopes, Scope{
			Name:        strings.ToLower(strings.TrimSpace(name)),
			Root:        root,
			TasksPath:   tasksPath,
			RoadmapPath: filepath.Join(root, ".ai", "roadmap.md"),
			Missing:     !isFile(tasksPath),
		})
	}
	return scopes, nil
}

func Context(cfg config.Config) (Scope, error) {
	if cfg.DefaultScope == "global" {
		return personal(cfg), nil
	}
	root, err := gitRoot()
	if err != nil {
		return Scope{}, err
	}
	if root == "" {
		if cfg.DefaultScope == "repo" {
			return Scope{}, fmt.Errorf("current directory is not inside a Git repository")
		}
		return personal(cfg), nil
	}
	scopes, err := List(cfg)
	if err != nil {
		return Scope{}, err
	}
	for _, item := range scopes[1:] {
		if samePath(item.Root, root) {
			if !item.Missing {
				return item, nil
			}
			break
		}
	}
	repo := Scope{
		Name:        strings.ToLower(filepath.Base(root)),
		Root:        root,
		TasksPath:   filepath.Join(root, ".ai", "tasks"),
		RoadmapPath: filepath.Join(root, ".ai", "roadmap.md"),
		Missing:     !isFile(filepath.Join(root, ".ai", "tasks")),
	}
	if !repo.Missing {
		return repo, nil
	}
	if cfg.DefaultScope == "repo" {
		return Scope{}, fmt.Errorf("repository task list does not exist: %s", repo.TasksPath)
	}
	return personal(cfg), nil
}

func ResolveProject(cfg config.Config, name string) (Scope, error) {
	scopes, err := List(cfg)
	if err != nil {
		return Scope{}, err
	}
	var found *Scope
	for _, item := range scopes[1:] {
		base := strings.ToLower(filepath.Base(item.Root))
		if ProjectKey(item.Name) != ProjectKey(name) && ProjectKey(base) != ProjectKey(name) {
			continue
		}
		if found != nil && !samePath(found.Root, item.Root) {
			return Scope{}, fmt.Errorf("project name %q is ambiguous in the project registry", name)
		}
		copy := item
		found = &copy
	}
	if found == nil {
		return Scope{}, fmt.Errorf("project %q was not found in the project registry", name)
	}
	return *found, nil
}

func ProjectKey(name string) string {
	return strings.Trim(invalidProjectName.ReplaceAllString(strings.ToLower(name), "-"), "-_")
}

func Current(cfg config.Config, forceGlobal, forceRepo bool) (Scope, error) {
	if forceGlobal || (!forceRepo && cfg.DefaultScope == "global") {
		return personal(cfg), nil
	}
	root, err := gitRoot()
	if err != nil {
		return Scope{}, err
	}
	if root == "" {
		if forceRepo || cfg.DefaultScope == "repo" {
			return Scope{}, fmt.Errorf("current directory is not inside a Git repository")
		}
		return personal(cfg), nil
	}
	repo := Scope{
		Name:        strings.ToLower(filepath.Base(root)),
		Root:        root,
		TasksPath:   filepath.Join(root, ".ai", "tasks"),
		RoadmapPath: filepath.Join(root, ".ai", "roadmap.md"),
		Missing:     !isFile(filepath.Join(root, ".ai", "tasks")),
	}
	if forceRepo || cfg.DefaultScope == "repo" {
		if repo.Missing {
			return Scope{}, fmt.Errorf("repository task list does not exist: %s", repo.TasksPath)
		}
		return repo, nil
	}
	if !repo.Missing {
		return repo, nil
	}
	if isFile(personal(cfg).TasksPath) {
		return personal(cfg), nil
	}
	return repo, nil
}

func samePath(left, right string) bool {
	if os.PathSeparator == '\\' {
		return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
	}
	return filepath.Clean(left) == filepath.Clean(right)
}

func Repository() (Scope, error) {
	root, err := gitRoot()
	if err != nil {
		return Scope{}, err
	}
	if root == "" {
		return Scope{}, fmt.Errorf("current directory is not inside a Git repository")
	}
	tasksPath := filepath.Join(root, ".ai", "tasks")
	return Scope{
		Name:        strings.ToLower(filepath.Base(root)),
		Root:        root,
		TasksPath:   tasksPath,
		RoadmapPath: filepath.Join(root, ".ai", "roadmap.md"),
		Missing:     !isFile(tasksPath),
	}, nil
}

func personal(cfg config.Config) Scope {
	return Scope{
		Name:        "personal",
		Root:        cfg.PersonalDir,
		TasksPath:   filepath.Join(cfg.PersonalDir, "tasks"),
		RoadmapPath: filepath.Join(cfg.PersonalDir, "roadmap.md"),
		Missing:     !isFile(filepath.Join(cfg.PersonalDir, "tasks")),
	}
}

func gitRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("find current directory: %w", err)
	}
	cmd := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return "", nil
		}
		return "", fmt.Errorf("find Git repository: %w", err)
	}
	return filepath.Clean(strings.TrimSpace(string(output))), nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
