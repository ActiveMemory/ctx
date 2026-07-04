# Spec: docs link integrity

## Problem

A user hit a 404 on the published site: the nav still pointed at
`recipes/activating-context.md`, deleted along with
`recipes/external-context.md` when `ctx activate` was dropped
(cwd-anchored change). A full sweep found three more breakage
classes: absolute `ctx.ist` URLs pointing at pages that moved in
old restructures (README's cli-reference/context-files/
integrations/manifesto links, a blog autonomous-loop link, a
self-referencing configuration URL), ten broken intra-doc anchors
(commands re-homed to new CLI pages, a removed `/ctx-prompt`
skill still listed), and the checked-in `site/` build predating
all of it.

## Fix Policy

- `zensical.toml` nav entries must resolve to files in `docs/`;
  entries for deliberately deleted pages are removed, not
  redirected.
- Absolute `https://ctx.ist/...` references must map to a current
  `docs/` page; each replacement target verified live (HTTP 200)
  before the edit.
- Anchor references follow the content: link to the page that
  documents the command today (`cli/bootstrap.md`,
  `init-status.md`), not the page that used to.
- Rows for removed skills are deleted rather than left as dead
  anchors.
- `site/` is rebuilt (`zensical build` + `ctx site feed`) so the
  deployed HTML carries the fixes; zensical's own anchor checker
  must report zero issues.

## Non-Goals

- Rewriting historical URLs inside `specs/` archives (they
  document past state)
- The `$id` URL in `ctxrc.schema.json` (a JSON Schema
  identifier, not a hyperlink; publishing the schema at that
  path is a separate decision)

## Acceptance

- [ ] All zensical.toml nav entries resolve
- [ ] `zensical build` reports zero anchor issues
- [ ] Internal link sweep over `docs/` reports zero broken file
      targets
- [ ] Replaced absolute URLs return HTTP 200
