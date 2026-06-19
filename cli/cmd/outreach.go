package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var outreachCmd = &cobra.Command{
	Use:   "outreach",
	Short: "Manage the outreach and CRM pipeline",
}

func agentStub(agentName string) error {
	fmt.Printf("This command requires the '%s' agent.\n", agentName)
	fmt.Println("Configure agent integration via cli/AGENT-INTEGRATION.md.")
	return nil
}

func init() {
	for _, sub := range []*cobra.Command{
		{Use: "research", Short: "Find new outreach opportunities (requires outreach-crm agent)",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "review", Short: "Review opportunities with pending draft approvals",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: cmdUseDraft, Short: "Draft outreach for a specific opportunity",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "send", Short: "Send an approved outreach message",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "follow-up", Short: "Draft a follow-up for an unanswered opportunity",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "log-response", Short: "Record a response from an outreach contact",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "score", Short: "Score an opportunity for fit",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "intake", Short: "Record a new inbound inquiry as an opportunity",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "close", Short: "Record the final outcome of an opportunity",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
		{Use: "log", Short: "Manually log an action on an opportunity",
			RunE: func(*cobra.Command, []string) error { return agentStub("outreach-crm") }},
	} {
		outreachCmd.AddCommand(sub)
	}
	rootCmd.AddCommand(outreachCmd)
}
