# Playlist Pitch

Researching and submitting tracks to independent playlist curators. Playlist placement is one of the highest-leverage activities for ambient and instrumental music discovery — a well-placed track can compound streams for months without ongoing effort.

This workflow targets **independent curators**, not algorithmic playlists (Spotify's Release Radar, Discover Weekly, etc.) which cannot be pitched directly. Editorial playlists (Spotify editorial, Apple Music editorial) have their own submission tools covered separately.

---

## Prerequisites

- At least one track with `status: live` (must be streaming before pitching)
- ISRC for the track being pitched
- Spotify URI or link for the track being pitched
- `workspace/outreach/opportunities.json` exists

---

## When to Run This Workflow

Playlist pitching is most effective:

- Within the first 4 weeks of a track going live
- When a track has an identifiable mood or use case (study music, sleep music, meditation, focus)
- After accumulating at least a few hundred streams organically (shows some listener validation)

It is least effective:

- Before a track is live on streaming platforms
- For tracks with no clear playlist home (avoid pitching every track to every playlist)

---

## Stage 1: Identify Target Playlists

```bash
ovl outreach research --type playlist-pitch --track <track-id>
```

→ `outreach-crm`

The CRM agent researches playlists suited to the track. For each candidate it checks:

- Playlist size (follower count) — aim for a mix of sizes; small engaged playlists often convert better than large passive ones
- Curator activity (when was the last track added? Is it still being updated?)
- Aesthetic fit (does the existing track list match the style and energy of the pitch track?)
- Submission policy (do they accept pitches? Is there a preferred method — SubmitHub, email, social DM?)

The agent creates `opportunity.json` records with `type: playlist-pitch` for each viable playlist.

**[REVIEW GATE]** The artist reviews the candidate list and confirms which to pursue.

---

## Stage 2: Prepare the Pitch Asset

Before drafting any outreach, confirm the pitch package is ready:

- Track is live on Spotify with its correct metadata
- Spotify for Artists pre-save or pitch link available
- Short track description exists (2–3 sentences on the mood, inspiration, or use case)
- `track.musical_key` and approximate BPM are populated if known (some curators filter by these)

If the track was released through a distributor with a Spotify editorial pitch tool (Amuse, DistroKid, etc.), submit to Spotify editorial first — that submission can be made up to 7 days before release, which is earlier than independent curator outreach.

---

## Stage 3: Draft Pitches

Playlist pitches are brief. The CRM agent drafts personalised pitches for each curator — each one references something specific about the playlist to show genuine familiarity.

```bash
ovl outreach draft --opportunity <opportunity-id>
```

→ `outreach-crm`

A good playlist pitch:

- Opens with the specific playlist name and why this track fits it
- Includes the Spotify link, track duration, and genre tags
- Mentions any relevant context (release date, CC licensing if applicable)
- Is under 150 words
- Has a clear, low-pressure close ("happy to share more" or "no problem if it's not a fit")

A poor playlist pitch:

- Opens with artist biography
- Uses generic flattery ("love your playlist!")
- Makes multiple asks or pitches multiple tracks at once
- Is longer than 200 words

**[APPROVAL GATE]** Artist reviews all drafts before any are marked as ready to send. Multiple pitches can be reviewed in one session.

---

## Stage 4: Submit

Submission method varies by curator:

**SubmitHub / Groover:** Submit through the platform's interface. The platform handles the delivery and any associated cost. Log the submission:

```bash
ovl outreach log --opportunity <opportunity-id> \
  --action sent \
  --note "Submitted via SubmitHub, 2 credits used"
```

**Direct email:** Approved pitch is sent via email MCP or manually. Confirm send before logging.

**Social DM:** Approved pitch adapted for direct message. The CRM agent may suggest adjustments for the character limit or conversational tone. Log when sent.

`opportunity.status → sent`, `follow_up_due` set to 14 days.

---

## Stage 5: Track Responses

Most curators do not reply. That is normal. A response rate of 10–20% is typical for cold playlist outreach; placement rate from responses varies widely.

When a response arrives:

```bash
ovl outreach log-response --opportunity <opportunity-id>
```

**Added to playlist:** `opportunity.status → won`. Log which playlist and date added. Monitor for stream uplift in next monthly metrics snapshot. Consider a thank-you reply and offer to stay in touch for future releases.

**Declined:** `opportunity.status → declined`. Note any reason given. If the curator gave specific feedback (wrong energy, not accepting new music), note it for future targeting.

**No response after follow-up:** `opportunity.status → stale`. The playlist may still add the track without notification — check the track's playlist reach metric in the next monthly snapshot for unexpected additions.

---

## Stage 6: Monitor Results

Playlist placement impact appears in the Spotify metrics:

- `playlist_reach` (number of listeners who heard the track through a playlist)
- Monthly listener uplift correlated with placement date

At the monthly review, the Metrics Analyst checks:

```bash
ovl metrics snapshot --period <YYYY-MM>
```

If a track shows a stream spike without a corresponding social push, it may indicate a playlist addition the artist was not notified of. Check `spotify.playlist_reach` against the prior month.

---

## Spotify Editorial Pitch (Separate Process)

Editorial pitches through Spotify for Artists (pitching to Spotify's own curators) are separate from independent curator outreach and must be submitted before the track goes live:

1. Upload release to distributor at least 7 days before release date
2. Log into Spotify for Artists
3. Navigate to the unreleased track → Pitch to playlists
4. Complete the pitch form: mood, genre, tempo, culture, instruments, and description

This is a manual step — no OVL agent handles it. Record the submission:

```bash
ovl outreach log --type playlist-pitch \
  --note "Spotify editorial pitch submitted [date] for [track]"
```

---

## Pitch Volume and Cadence

Recommended approach per release:

- 5–10 independent curator pitches per new track
- Focus on quality of fit over volume
- Space pitches out over the first 4 weeks rather than sending all at once

Pitching too broadly too quickly can signal low selectivity to curators. Pitching the same track to 50 playlists in one day will not outperform 10 well-chosen, genuinely personalised pitches.

Between releases, maintain curator relationships:

- Reply to curators who added previous tracks
- Mention them in social posts when a track is placed
- When the next release is ready, warm pitches to previous contacts convert significantly better than cold outreach
