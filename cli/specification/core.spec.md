# CLI Specification — Core

Cross-cutting concerns that apply to the entire `ovl` CLI: design principles,
installation, global options, exit codes, environment variables, and implementation notes.

---

## Design Principles

**Reads and writes workspace JSON, nothing else.** The CLI operates on files in
`workspace/`. It does not send emails, post content, submit to distributors, or upload
to archive services without an explicit confirmation prompt.

**Agents are invoked, not embedded.** When a command requires conversational guidance
(mastering, outreach drafting, QC review), the CLI loads the appropriate agent skill and
workspace context, then hands off. The agent's session output is written back to the
workspace records on completion.

**Every destructive or external action has a confirmation gate.** Commands that write to
external services, advance a release to a new status, or send messages on the artist's
behalf prompt for explicit confirmation before proceeding. Flags like `--yes` or
`--force` are available for scripting but must be documented and intentional.

**Validation is automatic.** Every write to a workspace record validates the result
against its schema before saving. A record that would fail schema validation is not
written; the error is reported with the field and constraint that failed.

**Single binary, no runtime required.** The CLI ships as a self-contained executable.
Users install it by placing the binary on their `PATH`; no language runtime or package
manager is needed.

---

## Installation

Download the appropriate binary for your platform from the project releases page and
place it on your `PATH`. Verify:

```bash
ovl --version
```

---

## Global Options

Available on all commands:

| Option | Description |
|---|---|
| `--workspace <path>` | Path to workspace directory. Defaults to `./workspace` relative to cwd, then walks up the directory tree. |
| `--schemas <path>` | Path to the JSON schemas directory. Defaults to `./schemas` relative to cwd, then walks up the directory tree. |
| `--artist <artist-id>` | Scopes the command to a specific artist when multiple artists exist in the workspace. |
| `--json` | Output result as JSON instead of formatted text. Useful for scripting. |
| `--quiet` | Suppress informational output. Errors still print to stderr. |
| `--yes` | Skip confirmation prompts, accepting the default action. Use with care in scripts. |
| `--help` | Print help for the command. |
| `--version` | Print the installed `ovl` version. |

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | General error (validation failure, missing required argument, etc.) |
| `2` | Workspace not found |
| `3` | Schema validation failure |
| `4` | Confirmation declined by user |
| `5` | MCP not configured for requested operation |
| `6` | External service error (upload failed, API error, etc.) |
| `7` | Artist not found |

---

## Environment Variables

| Variable | Description |
|---|---|
| `OVL_WORKSPACE` | Override the workspace path. Equivalent to `--workspace <path>` on every command. |
| `OVL_SCHEMAS_DIR` | Override the schemas directory path. Equivalent to `--schemas <path>` on every command. |
| `OVL_ARTIST` | Override the default artist. Equivalent to `--artist <id>` on every command. |
| `OVL_YES` | If set to `1`, skip non-critical confirmation prompts. Critical gates (release submit, archive push) are never bypassed. |
| `OVL_JSON` | If set to `1`, output all results as JSON. |

---

## Command Groups Summary

```text
Workspace:   init · status · validate · state show · state sync
Label:       label show · label set-name · label set-description · label set-license
             label set-distributor · label set-contact
Artist:      artist create · artist list · artist show · artist add-alias
             artist set-bio · artist set-contact · artist set-platform
             artist set-rights · artist set-distributor · artist set-license
             artist set-location
Release:     release create · release list · release show · release advance
             release set-profile · release set-live · release add-link · release submit
Track:       track add · track show · track set-file
ISRC:        isrc assign
Mastering:   mastering start · mastering profile create · mastering profile list
QC:          qc check
Archive:     archive push · archive status
Outreach:    outreach research · outreach review · outreach draft · outreach send
             outreach follow-up · outreach log-response · outreach score
             outreach intake · outreach close · outreach log
Finance:     finance add-revenue · finance add-expense · finance summary · finance quote
Metrics:     metrics snapshot
Content:     content brief
Social:      social draft
Agents:      agents list
MCP:         mcp list · mcp connect · mcp disconnect
Commission:  commission agreement
```

---

## Implementation Notes

**Argument parsing:** Subcommand-based, with global flags inherited by all commands.
Interactive prompts are used for required fields not supplied as flags. Long-running
operations (archive uploads, agent sessions) stream output rather than blocking.

**Schema validation:** JSON Schema Draft 7. Every workspace write validates against the
appropriate schema before the file is touched. Validation errors report the field path
and the constraint violated.

**Frontmatter parsing:** Agent skill files (`SKILL.md`) carry YAML frontmatter that the
CLI reads to build the agent registry (name, description, interaction pattern, required
integrations).

**Agent invocation:** When a command hands off to an agent, the CLI loads the agent's
`SKILL.md` and relevant workspace records, then opens a session in the current terminal
(or via a configured API if running non-interactively). The specific integration
mechanism depends on the deployment context and is specified separately in
`cli/AGENT-INTEGRATION.md`.

**Config file:** The CLI reads from `.ovlrc` in the workspace root for persistent local
configuration (preferred artist, output format, integration preferences). This file is
`.gitignored` alongside `workspace/`.
