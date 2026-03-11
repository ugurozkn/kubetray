package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugurozkn/kubetray/pkg/config"
	"github.com/ugurozkn/kubetray/pkg/k8s"
	"github.com/ugurozkn/kubetray/pkg/platform"
	"github.com/ugurozkn/kubetray/pkg/ui"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Kubernetes environment",
	Long: `Start a local Kubernetes cluster using k3s via k3d.

This will:
  1. Check Docker runtime (Docker Desktop, Colima, OrbStack, etc.)
  2. Create a lightweight k3s cluster via k3d (k3s in Docker)
  3. Configure kubeconfig for kubectl access

Resources default to 2 CPUs and 2G RAM. Override with --cpus and --memory flags.

Examples:
  kubetray start                       # Start with 2 CPUs, 2G RAM
  kubetray start --cpus 4 --memory 8G  # Start with custom resources`,
	RunE: runStart,
}

var (
	startCPUs   int
	startMemory string
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVar(&startCPUs, "cpus", 0, "CPU allocation (default: 2)")
	startCmd.Flags().StringVar(&startMemory, "memory", "", "Memory allocation, e.g. 2G, 4G (default: 2G)")
}

func runStart(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	fmt.Println()
	ui.Header("KubeTray")
	fmt.Println()

	// Detect platform
	spinner := ui.NewSpinner("Detecting platform")
	spinner.Start()

	plat, err := platform.Detect()
	if err != nil {
		spinner.Error("Platform detection failed")
		return err
	}
	spinner.Success(fmt.Sprintf("Platform: %s", plat.String()))

	// Check dependencies
	spinner = ui.NewSpinner("Checking dependencies")
	spinner.Start()

	statuses, err := platform.CheckAllDependencies(plat)
	if err != nil {
		spinner.Error("Missing dependencies")
		fmt.Println()
		for _, s := range statuses {
			if !s.Installed && s.Dependency.Required {
				ui.Error("%s not found", s.Dependency.Name)
				ui.Step("%s", s.Dependency.InstallHint)
			}
		}
		return err
	}
	spinner.Success("Dependencies: helm, kubectl, docker")

	// Check Docker runtime (macOS only)
	var runtimeName string
	if plat.OS == platform.OSMacOS {
		spinner = ui.NewSpinner("Checking Docker runtime")
		spinner.Start()

		runtime := platform.CheckDockerRuntime()
		if !runtime.DockerRunning {
			spinner.Error("Docker daemon not running")
			fmt.Println()
			fmt.Println(platform.GetDockerRuntimeHint())
			return fmt.Errorf("docker daemon not running")
		}

		runtimeName = runtime.RuntimeType
		if runtimeName == "unknown" {
			runtimeName = "docker"
		}
		spinner.Success(fmt.Sprintf("Container runtime: %s", runtimeName))
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Check if already running
	clusterMgr := k8s.NewClusterManager(cfg, plat)
	if clusterMgr.IsRunning() {
		fmt.Println()
		ui.Warning("Cluster '%s' is already running", cfg.ClusterName)
		return nil
	}

	// Start cluster
	fmt.Println()
	ui.Info("Creating k3s cluster via k3d...")
	fmt.Println()

	spinner = ui.NewSpinner(fmt.Sprintf("Creating cluster '%s'", cfg.ClusterName))
	spinner.Start()

	// Apply flag overrides
	cpus := cfg.DefaultCPUs
	memory := cfg.DefaultMemory
	if startCPUs > 0 {
		cpus = startCPUs
	}
	if startMemory != "" {
		memory = startMemory
	}

	if err := clusterMgr.StartCluster(cpus, memory); err != nil {
		spinner.Error("Failed to create cluster")
		return err
	}
	spinner.Success(fmt.Sprintf("Cluster '%s' created", cfg.ClusterName))

	// Print summary
	elapsed := time.Since(startTime).Round(time.Second)
	fmt.Println()
	ui.Success("Cluster is ready!")
	fmt.Println()

	ui.Header("Cluster Details")
	table := ui.NewTable("Property", "Value")
	table.AddRow("Cluster name", cfg.ClusterName)
	table.AddRow("Kubernetes", "k3s (via k3d)")
	if runtimeName != "" {
		table.AddRow("Container runtime", runtimeName)
	}
	table.AddRow("Resources", fmt.Sprintf("%d CPUs, %s RAM", cpus, memory))
	table.AddRow("Kubeconfig", cfg.KubeconfigPath())
	table.AddRow("Setup time", elapsed.String())
	table.Render()

	fmt.Println()
	ui.Header("Quick Start")
	fmt.Printf("  export KUBECONFIG=%s\n", cfg.KubeconfigPath())
	fmt.Println("  kubectl get nodes")
	fmt.Println("  kubectl get pods -A")

	return nil
}
