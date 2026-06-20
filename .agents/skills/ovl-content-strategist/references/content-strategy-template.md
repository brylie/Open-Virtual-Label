# Content Strategy

Artist-specific content preferences for the `ovl-content-strategist` skill.
Copy this file to `workspace/artists/<artist-id>/content-strategy.md`
and fill in the values.

---

## Platform priorities

Ordered list of platforms the artist actively uses for content. The skill
focuses effort here, in this order.

```yaml
platforms:
  - instagram        # community, BTS, release announcements
  - youtube          # full tracks, albums, process videos
  - facebook         # events, cross-posts
  # Add or remove as appropriate:
  # - tiktok
  # - bandcamp       # direct release posts
  # - substack       # newsletter / longer writing
```

## Time budget

Realistic hours available for content creation per month (not including
music production itself):

```yaml
content_hours_per_month: 4     # e.g. 2–3 for a full-time worker; 8–10 for part-time
```

## Content mix

Target ratio across all platforms:

```yaml
content_mix:
  promotional: 35%      # release announcements, streaming links, pre-saves
  engagement: 65%       # BTS, process, performance recaps, community
```

## Posting frequency (per platform)

Realistic targets — better to post less and maintain it than overcommit:

```yaml
posting_frequency:
  instagram_feed:    8       # posts per month
  instagram_stories: 10     # per month (casual, low-effort)
  youtube:           2       # uploads per month
  facebook:          4       # posts per month
```

## Batching preference

How the artist prefers to create content:

```yaml
batching:
  session_frequency: monthly        # monthly | weekly | per-release
  preferred_batch_day: [Saturday]   # day(s) of week for batch sessions
  batch_session_length_hours: 2
```

## Preferred content types

Which types of content the artist is comfortable producing (ranked or
listed in rough priority):

```yaml
preferred_content_types:
  - behind_the_scenes     # setup photos, session clips — lowest effort
  - track_spotlight       # highlight from catalog — low effort
  - process_post          # how something was made
  - performance_recap     # post-show photos and reflection
  - release_announcement  # per release, required
  - livestream            # if sustainable; remove if not
```

## Seasonal notes

Any artist-specific notes about how seasons affect content rhythm:

```yaml
seasonal_notes:
  winter: >
    [e.g. "Heavy studio season — lean into process and BTS content."]
  spring: >
    [e.g. "Performance season starts — increase announcement posts."]
  summer: >
    [e.g. "Peak activity — releases and performances."]
  autumn: >
    [e.g. "Preparation for end-of-year release — start building content buffer."]
```

## Tone notes

Voice and tone guidance for content planning:

```yaml
tone: >
  [e.g. "Authentic and direct. Share what is actually happening in the
  creative process. Avoid hype or promotional language. Community-first."]
```

## Avoid

Content types or topics to skip:

```yaml
avoid:
  - [e.g. "listicles or tips content — not authentic to this artist"]
  - [e.g. "trend-chasing — feels inauthentic"]
  - [e.g. "political or divisive topics"]
```
