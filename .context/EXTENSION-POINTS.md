# Extension Points

_Generated 2026-04-03. Enriched 2026-06-09 via GitNexus call graph
analysis + source verification (index @ 60d8e823)._

## Summary

| Pattern | Key Symbol | Domain | Count | Notes |
|---------|-----------|--------|-------|-------|
| Session parser | registeredParsers | journal/parser | 4 | Auto-detected via Matches() |
| CLI command | bootstrap group funcs | bootstrap | 42 | 40 grouped + 2 hidden, 8 group functions |
| MCP tool | def/tool.Defs | mcp/server | 15 | Up from 11 (added SessionEvent, SteeringGet, Search, SessionStart, SessionEnd...) |
| MCP prompt | def/prompt.Defs | mcp/server | 5 | Unchanged |
| MCP resource | catalog table | mcp/server | 8 | Per-file mappings; URIs built by catalog.Init() |
| Drift check | drift.Detect() | drift | 12 | Up from 7 |
| Agent setup | setup/core/* | cli/setup | 8 | Up from 5 (added cline, cursor, kiro, opencode) |
| File I/O guard | io.Safe* | io | n/a | 95 SafeWriteFile callers (see DANGER-ZONES.md) |
| Config constants | config/* sub-packages | config | 60+ | Pattern unchanged |
| Output writer | write/* packages | write | 46+ | Pattern unchanged; now includes handover, closeout, kb writers |
| Error constructor | err/* packages | err | 35+ | Pattern unchanged |
| Asset reader | assets/read/* | assets | 14 | Pattern unchanged |
| Entry type | entry.Validate() | entry | n/a | Type-specific validation rules |
| Exec wrapper | exec/* | exec | 5 | Pattern unchanged |

## By Pattern

### Session Parser Registration

Registration: `registeredParsers` slice in
`internal/journal/parser/parser.go:22` — auto-detection via
`Matches(path)`.

Registered implementations (verified 2026-06-09):
1. `NewClaudeCode()` - `internal/journal/parser/parser.go:23`
2. `NewCopilot()` - `internal/journal/parser/parser.go:24`
3. `NewCopilotCLI()` - `internal/journal/parser/parser.go:25`
4. `NewMarkdownSession()` - `internal/journal/parser/parser.go:26`

How to extend: implement the Session interface with `Matches()`
and `ParseFile()`, append to `registeredParsers`.

### CLI Command Registration

Registration: `internal/bootstrap/group.go` functions return
`[]registration` structs.

8 group functions (verified 2026-06-09, 42 total commands):
- `gettingStarted()` - 3 commands (`group.go:59`)
- `contextCmds()` - 7 commands (`group.go:79`)
- `artifacts()` - 7 commands (`group.go:102`) — now includes kb, handover
- `sessions()` - 4 commands (`group.go:118`)
- `runtimeCmds()` - 4 commands (`group.go:132`)
- `integrations()` - 9 commands (`group.go:150`)
- `diagnostics()` - 6 commands (`group.go:168`)
- `hiddenCmds()` - 2 commands: site, system (`group.go:190`)

How to extend: create new cli/ package following cmd/root + core/
taxonomy, add registration in the appropriate group function.

### MCP Tool Definitions

Registration: `internal/mcp/server/def/tool/tool.go:29` Defs()
function.

15 tools registered (verified 2026-06-09; dispatch switch in
`internal/mcp/server/route/tool/dispatch.go` has matching 15
cases): Status (`tool.go:32`), Add (`:39`), Complete (`:77`),
Drift (`:94`), JournalSource (`:101`), WatchUpdate (`:122`),
Compact (`:145`), Next (`:161`), CheckTaskCompletion (`:168`),
SessionEvent (`:184`), Remind (`:206`), SteeringGet (`:213`),
Search (`:229`), SessionStart (`:246`), SessionEnd (`:253`).

How to extend: add definition to Defs, add handler method in
`mcp/handler`, add case in route dispatch switch. Three-file
change.

### MCP Prompt Definitions

Registration: `internal/mcp/server/def/prompt/prompt.go` Defs.

5 prompts registered (verified 2026-06-09): SessionStart
(`prompt.go:21`), AddDecision (`:26`), AddLearning (`:53`),
Reflect (`:80`), Checkpoint (`:85`).

How to extend: add definition to Defs, add builder function in
`mcp/server/route/prompt/prompt.go`, add case in route dispatch.

### MCP Resource Catalog

Registration: `table` in `internal/mcp/server/catalog/data.go:17`;
URIs built by `catalog.Init()` (`internal/mcp/server/catalog/attr.go`).

8 per-file resources (verified 2026-06-09): Constitution, Tasks,
Conventions, Architecture, Decisions, Learnings, Glossary,
Playbook. catalog.Init() blast radius: d=1: 1 (server.New),
d=2: 1 (mcp/cmd).

How to extend: add a mapping to the table (file constant, URI
constant, DescKey description) — single-file change plus config
constants.

### Drift Checks

Registration: `internal/drift/detector.go` Detect() function,
checks called in sequence at `detector.go:47-80`.

12 checks (verified 2026-06-09; was 7 at last enrichment):
1. `checkPathReferences` - `detector.go:47`
2. `checkStaleness` - `detector.go:50`
3. `checkConstitution` - `detector.go:53`
4. `checkRequiredFiles` - `detector.go:56`
5. `checkFileAge` - `detector.go:59`
6. `checkEntryCount` - `detector.go:62`
7. `checkMissingPackages` - `detector.go:65`
8. `checkTemplateHeaders` - `detector.go:68` (new)
9. `checkSteeringTools` - `detector.go:71` (new)
10. `checkHookPerms` - `detector.go:74` (new)
11. `checkSyncStaleness` - `detector.go:77` (new)
12. `checkRCTool` - `detector.go:80` (new)

How to extend: add new check function, call it within Detect().
Single-file change.

### Agent Setup Deployers

Registration: `internal/cli/setup/core/*` packages.

8 deployer packages (verified 2026-06-09; was 5):
1. `agents/` - AGENTS.md deployment
2. `cline/` - Cline (new)
3. `copilot/` - GitHub Copilot (instructions + VS Code MCP)
4. `copilotcli/` - Copilot CLI (instructions, skills, agent, MCP)
5. `cursor/` - Cursor (new)
6. `kiro/` - Kiro (new)
7. `mcp/` - generic MCP config deployment
8. `opencode/` - OpenCode (skills + plugin) (new)

How to extend: create new `setup/core/<tool>/` package with
Deploy() function. Add case in setup command's Run() handler.
