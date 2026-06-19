# Monthly Review

A structured check-in combining analytics, finance, and strategic planning. Designed to be completed in under 90 minutes once the data is gathered. The output is an updated `label-state.md` with refreshed goal progress and a clear set of priorities for the coming month.

Run on a consistent cadence — the first Sunday of the month works well, as most platform analytics close on the first of the month.

---

## Prerequisites

- Platform analytics exports available (Spotify for Artists, YouTube Studio, Amuse.io royalty export, Bandcamp stats)
- Previous month's `metrics-snapshot.json` exists (or this is the first review)
- At least one `artist.json` in the workspace

---

## Part 1: Metrics Snapshot (30 minutes)

### Gather data

Before invoking the Metrics Analyst, collect raw data from each platform:

| Platform | Where to export | What to get |
|---|---|---|
| Spotify for Artists | spotify.com/artist → Home → Export | Monthly listeners, streams, saves, followers |
| YouTube Studio | studio.youtube.com → Analytics → Export | Views, watch time, subscribers |
| Amuse.io | Distribution → Royalty reports → Export CSV | Streams and revenue by platform |
| Bandcamp | bandcamp.com/fan_reports | Sales count, revenue, downloads, fans |
| SoundCloud (if active) | soundcloud.com/dashboard → Stats | Plays, followers |

Save exports to `workspace/metrics/[YYYY-MM]/raw/`.

### Invoke the Metrics Analyst

```bash
ovl metrics snapshot --period <YYYY-MM>
```

→ `metrics-analyst`

The Metrics Analyst reads the raw exports, populates a new `metrics-snapshot.json` for the period, and produces a written summary covering:

- Month-over-month changes for each key metric
- Top performing tracks and what drove their performance
- Platform trends (which platforms are growing, flat, or declining)
- Outreach outcomes if any opportunities closed this month
- Any anomalies (sudden spikes, unexpected drops) with possible explanations

**[APPROVAL GATE]** The artist reviews the snapshot and summary before it is written. Any factual corrections are made.

✓ Output: `workspace/metrics/[YYYY-MM]/[artist-id].json` created, `analyst_notes` populated

---

## Part 2: Finance Summary (15 minutes)

```bash
ovl finance summary --period <YYYY-MM>
```

→ `finance-manager`

The Finance Manager reads the period's revenue entries (populated from the Amuse.io export and Bandcamp report in Part 1), plus any manually logged revenue or expenses, and produces:

- Total revenue for the period by source
- Total expenses for the period by category
- Net position (revenue minus expenses)
- Progress toward the monthly revenue goal
- Rolling 3-month and 12-month trend

Any revenue entries not yet in the workspace are added now:

```bash
ovl finance add-revenue \
  --source <platform> \
  --amount <amount> \
  --currency EUR \
  --period <YYYY-MM> \
  --description "<description>"
```

Any expenses incurred during the month that are not yet logged:

```bash
ovl finance add-expense \
  --source <category> \
  --amount <amount> \
  --currency EUR \
  --date <YYYY-MM-DD> \
  --description "<description>"
```

**[APPROVAL GATE]** Finance summary reviewed and confirmed before `label-state.md` is updated with goal progress.

✓ Output: `label-state.md → ## Goal Progress` updated with current revenue figure and trend

---

## Part 3: Strategic Check-in (30 minutes)

With metrics and finance in front of both the artist and the orchestrator, the strategic check-in addresses four questions.

### What worked?

The orchestrator surfaces the top three positive signals from the metrics and finance summaries:

- Strongest metric improvement
- Best performing track or platform
- Any outreach wins or notable listener engagement

### What needs attention?

The orchestrator surfaces anything that warrants a response:

- Any metric declining for more than two consecutive months
- Revenue below target for three or more months
- Outreach pipeline stalling (no opportunities in `sent` or `follow-up` stage)
- Any open loops older than 30 days

### What is changing next month?

Based on the release pipeline status and the content calendar:

- Are any releases moving from mastering to QC or from ready to submitted?
- Are there campaign activities due?
- Any scheduled performances or collaboration sessions?

### What are the priorities?

The orchestrator proposes three priorities for the coming month. The artist refines and confirms them. These are written to `label-state.md`:

```text
## Goal Progress
...
Monthly priorities (July 2025):
1. Complete mastering tracks 5–8 of Spectra
2. Run outreach research for Spectra placement opportunities
3. Log expenses from studio session on June 20
```

---

## Part 4: Close and Update State

The orchestrator writes the session summary to `label-state.md`:

- Session Log entry: what was reviewed, what changed, priorities set
- Active Projects: updated release statuses
- Open Loops: new items added, resolved items removed
- Goal Progress: revenue figure, trend, monthly priorities

**[APPROVAL GATE]** Artist reviews the full state update before it is written.

```bash
ovl state sync
```

→ `label-state.md` updated

✓ **Review complete when:** Metrics snapshot written, finance summary confirmed, three monthly priorities set, state document updated

---

## Abbreviated Review (Under 30 Minutes)

When time is constrained, a shorter version focuses only on the most critical outputs:

```bash
ovl metrics snapshot --period <YYYY-MM> --brief
ovl finance summary --period <YYYY-MM> --brief
```

The `--brief` flag produces a one-paragraph summary of each rather than the full report. The artist confirms priorities verbally with the orchestrator, which writes a minimal state update. Full analysis can be run the following week.

---

## First Review Setup

On the first monthly review, there is no prior snapshot to compare against. The Metrics Analyst notes this and flags all figures as baselines rather than trends. The Finance Manager initialises the revenue and expense files if they do not yet exist. Goal progress section of `label-state.md` is populated for the first time with the artist's stated targets.
