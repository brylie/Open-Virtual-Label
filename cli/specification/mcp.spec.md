# CLI Specification — `ovl mcp`

Commands for managing MCP (Model Context Protocol) integrations. MCPs extend the CLI
with connections to external services such as email, calendar, distributors, and archive
storage. Credentials are stored locally and never written to workspace JSON records or
the OVL repository.

---

## Command Reference

---

### `ovl mcp list`

List available MCPs and their connection status.

```text
ovl mcp list
```

**Output:** Table of MCP names, descriptions, connection status, and the commands they
enable.

**Output example:**

```text
MCP                  Status         Enables
gmail                connected      outreach send, outreach follow-up
google-calendar      not connected  (scheduling)
internet-archive     not connected  archive push
amuse                not connected  release submit
```

---

### `ovl mcp connect <mcp-name>`

Connect an MCP integration.

```text
ovl mcp connect <mcp-name>
```

**Behaviour:** Launches the authentication flow for the specified MCP. Credentials are
stored locally and never written to workspace JSON records or the OVL repository.

**Available MCPs:** `gmail`, `google-calendar`, `internet-archive`, `amuse`
(read-only where API permits)

---

### `ovl mcp disconnect <mcp-name>`

Disconnect an MCP integration and remove stored credentials.

```text
ovl mcp disconnect <mcp-name>
```

**Behaviour:** Removes locally stored credentials for the specified MCP. Prompts for
confirmation before removing.
