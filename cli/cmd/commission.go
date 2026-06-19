package cmd

import (
	"github.com/spf13/cobra"
)

var commissionCmd = &cobra.Command{
	Use:   "commission",
	Short: "Commission agreement shortcuts",
}

func init() {
	agreementCmd := &cobra.Command{
		Use:   "agreement",
		Short: "Generate a commission agreement from the workspace template (requires outreach-crm agent)",
		RunE:  func(*cobra.Command, []string) error { return agentStub("outreach-crm") },
	}
	agreementCmd.Flags().String("opportunity", "", "opportunity ID (required)")
	_ = agreementCmd.MarkFlagRequired("opportunity")

	commissionCmd.AddCommand(agreementCmd)
	rootCmd.AddCommand(commissionCmd)
}
