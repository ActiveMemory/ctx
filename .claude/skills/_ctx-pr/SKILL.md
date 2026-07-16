---
name: _ctx-pr
description: "Scaffold a pull-request body from the current branch's commits, Spec: trailers, and closed TASKS into gitignored inbox/ for the human to paste. Use when preparing or opening a PR for a ctx branch. Never pushes and never creates the PR."
---

Draft a ready-to-paste PR body from what the branch actually
contains — its commits, the specs those commits cite, and the
TASKS it closed. Write it to `inbox/` (gitignored) for the human
to review and paste into the PR. **This skill never runs
`git push`, `gh pr create`, or any remote operation.** It only
writes a local file; the human opens the PR.

## When to Use

- Preparing or opening a PR for a feature/fix branch
- The user says "draft a PR", "write the PR body", "/_ctx-pr"
- Wrapping a branch and you want the narrative assembled from
  the real commits rather than hand-written

## When NOT to Use

- Work was committed straight to `main` (no branch to describe —
  `<base>..HEAD` is empty; say so and stop)
- The user wants the PR actually pushed/created — that is the
  human's job; this skill drafts text only
- Mid-branch, before the commits are final

## Before Running

1. Confirm the branch and base. Base is `main` unless the user
   names another. Current branch: `git rev-parse --abbrev-ref HEAD`.
2. Confirm there are commits to describe:
   `git rev-list --count <base>..HEAD`. If `0`, tell the user the
   branch has no commits over `<base>` and stop — do not fabricate
   a PR body.

## Gather (all derived — no guessing)

Run these and read the output; do not summarize from memory:

- **Commits, oldest first**, with full bodies:
  `git log --reverse --format='%H%n%s%n%b%n==END==' <base>..HEAD`
- **Specs cited** — the `Spec:` trailer of each commit, deduped in
  first-seen order:
  `git log --reverse --format='%(trailers:key=Spec,valueonly)' <base>..HEAD | sed '/^$/d' | sort -u`
  Every commit should have one (CONSTITUTION); if any lacks one
  (blank trailer output for that commit), note the gap rather than
  inventing a spec.
- **Closed tasks** — what this branch flipped `[ ]`→`[x]`:
  `git diff <base>..HEAD -- .context/TASKS.md` and read the added
  `- [x]` / `DONE` lines.
- **Verification** — pull any lint/test/audit or end-to-end drive
  notes out of the commit bodies; do not claim checks that the
  commits don't evidence.

## Write the PR body

Write to `inbox/pr-<branch-slug>-<UTCstamp>.md` where `<UTCstamp>`
is `date -u +%Y%m%dT%H%M%SZ`. Use this shape:

```markdown
# <imperative title synthesized from the commits, <= ~70 chars>

## Summary

<2-4 sentences: what the branch does and why. Ground it in the
commit subjects/bodies, not aspiration.>

## Changes

- <grouped bullet, by commit or theme> — serves `specs/<name>.md`
- <...>

## Specs

- `specs/<name>.md` — <one-line what it designs>
- <deduped, first-seen order>

## Closed tasks

- <TASKS item marked done on this branch>
- <...>  (omit the section if none)

## Verification

- <lint/test/audit result, end-to-end drive — only what the
  commits actually evidence>
```

## Hard constraints — self-check before finishing

Re-read the drafted file and confirm ALL of these. If any fails,
fix the file before reporting done:

- [ ] No `Co-Authored-By:` anywhere in the body
- [ ] No agent/tool sign-off and no "Generated with…" / "🤖"
      footer (CONSTITUTION "Process Invariants" — same rule
      `ctx-commit` enforces, applied to the PR narrative)
- [ ] Every `Changes` bullet ties to a spec listed in `Specs`
- [ ] `Closed tasks` entries actually appear as `[x]`/`DONE` in
      the branch's TASKS diff (no aspirational closes)
- [ ] `Verification` claims are backed by the commit bodies
- [ ] The file is under `inbox/` (gitignored) — never committed,
      never a tracked path
- [ ] No `git push` / `gh pr create` was run

## After Writing

Tell the user the file path and that it is theirs to paste — you
did not push or open anything. Offer to adjust tone or length.

## Usage Examples

```text
/_ctx-pr
/_ctx-pr (base develop)
/_ctx-pr (after finishing the ceremony-nudge branch)
```
