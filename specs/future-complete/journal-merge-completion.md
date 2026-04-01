# Journal-recall merge completion and cross-cutting cleanup

## Problem

The journal-recall merge (specs/journal-recall-merge.md) shipped
incomplete. The journal/core/ packages (plan, query, confirm, execute,
extract, validate) were copied from recall/core/ but never wired —
journal commands still delegate to the old recall packages. Additionally,
multiple convention violations accumulated across the changeset.

## Issues

### Structural

1. **journal/cmd/source/run.go** delegates to recall/cmd/list and
   recall/cmd/show instead of journal/core/query
2. **journal/core/plan.Import** imports from recall/core/ — should
   use journal/core/ siblings
3. **journal/core/query.FindSessions** unused — should be called by
   journal/cmd/source
4. **sourcefm** and **sourceformat** should be clustered as
   source/frontmatter and source/format
5. **extract.ExtractFrontMatter** → **extract.FrontMatter** (stuttery)

### Magic numbers and strings

6. **source/cmd.go:75** — `"project"`, `"p"` hardcoded
7. **check_context_size/run.go** — 30, 3, 15 are magic numbers
8. **post_commit/run.go** — regexes, violation points, and localizable
   strings all hardcoded
9. **state/state.go:29** — 0o750 hardcoded

### Naming and conventions

10. **state.StateDir** → **state.Dir** (stuttery in state package)
11. **session.go splitLines** — utility belongs elsewhere, private
    function should be in separate file per convention
12. **sourcefm docstrings** — reference private function names but
    functions are public

### Skill generalization

13. **/ctx-commit SKILL.md** assumes Go project (CGO_ENABLED, go build,
    go files). Must be language-agnostic.
14. `ctx add decision` signature outdated (requires flags now)
15. Doc drift check references ctx-internal skill `_ctx-update-docs`
16. Reflect step marked optional guarantees it will be skipped

## Non-goals

- Deleting recall/core/ packages (may still have callers outside journal)
- Changing journal file formats or frontmatter schema
