---
title: "Claude Code Permission Hygiene"
icon: lucide/shield
---

![ctx](../images/ctx-banner.png)

## The Problem

Claude Code's `.claude/settings.local.json` controls what the agent can do
without asking. Over time, this file accumulates one-off permissions from
individual sessions: Exact commands with hardcoded paths, duplicate entries,
and stale skill references. 

A noisy "*allowlist*" makes it harder to spot dangerous permissions and 
increases the surface area for unintended behavior.

Since `settings.local.json` is `.gitignore`d, it drifts independently of your
codebase. There is no PR review, no CI check: just whatever you clicked
"*Allow*" on.

This recipe shows what a well-maintained permission file looks like and how to
keep it clean.

## TL;DR

```bash
ctx init                            # seeds safe defaults
/ctx-drift                          # detects missing/stale permissions
/sanitize-permissions               # audits for dangerous patterns
```

See [Recommended Defaults](#recommended-defaults) for the full list.

## Commands and Skills Used

| Command/Skill           | Role in this workflow                            |
|-------------------------|--------------------------------------------------|
| `ctx init`              | Populates default ctx permissions                |
| `/ctx-drift`            | Detects missing or stale permission entries      |
| `/sanitize-permissions` | Audits for dangerous patterns (security-focused) |

## Recommended Defaults

After running `ctx init`, your `settings.local.json` will have the `ctx`
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
      "Skill(ctx-alignment-audit)",
      "Skill(ctx-archive)",
      "Skill(ctx-blog)",
      "Skill(ctx-blog-changelog)",
      "Skill(absorb)",
      "Skill(ctx-commit)",
      "Skill(ctx-context-monitor)",
      "Skill(ctx-drift)",
      "Skill(ctx-implement)",
      "Skill(ctx-journal-enrich)",
      "Skill(ctx-journal-enrich-all)",
      "Skill(ctx-journal-normalize)",
      "Skill(ctx-loop)",
      "Skill(ctx-next)",
      "Skill(ctx-pad)",
      "Skill(ctx-prompt-audit)",
      "Skill(ctx-recall)",
      "Skill(ctx-reflect)",
      "Skill(ctx-remember)",
      "Skill(ctx-status)",
      "Skill(ctx-worktree)",
      "WebSearch"
    ],
    "deny": [
      "Bash(sudo *)",
      "Bash(git push *)",
      "Bash(git push)",
      "Bash(rm -rf /*)",
      "Bash(rm -rf ~*)",
      "Bash(curl *)",
      "Bash(wget *)",
      "Bash(chmod 777 *)",
      "Read(**/.env)",
      "Read(**/.env.*)",
      "Read(**/*credentials*)",
      "Read(**/*secret*)",
      "Read(**/*.pem)",
      "Read(**/*.key)",
      "Edit(**/.env)",
      "Edit(**/.env.*)"
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

**Use wildcards for trusted binaries**: If you trust the binary (your own
project's CLI, `make`, `go`), a single wildcard like `Bash(ctx:*)` beats
twenty subcommand entries. It reduces noise and means new subcommands work
without re-prompting.

**Keep `git` commands granular**: Unlike `ctx` or `make`, git has both safe
commands (`git log`, `git status`) and destructive ones (`git reset --hard`,
`git clean -f`). Listing safe commands individually prevents accidentally
pre-approving dangerous ones.

**Pre-approve all `ctx-` skills**: Skills shipped with ctx (`Skill(ctx-*)`) are
safe to pre-approve. They are part of your project and you control their
content. This prevents the agent from prompting on every skill invocation.

### Default Deny Rules

`ctx init` automatically populates `permissions.deny` with rules that block
dangerous operations. Deny rules are evaluated before allow rules: A denied
pattern always prompts the user, even if it also matches an allow entry.

The defaults block:

| Pattern                  | Why                                          |
|--------------------------|----------------------------------------------|
| `Bash(sudo *)`           | Cannot enter password; will hang             |
| `Bash(git push *)`       | Must be explicit user action                 |
| `Bash(rm -rf /*)` etc.   | Recursive delete of system/home directories  |
| `Bash(curl *)` / `wget`  | Arbitrary network requests                   |
| `Bash(chmod 777 *)`      | World-writable permissions                   |
| `Read/Edit(**/.env*)`    | Secrets and credentials                      |
| `Read(**/*.pem, *.key)`  | Private keys                                 |

!!! note "Read/Edit Deny Rules"
    `Read()` and `Edit()` deny rules have known upstream enforcement issues
    (`claude-code#6631,#24846`). 

    They are included as defense-in-depth and intent documentation.

**Blocked by default deny rules**: no action needed, `ctx init` handles these:

| Pattern                         | Risk                                           |
|---------------------------------|------------------------------------------------|
| `Bash(git push:*)`              | Must be explicit user action                   |
| `Bash(sudo:*)`                  | Privilege escalation                           |
| `Bash(rm -rf:*)`                | Recursive delete with no confirmation          |
| `Bash(curl:*)` / `Bash(wget:*)` | Arbitrary network requests                     |

**Requires manual discipline**: **Never** add these to `allow`:

| Pattern                         | Risk                                      |
|---------------------------------|-------------------------------------------|
| `Bash(git reset:*)`             | Can discard uncommitted work              |
| `Bash(git clean:*)`             | Deletes untracked files                   |
| `Skill(sanitize-permissions)`   | Edits this file: self-modification vector |
| `Skill(release)`                | Runs the release pipeline: high impact    |

## Hooks: Regex Safety Net

Deny rules handle prefix-based blocking natively. Hooks complement them by
catching patterns that require regex matching: Things deny rules can't express.

The ctx plugin ships these blocking hooks:

| Hook                              | What it blocks                   |
|-----------------------------------|----------------------------------|
| `ctx system block-non-path-ctx`   | Running ctx from wrong path      |

Project-local hooks (not part of the plugin) catch regex edge cases:

| Hook                          | What it blocks                                                                    |
|-------------------------------|-----------------------------------------------------------------------------------|
| `block-dangerous-commands.sh` | Mid-command `sudo`/`git push` (after `&&`), copies to bin dirs, absolute-path ctx |


!!! warning "Pre-Approved + Hook-Blocked = Silent Block"
    If you pre-approve a command that a hook blocks, the user never sees
    the confirmation dialog. The agent gets a block response and must
    handle it, which is confusing.

    It's better not to pre-approve commands that hooks are designed to intercept.

## The Maintenance Workflow

### After busy sessions

Permissions accumulate fastest during debugging and exploration sessions.
After a session where you clicked "*Allow*" many times:

1. Open `.claude/settings.local.json` in your editor;
2. Look for entries at the bottom of the allowlist (*new entries append there*);
3. Delete anything that looks session-specific:
     * Exact commands with hardcoded paths,
     * Commands with literal string arguments,
     * Entries that duplicate an existing wildcard.

See [`hack/runbooks/sanitize-permissions.md`](https://github.com/ActiveMemory/ctx/blob/main/hack/runbooks/sanitize-permissions.md)
for a step-by-step runbook.

### Periodically

Run `/ctx-drift` to catch permission drift:

* Missing `Bash(ctx:*)` wildcard;
* Missing `Skill(ctx-*)` entries for installed skills;
* Stale `Skill(ctx-*)` entries for removed skills;
* Granular `Bash(ctx <subcommand>:*)` entries that should be consolidated.

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

If manual cleanup is too tedious, use a **golden image** to automate it: 

Snapshot a curated permission set, then restore at session start to automatically 
drop session-accumulated permissions. See the
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
* [`hack/runbooks/sanitize-permissions.md`](https://github.com/ActiveMemory/ctx/blob/main/hack/runbooks/sanitize-permissions.md):
  manual cleanup runbook
