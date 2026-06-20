# Platform Metrics Guide

What to track, what to watch, and how to interpret metrics on each platform.
For calculation formulas, see [metrics-calculations.md](metrics-calculations.md).

---

## Spotify for Artists

**Primary metrics to track monthly:**
- Monthly listeners (unique, last 28 days)
- Total streams
- Top 5 tracks by streams
- New playlist additions

**Secondary metrics:**
- Save rate (saves / listeners)
- Skip rate (exits before 30 seconds)
- New vs. returning listener ratio
- Top countries

**Good signs:**
- Consistent listener growth (any positive trend)
- Save rate >10%
- Skip rate <30%
- Appearing in user-generated playlists
- Balanced new/returning ratio (~30% returning = healthy retention)

**Warning signs:**
- Skip rate >60% (music not connecting with who's finding it)
- Save rate <5% (little lasting interest)
- Listeners declining despite continued releases
- No playlist presence (discoverability problem)

**Interpretation examples:**

> 47 monthly listeners, 32 new, 15 returning: 32% retention. Decent for an
> early-stage artist. Each release should convert some new listeners into
> returning ones over time.

> Track A: 120 streams, 18 saves (15% save rate). Track B: 85 streams,
> 4 saves (5% save rate). Track A resonates more strongly — promote it and
> consider creating more music in that style.

---

## Apple Music for Artists

**Primary metrics to track monthly:**
- Listeners
- Plays
- Plays per listener

**Secondary metrics:**
- Shazams (if any — indicates discoverability)
- Source breakdown (library, search, browse, radio)

**Good signs:**
- Plays per listener >3 (people returning)
- Growing Shazam count
- Diverse discovery sources

**Warning signs:**
- Plays per listener <2 (mostly casual or accidental plays)
- All traffic from a single source

---

## YouTube

**Primary metrics to track monthly:**
- Subscribers (total and change)
- Views this month
- Watch time (hours)
- Top videos

**Secondary metrics:**
- Click-through rate (CTR): impressions → clicks
- Average view duration / retention percentage
- Traffic sources (search vs. browse vs. suggested vs. external)

**Benchmarks:**
- CTR <2%: Thumbnail or title not working
- CTR 2–5%: Average
- CTR 5–10%: Good
- CTR >10%: Excellent
- Retention <30%: Content loses people early
- Retention 40–60%: Good for music content
- Retention >60%: Excellent

**Good signs:**
- Consistent subscriber growth
- Healthy CTR (>5%) on thumbnails
- Retention >40%
- Traffic from search (indicates discoverability)

**Warning signs:**
- Low CTR despite views (title/thumbnail disconnect)
- Poor retention (content doesn't engage quickly)
- Only browse traffic (not searchable)

**Interpretation examples:**

> 500 impressions, 50 views (10% CTR); 4:30 avg view of a 10:00 video
> (45% retention): Excellent CTR, good retention. This title and format
> are working — apply the same approach to future uploads.

---

## Instagram

**Primary metrics to track monthly:**
- Followers (total and monthly change)
- Average reach per post
- Engagement rate (engagements / reach)
- Number of posts published

**Secondary metrics:**
- Story views and completion rate
- Profile visits → website clicks conversion
- Best performing post (reach + engagement)
- Best posting days/times

**Good signs:**
- Follower growth (any positive trend)
- Engagement rate >3%
- Reach growing alongside (not just followers)
- Story completion rate >50%
- Profile visits converting to website clicks

**Warning signs:**
- Followers growing but engagement declining (wrong audience or inactive followers)
- Reach much lower than followers (algorithm suppression or inactive audience)
- Low story completion (content too long or unengaging)
- Profile visits but no website clicks (bio link or CTA problem)

**Interpretation examples:**

> Post reached 120 accounts (48% of 250 followers); 15 engagements (12.5%
> engagement rate): Excellent. Some discovery happening beyond followers —
> the algorithm is amplifying the content.

> Story views: 45, 38, 32, 28, 20. Completion: 44% (20/45). Losing half
> the audience — consider shorter stories or putting the most engaging
> content first.

---

## Facebook

**Primary metrics to track monthly:**
- Followers
- Average reach per post
- Engagement per post

**Secondary metrics:**
- Link clicks (more useful than impressions)
- Top posts

**Notes:**
- Facebook is generally low priority for independent musicians; treat it
  as a cross-posting channel and event promotion tool
- Organic reach is typically low (10–30% of followers is normal)
- Events and local announcements tend to perform better than general posts

---

## Bandcamp

**Primary metrics to track monthly:**
- Page visits
- Plays
- Sales and revenue (coordinate with `ovl-finance-manager`)
- Top referral sources

**Secondary metrics:**
- Conversion rate (visits → purchases)
- Geographic distribution
- Most-played tracks

**Good signs:**
- Visitors playing tracks (not just loading the page)
- Any sales at all
- Referrals from social media (content is driving traffic)

**Warning signs:**
- Visits but no plays (people leaving without engaging)
- Plays but no sales consistently (price point or value-framing issue)
- Single traffic source only

**Conversion rate benchmarks:**
- 0.5–2%: Normal for indie music
- 2–5%: Good
- 5%+: Excellent

---

## Cross-platform analysis

### Platform priorities

For each platform, assess four dimensions (↑ growing / → stable / ↓ declining):
- **Growth:** Follower/listener change rate
- **Engagement:** Interaction quality
- **Conversion:** Traffic to streams/sales
- **Efficiency:** Results per hour invested

Summarise in the report's Platform Priorities section.

### Comparing meaningfully

**Compare:** Growth rates (%), engagement rates (%), time ROI.

**Don't compare:** Raw follower counts across platforms; absolute stream
counts against social follower counts; metrics from different time periods.

### Anomaly detection

Flag any metric that moves >2× its recent average monthly change.
Investigate: release spike? algorithm change? external share? seasonal
pattern? One anomalous month isn't a trend — wait for a second data point
before drawing conclusions.

### Early-stage benchmarks (ambient/independent music)

These are rough orientation points, not targets:

| Platform | Early stage | Building momentum | Established small base |
| -------- | ----------- | ----------------- | ---------------------- |
| Spotify monthly listeners | 10–50 | 50–200 | 200–1,000 |
| Instagram followers | 50–200 | 200–500 | 500–2,000 |
| YouTube subscribers | 10–50 | 50–200 | 200–1,000 |
| Monthly revenue | €5–20 | €20–50 | €50–100 |

The most important thing at early stage is direction and consistency, not
the absolute numbers.
