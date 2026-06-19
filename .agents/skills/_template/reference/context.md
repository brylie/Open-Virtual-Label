# OVL Context Reference

This file is bundled with every OVL agent skill. It provides the shared
context any agent needs to operate within the Open Virtual Label system.
Load it when you need detail on workspace structure, schemas, interaction
patterns, or the label state document format.

---

## Workspace Structure

The `workspace/` directory is local and `.gitignored`. It never appears in
the OVL repository. Its layout:

```
workspace/
├── label/
│   ├── profile.json          # Label identity, license defaults, distributor
│   └── styleguide/           # Colors, fonts, logo assets — operator-defined
├── artists/
│   └── [artist-id]/
│       ├── artist.json       # Artist profile, platform IDs, PRO registration
│       ├── releases/
│       │   └── [release-id]/
│       │       ├── release.json
│       │       └── tracks/
│       │           └── [track-id].json
│       └── websites/
├── outreach/
│   ├── opportunities.json    # CRM: all opportunities in pipeline
│   └── contacts/
├── finance/
│   ├── revenue.json
│   └── expenses.json
├── metrics/
│   └── [YYYY-MM]/            # Monthly analytics snapshots
└── state/
    └── label-state.md        # Coordinator's session memory
```

---

## Core Schemas

Schemas live in `schemas/` in the OVL repo. Workspace JSON files validate
against them. Key fields by record type:

### `artist.json`
- `id` — slug, matches directory name
- `display_name` — performing name
- `legal_name`
- `pro` — performing rights org (e.g. Teosto, ASCAP, PRS)
- `distributor` — e.g. `amuse`, `distrokid`, `cdbaby`
- `default_license` — e.g. `CC BY 4.0`
- `platforms` — object with `spotify_artist_id`, `youtube_channel_id`, etc.

### `release.json`
- `id` — slug
- `title`
- `artist_id` — references artist directory
- `release_type` — `album | ep | single`
- `status` — `in-production | mastering | qc | ready | submitted | live`
- `target_release_date`
- `license`
- `tracks` — array of track IDs
- `mastering_profile` — references a profile ID
- `archive` — object with `internet_archive_id`, `object_storage_path`, flags
- `store_links` — object with platform URLs once live

### `track.json`
- `id`
- `release_id`
- `title`
- `isrc`
- `position` — track number
- `collaborators` — array of `{ artist_id, role, split_percentage }`
- `files` — object with `master_wav`, `stem_mix`, `project_file` paths
- `mastering` — object with `integrated_lufs`, `true_peak_dbtp`, `lra`, `profile_used`, `mastered_date`, `notes`
- `qc_passed` — boolean
- `qc_checklist_id`

### `opportunity.json` (item within `opportunities.json` array)
- `id`
- `type` — `commission | sync-license | playlist-pitch | collaboration | performance`
- `status` — `identified | researched | draft-ready | sent | follow-up | won | lost | declined`
- `contact` — object with `name`, `role`, `email`, `url`, `notes`
- `match_score` — 1–10
- `tracks_suggested` — array of track IDs
- `outreach_history` — array of `{ date, action, approved_by, response }`

### `mastering_profile.json`
- `id`
- `name`
- `targets` — object with `integrated_lufs` (min/max), `true_peak_dbtp`, `lra_min`, `sample_rate_hz`, `bit_depth`
- `platform_notes` — per-platform guidance
- `checklist` — ordered array of step strings
- `session_notes` — array of notes appended after each use

---

## Label State Document Format

`workspace/state/label-state.md` is the coordinator's persistent memory.
Agents read it at session start and the coordinator updates it at session end.

```markdown
# Label State

_Last updated: [ISO date] — Session [N]_

## Active Projects
<!-- One line per release: title — status — next action -->

## Open Loops
<!-- Unresolved items: started but not complete, decisions pending,
     approvals waiting, follow-ups due with dates -->

## Pending Approvals
<!-- Items blocked on artist sign-off. Format:
     - [Agent]: [what needs approval] — since [date] -->

## Recent Decisions
<!-- Significant choices in the last 30 days with brief rationale -->

## Goal Progress
<!-- Current status against tracked goals: revenue, releases, outreach -->

## Session Log
<!-- Reverse chronological. Format per entry:
     ### [ISO date] — [one-line summary]
     - What was done
     - What changed in workspace records
     - Open loops created or resolved
     - Next recommended action -->
```

---

## Interaction Patterns

Every OVL agent uses one of three patterns. When in doubt about which
applies to a given situation, default to requiring explicit confirmation
before any write or external action.

### approval-gate
Agent researches, prepares, and drafts. Artist reviews and explicitly
approves. No external action (send, upload, submit, post) occurs without
a named confirmation step. If the artist does not respond or is ambiguous,
the agent waits — it does not infer approval from silence.

### guided-session
Artist operates their own tools (DAW, browser, terminal). Agent reads
output the artist provides, interprets it against reference data, and
advises on what to do next. The agent never claims to be running the
tools itself. Session ends when the artist confirms the work is complete.

### review-and-refine
Agent produces a complete draft (copy, report, plan, email). Artist
reviews, requests edits if needed, and approves before the output is
used or filed. Multiple revision rounds are normal.

---

## Boundaries All Agents Share

Regardless of role, every OVL agent:

- Does not send emails, post content, or submit to external services
  without an explicit approval gate
- Does not modify `release.json`, `track.json`, or `opportunities.json`
  without confirming changes with the artist first
- Does not claim certainty about platform policies, royalty rates, or
  legal matters — flags these as "verify before acting"
- Hands off to the coordinator at session end with a summary of what
  was done and any new open loops
