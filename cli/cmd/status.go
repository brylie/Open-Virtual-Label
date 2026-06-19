package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display current label state and active work",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}

	// Print state document header
	statePath := ws.StateFile(wsPath)
	if stateData, err := os.ReadFile(statePath); err == nil {
		lines := strings.SplitN(string(stateData), "\n", 20)
		for i, line := range lines {
			if rest, ok := strings.CutPrefix(line, "### "); ok {
				fmt.Printf("Last session: %s\n\n", rest)
				break
			}
			if i > 15 {
				break
			}
		}
	}

	// Collect active releases
	artistIDs, err := ws.ListArtistIDs(wsPath)
	if err != nil {
		return err
	}
	if len(artistIDs) == 0 {
		fmt.Println("No artists found. Run 'ovl artist create' to add one.")
		return nil
	}

	filterArtist := cfg.ArtistID

	type releaseEntry struct {
		title    string
		artistID string
		id       string
		status   string
		tracks   int
	}
	var activeReleases []releaseEntry

	for _, aID := range artistIDs {
		if filterArtist != "" && aID != filterArtist {
			continue
		}
		rIDs, err := ws.ListReleaseIDs(wsPath, aID)
		if err != nil {
			continue
		}
		for _, rID := range rIDs {
			var rel models.Release
			if err := ws.ReadJSON(ws.ReleaseFile(wsPath, aID, rID), &rel); err != nil {
				continue
			}
			if rel.Status == statusLive || rel.Status == "archived" {
				continue
			}
			trackCount := len(rel.Tracks)
			activeReleases = append(activeReleases, releaseEntry{
				title:    rel.Title,
				artistID: aID,
				id:       rID,
				status:   rel.Status,
				tracks:   trackCount,
			})
		}
	}

	// Active releases
	fmt.Printf("Active releases (%d):\n", len(activeReleases))
	if len(activeReleases) == 0 {
		fmt.Println("  (none)")
	}
	for _, r := range activeReleases {
		fmt.Printf("  %s [%s] — %s (%d tracks)\n", r.title, r.artistID, r.status, r.tracks)
	}
	fmt.Println()

	// Pending outreach approvals
	oppDir := filepath.Join(wsPath, "opportunities")
	entries, _ := os.ReadDir(oppDir)
	var pendingApprovals []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		var opp models.Opportunity
		if err := ws.ReadJSON(filepath.Join(oppDir, e.Name()), &opp); err != nil {
			continue
		}
		if opp.Status == "draft-ready" {
			pendingApprovals = append(pendingApprovals, "outreach: draft for " + opp.Contact.Name)
		}
	}

	fmt.Printf("Pending approvals (%d):\n", len(pendingApprovals))
	if len(pendingApprovals) == 0 {
		fmt.Println("  (none)")
	}
	for _, p := range pendingApprovals {
		fmt.Println(" ", p)
	}
	fmt.Println()

	// Open loops: QC not run on ready-ish releases
	var openLoops []string
	for _, r := range activeReleases {
		if r.status == statusMastering || r.status == statusQC {
			openLoops = append(openLoops, "QC not yet complete on " + r.title)
		}
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		var opp models.Opportunity
		if err := ws.ReadJSON(filepath.Join(oppDir, e.Name()), &opp); err != nil {
			continue
		}
		if opp.FollowUpDue != nil && *opp.FollowUpDue != "" {
			openLoops = append(openLoops, fmt.Sprintf("Follow-up due: %s by %s", opp.Contact.Name, *opp.FollowUpDue))
		}
	}

	fmt.Printf("Open loops (%d):\n", len(openLoops))
	if len(openLoops) == 0 {
		fmt.Println("  (none)")
	}
	for _, l := range openLoops {
		fmt.Println(" ", l)
	}

	_ = output.IsJSON() // suppress unused import
	return nil
}
