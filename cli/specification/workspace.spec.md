# CLI Specification — `ovl` workspace commands

Top-level commands for initialising and inspecting the workspace, plus the
`state` subgroup for managing the label state document.

---

## Command Reference

---

### `ovl init`

Scaffold a new label workspace from the built-in template.

```text
ovl init [--workspace <path>] [--force]
```

**Behaviour:**

- Copies `workspace-scaffold/` into the target directory
- Runs an interactive setup to populate `workspace/label/profile.json` with label name,
  contact email, default license, and primary distributor
- Creates an initial `workspace/state/label-state.md` with the first session log entry
- Fails if a `workspace/` directory already exists, unless `--force` is passed

**Options:**

| Option | Description |
|---|---|
| `--force` | Overwrite an existing workspace. Prompts for confirmation. |

**Output:** Workspace directory structure created. Summary of files written.

**Errors:**

- `WORKSPACE_EXISTS` — workspace directory already exists; use `--force` to overwrite
- `SCHEMA_VALIDATION_FAILED` — generated `profile.json` failed validation (indicates a
  CLI bug; report it)

---

### `ovl status`

Display current label state: active projects, open loops, and pending approvals.

```text
ovl status [--artist <artist-id>]
```

**Behaviour:**

- Reads `workspace/state/label-state.md`
- Invokes the orchestrator agent with the current state as context
- Outputs a formatted summary of active releases, open loops, and pending approvals
- If `--artist` is specified, filters to that artist's releases and opportunities

**Output example:**

```text
Last session: 2025-06-14 — mastered tracks 1–3 of Spectra

Active releases (1):
  Spectra [brylie-christopher] — mastering (3 of 8 tracks)

Pending approvals (1):
  outreach-crm: draft to Calm Waters Podcast

Open loops (2):
  QC not yet run on Spectra
  Follow-up due: Calm Waters Podcast by 2025-06-28
```

---

### `ovl validate`

Validate workspace records against their schemas.

```text
ovl validate [<path>] [--all]
```

**Behaviour:**

- With a path argument: validates the single file at that path against its inferred
  schema (based on directory location and file naming convention)
- With `--all`: validates every JSON file in `workspace/` against the appropriate schema
- Reports each failure with the file path, field name, and constraint violated
- Exits with code 0 if all records pass, code 1 if any fail

**Options:**

| Option | Description |
|---|---|
| `--all` | Validate every record in the workspace |

**Output:**

```text
✓ workspace/label/profile.json
✓ workspace/artists/brylie-christopher/artist.json
✗ workspace/artists/brylie-christopher/releases/spectra/tracks/chromatic-drift.json
    mastering.true_peak_dbtp: must be number, got null
    qc.passed: required field missing
2 errors found.
```

---

### `ovl state show`

Print the full contents of `workspace/state/label-state.md`.

```text
ovl state show
```

**Output:** Raw markdown content of the state document printed to stdout.

---

### `ovl state sync`

Invoke the orchestrator to write a session summary to `label-state.md`.

```text
ovl state sync
```

**Behaviour:**

- Loads the current state document and recent workspace changes
- Invokes the orchestrator agent to produce a session summary
- Presents the summary for artist review
- **[CONFIRMATION GATE]** Writes to `label-state.md` only on approval
