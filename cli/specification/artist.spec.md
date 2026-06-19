# CLI Specification — `ovl artist`

Commands for creating and managing artist profiles
(`workspace/artists/[artist-id]/artist.json`).

All writes validate against `schemas/artist.schema.json` before saving.

---

## Command Reference

---

### `ovl artist create`

Create a new artist profile interactively.

```text
ovl artist create
```

**Behaviour:**

- Prompts for display name, legal name, default license, PRO, IPI number, distributor,
  and platform links
- Generates a slug from the display name (e.g. `brylie-christopher`)
- Creates `workspace/artists/[slug]/artist.json`
- Validates against `schemas/artist.schema.json` before writing

**Output:** Path to the created `artist.json` and the generated artist ID.

---

### `ovl artist list`

List all artist profiles in the workspace.

```text
ovl artist list
```

**Output:** Table of artist IDs, display names, and distributor.

---

### `ovl artist show <artist-id>`

Display an artist profile.

```text
ovl artist show <artist-id>
```

**Output:** Formatted display of all populated `artist.json` fields.

---

### `ovl artist add-alias <artist-id> --name "<alias>"`

Add a performing name alias to an existing artist profile.

```text
ovl artist add-alias <artist-id> --name "<alias>"
```

**Behaviour:** Appends to `artist.also_known_as[]` and saves.

---

### `ovl artist set-bio <artist-id>`

Set or update the bio for an artist.

```text
ovl artist set-bio <artist-id> [--short "<text>"] [--medium "<text>"] [--full "<text>"]
```

**Behaviour:**

- Updates one or more bio length variants in `artist.bio`
- Only the flags provided are updated; omitted flags leave existing values unchanged
- Validates character limits: short ≤ 280, medium ≤ 1024, full unbounded
- If no flags are provided, opens an interactive prompt for each length in sequence

**Options:**

| Option | Description |
|---|---|
| `--short "<text>"` | One or two sentences. For social profiles and streaming bios. Max 280 characters. |
| `--medium "<text>"` | One paragraph. For press kits and Bandcamp. Max 1024 characters. |
| `--full "<text>"` | Full bio. For website and detailed press materials. |

**Output:** Confirmation of which bio lengths were updated.

**Errors:**

- `ARTIST_NOT_FOUND` — no artist with the given ID exists in the workspace
- `VALIDATION_FAILED` — bio text exceeds length limit

---

### `ovl artist set-contact <artist-id>`

Update contact details for an artist.

```text
ovl artist set-contact <artist-id> [--email <email>] [--website <url>]
```

**Behaviour:** Updates one or more fields within `artist.contact`. Only the flags
provided are updated; omitted flags leave existing values unchanged.

**Options:**

| Option | Description |
|---|---|
| `--email <email>` | Public or licensing contact email |
| `--website <url>` | Artist website URL |

**Errors:**

- `VALIDATION_FAILED` — email format invalid or website is not a valid URI

---

### `ovl artist set-platform <artist-id> --platform <platform> --value <value>`

Set a platform identifier or URL on an artist profile.

```text
ovl artist set-platform <artist-id> --platform <platform> --value <value>
```

**Valid platforms:**

| Platform | Field | Value format |
|---|---|---|
| `spotify` | `spotify_artist_id` | Artist ID string |
| `apple-music` | `apple_music_artist_id` | Artist ID string |
| `youtube` | `youtube_channel_id` | Channel ID string (e.g. `UCxxxxxxxx`) |
| `youtube-music` | `youtube_music_artist_id` | Artist ID string |
| `bandcamp` | `bandcamp_url` | Full URL (e.g. `https://artist.bandcamp.com`) |
| `soundcloud` | `soundcloud_url` | Full URL (e.g. `https://soundcloud.com/artist`) |
| `instagram` | `instagram_handle` | Handle without `@` |
| `facebook` | `facebook_url` | Full URL |
| `tiktok` | `tiktok_handle` | Handle without `@` |
| `subvert-fm` | `subvert_fm_url` | Full URL (e.g. `https://www.subvert.fm/artist`) |

**Behaviour:**

- Sets the specified platform field on `artist.platforms`
- Validates URI format for URL fields before saving

**Output:** Confirmation showing the platform and value set.

**Errors:**

- `ARTIST_NOT_FOUND` — no artist with the given ID exists in the workspace
- `UNKNOWN_PLATFORM` — platform name is not in the valid list
- `VALIDATION_FAILED` — value fails format validation (e.g. invalid URI)

---

### `ovl artist set-rights <artist-id>`

Update performing rights registration details for an artist.

```text
ovl artist set-rights <artist-id> [--pro <pro>] [--ipi <number>] [--isni <number>]
```

**Options:**

| Option | Description |
|---|---|
| `--pro <pro>` | Performing rights organisation (e.g. `Teosto`, `ASCAP`, `PRS`) |
| `--ipi <number>` | IPI number assigned by the PRO |
| `--isni <isni>` | International Standard Name Identifier |

**Behaviour:** Updates one or more fields within `artist.rights`. Only the flags provided
are updated.

---

### `ovl artist set-distributor <artist-id>`

Set distribution details for an artist.

```text
ovl artist set-distributor <artist-id> --distributor <distributor> [--id <distributor-artist-id>]
```

**Options:**

| Option | Description |
|---|---|
| `--distributor <name>` | Distribution platform (e.g. `amuse`, `distrokid`, `cdbaby`) |
| `--id <id>` | Artist identifier within the distributor's system |

---

### `ovl artist set-license <artist-id> "<license>"`

Set the default license for an artist's releases.

```text
ovl artist set-license <artist-id> "<license>"
```

**Behaviour:** Updates `artist.default_license`. Overrides the label default for all
this artist's releases unless the release itself specifies a license.

---

### `ovl artist set-location <artist-id> "<location>"`

Set the artist's location.

```text
ovl artist set-location <artist-id> "<location>"
```

**Behaviour:** Updates `artist.location`. Used for press kits and venue outreach.
Recommended format: `City, Country` (e.g. `Helsinki, Finland`).
