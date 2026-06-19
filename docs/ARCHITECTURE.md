# ARCHITECTURE.md — Open Virtual Label

This document is the component inventory. Every subsystem in OVL is listed here with its responsibility, boundaries, inputs, outputs, and a link to its own directory for further detail. Read this after `docs/VISION.md` and before opening any subdirectory.

---

## System Map

```text
┌─────────────────────────────────────────────────────────┐
│                        Artist                           │
│            (creative input · approvals · review)        │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                  ovl CLI / TUI                          │
│         (primary interaction surface · Node.js)         │
└──────┬───────────────┬────────────────┬─────────────────┘
       │               │                │
       ▼               ▼                ▼
┌────────────┐  ┌─────────────┐  ┌───────────────┐
│  Workspace │  │   Agents    │  │  Schemas      │
│  (local,   │  │(markdown    │  │(JSON Schema   │
│  private)  │  │ skill files │  │ definitions)  │
└────────────┘  │ + LLM)      │  └───────────────┘
                └──────┬──────┘
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
   ┌─────────────┐ ┌────────┐ ┌──────────────┐
   │  Workflows  │ │  MCPs  │ │   External   │
   │ (playbooks) │ │(tools) │ │  Services    │
   └─────────────┘ └────────┘ └──────────────┘
```

Data flows in one direction: the workspace holds canonical records; agents read them, do work, and write updates back. The CLI is the primary interface for triggering this cycle. External services (Internet Archive, distributors, streaming platforms) are always at the boundary — always behind a human approval gate.

---

## Components

### 1. `cli/` — The `ovl` Command

**Responsibility:** Primary interaction surface. Wraps common operations into commands that validate input, update workspace records, and invoke agents when a process needs conversational guidance.

**Stack:** Node.js. CLI framework TBD (candidates: [ink](https://github.com/vadimdemedes/ink) for TUI components, [commander](https://github.com/tj/commander.js) or [citty](https://github.com/unjs/citty) for command parsing). Chosen to align with the Astro.js ecosystem and keep the toolchain consistent for contributors familiar with web development.

**What it reads:** Workspace JSON records, agent skill files (to know which agent to invoke for a given command), mastering profiles.

**What it writes:** Workspace JSON records (creating, updating status fields, adding archive paths, recording QC results).

**What it never does:** Sends emails, posts content, submits to distributors, or uploads files without explicit human confirmation at a named prompt step.

**Key commands (representative, not exhaustive):**

```bash
ovl init                              # scaffold a new workspace/
ovl status                            # show label-state summary
ovl artist create                     # interactive artist profile setup
ovl release create <title>            # new release record
ovl release list                      # show all releases and their stage
ovl track add <title> --release <id>  # add track, auto-link to release
ovl isrc assign --release <id>        # assign ISRCs to unregistered tracks
ovl mastering start --track <id>      # launch Mastering Companion session
ovl qc check --release <id>           # run pre-release QC checklist
ovl archive push --release <id>       # package and upload to archive
ovl release submit --release <id>     # distribution submission (gated)
ovl outreach research                 # trigger CRM research cycle
ovl outreach review                   # review pending opportunity drafts
ovl state show                        # print current label-state.md
ovl state sync                        # orchestrator writes session summary
```

**See:** `cli/README.md` for full command reference and installation.

---

### 2. `agents/` — Specialist Agent Definitions

**Responsibility:** Each agent is a markdown skill file that defines a role's purpose, interaction pattern, what records it reads and writes, and how it behaves with the artist. Agents are invoked by the CLI or directly in an LLM context (Claude Project, API, etc.).

**What they are:** Markdown documents, not code. The "execution" of an agent is an LLM reading its skill file and acting according to its instructions, with workspace data loaded as context.

**What they are not:** Autonomous processes. Agents do not run on a schedule, do not hold state between sessions, and do not take external actions without human approval.

**Interaction patterns** (each agent declares which it uses — see `docs/VISION.md`):
- `approval-gate` — prepares, human decides
- `guided-session` — human operates, agent advises in real time
- `review-and-refine` — agent drafts, human reviews and approves

#### Agent Roster

| Directory                     | Agent                   | Pattern           | Core Output                                     |
| ----------------------------- | ----------------------- | ----------------- | ----------------------------------------------- |
| `agents/orchestrator/`        | Orchestrator            | —                 | Routes requests, maintains `label-state.md`     |
| `agents/mastering-companion/` | Mastering Companion     | guided-session    | Mastering session notes written to `track.json` |
| `agents/archive/`             | Archive Agent           | approval-gate     | Archive paths written to `release.json`         |
| `agents/outreach-crm/`        | Outreach / CRM Agent    | approval-gate     | `opportunity.json` records, outreach drafts     |
| `agents/qc/`                  | QC Agent                | approval-gate     | QC report, `qc_passed` flag on each track       |
| `agents/metrics-analyst/`     | Metrics Analyst         | review-and-refine | Monthly metrics reports in `workspace/metrics/` |
| `agents/finance-manager/`     | Finance Manager         | review-and-refine | Revenue/expense summaries, goal progress        |
| `agents/content-strategist/`  | Content Strategist      | review-and-refine | Content calendars, release campaign briefs      |
| `agents/social-media/`        | Social Media Specialist | review-and-refine | Platform-specific copy drafts for human posting |

**Adding a new agent:** Create a directory under `agents/`, add a `README.md` skill file following the template in `agents/_template/`, and declare it in `agents/registry.json`. The orchestrator will route to it by name.

**See:** Individual `agents/[name]/README.md` files.

---

### 3. `schemas/` — JSON Schema Definitions

**Responsibility:** Canonical definitions for every structured data record in OVL. All agents, CLI commands, and workflow steps validate against these schemas. They are the shared contract between components.

**Format:** [JSON Schema Draft 7](https://json-schema.org/draft-07/schema). Each schema file includes a `$schema` field and a version string. Schemas are versioned independently; breaking changes increment the version and old records remain valid against their declared version.

**Validation:** The CLI validates workspace records against their schema on read and write. A standalone `ovl validate` command checks the entire workspace for schema compliance.

#### Schema Inventory

| File                                    | Describes                                                                         |
| --------------------------------------- | --------------------------------------------------------------------------------- |
| `schemas/artist.schema.json`            | Artist identity, platform IDs, PRO, distributor, default license                  |
| `schemas/release.schema.json`           | Release metadata, status, track list, archive paths, store links                  |
| `schemas/track.schema.json`             | Track metadata, ISRC, collaborators/splits, file paths, mastering data, QC status |
| `schemas/mastering-profile.schema.json` | Target LUFS/peak/LRA, platform notes, step checklist, session history             |
| `schemas/opportunity.schema.json`       | CRM record: type, status pipeline, contact, match score, outreach history         |
| `schemas/label.schema.json`             | Label identity, default settings, style guide reference                           |
| `schemas/finance-entry.schema.json`     | Revenue and expense records with period, source, amount, currency                 |
| `schemas/metrics-snapshot.schema.json`  | Periodic analytics data keyed by platform and period                              |

**See:** `schemas/README.md` for field-level documentation and validation examples.

---

### 4. `workflows/` — Process Playbooks

**Responsibility:** Step-by-step documentation of repeatable processes. Workflows describe the sequence of human actions, agent invocations, CLI commands, and approval gates for a given task. They are the operational manual.

**Format:** Markdown. Each workflow includes a prerequisites section, a numbered step sequence, notes on which agent handles each step, and what the expected output is. Where a CLI command executes a step, it is shown inline.

#### Workflow Inventory

| File                             | Process                                                                |
| -------------------------------- | ---------------------------------------------------------------------- |
| `workflows/release-pipeline.md`  | End-to-end: production → mastering → QC → archive → submit → campaign  |
| `workflows/mastering-session.md` | Step-by-step mastering with Mastering Companion for a target profile   |
| `workflows/outreach-loop.md`     | Research → score → draft → approve → send → follow-up → record outcome |
| `workflows/onboarding.md`        | First-time setup: `ovl init`, artist profile, first release scaffold   |
| `workflows/monthly-review.md`    | Metrics pull → analysis → finance summary → strategic check-in         |
| `workflows/commission-intake.md` | Inquiry → brief → quote → agreement → delivery → invoice               |
| `workflows/playlist-pitch.md`    | Research curators → prepare pitch → submit → track responses           |

**See:** `workflows/README.md` for an overview of how workflows relate to agents and CLI commands.

---

### 5. `workspace-scaffold/` — Starter Workspace Template

**Responsibility:** A template directory that `ovl init` copies into a new label's local `workspace/`. Every file contains placeholder values and inline comments explaining each field. This is the starting point for a new label operator — not a blank page.

**What it includes:**

```
workspace-scaffold/
├── label/
│   ├── profile.json              # label identity with all fields commented
│   └── styleguide/
│       └── README.md             # what to put here and why
├── artists/
│   └── _example-artist/
│       ├── artist.json           # all fields with explanatory comments
│       └── websites/
│           └── README.md
├── outreach/
│   ├── opportunities.json        # empty array with schema reference
│   └── contacts/
│       └── README.md
├── finance/
│   ├── revenue.json
│   └── expenses.json
├── metrics/
│   └── README.md                 # how metrics snapshots are structured
└── state/
    └── label-state.md            # orchestrator's session context template
```

The scaffold is intentionally minimal — enough to make the structure clear, not so much that it obscures what belongs to the operator.

**See:** `workspace-scaffold/README.md` for setup instructions.

---

### 6. `websites/` — Web Presence Components

**Responsibility:** Reusable Astro.js components, layouts, and configuration for label and artist websites. These live in the repo as shared infrastructure; each label's actual site content lives in `workspace/[label|artists]/websites/`.

**Stack:** [Astro.js](https://astro.build). Chosen for its content-first model (markdown/MDX sources, static output), alignment with the Node.js toolchain used by the CLI, and strong performance characteristics for largely static music sites.

**What it provides:**
- Base Astro project structure for a label site
- Base Astro project structure for an artist site
- Shared component library (track player, release card, bio block, contact form)
- Placeholder content and configuration referencing the workspace schema fields
- Build scripts that can read `artist.json` and `release.json` to generate catalog pages

**What it does not provide:** Design or visual identity. Each label's styleguide (colors, fonts, imagery) lives in `workspace/label/styleguide/` and is applied by the operator to the base components.

**Relationship to workspace:** The website components are wired to expect their content from the workspace. An artist site's catalog page, for example, is generated from `workspace/artists/[id]/releases/`. This keeps content in one place and the presentation layer separate.

**See:** `websites/README.md` for setup, build, and deployment instructions.

---

### 7. `mcp/` — Model Context Protocol Integrations

**Responsibility:** Optional MCP server definitions that extend agent capabilities with external service access. Each MCP is opt-in — agents function without them, but they unlock specific operations.

**What MCPs enable:**

| MCP                     | Capability                                                                   |
| ----------------------- | ---------------------------------------------------------------------------- |
| `mcp/internet-archive/` | Programmatic upload of release packages to archive.org via S3-compatible API |
| `mcp/amuse/`            | Distribution submission status checking (read-only where API permits)        |
| `mcp/gmail/`            | Outreach email drafting and sending (gated, requires approval)               |
| `mcp/calendar/`         | Scheduling follow-ups, release dates, review cycles                          |
| `mcp/google-drive/`     | Pulling analytics exports, syncing metrics snapshots                         |

MCPs that involve sending or publishing anything external always require an explicit human approval step defined in their configuration. The MCP cannot be invoked for a send action by an agent alone.

**Adding a new MCP:** Follow the MCP server definition format in `mcp/_template/` and register it in `mcp/registry.json`.

**See:** `mcp/README.md` for setup, authentication patterns, and the approval gate convention.

---

### 8. `docs/` — Project Documentation

**Responsibility:** Top-level conceptual documentation for the OVL project itself. Not operational instructions (those live in component directories) — these are the documents that explain the project to a new contributor or adopter.

| File                   | Purpose                                                          |
| ---------------------- | ---------------------------------------------------------------- |
| `docs/VISION.md`       | What OVL is, its philosophy, the agent model, workspace design   |
| `docs/ARCHITECTURE.md` | This document. Component inventory and cross-references          |
| `docs/CONTRIBUTING.md` | How to contribute: issues, PRs, adding agents, extending schemas |
| `docs/GOVERNANCE.md`   | Scope, license, what OVL is not, relationship to artists         |
| `docs/CHANGELOG.md`    | Version history, schema version bumps, breaking changes          |

---

## Data Flow: A Release, End to End

This trace shows how components interact for a typical release cycle. It is the architecture in motion.

```
1. Artist completes production in DAW
   └─ ovl release create "My Album"
      └─ Creates workspace/artists/brylie/releases/my-album/
         └─ release.json (status: in-production)

2. Artist adds tracks
   └─ ovl track add "Track One" --release my-album
      └─ Creates track.json, assigns local ID
   └─ ovl isrc assign --release my-album
      └─ Updates track.json with ISRC values

3. Mastering
   └─ ovl mastering start --track track-one
      └─ Loads mastering-profile.json for artist's genre
      └─ Launches Mastering Companion (guided-session)
         └─ Artist runs tools, agent advises on readings
         └─ Session ends → writes mastering{} block to track.json

4. QC
   └─ ovl qc check --release my-album
      └─ QC Agent checks all tracks: ISRC, mastering data, artwork specs
      └─ Produces report, sets qc_passed flags
      └─ Any failures block next step

5. Archive
   └─ ovl archive push --release my-album
      └─ [APPROVAL GATE] Shows package manifest, asks for confirmation
      └─ Archive Agent uploads to Internet Archive + object storage
      └─ Writes archive paths and checksums to release.json

6. Distribution submission
   └─ ovl release submit --release my-album --distributor amuse
      └─ [APPROVAL GATE] Shows full submission package for review
      └─ On approval, generates submission package or launches distributor flow
      └─ Updates release.json status: submitted

7. Post-release
   └─ ovl outreach research --release my-album
      └─ Outreach Agent researches placement opportunities
      └─ Creates opportunity.json records
      └─ [APPROVAL GATE] Artist reviews and approves outreach drafts
   └─ Content Strategist briefs Social Media Specialist
      └─ [REVIEW-AND-REFINE] Copy drafts delivered for artist to post
```

---

## Technology Decisions

| Decision          | Choice                         | Rationale                                                     |
| ----------------- | ------------------------------ | ------------------------------------------------------------- |
| CLI language      | Node.js                        | Aligns with Astro.js; single toolchain for contributors       |
| CLI framework     | TBD (ink / commander / citty)  | Evaluate for TUI quality and ESM compatibility                |
| Schema format     | JSON Schema Draft 7            | Wide tooling support, language-agnostic                       |
| Web framework     | Astro.js                       | Content-first, static output, Node ecosystem                  |
| Agent format      | Markdown skill files           | Human-readable, LLM-readable, no runtime dependency           |
| Archive primary   | Internet Archive               | Free, durable, aligns with CC/open values, public record      |
| Archive secondary | Object storage (S3-compatible) | Operator's choice of provider; bucket config in workspace     |
| License           | Apache 2.0                     | Permissive, patent grant, compatible with most downstream use |
| Package manager   | npm (default)                  | Widest compatibility; pnpm or yarn acceptable                 |

---

## Extension Points

OVL is designed to be extended without modifying core components.

**Adding an agent:** New directory in `agents/`, skill file following the template, register in `agents/registry.json`. No core code changes required.

**Adding a schema:** New file in `schemas/`, entry in `schemas/README.md`. CLI picks up new schemas automatically for the `ovl validate` command.

**Adding a workflow:** New markdown file in `workflows/`, entry in `workflows/README.md`. No code changes.

**Adding a CLI command:** New command file in `cli/commands/`, registered in `cli/index.js`. Should validate against schemas and follow the approval gate pattern for any external action.

**Adding an MCP:** New directory in `mcp/`, registered in `mcp/registry.json`. Must implement the approval gate interface for any write or send action.

**Forking for a new label:** Copy or clone the repo, run `ovl init` to create a `workspace/`, customize the scaffold. The repo itself never needs to be modified for label-specific content.

---

## What Is Intentionally Out of Scope

- **A hosted platform or SaaS layer.** OVL runs locally. There is no central server, no user accounts, no data collection.
- **Real-time collaboration between labels.** Each label runs its own instance. Sharing workflows or agent improvements happens through the open-source repo, not a shared runtime.
- **Automated social media posting.** The Social Media Specialist generates copy; the artist posts it. No platform API credentials are stored or used for posting.
- **Financial transactions.** The Finance Manager tracks records; it does not move money, issue invoices, or integrate with payment processors.
- **AI-generated music.** OVL supports human artists. It does not generate audio, lyrics, or compositions.

---

*Last updated: June 2026*
*Status: Pre-implementation draft — component interfaces subject to revision before v0.1*
