# Convergence Report

_Generated 2026-04-03 by /ctx-architecture principal_

## By Module

| Module | Confidence | Status | Blocker |
|--------|------------|--------|---------|
| internal/config/* | 0.85 | 🟡 Solid | Pattern understood; not all 60 sub-packages individually read |
| internal/assets | 0.85 | 🟡 Solid | Embed and read/* pattern understood; tpl/ migration pending |
| internal/io | 0.80 | 🟡 Solid | API understood; edge cases not all traced |
| internal/format | 0.80 | 🟡 Solid | All functions cataloged |
| internal/parse | 0.80 | 🟡 Solid | Small package, fully understood |
| internal/sanitize | 0.80 | 🟡 Solid | Small package, fully understood |
| internal/validate | 0.75 | 🟡 Solid | Enriched 2026-06-09: Symlinks() is a single-caller chokepoint on the context-load path (9 flows) |
| internal/inspect | 0.75 | 🟡 Solid | API cataloged from survey |
| internal/flagbind | 0.85 | 🟡 Solid | Pattern and all variants documented |
| internal/exec/* | 0.80 | 🟡 Solid | All 5 wrappers surveyed |
| internal/log/* | 0.80 | 🟡 Solid | Event + warn split understood |
| internal/crypto | 0.90 | ✅ Converged | Converged |
| internal/sysinfo | 0.85 | 🟡 Solid | Platform build tags understood |
| internal/rc | 0.90 | ✅ Converged | Converged |
| internal/entity | 0.85 | 🟡 Solid | All types cataloged; some methods not fully traced |
| internal/entry | 0.85 | 🟡 Solid | Flow understood; not all callers traced |
| internal/context/* | 0.80 | 🟡 Solid | 6 sub-packages surveyed |
| internal/drift | 0.85 | 🟡 Solid | 7 checks documented |
| internal/index | 0.85 | 🟡 Solid | Converged |
| internal/task | 0.85 | 🟡 Solid | Converged |
| internal/tidy | 0.80 | 🟡 Solid | Block parsing understood |
| internal/trace | 0.80 | 🟡 Solid | Flow documented; resolve helpers not all read |
| internal/journal/parser | 0.80 | 🟡 Solid | 4 parsers documented; Copilot parsers shallow |
| internal/journal/state | 0.85 | 🟡 Solid | Pipeline stages documented |
| internal/memory | 0.80 | 🟡 Solid | Discovery + sync flow understood |
| internal/notify | 0.85 | 🟡 Solid | Converged |
| internal/claude | 0.80 | 🟡 Solid | Thin wrapper understood |
| internal/mcp/* | 0.90 | ✅ Converged | All 15 sub-packages deeply read |
| internal/cli/* (34 cmds) | 0.75 | 🟡 Solid | Taxonomy understood; not all core/ packages deeply read |
| internal/bootstrap | 0.85 | 🟡 Solid | Registration and guards documented |
| internal/write/* | 0.80 | 🟡 Solid | Pattern understood; not all 46 packages individually read |
| internal/err/* | 0.80 | 🟡 Solid | Pattern understood; not all 35 packages individually read |
| internal/audit | 0.75 | 🟡 Solid | Purpose understood; individual test files not read |
| internal/compliance | 0.75 | 🟡 Solid | Enriched 2026-06-09: 29 checks cataloged across 4 test files; section added to DETAILED_DESIGN-cli.md |

## By Domain

| Domain | Modules | Converged | Avg Confidence |
|--------|---------|-----------|----------------|
| Foundation (config, assets, infra) | 12 | 2/12 | 0.82 |
| Domain (entity through claude) | 13 | 0/13 | 0.82 |
| MCP | 1 (15 sub-pkgs) | 1/1 | 0.90 |
| CLI | 2 (34 commands) | 0/2 | 0.78 |
| Output (write, err) | 2 | 0/2 | 0.80 |
| Quality (audit, compliance) | 2 | 0/2 | 0.75 |

## Overall

- Total module groups: 32
- Converged (>= 0.9): 3  ✅ (crypto, rc, mcp)
- Solid (0.7-0.89): 28  🟡
- Shallow (0.4-0.69): 0  🔶
- Stubbed (< 0.4): 0  🔴

## What Would Help Next

🟡 internal/cli/* (0.75) - Solid
  → Deep-read core/ packages for journal, agent, system, add
    (the largest/most complex commands)
  → Read test files to understand edge case behavior
  → Ask: "walk me through what happens when ctx journal site runs"

🟡 internal/audit (0.75) - Solid
  → Read individual audit test files to catalog all 40+ checks
  → Understand quarantine/allowlist management patterns

🟡 internal/compliance (0.70) - Solid
  → Read compliance test files to understand file-level checks
  → Understand relationship to Makefile lint targets

🟡 internal/journal/parser (0.80) - Solid
  → Deep-read Copilot and Copilot CLI parsers for parity assessment
  → Read query.go for session lookup capabilities

🟡 internal/trace (0.80) - Solid
  → Read all resolve_* helper files for ref resolution logic
  → Trace full commit lifecycle from collect to trailer

## Convergence Verdict

🟡 **MOSTLY CONVERGED** - All modules at 0.70+. Core modules (MCP,
config pattern, domain flow) well understood. Diminishing returns on
full re-run; use focus areas for journal, system hooks, and audit
tests to reach full convergence.

---

## Search Prompts

The right keyword changes everything. Based on what I found in
the codebase, here are targeted searches worth running — in your
internal docs, Confluence, Notion, Slack, or publicly:

### Fill the gaps (ranked by how much they'd help)

🟡 internal/cli/system (0.75) - 34 hook subcommands
  Try searching:
  - "ctx system check-persistence" behavior or "adaptive prompt counter"
  - "Claude Code hook lifecycle" or "UserPromptSubmit hook protocol"
  - "ctx hook throttle" or "daily throttle marker"

🟡 internal/audit (0.75) - 40+ AST audit checks
  Try searching:
  - "TestNoDeadExports quarantine" or "grandfathered map audit"
  - "Go AST audit pattern" or "go/packages test convention"

### Concepts worth understanding deeply

- "MCP protocol 2024-11-05 spec" — the protocol ctx implements;
  understanding the full spec reveals capability gaps
- "Claude Code JSONL session format" — reverse-engineered; the
  actual schema would replace guesswork in journal/parser
- "Copilot instructions.md specification" — understanding what
  Copilot expects would inform agent-agnostic setup
- "go:embed performance large binary" — ctx embeds skills, hooks,
  templates, YAML; understanding embed overhead at scale matters

### Architecture decision records

- "ctx context priority order rationale" or "FileReadOrder design"
  — why CONSTITUTION loads first is a design decision worth ADR
- "ctx YAML externalization decision" or "i18n readiness Go CLI"
  — the 879-key text externalization is a major architectural bet
- "ctx MCP vs CLI parity decision" — whether the capability gap
  between 34 CLI commands and 11 MCP tools is intentional

---

Note: I did not run these searches — you may have internal docs
where these are more useful than public results, and you know
which sources to trust.

---

## Enrichment Summary

_Last enrichment: 2026-06-09 via GitNexus (index @ 60d8e823,
27,927 symbols, 203 flows; prior enrichment 2026-04-03 @ bf42b1f6)_

| Phase | Items Processed | Key Findings |
|-------|----------------|--------------|
| Danger zones | 26 entries, 9 symbols re-verified | desc.Text d=1 now 940 (was "30+"); SafeWriteFile 95 (was 69); rc.RC new CRITICAL (d=1: 34, 8 modules); DiscoverPath downgraded CRITICAL->HIGH (d=1: 4, single domain); 3 undercount flags (load.Do, Serve, ValidateAndWrite) |
| Extension points | 14 patterns | CLI commands 34->42, MCP tools 11->15, drift checks 7->12, setup deployers 5->8 (cline, cursor, kiro, opencode); parsers/prompts unchanged (4/5) |
| Execution flows | 203 total, 11 indexed | 13 multi-flow hotspots; new top-tier: check.FullPreamble (40) and check.Preamble (39) — every hook routes through them |
| Clustering | 100 clusters vs 5 manual domains | Coverage gaps now dominate mismatches: Hub, Steering, Trigger, kb + ~8 command families have no manual baseline; Format cohesion collapsed 97%->26% |
| Shallow modules | 2 enriched | validate 0.70->0.75 (single-caller chokepoint), compliance 0.70->0.75 (29 checks cataloged) |

⚠ The codebase has outgrown the 2026-04 architecture baseline
(19.3k -> 27.9k symbols). Recommend `/ctx-architecture` refresh
before the next enrichment pass.
