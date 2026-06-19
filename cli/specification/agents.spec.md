# CLI Specification — `ovl agents`

Commands for inspecting the installed agent skill registry. Agent skills are defined
as `SKILL.md` files under `.agents/skills/` and loaded by the CLI when a command
requires agent handoff.

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
