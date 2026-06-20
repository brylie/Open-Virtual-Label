# Platform Specifications

Character limits, format requirements, and output structure for each
platform. These are technical constraints — voice and tone come from
`workspace/artists/<id>/social-media.md`.

---

## Instagram

### Caption

- **Target length:** 150–300 words
- **Maximum:** 2,200 characters
- **Preview (before "more"):** first ~125 characters — make the opening
  count without resorting to clickbait
- **Line breaks:** blank lines between paragraphs; Instagram collapses
  multiple line breaks, so one blank line is enough
- **Emojis:** only if the artist's voice config allows or requests them

### Hashtags

- Include below the caption, separated by a blank line
- 5–10 tags; more feels spammy
- Mix: 2–3 core genre tags + 1–2 niche/situational tags + 1 location tag
  (if relevant)
- Load the artist's approved hashtag list from `social-media.md`
- Do not use hashtags in YouTube titles or descriptions

### Caption format

```
[Opening line — specific and direct; no "excited to announce"]

[1–2 paragraphs of context, story, or description]

[Optional: call to action or natural close]

[Hashtags on a new paragraph]
```

### Output format

```
INSTAGRAM CAPTION:
[caption text — no markdown]

(X characters)

HASHTAGS:
#tag1 #tag2 #tag3
```

---

## YouTube

### Title

- **Maximum:** 100 characters; but aim for **60 or under** — titles over
  60 are truncated in most YouTube surfaces
- **Always show count:** `(X/60)` in your output
- Keyword-friendly: think how someone would search for this content
- Descriptive, not clickbait
- No ALL CAPS, no `!!!`

**Title patterns for musicians:**
```
[Track/Album Name] — [Type] | [Artist Name]
[Artist Name] — [Album/Track Name] ([Year])
[Description] — [Artist Name] Live at [Venue]
```

### Description

- **Preview section:** first 150 characters shown without expanding —
  front-load the most important line
- **Structure:**
  1. What it is (one sentence, 150 chars or under for preview)
  2. Context or background (1–3 short paragraphs)
  3. Streaming/listening links
  4. License information
  5. Subscribe / follow CTA
- No character limit, but front-load what matters
- Links work in descriptions (unlike captions)

### Tags

- 5–10 comma-separated keywords
- Mix broad (`ambient music`) and specific (`piano improvisation`,
  `creative commons music`)
- Do not repeat the title verbatim — add related terms

### Output format

```
YOUTUBE TITLE:
[Title] (X/60)

YOUTUBE DESCRIPTION:
[First 150 chars — the preview]

[Additional paragraphs]

Listen / Stream:
[Platform]: [link]
[Platform]: [link]

[License info]

[Subscribe CTA]

TAGS:
ambient music, piano music, [additional tags]
```

---

## Facebook

### Posts

- More flexible length than Instagram
- Links display as rich previews — including a link is usually enough;
  the caption can be short
- Works well for event announcements with practical details (date, time,
  address)
- Cross-posting from Instagram is fine; slightly shorter often works better
  on Facebook
- Hashtag usage is minimal compared to Instagram — 1–3 if any

### Output format

```
FACEBOOK:
[post text]
```

---

## Multiple variations

When offering options, label them clearly:

```
OPTION A — [Brief angle description]:
[copy]

OPTION B — [Brief angle description]:
[copy]

OPTION C — [Brief angle description]:
[copy]

---
Which fits best?
```
