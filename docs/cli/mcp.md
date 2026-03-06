---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: MCP Server
icon: lucide/plug
---

## `ctx mcp`

Run ctx as a [Model Context Protocol](https://modelcontextprotocol.io)
(MCP) server. MCP is a standard protocol that lets AI tools discover
and consume context from external sources via JSON-RPC 2.0 over
stdin/stdout.

This makes ctx accessible to **any MCP-compatible AI tool** without
custom hooks or integrations:

- Claude Desktop
- Cursor
- Windsurf
- VS Code Copilot
- Any tool supporting MCP

### `ctx mcp serve`

Start the MCP server. This command reads JSON-RPC 2.0 requests from
stdin and writes responses to stdout. It is intended to be launched
by MCP clients, not run directly.

```
ctx mcp serve
```

**Flags:** None. The server uses the configured context directory
(from `--context-dir`, `CTX_DIR`, `.ctxrc`, or the default `.context`).

---

## Configuration

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "ctx": {
      "command": "ctx",
      "args": ["mcp", "serve"]
    }
  }
}
```

### Cursor

Add to `.cursor/mcp.json` in your project:

```json
{
  "mcpServers": {
    "ctx": {
      "command": "ctx",
      "args": ["mcp", "serve"]
    }
  }
}
```

### VS Code (Copilot)

Add to `.vscode/mcp.json`:

```json
{
  "servers": {
    "ctx": {
      "command": "ctx",
      "args": ["mcp", "serve"]
    }
  }
}
```

---

## Resources

Resources expose context files as read-only content. Each resource
has a URI, name, and returns Markdown text.

| URI                          | Name           | Description                                  |
|------------------------------|----------------|----------------------------------------------|
| `ctx://context/constitution` | constitution   | Hard rules that must never be violated       |
| `ctx://context/tasks`        | tasks          | Current work items and their status          |
| `ctx://context/conventions`  | conventions    | Code patterns and standards                  |
| `ctx://context/architecture` | architecture   | System architecture documentation            |
| `ctx://context/decisions`    | decisions      | Architectural decisions with rationale       |
| `ctx://context/learnings`    | learnings      | Gotchas, tips, and lessons learned           |
| `ctx://context/glossary`     | glossary       | Project-specific terminology                 |
| `ctx://context/agent`        | agent          | All files assembled in priority read order   |

The `agent` resource assembles all non-empty context files into a
single Markdown document, ordered by the configured read priority.

---

## Tools

Tools expose ctx commands as callable operations. Each tool accepts
JSON arguments and returns text results.

### `ctx_status`

Show context health: file count, token estimate, and per-file summary.

**Arguments:** None.

### `ctx_add`

Add a task, decision, learning, or convention to the context.

| Argument       | Type   | Required | Description                              |
|----------------|--------|----------|------------------------------------------|
| `type`         | string | Yes      | Entry type: task, decision, learning, convention |
| `content`      | string | Yes      | Title or main content                    |
| `priority`     | string | No       | Priority level (tasks only): high, medium, low |
| `context`      | string | Conditional | Context field (decisions and learnings) |
| `rationale`    | string | Conditional | Rationale (decisions only)             |
| `consequences` | string | Conditional | Consequences (decisions only)          |
| `lesson`       | string | Conditional | Lesson learned (learnings only)        |
| `application`  | string | Conditional | How to apply (learnings only)          |

### `ctx_complete`

Mark a task as done by number or text match.

| Argument | Type   | Required | Description                              |
|----------|--------|----------|------------------------------------------|
| `query`  | string | Yes      | Task number (e.g. "1") or search text    |

### `ctx_drift`

Detect stale or invalid context. Returns violations, warnings, and
passed checks.

**Arguments:** None.


