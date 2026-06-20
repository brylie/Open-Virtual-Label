# Outreach Preferences

Artist-specific preferences for the `ovl-licensing-outreach` skill.
Copy this file to `workspace/artists/<artist-id>/outreach-preferences.md`
and fill in the values.

## Priority niches

Ranked list of creator types to focus on. The skill uses this to decide
where to spend research time each month.

```yaml
priority_niches:
  - podcasts          # e.g. meditation, philosophy, nature
  - youtube           # e.g. nature videography, study content
  - indie_games       # e.g. atmospheric, exploration genres
  - documentaries     # e.g. environmental, contemplative
  - apps              # e.g. meditation, wellness, focus
  - educational       # e.g. online courses, tutorial channels
```

## Niche notes

Free text notes on what to look for or avoid in each niche:

```yaml
podcasts:     Focus on mindfulness, philosophy, and contemplative topics.
              Avoid: news, comedy, true crime.

youtube:      Nature videography and study-with-me channels. Look for
              channels currently using generic library music.

indie_games:  Atmospheric exploration and narrative games. Avoid fast-paced
              or action genres.
```

## Music description for outreach

One or two sentences describing the music as you'd explain it to a prospect.
This is used to personalise the value proposition in emails.

```yaml
music_description: >
  [e.g. "I create ambient piano music — slow, atmospheric, and unobtrusive.
  It works well as background for contemplative content and fits naturally
  into spaces that benefit from a calm, focused mood."]
```

## Preferred contact platform

What to use when a website email isn't available:

```yaml
preferred_contact: email          # email | instagram | linkedin | twitter
```

## Preferred preview platforms

Which platforms to link to in outreach emails (ordered by preference):

```yaml
preview_platforms:
  - bandcamp
  - spotify
  - website
```

## Monthly outreach goal

Approximate number of personalised emails per month:

```yaml
monthly_emails: 5    # realistic for a solo artist with limited time
```

## Tone notes

Any guidance on voice and tone for this artist's outreach:

```yaml
tone: >
  [e.g. "Warm and genuine. Slightly informal but professional. Lead with
  curiosity about their work. Never pushy."]
```

## Avoid

Types of creators or topics to skip entirely:

```yaml
avoid:
  - news podcasts
  - comedy content
  - fast-paced action games
  - commercial / corporate content
  - large platforms (>100k subscribers — unlikely to respond)
```
