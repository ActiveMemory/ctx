# model-context-window

## [2026-07-03-182238] New Claude model families silently fall to the 200k default in ModelContextWindow

**Context**: claude-fable-5 sessions warned '104 percent full' while the 1M window was 79 percent free; ModelContextWindow only recognized the [1m] suffix and the opus substring

**Lesson**: Unmapped model families inherit rc.DefaultContextWindow (200k), and every EffectiveContextWindow consumer (check-context-size hook, heartbeat, nudges, provenance) misreports against the wrong window

**Application**: At every model launch: add the family substring to internal/config/claude, extend TestModelContextWindow, and sanity-check against a live session's /context output before trusting hook percentages

---

## [2026-03-01-124921] Model-to-window mapping requires ordered prefix matching

**Context**: Implementing modelContextWindow() for the three-tier context window
fallback. Claude model IDs use nested prefixes (claude-sonnet-4-5 vs
claude-sonnet-4-20250514).

**Lesson**: A switch with ordered HasPrefix cases (most specific first) is
cleaner and safer than iterating separate prefix lists. The catch-all 'claude-*'
returns 200k for unrecognized Claude models.

**Application**: When adding new model families to modelContextWindow() in
session_tokens.go, add the most specific prefix first to avoid shadowing shorter
prefixes.

---

## [2026-02-27-230741] Doctor token_budget vs context_window confusion

**Context**: ctx doctor reported context size against token_budget (8k) instead
of context_window (200k), making 22k tokens look alarming.

**Lesson**: token_budget (ctx agent output trim target) and context_window
(model capacity) serve different purposes. Health checks about context fitting
should use context_window, with warning threshold proportional (e.g., 20% of
window).

**Application**: Doctor now uses rc.ContextWindow() with 20% threshold and shows
per-file token breakdown for actionable insight into which files are heavy.

---

