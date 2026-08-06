package gxconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

type Config struct {
	User User `json:"user"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c Config) HasIdentity() bool {
	return strings.TrimSpace(c.User.Name) != "" && strings.TrimSpace(c.User.Email) != ""
}

func DefaultPath() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func DisplayPath() string {
	if override := os.Getenv("GX_HOME"); strings.TrimSpace(override) != "" {
		return filepath.Join(override, "config.json")
	}
	return "~/.gx/config.json"
}

func Load() (Config, error) {
	path, err := DefaultPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("read gx config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse gx config: %w", err)
	}
	return cfg, nil
}

func LoadAt(root string) (Config, error) {
	path, err := DefaultPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("read gx config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse gx config: %w", err)
	}
	return cfg, nil
}

func Save(cfg Config) error {
	return SaveAt("", cfg)
}

func SaveAt(root string, cfg Config) error {
	path, err := DefaultPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create gx config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal gx config: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write gx config: %w", err)
	}
	return nil
}
