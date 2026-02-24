# Detailed Design

Deep per-module architecture reference. NOT loaded at session start.
Consult specific sections when working on a module.

## internal/config

**Purpose**: Centralized constants, regex patterns, file names, read order, and permissions used across the codebase.

**Key types**: `Pattern` (glob-to-topic mapping)

**Exported API**:
- Constants: file permissions (`PermFile`, `PermExec`, `PermSecret`), file extensions, context file names (`FileConstitution`, `FileTask`, etc.), Claude API block types and field keys, directory names, heading/label/marker constants, limits/thresholds
- `FileType` map — maps entry type strings to filenames
- `FileReadOrder` slice — priority-ordered file loading sequence
- `FilesRequired` slice — essential files for drift detection
- `DefaultClaudePermissions` / `DefaultClaudeDenyPermissions` — permission lists
- `Packages` map — dependency manifest files to descriptions
- `UserInputToEntry(s string) string` — normalizes user input to canonical entry types
- `RegExFromAttrName(name string) *regexp.Regexp` — creates XML attribute extraction regex
- Pre-compiled regex patterns: `RegExEntryHeader`, `RegExTask`, `RegExDecision`, `RegExLearning`, `RegExPath`, `RegExCodeFenceInline`, etc.

**Data flow**: Pure constants package. Consumers import to access patterns, file names, and configuration values. Regex patterns compiled at init time.

**Edge cases**:
- Custom priority orders via `.ctxrc` override `FileReadOrder` defaults
- Obsidian vault output paths coexist with JSON site output
- Migration support for legacy `.scratchpad.key` → `.context.key`

**Dependencies**: None — foundation package with zero internal dependencies

---

## internal/assets

**Purpose**: Embedded templates, skills, tools, and configuration via Go's `//go:embed` directive.

**Key types**: `embed.FS` (embedded filesystem)

**Exported API**:
- `Template(name string) ([]byte, error)` — reads root template by name
- `List() ([]string, error)` — lists root template filenames
- `ListEntry() ([]string, error)` — lists entry template filenames
- `Entry(name string) ([]byte, error)` — reads entry template
- `ListSkills() ([]string, error)` — lists skill directory names
- `SkillContent(name string) ([]byte, error)` — reads SKILL.md for a skill
- `MakefileCtx() ([]byte, error)` — reads Makefile.ctx
- `RalphTemplate(name string) ([]byte, error)` — reads Ralph-mode template
- `ListTools() ([]string, error)` — lists tool script filenames
- `Tool(name string) ([]byte, error)` — reads tool script
- `PluginVersion() (string, error)` — extracts version from embedded plugin.json

**Data flow**: Assets embedded at build time → callers request by name → raw bytes returned or error if not found

**Edge cases**:
- Directory read failures return nil slice with error
- Plugin version requires valid JSON structure

**Dependencies**: `encoding/json` (for plugin.json parsing)

---

## internal/rc

**Purpose**: Runtime configuration loading from `.ctxrc` (YAML) with environment variable overrides and CLI flag precedence.

**Key types**: `CtxRC` (configuration container with ContextDir, TokenBudget, PriorityOrder, AutoArchive, etc.), `NotifyConfig` (webhook settings)

**Exported API**:
- `RC() *CtxRC` — returns cached configuration (lazy-loaded singleton via sync.Once)
- `ContextDir() string` — resolution: CLI override > env > .ctxrc > default
- `TokenBudget() int` — env > .ctxrc > 8000
- `PriorityOrder() []string` — custom file priority or nil
- `AutoArchive() bool`, `ArchiveAfterDays() int` — archive settings
- `ScratchpadEncrypt() bool` — encryption flag (default true)
- `EntryCountLearnings() int`, `EntryCountDecisions() int` — drift thresholds
- `ConventionLineCount() int` — convention line threshold
- `NotifyEvents() []string`, `KeyRotationDays() int` — notification settings
- `AllowOutsideCwd() bool` — boundary check flag
- `FilePriority(name string) int` — priority (1-9) or 100 for unknown
- `OverrideContextDir(dir string)` — sets CLI override
- `Reset()` — clears cache (testing only)

**Data flow**: First call triggers `loadRC()` via sync.Once → reads `.ctxrc` YAML → environment variables override → result cached → CLI overrides stored separately with RWMutex

**Edge cases**:
- Missing `.ctxrc` → uses defaults (not an error)
- Invalid YAML → warning to stderr, defaults used
- `ScratchpadEncrypt` uses nil-pointer triple-state (unset/true/false)

**Dependencies**: `internal/config`, `gopkg.in/yaml.v3`, `sync`

---

## internal/context

**Purpose**: Loads `.context/` directory contents with file metadata, token estimation, and content summarization.

**Key types**: `FileInfo` (Name, Path, Size, ModTime, Content, IsEmpty, Tokens, Summary), `Context` (Dir, Files, TotalTokens, TotalSize), `NotFoundError`

**Exported API**:
- `Load(dir string) (*Context, error)` — loads all .md files from directory
- `Exists(dir string) bool` — checks if directory exists
- `EstimateTokens(content []byte) int` — estimates tokens (1 per 4 chars)
- `EstimateTokensString(s string) int` — convenience wrapper
- `(*Context).File(name string) *FileInfo` — retrieves file by name

**Data flow**: `Load()` → validate directory (exists, no symlinks) → read all .md files → for each: estimate tokens, generate summary, check emptiness → aggregate totals → return `Context`

**Edge cases**:
- Empty directory → Context with empty Files slice
- `.md` files only (other extensions skipped)
- Read errors on individual files → file skipped, processing continues
- "Effectively empty" detected via heuristic (headers, comments, short dashes)
- Symlinks rejected for security (M-2 defense)

**Dependencies**: `internal/config`, `internal/rc`, `internal/validation`

---

## internal/crypto

**Purpose**: AES-256-GCM encryption for scratchpad files with key management.

**Key types**: None (functions only). Constants: `KeySize` = 32, `NonceSize` = 12

**Exported API**:
- `GenerateKey() ([]byte, error)` — generates 32 random bytes
- `LoadKey(path string) ([]byte, error)` — reads and validates key file (must be 32 bytes)
- `SaveKey(path string, key []byte) error` — writes key file with mode 0600
- `Encrypt(key, plaintext []byte) ([]byte, error)` — AES-256-GCM, returns [nonce][ciphertext+tag]
- `Decrypt(key, ciphertext []byte) ([]byte, error)` — extracts nonce, decrypts, authenticates

**Data flow**: `GenerateKey()` → crypto/rand → `SaveKey()` → disk (0600). `Encrypt()`: random nonce → GCM seal → [12-byte nonce + ciphertext + 16-byte tag]. `Decrypt()`: extract nonce → GCM open → plaintext.

**Edge cases**:
- Key size validation before any operation
- Ciphertext too short error (< 12 bytes)
- GCM tag automatically authenticated during decryption
- Random source failure propagated

**Dependencies**: `crypto/aes`, `crypto/cipher`, `crypto/rand` (standard library only)

---

## internal/sysinfo

**Purpose**: OS resource metrics (memory, swap, disk, load) with threshold-based alerting. Platform-specific via build tags.

**Key types**: `Severity` (OK/Warning/Danger), `MemInfo`, `DiskInfo`, `LoadInfo`, `Snapshot`, `ResourceAlert`

**Exported API**:
- `Collect(path string) Snapshot` — gathers metrics (path selects filesystem for disk)
- `Evaluate(snap Snapshot) []ResourceAlert` — checks thresholds (mem ≥80%/90%, swap ≥50%/75%, disk ≥85%/95%, load ≥0.8x/1.5x CPUs)
- `MaxSeverity(alerts []ResourceAlert) Severity` — highest severity in list
- `FormatGiB(bytes uint64) string` — formats bytes as GiB

**Data flow**: `Collect()` → platform-specific collectors (Linux: /proc/meminfo, /proc/loadavg, statfs; macOS: sysctl, vm_stat, statfs; Windows: syscall) → `Evaluate()` → alerts

**Edge cases**:
- Unsupported platform → `Supported=false` (graceful degradation)
- Zero total resources → skipped in Evaluate (prevents divide by zero)
- macOS uses command parsing (shell output errors → Supported=false)

**Dependencies**: Standard library only (platform-specific: `os`, `syscall`, `runtime`, `bufio`)

---

## internal/drift

**Purpose**: Context drift detection — identifies stale paths, completed-task buildup, potential secrets, missing required files, file age, and entry count growth.

**Key types**: `IssueType` (dead_path, staleness, potential_secret, missing_file, stale_age, entry_count), `StatusType` (ok, warning, violation), `CheckName`, `Issue`, `Report`

**Exported API**:
- `Detect(ctx *context.Context) *Report` — runs all six checks
- `(*Report).Status() StatusType` — computes overall status from violations/warnings

**Data flow**: Context files loaded → six sequential checks (path refs, staleness, constitution, required files, age, entry counts) → issues collected → Report returned

**Edge cases**:
- Path checks skip URLs, glob patterns, templates
- Secret detection verifies non-template content
- File age check excludes CONSTITUTION.md (expected to be static)
- Entry count thresholds configurable via rc (0 disables)

**Dependencies**: `internal/config`, `internal/context`, `internal/index`, `internal/rc`

---

## internal/index

**Purpose**: Parse entry headers and manage index tables in DECISIONS.md and LEARNINGS.md.

**Key types**: `Entry` (timestamp, date, title), `EntryBlock` (lines, start/end indices, superseded status)

**Exported API**:
- `ParseHeaders(content string) []Entry` — extracts `## [YYYY-MM-DD-HHMMSS] Title` headers
- `GenerateTable(entries []Entry, columnHeader string) string` — creates markdown index table
- `Update(content, fileHeader, columnHeader string) string` — regenerates index between markers
- `UpdateDecisions(content string) string` / `UpdateLearnings(content string) string` — file-specific wrappers
- `ReindexFile(w io.Writer, filePath, fileName string, updateFunc, entryType string) error` — full reindex workflow
- `ParseEntryBlocks(content string) []EntryBlock` — splits into self-contained entry blocks
- `(*EntryBlock).IsSuperseded() bool` — checks for superseded marker

**Data flow**: Content → regex parse headers → generate table between INDEX:START/END markers → preserve non-entry content

**Edge cases**:
- Pipe characters in titles escaped in table output
- Empty index removes markers and whitespace
- EntryBlocks trim trailing blank lines automatically

**Dependencies**: `internal/config`, `fatih/color`

---

## internal/task

**Purpose**: Domain logic for parsing task checkboxes independent of markdown representation.

**Key types**: Match index constants (`MatchFull`, `MatchIndent`, `MatchState`, `MatchContent`)

**Exported API**:
- `Completed(match []string) bool` — checks if `[x]`
- `Pending(match []string) bool` — checks if `[ ]` or empty
- `Indent(match []string) string` — extracts leading whitespace
- `Content(match []string) string` — extracts task text
- `SubTask(match []string) bool` — true if indent ≥ 2 spaces

**Data flow**: Uses `config.ItemPattern` regex for matching → capture groups → helper functions extract state/content/indent

**Edge cases**: Handles invalid matches gracefully (boundary checks on slice length)

**Dependencies**: `internal/config`

---

## internal/validation

**Purpose**: Input sanitization and path boundary validation.

**Key types**: None (utility functions only)

**Exported API**:
- `SanitizeFilename(s string) string` — converts topic to safe filename (lowercase, hyphenated, max 50 chars)
- `ValidateBoundary(dir string) error` — ensures resolved path stays within cwd
- `CheckSymlinks(dir string) error` — detects symlinks in directory or immediate children

**Data flow**: Sanitize: regex replace → trim → lowercase → limit length. Boundary: resolve symlinks → prefix check. Symlinks: lstat checks for ModeSymlink.

**Edge cases**:
- Non-existent targets fall back to absolute path for prefix check
- Path with separator appended to avoid false prefix matches
- Non-existent directory in CheckSymlinks returns nil

**Dependencies**: `internal/config`

---

## internal/recall/parser

**Purpose**: Parses AI session transcripts (JSONL, Markdown) into structured Go types. Extensible parser registry.

**Key types**: `SessionParser` (interface: ParseFile, ParseLine, Matches, Tool), `ToolUse`, `ToolResult`, `Message`, `Session` (ID, Slug, Tool, SourceFile, CWD, Project, Messages, TurnCount, TokenStats, etc.)

**Exported API**:
- `ParseFile(path string) ([]*Session, error)` — auto-detects format and parses
- `ScanDirectory(dir string) ([]*Session, error)` — recursively finds sessions, sorted newest first
- `ScanDirectoryWithErrors(dir string) ([]*Session, []error, error)` — returns sessions and parse errors
- `FindSessions(additionalDirs ...string) ([]*Session, error)` — searches default + custom locations
- `FindSessionsForCWD(cwd string, additionalDirs ...string) ([]*Session, error)` — filters by CWD (git remote, home path, exact match)
- `Parser(tool string) SessionParser` — gets parser for tool
- `RegisteredTools() []string` — lists supported tools
- `(*Session).UserMessages()`, `(*Session).AssistantMessages()`, `(*Session).AllToolUses()` — message filters
- `(*Message).Preview(maxLen int) string` — truncated text preview

**Data flow** (Claude Code): JSONL line-by-line → parse JSON → group by sessionId → sort by timestamp → convert to Session. Each message's content parsed as text or array of blocks.

**Data flow** (Markdown): Scan for H1 session header → extract H2 sections → build messages → infer project from path pattern.

**Edge cases**:
- Malformed JSONL lines skipped (doesn't fail entire file)
- Large JSONL lines: buffer expanded to 1MB max
- Subagents directory skipped to avoid duplicates
- Git remote matching preferred over path matching for CWD filtering

**Dependencies**: `internal/config`

---

## internal/claude

**Purpose**: Claude Code integration — permissions, hooks, and embedded skill management.

**Key types**: `HookConfig`, `HookMatcher`, `Hook`, `HookType`, `Matcher`, `PermissionsConfig`, `Settings`

**Exported API**:
- `Skills() ([]string, error)` — lists embedded skill directory names
- `SkillContent(name string) ([]byte, error)` — reads SKILL.md for a skill

**Data flow**: Thin wrapper over `internal/assets` — lists skills, retrieves content, wraps errors.

**Dependencies**: `internal/assets`

---

## internal/notify

**Purpose**: Fire-and-forget webhook notifications with encrypted URL storage.

**Key types**: `Payload` (Event, Message, SessionID, Timestamp, Project)

**Exported API**:
- `LoadWebhook() (string, error)` — reads/decrypts webhook URL from `.context/.notify.enc`
- `SaveWebhook(url string) error` — encrypts/writes webhook URL
- `EventAllowed(event string, allowed []string) bool` — checks event filter
- `Send(event, message, sessionID string) error` — fires webhook (silent noop on failure)

**Data flow**: Load: context dir → key file (migrate if needed) → decrypt `.notify.enc` → return URL. Send: check event filter → load URL → build payload → POST with 5s timeout → silent on error.

**Edge cases**:
- Missing key/encrypted file returns ("", nil) — silent noop
- Fire-and-forget: HTTP errors silently ignored
- Empty event list means no events pass (opt-in only)

**Dependencies**: `internal/config`, `internal/crypto`, `internal/rc`

---

## internal/journal/state

**Purpose**: Journal processing state via external JSON file tracking export/enrichment/normalization pipeline.

**Key types**: `JournalState` (Version, Entries map), `FileState` (Exported, Enriched, Normalized, FencesVerified, Locked as date strings)

**Exported API**:
- `Load(journalDir string) (*JournalState, error)` — reads `.state.json` (returns empty if missing)
- `(*JournalState).Save(journalDir string) error` — atomically writes state file
- `(*JournalState).MarkExported/Enriched/Normalized/FencesVerified(filename string)` — sets stage to today
- `(*JournalState).Mark(filename, stage string) bool` / `Clear(filename, stage string) bool` — generic stage ops
- `(*JournalState).Locked(filename string) bool` — checks lock status
- `(*JournalState).Rename(oldName, newName string)` — moves entry state
- `(*JournalState).IsExported/Enriched/Normalized/FencesVerified(filename string) bool` — stage checkers
- `(*JournalState).CountUnenriched(journalDir string) int` — counts .md files without enriched date

**Data flow**: JSON file read/write via atomic temp+rename → stages track processing pipeline → dates as YYYY-MM-DD strings

**Edge cases**:
- Missing file returns empty state (not error)
- CountUnenriched only counts .md files (skips directories)
- Mark/Clear return false for unrecognized stages

**Dependencies**: `internal/config`

---

## internal/bootstrap

**Purpose**: Create root Cobra command, register global flags, attach all subcommands.

**Key types**: None

**Exported API**:
- `RootCmd() *cobra.Command` — creates root command with global flags (--context-dir, --no-color, --allow-outside-cwd) and version
- `Initialize(cmd *cobra.Command) *cobra.Command` — registers all subcommands

**Data flow**: `RootCmd()` creates root → `Initialize()` attaches all CLI packages → `PersistentPreRun` applies global flags and validates context directory boundary

**Edge cases**:
- Context directory boundary validation can be overridden with `--allow-outside-cwd`
- Version injected at build time via ldflags

**Dependencies**: All `internal/cli/*` packages, `internal/rc`

---

## internal/cli/add

**Purpose**: Append entries (decisions, tasks, learnings, conventions) to context files.

**Key types**: `EntryParams` (type, content, Context, Rationale, Consequences, Lesson, Application)

**Exported API**:
- `Cmd() *cobra.Command` — returns "ctx add" command
- `ValidateEntry(params EntryParams) error` — validates required fields
- `WriteEntry(params EntryParams) error` — formats and writes entry

**Data flow**: Parse args → extract content from arg/--file/stdin → validate required fields → format entry → insert at correct location → update index for decisions/learnings

**Edge cases**:
- Tasks insert before first unchecked item or under --section
- Decisions require context+rationale+consequences; learnings require context+lesson+application

**Dependencies**: `internal/config`, `internal/index`, `internal/rc`

---

## internal/cli/agent

**Purpose**: Generate AI-ready context packets with token budgeting.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --budget, --format (md/json), --cooldown, --session

**Data flow**: Read context files → prioritize by recency/relevance → budget-cap → entries that don't fit get title-only summaries in "Also Noted" section → output markdown or JSON

**Edge cases**:
- Cooldown mechanism suppresses repeated output within specified duration per session
- Budget cap is approximate (token estimation)

**Dependencies**: `internal/config`, `internal/rc`

---

## internal/cli/compact

**Purpose**: Archive completed tasks, clean up context files.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --archive

**Data flow**: Read TASKS.md → move completed [x] tasks to "Completed (Recent)" section → if --archive: move to .context/archive/ → remove empty sections

**Dependencies**: `internal/config`, `internal/rc`, `internal/context`, `internal/task`

---

## internal/cli/complete

**Purpose**: Mark tasks as completed in TASKS.md.

**Exported API**:
- `Cmd() *cobra.Command` — args: task-id-or-text (by number, partial text, or full text)

**Data flow**: Accept identifier → read TASKS.md → find matching task → change `- [ ]` to `- [x]` → write back

**Edge cases**: Ambiguous partial matches require clarification

**Dependencies**: `internal/config`, `internal/rc`, `internal/task`

---

## internal/cli/decision

**Purpose**: Manage DECISIONS.md — reindex command.

**Exported API**:
- `Cmd() *cobra.Command` — subcommand: reindex

**Data flow**: Read DECISIONS.md → parse entries → generate compact index table → write back

**Dependencies**: `internal/config`, `internal/rc`, `internal/index`

---

## internal/cli/drift

**Purpose**: Detect stale, invalid, or broken context via CLI.

**Key types**: `JsonOutput` (Timestamp, Status, Warnings, Violations, Passed)

**Exported API**:
- `Cmd() *cobra.Command` — flags: --json, --fix

**Data flow**: Load context → run `drift.Detect()` → output report (human-readable or JSON) → if --fix: auto-fix supported issues

**Edge cases**: Auto-fix supports staleness and missing_file issues

**Dependencies**: `internal/config`, `internal/rc`, `internal/context`, `internal/drift`, `internal/task`

---

## internal/cli/hook

**Purpose**: Generate AI tool integration configurations (Claude Code, Cursor, Aider, Copilot, Windsurf).

**Exported API**:
- `Cmd() *cobra.Command` — flags: --write; args: tool name

**Data flow**: Accept tool name → generate tool-specific config snippet → if --write: write to config file, else print to stdout

**Dependencies**: Cobra only

---

## internal/cli/initialize

**Purpose**: Initialize `.context/` directory with templates, hooks, skills, and project configuration.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --force, --minimal, --merge, --ralph

**Data flow**: Check PATH → create .context/ → prompt if exists → load templates → write files → create entry templates + tools + sessions dir → init scratchpad → create/merge PROMPT.md + IMPLEMENTATION_PLAN.md → merge settings.local.json → handle CLAUDE.md → deploy Makefile.ctx → update .gitignore

**Edge cases**:
- Idempotent: existing files skipped unless --force
- --ralph uses different templates (one-task-per-iteration)
- --merge auto-merges ctx content into existing CLAUDE.md and PROMPT.md
- --minimal only creates essential files

**Dependencies**: `internal/assets`, `internal/config`, `internal/crypto`, `internal/rc`

---

## internal/cli/journal

**Purpose**: Analyze and publish exported AI session files to static sites or Obsidian vaults.

**Exported API**:
- `Cmd() *cobra.Command` — subcommands: site, obsidian

**Data flow**: Scan .context/journal/ → parse YAML frontmatter → generate static site (zensical) or Obsidian vault (wikilinks, MOC)

**Dependencies**: `internal/assets`, `internal/rc`

---

## internal/cli/learnings

**Purpose**: Manage LEARNINGS.md — reindex command.

**Exported API**:
- `Cmd() *cobra.Command` — subcommand: reindex

**Dependencies**: `internal/config`, `internal/rc`, `internal/index`

---

## internal/cli/load

**Purpose**: Output assembled context in priority order with token budgeting.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --budget, --raw

**Data flow**: Load context files → sort by FileReadOrder → truncate to budget → output markdown with assembly headers (or raw if --raw)

**Dependencies**: `internal/config`, `internal/rc`, `internal/context`

---

## internal/cli/loop

**Purpose**: Generate Ralph loop scripts for iterative autonomous development.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --prompt, --tool (claude/aider/generic), --max-iterations, --completion, --output

**Data flow**: Read prompt file → generate shell script with tool-specific invocation + completion signal check → write to output file

**Dependencies**: `internal/config`

---

## internal/cli/notify

**Purpose**: Send fire-and-forget webhook notifications via CLI.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --event, --session-id; subcommands: setup, test

**Data flow**: Accept event + message → call notify.Send() → silent noop if unconfigured or filtered

**Dependencies**: `internal/notify`

---

## internal/cli/pad

**Purpose**: Manage encrypted scratchpad for sensitive one-liners.

**Exported API**:
- `Cmd() *cobra.Command` — subcommands: show, add, rm, edit, mv, resolve, import, export, merge

**Data flow**: Entries encrypted with AES-256-GCM via .context/.context.key. File blobs stored as "label:::base64data". Subcommands: CRUD operations, merge with dedup, import/export for file blobs.

**Edge cases**:
- Blobs limited to 64KB pre-encoding
- Auto-detects encrypted/plaintext in merge
- Merge uses content-based deduplication

**Dependencies**: `internal/crypto`, `internal/rc`

---

## internal/cli/permissions

**Purpose**: Manage Claude Code permission snapshots (golden images).

**Exported API**:
- `Cmd() *cobra.Command` — subcommands: snapshot, restore

**Data flow**: Snapshot: copy settings.local.json → settings.golden.json. Restore: restore from golden, print diff of dropped permissions.

**Dependencies**: `internal/config`

---

## internal/cli/recall

**Purpose**: Browse, search, export, and manage AI session history.

**Exported API**:
- `Cmd() *cobra.Command` — subcommands: list, show, export, lock, unlock, sync; flags: --limit, --project, --tool, --all-projects, --latest, --full

**Data flow**: Parse JSONL session files → subcommands: list (sorted by date), show (by ID/slug/--latest), export (to journal with YAML frontmatter), lock/unlock (protect from overwrite), sync (frontmatter-to-state lock reconciliation)

**Dependencies**: `internal/config`, `internal/rc`, `internal/recall/parser`, `internal/journal/state`

---

## internal/cli/serve

**Purpose**: Serve static sites locally via zensical.

**Exported API**:
- `Cmd() *cobra.Command` — args: directory (default .context/journal-site)

**Edge cases**: Requires zensical installed (`pipx install zensical`)

**Dependencies**: `internal/rc` (external: zensical CLI)

---

## internal/cli/status

**Purpose**: Display context health and summary information.

**Key types**: `Output` (JSON structure), `FileStatus`

**Exported API**:
- `Cmd() *cobra.Command` — flags: --json, --verbose

**Data flow**: Scan .context/ → estimate tokens, check emptiness, generate summaries → output human-readable or JSON

**Dependencies**: `internal/config`, `internal/rc`, `internal/context`

---

## internal/cli/sync

**Purpose**: Reconcile context files with codebase changes.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --dry-run

**Data flow**: Scan codebase for undocumented changes (new dirs, manifest changes, config files) → identify stale references → suggest or apply updates

**Dependencies**: `internal/config`, `internal/context`

---

## internal/cli/system

**Purpose**: System diagnostics, resource monitoring, and Claude Code hook commands (plumbing).

**Exported API**:
- `Cmd() *cobra.Command` — visible subcommands: resources, bootstrap; hidden hook subcommands: check-context-size, check-persistence, check-journal, check-ceremonies, check-version, block-non-path-ctx, post-commit, cleanup-tmp, qa-reminder, check-resources, check-knowledge, mark-journal

**Data flow**: resources: display OS metrics with thresholds. bootstrap: print context directory. Hook commands: read JSON from stdin → perform logic → exit 0 (block commands output JSON with "decision" field).

**Edge cases**:
- Hook commands are plumbing (hidden, used by Claude Code plugin)
- Block commands can veto operations via JSON output
- Adaptive prompt counter with checkpoint messages

**Dependencies**: `internal/config`, `internal/rc`, `internal/sysinfo`, `internal/notify`, `internal/journal/state`

---

## internal/cli/task

**Purpose**: Task archival and snapshots.

**Exported API**:
- `Cmd() *cobra.Command` — subcommands: archive, snapshot

**Data flow**: Archive: read TASKS.md → move completed [x] to timestamped archive in .context/archive/ → preserve Phase structure. Snapshot: create point-in-time copy.

**Dependencies**: `internal/config`, `internal/rc`, `internal/task`, `internal/validation`

---

## internal/cli/watch

**Purpose**: Watch for `<context-update>` tags in AI output and apply them.

**Exported API**:
- `Cmd() *cobra.Command` — flags: --log, --dry-run

**Data flow**: Watch stdin/file for `<context-update type="...">` tags → parse attributes → validate required fields → apply updates (add entry, mark complete, etc.)

**Edge cases**:
- Learnings require: context, lesson, application attributes
- Decisions require: context, rationale, consequences attributes
- Simple types (task, convention, complete) need no attributes

**Dependencies**: `internal/config`, `internal/rc`, `internal/context`, `internal/task`
