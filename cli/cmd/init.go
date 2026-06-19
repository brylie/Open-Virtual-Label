package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/prompt"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var initForce bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a new label workspace",
	Long:  "Create the workspace directory structure and write the initial label profile and state document.",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initForce, "force", false, "overwrite an existing workspace (prompts for confirmation)")
	rootCmd.AddCommand(initCmd)
}

func runInit(_ *cobra.Command, _ []string) error {
	target := cfg.WorkspacePath
	if target == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		target = filepath.Join(cwd, "workspace")
	}

	if info, err := os.Stat(target); err == nil && info.IsDir() { //nolint:nestif // init flow is inherently sequential: stat→force-check→confirm→proceed
		if !initForce {
			return &ExitError{Code: 1, Msg: fmt.Sprintf(
				"WORKSPACE_EXISTS: %s already exists; use --force to overwrite", target)}
		}
		if !cfg.SkipConfirm {
			ok, err := prompt.Confirm(fmt.Sprintf("Overwrite existing workspace at %s?", target))
			if err != nil {
				return err
			}
			if !ok {
				return &ExitError{Code: 4, Msg: msgCanceled}
			}
		}
	}

	labelName, err := prompt.Ask("Label name", "")
	if err != nil {
		return err
	}
	if labelName == "" {
		return errors.New("label name is required")
	}
	contactEmail, err := prompt.Ask("Contact email", "")
	if err != nil {
		return err
	}
	defaultLicense, err := prompt.Ask("Default license", "CC BY 4.0")
	if err != nil {
		return err
	}
	distributor, err := prompt.Ask("Primary distributor (e.g. amuse, distrokid)", "")
	if err != nil {
		return err
	}

	labelID := ws.ToSlug(labelName)
	today := time.Now().Format("2006-01-02")

	label := models.Label{
		SchemaVersion:      "1",
		ID:                 labelID,
		Name:               labelName,
		DefaultLicense:     defaultLicense,
		DefaultDistributor: distributor,
		CreatedDate:        today,
	}
	if contactEmail != "" {
		label.Contact = &models.LabelContact{Email: contactEmail}
	}

	raw, err := json.Marshal(label)
	if err != nil {
		return err
	}
	errs, err := schema.Validate("label", raw)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: "SCHEMA_VALIDATION_FAILED: " + strings.Join(errs, "; ")}
	}

	dirs := []string{
		filepath.Join(target, "label"),
		filepath.Join(target, "state"),
		filepath.Join(target, "artists"),
		filepath.Join(target, "opportunities"),
		filepath.Join(target, "finance"),
		filepath.Join(target, "metrics"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	profilePath := ws.LabelFile(target)
	if err := ws.WriteJSON(profilePath, label); err != nil {
		return err
	}
	output.Success("Created %s", profilePath)

	statePath := ws.StateFile(target)
	stateContent := fmt.Sprintf(
		"# Label State\n\nLabel: %s\nLicense: %s\n\n## Session Log\n\n### %s — workspace initialized\n\nWorkspace created.\n",
		labelName, defaultLicense, today)
	if err := os.WriteFile(statePath, []byte(stateContent), 0o644); err != nil {
		return err
	}
	output.Success("Created %s", statePath)

	fmt.Printf("\nWorkspace initialized at %s\n\n", target)
	fmt.Println("Next steps:")
	fmt.Println("  ovl artist create    — add an artist profile")
	fmt.Println("  ovl status           — view label state")
	return nil
}
