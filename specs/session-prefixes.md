# Spec: Configurable Session Header Prefixes

Replace the hardcoded `session_prefix` / `session_prefix_alt` pair with a
user-extensible list of recognized session header prefixes.

## Problem

The Markdown session parser recognizes two session header prefixes:

- `Session:` (English, `parser.session_prefix` in text.yaml)
- `Oturum:` (Turkish, `parser.session_prefix_alt` in text.yaml)

This design has multiple problems:

1. **Doesn't scale**: adding Japanese requires `session_prefix_alt_2`
2. **Requires code changes**: new languages need a new constant in
   embed.go, a new YAML key, and a code change in the parser loop
3. **Wrong abstraction layer**: these aren't UI strings — they're
   parser recognition patterns. They don't belong in the i18n text
   asset system
4. **Conflates interface language and content language**: a user may
   prefer English UI but need to parse session files written in
   Turkish, Japanese, or Spanish by colleagues or other AI tools

## Solution

Move session header prefixes from text.yaml to `.ctxrc` as a
user-configurable list with built-in defaults.

### .ctxrc Configuration

```yaml
session_prefixes:
  - "Session:"     # English
  - "Oturum:"      # Turkish
```

When absent, defaults to `["Session:"]`. Users extend the list
without code changes:

```yaml
session_prefixes:
  - "Session:"
  - "Oturum:"
  - "セッション:"    # Japanese
  - "Sesión:"       # Spanish
```

### Internal Design

**New in `internal/rc`:**

- Add `SessionPrefixes []string` field to `CtxRC` struct
  (`yaml:"session_prefixes"`)
- Add `SessionPrefixes()` accessor function
- Default: `[]string{"Session:"}`

**New in `internal/config/parser`:**

- Add `DefaultSessionPrefixes` variable with the two built-in defaults
- The parser package owns the domain default; rc uses it as fallback

**Changes in `internal/recall/parser/markdown.go`:**

- `isSessionHeader()` and `parseSessionHeader()` call
  `rc.SessionPrefixes()` instead of `assets.TextDesc()`
- The loop over prefixes remains — it just reads from the configurable
  list instead of two hardcoded asset lookups

**Removals:**

- `TextDescKeyParserSessionPrefix` constant from `internal/assets/embed.go`
- `TextDescKeyParserSessionPrefixAlt` constant from `internal/assets/embed.go`
- `parser.session_prefix` entry from `text.yaml`
- `parser.session_prefix_alt` entry from `text.yaml`

### Data Flow

```
.ctxrc (user config)
  ↓
rc.SessionPrefixes()  ← falls back to config/parser.DefaultSessionPrefixes
  ↓
markdown.go: isSessionHeader() / parseSessionHeader()
  ↓
Loop over []string prefixes
```

### Test Changes

- `markdown_test.go`: Turkish test cases remain valid (default includes
  "Oturum:"). Add test case for a custom prefix (e.g., "セッション:")
  using `rc.SetOverride()` or test helper.
- Add unit test for `rc.SessionPrefixes()` default and override behavior.

### Edge Cases

- **Empty list in .ctxrc**: treat as "use defaults" — same pattern as
  `priority_order: []` falling back to `config.ReadOrder`
- **Duplicate prefixes**: harmless — the first match wins
- **Prefix with/without colon**: user's responsibility — document that
  prefixes should include the colon
- **Mixed-language directories**: fully supported — each file is matched
  independently against all prefixes

## Non-Goals

- Full i18n/locale system (out of scope — this is parser vocabulary)
- Changing how ctx *writes* session headers (ctx only parses, never
  generates session headers)
- Removing the date-only header pattern (`# 2026-01-15 — Topic`) —
  that's orthogonal and stays unchanged

## Documentation Updates

1. **docs/cli/index.md**: Add `session_prefixes` to .ctxrc reference
   table
2. **docs/recipes/multi-tool-setup.md**: Add section on multilingual
   session parsing
3. **docs/home/contributing.md**: Already has "How To Add a Session
   Parser" — verify prefix extensibility is mentioned
4. **ARCHITECTURE.md**: Update "Extensible Session Parsing" section
   to mention configurable prefixes
