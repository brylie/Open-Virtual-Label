package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/prompt"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var isrcCmd = &cobra.Command{
	Use:   "isrc",
	Short: "Manage ISRCs",
}

var isrcAssignRelease, isrcAssignTrack string
var isrcAssignCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign ISRCs to tracks that don't have one",
	Long: `Assign ISRCs to tracks that currently have no ISRC.
OVL does not register ISRCs. This command records codes already obtained
from the artist's distributor or national ISRC agency.`,
	RunE: runISRCAssign,
}

var isrcPattern = regexp.MustCompile(`^[A-Z]{2}-?[A-Z0-9]{3}-?\d{2}-?\d{5}$`)

func init() {
	isrcAssignCmd.Flags().StringVar(&isrcAssignRelease, "release", "", "release ID (required)")
	isrcAssignCmd.Flags().StringVar(&isrcAssignTrack, "track", "", "assign only to this track ID")
	_ = isrcAssignCmd.MarkFlagRequired("release")

	isrcCmd.AddCommand(isrcAssignCmd)
	rootCmd.AddCommand(isrcCmd)
}

func runISRCAssign(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	releaseID := isrcAssignRelease

	artistID, err := ws.FindRelease(wsPath, releaseID)
	if err != nil {
		return err
	}

	tIDs, err := ws.ListTrackIDs(wsPath, artistID, releaseID)
	if err != nil {
		return err
	}

	assigned := 0
	for _, tID := range tIDs {
		if isrcAssignTrack != "" && tID != isrcAssignTrack {
			continue
		}
		trackPath := ws.TrackFile(wsPath, artistID, releaseID, tID)
		var t models.Track
		if err := ws.ReadJSON(trackPath, &t); err != nil {
			continue
		}
		if t.ISRC != nil && *t.ISRC != "" {
			output.Info("  %s: already has ISRC %s", tID, *t.ISRC)
			continue
		}

		for {
			code, err := prompt.Ask(fmt.Sprintf("ISRC for %q (format CC-XXX-YY-NNNNN)", t.Title), "")
			if err != nil {
				return err
			}
			if code == "" {
				output.Info("  %s: skipped", tID)
				break
			}
			normalized := strings.ToUpper(strings.ReplaceAll(code, "-", ""))
			formatted := fmt.Sprintf("%s-%s-%s-%s",
				normalized[0:2], normalized[2:5], normalized[5:7], normalized[7:12])

			if !isrcPattern.MatchString(formatted) {
				fmt.Printf("  Invalid ISRC format: %q — expected CC-XXX-YY-NNNNN\n", code)
				continue
			}

			t.ISRC = &formatted
			raw, _ := json.Marshal(t)
			if errs, err := schema.Validate("track", raw); err != nil {
				return err
			} else if len(errs) > 0 {
				fmt.Println("  Validation error:", strings.Join(errs, "; "))
				t.ISRC = nil
				continue
			}

			if err := ws.WriteJSON(trackPath, t); err != nil {
				return err
			}
			output.Success("%s — ISRC set to %s", tID, formatted)
			assigned++
			break
		}
	}

	fmt.Printf("\n%d ISRC(s) assigned.\n", assigned)
	return nil
}
