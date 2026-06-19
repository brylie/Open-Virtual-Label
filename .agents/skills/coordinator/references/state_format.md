# Label State Format

Loaded by the orchestrator when it needs the full specification for
`workspace/state/label-state.md`. The orchestrator maintains this document
as the persistent memory of the label across sessions.

---

## Document Structure

```markdown
# Label State

_Last updated: [ISO date] — Session [N]_

## Active Projects
## Open Loops
## Pending Approvals
## Recent Decisions
## Goal Progress
## Session Log
```

---

## Section Definitions

### `## Active Projects`

One line per release or major initiative. Format:

```text
· [Release title] — [status] — next: [one specific action]
· [Initiative] — [status] — next: [action]
```

Status values mirror `release.json` status field:
`in-production | mastering | qc | ready | submitted | live`

Example:

```text
· Spectra — mastering — next: master tracks 5–8
· New Beginnings — live — monitoring streams
```

### `## Open Loops`

Things that were started but not resolved. Each item should have enough
context that someone reading it cold knows what is needed and why it matters.

Format:

```text
· [What is unresolved] — [why it matters or what it blocks] — [due date if any]
```

Example:

```text
· QC not run on Spectra — blocks archive and submission — no deadline yet
· Follow-up to Calm Waters Podcast — sent 2025-06-01, 14 days elapsed — due by 2025-06-20
· ISRC missing on track 7 "Luminance" — blocks QC
```

Remove an item when it is resolved, and note the resolution in the Session Log.

### `## Pending Approvals`

Items that are drafted and waiting for explicit artist sign-off before an
agent can proceed. The responsible agent and item must be named.

Format:

```text
· [agent]: [what needs approval] — since [date]
```

Example:

```text
· outreach-crm: outreach email to Calm Waters Podcast — since 2025-06-10
· archive: release package manifest for Spectra — since 2025-06-14
```

### `## Recent Decisions`

Significant choices made in the last 30 days with brief rationale. Helps
avoid relitigating resolved questions.

Format:

```text
· [date] [Decision]: [rationale, one sentence]
```

Example:

```text
· 2025-06-01 Chose ambient-streaming-v1 mastering profile for Spectra:
  preserves dynamic range, appropriate for contemplative genre.
· 2025-05-20 Deferred playlist pitching until 6 tracks available:
  insufficient catalog for effective pitch at current count.
```

Prune entries older than 30 days unless they remain actively relevant.

### `## Goal Progress`

Current status against any tracked goals. Format is flexible — update
to match what the artist is actually tracking.

Example:

```text
Revenue goal: €100/month
· Current: ~€12/month (June 2025)
· Trend: +€3/month over last quarter

Release goal: Complete Spectra by September 2025
· Status: 4 of 8 tracks mastered, on track

Outreach goal: 3 placements per quarter
· Q2 2025: 1 confirmed (Soundscapes for Sleep podcast)
```

### `## Session Log`

Reverse chronological. Most recent entry first. Each entry records what
happened and what comes next. This is the primary tool for orienting at
the start of the following session.

Format:

```markdown
### [ISO date] — [one-line summary of session]

- [What was done — specific, e.g. "Mastered tracks 5 and 6 of Spectra"]
- [What changed in workspace records — e.g. "track.json updated for tracks 5–6"]
- [Open loops created — e.g. "Track 7 ISRC missing, added to open loops"]
- [Open loops resolved — e.g. "Resolved: outreach to Calm Waters Podcast sent"]
- Next: [specific recommended action for next session]
```

Keep individual log entries concise. Detailed session notes belong in the
workspace `metrics/` snapshots or release records, not here.

Retain the last 10 session entries. Archive older entries to
`workspace/state/session-log-archive.md` if the document grows unwieldy.

---

## Editing the State Document

The artist may edit `label-state.md` directly at any time. The orchestrator
treats it as the source of truth — any manual edits are respected. If the
orchestrator notices inconsistencies between the state document and workspace
records (e.g. a release marked "live" in state but `submitted` in
`release.json`), it flags the discrepancy and asks which is correct.

---

## Initial State Document

Created by `ovl init` or the orchestrator on first cold start:

```markdown
# Label State

_Last updated: [date] — Session 1_

## Active Projects

(none yet — add releases with `ovl release create`)

## Open Loops

· Complete workspace setup: fill in workspace/label/profile.json
· Create first artist profile: `ovl artist create`

## Pending Approvals

(none)

## Recent Decisions

(none)

## Goal Progress

(not yet configured — discuss goals with the orchestrator to set these up)

## Session Log

### [date] — Initial setup

- Label workspace scaffolded via `ovl init`
- Next: complete profile.json and create first artist profile
```
