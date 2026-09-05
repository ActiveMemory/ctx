# Learnings

<!-- INDEX:START -->
| Date | Learning |
|----|--------|
| 2026-08-23 | Codex trust and hook wiring facts verified against codex 0.148 |
| 2026-08-23 | hack scripts must survive macOS /bin/bash 3.2 and BSD grep |
| 2026-08-23 | make lint SA5011 false positives mean a corrupted golangci-lint cache |
| 2026-07-25 | Using the proprietary sibling repo as design evidence leaks its internals into tracked files |
| 2026-07-25 | Skill and doc examples of a serialized structure must round-trip through the real parser |
| 2026-07-25 | A guard derived from a capability accessor silently lifts when the accessor is extended |
| 2026-07-19 | The disclosure parser is a deliberately dumb line-scanner (skips <!-- --> comments, not code fences) |
| 2026-07-19 | Measurement gates surface a real bug in every disclosure milestone |
<!-- INDEX:END -->

<!--
UPDATE WHEN:
- Discover a gotcha, bug, or unexpected behavior
- Debugging reveals non-obvious root cause
- External dependency has quirks worth documenting
- "I wish I knew this earlier" moments
- Production incidents reveal gaps

DO NOT UPDATE FOR:
- Well-documented behavior (link to docs instead)
- Temporary workarounds (use TASKS.md for follow-up)
- Opinions without evidence
-->


## [2026-08-23-204632] Windows git over HTTPS fails with SEC_E_NO_CREDENTIALS unless the openssl backend is pinned

**Context**: During the pi integration, every `git fetch`/`ls-remote` over an https remote on this Windows machine failed with `schannel: AcquireCredentialsHandle failed: SEC_E_NO_CREDENTIALS (0x8009030e)` - git's default TLS (schannel) has no usable credential state here. `gh` (own TLS stack) worked fine, and SSH keys were rejected, so the failure looked like a credentials/remote problem.

**Lesson**: git's schannel backend and the Windows credential manager can be broken independently of the remote; the fix is a per-invocation backend override, not key/credential surgery.

**Application**: On this machine, run every git https remote operation as `git -c http.sslBackend=openssl <cmd>` (git 2.51 supports it). Diagnose with `gh api` first to rule out remote-side auth before touching keys.

---

## [2026-08-23-204632] zensical site rebuilds are non-reproducible off the canonical toolchain (entity-escape drift)

**Context**: Rebuilding site/ locally with the pinned zensical 0.0.51 (fresh venv, Python 3.10) produced 120 files of pure escape drift (`&#39;` vs literal `'` inside highlighted code blocks) against the committed artifacts; pinning pygments older did not help.

**Lesson**: The zensical pin only pins the generator, not its transitive deps (markdown/pygments/jinja resolution drifts), so committed site/ output is not reproducible on arbitrary machines; the drift is user-invisible (entities render identically) but is pure diff churn that flips back on the next canonical build.

**Application**: After a local `zensical build`, run `git diff --ignore-cr-at-eol --numstat site/` and inspect representative hunks BEFORE committing; if the delta is escape drift rather than real content, revert site/ (`git checkout -- site/`, remove new untracked pages) and defer the rebuild to the canonical build machine; record the deferral in the task notes.

---

## [2026-08-23-170949] Hook commands must survive four shells and hostile cwds; hosts punish pre-ctx aborts

**Context**: Adversarial audit of every ctx hook surface (Claude/Codex/Copilot manifests, 16 Copilot wrapper scripts, OpenCode plugin, trace hook, plugin-reload) after the Codex non-repo-cwd anchor bug: 20 confirmed defects in 7 classes.

**Lesson**: Recurring classes: (1) ${VAR:?} aborts have SHELL-DEPENDENT exit codes (127 bash, 1 zsh, 2 dash) and exit 2 means BLOCK to Claude Code — never use :? in hook commands; guard with [ -d ... ] || { echo remedy >&2; exit 1; }. (2) Hosts may run hooks from non-repo cwds (Codex: plugin cache) — anchor with git rev-parse ... || pwd, or the host's schema-native cwd field (Copilot: "cwd": "."). (3) set -euo pipefail + jq/grep in command substitutions aborts whole hooks on non-JSON stdin — append || true inside the substitution. (4) INPUT=$(cat) hangs on host-held pipes — bound reads or < /dev/null. (5) Dead wrapper scripts accumulate real bugs invisibly; if a manifest calls ctx directly, ship no scripts.

**Application**: When adding any hook surface, test the command matrix under sh/bash3.2/zsh/dash from repo root, a subdir, a non-repo dir, with ctx absent, and with stdin held open. See specs/hook-surface-robustness.md.

---

## [2026-08-23-162206] Codex marketplace resolution diverges between CLI and TUI; dual-manifest plugin roots close it

**Context**: With both .agents/plugins/marketplace.json and legacy .claude-plugin/marketplace.json at the same root (same marketplace name, same plugin name), codex 0.148/0.149 CLI 'plugin add' resolved the .agents one, but the 0.149 TUI /plugins browser re-materialized the cache from the LEGACY one — silently swapping the installed ctx plugin back to the Claude variant whose ${CLAUDE_PROJECT_DIR:?} hooks all exit 1 under Codex (seen live as 13 UserPromptSubmit + 4 PreToolUse hook failures).

**Lesson**: Codex plugin-root ingestion prefers .codex-plugin/plugin.json when a root carries both manifest dirs (proven with a scratch legacy-only marketplace + marker hook). So a dual-manifest plugin root — .codex-plugin/plugin.json with hooks: ./hooks/codex.json next to the Claude manifest — yields working Codex hooks regardless of which marketplace file any Codex code path resolves.

**Application**: internal/assets/claude is now dual-manifest (codex.json synced by hack/sync-codex-skills.sh, byte-parity guarded by TestClaudeRootDualManifest). If Codex hooks suddenly fail with exit 1 en masse, check the cache root for a missing .codex-plugin/ and restart the codex session after reinstalling.

---

## [2026-08-23-154756] Codex marketplace add falls back to the legacy .claude-plugin marketplace

**Context**: User ran 'codex plugin marketplace add ActiveMemory/ctx' before the Codex marketplace landed on main. Codex 0.148 cloned GitHub main, found no .agents/plugins/marketplace.json, and silently used the legacy-compatible .claude-plugin/marketplace.json — installing the CLAUDE plugin variant (CLAUDE_PROJECT_DIR-anchored hooks that cannot run under Codex) into ~/.codex/plugins/cache under the same name and version.

**Lesson**: Codex marketplace resolution: .agents/plugins/marketplace.json is preferred when both exist (verified with a local dir containing both), but a source revision lacking it silently falls back to .claude-plugin/marketplace.json. The wrong variant is detectable by .claude-plugin/ in the installed cache root. Also: ctx's plugin-enabled detection cannot distinguish variants, so ctx setup codex --write short-circuits even when the wrong variant is installed.

**Application**: Until the branch is merged, install the Codex plugin from a local checkout of feat/codex-integration. When debugging 'plugin installed but hooks error', check the cache root for .claude-plugin/. Documented in docs/home/codex.md troubleshooting.

---

## [2026-08-23-125635] Codex trust and hook wiring facts verified against codex 0.148

**Context**: Live-tested the ctx Codex integration with codex exec on Codex CLI 0.148.0.

**Lesson**: (1) Project .codex/hooks.json loads only when the project path is trusted in the REAL ~/.codex/config.toml; a -c 'projects."...".trust_level="trusted"' CLI override is ignored for trust. (2) SessionStart plain-text stdout is injected verbatim as a developer message. (3) Codex's code-mode unified exec matches hook matcher 'Bash', and the legacy {"decision":"block"} shape blocks it. (4) SessionEnd hooks fire on codex exec process exit and ctx journal import completes within the 3 s cap. (5) trust for a parent dir (/Users/x) does NOT extend to subdirectories.

**Application**: When debugging 'ctx hooks not firing in Codex', check project trust in ~/.codex/config.toml first; do not suggest -c trust overrides.

---

## [2026-08-23-125635] hack scripts must survive macOS /bin/bash 3.2 and BSD grep

**Context**: make audit failed on macOS with 'unexpected EOF while looking for matching quote' in hack/lint-docstrings.sh (shebang #!/bin/bash = macOS bash 3.2.57). Root cause: bash 3.2's $( ) re-parser treats an apostrophe inside a COMMENT (didn't) as an open quote. Separately, the script's grep -cP (PCRE) silently fails on BSD grep, turning fieldcount empty and emitting 59 MISSING_FIELDS false positives.

**Lesson**: Two portability traps in hack/*.sh: (1) no apostrophes in comments inside command substitutions (bash 3.2 chokes); (2) no grep -P (BSD grep lacks PCRE) — use grep -E with a literal tab via TAB=$(printf '\t') and [[:space:]]. CI on Linux hides both.

**Application**: When adding hack scripts, test with /bin/bash (not Homebrew bash) on macOS; prefer 'did not' over contractions in comments inside $( ); use grep -E with POSIX classes.

---

## [2026-08-23-125635] make lint SA5011 false positives mean a corrupted golangci-lint cache

**Context**: make lint failed with 6 staticcheck SA5011 'possible nil pointer dereference' findings in test files untouched by the branch (if x == nil { t.Fatal } followed by x.Field). The flagged file set VARIED between runs (serve/compat one run, bootstrap/init the next).

**Lesson**: Nondeterministic staticcheck SA5011 on the guarded nil-check pattern is a corrupted golangci-lint build cache, not real findings. 'golangci-lint cache clean && make lint' returned 0 issues.

**Application**: Before chasing staticcheck findings in files a branch never touched, check whether the finding set is stable across two runs; if it varies, clean the golangci-lint cache first.

---

## [2026-07-25-124457] Using the proprietary sibling repo as design evidence leaks its internals into tracked files

**Context**: While deciding the pd-m4 add-path shape, I read the sibling repo's convention file to settle the question, then quoted its guide text and attributed the decision to it in a tracked plan file. An unrelated build warning prompted the sweep that caught it.

**Lesson**: The leak vector is not NAMING the sibling (allowed) — it is citing its designs or guide text as evidence in tracked specs/plans. Both the quote and the attribution are leaks.

**Application**: Cite ctx's OWN files as design evidence and describe designs on their merits; never quote or attribute to the sibling in tracked files; git grep the staged tree before committing any spec authored while consulting it. Per specs/public-repo-hygiene.md.

---

## [2026-07-25-124457] Skill and doc examples of a serialized structure must round-trip through the real parser

**Context**: The ctx-digest skill documented digest-plan entries as JSON strings, but Assignment.Entries is []StagedEntry{timestamp,title}; the documented form failed json.Unmarshal. Latent since M3 — any agent following the skill verbatim failed at the apply step.

**Lesson**: A hand-authored example of a serialized structure drifts from the code silently; nothing tests the prose.

**Application**: Verify skill/doc JSON/YAML examples against the actual unmarshal target — build the binary and run the exact documented command — before shipping the skill.

---

## [2026-07-25-124457] A guard derived from a capability accessor silently lifts when the accessor is extended

**Context**: pd-m4 T04 made ThemeDir(convention) return true; Apply's refusal of the convention kind was DERIVED from ThemeDir returning false. That one-line 'vocabulary' change silently opened the destructive mover to CONVENTIONS.md before the convention parse/validate/dup-title paths existed. Caught by TestApply_Refusals.

**Lesson**: When a refusal is derived from a capability check (if !capable { refuse }) instead of an explicit guard, extending that capability lifts the refusal as an invisible side effect — and the two tasks look independent in the plan.

**Application**: Before extending a capability accessor, grep for guards that read it. When a multi-task arc will open a destructive path, add an EXPLICIT gate that the opening task removes deliberately, so it cannot open early.

---

## [2026-07-19-210439] The disclosure parser is a deliberately dumb line-scanner (skips <!-- --> comments, not code fences)

**Context**: M4 unifies conventions into the ## -section entry model, where ## collides far more than ## [. Fence detection was considered and rejected.

**Lesson**: Fence detection (nested/tilde/indented fences) is a rabbit hole in a parser that does destructive moves. Safety comes from byte-conservation (a mis-cut is a reversible mis-grouping, not loss) plus the human-gated plan review (a false entry shows up at inspect).

**Application**: Don't add fence-awareness to the disclosure scanner; a stray ## inside a fenced code block is fixed at authoring time (rewrite that one line), never by the tool.

---

## [2026-07-19-210439] Measurement gates surface a real bug in every disclosure milestone

**Context**: pd-m1's add-path measurement gate flushed the insert.AfterHeader tail-truncation bug; pd-m3's T17 real rollout surfaced the DECISIONS-template comment-parse bug (a ## [ inside <!-- --> read as staging, so Validate returned ErrStagingUnparsable once the fold emptied real staging).

**Lesson**: For the disclosure work, driving the REAL fold — not unit tests — is where bugs hide. The measurement gate is not ceremony.

**Application**: Treat pd-m4 T20 (drive the digest on the real CONVENTIONS.md) as the authoritative gate; expect it to find something; verify conservation + all invariants on the real file before T21's human-gated rollout.

---

---

## Themes

- go-idioms-and-structure — Go language & package layout: import cycles, constant/helper placement, sync.Once smells, toolchain/build-tag pitfalls, error sentinels, test isolation → [go-idioms-and-structure](learnings/go-idioms-and-structure.md)
- audit-lint-compliance — Mechanical enforcement: the internal/audit AST gauntlet, compliance tests as command style guide, cmd/ purity, docstring floors, gosec triggers, linter quirks → [audit-lint-compliance](learnings/audit-lint-compliance.md)
- cli-command-design — cobra/CLI surface conventions: bare-invocation existence probing, legacyArgs silent-success on groups, commands named after input not output → [cli-command-design](learnings/cli-command-design.md)
- hooks-and-integration — Hook & editor-integration mechanics: output channels, key names, compliance wiring, Cursor/OpenCode plugins, project-local hooks, notify webhooks → [hooks-and-integration](learnings/hooks-and-integration.md)
- context-model-and-state — Context-dir resolution & on-disk state: single-source ContextDir, tombstones/logs/fs hygiene, managed-block guards, handover metadata, hub raft-lite → [context-model-and-state](learnings/context-model-and-state.md)
- text-markdown-serialization — Text/markdown/JSON hazards: line-sweep corruption, insert-helper tail drops, diacritic stripping, title-case brands, RFC3339 sort, key-dropping round-trips → [text-markdown-serialization](learnings/text-markdown-serialization.md)
- model-context-window — LLM model->context-window mapping: silent 200k fallback for new families, ordered prefix matching, doctor token_budget vs context_window → [model-context-window](learnings/model-context-window.md)
- git-and-signing — git CLI wrapping quirks: filter-branch leftover refs, hunk-level feature carving, GPG signing from non-TTY (pinentry) → [git-and-signing](learnings/git-and-signing.md)
- environment-and-platform — Cross-platform gotchas: macOS minimal PATH & /var symlink, Tauri rustc floor, GitNexus in-container pruning, host-pressure alerting → [environment-and-platform](learnings/environment-and-platform.md)
- skills-agents-and-tasks — Skill/agent workflow: skill lifecycle & shipping, agent context-loading/routing/behavior, task exit-criteria, ctx-dream design, refactor mechanics → [skills-agents-and-tasks](learnings/skills-agents-and-tasks.md)
- docs-and-templates — Docs/template/asset drift: tracked build output, magic-string discipline, contributor-PR reintroduction, docs/code divergence, redundant scripts → [docs-and-templates](learnings/docs-and-templates.md)
- editorial-and-product-signals — KB editorial & product signals: single topic-enumeration site, decision-recording, creator-confusion as doc-quality signal, editorial leverage, IDE-is-the-UI → [editorial-and-product-signals](learnings/editorial-and-product-signals.md)
