# Glossary Schema

Shape for `.context/kb/glossary.md`. Canonical terms scoped to
this kb's `## Scope` statement. Each entry is a declarative
definition with a confidence band and at least one
`evidence-index.md` row backing it.

Glossary confidence reflects definition quality ("do we agree
on what this term means?"), independent of whether the cited
evidence rows themselves carry `high` or `low` bands.

## Fields

| Field           | Description                                                  |
|-----------------|--------------------------------------------------------------|
| `term`          | Canonical term as a level-2 heading; lowercase preferred.    |
| `aliases`       | Optional comma-separated list of synonyms or short forms.    |
| `definition`    | One paragraph; declarative; no hedging beyond the band.      |
| `confidence`    | One of `high`, `medium`, `low`, `speculative`.               |
| `evidence`      | Comma-separated `EV-###` references; at least one required.  |
| `related-terms` | Optional links to other glossary entries in this file.       |

## Example

```markdown
## widget

**Aliases:** widget-the-primitive, the unit.

**Definition.** A widget is the smallest packaging boundary in
the system: a folder containing `SKILL.md` plus optional
`scripts/`, `references/`, and `assets/` subdirectories.

**Confidence:** high.

**Evidence:** EV-042, EV-043.

**See also:** [widget contract](#widget-contract).
```

## Rules

- Do not define a term outside the kb's declared `## Scope` in
  `.context/kb/index.md`; out-of-scope terms belong elsewhere.
- Every entry cites at least one `EV-###`; ungrounded terms
  open an `outstanding-questions.md` row instead.
- Aliases are additive, not renamings; the canonical term is
  the heading and never changes once published.
- Glossary confidence tracks definitional consensus, not the
  bands of the cited evidence rows; do not mechanically
  downgrade a definition just because its evidence was demoted.
- Cross-reference related terms within the file; do not link
  to topic pages from the definition body (use the topic page's
  `## Related concepts in this kb` block for that direction).
