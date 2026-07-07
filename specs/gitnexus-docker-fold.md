# Spec: vendor GitNexus Docker tooling; drop `npx gitnexus`

**Status:** accepted (impl 2026-07-05/06)

## Problem

GitNexus is the companion code-intelligence tool ctx recommends for
refactoring, impact analysis, and architecture enrichment. Two
operational facts made the previous `npx gitnexus ...` invocation
untenable on this box:

1. **`npx gitnexus` does not run.** GitNexus pins
   `tree-sitter@0.21.1`, a native addon that fails to build against
   Node 24 (a removed V8 API). Every `npx gitnexus` call errors out,
   so the docs and skills that told agents to run `npx gitnexus
   analyze` pointed at a broken command.
2. **No reproducible, host-safe way to run it.** Downgrading the
   host Node or hand-building native deps is invasive and per-machine.

The official GitNexus arm64 Docker image bakes Node 22 plus prebuilt
native deps, so it Just Works and stays off the host (no Node
downgrade, no extra host attack surface).

## Design

Vendor a thin Docker wrapper and wire it into the Makefile and docs:

- `hack/gitnexus-docker.sh` — `index <repo>` (analyze + register at
  the repo's real host path), `mcp` (air-gapped stdio MCP server that
  auto-mounts every registered repo), and passthrough to arbitrary
  `gitnexus <subcommand>`. Network is enabled only for the one-time
  `index` model fetch; the server and passthrough run `--network none`.
- `hack/gitnexus-index.sh` — convenience wrapper to index the current
  repo.
- `hack/strip-gitnexus.sh` — remove GitNexus-injected marker blocks
  from tracked files (keeps `.gitnexus/` out of the tree).
- `make gitnexus-index` / `make gitnexus-mcp` / `make strip-gitnexus`
  targets.
- Docs and skills updated to call `gitnexus ...` (via the wrapper) and
  the Docker flow, replacing every `npx gitnexus` reference
  (`GITNEXUS.md`, `CLAUDE.md`, `AGENTS.md`, `docs/home/getting-started.md`,
  `docs/recipes/multi-tool-setup.md`, and the `ctx-remember` /
  `ctx-architecture-enrich` skills that mention GitNexus).

## Data-loss guard (registry pruning)

GitNexus **prunes registry entries whose repo paths it cannot stat
in-container** — observed live with both `index` and `list`. The
wrapper therefore mounts *every* registered repo at its real path on
all registry-touching branches (a `registry_mounts` helper reads
`registry.json`). The helper fails loud rather than fail open: on a
registry-touching branch, a missing `python3` or an unparseable
non-empty registry aborts the run instead of silently mounting an
empty set (which would let GitNexus prune every entry). A legitimately
empty registry (fresh setup) proceeds with no mounts.

See LEARNINGS: "GitNexus prunes registry entries whose repo paths
don't resolve in-container." The sibling `os`/`orchestrator` copies of
these scripts still carry the older fail-open shape and must be
backported before use.
