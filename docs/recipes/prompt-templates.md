---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Reusable Prompt Templates"
icon: lucide/message-square-text
---

![ctx](../images/ctx-banner.png)

## The Problem

The gap between "raw instruction typed into chat" and "full SKILL.md skill"
has no lightweight option. A user who wants a consistent code review
checklist or refactoring guard rails has two bad choices:

- Author a full skill (high friction, needs frontmatter, rebuild to embed)
- Keep prompts in their head or a scratch file (no discoverability, no sharing)

## TL;DR

```bash
ctx init                         # stamps starter prompts
ctx prompt list                  # see available prompts
ctx prompt show code-review      # print a prompt
ctx prompt add my-prompt --stdin # create from stdin
ctx prompt rm my-prompt          # delete
```

Or in your AI assistant: `/ctx-prompt code-review`

## Commands and Skills Used

| Tool                          | Type    | Purpose                                    |
|-------------------------------|---------|--------------------------------------------|
| `ctx prompt list`             | Command | List available prompt templates             |
| `ctx prompt show <name>`      | Command | Print prompt content to stdout              |
| `ctx prompt add <name>`       | Command | Create from embedded template or stdin      |
| `ctx prompt rm <name>`        | Command | Delete a prompt template                    |
| `ctx init`                    | Command | Stamps starter prompts during initialization |
| `/ctx-prompt`                 | Skill   | List or apply prompt templates in-session   |

## The Workflow

### Creating Prompts

**From starter templates** — `ctx init` stamps three starters:

- `code-review` — review checklist anchored to project conventions
- `refactor` — refactoring with guard rails (tests first, preserve behavior)
- `explain` — explain code for onboarding and knowledge transfer

**From embedded templates:**

```bash
ctx prompt add code-review    # creates from built-in template
ctx prompt add refactor       # creates from built-in template
```

**Custom prompts from stdin:**

```bash
echo "# Debug Checklist

1. Reproduce the issue
2. Check error logs
3. Add targeted logging
4. Isolate the failing component" | ctx prompt add debug --stdin
```

### Using Prompts

**In your AI assistant** — invoke the skill:

```text
/ctx-prompt code-review
```

The agent retrieves the prompt and follows its instructions in your
current context. If no name is given, it lists available prompts.

**From the CLI** — pipe into other tools:

```bash
ctx prompt show code-review    # print to stdout
```

### Sharing Prompts

Prompt templates live in `.context/prompts/` and are **committed to git
by default**. Your whole team shares the same prompts. For private prompts,
add `.context/prompts/` to `.gitignore`.

## Tips

**Keep prompts short and focused.** A good prompt template is 5-15 lines.
If it's longer, it's probably a skill.

**Anchor to project context.** Reference `.context/CONVENTIONS.md` or
`.context/ARCHITECTURE.md` in your prompts — the AI will read those files
for project-specific patterns.

**Name prompts for the action, not the content.** Use `code-review` not
`code-review-checklist`. The `.md` extension is added automatically.

**Prompts are not skills.** They have no frontmatter, no trigger rules, no
allowed-tools. They are plain markdown instructions. If you need automation,
create a skill instead.

## See Also

* [CLI Reference](../cli/index.md): full command documentation
* [Context Files](../home/context-files.md): structure of `.context/`
* [Detecting and Fixing Drift](context-health.md): keeping context clean
