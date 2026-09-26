# Spec: fix CRLF knowledge files and malformed entry headers

Three defects in how ctx reads LEARNINGS.md / DECISIONS.md, all
field-observed on Windows, where a `core.autocrlf=true` checkout gives
every context file CRLF line endings. Two are CRLF sensitivities; the
third is a silent data-integrity hole in the progressive-disclosure
pass that CRLF debugging surfaced.

## Problem

### A. Entry titles keep the carriage return

`heading.ParseEntryBlocks` (`internal/heading/entry.go`) splits content
on `"\n"`, so on a CRLF file every line keeps a trailing `"\r"`.
`regex.EntryHeader` is `## \[(\d{4}-\d{2}-\d{2})-(\d{6})] (.+)`, and
`(.+)` captures that `"\r"` into the title. `heading.ParseHeaders`
(`internal/heading/index.go`) has the same bug: it runs
`FindAllStringSubmatch` over the whole file, and `.` stops at `"\n"`
but not at `"\r"`.

Observed: `ctx disclosure inspect .context/LEARNINGS.md --json` printed
titles such as `"Foo (consolidated)\r"`. A digest plan written from
those titles either carries the stray `"\r"` or, when a human or agent
types the title cleanly, fails `apply` with `ErrEntryNotInStaging`.
`ctx agent` overflow summaries and trace resolution show the same
dirty titles.

### B. `ctx drift` flags every CRLF file as a stale header

`drift.checkTemplateHeaders` (`internal/drift/check.go`) compares
`extractFirstComment(template)` with `extractFirstComment(live)`
byte-for-byte; `strings.TrimSpace` only trims the ends, not interior
newlines. A CRLF live file against the LF embedded template (or a
Windows-built binary, whose embedded templates are CRLF, against an LF
file) mismatches on line endings alone. Every context file warns
`comment header does not match template`, and `ctx init --reset`
cannot clear it: the rewritten file still differs from the embedded
bytes only in line endings. The false positives bury real drift.

### C. A date-only entry header is silently absorbed

`ParseEntryBlocks` starts a block only at a line matching the strict
`regex.EntryHeader`. An entry whose header lacks the time part —
`## [2026-09-26] Title`, from a hand edit or a pre-timestamp file — is
not a block start, so its heading and body become part of the
*previous* entry's block:

- `ctx disclosure inspect` never lists it;
- `ctx disclosure apply` moves it, inside the previous entry's span,
  into that entry's theme file — the wrong theme, with no trace in the
  plan the human approved.

`disclosure.Validate` already exists to refuse a structurally malformed
root before the pass mutates anything, but it only catches the case
where *no* staged entry parses (`ErrStagingUnparsable`). A malformed
header after at least one valid entry passes validation.

## Approach

Read-side fixes only. No writer changes, no file rewrites, no regex
changes, and no broadening of the accepted header format.

- **A.** Trim the captured title with `strings.TrimSpace` at both read
  points (`ParseEntryBlocks`, `ParseHeaders`). Block `Lines` stay
  verbatim, so `SplitStaging`'s byte-exact cuts and conservation are
  unaffected. Titles are already trimmed for convention sections
  (`conventionBlocks`), so every kind now yields clean titles.
- **B.** `extractFirstComment` normalizes `token.NewlineCRLF` to
  `token.NewlineLF` before extracting. Both sides of the comparison go
  through it, so the check becomes line-ending-insensitive while staying
  content-sensitive.
- **C.** `disclosure.Validate` refuses, for the timestamped kinds
  (learning, decision), any line outside an HTML comment that matches
  `regex.EntryHeading` (`^## \[`) but not `regex.EntryHeader`. The
  refusal is a typed error, `errDisc.MalformedEntryHeaderError`, naming
  the 1-based line number and the offending heading, and telling the
  user to give it a full `## [YYYY-MM-DD-HHMMSS] Title` timestamp.
  - HTML comments are skipped with the scanner's existing
    `htmlCommentSpans`: DECISIONS.md's template ships a
    `<!-- DECISION FORMATS ... -->` guide containing
    `## [YYYY-MM-DD] Decision Title`, which is an example, not an entry.
  - The check runs before `ErrStagingUnparsable`, and the typed error's
    `Is` matches `ErrStagingUnparsable`: a malformed header is the
    precise form of "staging could not be parsed into discrete
    entries". Callers matching the sentinel keep working; callers that
    want the location use `errors.AsType[*MalformedEntryHeaderError]`.
    A root whose *first* staged entry is malformed now gets the precise
    error too, instead of the location-free generic one.
  - Conventions are exempt: their sections open with `## ` and carry no
    timestamp, so `## [` is ordinary prose there.
  - `Apply` already calls `Validate` before any write, so a refused
    root is byte-identical. `Inspect` stays total (read-only, never
    fails), as documented.

## Tests

TDD, one failing test per defect before its fix:

- **A.** `internal/heading`: a CRLF case for `ParseHeaders` and
  `TestParseEntryBlocks_CRLFTitle`; both assert the title has no `"\r"`.
- **B.** `internal/drift`: `checkTemplateHeaders` against the real
  embedded template rendered with LF and with CRLF endings raises no
  `stale_header` warning; a genuinely edited comment still does. Both
  endings are exercised so the test fails before the fix on either an
  LF or a CRLF checkout.
- **C.** `internal/disclosure`:
  - a date-only header after a valid entry is refused with
    `MalformedEntryHeaderError` carrying the right line and heading, for
    both LF and CRLF content (the heading is reported without `"\r"`);
  - `errors.Is(err, ErrStagingUnparsable)` still holds;
  - the real embedded DECISIONS.md template (its commented
    `## [YYYY-MM-DD] Decision Title` example) plus a valid entry
    validates cleanly;
  - existing `TestValidate` cases, including "unparsable staging",
    are unchanged and pass.

## Acceptance

- On a CRLF LEARNINGS.md, `ctx disclosure inspect --json` titles
  contain no `"\r"`.
- Adding `## [2026-09-26] Legacy` below a valid entry makes
  `ctx disclosure apply` refuse with the line number and heading; the
  root is left untouched.
- `ctx drift` reports no `stale_header` for a CRLF context file whose
  header matches the template.
- `golangci-lint run` clean; no new test failures.

## Non-Goals

- Normalizing line endings at write time (`ctx init`, `add`, reindex).
- Accepting or auto-repairing date-only headers. The pass stays
  fail-loud with no auto-repair (progressive-disclosure Guards §4).
- Making `ctx disclosure inspect` refuse malformed roots; it stays a
  total, read-only view, and `apply` is the gate.
- Fenced-code awareness in the header scan; the disclosure scanners are
  deliberately fence-blind (see `conventionBlocks`).
