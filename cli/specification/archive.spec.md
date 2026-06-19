# CLI Specification — `ovl archive`

Commands for packaging and uploading releases to long-term archive storage.
The primary archive target is Internet Archive via the IA S3-compatible API.
A secondary object storage location can be configured in `workspace/label/profile.json`.

---

## Command Reference

---

### `ovl archive push --release <release-id>`

Package and upload a release to long-term archive storage.

```text
ovl archive push --release <release-id> [--skip-stems] [--skip-project-files]
```

**Behaviour:**

- Assembles a manifest of all files to upload: master WAVs, stems, project files,
  artwork, and JSON records
- **[CONFIRMATION GATE]** Presents the manifest for artist review. Lists each file, its
  size, and destination. Cannot be bypassed with `--yes`
- Uploads to Internet Archive (primary) via the IA S3-compatible API
- Uploads to secondary object storage if configured in `workspace/label/profile.json`
- Verifies SHA-256 checksums after upload
- Writes `release.archive{}` fields: IDs, URLs, paths, flags, checksums verified,
  archive date

**Options:**

| Option | Description |
|---|---|
| `--skip-stems` | Omit stems from the archive package. |
| `--skip-project-files` | Omit DAW project files from the archive package. |

**Errors:**

- `IA_MCP_NOT_CONFIGURED` — Internet Archive MCP not connected; run
  `ovl mcp connect internet-archive`
- `FILE_NOT_FOUND` — a file referenced in a track record does not exist at the stated
  path
- `UPLOAD_PARTIAL` — upload interrupted mid-way; re-run to resume (checksums prevent
  re-uploading completed files)

---

### `ovl archive status --release <release-id>`

Show archive status for a release.

```text
ovl archive status --release <release-id>
```

**Output:** Table of archive flags, URLs, and checksum verification status from
`release.archive{}`.
