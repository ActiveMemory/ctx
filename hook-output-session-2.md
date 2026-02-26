# Hook Output — Session 2 (2026-02-25)

**Model:** Claude Opus 4.6
**Trigger:** User ran `/plugin`, then asked "can you add a --verbose flag to the info command?"

---

## Hooks Received

**None.**

No `<user-prompt-submit-hook>` tags or hook output appeared in the session context.

---

## What DID Appear (non-hook context)

| # | Type | Source | Content |
|---|------|--------|---------|
| 1 | `<system-reminder>` | Skill list | ~50 available skills (both `_ctx-*` private and `ctx:*` public) |
| 2 | `<system-reminder>` | claudeMd | Full CLAUDE.md contents (session-start protocol, memory instructions, etc.) |
| 3 | `<local-command-caveat>` | `/plugin` command | `(no content)` — plugin command ran but produced empty output |
| 4 | Built-in | Git status snapshot | Branch, modified files, recent commits |
| 5 | Built-in | Auto memory | Path to persistent memory directory |

---

## Expected Hooks (from `internal/assets/claude/hooks/hooks.json`)

These are defined in the plugin but **none fired**:

### UserPromptSubmit (should fire on every user message)

| Hook Command | Expected Behavior |
|---|---|
| `ctx system check-bootstrap` | Emit "STOP, run bootstrap" nudge on first prompt |
| `ctx system check-context-size` | Warn if context window is filling up |
| `ctx system check-ceremonies` | Remind about wrap-up ceremony |
| `ctx system check-persistence` | Check if context persistence is healthy |
| `ctx system check-journal` | Journal-related checks |
| `ctx system check-reminders` | Surface pending reminders |
| `ctx system check-version` | Check for ctx version updates |
| `ctx system check-resources` | Resource usage checks |
| `ctx system check-knowledge` | Knowledge base checks |
| `ctx system check-map-staleness` | Check if architecture map is stale |

### PreToolUse (should fire before every tool call)

| Hook Command | Matcher | Expected Behavior |
|---|---|---|
| `ctx system context-load-gate` | `.*` | Gate context loading |
| `ctx system block-non-path-ctx` | `Bash` | Block non-path ctx invocations |
| `ctx system qa-reminder` | `Edit` | Remind about QA before edits |
| `ctx agent --budget 4000` | `.*` | Load context packet |

### PostToolUse (should fire after tool calls)

| Hook Command | Matcher | Expected Behavior |
|---|---|---|
| `ctx system post-commit` | `Bash` | Post-commit context capture |

### SessionEnd

| Hook Command | Expected Behavior |
|---|---|
| `ctx system cleanup-tmp` | Clean up temporary files |

---

## Comparison with Previous Session

| Aspect | Session 1 (hook-nudge-analysis.md) | Session 2 (this file) |
|---|---|---|
| Skills loaded | Yes | Yes |
| UserPromptSubmit hooks fired | Yes (check-bootstrap nudge received) | **No** |
| PreToolUse hooks fired | Unknown (agent jumped to task) | **No** |
| `/plugin` output | `(no content)` | `(no content)` |

---

## Conclusion

The plugin mechanism loaded **skills** successfully but did **not** wire up **hooks**.
Skills and hooks appear to follow different loading paths in the plugin system.
