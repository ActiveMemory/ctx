---
name: ctx-sanitize-permissions
description: "Audit settings.local.json for dangerous permissions. Use periodically, after granting permissions, or when security hygiene matters."
---

Audit `.claude/settings.local.json` permissions for entries that
bypass safety hooks, grant overly broad access, or create injection
vectors. This is a defense-in-depth measure — hooks block dangerous
commands at runtime, but pre-approved permissions skip the
confirmation step that makes hooks visible.

## When to Use

- Periodically (weekly or after busy sessions)
- After granting several permissions in a session
- Before committing settings changes
- When the user asks "are my permissions safe?"
- Proactively if you notice the allow list is growing fast

## When NOT to Use

- Mid-flow when the user is actively debugging permission issues
- When the user explicitly says "I know what I'm doing"

## Usage Examples

```text
/ctx-sanitize-permissions
/ctx-sanitize-permissions (after that long session)
```

## Execution

### Step 1: Read Current Permissions

```bash
cat .claude/settings.local.json
```

Parse the `permissions.allow` array.

### Step 2: Check for Dangerous Patterns

Flag any permission matching these categories:

#### Category A: Hook Bypass (Critical)

These pre-approve commands that safety hooks are designed to
intercept. The hook still runs, but the user never sees the
confirmation dialog — so they cannot reject it.

| Pattern | Why Dangerous |
|---------|---------------|
| `Bash(git push:*)` | Bypasses block-git-push.sh hook confirmation |
| `Bash(git push)` | Same — exact match variant |
| `Bash(git push --force:*)` | Force push with no confirmation |

#### Category B: Destructive Commands (High)

| Pattern | Why Dangerous |
|---------|---------------|
| `Bash(rm -rf:*)` | Recursive delete with no confirmation |
| `Bash(git reset --hard:*)` | Discards uncommitted work |
| `Bash(git checkout .:*)` | Discards all unstaged changes |
| `Bash(git clean -f:*)` | Deletes untracked files |
| `Bash(git branch -D:*)` | Force-deletes branches |
| `Bash(sudo:*)` | Escalated privileges |

#### Category C: Config Injection Vectors (High)

These allow the agent to modify files that control its own behavior
— a self-modification vector that could be exploited via prompt
injection.

| Pattern | Why Dangerous |
|---------|---------------|
| Any `Bash(...)` that could write to `.claude/settings.local.json` | Agent modifies its own permissions |
| Any `Bash(...)` that could write to `CLAUDE.md` | Agent modifies its own instructions |
| Any `Bash(...)` that could write to `.claude/hooks/*.sh` | Agent modifies safety hooks |
| Any `Bash(...)` that could write to `.context/CONSTITUTION.md` | Agent modifies its own hard rules |

These are harder to detect by pattern alone. Look for overly broad
permissions like `Bash(echo:*)`, `Bash(cat:*)`, `Bash(tee:*)`,
`Bash(cp:*)` that could be composed into writes to sensitive paths.
Flag them as **informational** — they have legitimate uses but are
worth noting.

#### Category D: Overly Broad (Medium)

| Pattern | Why Dangerous |
|---------|---------------|
| `Bash(*:*)` or `Bash(*)` | Allows any command |
| `Bash(curl:*)` | Arbitrary network requests |
| `Bash(wget:*)` | Arbitrary downloads |
| `Bash(pip install:*)` | Arbitrary package installation |
| `Bash(npm install:*)` | Arbitrary package installation |

### Step 3: Check for Duplicates

Look for permissions that are redundant:
- Exact duplicates
- Entries where a broader pattern already covers a narrower one
  (e.g., `Bash(git:*)` makes `Bash(git status:*)` redundant)

### Step 4: Report

Format findings by severity:

```
## Permission Audit Results

### Critical (hook bypass)
- `Bash(git push:*)` — bypasses block-git-push.sh

### High (destructive / injection vector)
- `Bash(rm -rf:*)` — recursive delete, no confirmation

### Medium (overly broad)
- `Bash(curl:*)` — arbitrary network access

### Informational
- `Bash(cat:*)` — could compose into config file writes
- 3 duplicate entries found

### Clean
- 45 permissions reviewed, no issues found
```

### Step 5: Offer to Fix

For each finding, offer a specific action:

- **Critical/High**: "Remove this permission? (y/n)"
- **Medium**: "This is broad — do you want to keep it?"
- **Duplicates**: "Remove N duplicate entries?"
- **Informational**: Note only, no action needed

When removing permissions, edit `.claude/settings.local.json`
directly. Show the diff before and after.

## Important Notes

- NEVER remove permissions without showing the user what will
  be removed and getting confirmation
- Permissions the user just granted in this session are more
  likely intentional — note them but do not alarm
- Some broad permissions are legitimate for development
  workflows (e.g., `Bash(go test:*)`) — use judgment
- The goal is awareness, not lockdown. Flag risks, let the
  user decide

## Quality Checklist

After running the audit, verify:
- [ ] Read the actual settings file (did not guess)
- [ ] Checked all four categories (bypass, destructive,
      injection, broad)
- [ ] Checked for duplicates
- [ ] Reported findings by severity
- [ ] Offered specific fix actions for Critical/High
- [ ] Did not remove anything without confirmation
