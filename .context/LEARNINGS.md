# Learnings

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-02-27 | Context injection and compliance strategy (consolidated) |
| 2026-02-26 | Webhook silence after ctxrc profile swap is the most common notify debugging red herring |
| 2026-02-26 | Documentation drift and auditing (consolidated) |
| 2026-02-26 | Agent context loading and task routing (consolidated) |
| 2026-02-26 | Go testing patterns (consolidated) |
| 2026-02-26 | PATH and binary handling (consolidated) |
| 2026-02-26 | Task management and exit criteria (consolidated) |
| 2026-02-26 | Agent behavioral patterns (consolidated) |
| 2026-02-26 | Hook compliance and output routing (consolidated) |
| 2026-02-26 | ctx add and decision recording (consolidated) |
| 2026-02-24 | CLI tools don't benefit from in-memory caching of context files |
| 2026-02-22 | Hook behavior and patterns (consolidated) |
| 2026-02-22 | UserPromptSubmit hook output channels (consolidated) |
| 2026-02-22 | Linting and static analysis (consolidated) |
| 2026-02-22 | Permission and settings drift (consolidated) |
| 2026-02-22 | Gitignore and filesystem hygiene (consolidated) |
| 2026-02-19 | Feature can be code-complete but invisible to users |
| 2026-01-28 | IDE is already the UI |
<!-- INDEX:END -->

---

## [2026-02-27-002830] Context injection and compliance strategy (consolidated)

**Consolidated from**: 3 entries (2026-02-26)

- Verbal summaries with linked diagram files cut ARCHITECTURE.md from ~12K to ~3.8K tokens. Extract diagrams to linked files outside FileReadOrder; keep prose summaries inline. The 4-chars-per-token estimator is accurate — optimize content, not the estimator.
- Soft instructions have a ~75-85% compliance ceiling because "don't apply judgment" is itself evaluated by judgment. When 100% compliance is required, don't instruct — inject via `additionalContext`. Reserve soft instructions for ~80% acceptable compliance.
- Once ~7K tokens are auto-injected (fait accompli), the agent's rationalization inverts from "skip to save effort" to "marginal cost is trivial." Front-load highest-value content as injection, then use sunk cost to motivate on-demand reads for the remainder.

---

## [2026-02-26-003854] Webhook silence after ctxrc profile swap is the most common notify debugging red herring

**Context**: Spent time investigating why webhooks weren't firing — checked binary version, hook configs, notify.Send internals. Actual cause was .ctxrc swapped to prod profile (notify commented out) earlier in session.

**Lesson**: When webhooks stop, check .ctxrc profile first (hack/ctxrc-swap.sh status). Also: not all tool uses trigger webhook-sending hooks — Read only triggers context-load-gate (one-shot) and ctx agent (no webhook). qa-reminder requires Edit matcher.

**Application**: Before debugging notify internals, run hack/ctxrc-swap.sh status and verify the event would actually match a hook with notify.Send.

---

## [2026-02-26-100000] Documentation drift and auditing (consolidated)

**Consolidated from**: 6 entries (2026-01-29 to 2026-02-24)

- CLI reference docs can outpace implementation: ctx remind had no CLI, ctx recall sync had no Cobra wiring, key file naming diverged between docs and code. Always verify with `ctx <cmd> --help` before releasing docs.
- Structural doc sections (project layouts, command tables, skill counts) drift silently. Add `<!-- drift-check: <shell command> -->` markers above any section that mirrors codebase structure.
- Agent sweeps for style violations are unreliable (8 found vs 48+ actual). Always follow agent results with targeted grep and manual classification.
- ARCHITECTURE.md missed 4 core packages and 4 CLI commands. The /ctx-drift skill catches stale paths but not missing entries — run /ctx-map after adding new packages or commands.
- Documentation audits must compare against known-good examples and pattern-match for the COMPLETE standard, not just presence of any comment.
- Dead link checking belongs in /consolidate's check list (check 12), not as a standalone concern. When a new audit concern emerges, check if it fits an existing audit skill first.

---

## [2026-02-26-100002] Agent context loading and task routing (consolidated)

**Consolidated from**: 5 entries (2026-01-20 to 2026-01-25)

- `ctx agent` is optimized for task execution (filters pending tasks, surfaces constitution, token-budget aware). Manual file reading is better for exploratory/memory questions (session history, timestamps, completed tasks).
- On "Do you remember?" questions, immediately read .context/ files and run `ctx recall list --limit 5`. Never ask "would you like me to check?" — that is the obvious intent.
- .context/ is NOT a Claude Code primitive. Only CLAUDE.md and .claude/settings.json are auto-loaded. The .context/ directory requires a hook or explicit CLAUDE.md instruction to be discovered.
- Orchestrator (IMPLEMENTATION_PLAN.md) and agent (.context/TASKS.md) task lists must be separate. The orchestrator says "check your mind" — it doesn't maintain a parallel ledger.
- Only CLAUDE.md is auto-loaded by Claude Code. Projects using ctx should rely on the CLAUDE.md -> AGENT_PLAYBOOK.md chain, not AGENTS.md.

---

## [2026-02-26-100005] Go testing patterns (consolidated)

**Consolidated from**: 7 entries (2026-01-19 to 2026-02-26)

- Compiler-driven refactoring misses test files: `go build ./...` catches production callsite breaks but not test files. Always run `go test ./...` after signature changes.
- All runCmd() returns must be consumed in tests: even setup calls need `_, _ = runCmd(...)` to satisfy errcheck.
- Set `color.NoColor = true` in a package-level init function to disable ANSI codes for CLI test string assertions.
- Recall CLI tests isolate via HOME env var: `t.Setenv("HOME", tmpDir)` with `.claude/projects/` structure gives full isolation from real session data.
- `formatDuration` accepts an interface with a Minutes method, not time.Duration directly. Use a stubDuration struct for testing.
- CI tests need `CTX_SKIP_PATH_CHECK=1` env var because init checks if ctx is in PATH.
- CGO must be disabled for ARM64 Linux (`CGO_ENABLED=0`) — CGO causes cross-compilation issues with `-m64` flag.

---

## [2026-02-26-100006] PATH and binary handling (consolidated)

**Consolidated from**: 3 entries (2026-01-21 to 2026-02-17)

- Always use `ctx` from PATH, never `./dist/ctx-linux-arm64` or `go run ./cmd/ctx`. Check `which ctx` if unsure.
- Hooks must use PATH, not hardcoded paths. `ctx init` checks if ctx is in PATH before proceeding. Tests can skip with `CTX_SKIP_PATH_CHECK=1`.
- Agent must never place binaries in any bin directory (not via cp, mv, or go install). Build with `make build`, then ask the user to run the privileged install step. Hooks in block-dangerous-commands.sh enforce this.

---

## [2026-02-26-100007] Task management and exit criteria (consolidated)

**Consolidated from**: 4 entries (2026-01-21 to 2026-02-17)

- Specs get lost without cross-references from TASKS.md. Three-layer defense: (1) playbook instruction, (2) spec reference in Phase header, (3) bold breadcrumb in first task.
- Subtask completion is implementation progress, not delivery. Parent tasks should have explicit deliverables; don't close until deliverable is verified.
- Exit criteria must include verification: integration tests (binary executes correctly), coverage targets, and smoke tests. "All tasks checked off" does not equal "implementation works."
- Reports graduate to ideas/done/ only after all items are tracked or resolved. Cross-reference every item against TASKS.md and the codebase before moving.

---

## [2026-02-26-100008] Agent behavioral patterns (consolidated)

**Consolidated from**: 5 entries (2026-01-25 to 2026-02-22)

- Interaction pattern capture risks softening agent rigor. Do not build implicit user-modeling from session history. Rely on explicit, human-reviewed context (learnings, conventions, hooks) for behavioral shaping.
- Chain-of-thought prompting improves agent reasoning accuracy (17.7% to 78.7%). Added "Reason Before Acting" to AGENT_PLAYBOOK.md and reasoning nudges to 7 skills.
- Say "project conventions" not "idiomatic X" to ensure Claude looks at project files first rather than triggering training priors (stdlib conventions).
- Autonomous "YOLO mode" is effective for feature velocity but accumulates technical debt (magic strings, monolithic tests, hardcoded paths). Schedule periodic consolidation sessions.
- Trust the binary output over source code analysis. A single ambiguous CLI output is not proof of absence — re-run the exact command before claiming something is missing.

---

## [2026-02-26-100009] Hook compliance and output routing (consolidated)

**Consolidated from**: 3 entries (2026-02-22 to 2026-02-25)

- Plain-text hook output is silently ignored by the agent. Claude Code parses hook stdout starting with `{` as JSON directives; plain text is disposable. All hooks should return JSON via `printHookContext()`.
- Hook compliance degrades on narrow mid-session tasks (~15-25% partial skip rate). Root cause: CLAUDE.md's "may or may not be relevant" system reminder competes with hook authority. Fix: CLAUDE.md explicitly elevates hook authority. The mandatory checkpoint relay block is the compliance canary.
- No reliable agent-side before-session-end event exists. SessionEnd fires after the agent is gone. Mid-session nudges and explicit /ctx-wrap-up are the only reliable persistence mechanisms.

---

## [2026-02-26-100010] ctx add and decision recording (consolidated)

**Consolidated from**: 4 entries (2026-01-27 to 2026-02-14)

- `ctx add learning` requires `--context`, `--lesson`, `--application` flags. `ctx add decision` requires `--context`, `--rationale`, `--consequences`. A bare string only sets the title and the command will fail without required flags.
- Structured entries with Context/Lesson/Application are more useful than one-liners. Agents are guided via AGENT_PLAYBOOK.md.
- Always complete decision record sections — placeholder text like "[Add context here]" is a code smell. Decisions without rationale lose their value over time.
- Slash commands using `!` bash syntax require matching permissions in settings.local.json. When adding new /ctx-* commands, ensure ctx init pre-seeds the required `Bash(ctx <subcommand>:*)` permissions.

---

## [2026-02-24-032945] CLI tools don't benefit from in-memory caching of context files

**Context**: Discussed whether ctx should read and cache LEARNINGS.md, DECISIONS.md etc. in memory

**Lesson**: ctx is a short-lived CLI process, not a daemon. Context files are tiny (few KB), sub-millisecond to read. Cache invalidation complexity exceeds the read cost. Caching only makes sense if ctx becomes a long-lived process (MCP server, watch daemon).

**Application**: Don't add caching layers to ctx's file reads. If an MCP server mode is ever added, revisit then.

---

## [2026-02-22-120000] Hook behavior and patterns (consolidated)

**Consolidated from**: 8 entries (2026-01-25 to 2026-02-17)

- Hook scripts receive JSON via stdin (not env vars); parse with `HOOK_INPUT=$(cat)` then jq
- Hook key names are case-sensitive: `PreToolUse` and `SessionEnd` (not `PreToolUseHooks`)
- Use `$CLAUDE_PROJECT_DIR` in hook paths, never hardcode absolute paths
- Hook regex can overfit: `ctx` as binary vs directory name differ; anchor patterns to command-start positions with `(^|;|&&|\|\|)\s*`
- grep patterns match inside quoted arguments — test with `ctx add learning "...blocked words..."` to verify no false positives
- Hook scripts can silently lose execute permission; verify with `ls -la .claude/hooks/*.sh` after edits
- Two-tier output is sufficient: unprefixed (agent context, may or may not relay) and `IMPORTANT: Relay VERBATIM` (guaranteed relay); don't add new severity prefixes
- Repeated injection causes agent repetition fatigue; use `--session $PPID --cooldown 10m` and pair with a readback instruction

---

## [2026-02-22-120001] UserPromptSubmit hook output channels (consolidated)

**Consolidated from**: 2 entries (2026-02-12)

- UserPromptSubmit hook stdout is prepended as AI context (not shown to user); stderr with exit 0 is swallowed entirely
- User-visible output requires `{"systemMessage": "..."}` JSON on stdout (warning banner) or exit 2 (blocks prompt)
- There is no non-blocking user-visible output channel for this hook type
- Design hooks for their actual audience: AI-facing = plain stdout, user-facing = systemMessage JSON

---

## [2026-02-22-120002] Linting and static analysis (consolidated)

**Consolidated from**: 7 entries (2026-01-25 to 2026-02-20)

- Full pre-commit gate: (1) `CGO_ENABLED=0 go build ./cmd/ctx`, (2) `golangci-lint run`, (3) `CGO_ENABLED=0 go test` — all three, every time
- Own the codebase: fix pre-existing lint issues even if you didn't introduce them
- gosec G301/G306: use 0o750 for dirs, 0o600 for files everywhere including tests
- gosec G304 (file inclusion): safe to suppress with `//nolint:gosec` in test files using `t.TempDir()` paths
- golangci-lint errcheck: use `cmd.Printf`/`cmd.Println` in Cobra commands instead of `fmt.Fprintf`
- `defer os.Chdir(x)` fails errcheck; use `defer func() { _ = os.Chdir(x) }()`
- golangci-lint Go version mismatch in CI: use `install-mode: goinstall` to build linter from source

---

## [2026-02-22-120006] Permission and settings drift (consolidated)

**Consolidated from**: 4 entries (2026-02-15)

- Permission drift is distinct from code drift — settings.local.json is gitignored, no review catches stale entries
- `Skill()` permissions don't support name prefix globs — list each skill individually
- Wildcard trusted binaries (`Bash(ctx:*)`, `Bash(make:*)`), but keep git commands granular (never `Bash(git:*)`)
- settings.local.json accumulates session debris; run periodic hygiene via `/sanitize-permissions` and `/ctx-drift`

---

## [2026-02-22-120008] Gitignore and filesystem hygiene (consolidated)

**Consolidated from**: 3 entries (2026-02-11 to 2026-02-15)

- Gitignored directories are invisible to `git status`; stale artifacts persist indefinitely — periodically `ls` gitignored working directories
- Add editor artifacts (*.swp, *.swo, *~) to .gitignore alongside IDE directories from day one
- Gitignore entries for sensitive paths are security controls, not documentation — never remove during cleanup sweeps

---

## [2026-02-19-215200] Feature can be code-complete but invisible to users

**Context**: ctx pad merge was fully implemented with 19 passing tests and binary support, but had zero coverage in user-facing docs (scratchpad.md, cli-reference.md, scratchpad-sync recipe). Only discoverable via --help.

**Lesson**: Implementation completeness \!= user-facing completeness. A feature without docs is invisible to users who don't explore CLI help.

**Application**: After implementing a new CLI subcommand, always check: feature page, cli-reference.md, relevant recipes, and zensical.toml nav (if new page).

---

## [2026-01-28-051426] IDE is already the UI

**Context**: Considering whether to build custom UI for .context/ files

**Lesson**: Discovery, search, and editing of .context/ markdown files works
better in VS Code/IDE than any custom UI we'd build. Full-text search,
git integration, extensions - all free.

**Application**: Don't reinvent the editor. Let users use their preferred IDE.

---

*Module-specific, niche, and historical learnings:
[learnings-reference.md](learnings-reference.md)*
