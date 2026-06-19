# CLI Specification — `ovl`

The `ovl` command-line tool is the primary interaction surface for the Open Virtual Label system. It manages workspace records, validates schemas, invokes agent sessions, and provides a consistent interface to every step of the release pipeline and outreach loop.

---

## Design Principles

**Reads and writes workspace JSON, nothing else.** The CLI operates on files in `workspace/`. It does not send emails, post content, submit to distributors, or upload to archive services without an explicit confirmation prompt.

**Agents are invoked, not embedded.** When a command requires conversational guidance (mastering, outreach drafting, QC review), the CLI loads the appropriate agent skill and workspace context, then hands off. The agent's session output is written back to the workspace records on completion.

**Every destructive or external action has a confirmation gate.** Commands that write to external services, advance a release to a new status, or send messages on the artist's behalf prompt for explicit confirmation before proceeding. Flags like `--yes` or `--force` are available for scripting but must be documented and intentional.

**Validation is automatic.** Every write to a workspace record validates the result against its schema before saving. A record that would fail schema validation is not written; the error is reported with the field and constraint that failed.

**Single binary, no runtime required.** The CLI ships as a self-contained executable. Users install it by placing the binary on their `PATH`; no language runtime or package manager is needed.

---

## Installation

Download the appropriate binary for your platform from the project releases page and place it on your `PATH`. Verify:

```bash
ovl --version
```

---

## Global Options

Available on all commands:

| Option | Description |
|---|---|
| `--workspace <path>` | Path to workspace directory. Defaults to `./workspace` relative to cwd, then walks up the directory tree. |
| `--artist <artist-id>` | Scopes the command to a specific artist when multiple artists exist in the workspace. |
| `--json` | Output result as JSON instead of formatted text. Useful for scripting. |
| `--quiet` | Suppress informational output. Errors still print to stderr. |
| `--yes` | Skip confirmation prompts, accepting the default action. Use with care in scripts. |
| `--help` | Print help for the command. |
| `--version` | Print the installed `ovl` version. |

---

## Command Reference

Commands are grouped by domain. Within each group, commands are listed in the order they are typically used.

---

### Workspace

#### `ovl init`

Scaffold a new label workspace from the built-in template.

```text
ovl init [--workspace <path>] [--force]
```

**Behaviour:**

- Copies `workspace-scaffold/` into the target directory
- Runs an interactive setup to populate `workspace/label/profile.json` with label name, contact email, default license, and primary distributor
- Creates an initial `workspace/state/label-state.md` with the first session log entry
- Fails if a `workspace/` directory already exists, unless `--force` is passed

**Options:**

| Option | Description |
|---|---|
| `--force` | Overwrite an existing workspace. Prompts for confirmation. |

**Output:** Workspace directory structure created. Summary of files written.

**Errors:**

- `WORKSPACE_EXISTS` — workspace directory already exists; use `--force` to overwrite
- `SCHEMA_VALIDATION_FAILED` — generated `profile.json` failed validation (indicates a CLI bug; report it)

---

#### `ovl status`

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

#### `ovl validate`

Validate workspace records against their schemas.

```text
ovl validate [<path>] [--all]
```

**Behaviour:**

- With a path argument: validates the single file at that path against its inferred schema (based on directory location and file naming convention)
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

#### `ovl state show`

Print the full contents of `workspace/state/label-state.md`.

```text
ovl state show
```

**Output:** Raw markdown content of the state document printed to stdout.

---

#### `ovl state sync`

Invoke the orchestrator to write a session summary to `label-state.md`.

```text
ovl state sync
```

**Behaviour:**

- Loads the current state document and recent workspace changes
- Invokes the orchestrator agent to produce a session summary
- Presents the summary for artist review
- **[CONFIRMATION GATE]** Writes to `label-state.md` only on approval

---

### Artist

#### `ovl artist create`

Create a new artist profile interactively.

```text
ovl artist create
```

**Behaviour:**

- Prompts for display name, legal name, default license, PRO, IPI number, distributor, and platform links
- Generates a slug from the display name (e.g. `brylie-christopher`)
- Creates `workspace/artists/[slug]/artist.json`
- Validates against `schemas/artist.schema.json` before writing

**Output:** Path to the created `artist.json` and the generated artist ID.

---

#### `ovl artist list`

List all artist profiles in the workspace.

```text
ovl artist list
```

**Output:** Table of artist IDs, display names, and distributor.

---

#### `ovl artist show <artist-id>`

Display an artist profile.

```text
ovl artist show <artist-id>
```

**Output:** Formatted display of `artist.json` fields.

---

#### `ovl artist add-alias <artist-id> --name "<alias>"`

Add a performing name alias to an existing artist profile.

```text
ovl artist add-alias <artist-id> --name "<alias>"
```

**Behaviour:** Appends to `artist.also_known_as[]` and saves.

---

### Release

#### `ovl release create "<title>"`

Create a new release record.

```text
ovl release create "<title>" [--artist <artist-id>] [--type album|ep|single|compilation]
```

**Behaviour:**

- Generates a slug from the title (e.g. `spectra`)
- Creates `workspace/artists/[artist-id]/releases/[slug]/release.json` with `status: in-production`
- Creates the `tracks/` subdirectory
- If `--artist` is omitted and only one artist exists in the workspace, uses that artist. If multiple artists exist, prompts.
- Prompts for target release date and distributor submission deadline

**Options:**

| Option | Default | Description |
|---|---|---|
| `--artist <id>` | (prompted if ambiguous) | Artist this release belongs to |
| `--type <type>` | `album` | Release type: `album`, `ep`, `single`, `compilation` |

**Output:** Path to created `release.json` and the generated release ID.

---

#### `ovl release list`

List all releases across all artists.

```text
ovl release list [--artist <artist-id>] [--status <status>]
```

**Options:**

| Option | Description |
|---|---|
| `--artist <id>` | Filter to one artist |
| `--status <status>` | Filter by release status (e.g. `--status mastering`) |

**Output:** Table of release IDs, titles, artists, types, statuses, and target dates.

---

#### `ovl release show <release-id>`

Display a release record and its track list with statuses.

```text
ovl release show <release-id>
```

**Output:** Formatted release details plus a table of tracks showing mastering and QC status.

---

#### `ovl release advance <release-id> --status <status>`

Manually advance a release to the next pipeline status.

```text
ovl release advance <release-id> --status <status>
```

**Behaviour:**

- Validates that the advance is legitimate (e.g. cannot advance to `ready` if `qc.passed` is not `true`)
- **[CONFIRMATION GATE]** Prompts before writing the status change
- Updates `release.status` in `release.json`

**Valid transitions:**

| From | To | Requires |
|---|---|---|
| `in-production` | `mastering` | At least one track exists |
| `mastering` | `qc` | All tracks have mastering data |
| `qc` | `ready` | `release.qc.passed: true` |
| `ready` | `submitted` | Via `ovl release submit` only |
| `submitted` | `live` | Via `ovl release set-live` only |

---

#### `ovl release set-profile <release-id> --profile <profile-id>`

Assign a mastering profile to a release.

```text
ovl release set-profile <release-id> --profile <profile-id>
```

**Behaviour:** Sets `release.mastering_profile_id`. The profile must exist in `workspace/artists/[artist-id]/mastering-profiles/`.

---

#### `ovl release set-live <release-id> --date <YYYY-MM-DD>`

Mark a submitted release as live.

```text
ovl release set-live <release-id> --date <YYYY-MM-DD>
```

**Behaviour:** Sets `release.status: live` and `release.dates.released`. Prompts for confirmation.

---

#### `ovl release add-link <release-id> --platform <platform> --url <url>`

Add a store link to a live release.

```text
ovl release add-link <release-id> --platform <platform> --url <url>
```

**Behaviour:** Sets the specified field in `release.store_links`. Valid platforms: `spotify`, `apple_music`, `youtube_music`, `bandcamp`, `soundcloud`, `tidal`, `amazon_music`.

---

#### `ovl release submit <release-id> --distributor <distributor>`

Prepare and submit a distribution package.

```text
ovl release submit <release-id> --distributor <distributor>
```

**Behaviour:**

- Assembles a submission summary: title, type, target date, track list with ISRCs and durations, artwork details, license, contributor credits
- **[CONFIRMATION GATE]** Presents the full summary for artist review. This gate cannot be bypassed with `--yes`
- On confirmation, either generates the distributor's required submission format for manual upload, or calls the distributor's API if an MCP is configured
- Sets `release.status: submitted` and `release.dates.submitted`

**Errors:**

- `QC_NOT_PASSED` — release has not passed QC; run `ovl qc check` first
- `ARCHIVE_NOT_COMPLETE` — masters not archived; run `ovl archive push` first
- `MISSING_ISRC` — one or more tracks are missing ISRC; assign before submitting

---

### Track

#### `ovl track add "<title>" --release <release-id>`

Add a track to a release.

```text
ovl track add "<title>" --release <release-id> [--position <n>]
```

**Behaviour:**

- Generates a slug from the title
- Creates `workspace/artists/[artist-id]/releases/[release-id]/tracks/[slug].json`
- If `--position` is omitted, appends after the last existing track

**Options:**

| Option | Description |
|---|---|
| `--position <n>` | Track number (1-based). If omitted, appends at end. |

---

#### `ovl track show <track-id>`

Display a track record.

```text
ovl track show <track-id> [--release <release-id>]
```

**Behaviour:** Track IDs are unique within a release. If the same slug exists in multiple releases, `--release` is required to disambiguate.

---

#### `ovl track set-file <track-id> --field <field> --path <path>`

Set a file path on a track record.

```text
ovl track set-file <track-id> --field <field> --path <path>
```

**Valid fields:** `master_wav`, `stems_zip`, `project_file`, `mp3_320`, `wav_for_distribution`

**Behaviour:** Validates that the referenced path exists before writing. Path is stored relative to the release directory.

---

#### `ovl isrc assign --release <release-id>`

Assign ISRCs to tracks that do not yet have one.

```text
ovl isrc assign --release <release-id> [--track <track-id>]
```

**Behaviour:**

- Lists all tracks in the release with `isrc: null`
- For each, prompts for the ISRC (format: `CC-XXX-YY-NNNNN`)
- Validates format against the schema pattern before saving
- If `--track` is specified, assigns only to that track

**Note:** OVL does not register ISRCs. They are issued by the artist's distributor or national ISRC agency. This command records codes that have already been obtained.

---

### Mastering

#### `ovl mastering start --track <track-id>`

Launch a mastering session with the Mastering Companion agent.

```text
ovl mastering start --track <track-id> [--remaster]
```

**Behaviour:**

- Loads the track record and the applicable mastering profile
- Invokes the `mastering-companion` agent in guided-session mode
- On session completion, writes mastering measurements to `track.mastering{}`
- Appends session notes to `mastering_profile.session_notes[]`

**Options:**

| Option | Description |
|---|---|
| `--remaster` | Re-master a track that already has mastering data. Previous measurements are preserved in `track.mastering.notes` before being overwritten. |

---

#### `ovl mastering profile create`

Create a new mastering profile interactively.

```text
ovl mastering profile create [--artist <artist-id>]
```

**Behaviour:**

- Prompts for profile name, LUFS targets, true peak ceiling, LRA guidance, sample rate, bit depth
- Offers platform-specific note prompts for Spotify, Apple Music, YouTube
- Generates a slug from the name
- Creates `workspace/artists/[artist-id]/mastering-profiles/[slug].json`

---

#### `ovl mastering profile list`

List mastering profiles for an artist.

```text
ovl mastering profile list [--artist <artist-id>]
```

**Output:** Table of profile IDs, names, LUFS targets, and session count.

---

### QC

#### `ovl qc check --release <release-id>`

Run the pre-release quality check on a release.

```text
ovl qc check --release <release-id>
```

**Behaviour:**

- Invokes the `qc` agent
- Checks all track-level requirements (ISRC, mastering data completeness, file paths) and release-level requirements (artwork, metadata, dates, license)
- Presents the full QC report
- **[CONFIRMATION GATE]** On a clean pass: asks artist to confirm before setting `release.qc.passed: true`
- On failures: lists each failure with the field and requirement; artist must resolve or record an override before QC can pass

**Output example:**

```text
QC Report — Spectra

Tracks (8/8 checked):
  ✓ chromatic-drift     ISRC ✓  Mastering ✓  File ✓
  ✓ passage             ISRC ✓  Mastering ✓  File ✓
  ✗ luminance           ISRC ✗  (null — assign before submitting)
  ...

Release:
  ✓ Artwork: 3000×3000px PNG
  ✓ License: CC BY 4.0
  ✓ Target date: 2025-09-01
  ✗ Submission deadline: not set

2 failures. Resolve before advancing to ready.
```

---

### Archive

#### `ovl archive push --release <release-id>`

Package and upload a release to long-term archive storage.

```text
ovl archive push --release <release-id> [--skip-stems] [--skip-project-files]
```

**Behaviour:**

- Assembles a manifest of all files to upload: master WAVs, stems, project files, artwork, and JSON records
- **[CONFIRMATION GATE]** Presents the manifest for artist review. Lists each file, its size, and destination. Cannot be bypassed with `--yes`
- Uploads to Internet Archive (primary) via the IA S3-compatible API
- Uploads to secondary object storage if configured in `workspace/label/profile.json`
- Verifies SHA-256 checksums after upload
- Writes `release.archive{}` fields: IDs, URLs, paths, flags, checksums verified, archive date

**Options:**

| Option | Description |
|---|---|
| `--skip-stems` | Omit stems from the archive package. |
| `--skip-project-files` | Omit DAW project files from the archive package. |

**Errors:**

- `IA_MCP_NOT_CONFIGURED` — Internet Archive MCP not connected; run `ovl mcp connect internet-archive`
- `FILE_NOT_FOUND` — a file referenced in a track record does not exist at the stated path
- `UPLOAD_PARTIAL` — upload interrupted mid-way; re-run to resume (checksums prevent re-uploading completed files)

---

#### `ovl archive status --release <release-id>`

Show archive status for a release.

```text
ovl archive status --release <release-id>
```

**Output:** Table of archive flags, URLs, and checksum verification status from `release.archive{}`.

---

### Outreach / CRM

#### `ovl outreach research`

Trigger the CRM agent to find new opportunities.

```text
ovl outreach research [--type <type>] [--release <release-id>] [--track <track-id>]
```

**Behaviour:**

- Invokes the `outreach-crm` agent in research mode
- Agent searches for opportunities matching the artist's genre tags and platform presence
- Presents a list of candidates with match scores for artist review
- **[REVIEW GATE]** Artist selects which opportunities to pursue
- Creates `opportunity.json` records for approved candidates with `status: identified`

**Options:**

| Option | Description |
|---|---|
| `--type <type>` | Limit to one opportunity type: `sync-license`, `commission`, `playlist-pitch`, `collaboration`, `press` |
| `--release <release-id>` | Focus research on placement opportunities for a specific release |
| `--track <track-id>` | Focus on a specific track (e.g. for playlist pitching) |

---

#### `ovl outreach review`

Review all opportunities with pending approvals.

```text
ovl outreach review [--type <type>]
```

**Behaviour:**

- Lists all opportunities with `status: draft-ready` (outreach drafts awaiting approval)
- Invokes the `outreach-crm` agent to present each draft for review
- Artist approves, edits, or declines each draft
- Approved drafts advance to `status: approved`

---

#### `ovl outreach draft --opportunity <opportunity-id>`

Draft outreach for a specific opportunity.

```text
ovl outreach draft --opportunity <opportunity-id>
```

**Behaviour:**

- Invokes the `outreach-crm` agent
- Agent reads opportunity record including contact notes and suggested tracks
- Produces a personalised outreach message
- Presents draft for artist review
- **[APPROVAL GATE]** Advances to `status: draft-ready` on approval

---

#### `ovl outreach send --opportunity <opportunity-id>`

Send an approved outreach message.

```text
ovl outreach send --opportunity <opportunity-id>
```

**Behaviour:**

- Confirms `opportunity.status` is `approved`
- **[CONFIRMATION GATE]** Shows recipient, subject, and message. Requires explicit confirmation. Cannot be bypassed with `--yes`
- Sends via Gmail MCP if configured, or outputs the message for manual sending
- Updates `status: sent`, logs `action: sent` in `outreach_history`, sets `follow_up_due`

**Errors:**

- `NOT_APPROVED` — draft has not been approved; run `ovl outreach draft` first
- `GMAIL_MCP_NOT_CONFIGURED` — email MCP not connected; message will be output for manual sending

---

#### `ovl outreach follow-up --opportunity <opportunity-id>`

Draft a follow-up for an opportunity with no response.

```text
ovl outreach follow-up --opportunity <opportunity-id>
```

**Behaviour:**

- Invokes the `outreach-crm` agent to draft a brief follow-up
- **[APPROVAL GATE]** Same approval and send flow as `ovl outreach draft` + `ovl outreach send`

---

#### `ovl outreach log-response --opportunity <opportunity-id>`

Record a response received from an outreach contact.

```text
ovl outreach log-response --opportunity <opportunity-id>
```

**Behaviour:**

- Prompts for the response content (paste or describe)
- Invokes the `outreach-crm` agent to interpret and advise on next steps
- Logs `action: response-received` in `outreach_history`
- Updates `status: responded`

---

#### `ovl outreach score --opportunity <opportunity-id>`

Score an opportunity for fit with the artist.

```text
ovl outreach score --opportunity <opportunity-id>
```

**Behaviour:** Invokes the `outreach-crm` agent to review the opportunity details and propose a match score (1–10) with rationale. Artist confirms or adjusts. Writes `opportunity.match{}`.

---

#### `ovl outreach intake --type <type>`

Record a new inbound inquiry as an opportunity.

```text
ovl outreach intake --type <type>
```

**Behaviour:**

- Prompts for contact details, inquiry description, source, and any deadline
- Creates an `opportunity.json` record with `status: identified`
- Logs `action: identified` in `outreach_history`

---

#### `ovl outreach close --opportunity <opportunity-id> --outcome <outcome>`

Record the final outcome of an opportunity.

```text
ovl outreach close --opportunity <opportunity-id> --outcome won|lost|declined
```

**Behaviour:**

- Sets `opportunity.status` to the specified outcome
- For `won`: prompts for confirmed value and which tracks were used
- For `lost` / `declined`: prompts for reason if known
- Logs the outcome action in `outreach_history`

---

#### `ovl outreach log --opportunity <opportunity-id> --action <action>`

Manually log an action on an opportunity.

```text
ovl outreach log --opportunity <opportunity-id> --action <action> [--note "<note>"]
```

**Behaviour:** Appends an entry to `opportunity.outreach_history[]`. For recording manual actions (platform-based submissions, phone calls, in-person conversations).

**Valid actions:** Any value from the `outreach_history.action` enum in `opportunity.schema.json`.

---

### Finance

#### `ovl finance add-revenue`

Log a revenue entry.

```text
ovl finance add-revenue \
  --source <source> \
  --amount <amount> \
  --currency <EUR|USD|...> \
  --period <YYYY-MM> \
  [--artist <artist-id>] \
  [--release <release-id>] \
  [--opportunity <opportunity-id>] \
  [--description "<text>"]
```

**Behaviour:** Creates a `finance-entry.json` record with `type: revenue` and appends it to `workspace/finance/revenue.json`.

---

#### `ovl finance add-expense`

Log an expense entry.

```text
ovl finance add-expense \
  --source <category> \
  --amount <amount> \
  --currency <EUR|USD|...> \
  --date <YYYY-MM-DD> \
  [--artist <artist-id>] \
  [--description "<text>"]
```

**Behaviour:** Creates a `finance-entry.json` record with `type: expense` and appends it to `workspace/finance/expenses.json`.

---

#### `ovl finance summary --period <YYYY-MM>`

Generate a financial summary for a period.

```text
ovl finance summary --period <YYYY-MM> [--brief]
```

**Behaviour:**

- Invokes the `finance-manager` agent
- Agent reads all revenue and expense entries for the period
- Produces a summary: total revenue by source, total expenses by category, net position, goal progress, 3-month and 12-month trends
- `--brief` produces a one-paragraph summary instead of the full report

---

#### `ovl finance quote --opportunity <opportunity-id>`

Generate a pricing quote for a commission opportunity.

```text
ovl finance quote --opportunity <opportunity-id>
```

**Behaviour:** Invokes the `finance-manager` agent to review the commission scope and propose a rate. Artist confirms before the quote is used.

---

### Metrics

#### `ovl metrics snapshot --period <YYYY-MM>`

Compile a metrics snapshot for a period.

```text
ovl metrics snapshot --period <YYYY-MM> [--artist <artist-id>] [--brief]
```

**Behaviour:**

- Invokes the `metrics-analyst` agent
- Agent reads platform export files from `workspace/metrics/[YYYY-MM]/raw/`
- Populates `workspace/metrics/[YYYY-MM]/[artist-id].json`
- Produces a written analysis with trends, top tracks, and anomalies
- `--brief` produces a one-paragraph summary

---

### Content

#### `ovl content brief --release <release-id>`

Generate a content campaign brief for a release.

```text
ovl content brief --release <release-id>
```

**Behaviour:**

- Invokes the `content-strategist` agent
- Agent reads release record, artist profile, and recent metrics snapshot
- Produces a campaign brief: announcement timing, platform-specific plan, track spotlight schedule
- **[APPROVAL GATE]** Artist reviews and approves the brief before social copy is generated

---

#### `ovl social draft --release <release-id>`

Generate social media copy for a release campaign.

```text
ovl social draft --release <release-id> [--platform instagram|youtube|facebook]
```

**Behaviour:**

- Invokes the `social-media` agent
- Agent reads the approved content brief and release record
- Produces platform-specific copy for artist review and posting
- If `--platform` is omitted, generates copy for all configured platforms

**Requires:** An approved content brief (`ovl content brief` run and approved first).

---

### Agents

#### `ovl agents list`

List all installed agent skills and their interaction patterns.

```text
ovl agents list
```

**Output:** Table of agent names, descriptions (from SKILL.md frontmatter), interaction patterns, and connection status (whether required MCPs are configured).

---

### MCP

#### `ovl mcp list`

List available MCPs and their connection status.

```text
ovl mcp list
```

**Output:** Table of MCP names, descriptions, connection status, and the commands they enable.

---

#### `ovl mcp connect <mcp-name>`

Connect an MCP integration.

```text
ovl mcp connect <mcp-name>
```

**Behaviour:** Launches the authentication flow for the specified MCP. Credentials are stored locally and never written to workspace JSON records or the OVL repository.

**Available MCPs:** `gmail`, `google-calendar`, `internet-archive`, `amuse` (read-only where API permits)

---

#### `ovl mcp disconnect <mcp-name>`

Disconnect an MCP integration.

```text
ovl mcp disconnect <mcp-name>
```

---

### Commission (Shortcut)

#### `ovl commission agreement --opportunity <opportunity-id>`

Generate a commission agreement from the workspace template.

```text
ovl commission agreement --opportunity <opportunity-id>
```

**Behaviour:**

- Reads the opportunity record for scope, rights, timeline, and payment terms
- Populates `workspace/label/templates/commission-agreement.md` with the agreed details
- Outputs the filled agreement for artist review
- **[APPROVAL GATE]** Artist approves before it is sent to the client

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

---

## Environment Variables

| Variable | Description |
|---|---|
| `OVL_WORKSPACE` | Override the workspace path. Equivalent to `--workspace <path>` on every command. |
| `OVL_ARTIST` | Override the default artist. Equivalent to `--artist <id>` on every command. |
| `OVL_YES` | If set to `1`, skip non-critical confirmation prompts. Critical gates (release submit, archive push) are never bypassed. |
| `OVL_JSON` | If set to `1`, output all results as JSON. |

---

## Command Groups Summary

```text
Workspace:   init · status · validate · state show · state sync
Artist:      artist create · artist list · artist show · artist add-alias
Release:     release create · release list · release show · release advance
             release set-profile · release set-live · release add-link · release submit
Track:       track add · track show · track set-file · isrc assign
Mastering:   mastering start · mastering profile create · mastering profile list
QC:          qc check
Archive:     archive push · archive status
Outreach:    outreach research · outreach review · outreach draft · outreach send
             outreach follow-up · outreach log-response · outreach score
             outreach intake · outreach close · outreach log
Finance:     finance add-revenue · finance add-expense · finance summary · finance quote
Metrics:     metrics snapshot
Content:     content brief · social draft
Agents:      agents list
MCP:         mcp list · mcp connect · mcp disconnect
Commission:  commission agreement
```

---

## Implementation Notes

**Argument parsing:** Subcommand-based, with global flags inherited by all commands. Interactive prompts are used for required fields not supplied as flags. Long-running operations (archive uploads, agent sessions) stream output rather than blocking.

**Schema validation:** JSON Schema Draft 7. Every workspace write validates against the appropriate schema before the file is touched. Validation errors report the field path and the constraint violated.

**Frontmatter parsing:** Agent skill files (`SKILL.md`) carry YAML frontmatter that the CLI reads to build the agent registry (name, description, interaction pattern, required integrations).

**Agent invocation:** When a command hands off to an agent, the CLI loads the agent's `SKILL.md` and relevant workspace records, then opens a session in the current terminal (or via a configured API if running non-interactively). The specific integration mechanism depends on the deployment context and is specified separately in `cli/AGENT-INTEGRATION.md`.

**Config file:** The CLI reads from `.ovlrc` in the workspace root for persistent local configuration (preferred artist, output format, integration preferences). This file is `.gitignored` alongside `workspace/`.
