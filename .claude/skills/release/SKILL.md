---
name: release
description: "Run the full release process. Use when cutting a new version of ctx."
---

Execute the release process for Context CLI.

## Before Running

All three prerequisites must be true:

1. **VERSION** is updated to the new version number
2. **`dist/RELEASE_NOTES.md`** exists (generate with `/release-notes`)
3. **Working tree is clean** (all changes committed)

If any prerequisite fails, stop and tell the user what to fix.

## When to Use

- When cutting a tagged release of ctx
- When the user says "release", "ship it", or "cut a release"

## When NOT to Use

- When only generating release notes (use `/release-notes`)
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
- Updates version references in `docs/index.md`
- Rebuilds the documentation site
- Commits the docs update
- Runs tests and smoke tests
- Builds binaries for all platforms
- Creates and pushes a signed git tag
- Updates the `latest` tag

3. **After completion**, tell the user to create the GitHub release
   at the URL shown in the script output and upload binaries from
   `dist/`.

## Quality Checklist

- [ ] VERSION updated before running
- [ ] `dist/RELEASE_NOTES.md` exists
- [ ] Working tree is clean
- [ ] Tests and smoke tests pass
- [ ] Tag is pushed to origin
- [ ] User reminded to create GitHub release
