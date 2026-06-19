---
name: project-cli-go
description: Go implementation of the ovl CLI under cli/ — structure, build, and status
metadata:
  type: project
---

The `ovl` CLI is implemented in Go (not Node.js as the spec originally said) under `cli/`.

**Module:** `github.com/open-virtual-label/ovl`  
**Build:** `mise exec -- go build -o ovl .` from `cli/` (requires Go 1.26 via mise)  
**Dependencies:** cobra, santhosh-tekuri/jsonschema/v5, olekukonko/tablewriter

**Structure:**
- `cli/cmd/` — one file per command group (root, init, status, validate, state, artist, release, track, isrc, mastering, qc, archive, outreach, finance, metrics, content, agents, mcp, commission)
- `cli/internal/workspace/` — workspace discovery, path helpers, JSON I/O
- `cli/internal/models/` — Go structs matching all JSON schemas
- `cli/internal/schema/` — JSON Schema Draft 7 validator with embedded `data/` copies of schemas
- `cli/internal/output/` — text/JSON output helpers, table rendering
- `cli/internal/prompt/` — stdin prompts and confirmation gates

**Fully implemented:** init, status, validate, state show/sync, artist CRUD, release CRUD + pipeline transitions, track CRUD, isrc assign, mastering profile create/list, qc check, finance add-revenue/expense/summary, agents list  
**Agent stubs:** mastering start, outreach all, metrics snapshot, content brief, social draft, commission agreement, archive push (needs IA MCP), mcp connect  

**Why:** User specified Go instead of the Node.js mentioned in the spec.
