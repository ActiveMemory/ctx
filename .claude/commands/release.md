---
description: "Run the full release process"
---

Execute the release process for ctx.

**Prerequisites:**
1. VERSION file is updated to the new version
2. Release notes exist at `dist/RELEASE_NOTES.md` (generate with `/release-notes`)
3. Working tree is clean (all changes committed)

**What this does:**
1. Validates release notes exist
2. Checks working tree is clean
3. Updates version references in docs/index.md
4. Rebuilds the documentation site
5. Commits the docs update
6. Runs tests and smoke tests
7. Builds binaries for all platforms
8. Creates and pushes a signed git tag
9. Updates the "latest" tag

```!
make release
```

After completion, create the GitHub release at the URL shown and upload the binaries from `dist/`.
