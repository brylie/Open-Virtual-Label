# Metrics Config

Artist-specific goals, platform list, and metric thresholds for the
`ovl-metrics-analyst` skill. Copy this file to
`workspace/artists/<artist-id>/metrics-config.md` and fill in the values.

---

## Revenue goal

```yaml
revenue_goal_per_month: 100    # target monthly revenue in currency below
currency: EUR
```

## Platforms to track

```yaml
platforms:
  streaming:
    - spotify
    - apple_music
    - youtube_music    # optional — omit if not tracking separately
    - bandcamp
  social:
    - instagram
    - youtube
    - facebook         # optional — omit if not maintaining
```

## Milestone targets

Thresholds the analyst should flag when reached:

```yaml
milestones:
  spotify_monthly_listeners:
    - 100
    - 250
    - 500
    - 1000
  instagram_followers:
    - 100
    - 250
    - 500
    - 1000
  youtube_subscribers:
    - 50
    - 100
    - 250
    - 500
  monthly_revenue_eur:
    - 10
    - 25
    - 50
    - 100
```

## Engagement rate expectations

Thresholds for flagging underperformance:

```yaml
engagement_thresholds:
  instagram_min_percent: 3.0     # flag if below this
  youtube_ctr_min_percent: 2.0
  youtube_retention_min_percent: 30
```

## Notes

Any context that helps interpret the numbers:

```yaml
notes: >
  [e.g. "Revenue goal is a sustainable monthly income from music alone.
  Bandcamp direct sales are the priority channel — streaming revenue is
  supplementary. Engagement rate matters more than follower count."]
```
