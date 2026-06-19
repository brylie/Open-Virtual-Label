# CLI Specification — `ovl metrics`

Commands for compiling and analysing platform performance metrics.
Raw export files from streaming platforms are placed in
`workspace/metrics/[YYYY-MM]/raw/` and processed by the `metrics-analyst` agent.

---

## Command Reference

---

### `ovl metrics snapshot --period <YYYY-MM>`

Compile a metrics snapshot for a period.

```text
ovl metrics snapshot --period <YYYY-MM> [--artist <artist-id>] [--brief]
```

**Behaviour:**

- Invokes the `metrics-analyst` agent
- Agent reads platform export files from `workspace/metrics/[YYYY-MM]/raw/`
- Populates `workspace/metrics/[YYYY-MM]/[artist-id].json`
- Produces a written analysis with trends, top tracks, and anomalies
- `--brief` produces a one-paragraph summary instead of the full report

**Options:**

| Option | Description |
|---|---|
| `--artist <id>` | Scope snapshot to one artist. If omitted, covers all artists in the workspace. |
| `--brief` | Output a one-paragraph summary instead of the full report. |

**Errors:**

- `NO_RAW_DATA` — no export files found in `workspace/metrics/[YYYY-MM]/raw/`
