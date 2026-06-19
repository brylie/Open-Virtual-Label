# CLI Specification — `ovl release`

Commands for managing release records
(`workspace/artists/[artist-id]/releases/[release-id]/release.json`).

All writes validate against `schemas/release.schema.json` before saving.

---

## Command Reference

---

### `ovl release create "<title>"`

Create a new release record.

```text
ovl release create "<title>" [--artist <artist-id>] [--type album|ep|single|compilation]
```

**Behaviour:**

- Generates a slug from the title (e.g. `spectra`)
- Creates `workspace/artists/[artist-id]/releases/[slug]/release.json` with
  `status: in-production`
- Creates the `tracks/` subdirectory
- If `--artist` is omitted and only one artist exists in the workspace, uses that
  artist. If multiple artists exist, prompts.
- Prompts for target release date and distributor submission deadline

**Options:**

| Option | Default | Description |
|---|---|---|
| `--artist <id>` | (prompted if ambiguous) | Artist this release belongs to |
| `--type <type>` | `album` | Release type: `album`, `ep`, `single`, `compilation` |

**Output:** Path to created `release.json` and the generated release ID.

---

### `ovl release list`

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

### `ovl release show <release-id>`

Display a release record and its track list with statuses.

```text
ovl release show <release-id>
```

**Output:** Formatted release details plus a table of tracks showing mastering and QC
status.

---

### `ovl release advance <release-id> --status <status>`

Manually advance a release to the next pipeline status.

```text
ovl release advance <release-id> --status <status>
```

**Behaviour:**

- Validates that the advance is legitimate (e.g. cannot advance to `ready` if
  `qc.passed` is not `true`)
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

### `ovl release set-profile <release-id> --profile <profile-id>`

Assign a mastering profile to a release.

```text
ovl release set-profile <release-id> --profile <profile-id>
```

**Behaviour:** Sets `release.mastering_profile_id`. The profile must exist in
`workspace/artists/[artist-id]/mastering-profiles/`.

---

### `ovl release set-live <release-id> --date <YYYY-MM-DD>`

Mark a submitted release as live.

```text
ovl release set-live <release-id> --date <YYYY-MM-DD>
```

**Behaviour:** Sets `release.status: live` and `release.dates.released`. Prompts for
confirmation.

---

### `ovl release add-link <release-id> --platform <platform> --url <url>`

Add a store link to a live release.

```text
ovl release add-link <release-id> --platform <platform> --url <url>
```

**Behaviour:** Sets the specified field in `release.store_links`.

**Valid platforms:** `spotify`, `apple_music`, `youtube_music`, `bandcamp`, `soundcloud`,
`tidal`, `amazon_music`, `fma` (Free Music Archive), `subvert_fm` (Subvert.fm co-op platform).

---

### `ovl release submit <release-id> --distributor <distributor>`

Prepare and submit a distribution package.

```text
ovl release submit <release-id> --distributor <distributor>
```

**Behaviour:**

- Assembles a submission summary: title, type, target date, track list with ISRCs and
  durations, artwork details, license, contributor credits
- **[CONFIRMATION GATE]** Presents the full summary for artist review. This gate cannot
  be bypassed with `--yes`
- On confirmation, either generates the distributor's required submission format for
  manual upload, or calls the distributor's API if an MCP is configured
- Sets `release.status: submitted` and `release.dates.submitted`

**Errors:**

- `QC_NOT_PASSED` — release has not passed QC; run `ovl qc check` first
- `ARCHIVE_NOT_COMPLETE` — masters not archived; run `ovl archive push` first
- `MISSING_ISRC` — one or more tracks are missing ISRC; assign before submitting
