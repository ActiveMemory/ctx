---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Cutting a Release
icon: lucide/package
---

![ctx](../images/ctx-banner.png)

## Prerequisites

Before you can cut a release you need:

* Push access to `origin` (GitHub)
* GPG signing configured (`make gpg-test`)
* [Go](https://go.dev/) installed (version in `go.mod`)
* [Zensical](https://github.com/zensical/zensical) installed (`make site-setup`)
* A clean working tree (`git status` shows nothing to commit)

---

## Step-by-Step

### 1. Update the VERSION File

```bash
echo "0.9.0" > VERSION
git add VERSION
git commit -m "chore: bump version to 0.9.0"
```

The VERSION file uses bare semver (`0.9.0`), no `v` prefix.
The release script adds the `v` prefix for git tags.

### 2. Generate Release Notes

In Claude Code:

```
/_ctx-release-notes
```

This analyzes commits since the last tag and writes
`dist/RELEASE_NOTES.md`. The release script refuses to
proceed without this file.

### 3. Commit Any Remaining Changes

```bash
git status          # must be clean
make audit          # full check: fmt, vet, lint, test
```

### 4. Run the Release

```bash
make release
```

Or, if you are in a Claude Code session:

```
/_ctx-release
```

The release script does everything in order:

| Step | What happens |
|------|-------------|
| 1 | Reads `VERSION`, verifies release notes exist |
| 2 | Verifies working tree is clean |
| 3 | Updates version in 4 config files (plugin.json, marketplace.json, VS Code package.json + lock) |
| 4 | Updates download URLs in 3 doc files (index.md, getting-started.md, integrations.md) |
| 5 | Adds new row to versions.md |
| 6 | Rebuilds the documentation site (`make site`) |
| 7 | Commits all version and docs updates |
| 8 | Runs `make test` and `make smoke` |
| 9 | Builds binaries for all 6 platforms via `hack/build-all.sh` |
| 10 | Creates a signed git tag (`v0.9.0`) |
| 11 | Pushes the tag to origin |
| 12 | Updates and pushes the `latest` tag |

### 5. GitHub CI Takes Over

Pushing a `v*` tag triggers `.github/workflows/release.yml`:

1. Checks out the tagged commit
2. Runs the full test suite
3. Builds binaries for all platforms
4. Creates a GitHub Release with auto-generated notes
5. Uploads binaries and SHA256 checksums

### 6. Verify

- [ ] [GitHub Releases](https://github.com/ActiveMemory/ctx/releases) shows the new version
- [ ] All 6 binaries are attached (linux/darwin x amd64/arm64, windows x amd64)
- [ ] SHA256 files are attached
- [ ] Release notes look correct

---

## What Gets Updated Automatically

The release script updates 8 files so you do not have to:

| File | What changes |
|------|-------------|
| `internal/assets/claude/.claude-plugin/plugin.json` | Plugin version |
| `.claude-plugin/marketplace.json` | Marketplace version (2 fields) |
| `editors/vscode/package.json` | VS Code extension version |
| `editors/vscode/package-lock.json` | VS Code lock version (2 fields) |
| `docs/index.md` | Download URLs |
| `docs/home/getting-started.md` | Download URLs |
| `docs/operations/integrations.md` | VSIX filename version |
| `docs/reference/versions.md` | New version row + latest pointer |

The Go binary version is injected at build time via `-ldflags`
from the VERSION file. No source file needs editing.

---

## Build Targets Reference

| Target | What it does |
|--------|-------------|
| `make release` | Full release (script + tag + push) |
| `make build` | Build binary for current platform |
| `make build-all` | Build all 6 platform binaries |
| `make test` | Unit tests |
| `make smoke` | Integration smoke tests |
| `make audit` | Full check (fmt + vet + lint + drift + docs + test) |
| `make site` | Rebuild documentation site |

---

## Troubleshooting

### "Release notes not found"

```
ERROR: dist/RELEASE_NOTES.md not found.
```

Run `/_ctx-release-notes` in Claude Code first, or write
`dist/RELEASE_NOTES.md` manually.

### "Working tree is not clean"

```
ERROR: Working tree is not clean.
```

Commit or stash all changes before running `make release`.

### "Tag already exists"

```
ERROR: Tag v0.9.0 already exists.
```

You cannot release the same version twice. Either bump VERSION
to a new version, or delete the old tag if the release was
incomplete:

```bash
git tag -d v0.9.0
git push origin :refs/tags/v0.9.0
```

### CI build fails after tag push

The tag is already published. Fix the issue, bump to a patch
version (e.g. `0.9.1`), and release again. Do not force-push
tags that others may have already fetched.
