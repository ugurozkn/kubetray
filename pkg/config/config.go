package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ClusterName   string
	DefaultCPUs   int
	DefaultMemory string
	DataDir       string
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dataDir := filepath.Join(home, ".kubetray")

	return &Config{
		ClusterName:   "kubetray",
		DefaultCPUs:   2,
		DefaultMemory: "2G",
		DataDir:       dataDir,
	}, nil
}

func (c *Config) EnsureDirectories() error {
	if err := os.MkdirAll(c.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", c.DataDir, err)
	}
	return nil
}

func (c *Config) StateFilePath() string {
	return filepath.Join(c.DataDir, "state.yaml")
}

func (c *Config) KubeconfigPath() string {
	return filepath.Join(c.DataDir, "kubeconfig")
}
