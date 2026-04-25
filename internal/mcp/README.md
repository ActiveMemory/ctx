# internal/mcp: MCP Server

JSON-RPC 2.0 server exposing ctx context to any MCP-compatible
AI tool over stdin/stdout. See `doc.go` for the full resource,
tool, and prompt catalog.

## Package Map

```
mcp/
  doc.go              Package docs (resource/tool/prompt catalog)
  proto/              JSON-RPC 2.0 message types, error codes
  handler/            Domain logic (testable without JSON-RPC)
    handler.go        Handler type, constructor
    tool.go           Tool method implementations (Status, Add, etc.)
    types.go          EntryOpts
    task/             Task list parsing for MCP (ForEachPending)
  server/             Protocol layer
    server.go         Main loop: stdin → parse → dispatch → stdout
    dispatch/         Method-based request routing
    catalog/          URI-to-file resource mapping (9 resources)
    poll/             File mtime polling for change notifications
    io/               Thread-safe JSON writer (mutex-protected)
    out/              Response builders (OkResponse, ErrResponse, etc.)
    parse/            JSON-RPC request parsing
    extract/          MCP argument extraction
    stat/             Lightweight analytics
    route/            Per-method handlers
      initialize/     Handshake with capability advertisement
      ping/           Pong
      fallback/       Unknown method → ErrCodeNotFound
      tool/           Tool invocation + governance warning append
      prompt/         Prompt rendering
    def/              Static definitions
      tool/           11 tool definitions with JSON Schema
      prompt/         5 prompt definitions with arguments
  session/            Per-session advisory state
    state.go          Tool call tracking, pending updates
    governance.go     Advisory warnings (drift, persist, etc.)
    violations.go     Extension integration via violations.json
```

## How To Add a New Tool

Three files, always:

1. **Define** in `server/def/tool/tool.go`: add entry to `Defs`
   array with name, description, and `InputSchema` (JSON Schema
   for parameters)

2. **Implement** in `handler/tool.go`: add method on `Handler`
   with signature `func (h *Handler) ToolName(args...) (string, error)`

3. **Route** in `server/route/tool/tool.go`: add case in the
   dispatch switch calling your handler method, wrap result with
   `out.ToolResult()`

## How To Add a New Prompt

Same pattern, three files:

1. **Define** in `server/def/prompt/prompt.go`: add entry to
   `Defs` array with name, description, and arguments

2. **Build** in `server/route/prompt/prompt.go`: add builder
   function returning `[]proto.PromptMessage`

3. **Route** in `server/route/prompt/dispatch.go`: add case in
   the dispatch switch

## How To Add a New Resource

1. **Register** in `server/catalog/data.go`: add URI-to-file
   mapping

2. **Handle** in `server/resource/resource.go`: if it needs
   special assembly (like the agent packet), add a reader function

## Key Design Decisions

- **handler/ has no JSON-RPC coupling**: all tool methods take
  typed args and return `(string, error)`. Protocol translation
  happens in server/route/. This makes handler/ testable without
  stdin/stdout.

- **Single-threaded main loop**: one request at a time. Poller
  runs in a background goroutine. Thread safety via mutex on
  stdout writer only.

- **Governance is advisory**: session state tracks tool calls and
  nudges (drift check, persist reminder) but never blocks execution.

- **Protocol version**: 2024-11-05. Capabilities advertised:
  resources (subscribe=true), tools, prompts.
