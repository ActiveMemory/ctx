# Detailed Design: Domain Layer

Modules: entity, entry, context/*, drift, index, task, tidy, trace,
journal/*, memory, notify, claude

## internal/entity

**Purpose**: Core domain types shared across CLI, MCP, and write
subsystems. Pure data — no I/O methods or business logic.

**Key types**:
- `Session`: reconstructed conversation (ID, tool source, CWD,
  timing, messages, tokens)
- `Message`: single message with role, text, thinking, tool uses
- `ToolUse`, `ToolResult`: tool invocation and response
- `Context`: loaded .context/ directory with files + token counts
- `FileInfo`: file metadata (name, path, size, mtime, content,
  isEmpty, tokens, summary)
- `TaskBlock`: task with nested content, completion and archival state
- `JournalEntry`: parsed journal file with frontmatter + metadata
- `JournalFrontmatter`: YAML structure (title, date, project,
  tokens, type, outcome, topics, keyfiles, summary)
- `HookInput`: JSON from Claude Code hook stdin
- `BlockResponse`: block/allow decision for hooks
- `EntryParams`: all fields for entry write operations
- `BootstrapOutput`: JSON output for bootstrap command

**Data flow**: Defined here, consumed everywhere. Entity types are
the shared vocabulary between CLI commands, MCP handlers, and
write packages.

**Edge cases**:
- TaskBlock.OlderThan() uses #done: timestamp extraction
- Session.UserMessages() and AssistantMessages() filter by role

**Danger zones**:
1. Adding fields to entity types affects serialization in MCP
   responses, JSON output, and YAML frontmatter simultaneously.
2. EntryParams is consumed by 3 callers (CLI add, MCP handler,
   watch) — field changes ripple widely.
3. JournalFrontmatter YAML tags must match journal state schema.

**Extension points**:
- Add new entity types for new domain concepts
- Session type is parser-agnostic — new parsers produce Sessions

**Dependencies**: config/* (constants), task (helper)

---

## internal/entry

**Purpose**: Domain API for adding entries to context files.
Validates required fields, formats per type, appends to file,
updates indices.

**Exported API**:
- `Validate(params, examplesFn)`: check required fields by type
  (Decision: context, rationale, consequence; Learning: context,
  lesson, application)
- `Write(params)`: format, append, update indices
- `ValidateAndWrite(params)`: orchestrate both

**Data flow**: `ValidateAndWrite()` -> `Validate()` -> `Write()` ->
reads existing file -> formats per type -> appends via insert ->
writes back -> updates Decisions/Learnings indices.

**Edge cases**:
- Tasks and conventions skip index update
- Empty content with --file flag reads from file path
- Stdin reading supported when no content or file specified

**Danger zones**:
1. Write() performs read-modify-write without locking — concurrent
   writes to the same file can lose data.
2. Index update failure after successful entry write leaves
   inconsistent state (entry added, index stale).

**Extension points**:
- New entry types need format template in tpl/ and validation
  rules in Validate()

**Dependencies**: cli/add/core/format, cli/add/core/insert, config/entry,
entity, err/add, index, io, rc

---

## internal/context/*

**Purpose**: Loads and manages .context/ files with token counting.

### context/load

**Exported API**: `Do(dir)` reads all .md files from .context/,
rejects symlinks (M-2 defense), calculates tokens per file and
totals. Returns `entity.Context`.

**Danger zones**:
1. Loads all files into memory — no streaming. Large .context/
   directories (many journal files) could use significant memory.

### context/resolve

**Exported API**: `JournalDir()`, `DirLine()`, `AppendDir(msg)`.
Path resolution helpers.

### context/sanitize

**Exported API**: `EffectivelyEmpty(content)`. Filters headers
and whitespace to detect placeholder files.

### context/summary

**Exported API**: `Generate(name, content)`. Creates brief summary
per file type (CONSTITUTION: counts invariants, TASKS: counts
active/completed, etc.).

### context/token

**Exported API**: `Estimate(content)`. Rough estimate using ~4
chars/token. Conservative overestimate (safer for budgeting).

### context/validate

**Exported API**: `Initialized(contextDir)` checks all required
files exist. `Exists(dir)` checks directory existence.

**Dependencies**: config/*, entity, err/context, io, rc, validate

---

## internal/drift

**Purpose**: Context quality validation with 7 checks.

**Exported API**: `Detect(context)` returns `Report` with
warnings and violations. `Report.Status()` computes overall
health (healthy/warning/error).

**Checks**: path refs, staleness, constitution compliance,
required files, file age, entry counts, missing packages.

**Edge cases**:
- Thresholds configurable via rc (age limits, count limits)
- Missing packages check uses go.mod/package.json detection

**Danger zones**:
1. Path ref check uses backtick-enclosed paths from markdown —
   false positives on code examples that happen to look like paths.
2. Staleness check compares git log timestamps — fails in repos
   with unusual date formats or clock skew.

**Extension points**:
- Add new drift checks by extending Detect() check list
- Thresholds are rc-configurable

**Dependencies**: context, index, rc, config/*

---

## internal/index

**Purpose**: Markdown index tables for DECISIONS.md and LEARNINGS.md.

**Exported API**: `Update(content)` regenerates index table between
INDEX:START/END markers. `ParseEntryBlocks(content)` splits file
into timestamped entry blocks.

**Edge cases**:
- Pipe characters in entry titles are escaped in index tables
- Superseded entries detected and marked
- Missing INDEX markers = no-op (not an error)

**Dependencies**: config/regex, config/marker

---

## internal/task

**Purpose**: Task checkbox parsing from TASKS.md.

**Exported API**: `Completed(line)`, `Pending(line)`,
`SubTask(line, indent)`. Match index constants for capture groups.

**Edge cases**: SubTask detection requires indent >= 2 spaces.

**Dependencies**: config/regex

---

## internal/tidy

**Purpose**: Shared helpers for context file maintenance.

**Key types**: `CompactResult` (tasks moved, skipped, file updates,
archivable blocks, sections cleaned)

**Exported API**: `parseBlockAt(lines, startIdx)` parses task
block with nested content. `indentLevel()`, `parseDoneTimestamp()`.

**Data flow**: Called by compact command. Parses TASKS.md into
blocks, identifies archivable blocks (no unchecked children),
moves them to archive section.

**Danger zones**:
1. Block boundary detection uses indentation — tabs vs spaces
   inconsistency causes wrong block boundaries.
2. #done: timestamp parsing is strict — missing or malformed
   timestamps prevent archival eligibility.

**Dependencies**: config/regex, config/time, entity, task

---

## internal/trace

**Purpose**: Link git commits back to decisions, tasks, learnings,
and sessions. Provides Spec: trailer for commits.

**Key types**:
- `PendingEntry`: staged context ref with timestamp
- `HistoryEntry`: commit + attached refs + message
- `OverrideEntry`: explicit post-hoc context association
- `ResolvedRef`: resolved reference (raw, type, number, title,
  detail, found)

**Exported API**: `Collect(contextDir)` gathers refs from pending
records, staged diffs, working state. `FormatTrailer(refs)`
formats as git trailer. `Deduplicate(refs)`.

**Data flow**: Pre-commit hook -> Collect() gathers refs ->
FormatTrailer() adds to commit message -> WriteHistory() persists.

**Danger zones**:
1. StagedRefs() reads git diff — if index changes between
   collect and commit, refs may be stale.
2. Working task refs use regex matching against TASKS.md —
   ambiguous task text causes wrong matches.

**Extension points**:
- New ref types (beyond decision/task/learning) need parser
  additions in resolve functions

**Dependencies**: config/dir, config/trace, exec/git, io

---

## internal/journal/parser

**Purpose**: Session transcript parsing with 4 registered formats.

**Parsers**: Claude Code JSONL, Copilot, Copilot CLI, Markdown.
Auto-detects format via `Matches(path)`.

**Exported API**: `ParseFile(path)` auto-detects and parses.
`ScanDirectory(dir)` recursively scans and aggregates.
`FindSessionsForCWD()` filters by git remote/CWD matching.

**Edge cases**:
- 1MB max buffer for large JSONL lines
- Subagent directories skipped
- Multipart continuations excluded from nav
- Session matching uses git remote URLs and relative paths

**Danger zones**:
1. Claude Code JSONL format is not documented — parser reverse-
   engineers the format. Schema changes break silently.
2. 1MB buffer limit — sessions with very large tool results
   are truncated without warning.
3. Session deduplication by path — same project in different
   locations creates duplicate imports.

**Extension points**:
- Register new parser by implementing SessionParser interface
- Session prefixes configurable via .ctxrc (Decision 2026-03-14)

**Dependencies**: config/parser, config/session, entity, exec/git

---

## internal/journal/state

**Purpose**: Journal processing pipeline state as external JSON.

**Exported API**: `Load(journalDir)` returns empty state if file
missing. `Save()` writes. `MarkImported()`, `MarkEnriched()`,
`MarkNormalized()`, `MarkLocked()`. `CountUnenriched()` for
nudge checks.

**5-stage pipeline**: exported -> enriched -> normalized ->
fences_verified -> locked. Tracked in `.context/journal/.state.json`.

**Edge cases**:
- Missing state file is not an error — returns empty state
- Atomic writes via temp file + rename
- Dates stored as YYYY-MM-DD strings (not timestamps)

**Dependencies**: config/journal, io

---

## internal/memory

**Purpose**: Bridges Claude Code auto memory (MEMORY.md) into
.context/ with discovery, mirroring, drift detection, and
bidirectional sync.

**Key types**:
- `Entry`: discrete block from MEMORY.md (text, startLine, kind)
- `SyncResult`: outcome (sourcePath, mirrorPath, archivedTo)
- `PublishResult`: selected content for publishing

**Exported API**:
- `DiscoverPath(projectRoot)`: locate Claude Code MEMORY.md
  using project slug derivation
- `ProjectSlug(absPath)`: encode path to slug format
- `SelectContent(contextDir, budget)`: read .context/ files,
  select within line budget (priority: tasks > decisions >
  conventions > learnings)
- `MergePublished()`: insert/replace marker block in MEMORY.md
- `Sync()`, `Mirror()`, `Discover()`: full sync operations

**Danger zones**:
1. DiscoverPath depends on Claude Code's project directory
   naming convention — changes to slug format break discovery.
2. MergePublished writes to files outside .context/ — this is
   the only package that modifies external state.
3. SelectContent budget is in lines, not tokens — inconsistent
   with token budgeting elsewhere.

**Extension points**:
- New memory sources beyond Claude Code (Cursor, Copilot) need
  DiscoverPath variants with different slug conventions

**Dependencies**: config/dir, config/memory, entity, err/memory, io

---

## internal/notify

**Purpose**: Fire-and-forget webhook notifications with encrypted
URL storage.

**Exported API**: `Send(event, message)`, `LoadWebhook()`,
`SaveWebhook(url)`.

**Edge cases**: 5s timeout, silent on all errors, event filtering
(opt-in only). Payload: Event, Message, SessionID, Timestamp,
Project.

**Dependencies**: crypto, config/event, io, rc

---

## internal/claude

**Purpose**: Thin wrapper for Claude Code integration. Skill
listing/content retrieval + settings types.

**Key types**: `HookConfig`, `HookMatcher`, `Hook`,
`PermissionsConfig`, `Settings`.

**Exported API**: `Skills()`, `SkillContent(name)`.

**Dependencies**: assets/read/skill, config/claude
