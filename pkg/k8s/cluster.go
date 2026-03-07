package k8s

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ugurozkn/kubetray/pkg/config"
	"github.com/ugurozkn/kubetray/pkg/platform"
)

// ClusterManager handles k3s cluster operations
type ClusterManager struct {
	config   *config.Config
	platform *platform.Platform
}

// NewClusterManager creates a new cluster manager
func NewClusterManager(cfg *config.Config, plat *platform.Platform) *ClusterManager {
	return &ClusterManager{
		config:   cfg,
		platform: plat,
	}
}

// StartCluster starts the k3s cluster using k3d
func (m *ClusterManager) StartCluster(cpus int, memory string) error {
	// Ensure k3d is installed
	if err := m.ensureK3d(); err != nil {
		return err
	}

	clusterName := m.config.ClusterName

	// Check if cluster already exists
	if m.clusterExists() {
		// Start existing cluster
		cmd := exec.Command("k3d", "cluster", "start", clusterName)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to start cluster: %s", string(out))
		}
	} else {
		// Create new cluster with port mappings for ingress
		args := []string{
			"cluster", "create", clusterName,
			"--agents", "0",
			"--servers", "1",
			"--k3s-arg", "--disable=traefik@server:0", // We install our own traefik
			"-p", "80:80@loadbalancer",                // HTTP
			"-p", "443:443@loadbalancer",              // HTTPS
			"--wait",
			"--timeout", "5m",
		}

		cmd := exec.Command("k3d", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create cluster: %w", err)
		}
	}

	// Export kubeconfig
	if err := m.exportKubeconfig(); err != nil {
		return fmt.Errorf("failed to export kubeconfig: %w", err)
	}

	// Wait for cluster to be ready
	return m.waitForKubernetes()
}

// StopCluster stops the k3s cluster
func (m *ClusterManager) StopCluster() error {
	if !m.clusterExists() {
		return nil
	}

	cmd := exec.Command("k3d", "cluster", "stop", m.config.ClusterName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop cluster: %s", string(out))
	}

	return nil
}

// DeleteCluster completely removes the cluster
func (m *ClusterManager) DeleteCluster() error {
	if !m.clusterExists() {
		return nil
	}

	cmd := exec.Command("k3d", "cluster", "delete", m.config.ClusterName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete cluster: %s", string(out))
	}

	// Clean up kubeconfig
	kubeconfigPath := m.config.KubeconfigPath()
	_ = os.Remove(kubeconfigPath)

	return nil
}

// IsRunning checks if the cluster is running
func (m *ClusterManager) IsRunning() bool {
	cmd := exec.Command("k3d", "cluster", "list", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	// Check if our cluster is in the list and running
	return strings.Contains(string(out), fmt.Sprintf(`"name":"%s"`, m.config.ClusterName)) &&
		strings.Contains(string(out), `"serversRunning":1`)
}

// GetKubeconfig returns the path to the kubeconfig file
func (m *ClusterManager) GetKubeconfig() string {
	return m.config.KubeconfigPath()
}

// clusterExists checks if the cluster exists (running or stopped)
func (m *ClusterManager) clusterExists() bool {
	cmd := exec.Command("k3d", "cluster", "list", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), fmt.Sprintf(`"name":"%s"`, m.config.ClusterName))
}

// ensureK3d installs k3d if not present
func (m *ClusterManager) ensureK3d() error {
	if _, err := exec.LookPath("k3d"); err == nil {
		return nil
	}

	// Install k3d
	var cmd *exec.Cmd
	switch m.platform.OS {
	case platform.OSMacOS:
		cmd = exec.Command("brew", "install", "k3d")
	case platform.OSLinux:
		cmd = exec.Command("sh", "-c", "curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash")
	default:
		return fmt.Errorf("unsupported platform for k3d installation")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install k3d: %w", err)
	}

	return nil
}

// exportKubeconfig exports the kubeconfig to our config directory
func (m *ClusterManager) exportKubeconfig() error {
	kubeconfigPath := m.config.KubeconfigPath()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(kubeconfigPath), 0755); err != nil {
		return err
	}

	cmd := exec.Command("k3d", "kubeconfig", "get", m.config.ClusterName)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	return os.WriteFile(kubeconfigPath, out, 0600)
}

// waitForKubernetes waits for the cluster to be ready
func (m *ClusterManager) waitForKubernetes() error {
	kubeconfigPath := m.config.KubeconfigPath()
	timeout := time.After(3 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for kubernetes to be ready")
		case <-ticker.C:
			cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "nodes", "--no-headers")
			out, err := cmd.Output()
			if err != nil {
				continue
			}

			if strings.Contains(string(out), "Ready") && !strings.Contains(string(out), "NotReady") {
				return nil
			}
		}
	}
}

// isKubernetesReady checks if kubernetes is ready
func (m *ClusterManager) isKubernetesReady() bool {
	kubeconfigPath := m.config.KubeconfigPath()
	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "nodes", "--no-headers")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), "Ready") && !strings.Contains(string(out), "NotReady")
}
