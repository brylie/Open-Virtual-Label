# Onboarding

First-time setup of an OVL label workspace. Run once when a new label operator adopts OVL. By the end of this workflow the workspace is populated, the first artist is registered, and the system is ready to manage its first release.

Estimated time: 30–60 minutes depending on how much information is ready to hand.

---

## Prerequisites

- OVL repository cloned or downloaded
- Node.js installed (v18 or higher)
- A text editor available for reviewing generated files
- Artist's platform URLs, PRO details, and distributor information to hand (look these up before starting if needed)

---

## Step 1: Install

```bash
npm install -g ovl
```

Verify:

```bash
ovl --version
```

---

## Step 2: Initialise the Workspace

```bash
ovl init
```

The CLI scaffolds the `workspace/` directory from `workspace-scaffold/`, creates an initial `workspace/state/label-state.md`, and opens a short interactive setup:

```
Welcome to Open Virtual Label.

This will set up your label workspace. You can change any of these
settings later by editing the JSON files in workspace/label/.

Label name: _
```

Provide:

- Label name (can be the artist's own name if this is a solo setup)
- Label contact email
- Default license (default: CC BY 4.0)
- Primary distributor (e.g. amuse, distrokid, cdbaby)

This populates `workspace/label/profile.json`.

→ `label-state.md` created with initial setup entry

---

## Step 3: Create the First Artist Profile

```bash
ovl artist create
```

The CLI walks through the artist profile interactively:

```
Artist display name (performing name): _
Legal name (not published — for contracts and rights registration): _
Default license for this artist [CC BY 4.0]: _
Performing rights organisation (PRO) [leave blank if not registered]: _
IPI number [leave blank if unknown]: _
Primary distributor [inherits from label: amuse]: _
```

Then platform links (all optional — skip any that do not yet exist):

```
Spotify artist URL or ID [optional]: _
YouTube channel URL or ID [optional]: _
Bandcamp URL [optional]: _
Instagram handle [optional]: _
```

This creates `workspace/artists/[artist-id]/artist.json`.

If the artist uses more than one performing name:

```bash
ovl artist add-alias --artist <artist-id> --name "<alias>"
```

---

## Step 4: Create the First Mastering Profile

At least one mastering profile is needed before beginning work on a release.

```bash
ovl mastering profile create
```

The CLI asks:

```
Profile name: _
  (e.g. "Ambient / Streaming", "Electronic / Club", "Acoustic / Podcast Bed")

Target integrated LUFS range:
  Minimum [−18]: _
  Maximum [−16]: _

True peak ceiling in dBTP [−1.0]: _

Minimum loudness range (LRA) in LU [8]: _

Sample rate [44100]: _
Bit depth [24]: _
```

Then platform notes — for each common platform, the CLI provides the standard target and asks if you want to include a custom note:

```
Spotify normalises to −14 LUFS. Add a note? [y/n]: _
YouTube normalises to −14 LUFS. Add a note? [y/n]: _
Apple Music Sound Check targets −16 LUFS. Add a note? [y/n]: _
```

This creates `workspace/artists/[artist-id]/mastering-profiles/[profile-id].json`.

For ambient and contemplative music, the defaults (−18 to −16 LUFS, −1.0 dBTP, LRA ≥ 8) are a good starting point.

---

## Step 5: Verify Setup

```bash
ovl status
```

→ `orchestrator`

The orchestrator reads the new state document and confirms what has been set up:

```
Label: [Label Name]
Artist: [Artist Name]
Mastering profiles: [profile name]

No active releases yet.

Workspace is ready. Next steps:
· Create your first release: ovl release create "<title>"
· Or explore available agents: ovl agents list
```

```bash
ovl validate --all
```

Validates all workspace records against their schemas. Should report no errors on a fresh workspace.

---

## Step 6 (Optional): Configure MCPs

MCPs extend agent capabilities with external service access. None are required to begin — the system works without them, but they enable operations like email outreach and calendar scheduling to be handled within sessions rather than manually.

```bash
ovl mcp list
```

Lists available MCPs and their connection status. Configure as needed:

```bash
ovl mcp connect gmail
ovl mcp connect google-calendar
ovl mcp connect internet-archive
```

Each MCP walks through its own authentication flow. MCP credentials are stored locally and never written to the workspace JSON records or the OVL repository.

---

## Step 7: Create the First Release

When ready to begin production:

```bash
ovl release create "<title>" --artist <artist-id>
```

Then follow the **Release Pipeline** workflow (`workflows/release-pipeline.md`).

---

## Multi-Artist Setup

If the label has more than one artist, repeat Steps 3–4 for each artist. Artist profiles are independent — each has their own directory, their own mastering profiles, and their own platform IDs.

```bash
ovl artist create   # repeat for each artist
```

The label coordinator (orchestrator) tracks all artists' releases in a single `label-state.md`. The Active Projects section lists releases across all artists.

```bash
ovl status --artist <artist-id>   # filter status to one artist
```

---

## Forking OVL for a New Label

If adopting OVL from an existing installation rather than installing fresh:

1. Clone or fork the repository
2. Confirm `workspace/` is in `.gitignore` — **do not commit workspace data**
3. Run `ovl init` — this creates a fresh workspace for this label instance
4. The upstream repo's agents, schemas, and workflows are available immediately
5. Customisations to agent skill files or schemas should be maintained as local overrides, not upstream commits, unless they are genuinely generic improvements worth contributing back

---

## Troubleshooting First Setup

**`ovl init` fails with a workspace already exists error:**
A `workspace/` directory already exists. Run `ovl init --force` to overwrite, or inspect the existing workspace first.

**Artist ID already taken:**
Artist IDs must be unique within the workspace. If `brylie-christopher` already exists, use a more specific slug such as `brylie-christopher-main`.

**Schema validation fails on first validate:**
Open the flagged file and compare against the corresponding schema in `schemas/`. The error message identifies the field and the constraint that failed.

**Mastering profile not appearing in release setup:**
Confirm the profile JSON is in `workspace/artists/[artist-id]/mastering-profiles/` and passes `ovl validate`.
