---
title: Migrate fmt.Fprint* output calls to write/ packages
date: 2026-03-22
status: ready
prerequisite: cmd-print-to-write-migration.md
---

# fmt.Fprint* Output Migration

## Problem

After migrating all `cmd.Print*` calls to write/, one `fmt.Fprint*`
call remains that writes directly to an output stream (`os.Stderr`)
instead of routing through write/:

- `internal/cli/pad/core/store.go:87` — `fmt.Fprintln(os.Stderr, ...)`
  for key creation notice

Additionally, `write/pad.KeyCreated` exists as dead code — it was
created for this purpose but never wired up.

## Not in scope

`fmt.Fprintf` calls that write to `strings.Builder` are string
construction, not output. These are in:
- `dep/core/format.go` — building Mermaid markup
- `journal/core/section.go` — building index/site content
- `recall/core/frontmatter.go` — building YAML frontmatter
- `system/cmd/check_freshness/run.go` — building warning text

`fmt.Fprintln` calls that write to an `io.Writer` parameter already
follow the migrated pattern (stats StreamStats, backup addEntry).

## Solution

### pad/core/store.go → write/pad.KeyCreated

Thread `io.Writer` through the call chain so the output goes through
write/:

1. Change `EnsureKey() error` → `EnsureKey(w io.Writer) error`
2. Replace `fmt.Fprintln(os.Stderr, ...)` with
   `writePad.KeyCreatedW(w, path)` (new write/ function taking
   `io.Writer`)
3. Change `WriteEntries(entries []string) error` →
   `WriteEntries(w io.Writer, entries []string) error`
4. Update all 12 callers of `WriteEntries` to pass
   `cmd.ErrOrStderr()`
5. Update `write/pad.KeyCreated` to take `io.Writer` instead of
   `*cobra.Command`, or add a `KeyCreatedW(w io.Writer, path string)`
   variant and have the `*cobra.Command` version delegate to it
