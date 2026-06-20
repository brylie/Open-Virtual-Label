# Routing Reference

Loaded by the orchestrator when routing logic needs detail beyond what is
held in the main SKILL.md. Not loaded at session start — only when a
routing decision requires this reference.

---

## Routing Table

| Request type                                           | Agent                           |
| ------------------------------------------------------ | ------------------------------- |
| Mastering questions, "master track X"                  | `mastering-companion`           |
| "Run QC", pre-release checks, release completeness     | `qc`                            |
| "Archive release", backup, Internet Archive upload     | `archive`                       |
| "Research outreach", find opportunities, review drafts | `outreach-crm`                  |
| "What are my numbers", metrics, analytics, trends      | `metrics-analyst`               |
| Revenue, expenses, income goals, financial summary     | `finance-manager`               |
| Content calendar, release campaign, posting schedule   | `content-strategist`            |
| Social copy, captions, announcements                   | `social-media`                  |
| Status, what to work on, planning                      | Orchestrator (direct)           |
| `ovl` CLI operations (create release, add track, etc.) | CLI — tell artist which command |

---

## Ambiguous Requests

If a request does not clearly map to one agent, ask one clarifying question
before routing. Do not guess and silently hand off.

Examples of ambiguous requests and how to clarify:

**"Help me with my release"**
→ Ask: "Are you working on production, mastering, preparing to submit, or
planning the content campaign?"

**"I need to do some outreach"**
→ Ask: "Are you looking to research new opportunities, review a draft that's
waiting, or follow up on something already sent?"

**"Can you look at my numbers?"**
→ This usually means metrics-analyst, but if the artist says "revenue" or
"how much have I earned", route to finance-manager instead.

---

## Multi-Agent Sequences

When a request spans multiple agents, propose the sequence before beginning.
Common sequences:

### New release to submission

1. `mastering-companion` — complete mastering for all tracks
2. `qc` — run pre-release QC checklist
3. `archive` — package and upload masters, stems, project files
4. Orchestrator presents distribution package → artist approves
5. `ovl release submit` — artist runs CLI command
6. `content-strategist` — brief release campaign
7. `outreach-crm` — research placement opportunities

### Monthly review
1. `metrics-analyst` — compile analytics snapshot
2. `finance-manager` — revenue and expense summary
3. Orchestrator — strategic check-in, update goals and open loops

### Commission intake

1. `outreach-crm` — record inquiry, create opportunity
2. Orchestrator — confirm scope and timeline with artist
3. `finance-manager` — agree rate and create expense/income projection
4. `outreach-crm` — draft agreement for artist review [GATED]

---

## Handoff Format

When handing off to a specialist, provide this context summary so the
agent does not need to re-read the full state document:

```
Handing off to [agent-name].

Context:
· Artist: [display name], working on [release title if relevant]
· Relevant state: [one or two lines from open loops or active projects]
· Session goal: [what the artist asked for]
· Any pending items from this agent: [if returning to it]
```

---

## After a Specialist Session

When a specialist completes and returns:

1. Receive the session output (what was done, what was written, new open loops).
2. Incorporate into the pending state update.
3. If the sequence has more steps, propose the next one.
4. If the session is complete, close with the full summary for artist review.
