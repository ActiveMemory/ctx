# Skill Frontmatter Validity Test

Every embedded `SKILL.md` (across all tool surfaces — claude,
opencode, copilot-cli) carries YAML frontmatter that the host
runtime reads to dispatch the skill. Today nothing in CI
prevents a malformed or incomplete frontmatter from shipping.

## Problem

`internal/assets/embed_test.go` checks exactly one SKILL.md
(`ctx-history`) for the bare presence of a `---` delimiter
(see `TestSkillContent`, line 163). It does not iterate the
106 embedded SKILL.md files across the three tool trees, does
not parse the frontmatter, and does not assert any required
keys. The naming convention (`specs/future-complete/skill-naming-convention.md`,
implemented) is enforced socially, not mechanically — there's
no test that catches the next contributor who copies a SKILL.md
and forgets the `description:` field, mistypes `name:`, or lets
the directory name drift from the declared `name`.

The 106-file baseline today is internally consistent (manual
survey: zero violations), so a test added now functions as a
ratchet, not a remediation.

## Scope

A single Go test that, for every embedded `SKILL.md` under
`assets.FS` matching `**/skills/*/SKILL.md`:

1. Extracts the frontmatter (between the first two `---`
   lines).
2. Parses it as YAML.
3. Asserts:
   - `name:` is present, a non-empty string, and equals the
     containing directory's basename. (Drift between
     directory and declared name breaks discovery.)
   - `description:` is present and a non-empty string. (The
     host runtime uses this to decide whether to surface the
     skill.)
4. Reports every violation in a single pass (`t.Errorf`, not
   `t.Fatalf`) so the contributor sees all problems at once.

## Out of Scope

- Validating `allowed-tools` syntax (claude-only, optional,
  currently varied — `Bash(git:*)`, `Read`, etc. — and the
  Anthropic spec for that grammar belongs in a separate test).
- Cross-surface parity (does every claude skill have a
  copilot-cli counterpart?). Separate concern; would require
  a different fixture model.
- Enforcing the kebab-case-with-`ctx-`-prefix convention
  beyond directory-name match. The naming-convention spec
  already covers this; a dedicated test would belong with
  that spec, not this one.
- Reorganizing `internal/assets/read/skill/skill.go` to expose
  all three trees. The test walks `assets.FS` directly; no
  reader-API change.

## Location

`internal/assets/read/skill/frontmatter_test.go`. The
`internal/assets/read/skill/` package is the natural home —
its job is reading SKILL.md content — and the existing
`embed_test.go` is already busy with broader embed concerns.
