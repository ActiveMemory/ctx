Generate release notes for the next release.

Steps:
1. Read the `VERSION` file to get the version number
2. Find the last git tag: `git describe --tags --abbrev=0 2>/dev/null`
3. Get commits since last tag (or all commits if no tags):
   - If tag exists: `git log <tag>..HEAD --pretty=format:"%h %s" --no-merges`
   - If no tags: `git log --pretty=format:"%h %s" --no-merges`
4. Get a summary of changed files: `git diff --stat <tag>..HEAD` (or `git diff --stat --root HEAD` if no tags)

Then synthesize release notes:
- Write a brief summary of what this release accomplishes (2-3 sentences)
- Group changes into logical sections (Features, Fixes, Documentation, etc.)
- Write human-friendly descriptions, not just commit messages
- Highlight breaking changes or important updates
- Skip trivial changes (typo fixes, minor refactors)

Write the output to `dist/RELEASE_NOTES.md` in this format:

```markdown
# Context CLI v<version>

<Brief summary paragraph>

## Highlights

- Key change 1
- Key change 2

## Features

- Description of feature

## Bug Fixes

- Description of fix

## Documentation

- Description of doc changes

---

Full changelog: https://github.com/ActiveMemory/ctx/compare/<last-tag>...v<version>
```

End by confirming: "Release notes written to dist/RELEASE_NOTES.md"
