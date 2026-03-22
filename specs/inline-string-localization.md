---
title: Externalize inline format strings and magic numbers
date: 2026-03-22
status: ready
---

# Inline String and Magic Number Externalization

## Problem

Multiple packages use inline format strings (e.g., `"%dm"`, `"%dB"`,
`"graph TD\n"`) and magic numbers (e.g., 1000, 1024, 50) instead of
routing through assets/config. This blocks localization and scatters
formatting conventions across the codebase.

## Scope

### Group 1: Duplicate formatting functions (consolidate to internal/format)

Three packages implement the same duration/token/size formatting with
different function names and slightly different thresholds:

| Package | Function | Formats |
|---------|----------|---------|
| `recall/core` | `FormatDuration` | `<1m`, `5m`, `1h30m` |
| `system/core` | `FormatAge` (prune.go) | `5m`, `3h`, `2d` |
| `system/core` | `FormatTokenCount` | `500`, `1.2k`, `52k` |
| `system/core` | `FormatWindowSize` | `200k` |
| `recall/core` | `FormatTokens` | `500`, `1.5K`, `2.3M` |
| `journal/core` | `FormatSize` | `512B`, `1.5KB`, `2.3MB` |
| `format` | `Number` | `500`, `1,500` |
| `format` | `Bytes` | `1.5 KB` |

The `format` package already has `TimeAgo`, `Duration`, `Bytes`,
`Number` — all using desc.Text(). The domain packages should delegate
to `format` or share its constants, not reimplement.

**Action**: Add text keys for the compact format suffixes (`m`, `h`,
`d`, `k`, `K`, `M`, `B`, `KB`, `MB`). Add config constants for SI/IEC
thresholds (1000, 1024). Consolidate where callers produce identical
output.

### Group 2: dep/core/format.go (Mermaid and table output)

| Line | String | Fix |
|------|--------|-----|
| 39 | `"graph TD\n"` | text key (protocol format, but externalizable) |
| 49 | `"    %s[\"%s\"] --> %s[\"%s\"]\n"` | text key |
| 65 | `"Package"`, `"Imports"` | text keys (table headers) |
| 66 | `strings.Repeat("-", 50)`, `strings.Repeat("-", 30)` | config constants for widths |
| 65,70 | `"%-50s %s\n"` | text key with config width |

**Action**: Create `config/dep/` constants for column widths. Add text
keys in `text/dep.go` + `write.yaml` for headers and format strings.

### Group 3: system/core/state.go (log format)

| Line | Issue | Fix |
|------|-------|-----|
| 82 | `len(short) > 8` | Use `journal.SessionIDShortLen` (already exists) |
| 86 | `"[%s] [session:%s] %s\n"` | text key |
| 87 | `"2006-01-02 15:04:05"` | Use `time2.DateTimePreciseFormat` (already exists) |

### Group 4: initialize backup format

| File | Line | String | Fix |
|------|------|--------|-----|
| `initialize/core/merge.go` | 40 | `"%s.%d.bak"` | text key or config constant |
| `initialize/core/plan.go` | 81,130 | `"%s.%d.bak"` | same — deduplicate via shared helper |

**Action**: Add backup format constant to `config/file/` or
`config/archive/`. Both merge.go and plan.go already call
`backupFile()` in merge.go — plan.go should too (DRY).

### Group 5: memory/state.go format strings

| Line | String | Fix |
|------|--------|-----|
| 74 | `fmt.Sprintf("%x", h[:8])` | Config constant for hash prefix length |
| 92 | `"%s:%s:%s"` (hash:target:date) | Config constant for separator |

### Group 6: ctximport inline error string

| File | Line | String | Fix |
|------|------|--------|-----|
| `write/ctximport/import.go` | 105 | `"  Error promoting to %s: %v"` | text key |

### Group 7: format.go magic numbers

`format.Number` and `format.Bytes` use inline `1000` and `1024`.
These are the canonical formatting functions — they should define or
reference shared constants.

**Action**: Add `config/format/` package with `SIThreshold = 1000`
and `IECUnit = 1024`. Reference from all formatting functions.

### Group 8: resources.go format widths

| Line | Number | Fix |
|------|--------|-----|
| 87 | `7` (label width) | `config/stats.ResourcesLabelWidth` |
| 113+ | `5` (value width) | `config/stats.ResourcesValueWidth` |

Check: `config/stats.ResourcesStatusCol` already exists for column
alignment. Add label/value widths alongside it.

## Execution order

1. **Config constants** — create `config/format/` and add missing
   constants to existing config packages
2. **Text keys** — add YAML entries and Go constants for format strings
3. **system/core/state.go** — quick fix using existing constants
4. **dep/core/format.go** — text keys + config widths
5. **Format consolidation** — wire domain formatters to shared constants
6. **initialize backup dedup** — consolidate plan.go backup to use merge.go helper
7. **ctximport error** — add text key
8. **memory/state.go** — config constants for hash format

## Non-goals

- Extracting Go format verbs (`%d`, `%s`, `%.1f`) — these are
  language-level primitives, not user-facing strings
- Extracting time.Duration constants (time.Hour, time.Minute) — these
  are stdlib constants
- Extracting percentage multiplier (100) — mathematical constant
- Changing the `format` package's `Bytes()` algorithm — only
  externalizing its threshold constant
