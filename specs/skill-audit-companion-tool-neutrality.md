# Skill Audit: Companion-Tool Neutrality in Body Text

Rewrite the prescriptive references to GitNexus and Gemini
Search in the embedded skill bodies as capability-first
descriptions with the canonical tools listed as examples.
Honors the manifesto's "ctx is unopinionated about the
agent's toolchain" stance without contradicting the
2026-05-23 anti-MCP-gateway decision.

## Problem

Eight skill files name GitNexus and Gemini Search directly:

- `claude/skills/ctx-refactor/SKILL.md`
- `claude/skills/ctx-explain/SKILL.md`
- `claude/skills/ctx-code-review/SKILL.md`
- `claude/skills/ctx-remember/SKILL.md`
- `claude/skills/ctx-architecture/SKILL.md`
- `claude/skills/ctx-architecture-enrich/SKILL.md`
- `claude/skills/ctx-architecture-failure-analysis/SKILL.md`
- `integrations/copilot-cli/skills/ctx-remember/SKILL.md`

The references span three shapes:

1. **External slash-command prescriptions** (e.g.,
   `/gitnexus-refactoring`) — assume the user has the
   specific upstream skill suite installed.
2. **Tool-name prescriptions in instructions** (e.g.,
   "Use Gemini throughout the skill", "GitNexus MCP:
   blast radius estimation") — assume one specific MCP
   server is the chosen implementation of a capability.
3. **`allowed-tools` frontmatter listing specific MCP
   server names** (`mcp__gitnexus__*`, `mcp__gemini-search__*`)
   — Claude Code permission system; scopes which tools
   the skill can invoke.

Shapes 1 and 2 are pure text and prescribe specific
tools. A user with Firecrawl + sourcegraph-cody (or vLLM,
or Exa, or Tavily — pick your stack) reads the skill and
sees instructions naming tools they don't have. The
skill's capability-level intent is identical regardless
of which underlying MCP server provides it; only the
text is opinionated.

Shape 3 is a real permission boundary and a different
problem class. Out of scope for this pass (see Out of
Scope below).

## Solution

Capability-first language with canonical tools as
examples. The pattern:

| Prescriptive (before) | Capability-first (after) |
|---|---|
| "Use Gemini Search to look up …" | "Use a web-search-with-citations MCP if available (Gemini Search is the typical choice; equivalents include Firecrawl, Exa, Tavily) …" |
| "GitNexus MCP: blast radius estimation" | "A code-intelligence MCP (GitNexus is the canonical choice; equivalents include sourcegraph-cody) provides blast-radius estimation …" |
| "Use `/gitnexus-refactoring` if available" | "If you have an external refactoring-aware skill (e.g., the GitNexus suite ships `/gitnexus-refactoring`), invoke it; otherwise proceed with built-in reasoning." |
| "Install via MCP settings if needed" | (removed — covered by the companion-fallback fix that preceded this commit) |

Per-file rewrite notes:

### ctx-refactor (1 line)

Replace `/gitnexus-refactoring` slash command with an
"if you have an external refactoring-aware skill, invoke
it" phrasing that names the GitNexus suite as the
canonical example.

### ctx-explain (2 lines)

Replace `/gitnexus-debugging` and `/gitnexus-exploring`
references with the same pattern.

### ctx-code-review (1 line)

Replace `/gitnexus-pr-review` reference with the same
pattern.

### ctx-remember (tool table + companion check)

Replace the prescriptive Gemini/GitNexus tool table with
a capability-first version. The Companion Tool Check
section already had its install-nag removed in the
preceding commit; this pass updates the table headers
above it.

### ctx-architecture

Rewrite the "Check if **Gemini Search** MCP is available"
section to "Check if a web-search-with-citations MCP is
available (Gemini Search is the typical choice)". Same
treatment for the principal-mode section that prescribes
Gemini specifically.

### ctx-architecture-enrich

Rewrite the GitNexus-required preamble and the install
procedure as "a code-intelligence MCP is required for
this skill (GitNexus is the typical choice; equivalents
work). If absent, the skill cannot run — its purpose IS
the code-graph verification pass." Keeps the
required-tool stance (the skill's purpose is unchanged)
while removing the GitNexus-specific install command.

### ctx-architecture-failure-analysis

Rewrite the "GitNexus MCP: blast radius" / "Gemini
Search: cross-reference" lines to name the capabilities,
list the canonical tools as examples, and let the agent
self-route based on what's connected.

## Doc updates (follow-on within scope)

After the SKILL.md rewrites landed, a check across `docs/`
surfaced six files that also named GitNexus and Gemini
Search prescriptively. Two categories:

**Operational / descriptive docs** (rewritten to
capability-first matching the SKILL.md tone):

- `docs/operations/runbooks/architecture-exploration.md` —
  "via GitNexus" → "via a code-intelligence MCP (canonical:
  GitNexus)"; preflight smoke-test wording; graceful-fail
  rationale.
- `docs/recipes/architecture-deep-dive.md` — "graph-backed
  data from GitNexus"; "Requires: GitNexus MCP server".
- `docs/reference/skills.md` — failure-analysis skill
  description.
- `docs/cli/index.md` — `companion_check` field
  parenthetical.

**Install-guide docs** (kept concrete; added one-liner
naming equivalents):

- `docs/home/getting-started.md` — Section 7's companion
  tool sublist gains a "canonical examples; equivalents
  work" preamble. The install commands stay because the
  doc's purpose IS install guidance.
- `docs/recipes/multi-tool-setup.md` — Companion Tools
  section preamble similarly clarifies that named tools
  are canonical examples. Individual tool sections still
  give concrete setup steps for the canonical impls.

This split is intentional: when a doc's job is "tell me
how to install something," that doc names specific tools.
When a doc's job is "describe what a skill does," it
describes capabilities. The decision recorded in
`DECISIONS.md` (2026-05-23 "Skill body text uses
capability-first language; install-guide docs name
canonical implementations") captures the rule for
future contributors.

## Out of Scope

- **`allowed-tools` frontmatter genericization.** The
  frontmatter at the top of each affected skill lists
  specific MCP server name patterns (`mcp__gitnexus__*`,
  `mcp__gemini-search__*`) to scope Claude Code's
  permission grant. Replacing with `mcp__*` would grant
  the skill access to *every* MCP server the user has —
  a permission expansion, not just a cosmetic change.
  Operators who use a different toolchain (sourcegraph,
  vLLM, Firecrawl, etc.) need to edit `allowed-tools`
  in their local skill copy or fork. A separate spec can
  consider whether to template the allowlist; this pass
  is body-text only.

- **Tested-companions catalog in docs.** Option C from
  the design discussion (a docs page listing known-good
  third-party MCP tools by capability) is a follow-up
  and out of this pass. Skills should be unopinionated
  first; the catalog informs without prescribing later.

- **Validating that named alternatives actually work.**
  Mentioning Firecrawl, sourcegraph-cody, Exa, etc. as
  examples is informational only — ctx does not guarantee
  these alternatives have parity with the canonical tool.
  Operators evaluate fit.

## Verification

- `grep -rE "use.* GitNexus|use.* Gemini" internal/assets/`
  returns no prescriptive matches (instances should all be
  "if available", "typical choice", or similar
  capability-first phrasing).
- The 8 SKILL.md files compile through the existing
  embed/audit tests unchanged (no metadata mutated).
- `make lint` clean; no Go code touched.
- `go test ./...` clean — no test depends on the specific
  phrasing being audited.

## Migration consequence

Existing users with GitNexus + Gemini Search configured
see no behavioral change — the canonical tools are still
the first-listed examples in every section. Users with
different toolchains read the skill text and understand
the capability requirement rather than the tool
prescription. The agent self-routes based on what's
connected.

For the architecture-enrich skill specifically, the
"required code-intelligence MCP" stance is unchanged —
the skill's entire purpose is code-graph enrichment, so
running it without ANY graph MCP makes no sense. The
text just names the capability instead of the specific
tool.
