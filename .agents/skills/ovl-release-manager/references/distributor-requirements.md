# Distributor Requirements

Lead times and upload requirements by distributor. The active distributor for
an artist is set in `workspace/artists/<id>/artist.json` under
`distribution.distributor`. Use that value to look up requirements here.

## Amuse

- **Lead time:** 28 days (use as safety minimum; 35 days recommended for
  editorial pitch eligibility)
- **File requirements:** WAV 16-bit or 24-bit, 44.1 kHz minimum
- **Artwork:** 3000×3000 px minimum, JPEG or PNG, RGB colour space
- **Metadata:** Title, artist, album, release date, genre, ISRC (auto-assigned
  if not provided), UPC (auto-assigned)
- **Spotify pitch:** Eligible if uploaded at least 7 days before Spotify's
  own pitch deadline (T−21 from release date)
- **Notes:** Free tier has limited release slots per year; Pro tier unlimited

## DistroKid

- **Lead time:** 1–5 business days (use 14 days as safety minimum for
  editorial pitch eligibility)
- **File requirements:** WAV or FLAC, 44.1 kHz or higher
- **Artwork:** 3000×3000 px minimum, JPEG or PNG
- **Notes:** Spotify pitch requires upload at least 7 days before T−21;
  Yearly subscription model

## TuneCore

- **Lead time:** 1–2 business days (use 14 days as safety minimum)
- **File requirements:** WAV, 44.1 kHz or higher, 16-bit or 24-bit
- **Artwork:** 3000×3000 px minimum, JPEG or PNG
- **Notes:** Per-release pricing model; Spotify pitch same window as above

## CD Baby

- **Lead time:** 3–5 business days (use 14 days as safety minimum)
- **File requirements:** WAV, 44.1 kHz, 16-bit minimum
- **Artwork:** 1400×1400 px minimum (3000×3000 recommended)
- **Notes:** One-time per-release fee; physical distribution also available

## SoundCloud for Artists

- **Lead time:** Varies; check current platform requirements
- **Notes:** Primarily for SoundCloud distribution; cross-platform
  distribution capabilities vary by plan

## Adding a new distributor

Add a section following the same structure. The `ovl-release-manager` skill
reads distributor from `artist.json` and looks it up here. If the distributor
is not listed, the skill will ask the artist for the lead time.
