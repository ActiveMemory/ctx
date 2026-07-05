# Spec: ctx-humanize skill

## Problem

Prose that ships from AI-assisted sessions (blog posts, docs,
READMEs, PR descriptions, announcements) carries recognizable
LLM writing tells: significance inflation, brochure language,
forced triplets, em-dash density, chatbot residue. Readers
increasingly discount text that pattern-matches to machine
output. The project has a typography detector
(`hack/detect-ai-typography.sh`) but no editorial pass that
fixes the writing itself.

The first draft of this skill was a single 1,024-line
SKILL.md. That size violates the skill-authoring guidance this
repo ships (`ctx-skill-create`: keep SKILL.md under 500 lines,
move detail to `references/`), and every invocation would pay
the full catalog's context cost even for a quick review
verdict.

## Design

Ship `ctx-humanize` with the following structural and
editorial decisions:

1. **Progressive disclosure.** SKILL.md (~200 lines) carries
   the operating rules: invariants, modes, protected content,
   process, detection guidance, and a one-line-per-pattern
   index. The full 28-pattern catalog with before/after
   examples lives in `references/pattern-catalog.md`, read
   when actually rewriting.
2. **Deterministic typography verification.** The em/en dash
   ban is checked with Grep before returning, not by
   eyeballing. In repos that have
   `hack/detect-ai-typography.sh`, prefer it (noting it takes
   a directory of markdown, not a file).
3. **Voice-injection guardrail.** An earlier draft's
   voice-guidance example invented first-person emotion not
   present in the source, contradicting its own "do not add
   opinions" invariant. The shipped version constrains voice
   additions to stances traceable to the author's material.
4. **Agent-surface exclusions.** In ctx projects, `.context/`
   knowledge files, `specs/`, and `SKILL.md` files are
   agent-facing operational text, not humanizer targets,
   unless the user explicitly says otherwise.

Refined after a first live test, which surfaced edge-case
rulings now encoded in the skill: zero-content filler does not
count toward the coverage invariant; removing a fake
attribution must not silently strengthen the claim; the
review-only mode reports what the rewrite steps would have
caught.

Ships as a plugin skill (`internal/assets/claude/skills/`), so
`ctx init` deploys it; instructions stay project-agnostic with
ctx-repo specifics marked as conditional.

## Non-Goals

- Style enforcement for code, comments, commit messages, or
  structured data
- A general copy-editing / grammar skill (this targets AI
  tells specifically)
- Automatic invocation from hooks (manual, user-triggered)

## Acceptance

- [ ] SKILL.md under 500 lines; catalog in `references/` with
      a table of contents
- [ ] Review mode returns findings without touching files
- [ ] Rewrite preserves meaning, certainty, and claim strength;
      invents no facts
- [ ] Final output passes an em/en dash scan
- [ ] Documented in `docs/reference/skills.md` (table row +
      Content Creation section)
