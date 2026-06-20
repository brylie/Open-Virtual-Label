# `ovl site` — website content sync

Manages and syncs websites whose content collections mirror the OVL workspace.
The workspace JSON files are the **canonical source of truth**. Websites are secondary
consumers, populated on demand by running `ovl site sync`.

This model supports any number of sites simultaneously — a label site showing all artists,
a personal artist site showing one artist's releases, or both:

```
OVL workspace (canonical)
  └─ ovl site sync
       ├─ sites/my-label.com  (label site — all artists)
       └─ sites/aria-nova     (artist site — aria-nova only)
```

---

## Commands

### `ovl site add <id> <path>`

Register a website target on the label profile.

| Argument | Description |
| --- | --- |
| `<id>` | Slug used to target this site with `--site`. e.g. `label-site`, `aria-nova-site`. |
| `<path>` | Path to the site root. Relative paths resolve from the workspace directory. |

**Flags**

| Flag | Description |
| --- | --- |
| `--artist <id>` | Restrict sync to a single artist's records. Omit for a label-wide site. |
| `--artists-dir <path>` | Override destination for artist JSON within the site root. Defaults to `src/content/artists`. |
| `--releases-dir <path>` | Override destination for release JSON within the site root. Defaults to `src/content/releases`. |
| `--description <text>` | Human-readable label shown in `ovl site list`. |

Use `--artists-dir` / `--releases-dir` when the site's existing content collections already occupy the default paths. This is common for artist sites that maintain hand-authored Markdown/MDX releases alongside the OVL-synced catalog.

**Examples**

```sh
# Label-wide site (all artists and releases, default dirs)
ovl site add label-site sites/my-label.com --description "Label site"

# Artist site — separate catalog dirs to avoid conflicting with existing MDX releases
ovl site add aria-nova-site sites/aria-nova \
  --artist aria-nova \
  --artists-dir src/content/catalog/artists \
  --releases-dir src/content/catalog/releases \
  --description "Aria Nova artist site"
```

---

### `ovl site list`

List all registered website targets and their configuration.

```
ID                    PATH                  ARTIST       DESCRIPTION
label-site            sites/my-label.com    (all)        Label site
aria-nova-site        sites/aria-nova       aria-nova    Aria Nova artist site
```

---

### `ovl site remove <id>`

Unregister a website target. Does not delete any files on disk.

---

### `ovl site sync [--site <id>]`

Copy artist and release JSON records from the workspace into each registered
site's Astro content collection directories.

Without `--site`, all registered sites are synced. With `--site <id>`, only
that site is updated.

**What is written**

| Source | Destination |
| --- | --- |
| `workspace/artists/{id}/artist.json` | `{site}/src/content/artists/{id}.json` |
| `workspace/artists/{id}/releases/{rel}/release.json` | `{site}/src/content/releases/{id}--{rel}.json` |

For sites registered with `--artist`, only that artist's records are synced.
For label-wide sites (no `--artist`), all artists and their releases are synced.

Destination directories are created if they do not exist. Existing files are
overwritten.

**Flags**

| Flag | Description |
| --- | --- |
| `--site <id>` | Sync only this site ID (default: all registered sites). |

**Example workflow**

```sh
# After any workspace change:
ovl release add-link my-release --platform bandcamp --url https://...
ovl site sync                     # push to all sites
ovl site sync --site aria-nova-site  # or push to one site only
```

**Exit codes**

| Code | Meaning |
| --- | --- |
| 0 | All targeted sites synced successfully. |
| 2 | No sites configured, site directory missing, or filtered artist not found. |
| 1 | Unexpected I/O error during file copy. |

---

## Path resolution

`path` in each site registration is resolved as follows:

- **Relative path** (no leading `/`): resolved from the workspace directory.
  `sites/my-label.com` → `{workspace}/sites/my-label.com`
- **Absolute path**: used as-is.

Storing relative paths is recommended — they remain valid if the workspace
directory is moved or cloned to a different machine.

---

## Astro content.config.ts

Each site requires a `src/content.config.ts` that defines Zod schemas matching
the OVL workspace record shapes. Two reference configs are maintained in this
repo:

| Site | Config |
| --- | --- |
| `workspace/sites/<label-site>` | Label site — syncs all artists and releases to `src/content/artists` and `src/content/releases` |
| `workspace/sites/<artist-site>` | Artist site — syncs one artist's records to `src/content/catalog/artists` and `src/content/catalog/releases`, preserving the existing MDX `releases` collection |

Both configs are identical in schema shape. The scoping happens at sync time
(via `--artist` on `ovl site add`), not in the Astro config. This keeps the
content schema consistent across site types and means a site can be re-purposed
between artist and label use without changing its config.

**Keeping configs in sync with OVL schemas**

When an OVL schema changes (new field added, enum value added), update the
corresponding Zod schema in each site's `content.config.ts`. The comment block
at the top of each config names the source schema files:

```ts
// Mirrors workspace/artists/*/releases/*/release.json
// Status values must match the OVL release schema exactly — OVL is canonical.
```

**Why schema generation is out of scope**

Auto-generating Zod schemas from OVL's JSON Schema definitions is intentionally
not provided. The reasons:

1. JSON Schema format validators (`email`, `uri`, `date`) have no direct Zod
   equivalent and require hand-written refinements.
2. Astro's content loader requires `glob()` configuration that a generator
   cannot infer from the schema alone.
3. The schema surface that a website actually needs is typically smaller than
   the full OVL schema — sites often omit internal fields like `qc`, `archive`,
   and `mastering_profile_id`.
4. A generator would add a build-time dependency between the CLI and
   TypeScript/Astro specifics, coupling two otherwise independent systems.

The maintained `content.config.ts` files in this repo serve as the reference.
Copy one as a starting point for a new site and trim to taste.

---

## Adding a new artist site

1. Create the Astro project at the desired path within the workspace (e.g. `sites/new-artist`).
2. Decide on content directories. If the site has no existing `releases` collection, the defaults (`src/content/artists`, `src/content/releases`) work fine. If it does, pick separate paths (e.g. `src/content/catalog/artists`, `src/content/catalog/releases`) and add matching collection definitions to `content.config.ts`.
3. Register it with OVL:
   ```sh
   # Fresh site — use defaults
   ovl site add new-artist sites/new-artist \
     --artist new-artist-id \
     --description "New Artist site"

   # Site with existing releases collection — use separate dirs
   ovl site add new-artist sites/new-artist \
     --artist new-artist-id \
     --artists-dir src/content/catalog/artists \
     --releases-dir src/content/catalog/releases \
     --description "New Artist site"
   ```
4. Run `ovl site sync --site new-artist` to populate initial content.

The site will receive only that artist's records on every subsequent sync.
