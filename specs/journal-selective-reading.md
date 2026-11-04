# Journal Selective Reading

## Problem

In `arXiv:2611.05306` (*You don't need long context*, Section 3.1 and
Appendix E), the authors identify a fundamental tension in compacting
harnesses: compaction is lossy, and a long-horizon task may occasionally
need to inspect exact values or tool outputs from an earlier stretch (the
"coin-toss history" problem). However, loading the entire past conversation
into the active context recreates the long-context penalty. The paper
resolves this by treating past history and long documents as external
filesystem resources queried through **selective reading**: searching with
`grep` to find relevant line/turn numbers, then reading a bounded window
with a continuation pointer (`next: ...`).

Today, `ctx journal source` (`internal/cli/journal/cmd/source/cmd.go` and
`internal/cli/journal/core/source/show.go`) provides only two extremes when
inspecting a session:

1. **Preview mode (default):** Prints session metadata, tool usage counts,
   and only the first 5 user turns truncated to 100 characters.
2. **Full mode (`--full`):** Dumps every turn, message, and tool result of
   the entire session to stdout with no bound.

An agent cannot search inside a session transcript or read a specific slice
of turns (for example, turns 15 to 25) without dumping the entire session
into its context window.

## Approach

Extend `ctx journal source` show mode with three selective-reading flags
modeled directly on Appendix E.1 of `arXiv:2611.05306`:

1. **`--turns <start:end>`:** 1-based inclusive turn range over the
   session's messages (for example, `--turns 12:18` or `--turns 12:12`,
   matching the 1-based `### Turn N` indices rendered by `--full`). Renders
   full message bodies, tool invocations, and tool outputs only for turns
   within `[start, end]`.
2. **`--grep <pattern>`:** Case-insensitive regular expression search
   across all turns in the target session (matching user text, assistant
   text, tool names/inputs, and tool output content). Emits a compact list
   of matching turn numbers, roles, timestamps, and a 200-character
   snippet around each match, instructing the caller to inspect surrounding
   context with `--turns <start:end>`.
3. **`--max-chars <n>`:** Character budget cap for bounded output
   (defaults to `6000` when `--turns` or `--grep` is active; `0` =
   unlimited when bare `--full` is explicitly requested). When output is
   truncated by `--max-chars` or when `--turns` ends before the final turn
   of the session, `ctx journal source` prints an explicit header/footer
   containing the exact command to read the next contiguous window:

```text
[session 96765858 turns 10:15 of 42 | 4120 chars] next: ctx journal source --show 96765858 --turns 16:21
```

## Behavior

### Happy Path

1. The agent needs to recall an exact error message or command output from
   a prior session (`96765858`) that was compacted away.
2. The agent runs:
   ```bash
   ctx journal source --show 96765858 --grep "cipher: message authentication"
   ```
3. `ctx journal source` prints session metadata and the matching turn hits:
   ```text
   Matches in session 96765858 (2 of 42 turns):
   - Turn 18 (Tool Output, 14:22:10): ...FAIL internal/cli/pad: cipher: message authentication failed...
   - Turn 19 (Assistant, 14:22:18): ...The pad test overwrote ~/.ctx/.ctx.key...
   Hint: read surrounding turns with: ctx journal source --show 96765858 --turns 17:20
   ```
4. The agent runs:
   ```bash
   ctx journal source --show 96765858 --turns 17:20
   ```
5. `ctx journal source` renders turns 17 through 20 in full, capped at
   `6000` characters, and prints the `next:` continuation line for turns
   `21:24`.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `--turns` has invalid syntax (e.g., `abc`, `10:5`, `0:5`, `-1:3`) | Reject with a typed validation error explaining the `START:END` 1-based inclusive format (`1 <= START <= END`). |
| `--turns 50:60` when the session has only 42 turns | Clamp `END` to `42`; if `START > 42`, print an informative empty-range notice (`session has 42 turns; requested start 50 is past the end`) and exit 0. |
| `--grep` pattern is an invalid regular expression | Return a typed regex compilation error naming the invalid pattern. |
| `--grep` finds zero matches in the session | Print `No turns matched pattern "<pattern>" in session <id> (<N> turns searched).` and exit 0. |
| Single turn's tool output exceeds `--max-chars 6000` | Truncate the turn's rendered content at the character cap on a rune boundary, append `\n... [truncated at 6000 chars]`, and emit the `next:` command pointing at the next turn. |
| Both `--grep` and `--turns` are supplied | Filter to the turn range `[start, end]` first, then apply `--grep` within that slice. |

### Validation Rules

- `--turns` format is `<start>:<end>` or `<turn>` (shorthand for
  `<turn>:<turn>`), where `start` and `end` are positive integers with
  `1 <= start <= end`.
- `--max-chars` must be `>= 0` (`0` disables character truncation).
- `--grep` must compile via `regex.CompileDynamic(regex.CaseInsensitivePrefix + pattern)` defined in `internal/config/regex/`.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Malformed `--turns` range | `Error: invalid --turns range "10:5": expected START:END with 1 <= START <= END` | Pass a valid 1-based range such as `--turns 5:10`. |
| Invalid `--grep` regex | `Error: invalid --grep regular expression "[unclosed": ...` | Escape regex metacharacters or fix the pattern. |
| Negative `--max-chars` | `Error: --max-chars must be non-negative` | Pass a positive integer (e.g., `6000`) or `0` for unlimited. |

## Interface

### CLI

```bash
ctx journal source [--show <id> | --latest] [--turns <start:end>] [--grep <regex>] [--max-chars <n>] [--full]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--turns` | `string` | `""` | 1-based inclusive turn range `START:END` (or single `TURN`) to display. |
| `--grep` | `string` | `""` | Case-insensitive regex pattern to search across session messages and tool outputs. |
| `--max-chars` | `int` | `6000` (bounded) / `0` (`--full`) | Maximum characters to emit in bounded show mode (`0` for unlimited). |

### Skill

Update `/ctx-history` and `/ctx-remember` skill instructions to prefer
`ctx journal source --show <id> --grep "<pattern>"` and
`ctx journal source --show <id> --turns <start:end>` over `--full` so
agents never dump an entire raw session transcript into a short context.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/config/flag/flag.go` | Add flag constants `Turns = "turns"`, `Grep = "grep"`, and `MaxChars = "max-chars"`. |
| `internal/config/embed/flag/journal.go` | Add description keys `DescKeyJournalSourceTurns`, `DescKeyJournalSourceGrep`, and `DescKeyJournalSourceMaxChars`. |
| `internal/assets/commands/flags.yaml` | Add flag descriptions for `journal.source.turns`, `journal.source.grep`, and `journal.source.max-chars`. |
| `internal/config/journal/len.go` & `limit.go` | Add `DefaultSelectiveMaxChars = 6000` and `GrepSnippetMaxChars = 200`. |
| `internal/config/regex/journal.go` | Add `CaseInsensitivePrefix = "(?i)"` and `CompileDynamic` for compiling dynamic CLI `--grep` patterns within `internal/config/regex/`. |
| `internal/config/embed/text/err_session.go` & `journal.go` | Add `DescKeyErrSession*` and `DescKeyJournalSource*` text keys for selective reading errors and output lines. |
| `internal/assets/commands/text/errors.yaml` & `write.yaml` | Add error templates and output format strings for `--turns`, `--grep`, and `--max-chars`. |
| `internal/err/session/session.go` | Add error constructors `InvalidTurnRange`, `InvalidGrepPattern`, and `NegativeMaxChars`. |
| `internal/write/journal/source.go` | Add output helpers for selective turn range banners, truncation footers, and grep match listings. |
| `internal/cli/journal/core/source/types.go` | Add `Turns string`, `Grep string`, and `MaxChars int` to `source.Opts`, plus `TurnMatch` struct. |
| `internal/cli/journal/cmd/source/cmd.go` | Bind `--turns`, `--grep`, and `--max-chars` flags and forward to `Run`. |
| `internal/cli/journal/cmd/source/run.go` | Trigger `RunShow` when `--turns` or `--grep` is non-empty (defaulting to `--latest` if neither `--show` nor positional ID is given). |
| `internal/cli/journal/core/source/selective.go` | Export `ParseTurnRange` and `MatchTurns` (keeping pure exported functions in `selective.go` and unexported rendering helpers in `selective_render.go`). |
| `internal/cli/journal/core/source/show.go` | Dispatch to selective reading when `opts.Turns != ""` or `opts.Grep != ""`. |
| `internal/cli/journal/core/source/slice_test.go` | Unit tests for `ParseTurnRange`, regex matching, character truncation, and `next:` continuation command generation. |
| `internal/assets/claude/skills/ctx-history/SKILL.md` & `ctx-remember/SKILL.md` | Document selective reading with `--grep` and `--turns` (and sync mirrored skill trees). |
| `docs/cli/journal.md` | Document `--turns`, `--grep`, and `--max-chars` with selective-reading examples. |

### Key Functions

```go
// ParseTurnRange parses a "START:END" or "N" turn selector into 1-based
// inclusive bounds clamped against totalTurns.
func ParseTurnRange(spec string, totalTurns int) (start, end int, err error)

// MatchTurns filters session messages within [start, end] against a
// case-insensitive regular expression.
func MatchTurns(
	msgs []entity.Message, start, end int, re *regexp.Regexp,
) []TurnMatch
```

### Helpers to Reuse

- `parse.StripLineNumbers` (`internal/parse`) when rendering or grepping
  tool result bodies.
- `writeRecall.ConversationTurn` and `writeRecall.CodeBlock`
  (`internal/write/journal`) for turn rendering.

## Configuration

No new `.ctxrc` keys required. `DefaultSelectiveMaxChars` (`6000`) matches
the 6,000-character window validated in `arXiv:2611.05306` Appendix E.1.

## Testing

- **Unit (`internal/cli/journal/core/source/slice_test.go`):**
  - `ParseTurnRange`: valid ranges (`"1:5"`, `"7"`, `"10:100"` clamped to
    total), invalid syntax (`"0:5"`, `"5:2"`, `"foo"`).
  - `MatchTurns`: matches user text, assistant text, tool input arguments,
    and tool result stdout/stderr; returns empty slice on non-match.
  - Bounded rendering: truncates at `--max-chars` on UTF-8 rune boundaries
    and emits accurate `next:` continuation commands.

## Non-Goals

- Full-text indexing or external vector/SQLite search over `.context/journal/`
  (standard regex scanning over a single session's parsed messages is
  sub-millisecond and dependency-free).
- Replacing `grep` or file reading tools for normal source code files in the
  repository.
