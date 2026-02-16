---
title: "Claude Code Permission Hygiene"
icon: lucide/shield
---

![ctx](../images/ctx-banner.png)

## The Problem

Claude Code's `.claude/settings.local.json` controls what the agent can do
without asking. Over time, this file accumulates one-off permissions from
individual sessions: Exact commands with hardcoded paths, duplicate entries,
and stale skill references. A noisy "*allowlist*" makes it harder to spot
dangerous permissions and increases the surface area for unintended behavior.

Since `settings.local.json` is `.gitignore`d, it drifts independently of your
codebase. There is no PR review, no CI check: just whatever you clicked
"*Allow*" on.

This recipe shows what a well-maintained permission file looks like and how to
keep it clean.

## Commands and Skills Used

| Command/Skill           | Role in this workflow                            |
|-------------------------|--------------------------------------------------|
| `ctx init`              | Populates default ctx permissions                |
| `/ctx-drift`            | Detects missing or stale permission entries      |
| `/sanitize-permissions` | Audits for dangerous patterns (security-focused) |

## Recommended Defaults

After running `ctx init`, your `settings.local.json` will have the ctx
defaults pre-populated. Here is an opinionated safe starting point for a Go
project using ctx:

```json
{
  "permissions": {
    "allow": [
      "Bash(/tmp/ctx-*:*)",
      "Bash(CGO_ENABLED=0 go build:*)",
      "Bash(CGO_ENABLED=0 go test:*)",
      "Bash(ctx:*)",
      "Bash(git add:*)",
      "Bash(git branch:*)",
      "Bash(git check-ignore:*)",
      "Bash(git checkout:*)",
      "Bash(git commit:*)",
      "Bash(git diff:*)",
      "Bash(git log:*)",
      "Bash(git remote:*)",
      "Bash(git restore:*)",
      "Bash(git show:*)",
      "Bash(git stash:*)",
      "Bash(git status:*)",
      "Bash(git tag:*)",
      "Bash(go build:*)",
      "Bash(go fmt:*)",
      "Bash(go test:*)",
      "Bash(go vet:*)",
      "Bash(golangci-lint run:*)",
      "Bash(grep:*)",
      "Bash(ls:*)",
      "Bash(make:*)",
      "Skill(ctx-add-convention)",
      "Skill(ctx-add-decision)",
      "Skill(ctx-add-learning)",
      "Skill(ctx-add-task)",
      "Skill(ctx-agent)",
      "Skill(ctx-archive)",
      "Skill(ctx-commit)",
      "Skill(ctx-drift)",
      "Skill(ctx-next)",
      "Skill(ctx-recall)",
      "Skill(ctx-reflect)",
      "Skill(ctx-remember)",
      "Skill(ctx-status)",
      "WebSearch"
    ]
  }
}
```

!!! note "This is a starting point, not a mandate"
    Your project may need more or fewer entries. 

    The goal is intentional permissions: Every entry should be there because
    **you** decided it belongs, not because you clicked "*Allow*" once during
    debugging.

### Design Principles

**Use wildcards for trusted binaries.** If you trust the binary (your own
project's CLI, `make`, `go`), a single wildcard like `Bash(ctx:*)` beats
twenty subcommand entries. It reduces noise and means new subcommands work
without re-prompting.

**Keep git commands granular.** Unlike `ctx` or `make`, git has both safe
commands (`git log`, `git status`) and destructive ones (`git reset --hard`,
`git clean -f`). Listing safe commands individually prevents accidentally
pre-approving dangerous ones.

**Pre-approve all ctx skills.** Skills shipped with ctx (`Skill(ctx-*)`) are
safe to pre-approve — they are part of your project and you control their
content. This prevents the agent from prompting on every skill invocation.

**Never pre-approve these:**

| Pattern                         | Risk                                           |
|---------------------------------|------------------------------------------------|
| `Bash(git push:*)`              | Bypasses `block-git-push` hook confirmation    |
| `Bash(git reset:*)`             | Can discard uncommitted work                   |
| `Bash(git clean:*)`             | Deletes untracked files                        |
| `Bash(rm -rf:*)`                | Recursive delete with no confirmation          |
| `Bash(sudo:*)`                  | Privilege escalation (also blocked by hook)    |
| `Bash(curl:*)` / `Bash(wget:*)` | Arbitrary network requests                     |
| `Skill(sanitize-permissions)`   | Edits this file — self-modification vector     |
| `Skill(release)`                | Runs release pipeline — high impact            |

## Hooks: Your Safety Net

Permissions and hooks work together. Even if a command is pre-approved, hooks
still run. The difference is that pre-approved commands skip the user
confirmation dialog: So if a hook blocks the command, the user never sees it.

The ctx plugin ships these blocking hooks:

| Hook                              | What it blocks                   |
|-----------------------------------|----------------------------------|
| `ctx system block-non-path-ctx`   | Running ctx from wrong path      |

Project-local hooks (not part of the plugin) can add more:

| Hook                          | What it blocks                   |
|-------------------------------|----------------------------------|
| `block-git-push.sh`           | All `git push` commands          |
| `block-dangerous-commands.sh` | `sudo`, copies to `~/.local/bin` |

!!! warning "Pre-approved + hook-blocked = silent block"
    If you pre-approve `Bash(git push:*)`, the hook still blocks it, but
    the user never sees the confirmation dialog. The agent gets a block
    response and must handle it, which is confusing.

    It's better not to pre-approve commands that hooks are designed to intercept.

## The Maintenance Workflow

### After busy sessions

Permissions accumulate fastest during debugging and exploration sessions.
After a session where you clicked "Allow" many times:

1. Open `.claude/settings.local.json` in your editor
2. Look for entries at the bottom of the allowlist (*new entries append there*)
3. Delete anything that looks session-specific:
   * Exact commands with hardcoded paths
   * Commands with literal string arguments
   * Entries that duplicate an existing wildcard

See [`hack/sanitize-permissions.md`](https://github.com/ActiveMemory/ctx/blob/main/hack/sanitize-permissions.md)
for a step-by-step runbook.

### Periodically

Run `/ctx-drift` to catch permission drift:

* Missing `Bash(ctx:*)` wildcard
* Missing `Skill(ctx-*)` entries for installed skills
* Stale `Skill(ctx-*)` entries for removed skills
* Granular `Bash(ctx <subcommand>:*)` entries that should be consolidated

Run `/sanitize-permissions` to catch security issues:

* Hook bypass patterns
* Destructive commands
* Overly broad permissions
* Injection vectors

### When adding new skills

If you create a custom `ctx-*` skill, add its `Skill()` entry to the
allowlist manually. 

`ctx init` only populates the default permissions: It won't pick up custom skills.

### Golden image snapshots

If manual cleanup is too tedious, use a golden image to automate it. Snapshot
a curated permission set, then restore at session start to automatically drop
session-accumulated permissions. See the
[Permission Snapshots](permission-snapshots.md) recipe for the full workflow.

## Adapting for Other Languages

The recommended defaults above are Go-specific. For other stacks, swap the
build/test tooling:

**Node.js / TypeScript:**
```json
"Bash(npm run:*)",
"Bash(npm test:*)",
"Bash(npx:*)",
"Bash(node:*)"
```

**Python:**
```json
"Bash(pytest:*)",
"Bash(python:*)",
"Bash(pip show:*)",
"Bash(ruff:*)"
```

**Rust:**
```json
"Bash(cargo build:*)",
"Bash(cargo test:*)",
"Bash(cargo clippy:*)",
"Bash(cargo fmt:*)"
```

The `ctx`, `git`, and skill entries remain the same across all stacks.

## Next Up

**[Permission Snapshots →](permission-snapshots.md)**: Save and restore
permission baselines for reproducible setups.

## See Also

* [Setting Up ctx Across AI Tools](multi-tool-setup.md): full setup recipe
  including `settings.local.json` creation
* [Context Health](context-health.md): keeping `.context/` files accurate
* [`hack/sanitize-permissions.md`](https://github.com/ActiveMemory/ctx/blob/main/hack/sanitize-permissions.md):
  manual cleanup runbook
