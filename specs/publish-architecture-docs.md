# Spec: Publish Architecture Docs

## Goal

Publish selected architecture artifacts from `.context/` to `docs/`
for public consumption. Keep internal-only documents (danger zones,
principal analysis) in `.context/`.

## Publish Matrix

| Source | Destination | Treatment |
|--------|-------------|-----------|
| `.context/ARCHITECTURE.md` | `docs/reference/architecture.md` | As-is or light edit (remove drift-check comments) |
| `.context/DETAILED_DESIGN.md` | `docs/reference/detailed-design/index.md` | As-is |
| `.context/DETAILED_DESIGN-*.md` (5 files) | `docs/reference/detailed-design/` | As-is |
| `.context/CHEAT-SHEETS.md` | `docs/reference/cheat-sheets.md` | As-is |
| ARCHITECTURE-PRINCIPAL.md "Intervention Points" | `docs/contributing/where-to-contribute.md` | Sanitized: positive framing, contributor-oriented |
| ARCHITECTURE-PRINCIPAL.md "Vision Alignment" | `docs/contributing/multi-agent-support.md` | Sanitized: "how ctx supports any agent", not gap analysis |
| Convergence Report "Search Prompts" | `docs/reference/further-reading.md` | "Concepts" and "ADR" sections only, no confidence scores |

## Do NOT Publish

- `.context/DANGER-ZONES.md` — internal risk assessment
- `.context/ARCHITECTURE-PRINCIPAL.md` — strategic analysis, silent choices, risks
- `.context/map-tracking.json` — internal tracking
- `.context/CONVERGENCE-REPORT.md` — internal convergence scores

## Skill Enhancement: GLOSSARY.md Updates

The `/ctx-architecture` skill should conditionally update GLOSSARY.md
during Phase 3 with terms discovered during analysis.

### Behavior

- **When**: Phase 3 (Update Documents), after ARCHITECTURE.md and
  DETAILED_DESIGN updates
- **What**: Append new domain terms discovered during codebase
  analysis that are not already in GLOSSARY.md
- **How**: Additive only — never remove or rewrite existing entries.
  New terms inserted alphabetically.
- **Skip if**: GLOSSARY.md doesn't exist or has no existing entries
  (user hasn't opted into glossary)

### Term Sources

During analysis, the skill encounters terms that deserve glossary
entries:
- Internal jargon (e.g., "token budget", "drift check", "entry block")
- Architecture concepts (e.g., "governance warning", "journal pipeline")
- Naming conventions (e.g., "slug format", "DescKey", "ReadOrder")
- Acronyms and protocol terms (e.g., "MCP", "JSON-RPC 2.0")

### Format

Follow existing GLOSSARY.md format. Typical:
```markdown
**Term**: One-line definition.
```

### Guard Rails

- Only add terms the skill has high confidence about (actually
  encountered in code, not inferred)
- Max 10 new terms per run to avoid overwhelming the glossary
- Skip terms that are standard industry vocabulary (e.g., "CLI",
  "JSON", "YAML") — only project-specific terms
- Print added terms in convergence report: "Added N terms to
  GLOSSARY.md: [list]"

## Sync Strategy

Two options (decide during implementation):

**Option A: Copy on publish** — manual `ctx publish-docs` command
or Makefile target copies files. Simple, explicit, user-controlled.

**Option B: Build-time sync** — `make docs` copies from .context/
before building the site. Automatic, but adds build dependency.

Recommendation: **Option A** — matches ctx's philosophy of explicit
user control. Add a `make sync-docs` target that copies and strips
internal-only content.

## Acceptance Criteria

- [ ] Published docs render correctly in docs site
- [ ] Internal-only documents stay in .context/ only
- [ ] Sanitized docs use positive framing (not gap analysis)
- [ ] GLOSSARY.md update behavior added to ctx-architecture skill
- [ ] Sync mechanism chosen and documented
- [ ] Navigation updated (mkdocs.yml or equivalent)
