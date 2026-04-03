# Detailed Design: MCP Server

Modules: mcp/proto, mcp/handler, mcp/server/*, mcp/session

## Overview

The MCP server is a JSON-RPC 2.0 implementation over stdin/stdout
that exposes ctx project context to any MCP-compatible AI tool
(Claude Desktop, Cursor, Windsurf, VS Code Copilot, etc.). It is
100% generic — no agent-specific coupling. Protocol version:
2024-11-05.

```
stdin --> Server.Serve()
  --> parse.Request() [JSON unmarshaling]
  --> dispatch.Do()
      |-- initialize     --> handshake
      |-- ping           --> pong
      |-- resources/*    --> catalog/read/subscribe
      |-- tools/call     --> handler.*()
      |-- prompts/get    --> prompt builders
      |-- [unknown]      --> ErrCodeNotFound
  --> out.*Response()
  --> io.Writer.WriteJSON()
--> stdout
```

## mcp/proto

**Purpose**: JSON-RPC 2.0 message types and MCP protocol constants.

**Key types**:
- `Request`: JSON-RPC request (JSONRPC, ID, Method, Params)
- `Response`: JSON-RPC response (JSONRPC, ID, Result, Error)
- `Notification`: JSON-RPC notification (no ID, no response)
- `RPCError`: error with code/message/data
- `Resource`, `Tool`, `Prompt`: MCP entity definitions
- `InputSchema`, `Property`: JSON Schema for tool parameters
- `ClientCaps`, `ServerCaps`: capability declarations

**Error codes**: ErrCodeParse (-32700), ErrCodeInvalidReq (-32600),
ErrCodeNotFound (-32601), ErrCodeInvalidArg (-32602),
ErrCodeInternal (-32603).

**Dependencies**: none (pure types)

---

## mcp/server

**Purpose**: Main loop: reads stdin, parses JSON-RPC, routes to
dispatch, writes responses to stdout.

**Key types**:
```
Server {
    handler      *handler.Handler
    version      string
    out          *mcpIO.Writer     // mutex-protected stdout
    in           io.Reader         // stdin
    poller       *poll.Poller
    resourceList proto.ResourceListResult  // immutable
}
```

**Exported API**: `New(contextDir, version)`, `Serve()`.

**Data flow**: `Serve()` blocks reading stdin line-by-line with
buffered scanner (configurable max: cfg.ScanMaxSize). Each line
parsed as JSON-RPC, routed via dispatch, response written to
stdout. Continues until stdin closes.

**Concurrency**: Main loop is single-threaded. Poller runs separate
goroutine for file change notifications. Thread-safe stdout writes
via mutex-protected `mcpIO.Writer`.

**Sub-packages**:

### server/dispatch
Routes by method name to specialized handlers. Falls back to
ErrCodeNotFound for unknown methods.

### server/catalog
URI-to-file resource mapping. 9 resources: 8 individual context
files + 1 assembled agent packet (`ctx://context/agent`).
`Init()` builds lookup map once; `ToList()` returns immutable list.

### server/poll
File mtime-based polling (5s interval). Lazy goroutine lifecycle:
starts on first Subscribe(), stops when all unsubscribed.
Emits `notifications/resources/updated` via callback.

### server/route/*
Method-specific handlers:
- `initialize/`: handshake with capability advertisement
- `ping/`: simple pong
- `fallback/`: unknown method error
- `tool/`: tool invocation router + governance warning append
- `prompt/`: prompt rendering router

### server/def/*
Static definitions:
- `def/tool/`: 11 tool definitions with JSON Schema
- `def/prompt/`: 5 prompt definitions with arguments

### server/extract
MCP argument extraction: `EntryArgs(args)` for required fields,
`Opts(args)` for optional entry attributes.

### server/io
Thread-safe JSON writer: `WriteJSON(v)` marshals, appends newline,
writes atomically under mutex.

### server/out
Response builders: `OkResponse()`, `ErrResponse()`, `ToolOK()`,
`ToolError()`, `ToolResult()`, `Call()`.

### server/parse
`Request(data)` unmarshals raw JSON to proto.Request. Returns
(nil, nil) for notifications; (nil, error) for malformed JSON.

### server/stat
Lightweight analytics: `TotalAdds(m)` sums entry add counts.

**Edge cases**:
- Parse errors return JSON-RPC error, loop continues
- Notifications (no ID) produce no response
- Scanner buffer is configurable for large payloads

**Performance considerations**:
- Single-threaded request processing — no concurrent tool calls
- Poller checks every 5s regardless of subscription count
- Resource list built once at startup, never recomputed

**Danger zones**:
1. Single-threaded main loop — a slow handler blocks all requests.
   No timeout on handler execution.
2. Poller uses file mtime — sub-second changes may be missed.
   Rapid writes between polls are coalesced.
3. Scanner buffer size is fixed at startup — payloads exceeding
   it cause silent truncation and parse errors.

**Extension points**:
- Add new tools: define in def/tool/, add handler method, add
  route in tool/tool.go dispatch switch
- Add new resources: add to catalog/data.go, add read handler
- Add new prompts: define in def/prompt/, add builder in
  prompt/prompt.go

**Improvement ideas**:
- Add request timeout to prevent handler hangs
- Consider concurrent tool execution for read-only tools
- Resource change detection could use fsnotify instead of polling

**Dependencies**: handler, proto, session, config/mcp/*

---

## mcp/handler

**Purpose**: Domain logic implementation, testable without JSON-RPC
coupling. All tool and prompt functionality lives here.

**Key types**:
```
Handler {
    ContextDir  string
    TokenBudget int
    Session     *session.State
}
```

**Exported API**:
- `Status()`: context health summary
- `Add(type, content, opts)`: validate boundary, write entry
- `Complete(query)`: mark task done by number/text match
- `Drift()`: detect violations/warnings
- `Recall(limit, since)`: query session history
- `WatchUpdate(type, content, opts)`: write + queue pending update
- `Compact(archive)`: move completed tasks to archive
- `Next()`: next pending task
- `CheckTaskCompletion(recentAction)`: match action to tasks
- `SessionEvent(eventType, caller)`: start/end lifecycle
- `Remind()`: list pending reminders

### handler/task
Task list parsing for MCP: `ForEachPending(lines, fn)` iterates
pending tasks. `ContainsOverlap(action, taskText)` matches by
word-set intersection (>= 2 significant words).

**Danger zones**:
1. Add() performs file I/O in the handler — no transaction
   semantics. Partial writes on failure leave inconsistent state.
2. Complete() by text match is fuzzy — ambiguous task text can
   match the wrong task.
3. CheckTaskCompletion word overlap threshold (2 words) is low —
   false positives on common words.

**Dependencies**: context/load, entry, tidy, drift, journal/parser,
entity, io, rc

---

## mcp/session

**Purpose**: Per-session advisory state and governance warnings.

**Key types**:
```
State {
    ToolCalls        int
    AddsPerformed    map[string]int
    PendingFlush     []PendingUpdate
    sessionStarted   bool
    contextLoaded    bool
    lastDriftCheck   time.Time
    callsSinceWrite  int
}
```

**Governance warnings** (appended to tool responses):
1. Session not started (if sessionStarted=false)
2. Context not loaded (if contextLoaded=false)
3. Drift not checked (after interval or min calls)
4. Persist nudge (after callsSinceWrite threshold)
5. Violations from extension (reads violations.json)

**Data flow**: Each tool call -> RecordToolCall() ->
CheckGovernance() -> warnings appended to response text.

**Edge cases**:
- violations.json is read-and-cleared (one-shot delivery)
- Governance is advisory only — never blocks tool execution

**Danger zones**:
1. Session state is in-memory only — server restart loses all
   tracking. No persistence across MCP reconnections.
2. Governance thresholds are compile-time constants in
   config/mcp/governance — not user-configurable.

**Extension points**:
- Add new governance rules in CheckGovernance()
- Violations file format is extensible

**Dependencies**: config/mcp/governance, proto

---

## Tools (11 total)

| Tool | Read-Only | Description |
|------|-----------|-------------|
| ctx_status | Yes | Context health summary |
| ctx_add | No | Add task/decision/learning/convention |
| ctx_complete | No | Mark task done (idempotent) |
| ctx_drift | Yes | Detect context violations |
| ctx_journal_source | Yes | Query session history |
| ctx_watch_update | No | Apply structured updates |
| ctx_compact | No | Archive completed tasks |
| ctx_next | Yes | Next pending task |
| ctx_check_task_completion | Yes | Match action to tasks |
| ctx_session_event | No | Signal session start/end |
| ctx_remind | Yes | List pending reminders |

## Resources (9 total)

| URI | Content |
|-----|---------|
| ctx://context/tasks | TASKS.md |
| ctx://context/decisions | DECISIONS.md |
| ctx://context/conventions | CONVENTIONS.md |
| ctx://context/constitution | CONSTITUTION.md |
| ctx://context/architecture | ARCHITECTURE.md |
| ctx://context/learnings | LEARNINGS.md |
| ctx://context/glossary | GLOSSARY.md |
| ctx://context/playbook | AGENT_PLAYBOOK.md |
| ctx://context/agent | Assembled packet (all files, token-budgeted) |

## Prompts (5 total)

| Prompt | Description |
|--------|-------------|
| ctx-session-start | Load full context at session start |
| ctx-add-decision | Format architectural decision entry |
| ctx-add-learning | Format learning entry |
| ctx-reflect | Guide end-of-session reflection |
| ctx-checkpoint | Report session statistics |
