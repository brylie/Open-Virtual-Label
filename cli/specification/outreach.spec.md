# CLI Specification — `ovl outreach`

Commands for managing the outreach pipeline: researching opportunities, drafting and
sending messages, logging responses, and closing outcomes. All outreach state is stored
in `workspace/artists/[artist-id]/opportunities/[opportunity-id].json`.

Outreach commands invoke the `outreach-crm` agent for research, drafting, scoring,
and response interpretation. Sending requires an approved draft and explicit confirmation.

---

## Command Reference

---

### `ovl outreach research`

Trigger the CRM agent to find new opportunities.

```text
ovl outreach research [--type <type>] [--release <release-id>] [--track <track-id>]
```

**Behaviour:**

- Invokes the `outreach-crm` agent in research mode
- Agent searches for opportunities matching the artist's genre tags and platform presence
- Presents a list of candidates with match scores for artist review
- **[REVIEW GATE]** Artist selects which opportunities to pursue
- Creates `opportunity.json` records for approved candidates with `status: identified`

**Options:**

| Option | Description |
|---|---|
| `--type <type>` | Limit to one opportunity type: `sync-license`, `commission`, `playlist-pitch`, `collaboration`, `press` |
| `--release <release-id>` | Focus research on placement opportunities for a specific release |
| `--track <track-id>` | Focus on a specific track (e.g. for playlist pitching) |

---

### `ovl outreach review`

Review all opportunities with pending approvals.

```text
ovl outreach review [--type <type>]
```

**Behaviour:**

- Lists all opportunities with `status: draft-ready` (outreach drafts awaiting approval)
- Invokes the `outreach-crm` agent to present each draft for review
- Artist approves, edits, or declines each draft
- Approved drafts advance to `status: approved`

---

### `ovl outreach draft --opportunity <opportunity-id>`

Draft outreach for a specific opportunity.

```text
ovl outreach draft --opportunity <opportunity-id>
```

**Behaviour:**

- Invokes the `outreach-crm` agent
- Agent reads opportunity record including contact notes and suggested tracks
- Produces a personalised outreach message
- Presents draft for artist review
- **[APPROVAL GATE]** Advances to `status: draft-ready` on approval

---

### `ovl outreach send --opportunity <opportunity-id>`

Send an approved outreach message.

```text
ovl outreach send --opportunity <opportunity-id>
```

**Behaviour:**

- Confirms `opportunity.status` is `approved`
- **[CONFIRMATION GATE]** Shows recipient, subject, and message. Requires explicit
  confirmation. Cannot be bypassed with `--yes`
- Sends via email MCP if configured, or outputs the message for manual sending
- Updates `status: sent`, logs `action: sent` in `outreach_history`, sets
  `follow_up_due`

**Errors:**

- `NOT_APPROVED` — draft has not been approved; run `ovl outreach draft` first
- `EMAIL_MCP_NOT_CONFIGURED` — email MCP not connected; message will be output for
  manual sending

---

### `ovl outreach follow-up --opportunity <opportunity-id>`

Draft a follow-up for an opportunity with no response.

```text
ovl outreach follow-up --opportunity <opportunity-id>
```

**Behaviour:**

- Invokes the `outreach-crm` agent to draft a brief follow-up
- **[APPROVAL GATE]** Same approval and send flow as `ovl outreach draft` +
  `ovl outreach send`

---

### `ovl outreach log-response --opportunity <opportunity-id>`

Record a response received from an outreach contact.

```text
ovl outreach log-response --opportunity <opportunity-id>
```

**Behaviour:**

- Prompts for the response content (paste or describe)
- Invokes the `outreach-crm` agent to interpret and advise on next steps
- Logs `action: response-received` in `outreach_history`
- Updates `status: responded`

---

### `ovl outreach score --opportunity <opportunity-id>`

Score an opportunity for fit with the artist.

```text
ovl outreach score --opportunity <opportunity-id>
```

**Behaviour:** Invokes the `outreach-crm` agent to review the opportunity details and
propose a match score (1–10) with rationale. Artist confirms or adjusts. Writes
`opportunity.match{}`.

---

### `ovl outreach intake --type <type>`

Record a new inbound inquiry as an opportunity.

```text
ovl outreach intake --type <type>
```

**Behaviour:**

- Prompts for contact details, inquiry description, source, and any deadline
- Creates an `opportunity.json` record with `status: identified`
- Logs `action: identified` in `outreach_history`

---

### `ovl outreach close --opportunity <opportunity-id> --outcome <outcome>`

Record the final outcome of an opportunity.

```text
ovl outreach close --opportunity <opportunity-id> --outcome won|lost|declined
```

**Behaviour:**

- Sets `opportunity.status` to the specified outcome
- For `won`: prompts for confirmed value and which tracks were used
- For `lost` / `declined`: prompts for reason if known
- Logs the outcome action in `outreach_history`

---

### `ovl outreach log --opportunity <opportunity-id> --action <action>`

Manually log an action on an opportunity.

```text
ovl outreach log --opportunity <opportunity-id> --action <action> [--note "<note>"]
```

**Behaviour:** Appends an entry to `opportunity.outreach_history[]`. For recording
manual actions (platform-based submissions, phone calls, in-person conversations).

**Valid actions:** Any value from the `outreach_history.action` enum in
`opportunity.schema.json`.
