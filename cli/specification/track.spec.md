# CLI Specification — `ovl track` and `ovl isrc`

Commands for managing track records within a release
(`workspace/artists/[artist-id]/releases/[release-id]/tracks/[track-id].json`)
and assigning ISRCs.

All writes validate against `schemas/track.schema.json` before saving.

---

## Command Reference

---

### `ovl track add "<title>" --release <release-id>`

Add a track to a release.

```text
ovl track add "<title>" --release <release-id> [--position <n>]
```

**Behaviour:**

- Generates a slug from the title
- Creates `workspace/artists/[artist-id]/releases/[release-id]/tracks/[slug].json`
- If `--position` is omitted, appends after the last existing track

**Options:**

| Option | Description |
|---|---|
| `--position <n>` | Track number (1-based). If omitted, appends at end. |

---

### `ovl track show <track-id>`

Display a track record.

```text
ovl track show <track-id> [--release <release-id>]
```

**Behaviour:** Track IDs are unique within a release. If the same slug exists in
multiple releases, `--release` is required to disambiguate.

---

### `ovl track set-file <track-id> --field <field> --path <path>`

Set a file path on a track record.

```text
ovl track set-file <track-id> --field <field> --path <path>
```

**Valid fields:** `master_wav`, `stems_zip`, `project_file`, `mp3_320`,
`wav_for_distribution`

**Behaviour:** Validates that the referenced path exists before writing. Path is stored
relative to the release directory.

---

### `ovl isrc assign --release <release-id>`

Assign ISRCs to tracks that do not yet have one.

```text
ovl isrc assign --release <release-id> [--track <track-id>]
```

**Behaviour:**

- Lists all tracks in the release with no ISRC assigned
- For each, prompts for the ISRC (format: `CC-XXX-YY-NNNNN`)
- Validates format against the schema pattern before saving
- If `--track` is specified, assigns only to that track

**Note:** OVL does not register ISRCs. They are issued by the artist's distributor or
national ISRC agency. This command records codes that have already been obtained.
