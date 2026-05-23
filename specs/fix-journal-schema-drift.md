# Absorb post-1.0 Claude Code Journal Field Drift

`ctx journal schema check` reported five unknown fields
in user JSONL records:

- `apiErrorStatus`
- `attributionPlugin`
- `attributionSkill`
- `errorDetails`
- `interruptedMessageId`

All five are optional fields Claude Code added to message
records in versions between 2.1.92 (end of the schema's
1.0.0 declared range) and 2.1.150 (the version users are
running today). The schema was ~60 minor versions stale
relative to the surface it was validating against.

## Problem

The 1.0.0 schema baseline was empirically derived from CC
versions 2.1.2–2.1.92 (`internal/config/schema/report.go`
constants `Version` and `CCVersionRange`). Field
expectations live in `internal/config/schema/field.go`
as `RequiredFields` and `OptionalFields` string slices.
When CC ships new optional fields, the schema check
flags them as drift until they're added to
`OptionalFields`. The check is a *passive* drift detector
— it does not auto-extend the schema — so each new field
requires a one-line addition plus a version bump.

The lag here was ~6 weeks of CC release cadence. Symptoms:

```
Schema drift detected in 27 file(s):
  Unknown fields: apiErrorStatus, attributionPlugin,
  attributionSkill, errorDetails, interruptedMessageId
```

End-users see this every time they run `ctx agent`,
`ctx journal source`, or any code path that triggers the
schema check. The noise is harmless (the parser still
ingests the records) but trains operators to ignore
schema-check output, which is exactly the failure mode
the check is meant to prevent.

## Solution

1. **Append the five fields to `OptionalFields`** in
   `internal/config/schema/field.go`, grouped by concept
   (error context, attribution, message metadata) with
   brief version-provenance comments.
2. **Bump `Version` 1.0.0 → 1.1.0** — MINOR per semver,
   because adding optional fields is backwards-compatible
   (records written against 1.0.0 still validate cleanly).
3. **Bump `CCVersionRange` end 2.1.92 → 2.1.150** to
   reflect the version range the new field set covers.
4. **Pin the new fields with a regression test**
   (`TestKnownField_PostV1FieldDrift` in
   `internal/journal/schema/schema_test.go`) so a future
   refactor that drops one of these from `OptionalFields`
   fires immediately rather than re-surfacing as noise in
   user terminals.
5. **Gitignore generated reports**
   (`.context/reports/`). The schema check writes a
   detailed report under `.context/reports/schema-drift.md`
   that's operator-local generated output, matching the
   established ignore pattern for `.context/logs/`,
   `.context/state/`, `.context/journal/`, etc. Without
   the ignore, a `git add -A` could leak the report into
   commits.

## Policy: additive field drift handling

For future CC field-drift fixes that match the same
shape (new optional field appears, schema check flags it):

- **MINOR bump** if the addition is backwards-compatible
  (the common case for new optional fields).
- **MAJOR bump** if a field is removed or made required
  retroactively (rare; would also need migration logic).
- **PATCH bump** for fixes to existing field handling
  that don't change the declared field set.
- **CCVersionRange end** always bumps to the latest CC
  version observed in the user-submitted drift report,
  even if we can't pin the exact CC version that
  introduced each field. (We don't run CC's full release
  archaeology; the range marks "we've seen records from
  CC versions up to N and they validate cleanly.")
- **Regression test** for each new field added: pin
  `KnownField(rt, field) == true` so a future "let's
  tidy" refactor can't silently drop the field and
  re-introduce the drift.

## Out of Scope

- **Auto-extending schema** on drift detection. The
  passive design is intentional — operators see drift,
  decide whether to upgrade ctx or pin to a specific CC
  version, and the maintainer team owns the per-version
  evaluation. Auto-extension would silently mask
  potentially-meaningful CC behavior changes.
- **Field-level type validation** (e.g., asserting
  `apiErrorStatus` is a string, not a number). The
  current check is field-presence only. A typed
  validation pass is a separate concern.
- **Schema-version negotiation** in the parser. The
  parser is permissive of unknown fields; only the
  `check` command flags them. No negotiation needed.

## Verification

- `go test ./internal/journal/schema/...` —
  `TestKnownField_PostV1FieldDrift` passes for each of
  the five new fields against both `user` and
  `assistant` record types.
- After install, `ctx journal schema check` reports zero
  drift on user journal directories.
- `make lint` clean (one line-length wrap was needed to
  fit a doc comment under 80 chars).
