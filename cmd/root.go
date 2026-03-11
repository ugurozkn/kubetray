package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kubetray %s (commit: %s, built: %s)\n", version, commit, date)
		},
	})
}

func Execute() error {
	return rootCmd.Execute()
}
