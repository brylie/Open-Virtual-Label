package cmd

import (
	"fmt"
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

func init() {
	agentsCmd.AddCommand(agentsListCmd)
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
