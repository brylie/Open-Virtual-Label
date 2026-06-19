package cmd

import (
	"fmt"

	"github.com/open-virtual-label/ovl/internal/models"
	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Package and upload releases to long-term archive storage",
}

var archivePushRelease string
var archiveSkipStems, archiveSkipProjectFiles bool
var archivePushCmd = &cobra.Command{
	Use:   "push",
	Short: "Package and upload a release to archive storage",
	RunE:  runArchivePush,
}

var archiveStatusRelease string
var archiveStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show archive status for a release",
	RunE:  runArchiveStatus,
}

func init() {
	archivePushCmd.Flags().StringVar(&archivePushRelease, "release", "", "release ID (required)")
	archivePushCmd.Flags().BoolVar(&archiveSkipStems, "skip-stems", false, "omit stems from archive package")
	archivePushCmd.Flags().BoolVar(&archiveSkipProjectFiles, "skip-project-files", false, "omit DAW project files")
	_ = archivePushCmd.MarkFlagRequired("release")

	archiveStatusCmd.Flags().StringVar(&archiveStatusRelease, "release", "", "release ID (required)")
	_ = archiveStatusCmd.MarkFlagRequired("release")

	archiveCmd.AddCommand(archivePushCmd, archiveStatusCmd)
	rootCmd.AddCommand(archiveCmd)
}

func runArchivePush(_ *cobra.Command, _ []string) error {
	fmt.Println("Archive push requires Internet Archive MCP integration.")
	fmt.Println("Run 'ovl mcp connect internet-archive' to configure it.")
	fmt.Println()
	fmt.Println("This command will:")
	fmt.Println("  1. Assemble a manifest of master WAVs, stems, project files, artwork, and JSON records")
	fmt.Println("  2. Present the manifest for review (this gate cannot be bypassed with --yes)")
	fmt.Println("  3. Upload to Internet Archive via the IA S3-compatible API")
	fmt.Println("  4. Verify SHA-256 checksums after upload")
	fmt.Println("  5. Write release.archive{} fields with IDs, URLs, and checksum status")
	return &ExitError{Code: 5, Msg: "IA_MCP_NOT_CONFIGURED"}
}

func runArchiveStatus(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	artistID, err := ws.FindRelease(wsPath, archiveStatusRelease)
	if err != nil {
		return err
	}

	var rel models.Release
	if err := ws.ReadJSON(ws.ReleaseFile(wsPath, artistID, archiveStatusRelease), &rel); err != nil {
		return err
	}

	if rel.Archive == nil {
		fmt.Printf("Release %s has not been archived yet.\n", archiveStatusRelease)
		fmt.Println("Run 'ovl archive push --release', then 'ovl mcp connect internet-archive'.")
		return nil
	}

	a := rel.Archive
	rows := [][]string{
		{"Masters archived", boolStr(a.MastersArchived)},
		{"Stems archived", boolStr(a.StemsArchived)},
		{"Project files archived", boolStr(a.ProjectFilesArchived)},
		{"Checksums verified", boolStr(a.ChecksumsVerified)},
	}
	if a.InternetArchiveID != nil && *a.InternetArchiveID != "" {
		rows = append(rows, []string{"Internet Archive ID", *a.InternetArchiveID})
	}
	if a.InternetArchiveURL != nil && *a.InternetArchiveURL != "" {
		rows = append(rows, []string{"Internet Archive URL", *a.InternetArchiveURL})
	}
	if a.ArchiveDate != nil && *a.ArchiveDate != "" {
		rows = append(rows, []string{"Archive date", *a.ArchiveDate})
	}
	output.Table([]string{"Field", "Value"}, rows)
	return nil
}

func boolStr(b bool) string {
	if b {
		return "✓"
	}
	return "✗"
}
