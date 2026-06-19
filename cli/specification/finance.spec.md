# CLI Specification — `ovl finance`

Commands for logging and summarising revenue and expenses.
Finance entries are stored in `workspace/finance/revenue.json` and
`workspace/finance/expenses.json`. Summaries and quotes invoke the `finance-manager`
agent.

---

## Command Reference

---

### `ovl finance add-revenue`

Log a revenue entry.

```text
ovl finance add-revenue \
  --source <source> \
  --amount <amount> \
  --currency <EUR|USD|...> \
  --period <YYYY-MM> \
  [--artist <artist-id>] \
  [--release <release-id>] \
  [--opportunity <opportunity-id>] \
  [--description "<text>"]
```

**Behaviour:** Creates a `finance-entry.json` record with `type: revenue` and appends
it to `workspace/finance/revenue.json`. Validates against
`schemas/finance-entry.schema.json` before writing.

---

### `ovl finance add-expense`

Log an expense entry.

```text
ovl finance add-expense \
  --source <category> \
  --amount <amount> \
  --currency <EUR|USD|...> \
  --date <YYYY-MM-DD> \
  [--artist <artist-id>] \
  [--description "<text>"]
```

**Behaviour:** Creates a `finance-entry.json` record with `type: expense` and appends
it to `workspace/finance/expenses.json`. Validates against
`schemas/finance-entry.schema.json` before writing.

---

### `ovl finance summary --period <YYYY-MM>`

Generate a financial summary for a period.

```text
ovl finance summary --period <YYYY-MM> [--brief]
```

**Behaviour:**

- Invokes the `finance-manager` agent
- Agent reads all revenue and expense entries for the period
- Produces a summary: total revenue by source, total expenses by category, net position,
  goal progress, 3-month and 12-month trends
- `--brief` produces a one-paragraph summary instead of the full report

---

### `ovl finance quote --opportunity <opportunity-id>`

Generate a pricing quote for a commission opportunity.

```text
ovl finance quote --opportunity <opportunity-id>
```

**Behaviour:** Invokes the `finance-manager` agent to review the commission scope and
propose a rate. Artist confirms before the quote is used in outreach.
