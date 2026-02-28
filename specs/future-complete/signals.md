# Signal System

## Overview

Signals are the communication channel between `ctx monitor` (writer) and
Claude Code sessions (reader via UserPromptSubmit hook).

**Design principle**: The user's prompt input is sacred. Only interrupt for
actionable, timely issues.

## Directory Structure

Session-scoped directories (not encoded in filenames):

```
.context/signals/
├── sessions/
│   ├── abc123/
│   │   ├── context-health.md
│   │   └── repetition.md
│   └── def456/
│       └── context-health.md
└── global/
    └── audit-reminder.md
```

- `sessions/{session_id}/` — signals for specific session
- `global/` — broadcast to all sessions

## Signal File Format

```markdown
---
generated_at: 2026-02-03T10:45:00Z
severity: warning
ttl: 300
auditor: context-health
code: CTX_HEALTH_85
---
Context: 87% full, repetition in last 8 turns.
→ Summarize current state in 5 bullets, then /clear and reload from DECISIONS.md
```

### Header Fields

| Field | Required | Description |
|-------|----------|-------------|
| `generated_at` | Yes | ISO 8601 timestamp |
| `severity` | Yes | `info`, `warning`, `critical` |
| `ttl` | Yes | Seconds until expiration |
| `auditor` | Yes | Which auditor generated this |
| `code` | Yes | Alert code for dedupe (e.g., `CTX_HEALTH_85`) |

No `session_id` needed — directory structure determines scope.

### Body Format

Two parts for human-centered messaging:
1. **Summary**: What's wrong (fact, not judgment)
2. **Action**: What to do (concrete, specific)

```
{Summary}
→ {Action}
```

## Atomic Write Protocol

Prevent partial reads:

```go
func WriteSignal(dir, filename string, signal Signal) error {
    finalPath := filepath.Join(dir, filename)
    tempPath := finalPath + ".tmp"

    content := formatSignal(signal)

    if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
        return err
    }

    // Atomic rename (POSIX guarantees atomicity)
    return os.Rename(tempPath, finalPath)
}
```

## Alert Fatigue Prevention

### 1. Hysteresis

Prevent threshold flapping:

```go
var Thresholds = []struct {
    AlertAt float64
    ClearAt float64
    Code    string
}{
    {70, 60, "CTX_HEALTH_70"},
    {85, 75, "CTX_HEALTH_85"},
    {95, 85, "CTX_HEALTH_95"},
}
```

Alert at 85%, don't clear until below 75%.

### 2. Cooldown Windows

Same alert type can't fire again within cooldown:

| Alert Type | Cooldown |
|------------|----------|
| context-health | 15 min |
| repetition | 10 min |
| audit-reminder | 4 hours |

```go
type AlertState struct {
    LastEmitted map[string]time.Time  // key: "{session}:{code}"
}

func (s *AlertState) ShouldEmit(session, code string, cooldown time.Duration) bool {
    key := session + ":" + code
    if last, ok := s.LastEmitted[key]; ok {
        if time.Since(last) < cooldown {
            return false
        }
    }
    s.LastEmitted[key] = time.Now()
    return true
}
```

### 3. Severity Filtering

Default: inject `warning` and above only.

| Severity | Default Behavior |
|----------|------------------|
| `info` | Written to file, not injected. Visible via `ctx status`. |
| `warning` | Injected into context |
| `critical` | Injected + optional terminal notification |

Configurable in `.ctxrc`:

```yaml
monitor:
  inject_min_severity: warning
```

### 4. Aggregation

Multiple signals → single injection block:

```
⚠️ [ctx monitor] 2 issues:
  • Context: 87% full → Summarize and /clear
  • Repetition in last 8 turns → Try different approach
```

## TTL Guidelines

| Signal Type | TTL | Rationale |
|-------------|-----|-----------|
| context-health | 300s (5 min) | Context changes rapidly |
| repetition | 300s (5 min) | May self-correct |
| audit-reminder | 3600s (1 hour) | Stable reminder |

## Hook Implementation

UserPromptSubmit hook reads, filters, aggregates, injects:

```bash
#!/bin/bash
# .claude/hooks/prompt-inject-signals.sh

HOOK_INPUT=$(cat)
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id')
NOW=$(date +%s)
MIN_SEV="${CTX_INJECT_MIN_SEVERITY:-warning}"

severity_rank() {
    case "$1" in
        info) echo 1 ;; warning) echo 2 ;; critical) echo 3 ;; *) echo 0 ;;
    esac
}

min_rank=$(severity_rank "$MIN_SEV")
declare -a MESSAGES

process_signals() {
    local dir="$1"
    for f in "$dir"/*.md; do
        [ -e "$f" ] || continue
        [[ "$f" == *.tmp ]] && continue

        # Parse frontmatter
        GENERATED=$(grep '^generated_at:' "$f" | cut -d' ' -f2-)
        SEVERITY=$(grep '^severity:' "$f" | cut -d' ' -f2-)
        TTL=$(grep '^ttl:' "$f" | cut -d' ' -f2-)

        # Check TTL
        GEN_TS=$(date -d "$GENERATED" +%s 2>/dev/null || echo 0)
        if [ "$NOW" -gt $((GEN_TS + TTL)) ]; then
            rm -f "$f"
            continue
        fi

        # Check severity
        if [ $(severity_rank "$SEVERITY") -lt "$min_rank" ]; then
            continue
        fi

        # Extract body (after second ---)
        BODY=$(awk 'BEGIN{p=0} /^---$/{p++; next} p==2{print}' "$f")
        MESSAGES+=("$BODY")

        rm -f "$f"
    done
}

# Process session-specific and global signals
process_signals ".context/signals/sessions/$SESSION_ID"
process_signals ".context/signals/global"

# Aggregate output
if [ ${#MESSAGES[@]} -gt 0 ]; then
    if [ ${#MESSAGES[@]} -eq 1 ]; then
        echo "⚠️ [ctx monitor] ${MESSAGES[0]}"
    else
        echo "⚠️ [ctx monitor] ${#MESSAGES[@]} issues:"
        for msg in "${MESSAGES[@]}"; do
            echo "  • $msg"
        done
    fi
fi
```

## Monitor State Persistence

```json
// .context/monitor-state.json
{
  "lastEmitted": {
    "abc123:CTX_HEALTH_85": "2026-02-03T10:30:00Z",
    "abc123:CTX_REPETITION": "2026-02-03T09:15:00Z"
  },
  "activeAlerts": {
    "abc123:CTX_HEALTH_85": true
  },
  "lastHealth": {
    "abc123": {
      "contextPercent": 87.5,
      "repetitionScore": 0.3
    }
  }
}
```

## Cleanup

Signals removed:
1. After injection (hook deletes)
2. On TTL expiration (hook deletes stale)
3. On session end (SessionEnd hook deletes directory)
4. Orphaned signals (monitor periodic cleanup)
