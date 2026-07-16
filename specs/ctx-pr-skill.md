# Spec: `_ctx-pr` — scaffold a PR body from the branch

## Problem

Opening a PR for a ctx branch means hand-assembling a body from the
branch's commits, the `Spec:` trailers they cite, and the TASKS closed
along the way — then remembering ctx's process rules (DCO-only, never a
`Co-Authored-By:` or "Generated with…" footer). This is repetitive and
easy to get subtly wrong (a stray agent sign-off, a missed spec link).
There is no skill for it; `ctx-commit` covers the per-commit ceremony
but nothing composes the branch-level PR narrative.

## Design

A repo-internal skill `_ctx-pr` (lives in `.claude/skills/_ctx-pr/`,
alongside `_ctx-release`/`_ctx-qa`; the `_` prefix marks it
ctx-repo-only, never shipped in the plugin). It reads the branch and
writes a ready-to-paste PR body to the gitignored `inbox/`. It never
pushes and never opens the PR — the human does that (per the project's
no-push rule); the skill only drafts text.

### Inputs (all derived, no args required)

- Base branch: `main` unless the user names another.
- Branch commits: `git log --reverse <base>..HEAD` — subject, body, and
  `Spec:` trailer of each.
- Distinct specs cited across those commits (deduped, in first-seen
  order) → a "Specs" section linking design rationale.
- Closed tasks: entries in `.context/TASKS.md` that this branch flipped
  `[ ]`→`[x]` (`git diff <base>..HEAD -- .context/TASKS.md`).

### Output

`inbox/pr-<branch-slug>-<UTCstamp>.md` containing:

1. A title line (imperative, ≤ ~70 chars) synthesized from the commits.
2. A short **Summary** paragraph — what the branch does and why.
3. **Changes** — grouped bullets (by commit or theme), each pointing at
   the spec it serves.
4. **Specs** — the deduped list of `specs/*.md` cited by the commits.
5. **Closed tasks** — the TASKS items the branch marked done.
6. **Verification** — how it was checked (lint/test/audit, any
   end-to-end drive), pulled from commit bodies.

### Hard constraints (enforced by the skill's own checklist)

- The PR body MUST NOT contain `Co-Authored-By:`, any agent/tool
  sign-off, or a "Generated with…" / "🤖" footer — per CONSTITUTION
  "Process Invariants". This mirrors the commit rule that `ctx-commit`
  enforces, applied to the PR narrative.
- The skill does not run `git push`, `gh pr create`, or any remote
  operation. Output is a local file only.
- If `<base>..HEAD` is empty (e.g. work committed straight to `main`),
  the skill says so and does not fabricate a PR body.

## Tests / verification

Skills are prose, not Go — verification is a dry run: invoke `_ctx-pr`
on a branch with ≥2 commits and confirm the drafted body lists every
spec and closed task, and contains no forbidden sign-off/footer. Add a
line to the skill's own checklist so each run self-checks the forbidden
strings.

## Acceptance

- `_ctx-pr` produces `inbox/pr-*.md` with title, summary, spec links,
  and closed-tasks list derived from the actual branch.
- The draft never contains `Co-Authored-By:` or a generated-with footer.
- The skill performs no push/PR-create; it only writes the file.
- `make audit` stays green (repo-internal skills are not Copilot-synced;
  no shipped-asset drift).

Closes the TASKS item "Create a /ctx-pr skill".
