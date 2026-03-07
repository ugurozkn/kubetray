package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ugurozkn/kubetray/pkg/config"
	"github.com/ugurozkn/kubetray/pkg/k8s"
	"github.com/ugurozkn/kubetray/pkg/platform"
	"github.com/ugurozkn/kubetray/pkg/ui"
)

var cleanForce bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete the Kubernetes cluster completely",
	Long: `Completely remove the cluster and all associated data.

This is irreversible. All cluster data, volumes, and configurations will be deleted.

Examples:
  kubetray clean           # Prompt for confirmation
  kubetray clean --force   # Skip confirmation`,
	RunE: runClean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().BoolVarP(&cleanForce, "force", "f", false, "Skip confirmation prompt")
}

func runClean(cmd *cobra.Command, args []string) error {
	plat, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	clusterMgr := k8s.NewClusterManager(cfg, plat)

	// Confirm deletion
	if !cleanForce {
		ui.Warning("This will permanently delete cluster '%s' and remove %s.", cfg.ClusterName, cfg.DataDir)
		fmt.Print("Are you sure? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			ui.Info("Cancelled.")
			return nil
		}
	}

	spinner := ui.NewSpinner(fmt.Sprintf("Deleting cluster '%s'", cfg.ClusterName))
	spinner.Start()

	if err := clusterMgr.DeleteCluster(); err != nil {
		spinner.Error("Failed to delete cluster")
		return err
	}
	spinner.Success(fmt.Sprintf("Cluster '%s' deleted", cfg.ClusterName))

	// Purge ~/.kubetray directory
	spinner = ui.NewSpinner(fmt.Sprintf("Removing %s", cfg.DataDir))
	spinner.Start()

	if err := os.RemoveAll(cfg.DataDir); err != nil {
		spinner.Error(fmt.Sprintf("Failed to remove %s", cfg.DataDir))
		return err
	}
	spinner.Success(fmt.Sprintf("Removed %s", cfg.DataDir))

	fmt.Println()
	ui.Success("Environment cleaned up.")

	return nil
}
