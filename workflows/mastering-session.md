# Mastering Session

A single-track mastering session with the Mastering Companion. This is a guided-session workflow: the artist operates their own tools while the agent interprets readings and advises on next steps in real time.

This workflow is called once per track. It is invoked as part of the release pipeline (Stage 2) but can also be run independently for re-masters or single releases.

---

## Prerequisites

- `track.json` exists for the target track with `release_id` set
- `release.mastering_profile_id` is set, or the track has its own `mastering.profile_id`
- The mastering profile exists in `workspace/artists/[id]/mastering-profiles/[profile-id].json`
- A lossless pre-master export is ready (24-bit WAV, exported from the DAW before any master bus limiting or loudness normalisation)
- The artist's metering tools are open and ready (spectrum analyser, LUFS meter, true peak limiter)

---

## Starting the Session

```bash
ovl mastering start --track <track-id>
```

→ `mastering-companion`

The Mastering Companion loads:

- `track.json` — title, collaborators, any existing mastering notes
- The applicable `mastering_profile.json` — targets, platform notes, step checklist, previous session notes

The agent opens by confirming the session context:

```text
Track: "Chromatic Drift" (track 3, Spectra)
Profile: Ambient / Streaming (Dynamic Preservation)
Targets: −18 to −16 LUFS integrated, −1.0 dBTP true peak, LRA ≥ 8 LU
Tools expected: LUFS meter, true peak limiter, reference comparison plugin

Ready to begin. First: do you have a lossless pre-master export ready,
before any limiting or loudness normalisation on the master bus?
```

---

## Session Steps

The Mastering Companion walks through the profile's checklist. The artist confirms each step as they complete it. Steps are not skipped — if a step cannot be completed, the reason is noted.

### Step 1: Pre-master check

The agent asks the artist to play through the full pre-master and report:

- Any obvious issues (clipping in source, unwanted noise, phase problems)
- Approximate level before the master chain

If issues are found, the agent advises whether to return to mixing or whether they can be addressed in mastering.

### Step 2: Reference comparison

Load a reference track in the same genre and approximate target loudness. The agent asks:

- How does the spectral balance compare (low end, midrange, high end)?
- Is the pre-master brighter, darker, or similar in character?

The artist describes what they hear. The agent advises on EQ approach if needed.

### Step 3: Mastering chain application

The artist applies their mastering chain (EQ, compression, limiting). The Mastering Companion does not prescribe specific plugin settings — it advises based on what the artist reports hearing and measuring. Typical conversational exchange:

```text
Artist: I've applied a high-shelf boost of about 1.5dB at 8kHz and a
        gentle low-shelf cut at 80Hz. Limiting is set to ceiling −1.0dBTP.

Agent:  Good start. What is your integrated LUFS reading at this point?

Artist: −17.4 LUFS

Agent:  That's within the target range (−18 to −16). True peak?

Artist: −1.1 dBTP

Agent:  Slightly below ceiling — you have a touch of headroom left if you
        want to bring up the limiting by 0.5dB, or you can leave it here.
        What is the LRA reading?

Artist: 10.2 LU

Agent:  Well above the minimum of 8 LU. The dynamics are well preserved.
        This is looking good for the profile. How does it compare to the
        reference track at matched loudness?
```

### Step 4: Low-frequency mono check

The agent prompts a mid/side check on the low end:

```text
Agent:  Before we finalise — can you check what's happening below 100Hz
        in the side channel? Solo the Side signal and look at the spectrum
        below 100Hz. Any significant energy down there?
```

If significant sub-100Hz side content exists, the agent advises making the low end mono at the master bus (M/S processing, low shelf on Side channel, or a dedicated mono maker).

### Step 5: Final measurements

Artist plays the full track through the finalised master chain and reports:

- Integrated LUFS
- True peak (dBTP)
- LRA

The agent checks all three against profile targets and confirms pass or flags any deviation.

**[APPROVAL GATE]** If any measurement falls outside the target range, the agent presents the deviation and asks how to proceed: adjust the chain, accept with a rationale note, or return to mixing.

### Step 6: Export

```text
Agent:  Measurements confirmed within target. Export the master as:
        · Format: WAV
        · Bit depth: 24-bit
        · Sample rate: 44100 Hz
        · No dithering applied (dithering only if reducing to 16-bit)

        Save as: [suggested filename based on track title and release]
        Update the file path in the track record once exported.
```

The artist exports and updates the `files.master_wav` path:

```bash
ovl track set-file --track <track-id> --field master_wav --path <path>
```

### Step 7: Session notes

The agent asks one closing question:

```text
        particularly well, any unusual characteristic of this recording,
        or anything you'd adjust on the next track?
```

The artist's response is appended to `mastering_profile.session_notes[]`.

---

## Session Output

On completion, the agent writes to `track.mastering{}`:

```json
{
  "profile_id": "<profile-id>",
  "integrated_lufs": <value>,
  "true_peak_dbtp": <value>,
  "lra": <value>,
  "sample_rate_hz": 44100,
  "bit_depth": 24,
  "mastered_date": "<YYYY-MM-DD>",
  "mastered_by": "<artist display name>",
  "notes": "<any session notes>"
}
```

And appends to `mastering_profile.session_notes[]`:

```json
{
  "date": "<YYYY-MM-DD>",
  "track_id": "<track-id>",
  "note": "<artist's closing note>"
}
```

→ `label-state.md` updated: track mastering progress noted

---

## When Measurements Fall Outside Target

The profile's target range is a guide, not a hard requirement. Ambient and dynamic music may legitimately fall outside standard ranges. When this happens:

- The agent explains which measurement is out of range and by how much
- The agent explains the practical consequence (e.g. "at −19.5 LUFS, this track will be turned up by Spotify's normalisation, which may introduce audible limiting artefacts")
- The artist decides whether to adjust or accept with a rationale

If accepted outside target, the deviation and rationale are written to `track.mastering.notes` and `track.qc.override` is set.

---

## Re-mastering an Existing Track

If a track already has mastering data and the session is a re-master:

```bash
ovl mastering start --track <track-id> --remaster
```

The previous mastering data is preserved in `track.mastering.notes` as a dated entry before being overwritten. The Mastering Companion surfaces the previous measurements at session start for comparison.
