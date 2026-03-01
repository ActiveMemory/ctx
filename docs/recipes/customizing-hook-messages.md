---
title: "Customizing Hook Messages"
icon: lucide/message-square-text
---

![ctx](../images/ctx-banner.png)

## The Problem

ctx hooks speak ctx's language, not your project's. The QA gate says
"lint the ENTIRE project" and "make build," but your Python project uses
`pytest` and `ruff`. The post-commit nudge suggests running lints, but
your project uses `npm test`. You could remove the hook entirely, but
then you lose the *logic* (counting, state tracking, adaptive frequency)
just to change the *words*.

**How do you customize what hooks say without removing what they do?**

## TL;DR

```bash
ctx system message list                     # see all hooks and their messages
ctx system message show qa-reminder gate    # view the current template
ctx system message edit qa-reminder gate    # copy default to .context/ for editing
ctx system message reset qa-reminder gate   # revert to embedded default
```

## Commands Used

| Tool                       | Type        | Purpose                                                  |
|----------------------------|-------------|----------------------------------------------------------|
| `ctx system message list`  | CLI command | Show all hook messages with category and override status |
| `ctx system message show`  | CLI command | Print the effective message template                     |
| `ctx system message edit`  | CLI command | Copy embedded default to `.context/` for editing         |
| `ctx system message reset` | CLI command | Delete user override, revert to default                  |

---

## How It Works

Hook messages use a **3-tier fallback**:

1. **User override**: `.context/hooks/messages/{hook}/{variant}.txt`
2. **Embedded default**: compiled into the `ctx` binary
3. **Hardcoded fallback**: belt-and-suspenders safety net

The hook *logic* (*when to fire, counting, state tracking, cooldowns*) is
unchanged. Only the *content* (*what text gets emitted*) comes from
the template. You customize what the hook says without touching how it
decides to speak.

### Finding the Original Templates

The default templates live in the ctx source tree at:

```
internal/assets/hooks/messages/{hook}/{variant}.txt
```

You can also browse them on GitHub:
[`internal/assets/hooks/messages/`](https://github.com/ActiveMemory/ctx/tree/main/internal/assets/hooks/messages)

Or use `ctx system message show` to print any template without digging
through source code:

```bash
ctx system message show qa-reminder gate        # QA gate instructions
ctx system message show check-persistence nudge  # persistence nudge
ctx system message show post-commit nudge        # post-commit reminder
```

The `show` output includes the template source and available variables --
everything you need to write a replacement.

### Template Variables

Some messages use Go `text/template` variables for dynamic content:

```
No context files updated in {{.PromptsSinceNudge}}+ prompts.
Have you discovered learnings, made decisions,
established conventions, or completed tasks
worth persisting?
```

The `show` and `edit` commands list available variables for each message.
When writing a replacement, keep the same `{{.VariableName}}` placeholders
to preserve dynamic content. Variables that you omit render as `<no value>`:
no error, but the output may look odd.

### Intentional Silence

An **empty template file** (0 bytes or whitespace-only) means "*don't
emit a message*". The hook still runs its logic but produces no output.
This lets you silence specific messages without removing the hook from
`hooks.json`.

---

## Example: Python Project QA Gate

The default QA gate says "*lint the ENTIRE project*" and references
`make lint`. For a Python project, you want `pytest` and `ruff`:

```bash
# See the current default
ctx system message show qa-reminder gate

# Copy it to .context/ for editing
ctx system message edit qa-reminder gate

# Edit the override
```

Replace the content in `.context/hooks/messages/qa-reminder/gate.txt`:

```
HARD GATE! DO NOT COMMIT without completing ALL of these steps first:
(1) Run the full test suite: pytest -x
(2) Run the linter: ruff check .
(3) Verify a clean working tree
Run tests and linter BEFORE every git commit, no exceptions.
```

The hook still fires on every `Edit` call. The logic is identical. Only
the instructions changed.

---

## Example: Silencing Ceremony Nudges

The ceremony check nudges you to use `/ctx-remember` and `/ctx-wrap-up`.
If your team has a different workflow and finds these noisy:

```bash
ctx system message edit check-ceremonies both
ctx system message edit check-ceremonies remember
ctx system message edit check-ceremonies wrapup
```

Then empty each file:

```bash
echo -n "" > .context/hooks/messages/check-ceremonies/both.txt
echo -n "" > .context/hooks/messages/check-ceremonies/remember.txt
echo -n "" > .context/hooks/messages/check-ceremonies/wrapup.txt
```

The hooks still track ceremony usage internally, but they no longer emit
any visible output.

---

## Example: JavaScript Project Post-Commit

The default post-commit nudge mentions generic "lints and tests." For a
JavaScript project:

```bash
ctx system message edit post-commit nudge
```

Replace with:

```
Commit succeeded. 1. Offer context capture to the user: Decision (design
choice?), Learning (gotcha?), or Neither. 2. Ask the user: "Want me to
run npm test and eslint before you push?" Do NOT push. The user pushes
manually.
```

---

## The Two Categories

Not all messages are equal. The `list` command shows each message's
category:

### Customizable (16 messages)

Messages that are **opinions**: project-specific wording that benefits
from customization. These are the primary targets for override.

| Hook                | Variant    | Description                              |
|---------------------|------------|------------------------------------------|
| check-backup-age    | warning    | Backup staleness warning                 |
| check-ceremonies    | both       | Both ceremonies missing                  |
| check-ceremonies    | remember   | Start-of-session ceremony                |
| check-ceremonies    | wrapup     | End-of-session ceremony                  |
| check-context-size  | checkpoint | Context capacity warning                 |
| check-context-size  | oversize   | Injection oversize nudge                 |
| check-context-size  | window     | Context window usage warning (>80%)      |
| check-journal       | both       | Unexported sessions + unenriched entries |
| check-journal       | unenriched | Unenriched journal entries               |
| check-journal       | unexported | Unexported sessions                      |
| check-knowledge     | warning    | Knowledge file growth                    |
| check-map-staleness | stale      | Architecture map staleness               |
| check-persistence   | nudge      | Context persistence nudge                |
| post-commit         | nudge      | Post-commit context capture              |
| qa-reminder         | gate       | Pre-commit QA gate                       |

### ctx-specific (10 messages)

Messages specific to ctx's own development workflow. You *can* customize
them, but `edit` will warn you first.

| Hook                     | Variant              | Description                    |
|--------------------------|----------------------|--------------------------------|
| block-dangerous-commands | cp-to-bin            | Block copy to bin dirs         |
| block-dangerous-commands | install-to-local-bin | Block copy to ~/.local/bin     |
| block-dangerous-commands | mid-git-push         | Block git push                 |
| block-dangerous-commands | mid-sudo             | Block sudo                     |
| block-non-path-ctx       | absolute-path        | Block absolute path invocation |
| block-non-path-ctx       | dot-slash            | Block ./ctx invocation         |
| block-non-path-ctx       | go-run               | Block go run invocation        |
| check-reminders          | reminders            | Pending reminders relay        |
| check-resources          | alert                | Resource pressure alert        |
| check-version            | key-rotation         | Key rotation nudge             |
| check-version            | mismatch             | Version mismatch               |

---

## Template Variables Reference

| Hook                     | Variant                | Variables                                      |
|--------------------------|------------------------|------------------------------------------------|
| check-backup-age         | warning                | `{{.Warnings}}`                                |
| check-context-size       | checkpoint             | *(none)*                                       |
| check-context-size       | oversize               | `{{.TokenCount}}`                              |
| check-context-size       | window                 | `{{.TokenCount}}`, `{{.Percentage}}`           |
| check-ceremonies         | both, remember, wrapup | *(none)*                                       |
| check-journal            | both                   | `{{.UnexportedCount}}`, `{{.UnenrichedCount}}` |
| check-journal            | unenriched             | `{{.UnenrichedCount}}`                         |
| check-journal            | unexported             | `{{.UnexportedCount}}`                         |
| check-knowledge          | warning                | `{{.FileWarnings}}`                            |
| check-map-staleness      | stale                  | `{{.LastRefreshDate}}`, `{{.ModuleCount}}`     |
| check-persistence        | nudge                  | `{{.PromptsSinceNudge}}`                       |
| check-reminders          | reminders              | `{{.ReminderList}}`                            |
| check-resources          | alert                  | `{{.AlertMessages}}`                           |
| check-version            | key-rotation           | `{{.KeyAgeDays}}`                              |
| check-version            | mismatch               | `{{.BinaryVersion}}`, `{{.PluginVersion}}`     |
| post-commit              | nudge                  | *(none)*                                       |
| qa-reminder              | gate                   | *(none)*                                       |
| block-dangerous-commands | all variants           | *(none)*                                       |
| block-non-path-ctx       | all variants           | *(none)*                                       |

Templates that reference undefined variables render `<no value>`: 
no error, graceful degradation.

---

## Tips

* **Override files are version-controlled**: they live in `.context/`
  alongside your other context files. Team members get the same
  customized messages.
* **Start with `show`**: always check the current default before editing.
  The embedded template is the baseline your override replaces.
* **Use `reset` to undo**: if a customization causes confusion, reset
  reverts to the embedded default instantly.
* **Empty file = silence**: you don't need to delete the hook. An empty
  override file silences the message while preserving the hook's logic.
* **JSON output for scripting**: `ctx system message list --json` returns
  structured data for automation.

## See Also

* [Hook Output Patterns](hook-output-patterns.md): understanding
  VERBATIM relays, agent directives, and hard gates
* [Auditing System Hooks](system-hooks-audit.md): verifying hooks
  are running and auditing their output
* [Configuration](../home/configuration.md): project-level settings
  via `.ctxrc`
