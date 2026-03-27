# Hook accountability: checkpoint gating and commit enforcement

## Problem

Two related issues with the hook system:

### 1. Counter-based checkpoints ignore context window usage

The `check-context-size` hook fires checkpoint nudges based on prompt
count alone (every 5 prompts after 15, every 3 after 30). On a 1M
context window, this produces noise at 5-8% usage — hundreds of
"consider wrapping up" messages when there's nothing to wrap up.

The irony: the checkpoint message already includes the percentage
("Context window: ~80k tokens (~8% of 1000k)"), proving the system
*knows* it's premature but nudges anyway.

The same problem affects `check-persistence`: it nudges "no context
files updated in 15+ prompts" even when the session is a simple
recurring loop with nothing to persist.

### 2. Context window tier ordering is wrong

`EffectiveContextWindow` prioritizes `.ctxrc` config over detected
reality:

1. `.ctxrc` `context_window` (manual, possibly stale)
2. Claude Code settings.json `[1m]` detection
3. JSONL model ID prefix
4. Default (200k)

Ground truth (JSONL model ID, settings.json) should outrank a config
file that may not have been updated.

### 3. No spec enforcement at commit time

Agents skip specs because nothing enforces them. The AGENT_PLAYBOOK
says "spec first, then task" but it's advisory. Agents take the
shortest path: create a one-line task, start coding, commit without
a spec. The spec is where thinking happens — without it, the agent
jumps from task to code without design.

## Approach

### A. Gate counter-based checkpoints behind minimum percentage

Add a minimum context window percentage threshold below which
counter-based checkpoints are suppressed. The `windowTrigger` (>80%)
remains independent and unchanged.

**Current logic** (`check_context_size/run.go:112-116`):
```go
if count > 30 {
    counterTriggered = count%3 == 0
} else if count > 15 {
    counterTriggered = count%5 == 0
}
```

**New logic**:
```go
if pct >= stats.ContextCheckpointMinPct {
    if count > 30 {
        counterTriggered = count%3 == 0
    } else if count > 15 {
        counterTriggered = count%5 == 0
    }
}
```

New constant: `ContextCheckpointMinPct = 20` in `config/stats/context.go`.

Below 20% usage, only the window warning (>80%) can fire. This
eliminates noise on 1M windows while preserving urgency on 200k
windows where 20% is reached faster.

**Same gating for `check-persistence`**: add `pct` awareness to
`PersistenceNudgeNeeded`. The persistence hook currently has no
access to token data — it will need to read session stats (already
written by `check-context-size` on every prompt) to get the current
percentage.

### B. Reorder context window tiers

New priority order (ground truth wins):

1. **JSONL model ID** — actual model running the session
2. **Claude Code settings.json** — configured model selection
3. **`.ctxrc` context_window** — manual override / escape hatch
4. **Default** (200k)

Change in `session_token.go:EffectiveContextWindow`:
```go
func EffectiveContextWindow(model string) int {
    // Tier 1: model-based detection (ground truth from session).
    if w := ModelContextWindow(model); w > 0 {
        return w
    }
    // Tier 2: auto-detect from Claude Code settings.
    if ClaudeSettingsHas1M() {
        return ContextWindow1M
    }
    // Tier 3: explicit .ctxrc override (fallback for non-Claude tools).
    if w := rc.RC().ContextWindow; w > 0 && w != rc.DefaultContextWindow {
        return w
    }
    // Tier 4: default.
    return rc.ContextWindow()
}
```

### C. Spec enforcement at commit time

**Philosophy**: trust first, consequence later. No friction during
coding. Accountability at the single choke point: commit.

#### C1. CONSTITUTION addition

Add to CONSTITUTION.md:
```
- [ ] Every commit references a spec (`Spec: specs/<name>.md` trailer)
```

#### C2. /ctx-commit skill update

Add a mandatory `Spec:` trailer requirement to the commit skill:

1. Before committing, check: does the commit message contain a `Spec:`
   trailer referencing an existing file in `specs/`?
2. If not, stop and tell the agent: "Create a spec first, or ask the
   human for a retroactive spec. This is a CONSTITUTION requirement.
   No exceptions — even one-liner fixes need a spec for traceability."
3. The wording must be absolute: no "non-trivial" qualifier, no
   "if applicable." The cost of a minimal spec is near-zero.

Commit message format:
```
add percentage gating to checkpoint hooks

Gate counter-based checkpoint nudges behind a 20% minimum context
window usage threshold to eliminate noise on large context windows.

Spec: specs/hook-accountability.md
Signed-off-by: Jose Alekhinne <jose@ctx.ist>
```

#### C3. Post-commit hook: bypass detection

Add telltale checks to `post-commit` for commits that bypassed
`/ctx-commit`. A commit scores violation points:

| Signal                                          | Points |
|-------------------------------------------------|--------|
| Missing `Spec:` trailer                         | 3      |
| Missing `Signed-off-by:` trailer                | 1      |
| No task reference in message body               | 1      |
| Source files changed but no TASKS.md in diff     | 1      |
| Single-line message (no structured body)         | 1      |

Scoring:
- 0-1: clean, no action
- 2-3: nudge to human ("commit looks informal")
- 4+: relay warning to human ("agent bypassed /ctx-commit")

The relay goes to the human verbatim, not to the agent. The human
is the gate; the hook is the sensor.

#### C4. Human relay at commit

After every commit (whether via `/ctx-commit` or detected by
post-commit hook), relay a structured summary to the human:

```
┌─ Commit Summary ─────────────────────────
│ Spec: specs/hook-accountability.md
│ Tasks closed: HA.3, HA.4
│ Files changed: 4
│ Violations: 0
└──────────────────────────────────────────
```

If violations > 0, the relay includes what's missing.

## Storage

- New constant `ContextCheckpointMinPct` in `config/stats/context.go`
- Tier reordering in `session_token.go` (no new storage)
- `Spec:` trailer in commit messages (git metadata, no new files)
- Post-commit violation scoring in `post-commit` hook logic

## CLI surface

No new commands. Changes are internal to existing hooks and the
`/ctx-commit` skill.

## Error cases

- **No JSONL available** (tier 1 fails): falls through to tier 2-4
  as before
- **Spec file referenced but doesn't exist**: `/ctx-commit` rejects
  the commit with a clear message
- **Agent bypasses /ctx-commit entirely**: post-commit hook catches it
  and relays to human
- **Percentage data unavailable** (tokens=0): counter-based checkpoints
  fire as before (no regression for sessions without token data)

## Non-goals

- Blocking edits during coding (no PreToolUse gates)
- Validating spec content quality (that's the human's job)
- Correlating tasks to specific files being edited
- Changing the 80% window warning threshold

## Testing

- Unit test: `EffectiveContextWindow` with all 4 tiers
- Unit test: counter gating suppresses below `ContextCheckpointMinPct`
- Unit test: counter gating allows at/above threshold
- Unit test: persistence nudge suppressed below threshold
- Unit test: post-commit violation scoring
- Integration: 1M window session produces no checkpoint noise below 20%
