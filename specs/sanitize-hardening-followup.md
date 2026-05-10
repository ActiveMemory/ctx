---
title: Sanitize hardening follow-up (UTF-8 truncation, Zl/Zp, opts caps)
date: 2026-05-10
status: ready
owner: jose
scope: bug fix — three surgical fixes to internal/sanitize and MCP extract
related:
  - specs/future-complete/context-hub.md
prior:
  - PR #76 (MCP-SAN + MCP-COV hardening) — landed at cf097a69
---

# Spec: Sanitize Hardening Follow-up

PR #76 added the `internal/sanitize` package and applied it across the
MCP server's untrusted input surface. Three issues surfaced during
post-merge review. This spec records the corrective fixes.

## Problem

### Issue 1 — Byte-level truncation can split UTF-8 runes

`internal/sanitize/truncate.go:18-23`, `reflect.go:21-23`, and
`path.go:32-34` all do `s = s[:maxLen]` on bytes. If the cut lands
inside a multi-byte UTF-8 rune, the result contains an invalid
trailing byte sequence. Existing tests only cover ASCII so the bug
is silently uncovered. Downstream consumers (JSON encoders, file
writers, log lines) handle invalid UTF-8 inconsistently — some
escape, some replace with U+FFFD, some pass through.

### Issue 2 — `StripControl` misses U+2028 and U+2029

`internal/sanitize/content.go:50-62` filters via `unicode.IsControl`.
That predicate returns **false** for U+2028 (LINE SEPARATOR) and
U+2029 (PARAGRAPH SEPARATOR), which are Unicode category `Zl`/`Zp`,
not `Cc`. Since `Content` writes into Markdown that some renderers
parse as line breaks, these can still inject visual newlines.

### Issue 3 — Secondary opts fields have no length cap

`extract.EntryArgs` enforces `MaxContentLen` on the primary
`content` field. `extract.SanitizedOpts` (extract.go:100-111)
sanitizes `Context`, `Rationale`, `Consequence`, `Lesson`, and
`Application` but applies no length limit. An attacker can send
10 MB in `rationale` and it will be sanitized and written.

## Approach

### Issue 1 — Rune-safe truncation

`truncate(s, maxLen)` cuts at `maxLen` bytes, then backs up to a
rune-start boundary using `utf8.RuneStart`. All three internal call
sites (`Reflect`, `SessionID`, internal helpers) delegate to this
shared helper instead of duplicating the slice. The function
remains unexported.

### Issue 2 — Explicit Zl/Zp handling

Add an explicit `r == ' ' || r == ' '` check to the
`StripControl` rune predicate before falling through to
`unicode.IsControl`. Both runes get dropped (return `-1` from the
`strings.Map` callback).

### Issue 3 — Length cap for opts fields

Add `MaxOptsFieldLen` constant (4 KB) to `internal/config/mcp/cfg`.
Change `SanitizedOpts` signature to `(entity.EntryOpts, error)`.
Reject any opts field exceeding the cap via `errMcp.InputTooLong`,
with the field name in the error. Update both call sites in
`internal/mcp/server/route/tool/tool.go` (`add`, `watchUpdate`) to
propagate the error to the MCP client.

Rationale for 4 KB rather than `MaxContentLen` (32 KB): the secondary
fields are qualifiers, not primary content. A tighter cap reduces
abuse surface and pushes large prose into the primary `content`
field where it belongs.

## Behavior

### Happy path (Issue 1)
- ASCII input within limit: unchanged.
- Multi-byte UTF-8 input within limit: unchanged.
- Multi-byte UTF-8 input exceeding limit where the cut lands on a
  rune start: result is exactly `maxLen` bytes, valid UTF-8.
- Multi-byte UTF-8 input exceeding limit where the cut lands inside
  a rune: result is fewer than `maxLen` bytes, terminates on a
  rune boundary, no trailing invalid bytes.

### Happy path (Issue 2)
- Input containing ` ` or ` ` is stripped of both runes.
- Existing behavior for tab, LF, CR, and `Cc` control chars unchanged.

### Happy path (Issue 3)
- Opts field within cap: passes through sanitization unchanged.
- Opts field exceeding cap: caller receives `InputTooLong(field, MaxOptsFieldLen)`.
  Field name is the canonical MCP arg key (e.g., `"rationale"`).

### Edge cases

| Case | Expected |
|------|----------|
| `truncate("", 100)` | `""` |
| `truncate("a", 0)` | `"a"` (no truncation when maxLen ≤ 0) |
| `truncate(<one 4-byte rune>, 3)` | `""` (back up below the rune start) |
| `StripControl(" ")` | `""` |
| `SanitizedOpts` with empty opts | no error |
| `SanitizedOpts` with `rationale = strings.Repeat("a", 4097)` | `InputTooLong("rationale", 4096)` error |

## Interface

No public API changes to the `sanitize` package — `truncate` stays
unexported.

`extract.SanitizedOpts` signature changes:

```go
// Before:
func SanitizedOpts(args map[string]interface{}) entity.EntryOpts

// After:
func SanitizedOpts(args map[string]interface{}) (entity.EntryOpts, error)
```

## Implementation

### Files to modify

| File | Change |
|------|--------|
| `internal/sanitize/truncate.go` | Rune-safe via `utf8.RuneStart` |
| `internal/sanitize/reflect.go` | Delegate to `truncate` |
| `internal/sanitize/path.go` | Delegate to `truncate` |
| `internal/sanitize/content.go` | Add Zl/Zp check in `StripControl` |
| `internal/sanitize/sanitize_test.go` | New tests for the above |
| `internal/config/mcp/cfg/config.go` | Add `MaxOptsFieldLen = 4_000` |
| `internal/mcp/server/extract/extract.go` | New error-returning `SanitizedOpts` |
| `internal/mcp/server/extract/extract_test.go` | New length-cap tests |
| `internal/mcp/server/route/tool/tool.go` | Propagate the new error |

## Testing

- Unit: rune-boundary cuts at 2-, 3-, and 4-byte runes.
- Unit: Zl/Zp stripping (single, mixed, repeated).
- Unit: opts-field length rejection per field.
- Integration: existing MCP server tests must continue to pass.

## Non-Goals

- Adding length caps for `Branch`, `Commit`, `Priority`, `Section`,
  `SessionID` (already capped via `sanitize.SessionID`). These are
  outside the review scope and will be tracked separately if needed.
- Adding fuzz tests (suggested as a nit, deferred).
- Refactoring `Content`'s overlapping regex passes (separate nit).
