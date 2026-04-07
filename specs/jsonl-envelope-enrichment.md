# JSONL Envelope Enrichment

## Problem

Claude Code JSONL session files contain rich envelope metadata beyond
the content blocks we currently parse. Our parser extracts message
text, thinking blocks, and tool use/results — but ignores envelope
fields that carry meaningful context about how the session unfolded.

This analysis was triggered by the schema validation work
(see `specs/jsonl-schema-validation.md`) which revealed that 6
envelope fields appear on 0.03–38.3% of message records and carry
data that could improve journal quality.

## Analysis: Envelope Fields We Don't Parse

### Field Inventory

Derived from scanning 2186 files, 693,149 lines across CC versions
2.1.2 → 2.1.92 (2026-01-14 → 2026-04-07).

| Field | Record type | Frequency | Content |
|-------|-------------|-----------|---------|
| `planContent` | user | 88 records | Full plan-mode document text (design specs, implementation plans) |
| `isApiErrorMessage` | assistant | 117 records | Boolean flag marking API error responses (rate limits, overload) |
| `sourceToolAssistantUUID` | user | 38.3% (114,659) | UUID linking tool result back to the assistant message that triggered it |
| `toolUseResult` | user | 23.3% (69,841) | CC-level tool error strings (EISDIR, file not found) — separate from API tool_result blocks |
| `entrypoint` | user/assistant | 36.5% (109,318) | How CC was launched: `cli`, `ide`, `sdk-ts`, `sdk-py` |
| `origin` | user | 106 records | Message injection source: `{"kind": "task-notification"}` |
| `error` / `apiError` | assistant | 106 / 6 records | Error objects on failed API responses |

### Field Value Samples

**`planContent`** (user records):
```
"# Plan: Add Quick Reference Index to DECISIONS.md\n\n## Goal\n
Add an auto-maintained index at the top of DECISIONS.md..."
```
The full plan document from plan mode. Currently our journal entries
show the conversation *about* the plan but lose the plan artifact itself.

**`toolUseResult`** (user records):
```
"Error: EISDIR: illegal operation on a directory, read"
"Error: File does not exist."
```
CC-level tool execution errors, attached directly to the user record
envelope. Distinct from `tool_result` content blocks which come from
the API. These represent failures that happen before the API sees them.

**`isApiErrorMessage`** (assistant records):
```
true
```
Marks assistant messages that are API error responses rather than real
completions. These render as regular assistant text in our journal but
are actually retry noise (rate limits, 500s, overloaded errors).

**`sourceToolAssistantUUID`** (user records):
```
"2e115417-a288-41d2-ad2c-1f07a9b5037c"
```
Links a user message (carrying a tool result) back to the assistant
message that requested the tool call. This is the threading chain
that would let us reconstruct call/response trees for subagent sessions.
Currently our subagent journal entries are flat sequences.

**`entrypoint`** (user/assistant records):
```
"cli"
```
Always `cli` in current data. Would differentiate sessions opened via
VS Code extension (`ide`), Agent SDK (`sdk-ts`, `sdk-py`), or web.
Useful as journal frontmatter metadata.

**`origin`** (user records):
```
{"kind": "task-notification"}
```
Distinguishes user-typed messages from system-injected ones (task
notifications, hook outputs). Helps tell "user asked this" from
"system nudged this" in journal entries.

**`attachment`** (attachment records):
```
{"type": "deferred_tools_delta", "addedNames": ["AskUserQuestion", ...]}
{"type": "companion_intro", "name": "Ingot", "species": "blob"}
```
Tool availability changes and companion metadata. No journal value.

## Approach

Capture all fields in the parser (they're just struct fields), but
only act on high-value fields in the near term. Three tiers:

### Tier 1: Act on now (improves journal quality)

- **`planContent`**: Render in journal entries as a collapsible
  "Plan" section when present. Plans are design artifacts worth
  preserving.
- **`isApiErrorMessage`**: Filter or collapse API error messages in
  journal output. These are retry noise, not conversation content.
  Mark them in the journal as `<!-- api error -->` or collapse into
  a single line: "⚠ API error (retried)".

### Tier 2: Capture as metadata (journal frontmatter)

- **`entrypoint`**: Add to `Session` entity and journal frontmatter.
  Currently always `cli` but future-proofs for multi-surface usage.
- **`origin`**: Add to `Message` entity. Use to annotate
  system-injected messages in journal rendering (lighter styling or
  a `[system]` prefix).

### Tier 3: Capture in parser, defer rendering

- **`sourceToolAssistantUUID`**: Add to `Message` entity. Enables
  future work on subagent call-tree reconstruction. Requires
  rethinking journal conversation rendering — out of scope for now.
- **`toolUseResult`**: Add to `Message` entity alongside existing
  `ToolResults`. CC-level errors are a distinct path from API
  tool_results; capturing both gives a complete error picture.
- **`error`** / **`apiError`**: Add to `Message` entity. Paired with
  `isApiErrorMessage` for full error context.

### Not captured

- **`attachment`**: Tool availability deltas and companion metadata.
  No journal value. Skip entirely.

## Behavior

### Happy Path

1. Parser encounters a user record with `planContent` field.
2. Field is deserialized into `Message.PlanContent string`.
3. Journal formatter checks `msg.PlanContent != ""` and renders:

   ```markdown
   <details>
   <summary>📋 Plan</summary>

   {plan content rendered as markdown}

   </details>
   ```

4. Parser encounters an assistant record with `isApiErrorMessage: true`.
5. Journal formatter collapses the message:

   ```markdown
   > ⚠ API error response (message omitted)
   ```

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `planContent` is empty string | Treat as absent — no plan section rendered |
| `planContent` contains markdown that conflicts with journal formatting | Render inside a fenced details block to isolate formatting |
| `isApiErrorMessage` on a message with useful text content | Still collapse — API errors are always retry noise, the useful response comes in the next message |
| Multiple consecutive API error messages | Collapse into one line: "⚠ N API errors (retried)" |
| `sourceToolAssistantUUID` references a UUID not in the current session | Store as-is — cross-session references are valid for subagent chains |
| `entrypoint` has an unknown value (not cli/ide/sdk-ts/sdk-py) | Store as-is — new entrypoints don't need parser changes |
| `toolUseResult` and `tool_result` content block both present for same tool call | Keep both — they represent different error paths (CC-level vs API-level) |

### Validation Rules

No validation needed — these are all optional fields. If present,
deserialize. If absent, zero value. The schema validation spec
(`jsonl-schema-validation.md`) handles field presence tracking.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `planContent` fails to render as markdown | Fall back to raw text in `<pre>` block | Automatic |
| Unknown `entrypoint` value | Store as-is, no warning | None needed |

## Interface

No new CLI commands or skills. Changes are internal to the parser and
journal formatter.

Journal entries gain:
- Plan sections (when `planContent` is present)
- Collapsed API error messages (when `isApiErrorMessage` is true)
- `entrypoint` in frontmatter (when present and not `cli`)

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/journal/parser/types.go` | Add `PlanContent`, `IsApiErrorMessage`, `Error`, `ApiError`, `SourceToolAssistantUUID`, `ToolUseResult`, `Entrypoint`, `Origin` to `claudeRawMessage` |
| `internal/entity/message.go` | Add `PlanContent string`, `IsApiError bool`, `SourceToolAssistantUUID string`, `ToolUseResult string`, `Origin string` to `Message` |
| `internal/entity/session.go` | Add `Entrypoint string` to `Session` |
| `internal/journal/parser/parse.go` | Map new raw fields to entity fields in `convertMessage` and `buildSession` |
| `internal/cli/journal/core/source/format/format.go` | Render plan sections, collapse API errors |
| `internal/entity/journal.go` | Add `entrypoint` to `JournalFrontmatter` |

### Key Functions

```go
// In parse.go — convertMessage additions
msg.PlanContent = raw.PlanContent
msg.IsApiError = raw.IsApiErrorMessage
msg.SourceToolAssistantUUID = raw.SourceToolAssistantUUID
msg.ToolUseResult = raw.ToolUseResult
msg.Origin = raw.Origin

// In format.go — plan rendering
if msg.PlanContent != "" {
    fmt.Fprintf(w, "<details>\n<summary>📋 Plan</summary>\n\n")
    fmt.Fprintf(w, "%s\n\n", msg.PlanContent)
    fmt.Fprintf(w, "</details>\n\n")
}

// In format.go — API error collapsing
if msg.IsApiError {
    fmt.Fprintf(w, "> ⚠ API error response (message omitted)\n\n")
    return // skip normal message rendering
}
```

### Helpers to Reuse

- `internal/journal/parser/parse.go` — existing `convertMessage` is
  the hook point for all new field mapping
- `internal/cli/journal/core/source/format/format.go` — existing
  message rendering loop is the hook point for plan/error rendering

## Configuration

None. All behavior is automatic based on field presence.

## Testing

- **Unit**: Parse JSONL lines with `planContent`, verify `Message.PlanContent`
  is populated. Parse lines with `isApiErrorMessage: true`, verify
  `Message.IsApiError` is set. Parse lines without these fields, verify
  zero values.
- **Integration**: Import a session containing plan mode usage, verify
  journal entry has plan section. Import a session with API errors,
  verify errors are collapsed.
- **Golden fixtures**: Add test JSONL snippets with envelope fields to
  `internal/journal/parser/testdata/`.

## Non-Goals

- Reconstructing subagent call trees from `sourceToolAssistantUUID`
  (captured but not rendered — future work)
- Parsing `attachment` records (no journal value)
- Validating envelope field values beyond type correctness
- Changing the journal format for existing entries (only new imports
  get the enriched rendering)
- Capturing `progress` record data (tool execution progress is
  transient, not archival)
