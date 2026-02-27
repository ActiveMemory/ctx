# Context Load Gate v2: Auto-Injection

Supersedes: `specs/context-load-gate.md`

## Problem

The current context-load-gate hook tells the agent to read context files.
The agent evaluates that instruction and can rationalize skipping it —
"don't apply judgment to this rule" is itself evaluated by judgment.

No amount of imperative framing ("STOP", "MANDATORY", "Do not assess
relevance") solves this, because every soft instruction passes through
the same attention/evaluation pipeline. The compliance ceiling for
advisory hooks is ~85% at session start, ~75% mid-session.

See: `docs/blog/2026-02-25-the-homework-problem.md`, "The Escape Hatch
Problem" and "The Residual Risk" sections.

## Approach

Eliminate the compliance step. Instead of telling the agent to read files,
the hook reads the files itself and injects the content directly into the
hook's `additionalContext` response. The agent never chooses whether to
read — the content is already in its context window.

This moves enforcement from the reasoning layer (soft instruction, subject
to judgment) to the infrastructure layer (content injection, not subject
to evaluation).

### Injection Strategy

Files are injected in `config.FileReadOrder` order. Not all files get
the same treatment:

| File | Strategy | Rationale |
|------|----------|-----------|
| CONSTITUTION | Verbatim | Hard rules. Every line is a guardrail. |
| CONVENTIONS | Verbatim | Code patterns. Prevents style violations. |
| ARCHITECTURE | Verbatim | System map. Prevents wrong-place code. |
| AGENT_PLAYBOOK | Verbatim | Behavioral guidance. Core operating manual. |
| DECISIONS | Index table only | Titles for grep. Full entry on demand. |
| LEARNINGS | Index table only | Titles for grep. Full entry on demand. |
| TASKS | One-liner mention | Read when discussing priorities. |
| GLOSSARY | Skip | Corpus + surrounding context covers it. |

Note: TASKS is placed second in `FileReadOrder` (position 2). The
injection reorders slightly — verbatim files first, then indexes, then
the mention — but the message header states the canonical read order
so the agent knows where to find full files if needed.

### Why This Works

1. **Zero compliance required for core files**: CONSTITUTION, CONVENTIONS,
   ARCHITECTURE, AGENT_PLAYBOOK are in the context window as a fait
   accompli. The agent can't avoid having them.

2. **Sunk cost leverages remaining reads**: With ~7k tokens of context
   already loaded, the marginal cost of reading a specific DECISIONS or
   LEARNINGS entry is trivial. The rationalization path inverts — "I
   already have 80% of the context, why skip the last 20%?"

3. **Index tables as lookup keys**: The DECISIONS and LEARNINGS index
   tables give the agent titles to grep for. Full entry bodies are
   demand-loaded when relevant to the task.

4. **Lost-in-the-middle mitigation**: CONSTITUTION is first in the
   injection (primacy position). At injection time (first tool use), the
   context window is fresh. As the session grows, the injection drifts
   toward the middle, but CLAUDE.md (true primacy) and the
   check-context-size nudge (wrap-up signal) bound the degradation.

### What Changes from v1

| Aspect | v1 (current) | v2 (this spec) |
|--------|-------------|----------------|
| Hook output | File paths + "read them" | File contents + indexes |
| Compliance model | Agent reads files (judgable) | Content injected (fait accompli) |
| Token cost at hook | ~200 tokens (paths) | ~7,700 tokens (content) |
| Agent effort | 7 Read tool calls | 0 tool calls for core files |
| Failure mode | Silent skip (no reads) | Content ignored (but still present) |
| VERBATIM relay | Skip message if agent skips | Not needed — nothing to skip |
| Webhook payload | Full directive text | Metadata only (see Security) |

## Behavior

### Happy Path

1. User sends first message
2. Agent decides to use a tool (Read, Bash, Grep, etc.)
3. PreToolUse fires → `ctx system context-load-gate` runs
4. Hook checks for session marker file
5. No marker → hook reads context files from disk:
   a. CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK: full content
   b. DECISIONS, LEARNINGS: extract `INDEX:START` to `INDEX:END` block
   c. TASKS: one-liner mention
   d. GLOSSARY: skip
6. Hook assembles content into `additionalContext` JSON
7. Hook estimates total tokens (heuristic: len/4)
8. Create marker file
9. Send webhook notification with metadata (not content)
10. Agent receives response — context is in the window

### Message Format

```
PROJECT CONTEXT (auto-loaded by system hook — already in your context window)
================================================================================

--- CONSTITUTION.md ---
[full file content]

--- CONVENTIONS.md ---
[full file content]

--- ARCHITECTURE.md ---
[full file content]

--- AGENT_PLAYBOOK.md ---
[full file content]

--- DECISIONS.md (index — read full entries by date when relevant) ---
[INDEX:START to INDEX:END content]

--- LEARNINGS.md (index — read full entries by date when relevant) ---
[INDEX:START to INDEX:END content]

================================================================================
Context: N files loaded (~XXXX tokens). Order follows config.FileReadOrder.

TASKS.md contains the project's prioritized work items. Read it when
discussing priorities, picking up work, or when the user asks about tasks.

For full decision or learning details, read the entry in DECISIONS.md or
LEARNINGS.md by timestamp.
```

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `.context/` not initialized | Silent — no gate on non-ctx projects |
| Session marker already exists | Silent — gate already fired |
| A context file is missing | Skip that file, include others |
| A context file is empty | Include with "(empty)" note |
| DECISIONS/LEARNINGS has no index markers | Include full file as fallback |
| Index extraction finds no content | Include "(no entries)" note |
| Total injection exceeds 15k tokens | Log warning in webhook, inject anyway |

### Validation Rules

- Marker file must include session ID to avoid cross-session leaks
- Marker lives in secure temp dir (volatile)
- Hook must complete within 2 seconds (same as all hooks)
- Files are read from `rc.ContextDir()` resolved path
- File read uses `os.ReadFile` — no validation boundary needed (files
  are in `.context/`, already within project)

### Error Handling

| Error condition | Behavior | Recovery |
|-----------------|----------|----------|
| Cannot create marker file | Emit injection every time | Noisy but functional |
| Cannot read `.context/` | Silent (not initialized) | None needed |
| Individual file read fails | Skip file, include others | Partial injection |
| stdin timeout | Silent (graceful degradation) | None needed |
| Index markers not found | Include full file content | Fallback, not error |

## Interface

### CLI

```
ctx system context-load-gate
```

Hidden subcommand. No user-facing flags. Same command name as v1.

### Hook Registration (hooks.json)

No change from v1. Same matcher, same position:

```json
{
  "matcher": ".*",
  "hooks": [{"type": "command", "command": "ctx system context-load-gate"}]
}
```

## Implementation

### Files to Modify

| File | Change |
|------|--------|
| `internal/cli/system/context_load_gate.go` | Rewrite — read files, inject content |
| `internal/cli/system/context_load_gate_test.go` | Update — test content injection |

No new files. No hook registration changes. No config changes.

### Key Implementation

```go
func runContextLoadGate(cmd *cobra.Command, stdin *os.File) error {
    if !isInitialized() {
        return nil
    }

    input := readInput(stdin)
    if input.SessionID == "" {
        return nil
    }

    tmpDir := secureTempDir()
    marker := filepath.Join(tmpDir, "ctx-loaded-"+input.SessionID)
    if _, err := os.Stat(marker); err == nil {
        return nil
    }

    touchFile(marker)

    dir := rc.ContextDir()
    var content strings.Builder
    var totalTokens int
    var filesLoaded int

    content.WriteString(
        "PROJECT CONTEXT (auto-loaded by system hook " +
        "— already in your context window)\n" +
        strings.Repeat("=", 80) + "\n\n")

    // Verbatim files: CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK
    // Index-only files: DECISIONS, LEARNINGS
    // Mention-only: TASKS
    // Skip: GLOSSARY
    for _, f := range config.FileReadOrder {
        if f == config.FileGlossary {
            continue
        }

        path := filepath.Join(dir, f)
        data, err := os.ReadFile(path)
        if err != nil {
            continue // file missing — skip gracefully
        }

        switch f {
        case config.FileTask:
            // One-liner mention, don't inject content
            continue

        case config.FileDecision, config.FileLearning:
            // Extract index table only
            idx := extractIndex(string(data))
            if idx == "" {
                idx = "(no index entries)"
            }
            content.WriteString(
                fmt.Sprintf("--- %s (index — read full entries "+
                    "by date when relevant) ---\n%s\n\n", f, idx))
            totalTokens += context.EstimateTokensString(idx)
            filesLoaded++

        default:
            // Verbatim injection
            content.WriteString(
                fmt.Sprintf("--- %s ---\n%s\n\n", f, string(data)))
            totalTokens += context.EstimateTokens(data)
            filesLoaded++
        }
    }

    // Footer
    content.WriteString(strings.Repeat("=", 80) + "\n")
    content.WriteString(fmt.Sprintf(
        "Context: %d files loaded (~%d tokens). "+
            "Order follows config.FileReadOrder.\n\n"+
            "TASKS.md contains the project's prioritized work items. "+
            "Read it when discussing priorities, picking up work, "+
            "or when the user asks about tasks.\n\n"+
            "For full decision or learning details, read the entry "+
            "in DECISIONS.md or LEARNINGS.md by timestamp.\n",
        filesLoaded, totalTokens))

    printHookContext(cmd, "PreToolUse", content.String())

    // Webhook: metadata only — no file content in payload
    webhookMsg := fmt.Sprintf(
        "context-load-gate: injected %d files (~%d tokens)",
        filesLoaded, totalTokens)
    _ = notify.Send("relay", webhookMsg, input.SessionID, "")

    return nil
}

// extractIndex returns the content between INDEX:START and INDEX:END
// markers, or empty string if markers are not found.
func extractIndex(content string) string {
    start := strings.Index(content, config.IndexStart)
    end := strings.Index(content, config.IndexEnd)
    if start < 0 || end < 0 || end <= start {
        return ""
    }
    startPos := start + len(config.IndexStart)
    return strings.TrimSpace(content[startPos:end])
}
```

### Helpers to Reuse

- `isInitialized()` — checks `.context/` exists
- `readInput()` — parses stdin JSON with session_id
- `printHookContext()` — emits JSON `additionalContext`
- `notify.Send()` — webhook notification
- `rc.ContextDir()` — resolves context directory
- `context.EstimateTokens()` / `context.EstimateTokensString()` —
  token estimation
- `config.FileReadOrder` — canonical file ordering
- `config.IndexStart` / `config.IndexEnd` — index markers
- `secureTempDir()` / `touchFile()` — marker file management

## Security

### Webhook Payload

The webhook payload MUST NOT contain file contents. Reasons:

1. **Payload size**: Injected content is ~7-8k tokens. Webhook payloads
   should be metadata, not bulk transfer.
2. **Sensitive content**: DECISIONS and LEARNINGS may contain
   project-specific information that should not leave the local machine
   via webhook.
3. **Logging exposure**: Webhook payloads are often logged by receiving
   services, intermediaries, and debug infrastructure.

The webhook sends metadata only:
```
context-load-gate: injected 6 files (~7742 tokens)
```

The previous v1 behavior passed the full directive text (including file
paths) to `notify.Send()`. This spec changes the fourth argument to
empty string, sending only the summary message.

## Configuration

None. The hook is always active when `.context/` is initialized.
No `.ctxrc` keys needed. Same as v1.

## Testing

### Unit Tests

- **No `.context/`**: silent (no output)
- **Empty session ID**: silent
- **First tool use**: injects content with all expected files
- **Second tool use**: silent (marker exists)
- **Different sessions**: each gets own injection
- **Verbatim files present**: CONSTITUTION, CONVENTIONS, ARCHITECTURE,
  AGENT_PLAYBOOK content appears in output
- **Index-only files**: DECISIONS and LEARNINGS show index table, not
  full content
- **TASKS not in content**: mentioned in footer only
- **GLOSSARY not in output**: excluded entirely
- **Missing file**: other files still injected
- **No index markers**: fallback to "(no index entries)"
- **Token count in footer**: matches heuristic estimate
- **Webhook payload**: contains metadata, NOT file content

### Manual Verification

Same protocol as v1 spec:

1. `make build && sudo make install`
2. Start a new Claude Code session in a ctx-initialized project
3. Ask a task that requires tool use but does NOT mention context:
   - Good: "Add a --verbose flag to ctx status"
   - Bad: "What are the pending tasks?" (probes TASKS.md directly)
4. Observe: the agent's FIRST response should reference project
   conventions, architecture, or constitution without having made
   any Read tool calls for `.context/` files
5. Ask a question about a specific decision — agent should read the
   full entry from DECISIONS.md (on-demand, from the index title)

## Non-Goals

- **Replacing `ctx agent`**: The context packet remains a separate,
  budget-optimized view. The gate provides raw files; `ctx agent`
  provides scored, prioritized summaries.
- **Token budgeting the injection**: All core files are injected
  verbatim. The value of having them always present outweighs the
  token cost. If a project's context files are too large, the correct
  fix is `/ctx-consolidate`, not injection budgeting.
- **VERBATIM relay / compliance canary**: No longer needed. There is
  no compliance step to monitor — the content is injected, not
  requested.
- **Modifying `config.FileReadOrder`**: The injection uses the existing
  order. The per-file strategy (verbatim/index/mention/skip) is
  specific to this hook, not a change to the canonical ordering.
