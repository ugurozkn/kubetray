package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubetray",
	Short: "Local K8s, served on a tray",
	Long: `KubeTray is a command-line tool that sets up complete Kubernetes
development environments in minutes, not hours.

Stop spending time on infrastructure setup and start building.`,
	SilenceUsage:      true,
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() error {
	return rootCmd.Execute()
}
