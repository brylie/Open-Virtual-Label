package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Manage website content sync",
	Long: `Register and sync websites that mirror OVL workspace records.

The OVL workspace is the canonical source of truth. Websites are
secondary consumers — their content collections are populated by
running 'ovl site sync'.`,
}

// site add
var (
	siteAddArtistID     string
	siteAddArtistsDir   string
	siteAddReleasesDir  string
	siteAddDescription  string
)

var siteAddCmd = &cobra.Command{
	Use:   "add <id> <path>",
	Short: "Register a website target",
	Long: `Register a website whose content collections will be populated by 'ovl site sync'.

<id>   Slug used to target this site with --site. e.g. label-site, aria-nova-site.
<path> Path to the site root. Relative to the workspace directory. e.g. sites/my-label.com.

Use --artist to scope sync to a single artist (artist sites). Omit for a
label-wide site that includes all artists.

Use --artists-dir / --releases-dir to override the destination directories
within the site when the defaults conflict with existing collections.`,
	Args: cobra.ExactArgs(2),
	RunE: runSiteAdd,
}

// site list
var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered website targets",
	RunE:  runSiteList,
}

// site remove
var siteRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Unregister a website target",
	Args:  cobra.ExactArgs(1),
	RunE:  runSiteRemove,
}

// site sync
var siteSyncTarget string

var siteSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync workspace records to website content collections",
	Long: `Copy artist and release JSON records into each registered website's
Astro content collection directories.

Without --site, all registered sites are synced. With --site <id>, only
that site is synced.

For sites registered with --artist, only that artist's records are copied.
For label-wide sites (no --artist), all artists and releases are synced.

After syncing, rebuild the site to pick up the changes.`,
	RunE: runSiteSync,
}

func init() {
	siteAddCmd.Flags().StringVar(&siteAddArtistID, "artist", "", "restrict sync to this artist ID (for artist sites)")
	siteAddCmd.Flags().StringVar(&siteAddArtistsDir, "artists-dir", "", "override destination dir for artist JSON (default: src/content/artists)")
	siteAddCmd.Flags().StringVar(&siteAddReleasesDir, "releases-dir", "", "override destination dir for release JSON (default: src/content/releases)")
	siteAddCmd.Flags().StringVar(&siteAddDescription, "description", "", "human-readable label for this site")

	siteSyncCmd.Flags().StringVar(&siteSyncTarget, "site", "", "sync only this site ID (default: all sites)")

	siteCmd.AddCommand(siteAddCmd, siteListCmd, siteRemoveCmd, siteSyncCmd)
	rootCmd.AddCommand(siteCmd)
}

func runSiteAdd(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	siteID := args[0]
	sitePath := args[1]

	if siteID == "" {
		return errors.New("VALIDATION_FAILED: site id is required")
	}

	label, err := readLabel(wsPath)
	if err != nil {
		return err
	}

	// Reject duplicate IDs.
	for _, s := range label.Sites {
		if s.ID == siteID {
			return &ExitError{Code: 1, Msg: fmt.Sprintf("site %q is already registered; remove it first with 'ovl site remove %s'", siteID, siteID)}
		}
	}

	// Validate --artist if provided.
	if siteAddArtistID != "" {
		artistIDs, _ := ws.ListArtistIDs(wsPath)
		found := false
		for _, id := range artistIDs {
			if id == siteAddArtistID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("artist %q not found in workspace", siteAddArtistID)
		}
	}

	site := models.LabelSite{
		ID:          siteID,
		Path:        sitePath,
		ArtistID:    siteAddArtistID,
		ArtistsDir:  siteAddArtistsDir,
		ReleasesDir: siteAddReleasesDir,
		Description: siteAddDescription,
	}
	label.Sites = append(label.Sites, site)

	if err := writeLabel(wsPath, &label); err != nil {
		return err
	}

	kind := "label site (all artists)"
	if siteAddArtistID != "" {
		kind = fmt.Sprintf("artist site (%s only)", siteAddArtistID)
	}
	output.Success("Registered site %q → %s [%s]", siteID, sitePath, kind)
	return nil
}

func runSiteList(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	label, err := readLabel(wsPath)
	if err != nil {
		return err
	}

	if output.IsJSON() {
		return output.JSON(label.Sites)
	}

	if len(label.Sites) == 0 {
		fmt.Println("No sites registered. Use 'ovl site add' to register one.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tPATH\tARTIST\tARTISTS DIR\tRELEASES DIR\tDESCRIPTION")
	for _, s := range label.Sites {
		artist := "(all)"
		if s.ArtistID != "" {
			artist = s.ArtistID
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			s.ID, s.Path, artist,
			s.ResolvedArtistsDir(), s.ResolvedReleasesDir(),
			s.Description)
	}
	return w.Flush()
}

func runSiteRemove(_ *cobra.Command, args []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	siteID := args[0]

	label, err := readLabel(wsPath)
	if err != nil {
		return err
	}

	found := false
	filtered := label.Sites[:0]
	for _, s := range label.Sites {
		if s.ID == siteID {
			found = true
		} else {
			filtered = append(filtered, s)
		}
	}
	if !found {
		return fmt.Errorf("site %q not found", siteID)
	}

	label.Sites = filtered
	if err := writeLabel(wsPath, &label); err != nil {
		return err
	}
	output.Success("Removed site %q", siteID)
	return nil
}

func runSiteSync(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}

	label, err := readLabel(wsPath)
	if err != nil {
		return err
	}
	if len(label.Sites) == 0 {
		return &ExitError{
			Code: 2,
			Msg:  "NO_SITES_CONFIGURED: no sites registered; run 'ovl site add <id> <path>' first",
		}
	}

	// Filter to the requested site, or sync all.
	targets := label.Sites
	if siteSyncTarget != "" {
		targets = nil
		for _, s := range label.Sites {
			if s.ID == siteSyncTarget {
				targets = []models.LabelSite{s}
				break
			}
		}
		if len(targets) == 0 {
			return fmt.Errorf("site %q not found; run 'ovl site list' to see registered sites", siteSyncTarget)
		}
	}

	for _, site := range targets {
		if err := syncSite(wsPath, site); err != nil {
			return fmt.Errorf("syncing site %q: %w", site.ID, err)
		}
	}
	return nil
}

// syncSite copies workspace records into a single site's content collections.
func syncSite(wsPath string, site models.LabelSite) error {
	var siteDir string
	if filepath.IsAbs(site.Path) {
		siteDir = site.Path
	} else {
		siteDir = filepath.Join(wsPath, site.Path)
	}

	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		return &ExitError{
			Code: 2,
			Msg:  fmt.Sprintf("SITE_NOT_FOUND: directory does not exist: %s", siteDir),
		}
	}

	artistsOut := filepath.Join(siteDir, filepath.FromSlash(site.ResolvedArtistsDir()))
	releasesOut := filepath.Join(siteDir, filepath.FromSlash(site.ResolvedReleasesDir()))

	if err := os.MkdirAll(artistsOut, 0o755); err != nil {
		return fmt.Errorf("cannot create artists content dir: %w", err)
	}
	if err := os.MkdirAll(releasesOut, 0o755); err != nil {
		return fmt.Errorf("cannot create releases content dir: %w", err)
	}

	allArtistIDs, err := ws.ListArtistIDs(wsPath)
	if err != nil {
		return err
	}

	// If site has artist_id set, filter to that artist only.
	artistIDs := allArtistIDs
	if site.ArtistID != "" {
		artistIDs = nil
		for _, id := range allArtistIDs {
			if id == site.ArtistID {
				artistIDs = []string{id}
				break
			}
		}
		if len(artistIDs) == 0 {
			return fmt.Errorf("artist %q not found in workspace", site.ArtistID)
		}
	}

	label := site.ID
	if site.Description != "" {
		label = fmt.Sprintf("%s (%s)", site.ID, site.Description)
	}
	output.Print("Syncing → %s", label)

	var artists, releases int
	for _, artistID := range artistIDs {
		srcArtist := ws.ArtistFile(wsPath, artistID)
		dstArtist := filepath.Join(artistsOut, artistID+".json")
		if err := copyFile(srcArtist, dstArtist); err != nil {
			return fmt.Errorf("copying artist %s: %w", artistID, err)
		}
		output.Print("  artists/%s.json", artistID)
		artists++

		releaseIDs, err := ws.ListReleaseIDs(wsPath, artistID)
		if err != nil {
			return err
		}
		for _, releaseID := range releaseIDs {
			srcRelease := ws.ReleaseFile(wsPath, artistID, releaseID)
			dstName := artistID + "--" + releaseID + ".json"
			dstRelease := filepath.Join(releasesOut, dstName)
			if err := copyFile(srcRelease, dstRelease); err != nil {
				return fmt.Errorf("copying release %s/%s: %w", artistID, releaseID, err)
			}
			output.Print("  releases/%s", dstName)
			releases++
		}
	}

	output.Success("Synced %d artist(s), %d release(s) → %s", artists, releases, siteDir)
	return nil
}

// copyFile copies src to dst, creating parent directories as needed.
func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

