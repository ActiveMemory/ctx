# Spec: Memory Bridge Foundation (`ctx memory`)

Mirror Claude Code's auto memory (MEMORY.md) into `.context/` with
drift detection, heuristic import, and bidirectional publish.

All three phases are implemented:

- **Foundation** (this spec): discovery, mirror, drift hook
- **Import**: `specs/memory-import.md` — classify and promote entries
- **Publish**: `specs/memory-publish.md` — push curated context back

## Problem

Claude Code maintains auto memory at `~/.claude/projects/<slug>/memory/MEMORY.md`.
This file is:

- **Outside the repo** — not version-controlled, not portable
- **Unstructured** — freeform Markdown, no separation of decisions vs
  learnings vs conventions
- **Ephemeral** — tied to one machine's `~/.claude/` directory
- **Invisible to ctx** — ctx hooks and context loading don't read it

Meanwhile, ctx maintains structured context files (DECISIONS.md,
LEARNINGS.md, CONVENTIONS.md, TASKS.md) that are git-tracked, portable,
and token-budgeted — but Claude Code doesn't read them unless hooks
inject them.

The two systems hold complementary knowledge with no bridge between them.

## Solution (This Phase)

Foundation layer that enables the bridge:

1. **Discovery**: locate MEMORY.md from the project root
2. **Mirror**: maintain a git-tracked copy + timestamped archives
3. **Drift detection**: hook nudge when MEMORY.md changes

Future phases (import and publish) build on this foundation.

## Discovery: Locating MEMORY.md

Claude Code encodes project paths as directory names under
`~/.claude/projects/`. The encoding replaces `/` with `-` and
prefixes with `-`:

```
/home/jose/WORKSPACE/ctx  →  ~/.claude/projects/-home-jose-WORKSPACE-ctx/
```

Discovery algorithm:

```go
func DiscoverMemoryPath(projectRoot string) (string, error) {
    // 1. Get absolute path of project root
    abs, _ := filepath.Abs(projectRoot)
    // 2. Replace "/" with "-", prefix with "-"
    slug := "-" + strings.ReplaceAll(abs[1:], "/", "-")
    // 3. Construct path
    home, _ := os.UserHomeDir()
    memPath := filepath.Join(home, ".claude", "projects", slug, "memory", "MEMORY.md")
    // 4. Check existence
    if _, err := os.Stat(memPath); err != nil {
        return "", fmt.Errorf("no auto memory found at %s", memPath)
    }
    return memPath, nil
}
```

Edge cases:
- Project root may differ across machines (different home dirs)
- Symlinks in the path may produce different slugs
- MEMORY.md may not exist yet (auto memory not triggered)

## Storage

### New directory: `.context/memory/`

```
.context/
├── memory/
│   ├── mirror.md                          # Raw copy of MEMORY.md (git-tracked)
│   └── archive/
│       ├── mirror-2026-03-05-143022.md    # Timestamped pre-sync snapshots
│       └── mirror-2026-03-04-220015.md
├── state/
│   └── memory-import.json                 # Import/sync tracking state
```

**mirror.md**: Exact copy of MEMORY.md at last sync. Git-tracked for
portability and diff visibility. Users who don't want this committed
can add `.context/memory/` to `.gitignore`.

**archive/**: Timestamped snapshots created before each mirror overwrite.
Git is the primary audit trail, but git is optional per project docs.
Archives provide a fallback.

**memory-import.json**: Tracks sync timestamps and (in future phases)
which entries have been imported.

```json
{
  "last_sync": "2026-03-05T14:30:22Z",
  "last_import": null,
  "last_publish": null,
  "imported_hashes": []
}
```

## CLI Surface

### `ctx memory sync`

Copy MEMORY.md to mirror. Archive previous mirror. Report drift.

```
$ ctx memory sync
✓ Archived previous mirror to archive/mirror-2026-03-05-143022.md
✓ Synced MEMORY.md → .context/memory/mirror.md
  Source: ~/.claude/projects/-home-jose-WORKSPACE-ctx/memory/MEMORY.md
  Lines: 47 (was 32)
  New content: 15 lines since last sync
```

Flags:
- `--dry-run`: show what would happen without writing

Exit codes:
- 0: synced successfully
- 1: MEMORY.md not found (auto memory not active)

### `ctx memory status`

Show drift, timestamps, and entry counts.

```
$ ctx memory status
Memory Bridge Status
  Source:      ~/.claude/projects/-home-jose-WORKSPACE-ctx/memory/MEMORY.md
  Mirror:      .context/memory/mirror.md
  Last sync:   2026-03-05 14:30 (2 hours ago)

  MEMORY.md:  47 lines (modified since last sync)
  Mirror:     32 lines
  Drift:      15 new lines detected
  Archives:   3 snapshots in .context/memory/archive/
```

Exit codes:
- 0: no drift
- 1: MEMORY.md not found
- 2: drift detected (MEMORY.md changed since last sync)

### `ctx memory diff`

Show what changed in MEMORY.md since last sync.

```
$ ctx memory diff
--- .context/memory/mirror.md (last sync: 2026-03-05 14:30)
+++ ~/.claude/projects/.../memory/MEMORY.md (current)
@@ -32,0 +33,15 @@
+## Session 2026-03-05: Memory Bridge Design
+
+### Key decisions
+- Import uses heuristic classification, not LLM
+- Publish happens on wrap-up, never on commit
+...
```

Implementation: unified diff between mirror.md and current MEMORY.md.

## Hook: `check-memory-drift`

### Registration

Add to `hooks.json` under `UserPromptSubmit`:

```json
{
  "command": "ctx system check-memory-drift",
  "timeout_ms": 3000
}
```

### Behavior

1. Discover MEMORY.md path for current project
2. Compare mtime against `.context/state/memory-last-sync`
3. If MEMORY.md is newer, output nudge via RESULT channel
4. Set session tombstone (`.context/state/memory-drift-nudged`)
   to suppress repeat nudges this session

### Output

```
┌─ Memory Drift ────────────────────────────────────
│ MEMORY.md has changed since last sync.
│ Run: ctx memory sync
│ Context: .context
└───────────────────────────────────────────────────
```

### Debounce

Once per session. The tombstone file is session-scoped (cleared on
session start, similar to existing nudge tombstones).

If MEMORY.md doesn't exist, skip silently (auto memory not active).

## Package Structure

```
internal/
├── memory/
│   ├── doc.go           # Package documentation
│   ├── discover.go      # DiscoverMemoryPath(), slug encoding
│   ├── discover_test.go
│   ├── mirror.go        # Sync(), Archive(), Diff()
│   ├── mirror_test.go
│   ├── state.go         # Sync state tracking (JSON)
│   └── state_test.go
├── cli/
│   └── memory/
│       ├── memory.go    # Parent command: ctx memory
│       ├── sync.go      # ctx memory sync
│       ├── status.go    # ctx memory status
│       └── diff.go      # ctx memory diff
```

## Changes Required

| File | Change |
|------|--------|
| `internal/config/dir.go` | Add `DirMemory = "memory"`, `DirMemoryArchive = "memory/archive"` |
| `internal/config/file.go` | Add `FileMemoryMirror = "mirror.md"`, `FileMemoryState = "memory-import.json"` |
| `internal/bootstrap/bootstrap.go` | Register `memory` parent command |
| `internal/assets/claude/hooks/hooks.json` | Add `check-memory-drift` hook |
| `internal/cli/system/memory_drift.go` | New: `ctx system check-memory-drift` |

## Error Cases

| Scenario | Behavior |
|----------|----------|
| MEMORY.md doesn't exist | `sync`: exit 1 with message. `status`: report "auto memory not active". Hook: skip silently. |
| `.context/memory/` doesn't exist | Create on first sync. |
| mirror.md doesn't exist | First sync creates it (no archive needed). |
| MEMORY.md is empty | Sync creates empty mirror. |
| No `.context/` (not initialized) | Init guard rejects (existing behavior). |

## Dependencies

- `internal/config`: new constants
- `internal/rc`: context dir resolution
- `internal/assets`: hook template registration

No new external dependencies.

## Testing Strategy

| Test | Type | Scope |
|------|------|-------|
| Slug encoding roundtrip | Unit | `discover.go` |
| Discovery with various home dirs | Unit | `discover.go` (HOME isolation) |
| Sync creates mirror + archive | Unit | `mirror.go` |
| Sync without prior mirror (first run) | Unit | `mirror.go` |
| Diff between mirror and source | Unit | `mirror.go` |
| State load/save roundtrip | Unit | `state.go` |
| CLI output formatting | Unit | each CLI file |

Test fixtures: create a realistic MEMORY.md in `testdata/` with
mixed entry types and session summaries.
