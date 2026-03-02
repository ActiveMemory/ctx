# Sanitize Permissions Runbook

Manual procedure for cleaning up `.claude/settings.local.json`.
The agent may analyze and recommend, but **you** make every edit.

## Why a Runbook, Not a Skill

`settings.local.json` controls what the agent can do without asking.
An agent that can edit its own permission file is a self-escalation
vector — especially if the skill is auto-accepted. Keep this manual.

## When to Run

- After busy sessions where you clicked "Allow" many times
- Weekly hygiene (pair with `ctx drift`)
- Before committing `.claude/settings.local.json`

## Step 1: Snapshot

```bash
cp .claude/settings.local.json /tmp/settings-backup-$(date +%Y%m%d).json
```

## Step 2: Extract the Allow List

```bash
jq '.permissions.allow[]' .claude/settings.local.json | sort
```

Eyeball it. You're looking for four categories:

## Step 3: Identify Problems

### A. Garbage / Nonsense

Entries that are clearly broken or meaningless:

```
Bash(done)
Bash(__NEW_LINE_aa838494a90279c4__ echo "")
```

**Action**: Delete.

### B. One-Off Commands (Session Debris)

Entries with hardcoded paths, literal arguments, or exact commands
that were accepted during a specific debugging session:

```
Bash(git -C /home/jose/WORKSPACE/ctx log --oneline --all -20)
Bash(/home/jose/WORKSPACE/ctx/ctx add decision "Use PostgreSQL" --context ...)
```

Signs of a one-off:
- Full absolute paths to specific files
- Literal string arguments (not wildcards)
- Very specific flag combinations
- Commands that look like they came from a single task

**Action**: Delete unless you want to promote to a wildcard pattern.

### C. Subsumed Entries (Redundant)

A narrow entry that's already covered by a broader one:

```
# Narrow (redundant):
Bash(ctx recall list)
Bash(git -C /home/jose/WORKSPACE/ctx log --oneline -5)

# Broad (already covers the above):
Bash(ctx recall list:*)
Bash(git -C:*)
```

To find these, look for entries where removing the specific args
would match an existing wildcard entry.

**Action**: Delete the narrow entry.

### D. Duplicate Intent, Different Spelling

Same command with env vars in different order, or slight variations:

```
Bash(CGO_ENABLED=0 CTX_SKIP_PATH_CHECK=1 go test:*)
Bash(CTX_SKIP_PATH_CHECK=1 CGO_ENABLED=0 go test:*)
```

**Action**: Keep one, delete the other.

## Step 4: Check for Security Concerns

While you're in here, also flag:

| Pattern | Risk |
|---------|------|
| `Bash(git push:*)` | Bypasses block-git-push.sh hook |
| `Bash(rm -rf:*)` | Recursive delete, no confirmation |
| `Bash(sudo:*)` | Privilege escalation |
| `Bash(echo:*)`, `Bash(cat:*)` | Can compose into writes to sensitive files |
| `Bash(curl:*)`, `Bash(wget:*)` | Arbitrary network access |
| Any write to `.claude/` paths | Agent self-modification |

See the `ctx-sanitize-permissions` skill SKILL.md for the full threat matrix.

## Step 5: Edit

Edit `.claude/settings.local.json` directly in your editor.
Remove flagged entries. Keep the JSON valid.

```bash
# Validate JSON after editing
jq . .claude/settings.local.json > /dev/null && echo "valid" || echo "BROKEN"
```

## Step 6: Verify

```bash
# Compare before/after
diff /tmp/settings-backup-$(date +%Y%m%d).json .claude/settings.local.json
```

## Step 7: Optionally Commit

```bash
git add .claude/settings.local.json
git commit -m "chore: sanitize agent permissions"
```

## Asking the Agent for Help

You can safely ask the agent to *analyze* the file:

> "Look at my settings.local.json and tell me which permissions
> look like one-offs or are redundant."

The agent can read and report. **You** do the edits.

Do **not** add these to your allow list:
- `Skill(ctx-sanitize-permissions)`
- `Edit(.claude/settings.local.json)`
- Any `Bash(...)` pattern that writes to `.claude/`
