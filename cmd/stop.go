package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ugurozkn/kubetray/pkg/config"
	"github.com/ugurozkn/kubetray/pkg/k8s"
	"github.com/ugurozkn/kubetray/pkg/platform"
	"github.com/ugurozkn/kubetray/pkg/ui"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Kubernetes cluster",
	Long: `Stop the running cluster without deleting it.

Data and configuration are preserved. Use 'kubetray start' to restart.

Examples:
  kubetray stop`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(cmd *cobra.Command, args []string) error {
	plat, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	clusterMgr := k8s.NewClusterManager(cfg, plat)

	if !clusterMgr.IsRunning() {
		ui.Warning("Cluster '%s' is not running", cfg.ClusterName)
		return nil
	}

	spinner := ui.NewSpinner(fmt.Sprintf("Stopping cluster '%s'", cfg.ClusterName))
	spinner.Start()

	if err := clusterMgr.StopCluster(); err != nil {
		spinner.Error("Failed to stop cluster")
		return err
	}
	spinner.Success(fmt.Sprintf("Cluster '%s' stopped", cfg.ClusterName))

	// Update state
	state, err := config.LoadState(cfg)
	if err == nil {
		state.Status = "stopped"
		state.K3sRunning = false
		_ = state.Save(cfg)
	}

	fmt.Println()
	ui.Info("Data preserved. Run 'kubetray start' to restart.")

	return nil
}
