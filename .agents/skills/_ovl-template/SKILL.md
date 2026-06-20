---
name: ovl-template
description: >
  Template for OVL agent skill files. Use this as the starting point when
  creating a new specialist agent for the Open Virtual Label system. Not
  intended to be invoked directly.
license: Apache-2.0
metadata:
  author: open-virtual-label
  version: "0.1.0"
  interaction-pattern: "approval-gate | guided-session | review-and-refine"
  reads: "workspace/state/label-state.md"
  writes: "workspace/state/label-state.md"
---

# [Agent Name]

One sentence: what this agent does and why it exists in OVL.

## Role

Two to four sentences describing this agent's function in the label.
Written in first person: "I am the X agent. My job is to Y."
Establishes tone and scope before any procedural detail.

## Interaction Pattern

State which pattern this agent uses and what it means concretely.

**Pattern:** `approval-gate` | `guided-session` | `review-and-refine`

- **approval-gate** — agent prepares everything; artist explicitly approves before any external action occurs
- **guided-session** — artist operates their own tools; agent interprets output and advises in real time
- **review-and-refine** — agent produces a complete draft; artist reviews, edits if needed, then approves

Describe what the artist does, what the agent does, and exactly where the gate or review moment falls.

## Inputs

What this agent reads before beginning work.

| Source                              | What It Provides                                 |
| ----------------------------------- | ------------------------------------------------ |
| `workspace/state/label-state.md`    | Current label context, open loops, pending items |
| `workspace/label/profile.json`      | Label identity and defaults                      |
| _(add rows specific to this agent)_ |                                                  |

If the CLI passes a specific target (e.g. `--release <id>`, `--track <id>`),
state what additional records are loaded.

## Outputs

What this agent produces or modifies. Every output traces to a schema record or named document.

| Output                 | Destination                         | Condition               |
| ---------------------- | ----------------------------------- | ----------------------- |
| _(e.g. session notes)_ | _(e.g. `track.json → mastering{}`)_ | After session completes |
| Session summary        | `label-state.md → ## Session Log`   | Always                  |

Mark any output requiring human approval before being written as `[GATED]`.

## Session Protocol

How a session unfolds from invocation to close.

1. **Load context.** Read `label-state.md` and relevant workspace records.
2. **Orient.** Confirm the session goal with the artist; surface any relevant state.
3. _(agent-specific steps — be concrete about what the agent says and does)_
4. **Confirm outputs.** Before writing any record, summarize changes and ask for confirmation.
5. **Write outputs.** Update workspace records.
6. **Update state.** Append session summary to `label-state.md`: what was done, open loops, next action.

## Boundaries

What this agent explicitly does not do.

- Does not take external actions (send, post, upload, submit) without artist approval
- Does not modify records outside its declared Outputs
- _(add agent-specific limits)_

When a request falls outside scope, name the appropriate agent and suggest invoking it.

## Related Agents

| Agent             | Relationship                                          |
| ----------------- | ----------------------------------------------------- |
| `orchestrator`    | Routes sessions here; receives state updates on close |
| _(related agent)_ | _(e.g. "hands off to archive after QC passes")_       |

---

For detailed reference on OVL data schemas, workspace structure, and interaction
pattern conventions, see [references/OVL-CONTEXT.md](references/OVL-CONTEXT.md).
