# Finance Config

Artist-specific financial settings for the `ovl-finance-manager` skill.
Copy this file to `workspace/artists/<artist-id>/finance-config.md` and
fill in the values.

---

## Revenue goal

```yaml
revenue_goal_per_month: 100      # target in currency below
currency: EUR
notes: >
  [Context on what this goal represents — e.g. "sustainable monthly income
  from music alone, not a living wage, but meaningful contribution"]
```

## Revenue sources

```yaml
revenue_sources:
  primary:
    - platform: bandcamp          # highest margin; direct fan payments
      payout: immediate
      priority: 1
    - platform: streaming         # via distributor; 2–3 month delay
      payout: monthly_delayed
      priority: 2
  secondary:
    - platform: youtube           # ad revenue once monetised
      payout: monthly
      priority: 3
    - platform: licensing         # track separately from streaming
      payout: varies
      priority: 4
```

## Distributor

```yaml
distributor: [e.g. Amuse, DistroKid, TuneCore]
notes: >
  [Any distribution notes — e.g. free tier vs. pro, payment timing,
  which catalog is under which name]
```

## Expense categories

```yaml
expense_categories:
  - software_vst
  - subscriptions
  - promotion
  - performance
  - equipment
  - distribution_admin
  - education            # optional
```

## Milestone targets (revenue)

```yaml
milestones:
  monthly_revenue:
    - 10
    - 25
    - 50
    - 100    # primary goal
```

## Notes

```yaml
notes: >
  [Any context relevant to interpreting the numbers — e.g. "Bandcamp Fridays
  spike sales and should not inflate rolling averages", payment delay
  patterns, seasonal variation, split between recurring and one-off income]
```
