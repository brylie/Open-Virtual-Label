package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	"github.com/open-virtual-label/ovl/internal/prompt"
	"github.com/open-virtual-label/ovl/internal/schema"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var artistCmd = &cobra.Command{
	Use:   "artist",
	Short: "Manage artist profiles",
}

var artistCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new artist profile interactively",
	RunE:  runArtistCreate,
}

var artistListCmd = &cobra.Command{
	Use:   cmdUseList,
	Short: "List all artist profiles",
	RunE:  runArtistList,
}

var artistShowCmd = &cobra.Command{
	Use:   "show <artist-id>",
	Short: "Display an artist profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistShow,
}

var artistAddAliasName string
var artistAddAliasCmd = &cobra.Command{
	Use:   "add-alias <artist-id>",
	Short: "Add a performing name alias to an artist",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistAddAlias,
}

func init() {
	artistAddAliasCmd.Flags().StringVar(&artistAddAliasName, "name", "", "alias to add (required)")
	_ = artistAddAliasCmd.MarkFlagRequired("name")

	artistCmd.AddCommand(artistCreateCmd)
	artistCmd.AddCommand(artistListCmd)
	artistCmd.AddCommand(artistShowCmd)
	artistCmd.AddCommand(artistAddAliasCmd)
	rootCmd.AddCommand(artistCmd)
}

func runArtistCreate(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}

	displayName, err := prompt.Ask("Display name", "")
	if err != nil {
		return err
	}
	if displayName == "" {
		return errors.New("display name is required")
	}
	legalName, err := prompt.Ask("Legal name (leave blank to skip)", "")
	if err != nil {
		return err
	}
	defaultLicense, err := prompt.Ask("Default license", "CC BY 4.0")
	if err != nil {
		return err
	}
	distributor, err := prompt.Ask("Distributor (e.g. amuse, distrokid)", "")
	if err != nil {
		return err
	}
	pro, err := prompt.Ask("PRO (e.g. Teosto, ASCAP) — leave blank to skip", "")
	if err != nil {
		return err
	}
	ipiNumber, err := prompt.Ask("IPI number — leave blank to skip", "")
	if err != nil {
		return err
	}

	slug := ws.ToSlug(displayName)
	today := time.Now().Format("2006-01-02")

	artist := models.Artist{
		SchemaVersion:  "1",
		ID:             slug,
		DisplayName:    displayName,
		DefaultLicense: defaultLicense,
		CreatedDate:    today,
	}
	if legalName != "" {
		artist.LegalName = legalName
	}
	if distributor != "" {
		artist.Distribution = &models.ArtistDistrib{Distributor: distributor}
	}
	if pro != "" || ipiNumber != "" {
		artist.Rights = &models.ArtistRights{PRO: pro, IPINumber: ipiNumber}
	}

	raw, err := json.Marshal(artist)
	if err != nil {
		return err
	}
	errs, err := schema.Validate("artist", raw)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: "SCHEMA_VALIDATION_FAILED:\n" + strings.Join(errs, "\n")}
	}

	artistPath := ws.ArtistFile(wsPath, slug)
	if err := ws.WriteJSON(artistPath, artist); err != nil {
		return err
	}

	output.Success("Created %s", artistPath)
	fmt.Printf("Artist ID: %s\n", slug)
	return nil
}

func runArtistList(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	ids, err := ws.ListArtistIDs(wsPath)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		fmt.Println("No artists found. Run 'ovl artist create' to add one.")
		return nil
	}

	type row struct {
		id          string
		displayName string
		distributor string
		license     string
	}
	var rows []row
	for _, id := range ids {
		var a models.Artist
		if err := ws.ReadJSON(ws.ArtistFile(wsPath, id), &a); err != nil {
			continue
		}
		dist := ""
		if a.Distribution != nil {
			dist = a.Distribution.Distributor
		}
		rows = append(rows, row{id, a.DisplayName, dist, a.DefaultLicense})
	}

	if output.IsJSON() {
		var artists []models.Artist
		for _, id := range ids {
			var a models.Artist
			if err := ws.ReadJSON(ws.ArtistFile(wsPath, id), &a); err == nil {
				artists = append(artists, a)
			}
		}
		return output.JSON(artists)
	}

	tableRows := make([][]string, len(rows))
	for i, r := range rows {
		tableRows[i] = []string{r.id, r.displayName, r.distributor, r.license}
	}
	output.Table([]string{"ID", "Display Name", "Distributor", "License"}, tableRows)
	return nil
}

func runArtistShow(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	var a models.Artist
	if err := ws.ReadJSON(ws.ArtistFile(wsPath, args[0]), &a); err != nil {
		return fmt.Errorf("artist %q not found", args[0])
	}

	if output.IsJSON() {
		return output.JSON(a)
	}

	fmt.Printf("Artist: %s (%s)\n", a.DisplayName, a.ID)
	if a.LegalName != "" {
		fmt.Printf("Legal name: %s\n", a.LegalName)
	}
	if len(a.AlsoKnownAs) > 0 {
		fmt.Printf("Also known as: %s\n", strings.Join(a.AlsoKnownAs, ", "))
	}
	fmt.Printf("License: %s\n", a.DefaultLicense)
	if a.Distribution != nil && a.Distribution.Distributor != "" {
		fmt.Printf("Distributor: %s\n", a.Distribution.Distributor)
	}
	if a.Rights != nil {
		if a.Rights.PRO != "" {
			fmt.Printf("PRO: %s\n", a.Rights.PRO)
		}
		if a.Rights.IPINumber != "" {
			fmt.Printf("IPI: %s\n", a.Rights.IPINumber)
		}
	}
	if a.CreatedDate != "" {
		fmt.Printf("Created: %s\n", a.CreatedDate)
	}
	return nil
}

func runArtistAddAlias(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	artistID := args[0]
	path := ws.ArtistFile(wsPath, artistID)

	var a models.Artist
	if err := ws.ReadJSON(path, &a); err != nil {
		return fmt.Errorf("artist %q not found", artistID)
	}

	for _, existing := range a.AlsoKnownAs {
		if strings.EqualFold(existing, artistAddAliasName) {
			return fmt.Errorf("alias %q already exists", artistAddAliasName)
		}
	}

	a.AlsoKnownAs = append(a.AlsoKnownAs, artistAddAliasName)

	raw, _ := json.Marshal(a)
	if errs, err := schema.Validate("artist", raw); err != nil {
		return err
	} else if len(errs) > 0 {
		return &ExitError{Code: 3, Msg: strings.Join(errs, "\n")}
	}

	if err := ws.WriteJSON(path, a); err != nil {
		return err
	}
	output.Success("Added alias %q to %s", artistAddAliasName, artistID)
	return nil
}
