package cmd

import (
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Compile and analyze streaming metrics",
}

var metricsSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Compile a metrics snapshot for a period (requires metrics-analyst agent)",
	RunE:  func(*cobra.Command, []string) error { return agentStub("metrics-analyst") },
}

func init() {
	metricsSnapshotCmd.Flags().String("period", "", "period YYYY-MM (required)")
	metricsSnapshotCmd.Flags().Bool("brief", false, "one-paragraph summary")
	_ = metricsSnapshotCmd.MarkFlagRequired("period")

	metricsCmd.AddCommand(metricsSnapshotCmd)
	rootCmd.AddCommand(metricsCmd)
}
