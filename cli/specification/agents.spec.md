# CLI Specification — `ovl agents`

Commands for inspecting and distributing agent skills. Agent skills are defined
as `SKILL.md` files under `.agents/skills/` and loaded by the CLI when a command
requires agent handoff. The `package` subcommand converts those skills into
installable files for AI toolkits that do not natively read the `.agents/skills/`
directory format.

---

## Command Reference

---

### `ovl agents list`

List all installed agent skills and their interaction patterns.

```text
ovl agents list
```

**Behaviour:**

- Scans `.agents/skills/*/SKILL.md` and reads YAML frontmatter from each
- Checks whether each agent's required MCPs are connected

**Output:** Table of agent names, descriptions, interaction patterns, and connection
status (whether required MCPs are configured).

**Output example:**

```text
Agent                 Pattern            Status
coordinator           review-and-refine  ready
mastering-companion   guided-session     ready
outreach-crm          approval-gate      ready (gmail MCP connected)
qc                    review-and-refine  ready
finance-manager       review-and-refine  ready
```

---

### `ovl agents package`

Package one or all agent skills as installable files for a target AI toolkit.

```text
ovl agents package [agent-id] [--target <toolkit>] [--output <dir>]
```

**Behaviour:**

- Reads skill directories from `.agents/skills/`. Directories prefixed with `_`
  (e.g. `_ovl-template`) are skipped when packaging all agents.
- For each skill, zips the entire skill directory (preserving sub-paths such as
  `references/`) into a single output file. `SKILL.md` is always at the zip root.
- Writes packaged files to `--output` (default: current directory).
- After packaging, prints step-by-step install instructions for the target toolkit.

If `[agent-id]` is omitted, all non-template skills are packaged.

**Options:**

| Option | Default | Description |
|---|---|---|
| `--target` | `cowork` | AI toolkit to package for. See Supported Targets below. |
| `--output` | `.` | Directory to write the packaged files into. Created if it does not exist. |

**Supported targets:**

| Target | Output format | Install method |
|---|---|---|
| `cowork` | `<agent-id>.skill` (zip) | Claude desktop → Settings → Capabilities → Install Skill |

**Output example:**

```text
✓ Packaged ovl-coordinator → ./ovl-coordinator.skill

To install in Claude desktop (Cowork mode):
  1. Open Claude desktop and go to Settings → Capabilities
  2. Click "Install Skill" (or drag the .skill file into the window)
  3. Select the packaged file(s):
       ./ovl-coordinator.skill
  4. The agent will appear in the / menu in your next conversation.
```

**Errors:**

- `AGENT_NOT_FOUND` — the specified `[agent-id]` does not match any directory
  under `.agents/skills/`
- `UNSUPPORTED_TARGET` — `--target` value is not in the supported targets list
