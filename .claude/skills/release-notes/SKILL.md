---
name: release-notes
description: "Generate release notes for dist/RELEASE_NOTES.md. Use when preparing a release or when hack/release.sh needs release notes."
---

Generate release notes for the next release of Context CLI.

## Before Writing

1. **"Are there commits since the last tag?"** If HEAD equals the last
   tag, there is nothing to release.
2. **"Does VERSION match the intended release?"** If VERSION still shows
   the previous version, ask the user to update it first.

## When to Use

- Before running `hack/release.sh` (it requires `dist/RELEASE_NOTES.md`)
- When the user asks to prepare or draft release notes

## When NOT to Use

- For blog posts about a release (use `/ctx-blog-changelog`)
- When no commits exist since the last tag

## Process

1. **Read version and find baseline**:
```bash
cat VERSION
git describe --tags --abbrev=0 2>/dev/null
```

2. **Gather commits since last tag** (filter out noise):
```bash
# Substantive commits (skip "docs." and "minor" one-liners)
git log <tag>..HEAD --oneline --no-merges \
  | grep -v "^[a-f0-9]* docs\.$" \
  | grep -v "^[a-f0-9]* Doc update\.$" \
  | grep -v "^[a-f0-9]* minor"

# Changed file stats
git diff --stat <tag>..HEAD | tail -5

# Go code changes
git log <tag>..HEAD --oneline --no-merges -- '*.go'

# CI changes
git log <tag>..HEAD --oneline --no-merges -- '.github/**'

# Dependency changes
git log <tag>..HEAD --oneline --no-merges -- 'go.mod' 'go.sum'
```

3. **Read detailed commit messages** for each substantive commit:
```bash
git show <hash> --format="%B" --stat | head -30
```

4. **Check previous release notes** for format reference:
```bash
gh release view <tag> --json body -q '.body' | head -60
```

5. **Synthesize release notes**:
   - Write a summary (2-3 sentences) of what this release accomplishes
   - Group changes into logical sections
   - Write human-friendly descriptions, not raw commit messages
   - Highlight breaking changes prominently
   - Skip trivial changes (typo fixes, minor refactors)

6. **Write to `dist/RELEASE_NOTES.md`** and confirm.

## Output Format

```markdown
<img src="https://ctx.ist/images/ctx-banner.png" />

# Context CLI v<version>

<Summary paragraph>

## Canonical Release Narrative

(coming soon) https://ctx.ist/blog/

## Highlights

- **Bold label**: 1-2 sentence description

## Features

- Description of feature

## Bug Fixes

- Description of fix

## Refactoring

- Description of refactor (only if significant)

## CI

- Description of CI change (only if present)

## Documentation

- Description of doc changes

---

Full changelog: https://github.com/ActiveMemory/ctx/compare/<last-tag>...v<version>
```

## Quality Checklist

- [ ] VERSION file version matches the heading
- [ ] Every substantive commit is represented
- [ ] Sections only appear if they have content
- [ ] No raw commit hashes in the prose
- [ ] Banner image included at the top
- [ ] Changelog URL uses correct tag range
- [ ] Output written to `dist/RELEASE_NOTES.md`
- [ ] Ends with: "Release notes written to dist/RELEASE_NOTES.md"

## Style

- Active voice: "Add X" not "X was added"
- No em-dashes; use `:`, `;`, or restructure
- Straight quotes only (`"`, `'`)
- One line per bullet when possible
