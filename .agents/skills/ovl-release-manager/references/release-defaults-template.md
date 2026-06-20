# Release Defaults

Artist-specific preferences for the `ovl-release-manager` skill.
Copy this file to `workspace/artists/<artist-id>/release-defaults.md`
and fill in the values.

## Timezone

```yaml
timezone: Europe/Helsinki    # IANA timezone name — used for all date calculations
```

## Default Output Format

Checklist output format when the artist does not specify:

```yaml
output_format: markdown    # markdown | google-keep | google-calendar | google-tasks
```

## PRO / Rights Registration Workflow

Name of the Performing Rights Organisation and the registration steps to
include in every release plan:

```yaml
pro: [e.g. Teosto, ASCAP, PRS, SOCAN, APRA]
```

Registration steps to add as checklist items:

- [ ] [Step 1 — e.g. "Log in to Teosto Online and register the work code"]
- [ ] [Step 2 — e.g. "Add co-writers and their IPI numbers if applicable"]
- [ ] [Step 3 — e.g. "Confirm ISRC is linked to the registration"]

Add or remove steps as appropriate for the PRO. The skill includes these
automatically in every release plan.

## Music Video Production Tools

What tools does the artist use for music videos? The skill uses this to
estimate production time when a video is part of the plan.

```yaml
video_tools:
  quick:    [e.g. "Resolume Avenue", "Synesthesia", "TouchDesigner"]  # 1–2 day turnaround
  advanced: [e.g. "Blender", "After Effects", "DaVinci Resolve"]      # 1–4 week turnaround
```

## Distribution Safety Lead Time Override

Override the default lead time from `distributor-requirements.md` if the
artist prefers a longer safety margin:

```yaml
distributor_lead_time_days: 35    # leave blank to use distributor default
```

## Preferred Promotional Platforms

Ordered list of platforms to include in the promotional timeline. Defaults
to the `platforms` keys in `artist.json` if not set here:

```yaml
promo_platforms:
  - instagram
  - youtube
  - bandcamp
```

## Notes

Any other artist-specific release planning context the skill should know:

```yaml
notes: |
  [free text]
```
