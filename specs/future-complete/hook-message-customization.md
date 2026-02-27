# Hook Message Customization

Phase 2 of hook message templates. Phase 1 (spec:
`specs/hook-message-templates.md`) externalized all 24 hook messages
into text templates with a 3-tier fallback. This spec covers
**discoverability**: a CLI command, a metadata registry, and
documentation so users can find and customize messages without reading
Go source code.

## Problem

The message override mechanism works but is invisible. A user who
wants to customize the QA gate for a Python project (`pytest` instead
of `make lint`) has no way to discover:

1. That customization is possible
2. Which hooks have customizable messages
3. What template variables are available
4. Where to place override files
5. What the current effective message is

The only path today is reading Go source code or knowing the
`.context/hooks/messages/{hook}/{variant}.txt` convention by heart.

## Approach

Three additions:

1. **Hook message registry** (`internal/assets/hooks/messages/registry.go`):
   static metadata for all 24 message entries. Each entry carries hook
   name, variant, category (customizable vs ctx-specific), a one-line
   description, and available template variables.

2. **`ctx system message` CLI command**: four subcommands (`list`,
   `show`, `edit`, `reset`) that expose the registry and override
   mechanism to users.

3. **Documentation**: a recipe, CLI reference updates, configuration
   page updates, and cross-links.

### Message Categories

Not all messages should be customized. Two categories:

| Category | Count | Rationale |
|----------|-------|-----------|
| **customizable** | 14 | Message is an opinion; logic is universal. Projects benefit from project-specific wording. |
| **ctx-specific** | 10 | Message is internal to ctx's development workflow. Customization is possible but unusual. |

**Customizable** (14):
- check-context-size/checkpoint
- check-persistence/nudge
- check-ceremonies/{both,remember,wrapup}
- check-journal/{both,unexported,unenriched}
- check-knowledge/warning
- check-map-staleness/stale
- check-backup-age/warning
- qa-reminder/gate
- post-commit/nudge

**ctx-specific** (10):
- check-reminders/reminders
- check-resources/alert
- check-version/{mismatch,key-rotation}
- block-dangerous-commands/{mid-sudo,mid-git-push,cp-to-bin,install-to-local-bin}
- block-non-path-ctx/{dot-slash,go-run,absolute-path}

The `edit` subcommand warns for ctx-specific messages but does not
refuse. Advanced users may have valid reasons to customize them.

## Behavior

### `ctx system message list`

Prints a table of all 24 hook message entries:

```
Hook                    Variant             Category        Override
──────────────────────  ──────────────────  ──────────────  ────────
check-backup-age        warning             customizable
check-ceremonies        both                customizable
check-ceremonies        remember            customizable
...
block-non-path-ctx      dot-slash           ctx-specific
block-non-path-ctx      go-run              ctx-specific
...
qa-reminder             gate                customizable    override
```

The Override column shows "override" when a user override file exists
at `.context/hooks/messages/{hook}/{variant}.txt`.

With `--json`, outputs a JSON array:

```json
[
  {
    "hook": "check-backup-age",
    "variant": "warning",
    "category": "customizable",
    "description": "Backup staleness warning",
    "template_vars": ["Warnings"],
    "has_override": false
  },
  ...
]
```

### `ctx system message show <hook> <variant>`

Prints the effective message template (raw, not rendered):

```
Source: embedded default
Template variables: (none)

HARD GATE — DO NOT COMMIT without completing ALL of these steps...
```

When a user override exists:

```
Source: user override (.context/hooks/messages/qa-reminder/gate.txt)
Template variables: (none)

Run the test suite before committing.
Tests: pytest -x
Lint: ruff check .
```

### `ctx system message edit <hook> <variant>`

1. Validates hook/variant exist in the registry.
2. If override file already exists, refuses with a message:
   "Override already exists at {path}. Edit it directly or use
   `ctx system message reset` first."
3. If the message is ctx-specific, warns:
   "This message is ctx-specific (intended for ctx development).
   Customizing it may produce unexpected results."
   Then proceeds.
4. Copies the embedded default to
   `.context/hooks/messages/{hook}/{variant}.txt`,
   creating directories as needed.
5. Prints: "Override created at {path}. Edit this file to customize."
6. If template variables exist, lists them:
   "Available template variables: {{.PromptsSinceNudge}}"

### `ctx system message reset <hook> <variant>`

1. Validates hook/variant exist in the registry.
2. Deletes `.context/hooks/messages/{hook}/{variant}.txt`.
3. If the file didn't exist, prints: "No override found for
   {hook}/{variant}. Already using embedded default."
4. Cleans up empty parent directories.
5. Prints: "Override removed. Using embedded default."

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Invalid hook name | Error: "unknown hook: {name}" |
| Invalid variant for valid hook | Error: "unknown variant '{variant}' for hook '{hook}'" |
| `show` with no override | Shows embedded default, source line says "embedded default" |
| `edit` when override exists | Refuses, tells user to edit directly or reset first |
| `edit` for ctx-specific message | Warns but proceeds |
| `reset` when no override exists | Informational message, exit 0 |
| `.context/` dir doesn't exist | `edit` creates it; `list` shows no overrides |
| `list` with no `.context/` | All entries show no override |

### Validation Rules

- Hook and variant names validated against the registry, not the
  filesystem. Unknown names are errors even if a matching directory
  happens to exist.
- File paths constructed with `filepath.Join`, never string
  concatenation.

### Error Handling

| Error condition | Behavior | Recovery |
|-----------------|----------|----------|
| Unknown hook | Error message, exit 1 | User checks `list` output |
| Unknown variant | Error message, exit 1 | User checks `list` output |
| Override already exists (edit) | Informational refusal, exit 1 | User edits directly or resets |
| Failed to create directory (edit) | Error with path, exit 1 | Check permissions |
| Failed to write file (edit) | Error with path, exit 1 | Check permissions |
| Failed to read embedded (show) | Error, exit 1 | Should not happen; build issue |

## Interface

### CLI

```
ctx system message list [--json]
ctx system message show <hook> <variant>
ctx system message edit <hook> <variant>
ctx system message reset <hook> <variant>
```

| Subcommand | Args | Flags | Description |
|------------|------|-------|-------------|
| `list` | *(none)* | `--json` | Table of all hooks/variants with category and override status |
| `show` | `hook variant` | *(none)* | Print effective message template with source |
| `edit` | `hook variant` | *(none)* | Copy embedded default to .context/ for editing |
| `reset` | `hook variant` | *(none)* | Delete user override, revert to embedded default |

## Implementation

### Files to Create

| File | Purpose |
|------|---------|
| `internal/assets/hooks/messages/registry.go` | Static metadata for all 24 entries |
| `internal/cli/system/message_cmd.go` | CLI command implementation |
| `internal/cli/system/message_cmd_test.go` | Tests |
| `docs/recipes/customizing-hook-messages.md` | Recipe |

### Files to Modify

| File | Change |
|------|--------|
| `internal/assets/embed.go` | Add `ListHookMessages()`, `ListHookVariants()` |
| `internal/cli/system/system.go` | Register `messageCmd()` |
| `docs/cli/system.md` | Add `ctx system message` reference |
| `docs/home/configuration.md` | Add "Hook Message Overrides" subsection |
| `docs/recipes/index.md` | Add recipe entry |
| `docs/recipes/hook-output-patterns.md` | Add See Also cross-link |
| `docs/recipes/system-hooks-audit.md` | Add See Also cross-link |
| `zensical.toml` | Add recipe to nav |

### Key Functions

```go
// registry.go
type HookMessageInfo struct {
    Hook         string
    Variant      string
    Category     string   // "customizable" or "ctx-specific"
    Description  string
    TemplateVars []string
}

func Registry() []HookMessageInfo

// embed.go
func ListHookMessages() ([]string, error)
func ListHookVariants(hook string) ([]string, error)
```

### Helpers to Reuse

- `assets.HookMessage()` — read embedded template
- `assets.FS` — embedded filesystem
- `rc.ContextDir()` — resolve `.context/` path
- `loadMessage()` — existing template loading (reference, not called by CLI)
- `resources.go` — pattern for `--json` flag and table output

## Configuration

No new `.ctxrc` keys. Override mechanism remains convention-based.

## Testing

### Unit Tests (~17)

**list**:
- All 24 entries shown
- Override status detected when override file exists
- JSON output valid and complete
- Categories correct

**show**:
- Embedded default shown with source
- User override shown with source
- Template variables listed
- Invalid hook/variant errors

**edit**:
- Creates override file and directories
- Refuses when override exists
- ctx-specific warning present
- Template variables listed in output

**reset**:
- Deletes override file
- Informational message when no override exists

**Registry validation**:
- Every registry entry has a matching embedded file

Test pattern: `t.TempDir()` + `os.Chdir()` + `rc.Reset()`.

## Non-Goals

- **Templating the structural framing**: Box drawing, VERBATIM
  preamble, JSON encoding stay in Go code.
- **Template inheritance or composition**: Each variant is standalone.
- **Live preview / rendered output**: `show` displays raw template,
  not rendered output. Variables appear as `{{.Name}}`.
- **Bulk edit/reset**: One hook/variant at a time.
- **Validation of template syntax**: Users edit text files; syntax
  errors fall back to embedded defaults at runtime (existing behavior).
