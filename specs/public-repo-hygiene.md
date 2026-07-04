# Spec: public-repo hygiene for internal identifiers

## Problem

ctx is developed alongside internal, employer-proprietary
tooling, and this repository is public. Tracked files had
accumulated references that do not belong in a public tree:
internal hostnames, an internal mirror-repo's name, an internal
Jira server, and lineage notes naming a proprietary sibling
project in spec cross-references, task notes, and the
experimental-skill headers.

The operator's calibrated policy (2026-07-04): naming the
sibling project is tolerable, but internal *details* (designs,
hostnames, field IDs, org infrastructure) must never ship in
tracked files. Working-tree scrubs do not rewrite history; the
goal is that every commit from here forward is clean.

## Policy

- Tracked files describe designs on their own merits, without
  internal lineage or attribution (see specs/statusline.md's
  Decisions section for the pattern).
- Internal hostnames, mirror-repo names, Jira servers/field IDs,
  and employer-infrastructure details are redacted to neutral
  phrasing that preserves the sentence's meaning ("an internal
  mirror-repo issue", "the work vm").
- Comparative analyses of internal tooling live only in
  gitignored locations (`inbox/`, `ideas/`); `.context/journal/`
  is gitignored and therefore safe.
- Redactions preserve document structure: tasks are edited in
  place, never moved or deleted; decision records keep their
  meaning.
- The plain English word "anchor" (HTML anchors, cwd-anchored,
  Edit anchors) is never sweep-replaced.

## Non-Goals

- Git history rewriting (operator judged the historical
  name-drops acceptable; nothing design-level was ever
  committed)
- Scrubbing gitignored working notes

## Acceptance

- [ ] `git grep` over tracked files finds no internal hostnames,
      mirror-repo names, or Jira server references
- [ ] Redacted sentences still read naturally and truthfully
- [ ] TASKS.md/DECISIONS.md structure invariants intact
