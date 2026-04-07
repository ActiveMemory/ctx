# JSONL Schema Validation

## Problem

ctx parses Claude Code's JSONL session files to produce journal entries.
This format is undocumented, unversioned, and changes whenever Anthropic
ships a new Claude Code release. Our parser (`internal/journal/parser/`)
uses Go's `json.Unmarshal` which silently ignores unknown fields and
zero-fills missing ones — meaning schema drift produces thinner journal
entries with no warnings, no errors, and no indication that data was lost.

This has already happened in the wild: the `slug` field was removed in
newer CC versions (noted in `claude.go:94`), and a third-party script
(`read-cycle-log.sh`) broke silently when CC changed record types from
`type: "message"` to `type: "user"|"assistant"`.

We need a declared schema we can validate against, so drift is detected
before users notice degraded journal quality.

## Approach

1. **Derive a schema** from empirical data (all existing JSONL files
   across CC versions 2.1.x) and CC source types (`types/logs.ts`,
   `types/message.js`).

2. **Embed the schema** in the ctx binary as a Go struct with known
   field sets, known record types, and known content block types.

3. **Validate at import time** (mechanism A) — during `ctx journal import`,
   run each JSONL line through the schema. Warn on unknown fields,
   missing expected fields, and unknown content block types. Never block
   the import.

4. **Standalone check command** (mechanism B) — `ctx journal schema check`
   scans JSONL files without importing. Write drift report to
   `.context/reports/schema-drift.md` only when drift is detected. Delete
   the report when drift is resolved.

5. **Nightly runner** — schedule mechanism B to run daily (midnight),
   producing a deterministic, diffable report. Since drift reports contain
   no sensitive data (just field names and counts), the report is safe to
   commit to git. The report file appearing in `git status` is itself the
   signal that something changed.

## Behavior

### Happy Path

1. User runs `ctx journal import` (or nightly cron triggers
   `ctx journal schema check`).
2. Parser reads each JSONL line and validates against embedded schema.
3. All fields, record types, and content block types are recognized.
4. Import proceeds normally. No drift report written (or existing one
   deleted if drift was previously resolved).

### Drift Detected Path

1. Parser encounters an unknown top-level key (e.g., `compactMetadata`),
   unknown record type, unknown content block type, or a previously
   required field that is now absent.
2. Validation accumulates drift findings per-file (not per-line) to
   avoid flooding output.
3. At end of import, a summary is printed:

   ```
   ⚠ Schema drift detected in 3 session files:
     Unknown fields: compactMetadata, isCompactSummary (2 files)
     Unknown block types: code_execution (1 file)
     Missing expected fields: slug (3 files)
   Run `ctx journal schema check` for full report.
   ```

4. Import still completes — drift warnings are informational.
5. For mechanism B, drift report is written to
   `.context/reports/schema-drift.md`.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Zero JSONL files found | Print "no session files found", exit 0 |
| File contains only non-message records (progress, file-history-snapshot) | Skip file silently — no drift reported for records we intentionally ignore |
| Mixed versions in one file (e.g., session started on 2.1.23, resumed on 2.1.25) | Report per-file, noting version range |
| Field present in some lines, absent in others within same file | Report as "intermittent field: X (present in N/M lines)" — this is normal for optional fields, only flag if a previously-always-present field becomes intermittent |
| Malformed JSON line | Count and report separately from schema drift ("N malformed lines skipped") |
| Schema itself is outdated (all files have "unknown" fields) | Suggest schema update: "All files contain field X — consider adding to schema" |
| `.context/reports/` directory does not exist | Create it |
| Drift report exists but no drift detected on this run | Delete the report file |

### Validation Rules

The schema defines three categories of expectations:

**Record types** (top-level `type` field):
- Known message types: `user`, `assistant`
- Known metadata types: `last-prompt`, `custom-title`, `ai-title`,
  `attachment`, `permission-mode`, `agent-name`, `agent-color`,
  `agent-setting`, `tag`, `pr-link`, `mode`, `worktree-state`,
  `content-replacement`, `speculation-accept`, `task-summary`
- Known infrastructure types: `progress`, `file-history-snapshot`,
  `attribution-snapshot`, `system`, `summary`, `queue-operation`
- Any other value triggers drift warning

**Top-level fields** per record type:
- Required for `user`/`assistant`: `uuid`, `parentUuid`, `sessionId`,
  `timestamp`, `type`, `cwd`, `version`, `message`, `isSidechain`,
  `userType`
- Optional (known): `gitBranch`, `slug`, `requestId`, `thinkingMetadata`,
  `todos`, `permissionMode`, `logicalParentUuid`, `isMeta`,
  `compactMetadata`, `isVisibleInTranscriptOnly`, `isCompactSummary`,
  `agentId`, `teamName`, `agentName`, `agentColor`, `promptId`,
  `entrypoint`, `agentSetting`, `sourceToolAssistantUUID`,
  `toolUseResult`, `sourceToolUseID`, `origin`, `planContent`,
  `isApiErrorMessage`, `error`, `apiError`
- Unknown fields are flagged

**Content block types** (`message.content[].type`):
- Known: `text`, `thinking`, `tool_use`, `tool_result`
- Known but not parsed: `server_tool_use`, `mcp_tool_use`,
  `mcp_tool_result`, `code_execution_tool_result`, `container_upload`
- Unknown types are flagged

## Historical Baseline (2026-01-14 → 2026-04-07)

Derived from scanning 2187 files, 700,192 lines, 299,423 message
records across 54 CC versions (2.1.2 → 2.1.92).

**Schema stability:**
- Zero missing required fields — all 10 core fields at 100% presence
- Zero unknown content block types — `text`, `thinking`, `tool_use`,
  `tool_result` have been stable across all versions
- `gitBranch` is effectively required (100%) despite being optional
- `slug` is at 98.7% — gradual deprecation, not a cliff

**Record types — fully accounted for:**
- Zero unknown record types across the entire history
- 6 metadata types we weren't tracking but now declare:
  `last-prompt` (123), `attachment` (52), `permission-mode` (47),
  `custom-title` (16), `agent-name` (16), `agent-color` (10)
- All are CC session state records, not message content

**Envelope fields we weren't tracking:**
- `sourceToolAssistantUUID` on 38.3% of user records (agent linkage)
- `toolUseResult` on 23.3% of user records (inline tool results)
- `entrypoint` on 36.5% of messages (CLI/SDK/IDE source)
- Error metadata (`isApiErrorMessage`, `error`, `apiError`) on
  assistant records
- Plan mode data (`planContent`) on user records
- Rich `system` record fields: `hookCount`, `hookInfos`, `hookErrors`,
  `stopReason`, retry metadata

**Implication:** The format is remarkably stable. Content blocks
(where journal text comes from) haven't drifted at all. Envelope
fields expand frequently but core fields never disappear. The drift
that exists is in optional metadata we could capture for richer
journal entries.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Cannot read JSONL file (permissions) | `"cannot read session file: {path}: {err}"` | Fix permissions, re-run |
| Schema asset missing from binary | `"internal error: embedded schema not found"` | Rebuild ctx |
| Report directory not writable | `"cannot write drift report: {err}"` | Check `.context/reports/` permissions |

## Interface

### CLI

```
ctx journal schema check [--dir path] [--all-projects] [--quiet]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | CC project sessions dir | Directory to scan |
| `--all-projects` | bool | false | Scan all CC project directories |
| `--quiet` | bool | false | Exit code only (0 = clean, 1 = drift) |

```
ctx journal schema dump
```

Emits the embedded schema as YAML/JSON for inspection.

### Import Integration

`ctx journal import` gains drift validation automatically — no new
flags. Warnings printed to stderr after import summary.

### Nightly Runner

A cron entry or systemd timer running:

```bash
ctx journal schema check --all-projects
```

The command writes `.context/reports/schema-drift.md` only when drift
is detected. The report is deterministic and diffable (sorted fields,
stable counts).

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/journal/schema/schema.go` | Schema definition: known fields, types, block types per record type |
| `internal/journal/schema/validate.go` | Validate a parsed JSONL line against schema, accumulate findings |
| `internal/journal/schema/report.go` | Format drift findings as Markdown report |
| `internal/journal/schema/embed.go` | Embed schema version metadata (date, CC version range tested) |
| `internal/journal/parser/claude.go` | Call schema validation during parse, collect findings |
| `internal/cli/journal/cmd/schema/cmd.go` | `ctx journal schema check` and `ctx journal schema dump` commands |
| `internal/cli/journal/cmd/importer/run.go` | Print drift summary after import |
| `internal/journal/schema/testdata/` | Golden JSONL fixtures from real sessions |

### Key Functions

```go
// Schema declares the expected shape of Claude Code JSONL records.
type Schema struct {
    Version          string
    TestedCCVersions []string
    RecordTypes      map[string]RecordSchema
    BlockTypes       map[string]BlockSchema
}

// RecordSchema declares expected fields for a record type.
type RecordSchema struct {
    Required []string
    Optional []string
}

// Finding represents a single schema drift observation.
type Finding struct {
    Type     FindingType // UnknownField, MissingField, UnknownRecordType, UnknownBlockType
    Name     string      // The field/type name
    Files    []string    // Which files exhibited this
    Count    int         // How many lines
}

// Validate checks a raw JSONL line against the schema.
func (s *Schema) Validate(line json.RawMessage, recordType string) []Finding

// Report formats accumulated findings as a Markdown drift report.
func Report(findings []Finding, scanMeta ScanMeta) string
```

### Helpers to Reuse

- `internal/journal/parser/claude.go` — already parses JSONL lines,
  validation hooks into the same loop
- `internal/journal/parser/types.go` — `claudeRawMessage` struct
  defines what we currently expect
- `internal/io/` — `SafeOpenUserFile`, `SafeWriteFile`
- `internal/config/claude/` — role constants (`RoleUser`, `RoleAssistant`)

## Configuration

No `.ctxrc` keys needed. The schema is embedded at build time.

The nightly runner is configured outside ctx (cron/systemd). Example
crontab entry:

```
0 0 * * * cd /home/jose/WORKSPACE/ctx && ctx journal schema check --all-projects
```

## Testing

- **Unit**: Validate known-good JSONL lines return zero findings.
  Validate lines with unknown fields/types return correct findings.
  Validate report formatting is deterministic.
- **Golden fixtures**: Pin real JSONL snippets (one per CC version we've
  seen: 2.1.23, 2.1.25) as testdata. Parser + schema tests run against
  these.
- **Integration**: Run `ctx journal schema check` against the test
  fixtures directory, verify exit code and report content.
- **Edge cases**: Empty file, file with only `progress` records, mixed
  versions, malformed lines interspersed with valid ones.

## Non-Goals

- Validating tool_use input shapes (too tool-specific, changes per tool)
- Parsing or validating `progress`, `file-history-snapshot`, or other
  non-message record content (we skip these intentionally)
- Full JSON Schema (RFC 8927) compliance — a Go struct with field sets
  is sufficient and more maintainable
- Upstream engagement with Anthropic about schema stability
- Blocking imports on drift — always warn, never fail
- Validating message *content* semantics (e.g., that tool_result
  references a valid tool_use ID)
