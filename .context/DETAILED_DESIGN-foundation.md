# Detailed Design: Foundation Layer

Modules: config/*, assets/*, io, format, parse, sanitize, validate,
inspect, flagbind, exec/*, log/*, crypto, sysinfo, rc

## internal/config/*

**Purpose**: 60+ sub-packages of pure constants, compiled regexes,
text keys, and type definitions. Zero internal dependencies.

**Organizational principle**: Domain-driven. Each sub-package groups
related constants by concept (agent, entry, file, mcp, regex, etc.).

**Key sub-package categories**:

- Flat single-file: agent/, cli/, dir/, entry/, env/, flag/, fmt/,
  marker/, version/, runtime/, wrap/ — pure `const` blocks
- Multi-file thematic: regex/ (14 files of compiled regexps),
  file/ (extensions, ignore, names, limits)
- Hierarchical: embed/ (cmd/, flag/, text/ — 130+ files of
  user-facing strings), mcp/ (12 sub-packages for protocol constants)

**Key types**: `entry.FromUserInput()` is the only function; all
else is const/var declarations.

**Data flow**: Imported granularly by all upper layers via
`config/agent`, `config/mcp/tool`, `config/embed/text/`, etc.

**Edge cases**:
- embed/text/ contains 879+ description keys; exhaustive test
  verifies all resolve to non-empty YAML values
- regex/ patterns are compiled at package init time via MustCompile

**Performance considerations**:
- All constants resolved at compile/init time — zero runtime cost
- Regex compilation happens once at package load

**Danger zones**:
1. `config/embed/text/` — adding a DescKey without a YAML entry
   causes runtime panic (MustCompile-style failure). Caught by
   TestDescKeyYAMLLinkage audit.
2. `config/regex/` — changing a regex pattern affects every consumer
   silently. No type safety on match group indices.
3. `config/file/` — FileReadOrder determines context priority;
   reordering changes what agents see first.

**Extension points**:
- Add new config domain by creating a new sub-package with doc.go
- Add new DescKey constants + YAML entries for new text

**Improvement ideas**:
- Consider code-generating DescKey constants from YAML source
  to eliminate manual sync

**Dependencies**: stdlib only (regexp, time, strings, path)

---

## internal/assets

**Purpose**: Single `go:embed` declaration exposing all embedded
files. 14 typed reader sub-packages under `assets/read/*` provide
domain-specific accessors.

**Key types**: `var FS embed.FS` (the single embed)

**Exported API** (via read/* sub-packages):
- `desc.Text(key)`, `desc.Flag(name)`, `desc.Command(key)`: YAML lookups
- `skill.List()`, `skill.Content(name)`: skill definitions
- `entry.List()`, `entry.ForName(name)`: entry templates
- `claude.Md()`: CLAUDE.md content
- `hook.Message(name)`, `hook.MessageRegistry()`: hook messages
- `project.Readme(dir)`: directory README templates
- `schema.Schema()`: .ctxrc JSON Schema
- `catalog.List()`: template file listing
- `lookup.Init()`: pre-loads all YAML maps at startup

**Data flow**: `assets.FS` -> `read/*` sub-packages -> all consumers.
Dependencies flow one way: read/* imports assets.FS; assets never
imports read/*.

**Edge cases**:
- `lookup.Init()` must be called before any desc.Text() calls
- tpl/ contains Sprintf templates marked for migration to
  text/template (tracked in TASKS.md)
- hooks/messages/ has 38 entries in registry.yaml with template
  variables (%s placeholders)

**Performance considerations**:
- YAML maps loaded once via Init(); subsequent lookups are map reads
- Embedded FS is compiled into binary — no runtime file I/O

**Danger zones**:
1. `lookup.Init()` ordering — calling desc.Text() before Init()
   returns empty strings silently (no panic).
2. `tpl/` format strings — mismatched %s/%d placeholders cause
   runtime panics. No compile-time checking.
3. `hooks/messages/registry.yaml` — adding a hook entry without
   a corresponding .txt file causes silent empty message.

**Extension points**:
- Add new read/* sub-package for new embedded asset type
- Add new hook messages by creating .txt + registry.yaml entry

**Improvement ideas**:
- Migrate tpl/ to text/template for safer templating
- Add compile-time template validation

**Dependencies**: embed (stdlib)

---

## internal/io

**Purpose**: Guarded file I/O wrappers applying path validation,
symlink protection, and deny-list filtering on every operation.

**Key types**: None (pure function package)

**Exported API**:
- `SafeReadFile(path)`: bounded, containment-checked read
- `SafeReadUserFile(path)`: deny-list filtered read
- `SafeWriteFile(path, data, perm)`: write with validation
- `SafeCreateFile(path, perm)`: create with validation
- `SafeAppendFile(path, data, perm)`: append with validation
- `SafePost(url, body)`: HTTP POST with scheme/redirect limits
- `SafeMkdirAll(path, perm)`: directory creation with validation
- `SafeStat(path)`: stat with validation
- `TouchFile(path)`: best-effort marker creation
- `SafeFprintf(w, format, args)`: guarded formatted write

**Data flow**: All file operations in the codebase route through
this package. Direct os.ReadFile/os.WriteFile calls are banned
(enforced by TestNoRawFileIO audit).

**Edge cases**:
- SafePost limits redirects and rejects non-http(s) schemes
- TouchFile is best-effort — errors logged via warn, not returned
- Nil writer in SafeFprintf is a silent no-op

**Performance considerations**:
- Every file op pays validation overhead (path resolve + prefix
  check). Negligible for CLI workloads.
- SafeReadFile has no size limit — large files load fully

**Danger zones**:
1. Path validation relies on resolved prefix matching — if the
   project root changes after validation, the check is stale.
2. SafeReadUserFile deny-list is static — new sensitive patterns
   need code changes.
3. SafePost timeout is not configurable — hardcoded limits.

**Extension points**:
- Add new Safe* variants for new I/O patterns
- Deny-list is extensible via config/fs constants

**Dependencies**: config/fs, config/warn, err/fs, err/http, log/warn

---

## internal/format

**Purpose**: Converts typed values into human-readable display
strings for terminal output.

**Exported API**:
- `TimeAgo(t)`: relative time (e.g., "3 hours ago")
- `Duration(d)`, `DurationAgo(t)`: duration formatting
- `Today()`: YYYY-MM-DD date string
- `TruncateFirstLine(s, n)`, `Truncate(s, n)`: with ellipsis
- `Number(n)`: thousands-separated
- `Bytes(n)`: binary units (B, KB, MB, GB)
- `Tokens(n)`: SI-suffix (K, M)

**Dependencies**: assets/read/desc, config/format, config/time,
config/token

---

## internal/exec/*

**Purpose**: Centralized external command wrappers. All os/exec
calls in the codebase go through these packages (enforced by audit).

| Sub-package | Wraps | Key Functions |
|-------------|-------|---------------|
| `exec/git` | git | `Run()`, `Root()`, `RemoteURL()`, `LogSince()`, `DiffTreeHead()` |
| `exec/dep` | go list | `GoListPackages()` |
| `exec/sysinfo` | sysctl, vm_stat | `Sysctl()`, `VMStat()` (darwin only) |
| `exec/gio` | gio mount | `Mount(url)` |
| `exec/zensical` | zensical | `Run(dir, command)` |

**Danger zones**:
1. `exec/git` — LookPath checked per call; missing git is a
   runtime error, not a startup check.
2. `exec/zensical` — inherits process stdin/stdout/stderr, which
   can interfere with MCP server if called from that context.

**Dependencies**: config/git, config/dep, config/archive,
config/zensical, err/git, err/site

---

## internal/log/*

**Purpose**: Diagnostic logging split into two packages to avoid
import cycles (Decision 2026-03-31).

**log/event**: Timestamped JSONL event logging with size-based
rotation. `Append()` writes, `Query()` filters. Best-effort —
never blocks operations on failure.

**log/warn**: Stderr warning sink. Single function `Warn()` for
errors that should not block execution (file close, state write).

**Dependencies**: entity, config/event, io, rc (event); none (warn)

---

## internal/crypto

**Purpose**: AES-256-GCM encryption for scratchpad and webhook URLs.

**Exported API**: `Encrypt(key, plaintext)`, `Decrypt(key, ciphertext)`,
`GenerateKey()`, `LoadKey(path)`, `SaveKey(path, key)`.

**Edge cases**: Key files must be 0600 permissions. 12-byte random
nonce prepended to ciphertext. Global key at `~/.ctx/.ctx.key`.

**Dependencies**: stdlib only (crypto/aes, crypto/cipher, crypto/rand)

---

## internal/sysinfo

**Purpose**: Platform-specific resource monitoring via build tags.

**Exported API**: `Collect()` gathers metrics, `Evaluate()` applies
thresholds (mem 80/90%, swap 50/75%, disk 85/95%, load 0.8x/1.5x
CPUs). Graceful degradation on unsupported platforms.

**Dependencies**: exec/sysinfo, config/sysinfo

---

## internal/rc

**Purpose**: Runtime configuration resolution with sync.Once
singleton caching.

**Exported API**: `RC()` (cached global), `ContextDir()`,
`TokenBudget()`, `PriorityOrder()`, `AutoArchive()`,
`AllowOutsideCwd()`.

**Resolution order** (highest wins): CLI flags -> env vars
(CTX_DIR, CTX_TOKEN_BUDGET) -> .ctxrc YAML -> defaults from
config/runtime.

**Danger zones**:
1. sync.Once means the first call to RC() locks the config for
   the process lifetime — late flag overrides have no effect.
2. Non-fatal YAML parse errors are printed as warnings and
   silently use defaults.

**Dependencies**: config/runtime, yaml.v3

---

## internal/parse

**Purpose**: Shared text-to-typed-value conversion.

**Exported API**: `Date(s string)` parses YYYY-MM-DD to time.Time.

**Dependencies**: config/time

---

## internal/sanitize

**Purpose**: Mutate untrusted input to conform to constraints.

**Exported API**: `Filename(s)` converts topic string to safe
filename (lowercase, hyphenated, max 50 chars).

**Dependencies**: config/file, config/regex, config/session

---

## internal/validate

**Purpose**: Path validation and symlink detection.

**Dependencies**: config/fs, rc

---

## internal/inspect

**Purpose**: String predicates and position queries for text parsing.

**Exported API**: `SkipNewline()`, `FindNewline()`, `Contains()`,
`StartsWithCtxMarker()`, etc.

**Dependencies**: config/marker, config/token

---

## internal/flagbind

**Purpose**: Cobra flag binding enforcing YAML-backed descriptions.
All flag descriptions go through `assets/read/desc.Flag()` lookup.

**Exported API**: `BoolFlag()`, `StringFlag()`, `IntFlag()`,
`DurationFlag()`, `PersistentBoolFlag()`, `LastJSON()`, and P
(shorthand) variants.

**Danger zones**:
1. Using cobra flag methods directly bypasses the YAML description
   system — enforced by TestNoFlagBindOutsideFlagbind audit.

**Dependencies**: spf13/cobra, assets/read/desc, config/flag
