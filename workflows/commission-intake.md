# Commission Intake

Handling an inbound request for original commissioned music — from the initial inquiry through brief, agreement, production, delivery, and invoicing. This workflow ensures every commission is scoped, agreed, and documented before production begins.

---

## Prerequisites

- `artist.json` exists with contact email
- `workspace/outreach/opportunities.json` exists

---

## Stage 1: Receive and Record the Inquiry

When a commission inquiry arrives (by email, social message, or contact form), record it immediately before responding:

```bash
ovl outreach intake --type commission
```

→ `outreach-crm`

The CRM agent asks for:

- Who is making the request (name, role, organisation, contact details)
- What they are describing (brief description in their own words)
- How they found the artist
- Any deadline they mentioned

An `opportunity.json` record is created with `type: commission`, `status: identified`, and a history entry `action: identified`.

---

## Stage 2: Assess Fit

Before responding, assess whether the commission is worth pursuing:

**Questions to consider:**

- Does the project align with the artist's values and mission?
- Is the timeline realistic given current commitments?
- Is the scope clear enough to quote, or does it need significant clarification first?
- Is there a budget, and if so, is it in range?

```bash
ovl outreach score --opportunity <opportunity-id>
```

→ `outreach-crm`

The CRM agent reviews the inquiry details and proposes a match score (1–10) with rationale. The artist confirms or adjusts. If the score is below 5, the agent drafts a polite decline for artist approval.

If pursuing: `opportunity.status → researched`

---

## Stage 3: Brief the Commission

If the initial inquiry lacked detail, gather a proper brief before quoting. The CRM agent drafts a brief request message:

**[APPROVAL GATE]** Artist approves the brief request before it is sent.

Key questions a music commission brief should answer:

- What is the project? (game, film, podcast, commercial, installation)
- What is the intended use and distribution? (affects licensing scope)
- How long is the piece, or how many pieces are needed?
- What is the mood, energy, and style?
- Are there reference tracks?
- What is the delivery format? (WAV stems, MP3, specific sample rate)
- What is the deadline?
- What is the budget, or is there flexibility to discuss?

Once the brief is received, it is recorded in `opportunity.contact.notes` and the relevant track suggestions are updated in `opportunity.tracks_suggested[]`.

---

## Stage 4: Scope and Quote

With a clear brief, the Finance Manager helps establish a rate.

```bash
ovl finance quote --opportunity <opportunity-id>
```

→ `finance-manager`

The Finance Manager reviews:

- Estimated production time (hours × rate)
- Licensing scope (how broadly will this music be used?)
- Revision rounds included
- Any rights transferred or retained

It produces a proposed quote range and a licensing scope summary. The artist adjusts and confirms the final figure.

Factors that affect pricing upward:

- Exclusive rights (artist cannot re-license the music)
- Broad distribution (commercial, broadcast, international)
- Short deadline requiring prioritisation
- Complex delivery requirements (multiple stems, many edits)

Factors that may adjust pricing:

- Non-commercial or charitable use
- Long-term relationship with an established contact
- Interesting creative brief that has portfolio value

**[APPROVAL GATE]** Artist confirms the quote and licensing scope before it is sent.

The CRM agent drafts a quote and brief summary for the client:

```
Dear [name],

Thank you for the details on [project]. Here's what I'm proposing:

Scope: [x] pieces, approximately [duration] each
Delivery: [format], [date]
Revisions: [n] rounds included
Licensing: [scope — e.g. non-exclusive, perpetual, for use in [project name]]
Investment: €[amount]

[Any questions or clarifications]

Let me know if you'd like to proceed or discuss further.
```

**[APPROVAL GATE]** Quote message approved before sending. History entry: `action: draft-approved`.

---

## Stage 5: Agreement

If the client accepts the quote, formalise the agreement before production begins.

For straightforward commissions with clear scope, a brief written confirmation by email is sufficient — the client's written acceptance of the quote terms serves as the agreement. Store the confirmation in `opportunity.contact.notes`.

For larger commissions (higher value, exclusive rights, significant production time):

```bash
ovl commission agreement --opportunity <opportunity-id>
```

→ `outreach-crm`

The CRM agent generates a simple commission agreement from the `workspace/label/templates/commission-agreement.md` template, populated with the agreed scope, rights, timeline, and payment terms. The artist reviews it.

**[APPROVAL GATE]** Agreement reviewed and approved before sending to client.

`opportunity.status → approved` (repurposed here to mean: commission accepted and agreed)

A revenue entry is created in advance for tracking purposes:

```bash
ovl finance add-revenue \
  --source commission \
  --amount <amount> \
  --currency EUR \
  --date <expected-payment-date> \
  --opportunity <opportunity-id> \
  --description "<client name> — <project name>"
```

This is flagged as `pending` until payment is received.

---

## Stage 6: Production

The commission is produced like any other track or release. Create a release record if the work will be archived:

```bash
ovl release create "<commission title>" --artist <artist-id> --type single
ovl track add "<track title>" --release <release-id>
```

The Open Loops section of `label-state.md` carries the commission with its delivery deadline:

```
· Commission: [client] "[project]" — delivery due [date]
```

Track production progress through regular `ovl status` checks with the orchestrator.

---

## Stage 7: Delivery

When the commission is complete, prepare delivery files according to the agreed spec.

```bash
ovl mastering start --track <track-id>
```

→ `mastering-companion`

Master to the agreed format (which may differ from standard streaming targets — e.g. a podcast bed may need different LUFS treatment than a streaming release).

The CRM agent drafts a delivery message for the artist's approval:

```
Hi [name],

[Project name] is ready — files attached / linked below.

Included:
· [file list]

[Any notes about the delivery — stems, edit versions, etc.]

Please let me know if anything needs adjusting within the agreed revision scope.
```

**[APPROVAL GATE]** Delivery message approved before sending.

`outreach_history` entry: `action: sent` (repurposed for delivery). `opportunity.status → responded` (awaiting client confirmation).

---

## Stage 8: Revision and Sign-off

If the client requests revisions within the agreed scope, note them in `opportunity.contact.notes` and update the track accordingly. Re-run the mastering session if the changes affect the audio.

When the client confirms acceptance:

- `opportunity.status → won`
- History entry: `action: won`
- Revenue entry updated from `pending` to confirmed

---

## Stage 9: Archive

Commission work is archived the same as any release:

```bash
ovl archive push --release <release-id>
```

Commissioned masters, stems, and project files are stored permanently. Even if the rights are fully transferred to the client, keeping an archival copy protects the artist against future disputes.

---

## Rights Retained vs. Transferred

By default, OVL tracks use CC BY 4.0. A commission may require different terms:

**Non-exclusive license** (recommended default): The artist retains copyright and licenses use for the specific project. The music can be re-licensed for other uses. Note in `opportunity.notes`.

**Exclusive license**: The client has exclusive use rights for a defined period or in perpetuity. The artist retains copyright but cannot re-license during the exclusivity term. Reflected in `opportunity.value_estimate.basis`.

**Full rights transfer (work for hire)**: The client owns the copyright. Rare; requires significantly higher compensation. Not recommended without legal advice.

Any rights arrangement other than the label's default license should be documented in the agreement and noted on the track record.
