# CLI Specification — `ovl mastering`

Commands for managing mastering sessions and mastering profiles.
Mastering sessions invoke the `mastering-companion` agent and write results to
track records. Profiles define loudness and format targets applied across a release.

---

## Command Reference

---

### `ovl mastering start --track <track-id>`

Launch a mastering session with the Mastering Companion agent.

```text
ovl mastering start --track <track-id> [--remaster]
```

**Behaviour:**

- Loads the track record and the applicable mastering profile
- Invokes the `mastering-companion` agent in guided-session mode
- On session completion, writes mastering measurements to `track.mastering{}`
- Appends session notes to `mastering_profile.session_notes[]`

**Options:**

| Option | Description |
|---|---|
| `--remaster` | Re-master a track that already has mastering data. Previous measurements are preserved in `track.mastering.notes` before being overwritten. |

---

### `ovl mastering profile create`

Create a new mastering profile interactively.

```text
ovl mastering profile create [--artist <artist-id>]
```

**Behaviour:**

- Prompts for profile name, LUFS targets, true peak ceiling, LRA guidance, sample rate,
  bit depth
- Offers platform-specific note prompts for Spotify, Apple Music, YouTube
- Generates a slug from the name
- Creates `workspace/artists/[artist-id]/mastering-profiles/[slug].json`
- Validates against `schemas/mastering-profile.schema.json` before writing

---

### `ovl mastering profile list`

List mastering profiles for an artist.

```text
ovl mastering profile list [--artist <artist-id>]
```

**Output:** Table of profile IDs, names, LUFS targets, and session count.
