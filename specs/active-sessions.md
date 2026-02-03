# Active Session Tracking

## Overview

Track which Claude Code sessions are currently active so `ctx monitor` knows
what to watch. Handle crashes gracefully via heartbeat-based lifecycle.

## Directory Structure

```
.context/
├── active-sessions/
│   ├── abc123--9f2c1a7e.tomb    # mtime = last heartbeat
│   └── def456--3b8c2d1f.tomb
└── session-index/
    ├── 9f2c1a7e.json
    └── 3b8c2d1f.json
```

## Tombstone Design

### Why Tombstones + Index?

Encoding transcript paths in filenames has problems:
- Ambiguity (`/` vs `_` collisions)
- Length limits (255 bytes)
- Unicode edge cases
- Path disclosure

Solution: Hash the path, store metadata separately.

### Filename Format

```
{session_id}--{path_hash}.tomb
```

- `session_id`: UUID from Claude Code
- `path_hash`: First 8 chars of SHA256(transcript_path)
- `.tomb`: Extension for clarity
- Content: Empty (all info in filename + mtime)

### Index File

```json
{
  "sessionId": "abc123-def4-5678-...",
  "transcript": "~/.claude/projects/-home-jose-ctx/abc123.jsonl",
  "projectDir": "/home/jose/WORKSPACE/ctx",
  "startTime": "2026-02-03T10:30:00Z",
  "model": "claude-opus-4-5-20251101"
}
```

## Session Lifecycle

```
┌─────────┐  SessionStart   ┌─────────┐  prompt submit  ┌─────────┐
│  None   │────────────────▶│  Active │◀───────────────▶│  Active │
└─────────┘  (create tomb)  └────┬────┘  (touch mtime)  └─────────┘
                                 │
           ┌─────────────────────┼─────────────────────┐
           │                     │                     │
           ▼                     ▼                     ▼
     SessionEnd            mtime stale           mtime stale
     (graceful)           transcript ok        transcript gone
           │                     │                     │
           ▼                     ▼                     ▼
     ┌─────────┐           ┌─────────┐           ┌─────────┐
     │ (gone)  │           │ Suspect │           │  Dead   │
     └─────────┘           └────┬────┘           └────┬────┘
                                │     timeout         │
                                └──────────┬──────────┘
                                           ▼
                                     ┌──────────┐
                                     │ Cleanup  │
                                     └──────────┘
```

### States

| State   | Tombstone mtime | Transcript           | Action         |
|---------|-----------------|----------------------|----------------|
| Active  | fresh (< 5 min) | —                    | Monitor health |
| Suspect | stale           | exists, mtime recent | Wait, log      |
| Dead    | stale           | missing or stale     | Cleanup        |

**Suspect** handles long Claude operations where user hasn't typed but work continues.

### Thresholds

```go
const (
    HeartbeatThreshold = 5 * time.Minute
    SuspectTimeout     = 30 * time.Minute
    TranscriptStale    = 30 * time.Minute
)
```

## Hooks

### SessionStart

```bash
#!/bin/bash
# .claude/hooks/session-start-track.sh

HOOK_INPUT=$(cat)
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id')
TRANSCRIPT=$(echo "$HOOK_INPUT" | jq -r '.transcript_path')
MODEL=$(echo "$HOOK_INPUT" | jq -r '.model // "unknown"')

PATH_HASH=$(echo -n "$TRANSCRIPT" | sha256sum | cut -c1-8)

mkdir -p .context/active-sessions \
         .context/session-index \
         .context/signals/sessions/"$SESSION_ID" \
         .context/signals/global

touch ".context/active-sessions/${SESSION_ID}--${PATH_HASH}.tomb"

cat > ".context/session-index/${PATH_HASH}.json" << EOF
{
  "sessionId": "$SESSION_ID",
  "transcript": "$TRANSCRIPT",
  "projectDir": "$(pwd)",
  "startTime": "$(date -Iseconds)",
  "model": "$MODEL"
}
EOF
```

### UserPromptSubmit (Heartbeat)

```bash
#!/bin/bash
# Combined with signal injection hook

HOOK_INPUT=$(cat)
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id')

# Heartbeat: update tombstone mtime
touch .context/active-sessions/${SESSION_ID}--*.tomb 2>/dev/null || true

# Signal injection follows (see signals.md)
```

### SessionEnd

```bash
#!/bin/bash
# .claude/hooks/session-end-cleanup.sh

HOOK_INPUT=$(cat)
SESSION_ID=$(echo "$HOOK_INPUT" | jq -r '.session_id')
TRANSCRIPT=$(echo "$HOOK_INPUT" | jq -r '.transcript_path')
PATH_HASH=$(echo -n "$TRANSCRIPT" | sha256sum | cut -c1-8)

rm -f ".context/active-sessions/${SESSION_ID}--${PATH_HASH}.tomb"
rm -f ".context/session-index/${PATH_HASH}.json"
rm -rf ".context/signals/sessions/${SESSION_ID}"
```

## Monitor Implementation

```go
type Session struct {
    ID             string
    PathHash       string
    TombstonePath  string
    TranscriptPath string
    ProjectDir     string
    StartTime      time.Time
    Model          string
}

func (m *Monitor) checkSessions() {
    sessions := m.listTombstones()

    for _, s := range sessions {
        index := m.loadIndex(s.PathHash)
        state := m.getState(s, index)

        switch state {
        case StateActive:
            health := m.analyzeHealth(index.Transcript)
            alerts := m.runAuditors(s, health)
            m.emitAlerts(s.ID, alerts)

        case StateSuspect:
            m.log.Debug("Session %s suspect", s.ID)

        case StateDead:
            m.log.Info("Cleaning up dead session %s", s.ID)
            m.cleanup(s)
        }
    }
}

func (m *Monitor) cleanup(s Session) {
    os.Remove(s.TombstonePath)
    os.Remove(filepath.Join(m.indexDir, s.PathHash+".json"))
    os.RemoveAll(filepath.Join(m.signalsDir, "sessions", s.ID))
}
```

## Edge Cases

| Scenario              | Behavior                                   |
|-----------------------|--------------------------------------------|
| Graceful exit         | SessionEnd cleans everything               |
| Terminal killed       | Heartbeat stops → Suspect → Dead → Cleanup |
| Long Claude operation | Transcript mtime fresh → stays Suspect     |
| Multiple sessions     | Each tracked independently                 |
| Orphan signals        | Cleaned with session directory             |
