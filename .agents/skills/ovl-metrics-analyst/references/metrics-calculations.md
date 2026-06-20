# Metrics Calculations

Formulas and worked examples for key performance metrics. All Python
examples can be run in a sandbox when the artist provides CSV data.

---

## Growth rate

### Month-over-month (absolute)

```
Growth = Current – Previous
```

**Example:** 47 − 38 = +9 listeners

### Month-over-month (percentage)

```
Growth % = ((Current − Previous) / Previous) × 100
```

**Example:** ((47 − 38) / 38) × 100 = 23.7%

### Average monthly growth rate (rolling)

Calculate percentage growth for each month, then average them.

```python
rates = [0.15, 0.18, 0.24, 0.12]          # recent monthly growth rates
avg_rate = sum(rates) / len(rates)          # 0.1725 → 17.25% per month
```

---

## Engagement rate

```
Engagement Rate = (Likes + Comments + Shares + Saves) / Reach × 100
```

**Example:** (8 + 2 + 1) / 120 × 100 = 9.2%

**Benchmarks for music content:**
- Below 1%: Poor
- 1–3%: Average
- 3–5%: Good
- 5%+: Excellent

---

## Streaming metrics

### Average streams per listener

```
Avg = Total Streams / Monthly Listeners
```

**Interpretation:** 1–2 = casual; 3–5 = engaged; 6+ = highly engaged

### Save rate

```
Save Rate = Saves / Listeners × 100
```

**Benchmarks:** <5% low; 5–15% average; 15–30% good; 30%+ excellent

### Stream completion rate

```
Completion = Complete Listens / Total Streams × 100
```

**Benchmarks:** <30% poor; 30–60% average; 60–80% good; 80%+ excellent

---

## Revenue metrics

### Revenue per stream

```
Per Stream = Total Streaming Revenue / Total Streams
```

**Typical rates:** Spotify €0.003–0.005; Apple Music €0.007–0.010;
YouTube Music €0.001–0.003

### Bandcamp conversion rate

```
Conversion = Purchases / Page Visits × 100
```

**Benchmarks:** 0.5–2% normal; 2–5% good; 5%+ excellent

### Direct support ratio

```
Direct Ratio = Bandcamp Revenue / Total Revenue × 100
```

Higher percentage = more sustainable income.

---

## Forecasts

### Linear projection (months to target)

```
Months = (Target − Current) / Average Monthly Growth
```

**Example:** (100 − 47) / 9 = 5.9 months

### Percentage-based projection

```
Future = Current × (1 + Monthly Rate)^Months
```

To find months to target:

```python
import math
months = math.log(target / current) / math.log(1 + rate)
```

**Example:** math.log(100 / 47) / math.log(1.18) = 4.9 months

### Conservative / realistic / optimistic range

Always present three scenarios:

```python
rates    = [0.15, 0.18, 0.24, 0.12]
conservative = min(rates)           # 0.12
realistic    = sum(rates)/len(rates)  # 0.1725
optimistic   = max(rates)           # 0.24
```

---

## Content performance

### Comparing content types

```
Score = Engagement Rate × Reach
Average = sum(Scores) / number of posts in type
```

**Example:**
- BTS posts average score: 6.95
- Announcement posts average score: 2.2
- BTS performs 3.2× better

### Best posting time

Group posts by time slot, calculate average reach per slot:

```
Avg Reach = Total Reach for Slot / Posts in Slot
```

---

## Platform efficiency

### Followers per post

```
Followers/Post = New Followers / Posts Made
```

### Time investment ROI

```
ROI = Results / Hours Invested
```

Use consistently — compare platforms on the same result metric.

---

## Insight structure

Every insight should follow this pattern:

1. **Observation:** what the data shows (number + context)
2. **Insight:** what it means
3. **Action:** what to do about it
4. **Expected outcome:** why it matters

**Example:**

```
Observation: BTS posts average 85 reach vs. 45 for announcements (+89%).

Insight: The audience prefers authentic process content over promotional
posts.

Action: Shift Instagram mix to 70% BTS / 30% announcements
(currently 50/50).

Expected outcome: Higher overall reach and engagement, driving more
profile visits and streaming traffic.
```

---

## Python reference for CSV data

When the artist provides a CSV export:

```python
import pandas as pd
import math

df = pd.read_csv('streaming_data.csv')

# Monthly totals
df['date'] = pd.to_datetime(df['date'])
monthly = df.groupby(df['date'].dt.to_period('M'))['streams'].sum()

# Growth rates
growth_rates = monthly.pct_change().dropna().tolist()

# Forecast
current = monthly.iloc[-1]
target = 100
avg_rate = sum(growth_rates) / len(growth_rates)
months_to_target = math.log(target / current) / math.log(1 + avg_rate)
print(f"Estimated months to {target}: {months_to_target:.1f}")
```
