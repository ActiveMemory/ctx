# Project Skills Reference

This project uses Claude Code skills (`.claude/skills/`) to standardize
recurring workflows. Skills are invoked with `/skill-name` in a Claude Code
session. This page is for contributors who want to understand what each
skill does without reading individual `SKILL.md` files.

For ctx plugin skills (`/ctx-status`, `/ctx-recall`, etc.), see the
[ctx skills documentation](docs/skills.md).

## Skill Catalog

| Skill | Purpose | Invoke |
|-------|---------|--------|
| [absorb](#absorb) | Merge changes from a parallel directory | `/absorb` |
| [audit](#audit) | Detect and fix code-level drift | `/audit` |
| [backup](#backup) | Backup context to SMB share | `/backup` |
| [brainstorm](#brainstorm) | Design before implementation | `/brainstorm` |
| [check-links](#check-links) | Audit docs for dead links | `/check-links` |
| [qa](#qa) | Run QA checks before committing | `/qa` |
| [release](#release) | Full release process | `/release` |
| [release-notes](#release-notes) | Generate release notes | `/release-notes` |
| [sanitize-permissions](#sanitize-permissions) | Audit granted permissions | `/sanitize-permissions` |
| [skill-creator](#skill-creator) | Create or improve a skill | `/skill-creator` |
| [update-docs](#update-docs) | Sync docs with code changes | `/update-docs` |
| [verify](#verify) | Verify before claiming completion | `/verify` |

## Details

### absorb

Extracts a diff between two copies of the same project and applies it as a
patch. The companion to `/ctx-worktree` — worktree splits work apart, absorb
merges it back. Useful when `git push/pull` isn't practical (USB copies,
disconnected machines, worktrees without a shared remote).

- **Use when**: merging from a parallel worktree or separate checkout
- **Skip when**: directories share a git remote (just `git pull`)

### audit

Detects code-level drift: predicate naming, magic strings, hardcoded values,
missing godoc. Follows the 3:1 consolidation ratio — run after every ~3
rapid implementation sessions.

- **Use when**: after YOLO sprints, before releases
- **Skip when**: mid-feature with intentionally incomplete code

### backup

Backs up `.context/`, `.claude/`, and global Claude data to a configured SMB
share. Requires `CTX_BACKUP_SMB_URL` environment variable.

- **Use when**: before risky operations, end of productive sessions
- **Skip when**: SMB is not configured, or no changes since last backup

### brainstorm

Structured design thinking before implementation. Transforms vague ideas into
validated designs with trade-offs, constraints, and a recommended approach.

- **Use when**: before new features, architectural changes, behavior modifications
- **Skip when**: bug fixes with clear solutions, well-defined requirements

### check-links

Crawls documentation files for broken internal and external links.

- **Use when**: before releases, after restructuring docs
- **Skip when**: editing a single doc (just eyeball it)

### qa

Runs `go vet`, `go test`, build checks, and convention compliance before
committing. Catches issues locally instead of in CI.

- **Use when**: after writing Go code, before committing
- **Skip when**: only docs/markdown/config changed (no Go code touched)

### release

Runs the full release pipeline: version bump, build, release notes, tag, push.

- **Use when**: cutting a tagged release
- **Skip when**: only generating release notes (use `/release-notes` instead)

### release-notes

Generates `dist/RELEASE_NOTES.md` from commits since the last tag.
Required by `hack/release.sh` before it can run.

- **Use when**: preparing a release, drafting changelog
- **Skip when**: writing a blog post (use `/ctx-blog-changelog` instead)

### sanitize-permissions

Audits `settings.local.json` for overly broad or dangerous permission grants
that accumulated during development.

- **Use when**: periodically, after granting many permissions in a session
- **Skip when**: actively debugging permission issues

### skill-creator

Guides creation of new skills or evaluation of existing ones. Ensures skills
follow the project's structure and quality standards.

- **Use when**: building a new skill, improving an existing one
- **Skip when**: one-off instructions that don't warrant a reusable skill

### update-docs

Checks that documentation stays consistent with code after changes. Covers
CLI flags, file formats, defaults, and conventions.

- **Use when**: after modifying user-facing behavior, before committing
- **Skip when**: purely internal changes with no docs impact

### verify

Runs the actual commands (build, test, lint) to verify a claim before
reporting it. Prevents "tests pass" without having run the tests.

- **Use when**: before claiming anything is done, working, or fixed
- **Skip when**: documentation-only changes with no testable outcome
