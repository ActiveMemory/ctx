# integrations-and-assets

## [2026-06-07-180005] YAML text externalization, init, and drift guards (consolidated)

**Consolidated from**: 5 entries (2026-03-13 to 2026-04-03)

- All user-facing text externalizes to embedded YAML domain files
  (commands/flags/text/examples split via dedicated loaders), justified by agent
  legibility (named DescKey constants as traversable graphs) and drift
  prevention, not i18n; the 3-file ceremony (DescKey + YAML + write/err fn) is
  the accepted cost.
- Static embedded data and resource lookups use an explicit Init() called
  eagerly at startup, not per-accessor sync.Once or package-level init() —
  makes the startup dependency visible and testable; maps unexported, accessors
  are plain lookups.
- A Go↔YAML linkage check (lint-drift check 5, shell grep+comm) catches
  orphaned/broken DescKey↔YAML links and cross-namespace duplicates at CI
  time, preventing silent runtime failures.
- The build target depends on sync-why so derived assets/why/ files cannot drift
  from their docs/ sources — build fails without sync.
- MCP resource name constants live in config/mcp/resource (parallel to
  config/mcp/tool); the resource→file mapping stays in server/resource (too
  many cross-cutting deps for a config package), pre-built once at server init
  for O(1) lookup.

---

## [2026-06-07-180008] ctxctl maintainer binary and out-of-band audit channel (consolidated)

**Consolidated from**: 4 entries (2026-05-24 to 2026-05-28)

- Discipline enforcement belongs on the verbatim-relay channel, run out-of-band:
  relay is the one discipline channel that survives tunnel vision; run the
  auditor in a separate Claude Code session for fresh-context judgment and cost
  control. New generic channel: a skill writes .context/audit/<kind>.md, a
  check-audit hook relays unread reports verbatim, ctx audit list/show/dismiss
  manages lifecycle (digest-bound dismissal).
- [Superseded] ctxctl first placed at cmd/ctxctl in the same Go module:
  binary-level isolation via transitive-import exclusion, zero relocation of
  existing internal/audit files, on the belief a separate go.mod couldn't import
  the parent's internal/.
- That belief was empirically disproved: a nested module lexically under the
  parent path CAN import internal/. So ctxctl became a separate Go module at
  tools/ctxctl (own go.mod) — a hard module boundary guarantees ctx can never
  import ctxctl (the asymmetric requirement that matters); one-directional
  ctxctl→ctx coupling is acceptable for disposable maintainer tooling. A
  go.work wires the workspace; a guard test asserts cmd/ctx never imports
  internal/ctxctl.
- ctxctl is PATH-installed alongside ctx (build to dist/, install to
  /usr/local/bin/ctxctl) for clean repo roots and one binary across all
  worktrees, mirroring ctx's install pattern; the local hook calls ctxctl from
  PATH.

---

## [2026-06-07-180010] Companion-tool integration: peer-MCP, no gateway (consolidated)

**Consolidated from**: 6 entries (2026-03-06 to 2026-05-23)

- Peer MCP model for external tools (GitNexus, context-mode): side-by-side
  servers each queried independently by the agent, chosen over orchestrator/hub
  models to respect ctx's markdown-on-filesystem invariant and avoid
  coupling/plugin registries.
- Skills stay CLI-based; MCP Prompts are the protocol equivalent: CLI is always
  available (PATH prereq), MCP is optional config, hooks are always CLI — two
  access patterns in one tool is gratuitous complexity.
- Recommend companion RAGs as peer MCP servers, not bridged through ctx: MCP is
  the composition layer; ctx is context, RAGs are intelligence — no bridging,
  plugin system, or schema abstraction.
- Companion tools documented as optional MCP enhancements with a runtime check
  (/ctx-remember smoke-tests MCPs at session start; companion_check:false
  suppresses) so users learn what enhances their workflow without being forced
  to install.
- MCP gateway not worth the coupling cost: a gateway would make ctx own
  install/uninstall/version/error-surface for tools it doesn't ship
  (bidirectional ownership coupling); composition is already MCP's job and the
  skills already work peer-to-peer. The pluggable-graph-tool task was skipped as
  a direct consequence (pluggability without ownership is incoherent).
- Skill body text uses capability-first language with canonical tools as
  examples; install-guide docs name canonical implementations directly
  (newcomers need a recommendation); allowed-tools frontmatter stays
  MCP-specific (genericizing to mcp__* is a permission expansion). Pure text
  rewrite, no new abstraction layer.

---

## [2026-06-07-180012] Embedded assets and editor-integration harnesses (consolidated)

**Consolidated from**: 7 entries (2026-04-01 to 2026-05-22)

- Embedded foreign-language assets (TS/Bash/PowerShell/YAML) under
  internal/assets/ are intentional, not a smell: every file is //go:embed'd into
  the ctx binary and written at ctx setup; internal/ is about import privacy,
  not source language. The fix for the legibility gap was a contract README, not
  relocation (//go:embed can't reference ../).
- assets/hooks/ split into assets/integrations/ (tool-integration assets:
  Copilot instructions, AGENTS.md, CLI scripts/skills) + assets/hooks/messages/
  (hook-system templates) — integration assets are not hooks.
- Embedded harnesses (//go:embed'd, shipped via ctx setup) and
  separately-published harnesses (e.g. VS Code extension → marketplace, own
  cadence) are first-class peers with distinct CI/release pipelines; a new
  harness declares which pattern it follows before placing files.
- OpenCode plugin ships without a tool.execute.before hook: the natural fit
  (block-dangerous-commands) isn't a ctx Go subcommand and shimming would brick
  the editor (Cobra exit-1 read as {blocked:true}) on installs without the
  Claude wrapper. This omission is permanent — block-dangerous-commands will
  not be promoted to a ctx Go subcommand; the perpetually-pending re-add task is
  closed.
- Under cwd-anchored, the OpenCode plugin's agent shell tool can't be anchored
  to project root (the @opencode-ai/plugin SDK exposes only env, not cwd on
  shell.env); drop the shell.env handler and document launch-from-root.
  Plugin-internal ceremony calls stay anchored; the cwd-anchored error message
  is self-fixing.
- Editor-integration plugins must filter post-commit to actual git commit
  invocations (regex on the extracted command), not fire on every shell call —
  firing on noise trains users to ignore nudges.

---

## [2026-05-31-094649] ctx Desktop shells out via std::process::Command, not tauri-plugin-shell

**Status**: Accepted

**Context**: The GUI must run the ctx binary for every read and write. Tauri 2 offers tauri-plugin-shell with capability-scoped command allowlists.

**Decision**: ctx Desktop shells out via std::process::Command, not tauri-plugin-shell

**Rationale**: Running ctx inside our own #[tauri::command] via std::process::Command avoids all shell-plugin permission/capability wiring and keeps the adapter a single Rust module. ctx resolves its context from $PWD/.context, so each call sets current_dir to the selected project root.

**Consequence**: No shell capability in capabilities/default.json; the adapter owns PATH augmentation and git provenance synthesis; a CLI/output change is a one-file fix in ctx_adapter.rs.

---

