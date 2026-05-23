# Chore-Class Commits: The Trailer Escape Hatch

A commit cites this spec **only** if its entire diff fits
one of the explicitly listed chore categories below.
Anything that doesn't fit a category needs its own spec
(or bundling into a functional commit that has one).

## Eligible chore categories

A commit may cite `specs/meta/chores.md` if and only if
the diff is entirely one of:

1. **Ignore-file additions.** New entries in `.gitignore`,
   `.dockerignore`, or similar — adding patterns to filter
   generated output that shouldn't be versioned. Removing
   entries is *not* a chore (it's a policy change about
   what becomes tracked).
2. **Dependency manifest bumps with no code change.**
   `go.mod`/`go.sum`/`package-lock.json` updates from
   `go get -u`, `npm update`, dependabot, or equivalent —
   no companion code changes required. If a bump requires
   a code change (API rename, breaking change adaptation),
   it's a functional commit and needs its own spec.
3. **Formatting passes.** `gofmt`, `prettier`, `goimports`
   normalization across files. No logic changes. Whitespace
   only.
4. **Typo / spelling fixes in comments or docs.** Single-
   letter or word-level corrections to comments,
   docstrings, README copy, error messages. *Not* in spec
   content (specs warrant their own commit), *not* in
   identifier names (renames are functional changes).
5. **Mechanical file moves with no logic change.** `git mv`
   operations where the moved file's content is unchanged
   except for import-path updates required by the new
   location. Renames that change behavior or interface
   are functional.
6. **License header / copyright year updates.** Bulk
   replacement of the project license header or copyright
   year across files. No code change beyond the header.

## Ineligible (require their own spec or bundling)

- Bug fixes of any size — even one-liners.
- Test additions, even regression-pin tests for already-
  working behavior.
- Configuration changes that alter runtime behavior
  (env vars, YAML defaults, RC keys).
- Documentation that adds new content (vs. fixing typos
  in existing content).
- Refactors, however small. "Renamed for clarity" is a
  functional change.

## Usage pattern

```
chore(gitignore): ignore .context/reports/

<commit body explaining what and why>

Spec: specs/meta/chores.md
```

The body still needs to explain *what* the change is and
*why* — citing the meta spec doesn't excuse a vague
commit message. The trailer just resolves the
"every-commit-needs-a-spec" requirement honestly when
there's no specific design rationale to point at.

## Anti-pattern

If you find yourself splitting a functional change into
"a fix + a chore that's needed for the fix to work,"
that's a sign the chore should be bundled into the
functional commit, not standalone. The functional
commit's spec covers the chore. The chore class is for
genuinely-standalone changes that don't have a parent.

## Audit

Periodically — at release time or during PR review —
scan `git log --grep="Spec: specs/meta/chores.md"` and
confirm each cited commit's diff is in the eligible
categories. Misuse drives the threshold for cleanup; if
the meta spec gets cited for non-chore commits, the
fix is to push back at PR time, not to expand the
category list.

## See Also

- `CONSTITUTION.md` — Process Invariants, "Every commit
  references a spec" rule with the chore escape hatch.
- `AGENT_PLAYBOOK.md` — Spec Verification Step procedure
  that gates use of this escape hatch.
- `specs/spec-trailer-discipline.md` — design rationale
  for why this escape hatch exists.
