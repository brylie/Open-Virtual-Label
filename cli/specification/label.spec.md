# CLI Specification — `ovl label`

Commands for viewing and updating the label profile (`workspace/label/profile.json`).
The label record is created by `ovl init` and is the only record that exists before any
artist is added. These commands allow the profile to be inspected and updated without
editing JSON directly.

All writes validate against `schemas/label.schema.json` before saving.

---

## Command Reference

---

### `ovl label show`

Display the current label profile.

```text
ovl label show [--json]
```

**Behaviour:**

- Reads `workspace/label/profile.json`
- Outputs a formatted summary of all populated fields
- Omits fields that are not set rather than printing nulls

**Output example:**

```text
Label:        OurMusi.cc
ID:           ourmusicc
Description:  Independent Creative Commons label for ambient, piano, and experimental music.
License:      CC BY 4.0
Distributor:  amuse
Contact:      brylie.oxley@gmail.com
Website:      https://ourmusi.cc
Created:      2026-06-19
```

**Errors:**

- `WORKSPACE_NOT_FOUND` — no workspace found; run `ovl init` first

---

### `ovl label set-name <name>`

Rename the label.

```text
ovl label set-name "<name>"
```

**Behaviour:**

- Updates `label.name` in `label/profile.json`
- Does **not** change `label.id` — the ID is set at init and is not renamed (it matches
  the directory name convention and may be referenced in other records)
- Validates and saves

**Output:** Confirmation line showing the old and new name.

**Errors:**

- `VALIDATION_FAILED` — name exceeds 128 characters or is empty

---

### `ovl label set-description "<text>"`

Set the label description.

```text
ovl label set-description "<text>"
```

**Behaviour:** Updates `label.description`. Maximum 1024 characters.

---

### `ovl label set-license "<license>"`

Set the label-wide default license.

```text
ovl label set-license "<license>"
```

**Behaviour:**

- Updates `label.default_license`
- This is the fallback used when an artist or release does not specify a license
- Common values: `CC BY 4.0`, `CC BY-SA 4.0`, `CC0 1.0`, `All Rights Reserved`

**Output:** Confirmation showing the old and new license.

---

### `ovl label set-distributor <distributor>`

Set the label-wide default distributor.

```text
ovl label set-distributor <distributor>
```

**Behaviour:** Updates `label.default_distributor`. Used as the default when creating
new releases unless overridden at the artist or release level.

---

### `ovl label set-contact`

Update label contact details.

```text
ovl label set-contact [--email <email>] [--website <url>] [--location <location>]
```

**Behaviour:**

- Updates one or more fields within `label.contact`
- Only the flags provided are updated; omitted flags leave existing values unchanged
- Validates email format and URI format before saving

**Options:**

| Option | Description |
|---|---|
| `--email <email>` | Public contact email for the label |
| `--website <url>` | Label website URL |
| `--location <location>` | City and country (e.g. `Helsinki, Finland`) |

**Errors:**

- `VALIDATION_FAILED` — email format invalid or website is not a valid URI

---

## Exit Codes

Inherits the standard OVL exit codes. In addition:

| Code | Meaning |
|---|---|
| `2` | Workspace not found |
| `3` | Schema validation failure |
