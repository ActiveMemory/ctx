# Context Load Gate

## Problem

When a Claude Code session starts, the agent receives CLAUDE.md instructions
telling it to load context files (AGENT_PLAYBOOK, TASKS.md, etc.). But the
agent routinely skips this — it jumps straight to answering the user's
question instead.

### Observed Agent Behavior (First-Person Analysis)

What the agent actually sees at session start:

| Source | What arrives | Impact |
|--------|-------------|--------|
| `claudeMd` system-reminder | Full CLAUDE.md contents | Agent *knows* what it should do |
| Skill list system-reminder | Available `/skills` | Informational |
| `UserPromptSubmit` hook outputs | `ctx system` help banner (4x) | **Zero** — just noise |
| Git status snapshot | Branch, dirty files, recent commits | Mildly useful |

What the agent does NOT have:
- AGENT_PLAYBOOK.md
- TASKS.md, DECISIONS.md, LEARNINGS.md, CONVENTIONS.md
- CONSTITUTION.md
- Any output from `ctx agent --budget 4000`

**Why the agent skips loading**: The user's message is the primary directive.
CLAUDE.md says "load context on session start" but the user's actual question
says "do X". The agent optimizes for responsiveness — answering the question
wins over a multi-step bootstrap ceremony that produces no visible output.
The instructions compete with the task, and the task always wins.

### Experimental Evidence: Delegation Kills Urgency

A controlled experiment (2026-02-25, Opus 4.6) tested a UserPromptSubmit
hook that said: "STOP. Before answering the user's question, run
`ctx system bootstrap` and follow its instructions."

The user asked: "can you add --verbose to the info command?"

**Results:**

| Step | Expected | Actual |
|------|----------|--------|
| Notice hook nudge | Yes | **Yes** — improvement over previous agent that ignored it |
| Run `ctx system bootstrap` | Yes | **Yes** — the "STOP" language worked |
| Read AGENT_PLAYBOOK.md (bootstrap told it to) | Yes | **No** — skipped |
| Run `ctx agent --budget 4000` (bootstrap told it to) | Yes | **No** — skipped |
| Wait for bootstrap before starting task | Yes | **No** — parallelized task exploration with bootstrap |

**Root cause — the "decaying urgency chain":**

```
Hook says "STOP"          → agent complies (high authority)
Hook says "run bootstrap" → agent runs it (direct instruction)
Bootstrap output says "Next steps: read AGENT_PLAYBOOK" → agent skips (feels like a suggestion)
Bootstrap output says "Run ctx agent" → agent skips (even weaker)
```

Each link in the delegation chain loses enforcement power. The hook's
authority doesn't transfer to the commands it delegates to. The agent
treats the hook itself as the obligation and the rest as optional.

**Key insight for this spec**: The hook message must contain the
**complete, actionable instruction** — not delegate to a command
whose output contains further instructions. "Read these files: [list]"
is one link. "Run bootstrap and follow its output" is three links,
and the agent drops the chain after the first.

Full analysis: `hook-nudge-analysis.md` (project root).

### Why Existing Hooks Don't Solve This

1. **`ctx agent --budget 4000` PreToolUse hook exists** (in embedded
   `hooks.json`, matcher `.*`) — but it has `2>/dev/null || true`, so
   failures are silent. More importantly, it has a **10-minute cooldown**,
   so after the first tool use it goes quiet. And the agent's first tool
   call may already be task-focused, not context-focused.

2. **`UserPromptSubmit` hooks are crowded** — 9+ checks fire on every
   prompt (ceremonies, persistence, journal, version, resources, knowledge,
   map staleness, reminders, context size). When everything fires at once,
   the agent treats them as background noise. Adding more text here
   would worsen the signal-to-noise ratio.

3. **CLAUDE.md instructions are aspirational, not enforced** — "Run
   `ctx system bootstrap`" is a request, not a gate. The agent weighs
   it against the user's actual question and decides bootstrapping can
   wait.

## Approach

A new `ctx system context-load-gate` command that fires as a **PreToolUse
hook** on the first tool invocation of a session. It checks whether the
agent has loaded context, and if not, emits a short, sharp directive:

```
STOP. Read your context files before proceeding: .context/CONSTITUTION.md,
.context/TASKS.md, .context/CONVENTIONS.md, .context/DECISIONS.md,
.context/LEARNINGS.md, .context/AGENT_PLAYBOOK.md
```

### Why This Works

1. **Timing**: The message arrives at the moment of action — when the
   agent is about to use its first tool. This is the highest-salience
   moment in the session. The agent's attention is focused, not scattered.

2. **Brevity**: A single clear line beats a wall of text. The agent won't
   skim or ACK-and-ignore a one-liner the way it does with multi-paragraph
   relay messages.

3. **Non-blocking**: The hook uses `additionalContext` (JSON directive),
   not `decision: block`. The agent receives the instruction and can
   choose to batch the reads into its next response. No deadlock risk.

4. **One-shot**: Fires only on the first tool use per session (tracked
   via a session-scoped marker file). Subsequent tool calls proceed
   without the gate.

### Why Not Hard-Block?

A blocking gate (`decision: block`) creates a chicken-and-egg problem:
the agent needs to Read files to satisfy the gate, but Read is a tool
that would also be blocked. Whitelisting specific tools adds complexity.
The non-blocking directive is sufficient because:
- The context window is fresh (low noise)
- The agent is not yet deep in a task (low switching cost)
- A clear, authoritative instruction at this moment has very high
  compliance rates

## Behavior

### Happy Path

1. User sends first message
2. Agent decides to use a tool (Read, Bash, Grep, etc.)
3. PreToolUse fires → `ctx system context-load-gate` runs
4. Hook checks for session marker file (e.g., `/tmp/ctx-loaded-{session_id}`)
5. No marker → emit JSON directive with the file list
6. Create marker file
7. Agent receives directive, reads the listed files
8. Subsequent tool calls: marker exists → hook is silent

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `.context/` not initialized | Silent — no gate on non-ctx projects |
| Session marker already exists | Silent — gate already fired |
| Agent ignores directive | No enforcement — same as status quo, no worse |
| `ctx agent` hook also fires | Both fire; agent gets context packet AND file list. Redundant but not harmful |
| Very short session (one tool call) | Gate fires once, agent may or may not read files. Acceptable. |

### Validation Rules

- Marker file must include session ID to avoid cross-session leaks
- Marker lives in `/tmp/` (volatile) — never persists across reboots
- Hook must complete within 2 seconds (same stdin timeout as all hooks)

### Error Handling

| Error condition | Behavior | Recovery |
|-----------------|----------|----------|
| Cannot create marker file | Emit directive every time | Noisy but functional |
| Cannot read `.context/` | Silent (not initialized) | None needed |
| stdin timeout | Silent (graceful degradation) | None needed |

## Interface

### CLI

```
ctx system context-load-gate
```

Hidden subcommand. No user-facing flags.

### Hook Registration (hooks.json)

```json
{
  "matcher": ".*",
  "hooks": [{"type": "command", "command": "ctx system context-load-gate"}]
}
```

Position: **first** in the PreToolUse array, before `block-non-path-ctx`
and `qa-reminder`. The gate should fire before any blocking hooks, so the
agent receives the "read your files" directive even if a subsequent hook
blocks the tool call.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/cli/system/context_load_gate.go` | New — hook implementation |
| `internal/cli/system/context_load_gate_test.go` | New — unit tests |
| `internal/cli/system/system.go` | Register `contextLoadGateCmd()` |
| `internal/assets/claude/hooks/hooks.json` | Add `.*` matcher entry (first position) |

### Key Implementation

```go
func contextLoadGateCmd() *cobra.Command {
    return &cobra.Command{
        Use:    "context-load-gate",
        Short:  "Emit context-load directive on first tool use",
        Hidden: true,
        RunE: func(cmd *cobra.Command, _ []string) error {
            if !isInitialized() {
                return nil
            }

            input := readInput(os.Stdin)
            sid := input.SessionID
            if sid == "" {
                return nil
            }

            marker := filepath.Join(os.TempDir(),
                fmt.Sprintf("ctx-loaded-%s", sid))
            if _, err := os.Stat(marker); err == nil {
                return nil // already fired this session
            }

            // Create marker before emitting — ensures one-shot
            _ = os.WriteFile(marker, []byte(time.Now().Format(time.RFC3339)), 0600)

            dir := rc.ContextDir()
            msg := fmt.Sprintf(
                "STOP. Read your context files before proceeding: "+
                    "%s/CONSTITUTION.md, %s/TASKS.md, %s/CONVENTIONS.md, "+
                    "%s/DECISIONS.md, %s/LEARNINGS.md, %s/AGENT_PLAYBOOK.md",
                dir, dir, dir, dir, dir, dir,
            )
            printHookContext(cmd, "PreToolUse", msg)
            _ = notify.Send("relay",
                "context-load-gate: directed agent to read context files", "")
            return nil
        },
    }
}
```

### Helpers to Reuse

- `isInitialized()` — checks `.context/` exists with required files
- `readInput()` — parses stdin JSON with session_id
- `printHookContext()` — emits JSON `additionalContext` directive
- `notify.Send()` — relay notification
- `rc.ContextDir()` — resolves context directory path

## Configuration

None. The hook is always active when `.context/` is initialized.
No `.ctxrc` keys needed.

## Testing

### Unit Tests

- Session marker absent → directive emitted, marker created
- Session marker present → silent (no output)
- `.context/` not initialized → silent
- Empty session ID → silent (graceful degradation)

### Integration / Manual Verification

1. `make build && sudo make install`
2. Start a new Claude Code session in a ctx-initialized project
3. Ask a task that requires tool use but does NOT mention context files:
   - Good: "Add a --verbose flag to ctx status"
   - Good: "Fix the lint warning in cmd/ctx/main.go"
   - Bad: "What are the pending tasks?" (probes agent to read TASKS.md
     regardless of gate — tests the question, not the mechanism)
4. Observe the **ordering of tool calls**:
   - Gate worked: first tool calls are Read for CONSTITUTION.md, TASKS.md,
     AGENT_PLAYBOOK.md, etc. — before any task-related grep/glob
   - Gate failed: first tool call is task-focused (grep, glob, read of
     source files) — agent skipped context loading
5. Ask a follow-up question → verify the gate does NOT fire again
6. Verify marker file: `ls /tmp/ctx-loaded-*` shows one file per session

## Non-Goals

- **Hard blocking**: Not blocking tool calls, only injecting a directive.
- **ACK verification**: Not hard-verifying file reads. The gate uses an
  unconditional checkpoint block ("Context Loaded" with Read/Skipped fields)
  that the agent must always output. This is a fill-in-the-blank template,
  not a conditional — models are more likely to complete templates than
  evaluate conditionals.
- **Replacing `ctx agent`**: The existing `ctx agent --budget 4000`
  PreToolUse hook serves a different purpose (context packet summary).
  The gate complements it with explicit file paths.
- **Modifying UserPromptSubmit hooks**: The existing hooks are already
  crowded. This deliberately uses PreToolUse instead.

## Open Questions

1. **Message wording**: Is "STOP." too aggressive? Alternatives:
   "Before proceeding, read your context files:" — but stronger wording
   has higher compliance. The qa-reminder uses "HARD GATE" successfully.
   Recommendation: keep "STOP." — it mirrors the authority pattern that
   `qa-reminder` established.

2. **Interaction with `ctx agent --budget 4000` hook**: Both fire on
   `.*` matcher. The agent receives both the context packet (summary)
   and the file list (explicit paths). This is redundant but the packet
   has a 10-minute cooldown so it won't always fire. The gate should
   always fire on first tool use regardless.

3. **Marker cleanup**: Markers in `/tmp/` accumulate. The existing
   `cleanup-tmp` SessionEnd hook handles `/tmp/ctx-*` files older than
   15 days. The `ctx-loaded-{session_id}` pattern fits this cleanup.
   No additional cleanup needed.
