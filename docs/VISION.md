# VISION.md — Open Virtual Label

> A framework for independent musicians and small collectives to operate with the support of a coordinated team of AI agents — without ceding creative control, revenue, or identity to a traditional label.

---

## What This Is

Open Virtual Label (OVL) is an open-source system of agents, workflows, schemas, and tooling that performs the operational work a record label or management company would typically provide. The artist stays focused on their craft. The agents handle the rest — with the artist's approval at every meaningful gate.

The system is built around three convictions:

**Agents should earn trust, not assume it.** Every consequential action — sending an email, submitting a release, publishing content — requires explicit human approval before it happens. Agents research, draft, and prepare; humans decide.

**Structure should travel.** The workflows, schemas, and agent definitions are generic and reusable. An artist adopting OVL brings their own workspace — their catalog, their contacts, their style — but the underlying system belongs to no one's specific situation.

**Transparency compounds over time.** By keeping decisions, drafts, metrics, and process notes in structured, version-controlled files, the system builds institutional memory. Releases get better. Outreach gets sharper. The mastering chain improves. A solo artist can accumulate the kind of operational knowledge that usually lives inside a label's A&R team.

---

## What This Is Not

- A label that takes a cut of revenue or owns any rights
- A replacement for human creative judgment
- An autonomous system that acts without artist approval
- A platform or SaaS product
- Specific to any artist, genre, or distributor

OVL is tooling. It has no contractual relationship with the artists who use it.

---

## The Workspace Model

The project repository contains only the reusable system — agent definitions, schemas, workflow playbooks, and CLI tooling. No artist-specific content lives in the repo.

Each label instance maintains a `workspace/` directory that is `.gitignored` by default. This is where the specific, private, operational content lives:

```text
workspace/
├── label/
│   ├── profile.json          # Label identity, contact, license defaults
│   ├── styleguide/           # Colors, fonts, tone of voice, logo assets
│   └── websites/
│       └── label-site/       # Label's own web presence
├── artists/
│   └── [artist-id]/
│       ├── artist.json       # Artist profile, platform IDs, PRO registration
│       ├── releases/         # One subdirectory per release
│       ├── catalog/          # Master track registry
│       └── websites/
│           └── artist-site/  # Artist-specific web presence
├── outreach/
│   ├── opportunities.json    # CRM: commission, sync, playlist opportunities
│   └── contacts/             # Individual contact records
├── finance/
│   ├── revenue.json          # Income records by source and period
│   └── expenses.json         # Tools, services, production costs
├── metrics/
│   └── [YYYY-MM]/            # Monthly analytics snapshots
└── state/
    └── label-state.md        # Orchestrator's living session context document
```

A `workspace/` scaffold — with placeholder files and inline comments explaining each field — ships with the repo so new label operators have a clear starting point, not a blank page.

---

## The Agent Model

OVL is organized around a central orchestrator that delegates to a roster of specialist agents. Each agent is defined as a markdown skill file containing its role, responsibilities, interaction patterns, and what it reads and writes.

### Orchestrator

The orchestrator is the entry point for any session. It reads `workspace/state/label-state.md` at the start of each conversation to restore context, then routes requests to the appropriate specialist. At the end of a session it summarizes any decisions, pending approvals, and open loops back into the state document.

The orchestrator does not specialize. It coordinates.

### Specialist Agents

| Agent                       | Core Responsibility                                                                                                                                                                                                                                                         |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Mastering Companion**     | Conversational expert that walks an artist through the mastering chain for a specific target profile (LUFS, true peak, LRA, platform norms). Reads from a `mastering_profile` record; writes measurements and notes back to the track record on completion.                 |
| **Archive Agent**           | Packages and uploads completed releases — lossless masters, stems, project files — to Internet Archive and secondary object storage. Updates the release record with confirmed archive paths and checksums.                                                                 |
| **Outreach / CRM Agent**    | Researches commission, sync, playlist, and collaboration opportunities. Scores each against the artist profile. Drafts outreach messages for human review. Tracks status through the full pipeline from identified to won or closed. Never sends without explicit approval. |
| **QC / Pre-Release Agent**  | Runs a completeness and consistency check on a release before distribution submission. Verifies every track has ISRC, mastering data, and a passed QC flag. Checks artwork specs, metadata consistency, and license fields. Blocks or flags any gaps.                       |
| **Metrics Analyst**         | Compiles periodic analytics from platform exports. Identifies trends, surfaces anomalies, and produces structured reports that feed into strategic reviews.                                                                                                                 |
| **Finance Manager**         | Tracks revenue by source and period. Monitors expenses. Reports progress toward income goals.                                                                                                                                                                               |
| **Content Strategist**      | Plans release timelines, content calendars, and cross-platform campaigns. Produces schedules and briefs that other agents and the artist can execute against.                                                                                                               |
| **Social Media Specialist** | Generates platform-specific copy — captions, descriptions, announcements — for human review and posting. Does not post autonomously.                                                                                                                                        |
| **Release Manager CLI**     | Command-line interface (`ovl`) for common operations: creating releases, adding tracks, assigning ISRCs, triggering agent sessions, running QC, pushing to archive, and submitting to distribution.                                                                         |

New agents can be added by creating a new skill file in `/agents/`. The orchestrator can route to any named agent it finds in that directory.

---

## The Data Layer

Every entity in OVL has a canonical JSON schema. Agents read from and write to these records. The schemas are versioned so the format can evolve without breaking existing data.

### Core Schemas

**`artist.json`** — Identity, platform IDs, PRO registration, default license, contact details, distribution account reference.

**`release.json`** — Title, type, status, target dates, track list, distributor, license, archive paths, store links. The release record is the spine that other records attach to.

**`track.json`** — Title, ISRC, ISWC, duration, collaborators and split percentages, file paths (master WAV, stems, project file, lossy exports), mastering measurements, QC status.

**`mastering_profile.json`** — Target LUFS range, true peak ceiling, LRA guidance, sample rate and bit depth, platform-specific notes, and a checklist of steps. Profiles are named and reusable across tracks. After each use, the Mastering Companion can append session notes for continuous improvement.

**`opportunity.json`** — CRM record for outreach targets. Type (commission, sync, playlist, collaboration, performance), status through the pipeline, contact details, match score and rationale, suggested tracks, full outreach history with draft approval records, and estimated value.

**`label-state.md`** — A human-readable markdown document maintained by the orchestrator. Contains current project statuses, recent decisions, pending approvals, open loops, and goal progress. This is the orchestrator's memory between sessions.

Schema files live in `/schemas/` and are validated using standard JSON Schema tooling. The workspace scaffold includes starter data files with comments explaining each field.

---

## Human-in-the-Loop Patterns

Three distinct interaction patterns govern how agents and humans share work.

### 1. Approval Gate

The agent prepares everything; the human decides whether it happens. No action crosses this gate without an explicit yes.

Used for: outreach emails before sending, distribution submissions after QC, commission quotes before delivery, any communication sent on the artist's behalf.

### 2. Guided Session

The human operates tools while the agent provides real-time expert guidance. The agent reads instrument output, interprets it against the target profile, and advises on what to adjust next. The human makes every move in their DAW or tool chain.

Used for: mastering sessions (Mastering Companion reads your LUFS measurement and tells you what to do next), onboarding new releases, walking through complicated processes for the first time.

### 3. Review and Refine

The agent produces a complete draft. The human reviews, edits if needed, and approves before the output is used.

Used for: social media copy, press kit text, outreach message drafts, release descriptions, content calendar proposals.

Every interactive process in OVL falls into one of these three patterns. When adding a new agent or workflow, the pattern should be declared explicitly in the agent's skill file.

---

## The Release Pipeline

A release in OVL moves through a defined set of stages, each with a responsible agent and a human checkpoint.

```text
Production (artist)
    ↓
Track Registration → release.json + track.json records created
    ↓
Mastering Session → Mastering Companion (guided)
    ↓
QC Check → QC Agent (blocking — flags must resolve before proceeding)
    ↓
Archive → Archive Agent uploads masters, stems, project files
    ↓
[HUMAN APPROVAL] Distribution submission package reviewed
    ↓
Submission → distributor
    ↓
Store Links → release.json updated with live URLs
    ↓
Content Campaign → Content Strategist briefs Social Media Specialist
    ↓
Outreach → CRM Agent surfaces placement opportunities
```

The CLI (`ovl`) can advance a release through this pipeline one step at a time, or an artist can trigger individual stages in any order once records exist.

---

## The Outreach Loop

Outreach is the process most likely to fall through the cracks without structure. OVL treats it as a continuous loop with named stages and a human gate before any message is sent.

```text
Research → Score → [HUMAN: review opportunity list]
    ↓
Draft → [HUMAN: review and approve draft]
    ↓
Send → logged in opportunity.json
    ↓
Follow-up reminder → [HUMAN: approve follow-up]
    ↓
Outcome recorded → won / lost / declined / stale
```

The CRM Agent maintains a pipeline view of all open opportunities and surfaces follow-ups on a configurable schedule. It never infers approval from silence.

---

## The CLI — `ovl`

The `ovl` command-line tool wraps the most common operations so they can be run without manually editing JSON files. It validates input against the schemas and invokes the appropriate agent when a process requires conversational guidance.

Representative commands:

```bash
ovl init                          # scaffold a new workspace/
ovl artist create                 # interactive artist profile setup
ovl release create "Spectra"      # new release record
ovl track add "Chromatic Drift"   # add track to a release
ovl isrc assign                   # assign ISRC to unregistered tracks
ovl mastering start               # launch Mastering Companion session
ovl qc check --release spectra    # run pre-release QC
ovl archive push --release spectra
ovl release submit --distributor amuse
ovl outreach research             # trigger CRM research cycle
ovl state show                    # display current label-state.md
```

Full command reference lives in `/cli/README.md`.

---

## Repository Structure

```text
/                         # System-level files (README, LICENSE, CONTRIBUTING)
├── VISION.md             # This document
├── ARCHITECTURE.md       # Component inventory and cross-references
├── agents/               # Agent skill files
│   ├── orchestrator/
│   ├── mastering-companion/
│   ├── archive/
│   ├── outreach-crm/
│   ├── qc/
│   ├── metrics-analyst/
│   ├── finance-manager/
│   ├── content-strategist/
│   └── social-media/
├── schemas/              # JSON Schema definitions
│   ├── artist.schema.json
│   ├── release.schema.json
│   ├── track.schema.json
│   ├── mastering-profile.schema.json
│   └── opportunity.schema.json
├── workflows/            # Playbook documents
│   ├── release-pipeline.md
│   ├── mastering-session.md
│   ├── outreach-loop.md
│   └── onboarding.md
├── cli/                  # `ovl` CLI source and documentation
├── workspace-scaffold/   # Template workspace/ with placeholder files
└── docs/                 # Contributing guide, governance, FAQ
```

The `workspace/` directory used by any given label instance is not part of this repository. It is local, private, and `.gitignored`.

---

## Design Principles

**Co-locate documentation with code.** Every component has a `README.md` in its own directory. There is no separate wiki. Documentation that drifts from the implementation becomes fiction.

**Schemas are the contract.** Agents read and write JSON records. If two agents need to share data, that data must be in a schema. Undocumented conventions become bugs.

**Human approval is not optional.** It is a named step in every pipeline that involves external communication, publication, or financial action. Removing it requires a deliberate change to the workflow definition, not just skipping a prompt.

**Continuous improvement is built in.** The Mastering Companion appends session notes to mastering profiles. The CRM Agent records win/loss outcomes and rationale. The metrics snapshots accumulate into trend data. The system should get measurably better at its job the longer it runs.

**The workspace is sovereign.** Artists own their data. The workspace directory structure is a recommendation, not a requirement enforced by tooling. The schemas are strict about field names and types; the directory layout can be adapted.

---

## Governance

OVL is released under an open-source license (Apache 2.0 for tooling and schemas is the current leaning). Contributions are welcome via pull request.

The project has no affiliation with any distributor, PRO, streaming platform, or commercial label. It takes no revenue share. It creates no contractual relationship between contributors and artists who use the system.

Artist workspaces are private. No workspace data is collected, transmitted, or shared with the project.

---

*Last updated: June 2026*
*Status: Pre-implementation draft — structure and terminology subject to revision*
