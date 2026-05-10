# Spec: Architecture Refresh (April 2026)

## Goal

Full architecture re-analysis after the v0.8.0 restructuring
(527 packages, 34 commands, MCP server, cmd/core taxonomy).
Principal mode with agent-agnosticism as strategic lens.
GitNexus enrichment for verified blast radius data.

## Scope

- Rewrite ARCHITECTURE.md (layered architecture, dual entry points)
- Split DETAILED_DESIGN.md into 5 domain files
- Write CHEAT-SHEETS.md (6 lifecycle flows)
- Write ARCHITECTURE-PRINCIPAL.md (strategic analysis)
- Write DANGER-ZONES.md (23 zones, 4 CRITICAL after enrichment)
- Write CONVERGENCE-REPORT.md with search prompts
- Write EXTENSION-POINTS.md (14 patterns via GitNexus)
- Update GLOSSARY.md (9 new terms)
- Write READMEs for config/, mcp/, cli/system/
- Update config/doc.go

## Acceptance Criteria

- [x] All architecture artifacts regenerated
- [x] Principal analysis with agent-agnosticism lens
- [x] GitNexus enrichment with verified blast radius
- [x] Danger zones enriched (desc.Text 53 flows, SafeWriteFile 69 callers)
- [x] Clustering comparison (94 auto vs 5 manual domains)
- [x] READMEs for high-navigation packages
- [x] Glossary updated with architecture terms
- [x] Decisions and learnings persisted
