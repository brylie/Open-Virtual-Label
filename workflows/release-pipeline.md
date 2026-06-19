# Release Pipeline

The complete sequence from finishing production to a release going live on platforms, with archival and campaign launch. This is the master workflow that references the more detailed sub-workflows for mastering and outreach.

---

## Prerequisites

- At least one `artist.json` exists in the workspace
- All tracks are recorded and exported from the DAW as lossless WAV files
- Artwork exists at minimum 3000×3000px

---

## Stages at a Glance

```
1. Register          Create release and track records
2. Mastering         Master each track to target profile
3. QC                Verify release completeness
4. Archive           Upload masters, stems, project files
5. Submit            Send distribution package
6. Campaign          Content and outreach for the release
```

Release `status` field advances through:
`in-production → mastering → qc → ready → submitted → live`

---

## Stage 1: Register

Create the release and track records so the rest of the pipeline has something to write to.

```bash
ovl release create "<title>" --artist <artist-id>
```

This creates `workspace/artists/[artist-id]/releases/[release-id]/release.json` with `status: in-production`.

For each track:

```bash
ovl track add "<title>" --release <release-id> --position <n>
```

This creates `workspace/artists/[artist-id]/releases/[release-id]/tracks/[track-id].json`.

Assign ISRCs. If the distributor assigns them, leave `isrc: null` until they are issued and then fill them in. If the label registers them independently:

```bash
ovl isrc assign --release <release-id>
```

Set the mastering profile for the release:

```bash
ovl release set-profile --release <release-id> --profile <profile-id>
```

If no profile exists yet for this artist's genre, create one first:

```bash
ovl mastering profile create
```

→ `label-state.md` updated: release added to Active Projects

✓ **Stage complete when:** `release.json` exists with all tracks listed in order, mastering profile set, `status: in-production`

---

## Stage 2: Mastering

Master each track in sequence. Each mastering session is its own conversation with the Mastering Companion. See `mastering-session.md` for the full per-track protocol.

```bash
ovl mastering start --track <track-id>
```

→ `mastering-companion`

The Mastering Companion loads the track record and the release's mastering profile, then guides the artist through the session. On completion it writes to `track.mastering{}` and prompts updating `track.qc.passed` once measurements are within target.

Repeat for every track in the release. Track mastering sessions can happen across multiple days — the release status stays `mastering` until all tracks are complete.

When all tracks have mastering data populated:

```bash
ovl release advance --release <release-id> --status qc
```

✓ **Stage complete when:** Every `track.json` has `mastering.integrated_lufs`, `mastering.true_peak_dbtp`, and `mastering.mastered_date` populated; `release.status: qc`

---

## Stage 3: QC

Run the pre-release quality check. The QC agent inspects every track and the release record for completeness and consistency.

```bash
ovl qc check --release <release-id>
```

→ `qc`

The QC agent checks:

**Track-level (each track must pass):**

- `isrc` is not null
- `mastering.integrated_lufs` is within the profile's target range
- `mastering.true_peak_dbtp` is at or below the profile's ceiling
- `mastering.bit_depth` and `mastering.sample_rate_hz` are set
- `files.master_wav` path exists

**Release-level:**
- All tracks in `release.tracks[]` have a corresponding `track.json`
- Track positions are sequential with no gaps or duplicates
- `artwork.primary_file` path exists
- `artwork.dimensions_px` ≥ 3000
- `artwork.format` is `png`, `jpg`, or `tiff`
- `artwork.qc_passed` is true
- `release.license` is set
- `dates.target_release` is set
- `dates.distributor_submission_deadline` is set and not in the past

Any failures are listed in `release.qc.notes` and the relevant `track.qc.failures[]`. The QC agent surfaces these to the artist before marking anything as passed.

**[APPROVAL GATE]** The artist reviews the QC report. Failures must be resolved or explicitly overridden with a written reason (`track.qc.override`). The QC agent will not set `release.qc.passed: true` without confirmation.

On a clean pass:

```bash
ovl release advance --release <release-id> --status ready
```

→ `label-state.md` updated

✓ **Stage complete when:** `release.qc.passed: true`, all `track.qc.passed: true` or overrides recorded, `release.status: ready`

---

## Stage 4: Archive

Package and upload the release to long-term storage before submitting to distribution. Archival happens before submission — not after — so the canonical record exists independently of any platform.

```bash
ovl archive push --release <release-id>
```

→ `archive`

The Archive agent assembles a manifest of all files to upload:
- Lossless master WAV for each track (`files.master_wav`)
- Stems ZIP for each track that has one (`files.stems_zip`)
- DAW project file for each track (`files.project_file`)
- Artwork file (`artwork.primary_file`)
- `release.json` and all `track.json` records

**[APPROVAL GATE]** The agent presents the manifest to the artist for review before uploading. Upload cannot proceed without confirmation.

On approval, the agent uploads to:

1. Internet Archive (primary) — via the IA S3-compatible API
2. Object storage (secondary) — operator-configured bucket

After successful upload, it verifies checksums and writes to `release.archive{}`:
- `internet_archive_id`
- `internet_archive_url`
- `object_storage_path`
- `masters_archived: true`
- `stems_archived: true` (if stems were present)
- `project_files_archived: true` (if project files were present)
- `checksums_verified: true`
- `archive_date`

→ `label-state.md` updated

✓ **Stage complete when:** `release.archive.masters_archived: true`, `release.archive.checksums_verified: true`, both archive locations confirmed

---

## Stage 5: Submit

Prepare and submit the distribution package. This is the only stage with direct external consequences — a submission cannot be recalled once accepted by the distributor.

```bash
ovl release submit --release <release-id> --distributor <distributor>
```

The CLI assembles a submission summary:

- Release title, type, and target release date
- Track list with titles, ISRCs, and durations
- Artwork file and dimensions
- License and contributor credits
- Distributor-specific metadata requirements

**[APPROVAL GATE]** The full submission summary is presented to the artist. The artist confirms every field is correct. This gate cannot be bypassed.

On confirmation, the CLI either:

- Generates the submission package in the distributor's required format for manual upload, or
- Interfaces with the distributor's API if available

`release.status` advances to `submitted`, `release.dates.submitted` is set.

When the release goes live, update:

```bash
ovl release set-live --release <release-id> --date <YYYY-MM-DD>
```

Then populate store links as they become available:

```bash
ovl release add-link --release <release-id> --platform spotify --url <url>
```

→ `label-state.md` updated: release moved from Active Projects to "live"

✓ **Stage complete when:** `release.status: submitted`, submission date recorded, distributor confirmation received

---

## Stage 6: Campaign

Begin the release campaign once submission is confirmed. This stage runs in parallel with the distributor review window.

### Content campaign

```bash
ovl content brief --release <release-id>
```

→ `content-strategist`

The Content Strategist produces a campaign brief covering:
- Announcement timing relative to release date
- Platform-specific content plan (YouTube, Instagram, etc.)
- Track spotlight schedule
- Any behind-the-scenes or making-of content

**[APPROVAL GATE]** Artist reviews the brief before the Social Media Specialist generates copy.

```bash
ovl social draft --release <release-id>
```

→ `social-media`

Copy is delivered for artist review and posting. The Social Media Specialist does not post autonomously.

### Placement outreach

```bash
ovl outreach research --release <release-id>
```

→ `outreach-crm`

The CRM agent researches placement opportunities suited to the new release: sync licensing, playlist pitching, podcast placement. Creates `opportunity.json` records for each viable prospect.

See `outreach-loop.md` for the full outreach sequence from this point.

→ `label-state.md` updated with campaign status and any new opportunities

✓ **Stage complete when:** Campaign brief approved, social copy drafted, outreach research complete and in pipeline

---

## State Through the Pipeline

The orchestrator maintains `label-state.md` throughout. At each stage transition, the Active Projects entry for the release updates:

```
· Spectra — mastering (4 of 8 tracks) — next: master tracks 5–8
· Spectra — qc — next: run ovl qc check
· Spectra — ready — next: archive before submitting
· Spectra — submitted 2025-08-01, live 2025-09-01 — monitoring
```

---

## Common Interruptions

**A track fails QC.** Return to mastering for that track, then re-run QC. The release does not advance to `ready` until all failures are resolved.

**Distributor submission deadline missed.** Update `dates.distributor_submission_deadline` and `dates.target_release` accordingly. Notify any scheduled social content to adjust.

**Archive upload fails mid-way.** Re-run `ovl archive push` — the agent checks which files were already uploaded via checksum before re-uploading. Partial uploads do not advance `release.archive` flags.

**A collaborator's split is disputed before submission.** Update `track.collaborators[]` and re-run QC. Do not submit with an unresolved split disagreement.
