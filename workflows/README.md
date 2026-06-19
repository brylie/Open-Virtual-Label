# Workflows

Step-by-step playbooks for repeatable OVL processes. Each document describes a complete workflow: the prerequisite state, the sequence of human actions and agent invocations, the approval gates, and the expected output.

Workflows reference specific agents by name and CLI commands by exact syntax. When a workflow step says `ovl [command]`, the artist runs that command. When it says `→ [agent]`, the orchestrator routes to that agent.

---

## Workflow Index

| File                   | Process                                                         | Primary Agents                                 |
| ---------------------- | --------------------------------------------------------------- | ---------------------------------------------- |
| `release-pipeline.md`  | End-to-end release: production through archive and distribution | orchestrator, mastering-companion, qc, archive |
| `mastering-session.md` | Single-track mastering session with target profile              | mastering-companion                            |
| `outreach-loop.md`     | Research → draft → approve → send → follow-up cycle             | outreach-crm                                   |
| `monthly-review.md`    | Periodic metrics, finance, and strategic check-in               | metrics-analyst, finance-manager, orchestrator |
| `commission-intake.md` | Handling an inbound commission inquiry end to end               | outreach-crm, finance-manager                  |
| `playlist-pitch.md`    | Researching and submitting to playlist curators                 | outreach-crm, metrics-analyst                  |
| `onboarding.md`        | First-time label setup and first release scaffold               | orchestrator                                   |

---

## Reading a Workflow

Each workflow uses consistent notation:

- `ovl [command]` — run this CLI command
- `→ [agent-name]` — the orchestrator routes to this agent
- `[APPROVAL GATE]` — no action proceeds past this point without explicit artist confirmation
- `→ label-state.md` — this step updates the state document
- `✓` — this step produces a verifiable output (a record written, a file created, a confirmation received)

Prerequisites are listed at the top of each document. A workflow should not be started if its prerequisites are unmet — the orchestrator will flag this.

---

## Adding a New Workflow

1. Create `workflows/[name].md` following the structure of an existing workflow
2. Add an entry to this README
3. If the workflow introduces a new multi-agent sequence, add it to `agents/orchestrator/references/ROUTING.md`
