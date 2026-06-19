package cmd

import (
	"fmt"

	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Manage MCP integrations",
}

var mcpListCmd = &cobra.Command{
	Use:   cmdUseList,
	Short: "List available MCPs and their connection status",
	RunE:  runMCPList,
}

var mcpConnectCmd = &cobra.Command{
	Use:   "connect <mcp-name>",
	Short: "Connect an MCP integration",
	Args:  cobra.ExactArgs(1),
	RunE:  runMCPConnect,
}

var mcpDisconnectCmd = &cobra.Command{
	Use:   "disconnect <mcp-name>",
	Short: "Disconnect an MCP integration",
	Args:  cobra.ExactArgs(1),
	RunE:  runMCPDisconnect,
}

func init() {
	mcpCmd.AddCommand(mcpListCmd, mcpConnectCmd, mcpDisconnectCmd)
	rootCmd.AddCommand(mcpCmd)
}

var knownMCPs = []struct {
	name        string
	description string
	enables     string
}{
	{"gmail", "Gmail integration for sending outreach", "ovl outreach send"},
	{"google-calendar", "Google Calendar for scheduling follow-ups", "ovl outreach follow-up"},
	{mcpInternetArchive, "Internet Archive S3-compatible upload", "ovl archive push"},
	{mcpAmuse, "Amuse distributor API (read-only where permitted)", "ovl release submit"},
}

func runMCPList(_ *cobra.Command, _ []string) error {
	rows := make([][]string, len(knownMCPs))
	for i, m := range knownMCPs {
		rows[i] = []string{m.name, m.description, "not configured", m.enables}
	}
	output.Table([]string{"MCP", "Description", "Status", "Enables"}, rows)
	fmt.Println()
	fmt.Println("Note: MCP credentials are stored locally and never written to workspace records.")
	return nil
}

func runMCPConnect(_ *cobra.Command, args []string) error {
	name := args[0]
	fmt.Printf("MCP '%s' connection requires Claude Code with MCP support.\n", name)
	fmt.Println("Credentials are stored locally and never written to workspace JSON records.")
	fmt.Println()
	fmt.Println("Configure MCPs in your Claude Code settings or .claude/settings.json:")
	fmt.Printf("  See: https://docs.anthropic.com/claude-code/mcp\n")
	return &ExitError{Code: 5, Msg: fmt.Sprintf("MCP '%s' not configured", name)}
}

func runMCPDisconnect(_ *cobra.Command, args []string) error {
	fmt.Printf("To disconnect '%s', remove it from your Claude Code MCP configuration.\n", args[0])
	return nil
}
