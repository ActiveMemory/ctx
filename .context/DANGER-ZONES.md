# Danger Zones

_Generated 2026-04-03 from DETAILED_DESIGN module analysis.
Enriched 2026-06-09 via GitNexus (index @ 60d8e823, 27,927 symbols,
203 flows; supersedes 2026-04-03 enrichment)._

## Summary

| Module | Zone | Risk | d=1 | Flows | Why |
|--------|------|------|-----|-------|-----|
| assets/read/desc | desc.Text() blast radius | CRITICAL | 940 | 183 | Single highest-connectivity symbol in codebase |
| io | SafeWriteFile blast radius | CRITICAL | 95 | 15 | Every file write in the system routes through this |
| context/load | load.Do() context hub | CRITICAL | 20 | 28 | Every context consumer (CLI + MCP) routes through this |
| config/embed/text | DescKey-YAML sync | CRITICAL | n/a | 183 | Missing key = empty output everywhere desc.Text() reaches |
| rc | RC() accessor surface | CRITICAL | 34 | 5+ | 26 in-package accessors + 8 external modules read config through it |
| memory | DiscoverPath coupling | HIGH | 4 | 4 | 4 modules depend; slug format is undocumented |
| config/file | FileReadOrder | HIGH | n/a | 100+ | Reordering changes what agents see first |
| assets/read/lookup | Init() ordering | HIGH | n/a | n/a | desc.Text() before Init() = silent empty strings |
| entry | Concurrent writes | HIGH | n/a | n/a | No locking on read-modify-write |
| journal/parser | JSONL format dependency | HIGH | n/a | n/a | Undocumented format; schema changes break silently |
| memory | External file writes | HIGH | n/a | n/a | Only package that modifies files outside .context/ |
| mcp/server | Single-threaded main loop | HIGH | 1 | 6+ | Slow handler blocks all requests; no timeout |
| system | 34 hook subcommands | HIGH | n/a | n/a | Hook behavior changes affect all agent integrations |
| config/regex | Pattern changes | MEDIUM | n/a | n/a | Silent behavior change across all consumers |
| assets/tpl | Sprintf templates | MEDIUM | n/a | n/a | Mismatched placeholders = runtime panic |
| io | Path validation timing | MEDIUM | n/a | n/a | Stale validation if root changes post-check |
| entry | Index update failure | MEDIUM | n/a | n/a | Write succeeds but index left stale |
| journal/parser | 1MB buffer limit | MEDIUM | n/a | n/a | Large tool results truncated without warning |
| memory | Slug format dependency | MEDIUM | 4 | 4 | Claude Code naming convention change breaks discovery |
| drift | Path ref false positives | MEDIUM | n/a | n/a | Code examples in markdown trigger false warnings |
| mcp/server/dispatch/poll | Mtime granularity | MEDIUM | n/a | n/a | Sub-second changes between polls are missed |
| entity.MCPSession | In-memory only state | MEDIUM | n/a | n/a | Server restart loses governance tracking |
| mcp/handler | Fuzzy task matching | MEDIUM | n/a | n/a | Word overlap threshold (2) causes false positives |
| rc | sync.Once lock-in | MEDIUM | n/a | n/a | First RC() call locks config for process lifetime |
| bootstrap | PersistentPreRunE | MEDIUM | n/a | n/a | New commands without SkipInit fail pre-.context/ |
| tidy | Indentation sensitivity | MEDIUM | n/a | n/a | Tab/space inconsistency = wrong block boundaries |
| trace | Stale staged refs | MEDIUM | n/a | n/a | Git index changes between collect and commit |

## By Module

### internal/assets/read/desc (enriched 2026-06-09 via GitNexus)

1. **desc.Text() blast radius** - 940 direct callers spanning every
   layer of the codebase; the new index resolves call sites the
   2026-04 index missed (previously reported as "30+"). Top
   affected modules: Skill (157), Initialize (62), Journal (59),
   Pad (42), Memory (36), Trigger (34), Nudge (31), Steering (31).
   58 affected process groups totaling 183 execution-flow hits.
   - Blast radius: d=1: 940, flows: 183 (58 process groups), modules: 20
   - Risk: CRITICAL (enriched 2026-06-09 via GitNexus)
   - Modification advice: treat as a frozen API. Any signature
     change cascades through ~940 call sites. Add new functions
     rather than modifying existing ones.

### internal/io (enriched 2026-06-09 via GitNexus)

1. **SafeWriteFile blast radius** - 95 direct callers across the
   entire codebase (up from 69 in the 2026-04 enrichment):
   entry.Write, index.Reindex, journal state.Save, crypto.SaveKey,
   tidy.WriteArchive, memory (Sync, Publish, Archive, SaveState),
   initialize/* (templates, vscode, plugin, kb, backup, merge),
   system/* (persistence, counter, heartbeat, load gate), all
   setup/* deployers (now including opencode and copilotcli), pad
   store/history, trace hooks, task archive/complete/snapshot,
   compact, config profile — plus the newer kb (scaffold, reindex,
   source-coverage), hub (persist, daemon, admin), steering,
   trigger, connection (config, sync state), handover/closeout
   writers, and ctxctl audit store.
   - Blast radius: d=1: 95, processes: 15, modules: 20
   - Risk: CRITICAL (enriched 2026-06-09 via GitNexus)
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
   output everywhere that key is used. Since desc.Text() has 940
   direct callers across 183 execution-flow hits (verified
   2026-06-09), a missing YAML entry creates invisible missing
   text across the entire system.
   - Risk: CRITICAL (re-verified 2026-06-09 via GitNexus)
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
   to the same file can lose data. ValidateAndWrite blast radius is
   small (d=1: 2 — MCP handler Add and WatchUpdate; d=2: tool route
   add/watchUpdate; 6 flow hits via DispatchCall), so the risk is
   operational (race), not structural. (enriched 2026-06-09 via
   GitNexus)
   ⚠ Possible undercount - CLI add path is a known caller from
     reading but does not appear at d=1 in the graph
   - Risk: HIGH
   - Modification advice: consider adding file-level locking for
     write operations; current risk is low (single-user tool)

2. **Index update after write** - If index update fails after
   successful entry write, the entry exists but the index table
   is stale. No rollback mechanism.
   - Risk: MEDIUM

### internal/memory (enriched 2026-06-09 via GitNexus)

1. **DiscoverPath coupling** - 4 direct callers (checkmemorydrift
   hook, memory/core/resolve.DiscoverSource, memory status, memory
   diff) plus 4 indirect dependents at d=2 (memory sync, publish,
   importer, unpublish — all route through resolve.DiscoverSource).
   All memory subcommands + the drift hook still depend on this;
   the call graph now funnels most of them through one resolver.
   4 affected modules (Memory, Publish, Nudge, Ctximport), 4 flows.
   - Blast radius: d=1: 4, d=2: 4, flows: 4, modules: 4
   - Risk: HIGH (enriched 2026-06-09 via GitNexus; was CRITICAL —
     blast radius is single-domain, but slug format remains
     undocumented upstream)
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

### internal/mcp/server (enriched 2026-06-09 via GitNexus)

1. **Single-threaded main loop** - server.New() calls catalog.Init()
   (d=1: 1 — server.New; d=2: 1 — mcp/cmd Cmd). Handler execution
   has no timeout. A blocking handler freezes all tools.
   - Blast radius: catalog.Init d=1: 1, d=2: 1; Serve upstream not
     resolved by the graph (method-call edge gap) — entry point is
     mcp/cmd
   ⚠ Possible undercount - Serve is a top-level entry point; the
     operational risk is in everything downstream of it, not its
     callers
   - Risk: HIGH (operational, not blast-radius; enriched 2026-06-09
     via GitNexus)
   - Modification advice: add context.WithTimeout to handler calls

2. **Poller mtime granularity** - 5s interval. Sub-second changes
   between polls are coalesced.
   - Risk: MEDIUM

### entity.MCPSession / mcp/handler CheckGovernance

1. **In-memory only state** - Session governance lost on restart.
   `handler.CheckGovernance` has clean call chain (d=1: 1 ->
   appendGovernance -> DispatchCall -> dispatch.Do) but the advisory
   data it tracks is ephemeral. Data lives in `entity.MCPSession`;
   the I/O-touching CheckGovernance free function lives in
   `mcp/handler` (now `governance.go`) because it drains
   `.context/state/violations.json`. 6 execution-flow hits via
   DispatchCall.
   - Blast radius: d=1: 1, d=2: 1, d=3: 1 (clean chain), flows: 6
   - Risk: MEDIUM (enriched 2026-06-09 via GitNexus)

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

### internal/context/load (enriched 2026-06-09 via GitNexus)

1. **load.Do() is the context hub** - ~20 non-test direct callers:
   CLI commands (status, sync, load, drift, compact, agent),
   doctor checks (Drift, ContextTokenSize), MCP handlers (Status,
   Drift), MCP resource DispatchRead, MCP prompt sessionStart,
   rc.RC, and assets/hooks/messages Registry. Participates in 28
   auto-detected execution flows. Calls: validate.Symlinks,
   rc.ContextDir, io.SafeReadUserFile, token.Estimate,
   summary.Generate, sanitize.EffectivelyEmpty, err/context.NotFound.
   - Blast radius: d=1: ~20 non-test (via context query), flows: 28
   ⚠ Possible undercount - impact() returned 0 callers with
     partial=true on the fresh index; caller counts above come from
     the context() incoming-call list instead
   - Risk: CRITICAL (enriched 2026-06-09 via GitNexus)
   - Modification advice: any change to load behavior (file
     filtering, sort order, error handling) affects every context
     consumer. The return type (entity.Context) is a shared
     contract across CLI and MCP.

### internal/rc (enriched 2026-06-09 via GitNexus)

1. **RC() blast radius** - 34 direct callers: 26 in-package
   accessor wrappers (TokenBudget, PriorityOrder, HooksEnabled,
   SteeringDir, ClassifyRules, ...) plus external callers in
   system/session (EffectiveContextWindow), config profile Detect,
   and others across 8 modules (Rc, Drift, Nudge, Token, Skill,
   Switchcmd, Memory, Add). Combined with the sync.Once lock-in
   zone below, any change to RC() initialization semantics
   propagates through every config read in the process.
   - Blast radius: d=1: 34, modules: 8, processes: 5
   - Risk: CRITICAL (enriched 2026-06-09 via GitNexus)
   - Modification advice: never add I/O or fallible logic to the
     accessor path; changes to load order or defaults are
     architectural decisions (record in DECISIONS.md).
