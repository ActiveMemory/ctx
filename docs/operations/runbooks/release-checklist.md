---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Release Checklist
icon: lucide/clipboard-check
---

![ctx](../../images/ctx-banner.png)

# Release Checklist

The canonical pre-release sequence. This runbook ties together
the audits, tests, and release steps that are otherwise scattered
across docs and the operator's head.

**When to run**: Before every release. No exceptions.

**Companion**: The [`/_ctx-release`](../release.md) skill
automates the tag-and-push portion; this checklist covers
everything *before* and *after* that automation.

---

## Pre-Release

### 1. Run the Codebase Audit

Use the [codebase audit runbook](codebase-audit.md) prompt with
your agent. Focus on analyses 1-4 (extractable patterns,
documentation drift, maintainability, security). Triage findings
into TASKS.md — anything blocking ships before the release.

### 2. Run the Docs Semantic Audit

Use the [docs semantic audit runbook](docs-semantic-audit.md)
prompt. Fix high-severity findings (weak pages, broken narrative
arcs). Medium-severity items can be deferred.

### 3. Sanitize Permissions

Follow the [sanitize permissions runbook](sanitize-permissions.md).
Clean up `.claude/settings.local.json` before it gets committed
as part of the release.

### 4. Run the Full Test Suite

```bash
make audit    # fmt + vet + lint + drift + docs + test
make smoke    # integration smoke tests
```

All tests must pass. No exceptions.

### 5. Check Context Health

```bash
ctx drift          # broken references, stale patterns
ctx status         # context file health
/ctx-link-check    # dead links in docs
```

Fix anything flagged.

### 6. Review TASKS.md

Scan for incomplete tasks tagged as release-blocking. Either
finish them or explicitly defer with a reason in the task note.

---

## Release

### 7. Bump Version

```bash
echo "0.X.0" > VERSION
git add VERSION
git commit -m "chore: bump version to 0.X.0"
```

### 8. Generate Release Notes

In Claude Code:

```
/_ctx-release-notes
```

Review `dist/RELEASE_NOTES.md`. Ensure it captures all
user-visible changes.

### 9. Cut the Release

```bash
make release
```

Or in Claude Code: `/_ctx-release`. See
[Cutting a Release](../release.md) for the full step-by-step.

---

## Post-Release

### 10. Verify the GitHub Release

- [ ] [GitHub Releases](https://github.com/ActiveMemory/ctx/releases) shows the new version
- [ ] All 6 binaries are attached
- [ ] SHA256 checksums are attached
- [ ] Release notes render correctly

### 11. Update the Plugin Marketplace

If the plugin version changed, verify the marketplace entry:

```bash
claude /plugin list   # shows updated version
```

### 12. Announce

Post in the project's communication channels. Reference the
release notes.

### 13. Clean Up

```bash
rm dist/RELEASE_NOTES.md   # consumed by the release script
git stash pop              # if you stashed earlier
```
