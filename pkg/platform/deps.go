package platform

import (
	"fmt"
	"os/exec"
	"strings"
)

// Dependency represents a required system dependency
type Dependency struct {
	Name        string
	Command     string
	VersionFlag string
	Required    bool
	InstallHint string
}

// DependencyStatus represents the status of a dependency check
type DependencyStatus struct {
	Dependency *Dependency
	Installed  bool
	Version    string
	Error      error
}

// RuntimeStatus represents the status of container/k8s runtime
type RuntimeStatus struct {
	DockerRunning bool
	DockerInfo    string
	RuntimeType   string // "docker-desktop", "colima", "orbstack", "rancher", "unknown"
}

// GetRequiredDependencies returns the list of required dependencies for the platform
func GetRequiredDependencies(p *Platform) []Dependency {
	deps := []Dependency{
		{
			Name:        "helm",
			Command:     "helm",
			VersionFlag: "version --short",
			Required:    true,
			InstallHint: "Install with: brew install helm (macOS) or curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash",
		},
		{
			Name:        "kubectl",
			Command:     "kubectl",
			VersionFlag: "version --client -o yaml",
			Required:    false,
			InstallHint: "Install with: brew install kubectl (macOS) or see https://kubernetes.io/docs/tasks/tools/",
		},
	}

	if p.OS == OSMacOS {
		deps = append(deps, Dependency{
			Name:        "docker",
			Command:     "docker",
			VersionFlag: "--version",
			Required:    true,
			InstallHint: "Install with: brew install docker",
		})
	}

	return deps
}

// CheckDependency checks if a single dependency is installed
func CheckDependency(dep *Dependency) *DependencyStatus {
	status := &DependencyStatus{
		Dependency: dep,
	}

	path, err := exec.LookPath(dep.Command)
	if err != nil {
		status.Installed = false
		status.Error = fmt.Errorf("%s not found in PATH", dep.Name)
		return status
	}

	status.Installed = true

	// Try to get version
	args := strings.Fields(dep.VersionFlag)
	cmd := exec.Command(path, args...)
	out, err := cmd.Output()
	if err != nil {
		status.Version = "unknown"
	} else {
		// Extract first line and clean up
		version := strings.TrimSpace(string(out))
		if idx := strings.Index(version, "\n"); idx != -1 {
			version = version[:idx]
		}
		status.Version = version
	}

	return status
}

// CheckAllDependencies checks all required dependencies
func CheckAllDependencies(p *Platform) ([]*DependencyStatus, error) {
	deps := GetRequiredDependencies(p)
	statuses := make([]*DependencyStatus, 0, len(deps))
	var missingRequired []string

	for i := range deps {
		status := CheckDependency(&deps[i])
		statuses = append(statuses, status)

		if !status.Installed && status.Dependency.Required {
			missingRequired = append(missingRequired, status.Dependency.Name)
		}
	}

	if len(missingRequired) > 0 {
		return statuses, fmt.Errorf("missing required dependencies: %s", strings.Join(missingRequired, ", "))
	}

	return statuses, nil
}

// CheckDockerRuntime checks if Docker daemon is running and identifies the runtime
func CheckDockerRuntime() *RuntimeStatus {
	status := &RuntimeStatus{}

	// Check if docker daemon is responding
	cmd := exec.Command("docker", "info", "--format", "{{.Name}}")
	out, err := cmd.Output()
	if err != nil {
		status.DockerRunning = false
		return status
	}

	status.DockerRunning = true
	status.DockerInfo = strings.TrimSpace(string(out))

	// Try to identify the runtime
	status.RuntimeType = identifyRuntime()

	return status
}

func identifyRuntime() string {
	// Check for Colima
	cmd := exec.Command("colima", "status")
	if out, err := cmd.Output(); err == nil && strings.Contains(string(out), "Running") {
		return "colima"
	}

	// Check Docker context for clues
	cmd = exec.Command("docker", "context", "show")
	if out, err := cmd.Output(); err == nil {
		context := strings.TrimSpace(string(out))
		switch {
		case strings.Contains(context, "colima"):
			return "colima"
		case strings.Contains(context, "orbstack"):
			return "orbstack"
		case strings.Contains(context, "rancher"):
			return "rancher-desktop"
		case context == "default" || context == "desktop-linux":
			return "docker-desktop"
		}
	}

	return "unknown"
}

// GetDockerRuntimeHint returns instructions to start a Docker runtime
func GetDockerRuntimeHint() string {
	return `Docker daemon is not running. Start your container runtime:

  Colima:       colima start
  Docker Desktop: open -a Docker
  OrbStack:     open -a OrbStack

Then run 'kubetray start' again.`
}
