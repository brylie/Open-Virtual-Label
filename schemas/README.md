# Schemas

JSON Schema definitions for all structured data in OVL. Every workspace record validates against one of these schemas. Agents read them to understand field names and types; the CLI validates against them on read and write.

All schemas use [JSON Schema Draft 7](https://json-schema.org/draft-07/schema).

---

## Schema Index

| File                            | Title            | Stored At                                                         |
| ------------------------------- | ---------------- | ----------------------------------------------------------------- |
| `label.schema.json`             | Label            | `workspace/label/profile.json`                                    |
| `artist.schema.json`            | Artist           | `workspace/artists/[id]/artist.json`                              |
| `release.schema.json`           | Release          | `workspace/artists/[id]/releases/[id]/release.json`               |
| `track.schema.json`             | Track            | `workspace/artists/[id]/releases/[id]/tracks/[id].json`           |
| `mastering-profile.schema.json` | MasteringProfile | `workspace/artists/[id]/mastering-profiles/[id].json`             |
| `opportunity.schema.json`       | Opportunity      | Items within `workspace/outreach/opportunities.json`              |
| `finance-entry.schema.json`     | FinanceEntry     | Items within `workspace/finance/revenue.json` and `expenses.json` |
| `metrics-snapshot.schema.json`  | MetricsSnapshot  | `workspace/metrics/[YYYY-MM]/[artist-id].json`                    |

---

## Versioning

Every schema has a `schema_version` field. Workspace records declare the version they conform to. The current version of every schema is `"1"`.

When a schema changes in a backwards-incompatible way:
- The `const` value in `schema_version` increments
- A migration note is added to `CHANGELOG.md`
- Old records remain valid against their declared version
- The CLI `ovl validate` reports records using older schema versions

Additive changes (new optional fields) do not increment the version.

---

## Validation

```bash
# Validate a single file against its schema
ovl validate workspace/artists/brylie/artist.json

# Validate all workspace records
ovl validate --all

# Validate a specific schema against the spec (requires skills-ref)
npx skills-ref validate ./schemas/track.schema.json
```

All workspace records must pass validation before QC can be marked as passed on a release.

---

## Common Field Conventions

**`id`** — Always a slug: lowercase letters, numbers, hyphens only. Must match the record's filename or directory name. Pattern: `^[a-z0-9-]+$`.

**`schema_version`** — Always `"1"` until a breaking change is made. Stored as a string, not a number, to keep comparison unambiguous.

**Nullable fields** — Fields typed as `["string", "null"]` or `["integer", "null"]` are optional data that will be populated over time (e.g. `isrc`, `internet_archive_id`). Set them to `null` until the value is known, rather than omitting them.

**Dates** — All dates use ISO 8601 format (`YYYY-MM-DD`). No datetimes; the label state document uses plain dates throughout.

**Paths** — File paths in schema records are relative to the record's own directory unless otherwise noted.

**Arrays as top-level containers** — `opportunities.json`, `revenue.json`, and `expenses.json` are JSON files whose root element is an array of the respective schema objects, not a wrapper object. This keeps them easy to append to and diff.

---

## Adding a New Schema

1. Create `schemas/[name].schema.json` following the patterns above
2. Set `$id` to `https://openvirtuallabel.org/schemas/[name]/v1`
3. Include `schema_version` as a required `const: "1"` field
4. Add an entry to this README
5. Update `OVL-CONTEXT.md` in the agent skills references if agents need to know about the new record type
