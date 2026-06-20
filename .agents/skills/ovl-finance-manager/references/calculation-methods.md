# Financial Calculation Methods

Formulas and worked examples for key financial metrics.
For report formats, see [report-templates.md](report-templates.md).

---

## Basic revenue metrics

### Monthly total

```
Total = sum of all platform payments received in the month
```

Note: Streaming payments typically arrive 2–3 months after streams occur.
Record in the month received, and note the delay in the report.

### Platform percentage

```
Platform % = (Platform Revenue / Total Revenue) × 100
```

**Example:** €12.50 Bandcamp / €20.00 total = 62.5%

### 3-month rolling average

```
Average = (Month N + Month N−1 + Month N−2) / 3
```

Use the rolling average (not a single month) as the canonical "current
monthly revenue" figure when assessing progress toward the goal. This
smooths release-month spikes and payment-timing anomalies.

---

## Growth metrics

### Month-over-month (absolute)

```
Growth = Current − Previous
```

### Month-over-month (percentage)

```
Growth % = ((Current − Previous) / Previous) × 100
```

### Average monthly growth rate (rolling)

```python
rates = [0.08, 0.15, 0.10, 0.12]   # recent monthly growth rates
avg_rate = sum(rates) / len(rates)  # → 11.25%
```

### Quarter-over-quarter

```
QoQ % = ((Current Q − Previous Q) / Previous Q) × 100
```

---

## Progress tracking

### Distance to goal (absolute)

```
Distance = Goal − 3-Month Average
```

### Percentage of goal reached

```
% = (3-Month Average / Goal) × 100
```

### Compound growth — months to target

```
Months = log(Goal / Current) / log(1 + Monthly Rate)
```

```python
import math
months = math.log(goal / current) / math.log(1 + rate)
```

**Example:** current = €21.93, goal = €100, rate = 0.1125
→ `math.log(100/21.93) / math.log(1.1125)` ≈ 14.5 months

### Three-scenario projection

```python
rates = [0.08, 0.15, 0.10, 0.12]
conservative = min(rates)           # worst recent month
realistic    = sum(rates)/len(rates) # average
optimistic   = max(rates)           # best recent month
```

Always present all three; single-point forecasts mislead.

---

## Expense metrics

### Net income

```
Net = Total Revenue − Total Expenses
```

### Expense ratio (target: under 50%)

```
Expense Ratio = (Total Expenses / Total Revenue) × 100
```

Under 30% is excellent; 30–50% is healthy; above 50% warrants review.

### ROI on promotional spend

```
ROI % = ((Attributed Revenue − Cost) / Cost) × 100
```

Attribution is difficult — use conservative estimates. Run the smallest
possible test before scaling any promotional spend.

### Payback period

```
Payback (months) = Cost / Monthly Benefit
```

**Example:** €60 software; saves 2 hours/month valued at €15/hour → 2 months

---

## Revenue diversification

### Diversification index (Herfindahl-style)

```
Index = 1 − Σ(Platform share²)
```

**Example:** Bandcamp 60%, streaming 30%, other 10%
→ `1 − (0.36 + 0.09 + 0.01)` = 0.54

Interpretation: 0.0 = single source; 1.0 = perfectly spread; >0.5 = good.

### Recurring vs. one-time ratio

```
Recurring % = Passive/Subscription Revenue / Total Revenue × 100
```

Track this over time — increasing recurring % means more predictable income.

---

## Goal pacing

### Required monthly increase (absolute)

```
Required = (Goal − Current) / Months Remaining
```

### Required monthly growth rate

```python
import math
required_rate = (goal / current) ** (1 / months) - 1
```

### Gap analysis

```
Gap = Required Rate − Actual Rate
```

Positive gap = behind pace; negative gap = ahead of pace.

---

## Python reference for CSV data

When the artist provides a CSV export:

```python
import pandas as pd
import math

df = pd.read_csv('revenue.csv')
df['date'] = pd.to_datetime(df['date'])
df['amount'] = df['amount'].str.replace('€', '').str.replace(',', '.').astype(float)

monthly = df.groupby(df['date'].dt.to_period('M'))['amount'].sum()

# Growth rates
growth_rates = monthly.pct_change().dropna().tolist()

# Forecast
current = monthly.rolling(3).mean().iloc[-1]
goal = 100
avg_rate = sum(growth_rates) / len(growth_rates)
months_to_goal = math.log(goal / current) / math.log(1 + avg_rate)
print(f"Estimated months to goal: {months_to_goal:.1f}")
```

---

## Common mistakes to avoid

- **Comparing different timeframes** — one month vs. three-month average
- **Ignoring payment delays** — streaming pays 2–3 months late; note this
- **Over-attributing** — not all growth is caused by one action
- **Ignoring outliers** — Bandcamp Friday or release month spikes
- **Mixing currencies** — always convert to the artist's base currency
- **Too much precision** — €21.93, not €21.9278
