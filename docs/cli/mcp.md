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

### Resource Subscriptions

Clients can subscribe to resource changes via `resources/subscribe`.
The server polls for file mtime changes (default: 5 seconds) and
emits `notifications/resources/updated` when a subscribed file
changes on disk.

---

## Tools

Tools expose ctx commands as callable operations. Each tool accepts
JSON arguments and returns text results.

### `ctx_status`

Show context health: file count, token estimate, and per-file summary.

**Arguments:** None. **Read-only.**

### `ctx_add`

Add a task, decision, learning, or convention to the context.

| Argument      | Type   | Required    | Description                                      |
|---------------|--------|-------------|--------------------------------------------------|
| `type`        | string | Yes         | Entry type: task, decision, learning, convention |
| `content`     | string | Yes         | Title or main content                            |
| `priority`    | string | No          | Priority level (tasks only): high, medium, low   |
| `context`     | string | Conditional | Context field (decisions and learnings)          |
| `rationale`   | string | Conditional | Rationale (decisions only)                       |
| `consequence` | string | Conditional | Consequence (decisions only)                     |
| `lesson`      | string | Conditional | Lesson learned (learnings only)                  |
| `application` | string | Conditional | How to apply (learnings only)                    |

### `ctx_complete`

Mark a task as done by number or text match.

| Argument | Type   | Required | Description                              |
|----------|--------|----------|------------------------------------------|
| `query`  | string | Yes      | Task number (e.g. "1") or search text    |

### `ctx_drift`

Detect stale or invalid context. Returns violations, warnings, and
passed checks.

**Arguments:** None. **Read-only.**

### `ctx_journal_source`

Query recent AI session history (summaries, decisions, topics).

| Argument | Type   | Required | Description                                            |
|----------|--------|----------|--------------------------------------------------------|
| `limit`  | number | No       | Max sessions to return (default: 5)                    |
| `since`  | string | No       | ISO date filter: sessions after this date (YYYY-MM-DD) |

**Read-only.**

### `ctx_watch_update`

Apply a structured context update to `.context/` files. Supports
task, decision, learning, convention, and complete entry types.
Human confirmation is required before calling.

| Argument      | Type   | Required    | Description                                                |
|---------------|--------|-------------|------------------------------------------------------------|
| `type`        | string | Yes         | Entry type: task, decision, learning, convention, complete |
| `content`     | string | Yes         | Main content                                               |
| `context`     | string | Conditional | Context background (decisions/learnings)                   |
| `rationale`   | string | Conditional | Rationale (decisions only)                                 |
| `consequence` | string | Conditional | Consequence (decisions only)                               |
| `lesson`      | string | Conditional | Lesson learned (learnings only)                            |
| `application` | string | Conditional | How to apply (learnings only)                              |

### `ctx_compact`

Move completed tasks to the archive section and remove empty
sections from context files. Human confirmation required.

| Argument  | Type    | Required | Description                                              |
|-----------|---------|----------|----------------------------------------------------------|
| `archive` | boolean | No       | Also write tasks to `.context/archive/` (default: false) |

### `ctx_next`

Suggest the next pending task based on priority and position.

**Arguments:** None. **Read-only.**

### `ctx_check_task_completion`

Advisory check: after a write operation, detect if any pending tasks
were silently completed. Returns nudge text if a match is found.

| Argument        | Type   | Required | Description                             |
|-----------------|--------|----------|-----------------------------------------|
| `recent_action` | string | No       | Brief description of what was just done |

**Read-only.**

### `ctx_session_event`

Signal a session lifecycle event. Type `end` triggers the session-end
persistence ceremony - human confirmation required.

| Argument | Type   | Required | Description                                                  |
|----------|--------|----------|--------------------------------------------------------------|
| `type`   | string | Yes      | Event type: start, end                                       |
| `caller` | string | No       | Caller identifier (cursor, windsurf, vscode, claude-desktop) |

### `ctx_steering_get`

Retrieve applicable steering files for a prompt. Without a prompt,
returns always-included files only.

| Argument | Type   | Required | Description                                                |
|----------|--------|----------|------------------------------------------------------------|
| `prompt` | string | No       | Prompt text to match against steering file descriptions    |

**Read-only.**

### `ctx_search`

Search across `.context/` files for a query string. Returns matching
lines with file paths and line numbers.

| Argument | Type   | Required | Description                    |
|----------|--------|----------|--------------------------------|
| `query`  | string | Yes      | Search string to match against |

**Read-only.**

### `ctx_session_start`

Execute session-start hooks and return aggregated context from hook
outputs.

**Arguments:** None.

### `ctx_session_end`

Execute session-end hooks with an optional summary. Returns aggregated
context from hook outputs.

| Argument  | Type   | Required | Description                          |
|-----------|--------|----------|--------------------------------------|
| `summary` | string | No       | Session summary passed to hook scripts |

### `ctx_remind`

List pending session-scoped reminders.

**Arguments:** None. **Read-only.**

---

## Prompts

Prompts provide pre-built templates for common workflows. Clients
can list available prompts via `prompts/list` and retrieve a
specific prompt via `prompts/get`.

### `ctx-session-start`

Load full context at the beginning of a session. Returns all context
files assembled in priority read order with session orientation
instructions.

### `ctx-add-decision`

Format an architectural decision entry with all required fields.

| Argument       | Type   | Required | Description                    |
|----------------|--------|----------|--------------------------------|
| `content`      | string | Yes      | Decision title                 |
| `context`      | string | Yes      | Background context             |
| `rationale`    | string | Yes      | Why this decision was made     |
| `consequence`  | string | Yes      | Expected consequence           |

### `ctx-add-learning`

Format a learning entry with all required fields.

| Argument      | Type   | Required | Description                     |
|---------------|--------|----------|---------------------------------|
| `content`     | string | Yes      | Learning title                  |
| `context`     | string | Yes      | Background context              |
| `lesson`      | string | Yes      | The lesson learned              |
| `application` | string | Yes      | How to apply this lesson        |

### `ctx-reflect`

Guide end-of-session reflection. Returns a structured review prompt
covering progress assessment and context update recommendations.

### `ctx-checkpoint`

Report session statistics: tool calls made, entries added, and
pending updates queued during the current session.


