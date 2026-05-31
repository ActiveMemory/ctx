# tpl-text-template-migration

Covers TASKS.md task 252 ("Migrate Sprintf-based templates
(`tpl_*.go`) to Go `text/template` or embedded template files").

## Problem

The multi-line block templates in `internal/assets/tpl/tpl_*.go` are
stored as `fmt.Sprintf` format-string constants. For documents,
scripts, and config blocks this is the wrong tool:

- **Positional `%s`/`%d` verbs are unreadable and unsafe at scale.**
  `LoopScript` takes six positional args assembled in a precise order
  (`script.go:61`); `Decision` passes `title` twice
  (`fmt.go:100`). A reordered argument is a silent corruption, not a
  compile error.
- **The copy lives inside `.go` source**, so editing a generated
  README or a TOML block means editing Go string literals with
  backtick-escaping gymnastics (`tpl_obsidian.go` interleaves
  `` ` `` + `"..."` + `` ` `` just to embed a fenced block).
- **HTML is assembled by scattered paired-tag writes.** The recall
  formatter (`source/format/format.go`) builds `<details>`/`<table>`
  blocks by emitting an open constant, looping rows, then a close
  constant across ~25 lines. The open/close pair is a structural
  invariant smeared across the call site — its own code smell.
- **`tpl_obsidian.go`'s own docstring already prescribes the fix**:
  "should migrate to a Go text/template or an embedded template file
  when the template rendering pipeline is implemented (see
  TASKS.md)." This spec is that pipeline.

These templates can't move to the YAML `desc.Text` system (which is
for short/long single-string descriptions); they need real template
rendering.

## Settled Decisions

Resolved during spec review (2026-05-30):

1. **Tier-3 stays `fmt.Sprintf`.** Pure positional joins
   (`RecallFencedBlock = "%s\n%s\n%s"`, `Fm*`, `ToolDisplay`) and the
   `RecallListRow` meta-format are not templates; converting them adds
   indirection (and a name surface) for no readability gain.
2. **Tier-2 is refactored, not demoted.** The interleaved paired-tag
   call sites are the smell; the fix is two data-driven block
   templates that own the structure, per the no-broken-windows
   invariant — not leaving them as scattered `Sprintf` because it is
   easier.
3. **No `panic` on init parse.** Parse-at-init + a `TestTemplatesParse`
   CI guard + an error-returning `Render`. No `template.Must` (it
   panics, and has no precedent in this repo).

## Approach

Move multi-line template **text out of `.go` into embedded files**
under `internal/assets/tpl/templates/`, parsed once via Go
`text/template`, following the existing pattern in
`internal/cli/system/core/message/render.go`. The embedded `assets.FS`
(`internal/assets/embed.go`) is the delivery mechanism — it already
embeds an `entry-templates/*.md` set, so a sibling `tpl/templates/`
glob is idiomatic.

**No magic strings (hard constraint).** The exported identifier is
preserved but retyped: `tpl.ObsidianReadme` changes from a
`string` format constant to a parsed `*template.Template` handle.
Call sites reference the **handle**, never a name literal:

```go
// before
[]byte(fmt.Sprintf(tpl.ObsidianReadme, journalDir))
// after
out, err := tpl.Render(tpl.ObsidianReadme, obsidianData{JournalDir: journalDir})
```

The string filename (`"tpl/templates/obsidian-readme.md.tmpl"`) appears
**exactly once**, in the parse table inside the `tpl` package — never
at a call site. This satisfies `audit/magic_strings_test.go`, which
would have rejected an earlier `Render("obsidian-readme", …)` sketch.

### Three tiers (full inventory below)

| Tier | What | Treatment |
|------|------|-----------|
| **1 — Blocks** | Multi-line documents/scripts/config | One embedded file each; `*.tmpl` (interpolated) or static (`Zensical*`) |
| **2 — HTML assembly** | Recall `<details>`/`<table>` blocks built from paired-tag constants | Refactor into two data-driven block templates (`metaTable`, `details`); the paired constants are deleted |
| **3 — Joins** | Single-line format strings + pure positional joins + the meta-format | **Stay `fmt.Sprintf` consts** (not templates) |

### Rendering helper

Generalize `message/render.go` into the `tpl` package:

```go
// Render executes a parsed template handle against data, returning
// the rendered string. A non-nil error means a programmer bug (bad
// field, malformed embedded template) — golden tests gate against it.
func Render(t *template.Template, data any) (string, error)
```

Templates are parsed at package init from the embedded FS into the
exported handles. Parse failures are collected (not panicked) and
asserted empty by `TestTemplatesParse`, so a malformed embedded
template fails CI rather than reaching production.

### Tier-2 refactor detail

Two block templates replace six paired-tag constants
(`MetaDetailsOpen/Close`, `MetaRow`, `RecallDetailsOpen/Close`,
`RecallPlanOpen/Close`):

- **`metaTable`** — input `{Summary string; Rows []struct{Label, Value string}}`.
  Replaces `format.go:255-276` and `280-293`: build the rows slice
  (conditional rows like `GitBranch`/`Model`/`Parts` become
  conditional appends), render once. `MetaRow` becomes a `{{range}}`
  body, not a standalone const.
- **`details`** — input `{Summary, Body string}`. Replaces the three
  open/close pairs (`format.go:357-359`, `396-400`,
  `collapse.go:92-100`): the caller builds the inner body string
  (e.g. `<pre>`-escaped content) and the template wraps it.

## Behavior

### Happy Path

1. At `tpl` init, each `*.tmpl` file is read from `assets.FS` and
   parsed into its exported `*template.Template` handle.
2. A call site builds a typed data struct and calls
   `tpl.Render(tpl.X, data)`.
3. `Render` executes into a `bytes.Buffer` and returns the string —
   **byte-for-byte identical** to today's output, trailing newlines
   included.
4. Static blocks (`ZensicalProject`, `ZensicalTheme`) are exposed as
   `string` values loaded from their embedded files at init; their
   `sb.WriteString(...)` call sites (`generate.go:182,242`) are
   unchanged.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Empty data field (e.g. empty `journalDir`) | Renders the empty string into the placeholder — same as `Sprintf("%s","")`. No special-casing. |
| `LoopScript` with `maxIterations == 0` | `{{if .MaxIter}}…{{end}}` renders nothing — replaces the "inject empty `maxIterCheck`" composition (`script.go:53-59`). Output identical. |
| `LoopScript` tool selection | `aiCommand` is chosen in Go (small `LoopCmd*` consts stay) and passed as `{{.AICommand}}`; the template does not branch on tool. |
| `metaTable` conditional rows | Absent `GitBranch`/`Model`/`Parts` append no row — matches the current `if s.X != ""` guards exactly. |
| **Whitespace fidelity (the chief hazard)** | `MetaDetailsOpen` ends `<table>` with *no* newline; the first `<tr>` follows on the same line. The `{{range}}`/`{{define}}` blocks need explicit `{{-`/`-}}` trimming to reproduce exact bytes. Golden tests assert this. |
| Malformed embedded template ships | `init` records the parse error; `TestTemplatesParse` fails in CI. Cannot reach a release. |
| Exec error (missing/renamed field) | `Render` returns non-nil error; the call site's golden test fails pre-merge. |

### Validation Rules

Template data is passed as typed structs (one per template), so field
presence is compile-checked. No runtime input validation is added —
inputs are already-validated values from existing call sites.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Init parse failure (malformed `.tmpl`) | None in prod (CI-gated); dev sees `TestTemplatesParse` failure naming the file | Fix the template file |
| Exec error (bad field) | Propagated by `Render`; call sites that today return a bare `string` (`SiteReadme`, `script.Generate`) gain an `error` return or a test-guaranteed wrapper | Golden test catches pre-merge |

## Interface

Internal refactor — **no CLI, no skill, no user-visible surface
change**. The "interface" is the `tpl` package API: exported
`*template.Template` handles + static `string`s + `Render`. Output of
every affected command is byte-identical.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/tpl/templates/*.tmpl`, `*.toml` | **New** — extracted Tier-1 bodies + `blocks.tmpl` holding the `metaTable` and `details` `{{define}}`s |
| `internal/assets/tpl/render.go` | **New** — `Render(t, data)`, init parse table (the only place filenames appear), parse-error collection, `TestTemplatesParse` target |
| `internal/assets/embed.go` | Add `//go:embed tpl/templates/*` glob |
| `internal/assets/tpl/tpl_*.go` | Retype migrated consts → handles / FS-loaded strings; delete migrated bodies + the six Tier-2 paired-tag consts; Tier-3 consts stay |
| `internal/cli/journal/core/source/format/format.go` | Tier-2 refactor: build `metaTable` rows + `details` bodies, render via handles (replaces `255-293`, `357-359`, `394-400`) |
| `internal/cli/journal/core/collapse/collapse.go` | Tier-2: `92-100` → `details` render |
| `internal/cli/journal/core/obsidian/vault.go:91` | `Sprintf(tpl.ObsidianReadme,…)` → `Render` |
| `internal/cli/journal/core/generate/generate.go:37` | `SiteReadme` → `Render`; `Zensical*` `WriteString` unchanged (FS-loaded strings) |
| `internal/cli/loop/core/script/script.go:61` | Replace 6-arg `Sprintf` + `maxIterCheck` pre-format with one `Render(tpl.LoopScript, loopData{…})` |
| `internal/cli/trigger/cmd/add/cmd.go:93` | `Sprintf(tpl.TriggerScript,…)` → `Render` |
| `internal/cli/add/core/format/fmt.go:63-101` | `Learning`/`Decision` → `Render` (removes the double-`title` positional surface) |

### Helpers to Reuse

- `internal/cli/system/core/message/render.go` — the parse+execute+buffer
  pattern to generalize (don't reinvent).
- `internal/assets` `embed.FS` — existing embed delivery.
- `internal/io.SafeWriteFile` / `SafeFprintf` — unchanged where Tier-3
  consts remain.

### Full Inventory (every `tpl_*.go` constant)

**Tier 1 — embedded files:** `ObsidianReadme`, `JournalSiteReadme`,
`LoopScript` (absorbs `LoopMaxIter` as `{{if .MaxIter}}` and
`LoopNotify` as a `{{define}}`), `TriggerScript`, `Learning`,
`Decision`; static: `ZensicalProject`, `ZensicalTheme`.

**Tier 2 — absorbed into block templates (consts deleted):**
`MetaDetailsOpen`, `MetaDetailsClose`, `MetaRow` → `metaTable`;
`RecallDetailsOpen`, `RecallDetailsClose`, `RecallPlanOpen`,
`RecallPlanClose` → `details`.

**Tier 3 — stay `fmt.Sprintf`:** single-line format strings
(`LoadBudget`, `LoadSectionHeading`, `RecallTurnHeader`,
`RecallDetailsSummary`, `JournalMonthHeading`, `Task*`, `Convention`,
`HubEntryMarkdown`, `JournalNav*`, stats lines, `LoopCmd*`, …); pure
positional joins (`RecallFencedBlock`, `Fm{Quoted,String,Int}`,
`ToolDisplay`, `RecallFilename`, `RecallPartFilename`); the meta-format
`RecallListRow`.

## Configuration

None. No `.ctxrc` keys, environment variables, or settings.

## Testing

- **Golden equivalence (the core guarantee):** for every migrated
  template, assert `Render(handle, data)` is byte-for-byte equal to the
  legacy `fmt.Sprintf(oldConst, args)` output for representative
  inputs. Capture legacy output as a golden fixture *before* deleting
  the old const.
- **Tier-2 assembly goldens:** full-output tests for the two metadata
  tables (with/without `GitBranch`/`Model`/`Parts`), the plan block,
  the tool-result `<details>` (collapsed and fenced branches), and
  `collapse.go` (wrapped and already-wrapped). These guard the
  whitespace-fidelity hazard.
- **`TestTemplatesParse`:** asserts the init parse-error set is empty.
- **Per-call-site tests:** `loop/core/script` (with/without
  max-iterations, each tool), `trigger/cmd/add`, `add/core/format`,
  `journal/core/generate` (SiteReadme + full `ZensicalToml`).
- **Compliance:** `internal/audit/magic_strings_test.go` and the
  `compliance` suite stay green (no name literals at call sites).

## Non-Goals

- **Not** migrating Tier-3 format strings, pure joins, or the
  `RecallListRow` meta-format — they are not templates.
- **Not** changing any rendered output — behavior-preserving, asserted
  by golden tests.
- **Not** touching the YAML `desc.Text` system or moving anything into
  YAML.
- **Not** adding caching/perf work; init-time parse is sufficient.
- **Not** restructuring the recall formatter beyond the `<details>`/
  `<table>` assembly — only the paired-tag smell is in scope.
