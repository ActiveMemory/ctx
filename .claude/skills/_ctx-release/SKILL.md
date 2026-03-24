---
name: _ctx-release
description: "Run the full release process. Use when cutting a new version of ctx."
---

Execute the release process for Context CLI.

## Before Running

All three prerequisites must be true:

1. **VERSION** is updated to the new version number
2. **`dist/RELEASE_NOTES.md`** exists (generate with `/_ctx-release-notes`)
3. **Working tree is clean** (all changes committed)

If any prerequisite fails, stop. Running the release script with
missing notes or a dirty tree produces an incomplete or unsigned
tag that must be manually deleted.

## When to Use

- When cutting a tagged release of ctx
- When the user says "release", "ship it", or "cut a release"

## When NOT to Use

- When only generating release notes (use `/_ctx-release-notes`)
- When doing a dry run or preview

## Process

1. **Verify prerequisites**:
```bash
cat VERSION
test -f dist/RELEASE_NOTES.md && echo "Release notes: OK" || echo "MISSING"
git status --porcelain
```

2. **Run the release script**:
```bash
make release
```

This script:
- Updates version in 4 config files (plugin.json, marketplace.json, VS Code package.json + lock)
- Updates download URLs in 3 doc files (index.md, getting-started.md, integrations.md)
- Adds new row to versions.md
- Rebuilds the documentation site
- Commits the version and docs update
- Runs tests and smoke tests
- Builds binaries for all 6 platforms
- Creates and pushes a signed git tag
- Updates the `latest` tag

3. **After completion**, verify the GitHub release was created by CI
   at `https://github.com/ActiveMemory/ctx/releases`.

## Full Runbook

See [Cutting a Release](https://ctx.ist/operations/release/) for the
complete step-by-step guide including troubleshooting.

## Quality Checklist

- [ ] VERSION updated before running
- [ ] `dist/RELEASE_NOTES.md` exists
- [ ] Working tree is clean
- [ ] Tests and smoke tests pass
- [ ] Tag is pushed to origin
- [ ] GitHub release created by CI with all 6 binaries
