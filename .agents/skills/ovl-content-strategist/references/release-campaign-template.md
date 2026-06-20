# Release Campaign Template

Content campaign milestones for album and single releases. Works alongside
`ovl-release-manager` milestones — this file covers content and promotional
tasks specifically.

Pull the distributor lead time from
`workspace/artists/<id>/release-defaults.md` to calculate upload deadlines.

---

## Album Campaign Timeline

Work backward from the release date. T = release day.

| Milestone | Timing | Content tasks | OVL CLI steps |
| --------- | ------ | ------------- | ------------- |
| Content planning | T−42 | Identify key tracks for spotlights; plan campaign arc | `ovl release show <id>` to get track list |
| Announcement | T−21 | Announcement post + pre-save links go live | `ovl release add-link` when pre-save available |
| First teaser | T−14 | Track preview or snippet (story / short video) | — |
| Ramp up | T−7 | Behind-the-scenes post; reminder post | `ovl validate` |
| Pre-release push | T−3 | Final reminder; save links prominent | `ovl site sync` |
| **Release day** | **T−0** | **Release announcement (morning); thank supporters** | **`ovl site sync`** |
| Post-release week 1 | T+3 to T+7 | Track spotlight × 1; engagement with comments | — |
| Post-release week 2 | T+8 to T+14 | Track spotlight × 1; BTS from recording | — |
| Post-release weeks 3–4 | T+15 to T+28 | Remaining track spotlights; share placements if any | — |

**Total campaign content from one album:**
- 1 announcement post
- 2–3 teasers / snippets
- 2 reminder posts
- 1 release-day post
- n track spotlights (one per track)
- 2–4 behind-the-scenes posts from recording
- 1–2 supporter thank-you posts

**= 15–25 pieces of content across 4–6 weeks**

---

## Single Campaign Timeline

Simpler than album, but still benefits from a planned arc.

| Milestone | Timing | Content tasks | OVL CLI steps |
| --------- | ------ | ------------- | ------------- |
| Announcement | T−14 | "New single coming" post + pre-save if available | `ovl release add-link` |
| Teaser | T−7 | Short clip or snippet | — |
| Reminder | T−3 | Final push | `ovl validate` |
| **Release day** | **T−0** | **Announcement; streaming links** | **`ovl site sync`** |
| Post-release | T+3 to T+14 | 1–2 follow-up posts (BTS, context, story behind track) | — |

---

## Single from Upcoming Album

Use the single to build anticipation for the full release.

**Approach:**
1. Release single 3–6 weeks before album
2. Single gets its own abbreviated campaign (T−14 to T+7)
3. Album announcement drops at T−21 for the album, typically overlapping
   with post-release period of the single
4. Single becomes the "lead track" in album promotion

**Content flow:**
```
Single announcement → Single release → Post-release (single) →
Album announcement ("from the upcoming album...") → Album release
```

---

## Release Day Checklist

- [ ] Announcement post live (morning, not overnight)
- [ ] Streaming links working on all platforms in `ovl site list`
- [ ] `ovl site sync` run to push release links to all sites
- [ ] Bandcamp page updated
- [ ] YouTube upload live (if applicable)
- [ ] Stories / real-time updates through the day
- [ ] Thank supporters in comments or a follow-up story

---

## Campaign Content Format

```markdown
# Content Campaign — [Release Title] — [Artist Name]

**Release date:** [Date]
**Distributor:** [Amuse / DistroKid / etc.]
**Lead time:** [n] days → upload by [Date]

## Campaign arc
[One sentence: what is the story of this release?]

## Key tracks to spotlight
1. [Track] — [Why this one first?]
2. [Track]
...

## Pre-release content (T−21 to T−1)
| Date | Platform | Content | Notes |
| ---- | -------- | ------- | ----- |
| | | | |

## Release day (T−0)
| Time | Platform | Content |
| ---- | -------- | ------- |
| Morning | All | Announcement post |
| Afternoon | Primary | Stories / real-time |

## Post-release content (T+1 to T+28)
| Date | Platform | Content | Notes |
| ---- | -------- | ------- | ----- |
| | | | |

## Batching opportunities
- [Date]: Record session clips during [planned session]
- [Date]: Performance — capture before/after

## OVL CLI steps in this campaign
- [ ] `ovl release add-link` — pre-save URL when available
- [ ] `ovl validate` — before release day
- [ ] `ovl site sync` — release day morning
```
