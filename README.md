# Open Virtual Label

A framework for independent musicians and small collectives to operate with the support of a coordinated team of AI agents — without ceding creative control, revenue, or identity to a traditional label.

OVL is three things working together:

- **`ovl`** — a Go CLI that manages a local workspace of artist, release, and track records, validates them against JSON schemas, and hands off to AI agents for anything that needs judgment.
- **Agent skills** (`.agents/skills/`) — markdown-defined specialists (mastering companion, outreach CRM, QC, finance manager, etc.) that an AI assistant loads to do the conversational parts of the work.
- **Workflows and schemas** (`workflows/`, `schemas/`) — the playbooks and data contracts that keep agents and CLI commands consistent.

Read [docs/VISION.md](docs/VISION.md) for the why, and [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the full component map.

---

## Repository layout

| Path | Contents |
| --- | --- |
| `cli/` | The `ovl` CLI, written in Go. See [cli/specification.md](cli/specification.md) for the full command reference. |
| `.agents/skills/` | Agent skill definitions (`SKILL.md` per agent), loaded by the CLI or an AI assistant. |
| `schemas/` | JSON Schema (Draft 7) definitions for every workspace record. See [schemas/README.md](schemas/README.md). |
| `workflows/` | Step-by-step playbooks tying CLI commands and agents together. See [workflows/README.md](workflows/README.md). |
| `docs/` | Vision and architecture documentation. |
| `workspace/` | Where a label's actual data lives once you run `ovl init`. Gitignored — this is private artist data, not part of the framework. |

> Note: `cli/specification.md` describes the CLI design in Node.js terms from the original spec, but the implementation is Go. Go was chosen for a single static binary with no runtime dependency for end users.

---

## Getting started

### Prerequisites

This project uses [mise](https://mise.jdx.dev/) to pin tool versions (Go, markdownlint-cli2). Install mise, then from the repo root:

```bash
mise install
```

This provisions Go 1.26 and `markdownlint-cli2` as declared in [mise.toml](mise.toml). `golangci-lint` is also configured via mise but resolves to `latest`.

### Build the CLI

```bash
cd cli
mise exec -- go build -o ovl .
```

Verify it runs:

```bash
./ovl --help
```

> If you build with a system `go` instead of `mise exec -- go`, you may hit a `dyld: missing LC_UUID load command` error on macOS. Always build and test through `mise exec --`.

### Run the test suite

```bash
cd cli
mise exec -- go test ./...
```

### Cross-compile release binaries

```bash
cd cli
./scripts/build.sh
```

Builds `ovl` for `linux/amd64`, `linux/arm64`, `darwin/amd64`, and `darwin/arm64`, writing each binary to `build/` at the repo root (e.g. `build/ovl-linux-amd64`). `build/` is gitignored.

`ovl` reads JSON schemas from a `schemas/` directory at runtime rather than embedding them, discovered the same way as `workspace/`: an explicit `--schemas <path>` flag, the `OVL_SCHEMAS_DIR` environment variable, or by walking up from the current directory looking for a `schemas/` folder. Binaries run from inside this repo find it automatically; if you copy a binary elsewhere, pass `--schemas` or set `OVL_SCHEMAS_DIR` to point at a copy of this repo's `schemas/` directory.

### Lint

```bash
cd cli
mise exec -- golangci-lint run ./...
```

Markdown files (docs, workflows, schemas, agent skills) are linted with `markdownlint-cli2` using the rules in [.markdownlint.yaml](.markdownlint.yaml):

```bash
mise exec -- markdownlint-cli2 "**/*.md"
```

### Try it on a real workspace

```bash
cd cli
mise exec -- go run . init --workspace ../workspace
```

`init` prompts interactively for label name, contact email, default license, and primary distributor, then scaffolds `workspace/`. After that:

```bash
mise exec -- go run . status --workspace ../workspace
```

`ovl` walks up from the current directory looking for a `workspace/` directory by default, so once one exists you can drop `--workspace` for subsequent commands run from the repo root.

---

## Working with agents

Agent skills under `.agents/skills/` are markdown files (`SKILL.md`) describing a specialist's role, inputs, outputs, and interaction pattern (`approval-gate`, `guided-session`, or `review-and-refine`). They're written to be loaded by an AI coding assistant (e.g. Claude Code) alongside the CLI — the CLI manages workspace state, the agent provides judgment and drafts that a human approves before anything external happens.

Start with the `coordinator` skill ([.agents/skills/coordinator/SKILL.md](.agents/skills/coordinator/SKILL.md)) — it's the session entry point that reads label state and routes to the right specialist. Use [.agents/skills/_template/SKILL.md](.agents/skills/_template/SKILL.md) as the starting point for writing a new specialist.

---

## Contributing

This is an early-stage framework. Schemas, the CLI, and agent skills are expected to evolve together — if you change a schema, check whether the matching Go model in `cli/internal/models/`, the CLI commands that touch it, and any agent skill that reads or writes it need updating too.
