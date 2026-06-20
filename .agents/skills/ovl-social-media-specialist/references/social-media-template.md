# Social Media Config

Artist-specific voice, hashtag strategy, and platform preferences for the
`ovl-social-media-specialist` skill. Copy this file to
`workspace/artists/<artist-id>/social-media.md` and fill in the values.

---

## Voice description

How the artist's copy should feel. Be specific — "authentic" is not enough.

```yaml
voice: >
  [e.g. "Direct and honest. Conversational but not casual to the point of
  carelessness. Describes what something is rather than how great it is.
  Specific details over abstract claims. Humble — invites rather than sells."]
```

## Banned phrases

Words and phrases to never use. Check every draft against this list.

```yaml
banned_phrases:
  # LLM clichés — the most common offenders
  - journey
  - landscape
  - tapestry
  - symphony
  - testament
  - nexus
  - delve
  - dive
  - embark
  - navigating
  - realm
  - bustling
  # Hollow opener
  - "excited to announce"
  - "I'm passionate about"
  # Marketing language
  - "Join me on this"
  - "Transform your"
  - "Elevate your"
  - "Explore the world of"
  # Punctuation
  - em-dash  # use period, comma, or semicolon instead
  # Overhyped
  - amazing
  - incredible
  - mind-blowing
  - life-changing
  - must-see
  - can't miss
```

## Platform priorities

Which platforms to prioritise and how actively to use each:

```yaml
platforms:
  instagram:
    active: true
    feed_posts_per_month: 8
    stories_per_week: 3
    primary_use: community, BTS, releases
  youtube:
    active: true
    uploads_per_month: 2
    primary_use: full tracks, albums, process videos
  facebook:
    active: true
    posts_per_month: 4
    primary_use: events, cross-posts
```

## Hashtag strategy

Core tags (use in almost every post) and situational tags:

```yaml
hashtags:
  core:
    - "#[genre]music"
    - "#[instrument]"
    - "#[mood]music"
    - "#creativecommons"
    - "#independentmusic"
  situational:
    release:
      - "#newmusic"
      - "#newalbum"
    performance:
      - "#livemusic"
      - "#[location]"
    process:
      - "#musicproduction"
      - "#[daw or tool]"
    genre_specific:
      - "#[subgenre]"
      - "#[influence style]"
  max_per_post: 10
  notes: >
    [e.g. "Location tag only for local content. Avoid generic tags like #music."]
```

## CC licensing framing

Approved ways to mention Creative Commons in copy. Use one of these
phrasings; do not invent new ones.

```yaml
cc_framing:
  short: "[e.g. CC BY 4.0 — free to use with attribution.]"
  medium: "[e.g. Released under Creative Commons — free for your podcasts, videos, and games.]"
  long: "[e.g. All my music is released under Creative Commons (CC BY 4.0), which means you can use it in your projects as long as you credit me. No fees, no complicated contracts.]"
```

## Recurring contexts

Things that appear regularly in this artist's content — fill these in
so the skill uses consistent phrasing:

```yaml
recurring_contexts:
  home_city: "[city, country]"
  home_studio_phrase: "[e.g. 'my home studio in [city]']"
  regular_venues:
    - name: "[Venue Name]"
      description: "[brief description — one clause]"
  collaborating_projects:
    - name: "[Project or ensemble name]"
      description: "[brief]"
```

## Good / bad examples

Provide 2–3 examples of copy in this artist's voice and 2–3 that fail it.
The skill uses these to calibrate before drafting.

```yaml
good_examples:
  - >
    [e.g. "New album 'Solo' out now. Piano improvisations recorded in my
    home studio in Tampere."]
  - >
    [e.g. "Playing at Teon Tupa tomorrow morning. Piano improvisation
    during coffee hour."]

bad_examples:
  - >
    [e.g. "I'm excited to announce my new album! Join me on this musical
    journey through the landscapes of ambient piano improvisation!"]
  - >
    [e.g. "Come experience the transformative power of healing music as I
    delve into sonic exploration at Teon Tupa!"]
```

## Notes

Any other guidance specific to this artist:

```yaml
notes: >
  [e.g. "Emojis only if explicitly requested. Always mention CC licensing
  in release posts. Keep performance posts short — venue, date, time,
  one sentence description."]
```
