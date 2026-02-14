---
name: ctx-borrow
description: "Extract and apply deltas between two directories of the same project. Use when merging results from a parallel worktree or separate checkout where git push/pull isn't practical."
allowed-tools: Bash(git:*), Bash(diff:*), Bash(ls:*), Bash(patch:*), Read
---

Extract changes from a "future" directory and apply them to the "present"
working tree. Complement to `ctx-worktree`: worktree splits work apart,
borrow merges results back.

## When to Use

- Merging results from a parallel worktree back into the main checkout
- Applying changes from a separate checkout (different machine, USB copy)
- User says "borrow", "bring changes from", "merge from folder"
- Two copies of the same project need their deltas reconciled

## When NOT to Use

- Directories are on a shared git remote — use `git pull` / `git merge`
- Worktree is managed by git and can use `ctx-worktree teardown`
- Directories are different projects (no common ancestry)

## Pre-flight Checks

1. **Validate directories exist**:
   ```bash
   ls "<present>" "<future>"
   ```

2. **Same-project fingerprint** — at least one must match:
   - Same `go.mod` module path, or same `package.json` name
   - Same `.context/CONSTITUTION.md` exists in both
   - Same basename after stripping worktree suffixes
   If none match, warn and ask the user to confirm before proceeding.

3. **Detect git status** in each directory:
   ```bash
   git -C "<present>" rev-parse --git-dir 2>/dev/null && echo "GIT" || echo "NOGIT"
   git -C "<future>" rev-parse --git-dir 2>/dev/null && echo "GIT" || echo "NOGIT"
   ```

## Strategy Selection

| Present | Future | Strategy |
|---------|--------|----------|
| git     | git    | **1** — `git fetch` + merge-base diff |
| git     | no-git | **2** — `git diff --no-index` |
| no-git  | git    | **2** — `git diff --no-index` |
| no-git  | no-git | **3** — `diff -rNu` fallback |

## Strategy 1: Both Git

```bash
# One-shot fetch — no permanent remote
git -C "<present>" fetch "<future>" HEAD --no-tags
# Find common ancestor
MERGE_BASE=$(git -C "<present>" merge-base HEAD FETCH_HEAD)
# Generate diff from merge-base to future HEAD
git -C "<present>" diff "$MERGE_BASE" FETCH_HEAD
```

If `merge-base` fails (no common history), fall back to Strategy 2.

## Strategy 2: `git diff --no-index`

```bash
git diff --no-index "<present>" "<future>" -- . \
  ':!.git' ':!.git/**'
```

This works whether either, both, or neither directory has `.git/`.

## Strategy 3: Plain `diff`

```bash
diff -rNu "<present>" "<future>" \
  --exclude='.git' --exclude='node_modules' --exclude='.venv'
```

## Delta Summary

After generating the diff, present a summary before applying:

```text
Delta: <future> → <present>
Strategy: <1|2|3>
Files changed: 12  (+340 / -85)
  modified:  internal/cli/pad.go, docs/getting-started.md, ...
  added:     internal/crypto/aes.go
  deleted:   hack/old-script.sh
```

**Always** show this summary and ask for confirmation before applying.

## Conflict Check

Cross-reference delta files with uncommitted changes in present:

```bash
git -C "<present>" status --porcelain
```

If any files appear in both the delta and the working tree changes,
warn explicitly:

```text
WARNING: These files have local changes AND incoming changes:
  - internal/cli/pad.go (modified locally + modified in delta)
Applying may cause conflicts or overwrite local work.
```

Ask the user how to proceed: apply anyway, skip conflicting files,
or abort.

## Applying the Delta

1. **Dry-run first** (Strategy 1/2):
   ```bash
   git apply --check --stat <patch-file>
   ```
   For Strategy 3:
   ```bash
   patch --dry-run -p0 < <patch-file>
   ```

2. **Apply** (after user confirms):
   ```bash
   git apply <patch-file>       # Strategy 1/2
   patch -p0 < <patch-file>     # Strategy 3
   ```

3. If apply fails with conflicts, report which hunks failed and
   offer to apply with `--reject` so the user can resolve manually.

## Selective Application

User can pick specific files. Regenerate the diff scoped to those files:

- Strategy 1: `git -C "<present>" diff "$MERGE_BASE" FETCH_HEAD -- file1 file2`
- Strategy 2: `git diff --no-index "<present>/file1" "<future>/file1"`
- Strategy 3: `diff -u "<present>/file1" "<future>/file1"`

## After Applying

Remind the user:
- Review changes with `git diff`
- Changes are unstaged — commit when satisfied
- If borrowing from a worktree, consider `ctx-worktree teardown` to clean up

## Guardrails

- **Always preview** — never apply without showing the delta summary first
- **No permanent remotes** — Strategy 1 uses one-shot fetch, no `git remote add`
- **Respect local work** — conflict check before every apply
- **Exclude noise** — always exclude `.git/`, `node_modules/`, `.venv/`
- **Reversible** — if present is a git repo, all changes can be reverted with `git checkout .`

## Quality Checklist

- [ ] Both directories validated as existing and same project
- [ ] Git status detected in both directories
- [ ] Correct strategy selected from the table
- [ ] Delta summary shown to user before applying
- [ ] Conflict check run against local working tree changes
- [ ] Dry-run passed before actual apply
- [ ] User confirmed before applying changes
