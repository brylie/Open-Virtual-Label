package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

const Version = "0.1.0"

// ExitError allows commands to return specific exit codes defined in the spec.
type ExitError struct {
	Code int
	Msg  string
}

func (e *ExitError) Error() string { return e.Msg }

// Config holds values from global flags and environment variables.
type Config struct {
	WorkspacePath string
	ArtistID      string
	OutputJSON    bool
	Quiet         bool
	SkipConfirm   bool
}

var cfg Config

var rootCmd = &cobra.Command{
	Use:   "ovl",
	Short: "Open Virtual Label CLI",
	Long: `ovl manages the Open Virtual Label workspace: artists, releases, tracks,
mastering, QC, archival, outreach, and the release pipeline.`,
	Version: Version,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		output.SetQuiet(cfg.Quiet)
		output.SetJSON(cfg.OutputJSON)
		return nil
	},
}

// Execute runs the root command and returns an OS exit code.
func Execute() int {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			return exitErr.Code
		}
		return 1
	}
	return 0
}

func init() {
	cobra.OnInitialize(applyEnvDefaults)

	rootCmd.PersistentFlags().StringVar(&cfg.WorkspacePath, "workspace", "",
		"path to workspace directory (default: ./workspace, walks up)")
	rootCmd.PersistentFlags().StringVar(&cfg.ArtistID, "artist", "",
		"scope to a specific artist ID")
	rootCmd.PersistentFlags().BoolVar(&cfg.OutputJSON, "json", false,
		"output result as JSON")
	rootCmd.PersistentFlags().BoolVar(&cfg.Quiet, "quiet", false,
		"suppress informational output (errors still print to stderr)")
	rootCmd.PersistentFlags().BoolVar(&cfg.SkipConfirm, "yes", false,
		"skip confirmation prompts (critical gates are never bypassed)")
}

func applyEnvDefaults() {
	if cfg.WorkspacePath == "" {
		if env := os.Getenv("OVL_WORKSPACE"); env != "" {
			cfg.WorkspacePath = env
		}
	}
	if cfg.ArtistID == "" {
		if env := os.Getenv("OVL_ARTIST"); env != "" {
			cfg.ArtistID = env
		}
	}
	if os.Getenv("OVL_JSON") == "1" {
		cfg.OutputJSON = true
	}
	if os.Getenv("OVL_YES") == "1" {
		cfg.SkipConfirm = true
	}
}

// resolveWorkspace returns the discovered workspace path or an exit error.
func resolveWorkspace() (string, error) {
	path, err := ws.Find(cfg.WorkspacePath)
	if err != nil {
		return "", &ExitError{Code: 2, Msg: err.Error()}
	}
	return path, nil
}
