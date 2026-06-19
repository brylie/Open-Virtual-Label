# CLI Specification — `ovl qc`

Commands for running pre-release quality checks on releases and tracks.
QC invokes the `qc` agent, which checks both track-level and release-level
requirements before a release can advance to `ready`.

---

## Command Reference

---

### `ovl qc check --release <release-id>`

Run the pre-release quality check on a release.

```text
ovl qc check --release <release-id>
```

**Behaviour:**

- Invokes the `qc` agent
- Checks all track-level requirements (ISRC, mastering data completeness, file paths)
  and release-level requirements (artwork, metadata, dates, license)
- Presents the full QC report
- **[CONFIRMATION GATE]** On a clean pass: asks artist to confirm before setting
  `release.qc.passed: true`
- On failures: lists each failure with the field and requirement; artist must resolve or
  record an override before QC can pass

**Output example:**

```text
QC Report — Spectra

Tracks (8/8 checked):
  ✓ chromatic-drift     ISRC ✓  Mastering ✓  File ✓
  ✓ passage             ISRC ✓  Mastering ✓  File ✓
  ✗ luminance           ISRC ✗  (not assigned — run ovl isrc assign)
  ...

Release:
  ✓ Artwork: 3000×3000px PNG
  ✓ License: CC BY 4.0
  ✓ Target date: 2025-09-01
  ✗ Submission deadline: not set

2 failures. Resolve before advancing to ready.
```
