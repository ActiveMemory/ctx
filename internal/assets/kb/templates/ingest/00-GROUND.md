# Ground Mode

External-grounding pass. This mode reaches **outside** the
project to pull authoritative facts that the kb depends on.
Internal grounding (reading the project's own code, specs,
transcripts) belongs in `ctx kb ingest`, not here.

Read `KB-RULES.md` first; it carries the editorial contract,
the demotion policy, evidence discipline, and the closeout
shape. This file describes the per-mode procedure only.

---

## Inputs

The skill reads `grounding-sources.md`. Each non-comment,
non-blank line is one of:

- `mcp:<server>:<resource>` — an MCP-reachable resource.
- `https://...` — a URL the skill should fetch (prefer a
  grounded-search MCP when available; web fetch is fallback).
- `./relative/path` or `/abs/path` — a local folder or file
  **outside** `.context/`.
- `symlink:<target>` — a stable symlink the project uses for
  cross-repo grounding.
- `NONE` — explicit no-op for this pass; re-prompts on next
  invocation. **Not** a project-wide opt-out.

Empty file (no non-comment lines) ⇒ the skill prompts once
for sources before doing any work and exits. It does **not**
fabricate sources, and it does **not** scan the project for
"likely" external references.

---

## Pre-Write Gates

Before the first fetch, confirm:

1. `.context/` and `.context/kb/` exist.
2. `.context/kb/index.md` has a non-placeholder `## Scope`
   section (no `<!-- TODO: ... -->` marker present).
3. `grounding-sources.md` has at least one resolvable line.

If any gate fails, refuse cleanly with the recovery hint
named in `KB-RULES.md` and `OPERATOR.md`.

---

## Constraints

- Ground does **not** rewrite kb prose by itself. It produces
  `evidence-index.md` rows; the human (or a follow-up
  `ctx kb ingest` pass) weaves new claims into glossary,
  domain-decisions, contradictions, timeline.
- A `low`-confidence claim grounded externally promotes one
  band **only** when the external source is independent of
  the original. Two transcripts of the same meeting do not
  count as independent.
- Out-of-repo citations omit the `sha:` field on their
  evidence rows; in-repo grounding is not this mode's job.
- If a source returns nothing (fetch failed, MCP unreachable,
  page 404), name the failure in the closeout's
  `Next pass hint` and continue with the remaining sources.

---

## Procedure

1. **Resolve sources.** Read `grounding-sources.md`; expand
   each non-comment line to a concrete fetchable target.
2. **Fetch.** Call the appropriate retrieval mechanism per
   source kind. Prefer a grounded-search MCP for URLs.
3. **Extract.** Pull the smallest set of atomic claims that
   bear on the kb's declared scope. Cite each by source
   short name (defined in `source-map.md`) plus a locator.
4. **Advance the ledger.** Update `source-coverage.md` for
   each source touched in this pass per the state-machine
   rules in `KB-RULES.md`. Honest state only — refusing to
   record incomplete coverage is the anti-pattern.
5. **Reconcile.** Compare new claims against `glossary.md`,
   `domain-decisions.md`, `contradictions.md`. Where external
   evidence diverges from internal claims, add a
   `contradictions.md` row and demote per the demotion
   policy.
6. **Append evidence.** Add new rows to `evidence-index.md`.
   Never renumber existing rows.
7. **Write the closeout** under
   `.context/ingest/closeouts/<TS>-ground-closeout.md` with
   the frontmatter shape from `KB-RULES.md` (`mode: ground`).
8. **Append a `SESSION_LOG.md` line** at the closeout phase
   boundary in the exact shape described in `KB-RULES.md`.
