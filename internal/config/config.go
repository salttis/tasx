package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	PersonalDir  string `toml:"personal_dir"`
	DefaultScope string `toml:"default_scope"`
	State        string `toml:"state"`
}

func Load() (Config, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, "", fmt.Errorf("find user home directory: %w", err)
	}

	personalDir := filepath.Join(home, ".personal")
	if value := os.Getenv("TASX_HOME"); value != "" {
		personalDir = value
	}
	cfg := Config{
		PersonalDir:  personalDir,
		DefaultScope: "auto",
		State:        "all",
	}

	path := os.Getenv("TASX_CONFIG")
	if path == "" {
		path = filepath.Join(home, ".tasxrc")
	}
	data, err := os.ReadFile(path)
	if err == nil {
		decoder := toml.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&cfg); err != nil {
			return Config{}, path, fmt.Errorf("read Tasx settings %q: %w", path, err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return Config{}, path, fmt.Errorf("read Tasx settings %q: %w", path, err)
	}

	if value := os.Getenv("TASX_HOME"); value != "" {
		cfg.PersonalDir = value
	} else if cfg.PersonalDir == "" {
		cfg.PersonalDir = filepath.Join(home, ".personal")
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, path, fmt.Errorf("invalid Tasx settings %q: %w", path, err)
	}
	return cfg, path, nil
}

func (c Config) Validate() error {
	if c.DefaultScope != "auto" && c.DefaultScope != "repo" && c.DefaultScope != "global" {
		return fmt.Errorf("default_scope must be auto, repo, or global")
	}
	if c.State != "all" && c.State != "open" && c.State != "done" {
		return fmt.Errorf("state must be all, open, or done")
	}
	return nil
}
