# Hook Message Templates

## Problem

All 16 system hook messages are hardcoded in Go source files. This creates
two problems:

1. **Project-specificity baked into the binary**: Messages like "lint the
   ENTIRE project" (qa-reminder) and "Use 'make build && sudo make install'"
   (block-dangerous-commands) are ctx development rules, not universal
   truths. A Python project needs different quality gates. A prototype
   needs none at all.

2. **No customization path**: Users who adopt ctx for their own projects
   inherit ctx's opinions as immutable behavior. The only escape hatch is
   removing the hook from hooks.json entirely — losing the hook's logic
   (counting, state tracking, adaptive frequency) just to change its words.

The hook *logic* is universal. The hook *messages* are opinions. These
should be separable.

## Approach

Externalize hook messages into text templates. Each hook loads its message
from a template file instead of a hardcoded string. Templates live in two
locations with a clear priority:

1. **User override**: `.context/hooks/messages/{hook}/{variant}.txt`
2. **Shipped default**: `internal/assets/hooks/messages/{hook}/{variant}.txt`

The hook code keeps all logic (when to fire, counter tracking, marker
files, adaptive frequency). Only the *content* — what text gets emitted —
comes from the template.

### Template Variables

Templates use Go `text/template` syntax for dynamic content:

```
No context files updated in {{.PromptsSinceNudge}}+ prompts.
Have you discovered learnings, made decisions,
established conventions, or completed tasks
worth persisting?

Run /ctx-wrap-up to capture session context.
```

Each hook defines a fixed set of available variables. Templates that
reference undefined variables render to the variable name literally
(no error, graceful degradation).

### Structural Framing vs. Content

The template contains only the **content** — the body text. The hook code
handles all **structural framing**:

- VERBATIM preamble (`"IMPORTANT: Relay this..."`)
- Box drawing (`┌─`, `│`, `└──`)
- JSON encoding (for block responses and additionalContext)
- `contextDirLine()` footer
- Webhook `notify.Send()` calls

This separation means:
- Users customize *what* is said, not *how* it's delivered
- Structural changes (new framing format, different box style) don't
  break user templates
- Templates are plain text, not JSON or Go code

### Empty Templates

An empty template file (0 bytes or whitespace-only) means "don't emit
a message." The hook still runs its logic (counting, state tracking) but
produces no output. This lets users silence specific messages without
removing the hook entirely.

### Hook Tiers

Not all hooks should be customizable:

| Tier | Hooks | Rationale |
|------|-------|-----------|
| **System** (hardcoded) | context-load-gate, cleanup-tmp, mark-journal | Infrastructure — the format *is* the feature |
| **Templated** (customizable) | All 13 others | Message is an opinion, logic is universal |

#### System Hooks — Not Templated

- **context-load-gate**: The injection format (file headers, separators,
  footer) is structural. The content comes from context files, not
  a message template.
- **cleanup-tmp**: Silent. No message to template.
- **mark-journal**: Plumbing. Informational output, not user-facing.

#### Templated Hooks — Full Inventory

**VERBATIM Relay (user-facing, box-drawn):**

| Hook | Variants | Template Variables |
|------|----------|-------------------|
| check-context-size | `checkpoint` | `PromptCount` |
| check-persistence | `nudge` | `PromptCount`, `PromptsSinceNudge` |
| check-ceremonies | `both`, `remember`, `wrapup` | *(none)* |
| check-journal | `both`, `unexported`, `unenriched` | `UnexportedCount`, `UnenrichedCount` |
| check-knowledge | `warning` | `FileWarnings` (formatted list) |
| check-map-staleness | `stale` | `LastRefreshDate`, `ModuleCount` |
| check-backup-age | `warning` | `Warnings` (formatted list) |
| check-reminders | `reminders` | `ReminderList` (formatted list) |
| check-resources | `alert` | `AlertMessages` (formatted list) |
| check-version | `mismatch`, `key-rotation` | `BinaryVersion`, `PluginVersion`, `KeyAgeDays` |

**Agent Directives (additionalContext):**

| Hook | Variants | Template Variables |
|------|----------|-------------------|
| qa-reminder | `gate` | *(none — static text)* |
| post-commit | `nudge` | *(none — static text)* |

**Block Responses (JSON decision:block):**

| Hook | Variants | Template Variables |
|------|----------|-------------------|
| block-dangerous-commands | `mid-sudo`, `mid-git-push`, `cp-to-bin`, `install-to-local-bin` | *(none — static text)* |
| block-non-path-ctx | `dot-slash`, `go-run`, `absolute-path` | *(none — static text)* |

## Behavior

### Happy Path

1. Hook fires (e.g., qa-reminder on Edit)
2. Hook runs its logic (state checks, counter updates, etc.)
3. Hook determines it should emit a message (variant: "gate")
4. Hook calls `loadMessage("qa-reminder", "gate", nil)`
5. `loadMessage` checks `.context/hooks/messages/qa-reminder/gate.txt`
6. File not found → falls back to embedded asset
   `internal/assets/hooks/messages/qa-reminder/gate.txt`
7. Template loaded, variables substituted (none in this case)
8. Hook wraps result in structural framing and emits

### User Override Path

1. User creates `.context/hooks/messages/qa-reminder/gate.txt`:
   ```
   Run the test suite before committing.
   Tests: pytest -x
   Lint: ruff check .
   ```
2. Next Edit triggers qa-reminder
3. `loadMessage` finds the user override first
4. Returns user's text instead of the default
5. Hook wraps it in the standard directive framing and emits

### Silence Path

1. User creates `.context/hooks/messages/check-ceremonies/both.txt`
   with empty content
2. Next session start triggers check-ceremonies
3. `loadMessage` finds the file, sees it's empty
4. Returns empty string
5. Hook sees empty message, skips emission
6. Hook logic (cooldown tracking, day counting) still runs normally

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Template file has bad syntax | Log warning, fall back to embedded default |
| Embedded default also missing | Fall back to hardcoded string in Go (belt and suspenders) |
| User template references unknown variable | Renders as `<no value>` (Go template default) |
| `.context/hooks/messages/` dir doesn't exist | Skip check, use embedded defaults |
| Template file is a directory | Skip, use embedded default |
| Template file permissions | Read-only required, no execution |
| Non-UTF8 content | Pass through as-is (hook output is bytes) |

### Error Handling

| Error condition | Behavior | Recovery |
|-----------------|----------|----------|
| User template parse error | Warn to stderr, use embedded default | Non-blocking |
| User template execute error | Warn to stderr, use embedded default | Non-blocking |
| Embedded asset missing | Use hardcoded Go fallback | Always works |
| Both user and embedded fail | Hardcoded fallback string | Belt and suspenders |

## Implementation

### Files to Create

| File | Purpose |
|------|---------|
| `internal/cli/system/message.go` | `loadMessage()` function and template loading |
| `internal/cli/system/message_test.go` | Tests for template loading priority and rendering |
| `internal/assets/hooks/messages/{hook}/{variant}.txt` | Default templates (13 hooks, ~25 variants) |

### Files to Modify

| File | Change |
|------|--------|
| `internal/assets/embed.go` | Add `hooks/messages/` to embed directive |
| `internal/cli/system/qa_reminder.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/post_commit.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_context_size.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_persistence.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_ceremonies.go` | Replace hardcoded strings with `loadMessage()` |
| `internal/cli/system/check_journal.go` | Replace hardcoded strings with `loadMessage()` |
| `internal/cli/system/check_knowledge.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_map_staleness.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_backup_age.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_reminders.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_resources.go` | Replace hardcoded string with `loadMessage()` |
| `internal/cli/system/check_version.go` | Replace hardcoded strings with `loadMessage()` |
| `internal/cli/system/block_dangerous_commands.go` | Replace hardcoded reasons with `loadMessage()` |
| `internal/cli/system/block_non_path_ctx.go` | Replace hardcoded reasons with `loadMessage()` |

### Key Implementation

#### message.go — Template Loader

```go
package system

import (
    "bytes"
    "os"
    "path/filepath"
    "strings"
    "text/template"

    "github.com/ActiveMemory/ctx/internal/assets"
    "github.com/ActiveMemory/ctx/internal/rc"
)

// loadMessage loads a hook message template by hook name and variant.
//
// Priority:
//  1. .context/hooks/messages/{hook}/{variant}.txt (user override)
//  2. internal/assets/hooks/messages/{hook}/{variant}.txt (embedded default)
//  3. fallback string (hardcoded, belt and suspenders)
//
// Returns empty string if the template file exists but is empty
// (intentional silence). The vars map provides template variables;
// nil is valid when no dynamic content is needed.
func loadMessage(hook, variant string, vars map[string]any, fallback string) string {
    relPath := filepath.Join("hooks", "messages", hook, variant+".txt")

    // 1. Check user override in .context/
    userPath := filepath.Join(rc.ContextDir(), relPath)
    if data, err := os.ReadFile(userPath); err == nil {
        return renderTemplate(string(data), vars, fallback)
    }

    // 2. Check embedded default
    if data, err := assets.FS.ReadFile(relPath); err == nil {
        return renderTemplate(string(data), vars, fallback)
    }

    // 3. Hardcoded fallback
    return renderTemplate(fallback, vars, fallback)
}

// renderTemplate executes a Go text/template with the given vars.
// Returns the fallback on any error. Returns empty string if the
// template content is empty (intentional silence).
func renderTemplate(tmpl string, vars map[string]any, fallback string) string {
    if strings.TrimSpace(tmpl) == "" {
        return "" // intentional silence
    }

    t, err := template.New("msg").Parse(tmpl)
    if err != nil {
        return fallback
    }

    var buf bytes.Buffer
    if err := t.Execute(&buf, vars); err != nil {
        return fallback
    }
    return buf.String()
}
```

#### Example Migration: qa-reminder

Before:
```go
msg := "HARD GATE — DO NOT COMMIT without completing ALL of these steps..."
```

After:
```go
msg := loadMessage("qa-reminder", "gate", nil,
    "HARD GATE — DO NOT COMMIT without completing ALL of these steps...")
```

The hardcoded string becomes the fallback — zero behavioral change if
templates are missing. Migration is mechanical: extract the string,
pass it as the last argument.

#### Example Migration: check-persistence (with variables)

Before:
```go
msg := fmt.Sprintf("...No context files updated in %d+ prompts...", sinceNudge)
```

After:
```go
msg := loadMessage("check-persistence", "nudge",
    map[string]any{
        "PromptCount":       state.Count,
        "PromptsSinceNudge": sinceNudge,
    },
    fmt.Sprintf("No context files updated in %d+ prompts...", sinceNudge))
```

Default template (`check-persistence/nudge.txt`):
```
No context files updated in {{.PromptsSinceNudge}}+ prompts.
Have you discovered learnings, made decisions,
established conventions, or completed tasks
worth persisting?

Run /ctx-wrap-up to capture session context.
```

#### Embed Directive Update

```go
//go:embed *.md Makefile.ctx entry-templates/*.md claude/skills/*/SKILL.md claude/.claude-plugin/plugin.json ralph/*.md tools/*.sh hooks/messages/*/*.txt
var FS embed.FS
```

### Directory Structure

```
internal/assets/hooks/messages/
├── block-dangerous-commands/
│   ├── mid-sudo.txt
│   ├── mid-git-push.txt
│   ├── cp-to-bin.txt
│   └── install-to-local-bin.txt
├── block-non-path-ctx/
│   ├── dot-slash.txt
│   ├── go-run.txt
│   └── absolute-path.txt
├── check-backup-age/
│   └── warning.txt
├── check-ceremonies/
│   ├── both.txt
│   ├── remember.txt
│   └── wrapup.txt
├── check-context-size/
│   └── checkpoint.txt
├── check-journal/
│   ├── both.txt
│   ├── unexported.txt
│   └── unenriched.txt
├── check-knowledge/
│   └── warning.txt
├── check-map-staleness/
│   └── stale.txt
├── check-persistence/
│   └── nudge.txt
├── check-reminders/
│   └── reminders.txt
├── check-resources/
│   └── alert.txt
├── check-version/
│   ├── mismatch.txt
│   └── key-rotation.txt
├── post-commit/
│   └── nudge.txt
└── qa-reminder/
    └── gate.txt
```

User overrides mirror this structure under `.context/hooks/messages/`.

### Helpers to Reuse

- `assets.FS` — existing `embed.FS` for compiled assets
- `rc.ContextDir()` — resolves `.context/` path
- `text/template` — Go stdlib, no new dependencies
- Existing hardcoded strings — become fallback arguments (zero-risk migration)

## Migration Strategy

### Phase 1: Add loadMessage() + defaults (no behavioral change)

1. Create `message.go` with `loadMessage()`
2. Extract all hardcoded strings into `internal/assets/hooks/messages/`
3. Migrate each hook to call `loadMessage()` with the hardcoded string
   as fallback
4. All tests pass with identical behavior — templates match current strings

### Phase 2: Documentation

1. Document the override mechanism in the prompting guide or a new
   "Customizing Hooks" page
2. Add `ctx init` awareness: optionally scaffold `.context/hooks/messages/`
   (empty dir, ready for overrides)

### Phase 3: Template variable expansion

1. For hooks with dynamic content, add template variables
2. Update default templates to use `{{.VarName}}` syntax
3. Update hook code to pass variable maps

Phase 1 and 3 can be combined if preferred — the mechanical extraction
is straightforward.

## Configuration

No new `.ctxrc` keys. The override mechanism is convention-based:
place a file at the right path, it gets picked up. No registration
or configuration needed.

## Testing

### Unit Tests — message_test.go

- **No override, no embedded**: returns fallback string
- **Embedded exists, no override**: returns embedded content
- **Override exists**: returns override content (not embedded, not fallback)
- **Empty override**: returns empty string (intentional silence)
- **Template with variables**: renders correctly
- **Template with unknown variable**: renders with `<no value>`, no error
- **Malformed template**: returns fallback, no panic
- **Override dir doesn't exist**: falls through to embedded
- **Variables map is nil**: works for static templates

### Integration Tests

- **End-to-end hook with default template**: qa-reminder produces
  expected output
- **End-to-end hook with user override**: qa-reminder produces
  custom message
- **Silence via empty template**: qa-reminder produces no output

## Non-Goals

- **Replacing hooks entirely via templates**: Templates control messages,
  not logic. To change *when* a hook fires or *what it checks*, modify
  hooks.json or write a custom hook script.
- **Template inheritance or composition**: Each variant is a standalone
  file. No includes, no partials, no base templates.
- **Localization / i18n**: Templates are English text files. The mechanism
  *could* support localization later (locale-prefixed dirs), but that's
  not a goal now.
- **Templating the structural framing**: Box drawing, VERBATIM preamble,
  JSON encoding — these stay in Go code. Templates are content only.
- **Runtime template reloading**: Templates are read on each hook
  invocation (hooks are short-lived processes). No caching or watch
  mechanism needed.
