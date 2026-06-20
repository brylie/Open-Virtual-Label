package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-virtual-label/ovl/internal/output"
	ws "github.com/open-virtual-label/ovl/internal/workspace"
	"github.com/spf13/cobra"
)

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "List and inspect installed agent skills",
}

var agentsListCmd = &cobra.Command{
	Use:   cmdUseList,
	Short: "List all installed agent skills",
	RunE:  runAgentsList,
}

var (
	packageTarget string
	packageOutput string
)

var agentsPackageCmd = &cobra.Command{
	Use:   "package [agent-id]",
	Short: "Package agent skills for installation in an AI toolkit",
	Long: `Package one or all agent skills as installable files for a target AI toolkit.

Supported targets:
  cowork      Claude desktop app (Cowork mode) — produces <agent-id>.skill files

If [agent-id] is omitted, all skills in .agents/skills/ are packaged.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAgentsPackage,
}

func init() {
	agentsPackageCmd.Flags().StringVar(&packageTarget, "target", "cowork",
		"AI toolkit to package for (cowork)")
	agentsPackageCmd.Flags().StringVar(&packageOutput, "output", ".",
		"directory to write packaged files into")

	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsPackageCmd)
	rootCmd.AddCommand(agentsCmd)
}

func runAgentsList(_ *cobra.Command, _ []string) error {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	// Look for .agents/skills/ relative to workspace root
	root := ws.Root(wsPath)
	skillsDir := filepath.Join(root, ".agents", "skills")

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		fmt.Println("No agent skills directory found at", skillsDir)
		fmt.Println("Expected: .agents/skills/<agent-name>/SKILL.md")
		return nil //nolint:nilerr // missing skills dir is not an error; print guidance and exit cleanly
	}

	type agentInfo struct {
		id          string
		description string
	}
	var agents []agentInfo

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillFile := filepath.Join(skillsDir, e.Name(), "SKILL.md")
		desc := parseSkillDescription(skillFile)
		agents = append(agents, agentInfo{id: e.Name(), description: desc})
	}

	if len(agents) == 0 {
		fmt.Println("No agent skills found.")
		return nil
	}

	rows := make([][]string, len(agents))
	for i, a := range agents {
		rows[i] = []string{a.id, a.description}
	}
	output.Table([]string{"Agent ID", "Description"}, rows)
	return nil
}

func runAgentsPackage(_ *cobra.Command, args []string) error {
	if packageTarget != "cowork" {
		return fmt.Errorf("unsupported target %q — currently supported: cowork", packageTarget)
	}

	wsPath, err := resolveWorkspace()
	if err != nil {
		return err
	}
	root := ws.Root(wsPath)
	skillsDir := filepath.Join(root, ".agents", "skills")

	// Collect agent IDs to package.
	var agentIDs []string
	if len(args) == 1 {
		agentIDs = []string{args[0]}
	} else {
		entries, err := os.ReadDir(skillsDir)
		if err != nil {
			return fmt.Errorf("no agent skills directory found at %s", skillsDir)
		}
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), "_") {
				agentIDs = append(agentIDs, e.Name())
			}
		}
		if len(agentIDs) == 0 {
			fmt.Println("No agent skills found.")
			return nil
		}
	}

	outDir, err := filepath.Abs(packageOutput)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("cannot create output directory: %w", err)
	}

	var packaged []string
	for _, id := range agentIDs {
		skillDir := filepath.Join(skillsDir, id)
		if info, err := os.Stat(skillDir); err != nil || !info.IsDir() {
			return fmt.Errorf("agent %q not found at %s", id, skillDir)
		}
		outFile := filepath.Join(outDir, id+".skill")
		if err := zipSkillDir(skillDir, outFile); err != nil {
			return fmt.Errorf("packaging %s: %w", id, err)
		}
		output.Success("Packaged %s → %s", id, outFile)
		packaged = append(packaged, outFile)
	}

	printCoworkInstallInstructions(packaged)
	return nil
}

// zipSkillDir zips the contents of skillDir into destFile.
// Files are stored with paths relative to skillDir (SKILL.md at the zip root).
func zipSkillDir(skillDir, destFile string) error {
	f, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	return filepath.Walk(skillDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(skillDir, path)
		if err != nil {
			return err
		}
		// Use forward slashes inside the zip regardless of OS.
		rel = filepath.ToSlash(rel)

		w, err := zw.Create(rel)
		if err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(w, src)
		return err
	})
}

func printCoworkInstallInstructions(files []string) {
	fmt.Println()
	fmt.Println("To install in Claude desktop (Cowork mode):")
	fmt.Println("  1. Open Claude desktop and go to Settings → Capabilities")
	fmt.Println("  2. Click \"Install Skill\" (or drag the .skill file into the window)")
	fmt.Println("  3. Select the packaged file(s):")
	for _, f := range files {
		fmt.Println("       " + f)
	}
	fmt.Println("  4. The agent will appear in the / menu in your next conversation.")
}

// parseSkillDescription extracts the description from a SKILL.md frontmatter block.
func parseSkillDescription(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "(no SKILL.md)"
	}
	lines := strings.Split(string(data), "\n")
	inFrontmatter := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter && strings.HasPrefix(line, "description:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		}
	}
	// Fallback: first non-empty non-heading line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "---") {
			if len(line) > 80 {
				line = line[:77] + "..."
			}
			return line
		}
	}
	return ""
}
