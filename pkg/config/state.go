package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ClusterState represents the current state of the kubetray environment
type ClusterState struct {
	// Cluster info
	ClusterName string    `yaml:"cluster_name"`
	Profile     string    `yaml:"profile"`
	Status      string    `yaml:"status"` // running, stopped, unknown
	StartedAt   time.Time `yaml:"started_at,omitempty"`

	// Resource allocation
	CPUs   int    `yaml:"cpus"`
	Memory string `yaml:"memory"`

	// Installed components
	Components []ComponentState `yaml:"components,omitempty"`

	// Runtime info
	ColimaProfile string `yaml:"colima_profile,omitempty"` // macOS only
	K3sRunning    bool   `yaml:"k3s_running"`
}

// ComponentState represents the state of an installed component
type ComponentState struct {
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Namespace   string    `yaml:"namespace"`
	Status      string    `yaml:"status"` // installed, pending, failed
	InstalledAt time.Time `yaml:"installed_at"`
	URL         string    `yaml:"url,omitempty"`
}

// LoadState reads the state from the state file
func LoadState(cfg *Config) (*ClusterState, error) {
	statePath := cfg.StateFilePath()

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty state if file doesn't exist
			return &ClusterState{
				ClusterName: cfg.ClusterName,
				Status:      "stopped",
			}, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state ClusterState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

// Save writes the state to the state file
func (s *ClusterState) Save(cfg *Config) error {
	if err := cfg.EnsureDirectories(); err != nil {
		return err
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to serialize state: %w", err)
	}

	statePath := cfg.StateFilePath()
	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// IsRunning returns true if the cluster is running
func (s *ClusterState) IsRunning() bool {
	return s.Status == "running" && s.K3sRunning
}

// GetComponent returns a component by name
func (s *ClusterState) GetComponent(name string) *ComponentState {
	for i := range s.Components {
		if s.Components[i].Name == name {
			return &s.Components[i]
		}
	}
	return nil
}

// AddComponent adds or updates a component in the state
func (s *ClusterState) AddComponent(comp ComponentState) {
	for i := range s.Components {
		if s.Components[i].Name == comp.Name {
			s.Components[i] = comp
			return
		}
	}
	s.Components = append(s.Components, comp)
}

// RemoveComponent removes a component from the state
func (s *ClusterState) RemoveComponent(name string) {
	for i := range s.Components {
		if s.Components[i].Name == name {
			s.Components = append(s.Components[:i], s.Components[i+1:]...)
			return
		}
	}
}

// Clear resets the state
func (s *ClusterState) Clear() {
	s.Status = "stopped"
	s.K3sRunning = false
	s.Components = nil
	s.StartedAt = time.Time{}
}

// Delete removes the state file
func DeleteState(cfg *Config) error {
	statePath := cfg.StateFilePath()
	if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete state file: %w", err)
	}
	return nil
}
