# CLI Specification — `ovl commission`

Shortcut commands for commission workflows. Commission opportunities are created and
managed through the standard `ovl outreach` commands; this group provides shortcuts
for the commission-specific steps that don't fit cleanly into the generic outreach flow.

---

## Command Reference

---

### `ovl commission agreement --opportunity <opportunity-id>`

Generate a commission agreement from the workspace template.

```text
ovl commission agreement --opportunity <opportunity-id>
```

**Behaviour:**

- Reads the opportunity record for scope, rights, timeline, and payment terms
- Populates `workspace/label/templates/commission-agreement.md` with the agreed details
- Outputs the filled agreement for artist review
- **[APPROVAL GATE]** Artist approves before the agreement is sent to the client

**Errors:**

- `OPPORTUNITY_NOT_FOUND` — no opportunity with the given ID exists
- `WRONG_TYPE` — opportunity is not of type `commission`
- `TEMPLATE_NOT_FOUND` — `workspace/label/templates/commission-agreement.md` does not
  exist; create it first
