package cmd

import (
	"fmt"
	"os"

	"github.com/open-virtual-label/ovl/internal/prompt"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Manage the label state document",
}

var stateShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the full contents of label-state.md",
	RunE:  runStateShow,
}

var stateSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Invoke the orchestrator to write a session summary",
	RunE:  runStateSync,
}

func init() {
	stateCmd.AddCommand(stateShowCmd)
	stateCmd.AddCommand(stateSyncCmd)
	rootCmd.AddCommand(stateCmd)
}

func runStateShow(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(ws.StateFile(wsPath))
	if err != nil {
		return fmt.Errorf("cannot read state file: %w", err)
	}
	fmt.Print(string(data))
	return nil
}

func runStateSync(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	statePath := ws.StateFile(wsPath)
	if _, err := os.Stat(statePath); err != nil {
		return fmt.Errorf("state file not found: %w", err)
	}

	fmt.Println("Agent integration required for state sync.")
	fmt.Println("This command invokes the orchestrator agent to produce a session summary.")
	fmt.Println("Configure agent integration via cli/AGENT-INTEGRATION.md.")
	fmt.Println()

	// Placeholder: manual note entry
	note, err := prompt.Ask("Add a manual note to the state (leave blank to cancel)", "")
	if err != nil {
		return err
	}
	if note == "" {
		return &ExitError{Code: 4, Msg: msgCanceled}
	}

	ok, err := prompt.Confirm(fmt.Sprintf("Append note to %s?", statePath))
	if err != nil {
		return err
	}
	if !ok {
		return &ExitError{Code: 4, Msg: msgCanceled}
	}

	f, err := os.OpenFile(statePath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, err = fmt.Fprintf(f, "\n- %s\n", note)
	return err
}
