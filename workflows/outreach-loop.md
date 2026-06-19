# Outreach Loop

The full cycle for finding, researching, contacting, and following up with licensing, sync, playlist, commission, and collaboration opportunities. This is a continuous loop — it does not end with a single outreach campaign but runs persistently as the label's relationship-building engine.

**Core constraint:** No message is sent on the artist's behalf without explicit approval at a named gate. The CRM agent drafts; the artist decides.

---

## Prerequisites

- At least one `artist.json` with genre tags and platform links
- At least one released track or release with `status: live` (for sync and playlist outreach)
- `workspace/outreach/opportunities.json` exists (created by `ovl init`)

---

## The Loop

```text
Research → Score → [REVIEW] → Draft → [APPROVE] → Send → Track → Follow-up → [REVIEW] → Outcome
```

Each opportunity moves through `opportunity.status` as the loop progresses:
`identified → researched → draft-ready → approved → sent → follow-up → responded → won | lost | declined | stale`

---

## Phase 1: Research

```bash
ovl outreach research
```

→ `outreach-crm`

The CRM agent searches for opportunities matching the artist's profile. For each opportunity type, the approach differs:

**Sync licensing / podcast placement**

- Search podcast directories (Podchaser, Listen Notes) for shows in relevant niches (meditation, wellness, ambient, study)
- Review YouTube channels using background music in the artist's genre
- Check indie game development forums and itch.io for developers seeking ambient scores
- Look for documentary or short film projects with open music calls

**Playlist pitching**

- Identify independent playlist curators (not algorithmic) on Spotify and Apple Music
- Check playlist submission platforms (SubmitHub, Groover) for active curators in genre
- Review playlists that feature artists similar in style

**Commissions**

- Check community boards, composer forums, and social media for open briefs
- Monitor past contacts for new projects

**Collaboration**

- Identify artists with complementary styles who release under open or CC licensing
- Note any mutual connections or shared venue history

For each candidate, the agent creates or updates an `opportunity.json` record with `status: identified` and populates `contact{}` with what it has found.

**[REVIEW GATE]** The agent presents the identified opportunities as a list with brief descriptions and proposed match scores. The artist reviews and indicates which to pursue:

```text
Found 5 new opportunities:

1. Calm Waters Podcast (meditation, ~8k listeners) — match: 9/10
   Uses ambient piano beds, CC-friendly, consistent release schedule
   → Pursue?

2. Nordic Indie Devs collective — match: 7/10
   Several members seeking atmospheric game scores
   → Pursue?

3. [3 more...]

Which would you like to pursue? (Enter numbers, or 'all', or 'none')
```

Approved opportunities advance to `status: researched`. Declined ones are marked `status: declined` with a note.

---

## Phase 2: Deep Research

For each approved opportunity, the CRM agent does deeper research before drafting:

- Listen to or watch a sample of their work
- Note specific moments where the artist's music would fit
- Check for existing music credits (are they already using CC music? what style?)
- Find the right contact person and preferred contact method
- Note any recent projects, announcements, or context that makes outreach timely

This research is written into `opportunity.contact.notes`. The more specific the notes, the more personal the outreach.

No status change at this phase — still `researched` until a draft exists.

---

## Phase 3: Draft

```bash
ovl outreach draft --opportunity <opportunity-id>
```

→ `outreach-crm`

The CRM agent writes a personalised outreach message based on:

- The contact's work and audience
- Specific tracks from the artist's catalog that are a good fit (`tracks_suggested[]`)
- The value proposition (CC licensing, professional quality, mission alignment)
- Any timely hook from the deep research notes

The draft is never a template with fields filled in. It is written for this specific person.

The agent presents the draft to the artist with the opportunity context:

```
Draft for: Calm Waters Podcast
Contact: Sarah Chen, producer
Match: 9/10

---
Subject: Music for Calm Waters — free to use, no strings

Hi Sarah,

[draft body]
---

Suggested tracks: "Evening Meal", "Self-care"
Tone: conversational, value-first, no ask for money

Changes? Or approve to mark as ready?
```

`opportunity.status` advances to `draft-ready`. An `outreach_history` entry is logged: `action: draft-created`.

---

## Phase 4: Approve

**[APPROVAL GATE]** The artist reviews the draft. Three outcomes:

**Approve as-is:** `opportunity.status → approved`. History entry: `action: draft-approved, approved_by: [artist]`.

**Edit and approve:** Artist makes changes, confirms. History entry: `action: draft-edited`, then `action: draft-approved`. The final approved text is stored in the opportunity notes.

**Decline:** Opportunity is marked `declined`. No outreach sent.

The agent does not advance past this gate without a clear approval signal. Silence or ambiguity is treated as "not yet approved."

---

## Phase 5: Send

**[APPROVAL GATE — second confirmation]** Immediately before sending, the agent confirms one more time:

```text
Ready to send to Sarah Chen (sarahchen@calmwaters.fm)?
This will be sent from [artist email]. Confirm? [yes / no]
```

On confirmation:

- Message is sent (via configured email MCP or flagged for manual send if MCP not configured)
- `opportunity.status → sent`
- History entry: `action: sent, date: [today], approved_by: [artist]`
- `follow_up_due` is set to 14 days from send date (configurable)
- → `label-state.md`: opportunity noted as sent, follow-up date recorded in Open Loops

If no email MCP is configured, the agent outputs the final message for the artist to send manually and asks for confirmation that it was sent before updating the status.

---

## Phase 6: Track and Follow-up

The orchestrator surfaces follow-ups in the status report when `follow_up_due` is reached:

```text
Open loops:
· Follow-up due: Calm Waters Podcast — sent 14 days ago, no response
```

```bash
ovl outreach follow-up --opportunity <opportunity-id>
```

→ `outreach-crm`

The agent drafts a brief, non-pushy follow-up. Typically one follow-up only — if there is no response after the follow-up, the opportunity is marked `stale` and the loop ends.

**[APPROVAL GATE]** Follow-up draft approved before sending. Same process as Phase 4–5.

If a response is received at any point:

```bash
ovl outreach log-response --opportunity <opportunity-id>
```

→ `outreach-crm`

The agent asks the artist to describe or paste the response, then records it and advises on how to reply. `opportunity.status → responded`.

---

## Phase 7: Outcome

Every opportunity reaches a terminal status. The CRM agent prompts the artist to record the outcome when one is known:

```bash
ovl outreach close --opportunity <opportunity-id> --outcome won|lost|declined
```

**Won:** `opportunity.status → won`. The agent prompts for:

- Actual value if financial (`value_estimate.amount` confirmed or corrected)
- Which tracks were used
- Whether to follow up for a testimonial or future use

**Lost / Declined:** `opportunity.status → lost` or `declined`. The agent asks one question: "Any reason given? (Helps improve future outreach — skip if not known)". Answer appended to notes.

**Stale:** No response after follow-up. `opportunity.status → stale`. The agent may suggest re-approaching in a future cycle if circumstances change.

All outcomes feed into the Metrics Analyst's outreach section of the monthly snapshot:

- `outreach.outreach_sent`
- `outreach.responses_received`
- `outreach.placements_won`

---

## Running Outreach as a Regular Practice

Outreach works through consistency more than volume. A sustainable cadence:

**Weekly (5–10 minutes):**

- `ovl outreach review` — check for pending approvals and overdue follow-ups

**Monthly (30–60 minutes):**

- `ovl outreach research` — find new opportunities
- Review and approve new drafts

**Quarterly:**

- Review win/loss outcomes with the Metrics Analyst
- Adjust targeting based on what has and hasn't worked

The outreach loop runs alongside releases, not only during them. Relationships built between releases are often the ones that convert when a new release drops.
