package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/prompt"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var masteringCmd = &cobra.Command{
	Use:   statusMastering,
	Short: "Mastering sessions and profiles",
}

// mastering start
var masteringStartTrack string
var masteringStartRemaster bool
var masteringStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch a mastering session (requires agent integration)",
	RunE:  runMasteringStart,
}

// mastering profile
var masteringProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage mastering profiles",
}

var masteringProfileCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new mastering profile interactively",
	RunE:  runMasteringProfileCreate,
}

var masteringProfileListCmd = &cobra.Command{
	Use:   cmdUseList,
	Short: "List mastering profiles for an artist",
	RunE:  runMasteringProfileList,
}

func init() {
	masteringStartCmd.Flags().StringVar(&masteringStartTrack, "track", "", "track ID (required)")
	masteringStartCmd.Flags().BoolVar(&masteringStartRemaster, "remaster", false, "re-master a track with existing mastering data")
	_ = masteringStartCmd.MarkFlagRequired("track")

	masteringProfileCmd.AddCommand(masteringProfileCreateCmd, masteringProfileListCmd)
	masteringCmd.AddCommand(masteringStartCmd, masteringProfileCmd)
	rootCmd.AddCommand(masteringCmd)
}

func runMasteringStart(_ *cobra.Command, _ []string) error {
	fmt.Println("Mastering sessions require agent integration.")
	fmt.Println("The mastering-companion agent guides you through measurements and")
	fmt.Println("writes results to the track record on completion.")
	fmt.Println()
	fmt.Println("Configure agent integration via cli/AGENT-INTEGRATION.md.")
	return nil
}

func runMasteringProfileCreate(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	artistID, err := resolveArtistID(wsPath)
	if err != nil {
		return err
	}

	name, err := prompt.Ask("Profile name (e.g. Ambient / Streaming)", "")
	if err != nil || name == "" {
		return errors.New("profile name is required")
	}
	desc, err := prompt.Ask("Description (when to use this profile)", "")
	if err != nil {
		return err
	}
	lufsMinStr, err := prompt.Ask("LUFS target min (e.g. -16)", "-16")
	if err != nil {
		return err
	}
	lufsMaxStr, err := prompt.Ask("LUFS target max (e.g. -14)", "-14")
	if err != nil {
		return err
	}
	peakStr, err := prompt.Ask("True peak ceiling dBTP (e.g. -1.0)", "-1.0")
	if err != nil {
		return err
	}

	lufsMin, _ := strconv.ParseFloat(lufsMinStr, 64)
	lufsMax, _ := strconv.ParseFloat(lufsMaxStr, 64)
	peak, _ := strconv.ParseFloat(peakStr, 64)

	slug := ws.ToSlug(name)
	today := time.Now().Format("2006-01-02")

	profile := models.MasteringProfile{
		SchemaVersion: "1",
		ID:            slug,
		Name:          name,
		Description:   desc,
		Targets: models.MasteringTargets{
			IntegratedLUFS: models.LUFSRange{Min: lufsMin, Max: lufsMax},
			TruePeakDBTP:   peak,
			SampleRateHz:   44100,
			BitDepth:       24,
		},
		CreatedDate: today,
	}

	raw, _ := json.Marshal(profile)
	if errs, err := schema.Validate("mastering-profile", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: strings.Join(errs, "\n")}
	}

	profilePath := ws.MasteringProfileFile(wsPath, artistID, slug)
	if err := ws.WriteJSON(profilePath, profile); err != nil {
		return err
	}
	output.Success("Created %s", profilePath)
	fmt.Printf("Profile ID: %s\n", slug)
	return nil
}

func runMasteringProfileList(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	artistID, err := resolveArtistID(wsPath)
	if err != nil {
		return err
	}

	ids, err := ws.ListMasteringProfileIDs(wsPath, artistID)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		fmt.Println("No mastering profiles found.")
		return nil
	}

	rows := make([][]string, 0, len(ids))
	for _, id := range ids {
		var p models.MasteringProfile
		if err := ws.ReadJSON(ws.MasteringProfileFile(wsPath, artistID, id), &p); err != nil {
			continue
		}
		sessions := strconv.Itoa(len(p.SessionNotes))
		rows = append(rows, []string{
			id, p.Name,
			fmt.Sprintf("%.0f to %.0f", p.Targets.IntegratedLUFS.Min, p.Targets.IntegratedLUFS.Max),
			fmt.Sprintf("%.1f", p.Targets.TruePeakDBTP),
			sessions,
		})
	}
	output.Table([]string{"ID", "Name", "LUFS range", "True peak", "Sessions"}, rows)
	return nil
}
