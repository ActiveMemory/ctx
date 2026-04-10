# Danger Zones

_Generated 2026-04-03 from DETAILED_DESIGN module analysis.
Enriched 2026-04-03 via GitNexus (blast radius verified)._

## Summary

| Module | Zone | Risk | d=1 | Flows | Why |
|--------|------|------|-----|-------|-----|
| assets/read/desc | desc.Text() blast radius | CRITICAL | 30+ | 53 | Single highest-connectivity symbol in codebase |
| io | SafeWriteFile blast radius | CRITICAL | 69 | n/a | Every file write in the system routes through this |
| config/embed/text | DescKey-YAML sync | CRITICAL | n/a | 53 | Missing key = empty output in all 53 flows |
| memory | DiscoverPath coupling | CRITICAL | 7 | 7 | 4 modules depend; slug format is undocumented |
| config/file | FileReadOrder | HIGH | n/a | 100+ | Reordering changes what agents see first |
| assets/read/lookup | Init() ordering | HIGH | n/a | n/a | desc.Text() before Init() = silent empty strings |
| entry | Concurrent writes | HIGH | n/a | n/a | No locking on read-modify-write |
| journal/parser | JSONL format dependency | HIGH | n/a | n/a | Undocumented format; schema changes break silently |
| memory | External file writes | HIGH | n/a | n/a | Only package that modifies files outside .context/ |
| mcp/server | Single-threaded main loop | HIGH | 1 | 20 | Slow handler blocks all requests; no timeout |
| system | 34 hook subcommands | HIGH | n/a | n/a | Hook behavior changes affect all agent integrations |
| config/regex | Pattern changes | MEDIUM | n/a | n/a | Silent behavior change across all consumers |
| assets/tpl | Sprintf templates | MEDIUM | n/a | n/a | Mismatched placeholders = runtime panic |
| io | Path validation timing | MEDIUM | n/a | n/a | Stale validation if root changes post-check |
| entry | Index update failure | MEDIUM | n/a | n/a | Write succeeds but index left stale |
| journal/parser | 1MB buffer limit | MEDIUM | n/a | n/a | Large tool results truncated without warning |
| memory | Slug format dependency | MEDIUM | 7 | 7 | Claude Code naming convention change breaks discovery |
| drift | Path ref false positives | MEDIUM | n/a | n/a | Code examples in markdown trigger false warnings |
| mcp/server/dispatch/poll | Mtime granularity | MEDIUM | n/a | n/a | Sub-second changes between polls are missed |
| entity.MCPSession | In-memory only state | MEDIUM | n/a | n/a | Server restart loses governance tracking |
| mcp/handler | Fuzzy task matching | MEDIUM | n/a | n/a | Word overlap threshold (2) causes false positives |
| rc | sync.Once lock-in | MEDIUM | n/a | n/a | First RC() call locks config for process lifetime |
| bootstrap | PersistentPreRunE | MEDIUM | n/a | n/a | New commands without SkipInit fail pre-.context/ |
| tidy | Indentation sensitivity | MEDIUM | n/a | n/a | Tab/space inconsistency = wrong block boundaries |
| trace | Stale staged refs | MEDIUM | n/a | n/a | Git index changes between collect and commit |

## By Module

### internal/assets/read/desc (enriched 2026-04-03 via GitNexus)

1. **desc.Text() blast radius** - 30+ direct callers spanning every
   layer: MCP handler (all 11 tool methods), format (TimeAgo,
   Duration, Tokens), index (GenerateTable, UpdateDecisions,
   UpdateLearnings), tidy, trace, memory, sysinfo, io (SafeFprintf),
   mcp/handler (CheckGovernance), mcp/server (Serve).
   Participates in 53 execution flows.
   - Blast radius: d=1: 30+, flows: 53
   - Risk: CRITICAL (enriched 2026-04-03 via GitNexus)
   - Modification advice: treat as a frozen API. Any signature
     change cascades through 30+ call sites and 53 flows. Add
     new functions rather than modifying existing ones.

### internal/io (enriched 2026-04-03 via GitNexus)

1. **SafeWriteFile blast radius** - 69 direct callers across the
   entire codebase: entry.Write, index.Reindex, journal state.Save,
   crypto.SaveKey, tidy.WriteArchive, memory (Sync, Publish,
   Archive, SaveState), 20+ initialize/* functions, 10+ system/*
   functions, all setup/* deploy functions, pad store, trace hooks,
   task archive/complete/snapshot, compact, config profile, and more.
   - Blast radius: d=1: 69
   - Risk: CRITICAL (enriched 2026-04-03 via GitNexus)
   - Modification advice: any change to SafeWriteFile semantics
     (validation rules, error handling, permissions) affects every
     write operation in the system. Test exhaustively.

2. **Path validation timing** - Path validation relies on resolved
   prefix matching. If the project root changes after validation,
   the check is stale.
   - Risk: MEDIUM
   - Modification advice: re-validate on use, not on construction

### internal/config/*

1. **config/embed/text DescKey-YAML sync** - Adding a DescKey
   constant without a corresponding YAML entry produces empty
   output everywhere that key is used. Since desc.Text() participates
   in 53 execution flows, a missing YAML entry creates invisible
   missing text across the entire system.
   - Risk: CRITICAL (upgraded from HIGH based on desc.Text() flow count)
   - Modification advice: always run TestDescKeyYAMLLinkage audit
     after adding/removing DescKey constants

2. **config/regex pattern changes** - Compiled regex patterns are
   consumed by every layer. Changing a pattern silently affects
   all match sites. No type safety on capture group indices.
   - Risk: MEDIUM
   - Modification advice: grep for all import sites of the specific
     regex sub-file before changing patterns

3. **config/file FileReadOrder** - This array determines context
   priority for all agents. context/load.Do() participates in 100+
   execution flows. Reordering changes what every AI agent sees
   first when context is loaded or budgeted.
   - Risk: HIGH
   - Modification advice: treat as an architectural decision; update
     DECISIONS.md and notify users

### internal/assets

1. **assets/read/lookup Init() ordering** - desc.Text() returns
   empty strings if called before lookup.Init(). Silent failure.
   No warning, no panic.
   - Risk: HIGH
   - Modification advice: Init() is called in bootstrap; ensure
     any new code paths that bypass bootstrap also call Init()

2. **assets/tpl Sprintf templates** - Format strings with %s/%d
   placeholders. Mismatched arg count = runtime panic. No
   compile-time checking.
   - Risk: MEDIUM
   - Modification advice: check all callers when modifying templates;
     migration to text/template is tracked in TASKS.md

### internal/entry

1. **Concurrent writes** - Read-modify-write to context files
   without file locking. Two concurrent callers (CLI + MCP) writing
   to the same file can lose data.
   - Risk: HIGH
   - Modification advice: consider adding file-level locking for
     write operations; current risk is low (single-user tool)

2. **Index update after write** - If index update fails after
   successful entry write, the entry exists but the index table
   is stale. No rollback mechanism.
   - Risk: MEDIUM

### internal/memory (enriched 2026-04-03 via GitNexus)

1. **DiscoverPath coupling** - 7 direct callers across 4 modules
   (Memory, Publish, Session, Ctximport). All 6 memory subcommands
   + check_memory_drift hook depend on this. 7 execution flows.
   - Blast radius: d=1: 7, flows: 7, modules: 4
   - Risk: CRITICAL (enriched 2026-04-03 via GitNexus)
   - Modification advice: slug format change breaks all memory
     operations. Add fallback discovery; abstract into agent-keyed
     registry for multi-agent support.

2. **External file writes** - MergePublished() writes to MEMORY.md
   outside .context/. This is the only package that modifies
   external state, bypassing boundary validation.
   - Risk: HIGH

### internal/journal/parser

1. **JSONL format dependency** - Claude Code's session format is
   reverse-engineered, not documented. Any upstream schema change
   breaks import silently.
   - Risk: HIGH

2. **1MB buffer limit** - Sessions with very large tool results
   are silently truncated at the scanner buffer boundary.
   - Risk: MEDIUM

### internal/mcp/server (enriched 2026-04-03 via GitNexus)

1. **Single-threaded main loop** - server.New() calls catalog.Init()
   which feeds 20 execution flows. Handler execution has no timeout.
   A blocking handler freezes all 11 tools.
   - Blast radius: Serve d=1: 1 (mcp/cmd), but Init d=1: 1 affects 20 flows
   - Risk: HIGH (operational, not blast-radius)
   - Modification advice: add context.WithTimeout to handler calls

2. **Poller mtime granularity** - 5s interval. Sub-second changes
   between polls are coalesced.
   - Risk: MEDIUM

### entity.MCPSession / mcp/handler CheckGovernance

1. **In-memory only state** - Session governance lost on restart.
   `handler.CheckGovernance` has clean call chain (d=1: 1 ->
   appendGovernance -> DispatchCall -> Do) but the advisory data
   it tracks is ephemeral. Data lives in `entity.MCPSession`;
   the I/O-touching CheckGovernance free function lives in
   `mcp/handler` because it drains `.context/state/violations.json`.
   - Blast radius: d=1: 1, d=2: 1, d=3: 1 (clean chain)
   - Risk: MEDIUM (enriched 2026-04-03 via GitNexus)

### internal/system

1. **34 hook subcommands** - Hidden plumbing commands that agent
   integrations depend on. Behavior changes affect all connected
   agents silently.
   - Risk: HIGH
   - Modification advice: treat hook commands as public API

### internal/tidy

1. **Indentation sensitivity** - Block boundary detection uses
   indentation. Tab/space inconsistency = wrong block boundaries.
   - Risk: MEDIUM

### internal/trace

1. **Stale staged refs** - Git index changes between collect and
   commit may cause refs to be stale.
   - Risk: MEDIUM

### internal/context/load (enriched 2026-04-03 via GitNexus)

1. **load.Do() is the context hub** - 12+ non-test callers: CLI
   commands (status, sync, load, drift, doctor), MCP handlers
   (Status, Drift, Compact, Next, CheckTaskCompletion), MCP
   resource/prompt handlers. Participates in 100+ execution flows.
   Calls: validate.Symlinks, rc.ContextDir, io.SafeReadUserFile,
   token.Estimate, summary.Generate, sanitize.EffectivelyEmpty.
   - Blast radius: d=1: 12+ non-test, flows: 100+
   - Risk: CRITICAL (enriched 2026-04-03 via GitNexus)
   - Modification advice: any change to load behavior (file
     filtering, sort order, error handling) affects every context
     consumer. The return type (entity.Context) is a shared
     contract across CLI and MCP.
