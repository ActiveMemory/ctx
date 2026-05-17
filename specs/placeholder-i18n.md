---
title: Placeholder set i18n + .ctxrc override
date: 2026-05-11
status: ready
owner: jose
scope: validate package + new asset bundle + .ctxrc accessor + Unicode-safe folding
design-ref: feedback from session 3c81f71b on commit 0ccc1a83 (validate body-flag refactor)
phase: 0 (prerequisite for any phase that ships locale-specific behaviour)
---

# Spec: Placeholder set i18n + .ctxrc override

Make the placeholder set used by `RejectPlaceholder` localizable and
user-extensible. Today the set is hardcoded English in
`internal/config/validate/placeholder.go`; users writing in any other
language can still slip "por definir", "iptal", or "TBD-yapılacak"
through the validator.

## Problem

Three concrete gaps:

1. **English-only set** baked as Go constants
   (`PlaceholderTBD`, `PlaceholderNA`, `PlaceholderSeeChat`, …). A
   Spanish, Turkish, German, or Japanese user has no shipped
   defense against their language's "to be defined" markers.
2. **No `.ctxrc` override hook**. Every other parser-vocabulary
   list in ctx (`session_prefixes`, `classify_rules`,
   `spec_signal_words`) is overridable. Placeholders are the
   outlier.
3. **Locale-naive case folding**. `RejectPlaceholder` does
   `strings.ToLower(trimmed)`. Today the set is pure ASCII so
   nothing breaks; the moment a Turkish or German placeholder
   lands in the YAML, byte-level `ToLower` will silently miss
   values typed with Unicode case differences:
   `strings.ToLower("İPTAL") == "i̇ptal"` (combining dot above),
   not `"iptal"`. Same trap with German `ß`/`SS`.

## Approach

Three coordinated moves, one commit:

1. **Extract defaults to an embedded YAML asset.** New file
   `internal/assets/commands/vocab/placeholders.en.yaml` (path
   subject to convention check — see Open Questions). Loaded at
   init time the same way `commands/text/*.yaml` is loaded. The
   `Placeholder*` Go consts in `internal/config/validate/placeholder.go`
   stay as identifier handles for tests and references; only the
   string values move to YAML.

2. **Add a `.ctxrc placeholders:` accessor with EXTEND semantics.**
   Modeled on `rc.SessionPrefixes()` but with the combine rule
   inverted:

   ```go
   // SessionPrefixes — REPLACE: empty user list → use defaults.
   // Placeholders   — EXTEND: user list is appended to defaults,
   //                  case-folded and de-duplicated.
   ```

   Reason for the difference: a Spanish user setting
   `placeholders: ["por definir"]` should still have `tbd`
   rejected — bilingual ("Tarzan Turkish": EN+TR intermingled in
   the same project) is the dominant case for this codebase, so
   replace would surprise. Replace can be reconsidered later via
   an opt-in flag if needed; extend is the safe default.

3. **Replace `strings.ToLower` with proper Unicode case folding.**
   Use `golang.org/x/text/cases` + `language.Und` for Unicode-aware
   folding that handles İ/I/i correctly across all locales.
   Applied at two sites: when loading user overrides from .ctxrc
   (normalize on ingest) and when comparing the input value
   against the set in `RejectPlaceholder`.

## Behavior

### Happy Path

User in a Turkish-language project adds to `.ctxrc`:

```yaml
placeholders:
  - "iptal"
  - "yapılacak"
  - "tbd-yapılacak"
```

`ctx decision add --rationale "İptal"` → rejected (case-folded
match against user-supplied `iptal`).
`ctx decision add --rationale "TBD"` → still rejected (shipped
default).
`ctx decision add --rationale "iptal edildi çünkü ..."` →
accepted (substring, not exact match — same rule as today).

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `.ctxrc` has `placeholders:` key but list is empty | Use shipped defaults only (no error) |
| `.ctxrc` user value duplicates a shipped default (`"tbd"`) | De-dupe silently after fold; no error |
| User value differs only by case (`"TBD"` vs `"tbd"`) | Same — de-dupe after fold |
| User value has surrounding whitespace | Trim on ingest; do not store raw |
| Malformed YAML (`placeholders: "foo"` instead of list) | RC loader rejects with typed error; do not silently ignore |
| Turkish dotted/dotless I (`İptal` typed; `iptal` in YAML) | Match (Unicode fold) |
| German sharp s (`STRASSE` typed; `strasse` in YAML) | Match (Unicode fold via `cases.Fold`) |
| Empty user value in list (`placeholders: ["", "tbd"]`) | Skip empty entries silently on ingest |

### Validation Rules

- Defaults YAML schema: top-level key `placeholders:`, list of
  non-empty strings, lowercase ASCII for the shipped `en` file.
- `.ctxrc` schema: same shape; values may be any Unicode.
- Validator pre-folds and trims both sides before exact-match
  comparison; substring matches remain accepted as legitimate
  prose.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Embedded YAML fails to parse at init | Panic with file path (build-time invariant) | Fix the YAML; CI catches via test |
| `.ctxrc placeholders:` is wrong type | `errCli.RCField("placeholders", "expected list of strings")` | User fixes .ctxrc |

## Interface

### Configuration

`.ctxrc` gains:

```yaml
placeholders:
  - "iptal"
  - "yapılacak"
```

No new CLI flags, no new commands. The validator change is
transparent — same rejection semantics, larger set of triggers.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/commands/vocab/placeholders.en.yaml` | NEW: list of shipped defaults |
| `internal/config/validate/placeholder.go` | Keep `Placeholder*` const identifiers; move string values to YAML; delete the `Placeholders` map |
| `internal/assets/read/vocab/vocab.go` | NEW: loader for the vocab namespace, modeled on `read/desc` |
| `internal/rc/rc.go` | NEW: `Placeholders() []string` accessor with EXTEND combine rule |
| `internal/rc/<schema>` | Add `Placeholders []string \`yaml:"placeholders"\`` to RC struct |
| `internal/validate/rejectplaceholder.go` | Switch from `cfgValidate.Placeholders` map lookup to `rc.Placeholders()` slice scan; replace `strings.ToLower` with `cases.Fold` |
| `internal/validate/rejectplaceholder_test.go` | Add Unicode-fold cases (İ/i, ß/ss); add .ctxrc-override test |
| `go.mod` / `go.sum` | Add `golang.org/x/text/cases` if not already a transitive |
| `docs/operations/configuration/ctxrc.md` (if exists) | Document the new `placeholders:` key |

### Key Functions

```go
// internal/rc/rc.go
func Placeholders() []string {
    defaults := vocab.Placeholders() // from embedded YAML
    user := RC().Placeholders
    return mergeFolded(defaults, user) // de-dupe after Fold
}

// internal/validate/rejectplaceholder.go
func RejectPlaceholder(flag, value string) error {
    trimmed := strings.TrimSpace(value)
    if trimmed == "" {
        return errCli.FlagEmpty(flag)
    }
    folded := cases.Fold().String(trimmed)
    for _, p := range rc.Placeholders() {
        if folded == p { // p is pre-folded at load time
            return errCli.FlagPlaceholder(flag, value)
        }
    }
    return nil
}
```

### Helpers to Reuse

- `internal/assets/embed.go` for the embed.FS pattern
- `internal/assets/read/desc/desc.go` for the lookup-map pattern
- `internal/rc/rc.go` `SessionPrefixes()` for the accessor shape
  (but invert the combine rule)

## Testing

- **Unit (`internal/validate/`)**: existing 4 tests stay; add cases
  for İ/i, ß/SS, mixed-script Tarzan Turkish, and one negative
  case (substring match still accepted).
- **Unit (`internal/rc/`)**: combine rule — user list extends
  defaults; user duplicates fold to the same key; empty user list
  → defaults only.
- **Integration**: write a temp `.ctxrc` with custom
  `placeholders:`, run `ctx decision add` with a value matching
  the user-supplied placeholder, assert non-zero exit and the
  flag-name in stderr.
- **Audit**: `desckey_namespace_test`-style audit that every YAML
  entry parses and every `Placeholder*` const has a corresponding
  YAML value in the `en` file.

## Non-Goals

- **Full error-message i18n.** That's the `internal/config/embed/text`
  package's job and is out of scope. Only the placeholder
  vocabulary moves; the `errCli.FlagEmpty` and
  `errCli.FlagPlaceholder` messages stay in their current form.
- **Per-flag placeholder sets.** One vocabulary applies to all
  body flags (decision/learning add today, more later). Per-flag
  customization can be added if a real need appears.
- **Shipping `tr` (or any other locale) defaults in this spec.**
  ctx has not shipped any locale-specific behavior yet (no `tr`
  files anywhere). Ship `en` only; the structure makes adding
  `tr.yaml`, `es.yaml`, etc. a copy-edit later.
- **Fuzzy / Levenshtein / stem matching.** Exact case-folded
  match after trim, same as today. Substring matches still
  accepted.
- **Replace semantics for `.ctxrc` overrides.** Extend only in
  v1. Reconsider via an explicit `placeholders_replace: true`
  toggle if a real need appears.

## Open Questions

- **Asset path.** Two reasonable homes: `internal/assets/commands/vocab/`
  (sibling to `commands/text/`) or a new
  `internal/assets/vocab/`. Pick whichever the existing convention
  audits prefer; if the audit is silent, default to
  `commands/vocab/` for consistency with the existing
  `commands/text/` and `commands/cmd/` neighbors.
- **Loader namespace.** Mirror `internal/assets/read/desc` as
  `internal/assets/read/vocab` or fold into `desc`? Separate
  package keeps the responsibility line clean (display text vs
  parser vocabulary).
