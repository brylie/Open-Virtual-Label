# CLI Specification — `ovl content` and `ovl social`

Commands for generating content campaign briefs and social media copy.
Content planning invokes the `content-strategist` agent; social copy invokes the
`social-media` agent. Social copy requires an approved content brief from
`ovl content brief` before it can be generated.

---

## Command Reference

---

### `ovl content brief --release <release-id>`

Generate a content campaign brief for a release.

```text
ovl content brief --release <release-id>
```

**Behaviour:**

- Invokes the `content-strategist` agent
- Agent reads release record, artist profile, and recent metrics snapshot
- Produces a campaign brief: announcement timing, platform-specific plan, track
  spotlight schedule
- **[APPROVAL GATE]** Artist reviews and approves the brief before social copy is
  generated

---

### `ovl social draft --release <release-id>`

Generate social media copy for a release campaign.

```text
ovl social draft --release <release-id> [--platform instagram|youtube|facebook]
```

**Behaviour:**

- Invokes the `social-media` agent
- Agent reads the approved content brief and release record
- Produces platform-specific copy for artist review and posting
- If `--platform` is omitted, generates copy for all configured platforms

**Options:**

| Option | Description |
|---|---|
| `--platform <name>` | Limit output to one platform: `instagram`, `youtube`, `facebook`. If omitted, generates for all. |

**Requires:** An approved content brief. Run `ovl content brief --release <release-id>`
first.

**Errors:**

- `NO_APPROVED_BRIEF` — no approved content brief found for this release
