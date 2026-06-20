# Specialist Agents

Describes each specialist sub-agent available to the Artist Manager, when to
invoke them, and what to expect back. Adapt this list to the skills actually
installed in `.agents/skills/`.

## `ovl-finance-manager`

**Purpose:** Tracks revenue and expenses; produces monthly financial
summaries; monitors progress toward the artist's revenue target. Reads
`finance-config.md` and historical records from
`workspace/artists/<id>/finances/`.

**Invoke when:**

- Artist wants to know how much they earned this month
- Reviewing progress toward a revenue goal
- Recording a new expense or income source
- Planning budget for a project or promotion
- Artist asks "can I afford X?"

**Returns:** Revenue by platform, expense summary, trend over time, progress
toward target, budget decision analysis.

## `ovl-social-media-specialist`

**Purpose:** Writes platform-specific copy for Instagram, YouTube, Facebook,
or other channels the artist uses. Reads artist voice config and EPK from
the workspace. Does not post — delivers copy for the artist to paste directly.

**Invoke when:**

- Album, EP, or single announcement needed
- Promoting a live performance or livestream
- Planning a content push around a release
- Track spotlight or behind-the-scenes caption needed
- Placement announcement copy (after a licensing win)

**Returns:** Ready-to-paste text for each requested platform, with character
counts and hashtags.

## `ovl-metrics-analyst`

**Purpose:** Compiles streaming stats, social media growth, and engagement
data across platforms. Identifies trends, forecasts goal timelines, and
surfaces actionable insights. Reads `metrics-config.md` and historical
snapshots from `workspace/artists/<id>/metrics/`.

**Invoke when:**

- Monthly metrics review
- Evaluating whether a content strategy is working
- Artist wants to know which releases or posts are performing best
- Preparing for a quarterly planning session
- Forecasting when a milestone or revenue goal will be reached

**Returns:** Metrics summary by platform, trend analysis, top performers,
recommendations.

## `ovl-content-strategist`

**Purpose:** Plans content calendars for the month or quarter. Coordinates
content around release windows, seasonal moments, and platform-specific
opportunities. Reads the artist's content-strategy.md and release schedule.

**Invoke when:**

- Planning next month's posting schedule
- Coordinating content around an upcoming release
- Artist wants a content batch plan for efficiency
- Quarterly planning session

**Returns:** Monthly or quarterly content calendar with themes, post types,
timing per platform, and batching session recommendations.

## `ovl-licensing-outreach`

**Purpose:** Researches and identifies opportunities for the artist's music
to be used in podcasts, YouTube videos, games, films, or other media.
Drafts personalised outreach emails. Manages the prospect pipeline via
`workspace/outreach/*.json`. Nothing is sent without explicit artist approval.

**Invoke when:**

- Artist wants to find new licensing opportunities
- Following up on a previous outreach thread
- Building or reviewing the prospect pipeline
- Drafting or refining an outreach email

**Returns:** Scored prospect list, personalised outreach drafts (for artist
approval), pipeline status summary, follow-up recommendations.

## Performance Coordinator

**Purpose:** Manages scheduling and logistics for live performances. Handles
venue outreach, preparation checklists, and post-show follow-up.

**Invoke when:**

- Planning or booking a live performance
- Reaching out to a new venue
- Preparing materials and repertoire for an upcoming show
- Post-performance debrief and documentation

**Returns:** Performance timeline, venue communication drafts, preparation
checklist, post-show summary.
