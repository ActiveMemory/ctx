# Extension Points

_Generated 2026-04-03. Enriched via GitNexus call graph analysis._

## Summary

| Pattern | Key Symbol | Domain | Notes |
|---------|-----------|--------|-------|
| Session parser | SessionParser interface | journal/parser | 4 registered formats |
| CLI command | bootstrap.Initialize() | bootstrap | 34 commands in 9 groups |
| MCP tool | def/tool.Defs | mcp/server | 11 tools defined |
| MCP prompt | def/prompt.Defs | mcp/server | 5 prompts defined |
| MCP resource | catalog.Init() | mcp/server | 9 resources mapped |
| File I/O guard | io.Safe* | io | All file ops route through |
| Config constants | config/* sub-packages | config | 60+ domain packages |
| Output writer | write/* packages | write | 46 command-specific writers |
| Error constructor | err/* packages | err | 35 domain-specific packages |
| Asset reader | assets/read/* | assets | 14 typed accessor packages |
| Drift check | drift.Detect() | drift | 7 pluggable checks |
| Agent setup | setup/core/* | cli/setup | 5 tool-specific deployers |
| Entry type | entry.Validate() | entry | Type-specific validation rules |
| Exec wrapper | exec/* | exec | 5 command wrappers |

## By Pattern

### Session Parser Registration

Registration: `journal/parser` auto-detection via `Matches(path)`

Registered implementations:
1. `ClaudeCode` - `internal/journal/parser/claude.go`
2. `Copilot` - `internal/journal/parser/copilot.go`
3. `CopilotCLI` - `internal/journal/parser/copilot_cli.go`
4. `MarkdownSession` - `internal/journal/parser/markdown.go`

How to extend: implement SessionParser interface with `Matches()`
and `ParseFile()` methods. New parsers are auto-detected.

### CLI Command Registration

Registration: `internal/bootstrap/group.go` functions return
`[]registration` structs.

9 group functions: gettingStarted(), contextCmds(), artifacts(),
sessions(), runtimeCmds(), integrations(), diagnostics(),
utilities(), hiddenCmds()

How to extend: create new cli/ package following cmd/root + core/
taxonomy, add registration in the appropriate group function.

### MCP Tool Definitions

Registration: `internal/mcp/server/def/tool/tool.go` Defs array

11 tools registered. Each tool has InputSchema (JSON Schema),
handler method in `mcp/handler/tool.go`, and route in
`mcp/server/route/tool/tool.go` dispatch switch.

How to extend: add definition to Defs, add handler method, add
case in route dispatch switch. Three-file change.

### MCP Prompt Definitions

Registration: `internal/mcp/server/def/prompt/prompt.go` Defs array

5 prompts registered. Each prompt has arguments and a builder
function in `mcp/server/route/prompt/prompt.go`.

How to extend: add definition to Defs, add builder function, add
case in route dispatch. Three-file change.

### Agent Setup Deployers

Registration: `internal/cli/setup/core/*` packages

5 tool-specific deployers:
1. `agents/` - AGENTS.md deployment
2. `copilot/` - GitHub Copilot (instructions + VS Code MCP)
3. `copilot_cli/` - Copilot CLI (instructions, skills, agent, MCP)
4. (Claude Code via initialize, not setup)
5. (Cursor/Aider/Windsurf via simpler paths)

How to extend: create new `setup/core/<tool>/` package with
Deploy() function. Add case in setup command's Run() handler.

### Drift Checks

Registration: `internal/drift/detector.go` Detect() function

7 checks executed in sequence. Each check is a function called
within Detect().

How to extend: add new check function, call it within Detect().
Single-file change.
