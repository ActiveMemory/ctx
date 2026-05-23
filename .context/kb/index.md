# Knowledge Base

The kb is the project's evidence-tracked knowledge surface.
It lives at `.context/kb/` and is governed by the editorial
contract in `../ingest/KB-RULES.md`. Topic pages live under
`topics/<slug>/index.md`. Cross-cutting artifacts
(`glossary.md`, `domain-decisions.md`, `contradictions.md`,
`outstanding-questions.md`, `evidence-index.md`,
`source-coverage.md`, `timeline.md`) live alongside this
file.

---

## Scope

This kb captures design lessons and operational patterns the
ctx project has learned from adjacent / inspirational AI
infrastructure projects, including vLLM, Claude Code,
zensical, GitNexus, and similar single-purpose tools that
manage local project state. In scope: serving architectures,
CLI ergonomics, doc/recipe patterns, hook and skill surfaces,
and editorial workflows that ctx can lift or contrast against.
Deliberately NOT in scope: production deployment guides for
any of those tools (we are studying them, not running them),
end-user tutorials, or marketing-shaped overviews. Audience:
Volkan (project owner) and the `/ctx-kb-ingest`,
`/ctx-kb-ask`, `/ctx-kb-ground`, `/ctx-kb-site-review`,
`/ctx-kb-note` skills.

---

## Topics

<!-- CTX:KB:TOPICS START -->
- [`vllm`](topics/vllm/)
<!-- CTX:KB:TOPICS END -->

---

## Conventions

This kb is governed by:

- **`../ingest/KB-RULES.md`** is the editorial contract:
  pass-mode discipline, topic-page circuit breaker,
  source-coverage state machine, topic-adjacency pre-flight,
  cold-reader rubric, life-stage check, evidence discipline,
  confidence bands, demotion policy, closeout shape.
- **`../ingest/schemas/`** holds field-level shape for each
  cross-cutting artifact (`evidence-index.md`, `glossary.md`,
  `contradictions.md`, `outstanding-questions.md`,
  `domain-decisions.md`, `timeline.md`, `source-map.md`,
  `source-coverage.md`, `relationship-map.md`). Each schema
  ships shape, not content.
- **`../../specs/kb-editorial-pipeline.md`** is the full spec,
  including the failure-analysis section and the v1
  non-goals. Read this when you need to understand *why* a
  rule exists, not just *what* it says.

The mode skills (`/ctx-kb-ingest`, `/ctx-kb-ask`,
`/ctx-kb-site-review`, `/ctx-kb-ground`, `/ctx-kb-note`) are
the canonical writers. Hand-edits to kb files are an escape
hatch, not the primary path.
