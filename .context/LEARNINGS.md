# Learnings

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-02-26 | Soft instructions cannot achieve deterministic compliance — compliance paradox |
| 2026-02-26 | Sunk cost leverage inverts rationalization path for context loading |
| 2026-02-26 | Webhook silence after ctxrc profile swap is the most common notify debugging red herring |
| 2026-02-26 | Documentation drift and auditing (consolidated) |
| 2026-02-26 | ctx init and CLAUDE.md behavior (consolidated) |
| 2026-02-26 | Agent context loading and task routing (consolidated) |
| 2026-02-26 | Blog and content publishing (consolidated) |
| 2026-02-26 | Worktrees and parallel agents (consolidated) |
| 2026-02-26 | Go testing patterns (consolidated) |
| 2026-02-26 | PATH and binary handling (consolidated) |
| 2026-02-26 | Task management and exit criteria (consolidated) |
| 2026-02-26 | Agent behavioral patterns (consolidated) |
| 2026-02-26 | Hook compliance and output routing (consolidated) |
| 2026-02-26 | ctx add and decision recording (consolidated) |
| 2026-02-26 | Plugin and marketplace architecture (consolidated) |
| 2026-02-24 | CLI tools don't benefit from in-memory caching of context files |
| 2026-02-24 | /ctx-journal-normalize is dangerous at scale on non-ctx projects |
| 2026-02-24 | url.Parse works for SMB URLs |
| 2026-02-22 | Loop script templates use positional format verbs — adding a verb requires updating all callers |
| 2026-02-22 | Hook behavior and patterns (consolidated) |
| 2026-02-22 | UserPromptSubmit hook output channels (consolidated) |
| 2026-02-22 | Linting and static analysis (consolidated) |
| 2026-02-22 | Journal site rendering pipeline (consolidated) |
| 2026-02-22 | Skills evaluation and design (consolidated) |
| 2026-02-22 | Claude Code session and JSONL format (consolidated) |
| 2026-02-22 | Permission and settings drift (consolidated) |
| 2026-02-22 | Superseded session-related learnings (consolidated) |
| 2026-02-22 | Gitignore and filesystem hygiene (consolidated) |
| 2026-02-21 | Zensical section icons require index pages |
| 2026-02-21 | zensical serve supports -a flag for dev_addr override |
| 2026-02-20 | Default export already preserves enrichment — T2.1 was partially stale |
| 2026-02-19 | Feature can be code-complete but invisible to users |
| 2026-02-19 | GCM authentication makes try-decrypt a reliable format discriminator |
| 2026-02-14 | normalizeCodeFences regex splits language specifiers |
| 2026-02-06 | PROMPT.md deleted — was stale project briefing, not a Ralph loop prompt |
| 2026-02-03 | User input often has inline code fences that break markdown rendering |
| 2026-02-03 | Claude Code injects system-reminder tags into tool results, breaking markdown export |
| 2026-01-28 | IDE is already the UI |
| 2026-01-26 | Go json.Marshal Escapes Shell Characters |
| 2026-01-21 | One Templates Directory, Not Two |
| 2026-01-20 | ctx and Ralph Loop Are Separate Systems |
<!-- INDEX:END -->

---

## [2026-02-26-200000] Soft instructions cannot achieve deterministic compliance — compliance paradox

**Context**: Context-load-gate v1 told the agent to read context files. The agent evaluated that instruction and rationalized skipping it — "don't apply judgment to this rule" is itself evaluated by judgment. Compliance ceiling: ~85% at session start, ~75% mid-session.

**Lesson**: Every soft instruction passes through the same attention/evaluation pipeline. No amount of imperative framing ("STOP", "MANDATORY", "Do not assess relevance") solves this because the instruction is processed by the same mechanism it's trying to override. Only infrastructure-level enforcement (content injection via `additionalContext`, or exit code 2 hard gates) operates outside the reasoning layer.

**Application**: When 100% compliance is required, don't instruct — inject. Use `additionalContext` for content that must be present, exit code 2 for actions that must be blocked. Reserve soft instructions for guidance where ~80% compliance is acceptable.

---

## [2026-02-26-200001] Sunk cost leverage inverts rationalization path for context loading

**Context**: When designing the v2 injection strategy, we considered whether to inject DECISIONS/LEARNINGS indexes or skip them entirely.

**Lesson**: Once ~7k tokens of core context are auto-injected (fait accompli), the agent's rationalization path inverts. Instead of "skip to save effort" (v1), the calculus becomes "I already have 80% of the context, the marginal cost of reading one more entry is trivial." This makes index-only injection effective — the agent sees titles, and the sunk cost makes on-demand reads near-certain.

**Application**: When designing multi-file context loading, front-load the highest-value content as injection (no compliance step), then use the sunk cost to motivate demand-loaded reads for the remainder.

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

## [2026-02-26-100001] ctx init and CLAUDE.md behavior (consolidated)

**Consolidated from**: 4 entries (2026-01-20 to 2026-02-14)

- ctx init is non-destructive: only creates .context/, CLAUDE.md, .claude/, PROMPT.md, and IMPLEMENTATION_PLAN.md. Zero awareness of .cursorrules, .aider.conf.yml, or other tools' configs.
- CLAUDE.md merge insertion is position-aware: findInsertionPoint() finds the first H1, skips trailing blank lines, and inserts there. Never appends to end.
- CLAUDE.md handling is a 3-state machine: no file (create), file without ctx markers (merge/prompt), file with `<!-- ctx:context -->` / `<!-- ctx:end -->` markers (skip or force-replace).
- Always backup before modifying user files: file.bak before modification, marker comments for idempotency, offer merge not overwrite, provide `--merge` escape hatch.

---

## [2026-02-26-100002] Agent context loading and task routing (consolidated)

**Consolidated from**: 5 entries (2026-01-20 to 2026-01-25)

- `ctx agent` is optimized for task execution (filters pending tasks, surfaces constitution, token-budget aware). Manual file reading is better for exploratory/memory questions (session history, timestamps, completed tasks).
- On "Do you remember?" questions, immediately read .context/ files and run `ctx recall list --limit 5`. Never ask "would you like me to check?" — that is the obvious intent.
- .context/ is NOT a Claude Code primitive. Only CLAUDE.md and .claude/settings.json are auto-loaded. The .context/ directory requires a hook or explicit CLAUDE.md instruction to be discovered.
- Orchestrator (IMPLEMENTATION_PLAN.md) and agent (.context/TASKS.md) task lists must be separate. The orchestrator says "check your mind" — it doesn't maintain a parallel ledger.
- Only CLAUDE.md is auto-loaded by Claude Code. Projects using ctx should rely on the CLAUDE.md -> AGENT_PLAYBOOK.md chain, not AGENTS.md.

---

## [2026-02-26-100003] Blog and content publishing (consolidated)

**Consolidated from**: 4 entries (2026-01-28 to 2026-02-18)

- Blog posts are living documents: periodic enrichment passes (cross-links, update admonitions, citations) improve reader experience. Schedule enrichment as part of consolidation sessions.
- Blog publishing from ideas/ requires a 7-step checklist: (1) update date/frontmatter, (2) fix relative paths, (3) add cross-links, (4) add "The Arc" section, (5) update blog index, (6) add "See also" in related posts, (7) verify all link targets exist.
- Changelogs document WHAT (for machines/audits); blogs explain WHY (for humans/narrative). When synthesizing session history, output both.
- Cross-repo links to published docs should use ctx.ist URLs, not GitHub blob URLs or relative paths. GitHub won't render zensical admonitions and readers lose navigation context.

---

## [2026-02-26-100004] Worktrees and parallel agents (consolidated)

**Consolidated from**: 5 entries (2026-02-12 to 2026-02-24)

- Git worktrees enable parallel Claude Code sessions without file conflicts. Create as sibling directories (`git worktree add ../ctx-docs -b work/docs`). 3-4 parallel worktrees is the practical limit before merge complexity outweighs productivity gains.
- Parallel agents work cleanly on disjoint file sets: partition by file domain (Go source, docs, specs/context) for zero-conflict parallel execution.
- Worktree agents lack key-dependent features by design: ctx pad fails gracefully (no key), ctx notify silently no-ops. The only real gap is notify — pad and journal are naturally avoided by workflow.
- Multi-agent work requires commit auditing: changes intermingle in the working tree with no commit boundaries. Run `git diff --stat` and group by logical feature before committing.
- rsync between worktrees can clobber permissions and gitignored files. Use `--no-perms` or `--chmod=+x` for scripts, and `--exclude` gitignored paths. The /absorb skill handles these edge cases.

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

## [2026-02-26-100011] Plugin and marketplace architecture (consolidated)

**Consolidated from**: 3 entries (2026-02-16)

- When repo-local .claude/skills/ and a marketplace plugin both define the same skill name, Claude Code lists both: local unprefixed and plugin with `ctx:` namespace prefix. Ensure distributed skills live only in the plugin source.
- Claude Code marketplace plugins source from the repo root where `.claude-plugin/marketplace.json` lives. Edits to skills and hooks under the plugin path take effect on next Claude Code load — no reinstall needed.
- Security docs are most vulnerable to stale paths after architecture migrations. After any file-layout migration, grep security docs for old paths first — stale paths in security guidance give users a false sense of protection.

---

## [2026-02-24-032945] CLI tools don't benefit from in-memory caching of context files

**Context**: Discussed whether ctx should read and cache LEARNINGS.md, DECISIONS.md etc. in memory

**Lesson**: ctx is a short-lived CLI process, not a daemon. Context files are tiny (few KB), sub-millisecond to read. Cache invalidation complexity exceeds the read cost. Caching only makes sense if ctx becomes a long-lived process (MCP server, watch daemon).

**Application**: Don't add caching layers to ctx's file reads. If an MCP server mode is ever added, revisit then.

---

## [2026-02-24-022214] /ctx-journal-normalize is dangerous at scale on non-ctx projects

**Context**: Discussed whether to keep normalize in the default journal pipeline

**Lesson**: On projects with large session JSONL files (millions of lines), the normalize skill blows up subagent context windows, consumes excessive tokens, and produces nondeterministic half-baked outputs

**Application**: Keep expensive AI skills out of batch pipelines; offer them as targeted per-file tools instead

---

## [2026-02-24-015713] url.Parse works for SMB URLs

**Context**: Shell script used sed to extract host/share from smb://host/share. Needed a Go equivalent.

**Lesson**: Go net/url.Parse handles smb:// scheme correctly — u.Host gives hostname, u.Path gives share path with leading slash.

**Application**: Use url.Parse for any custom-scheme URL parsing instead of hand-rolled regex.

---

## [2026-02-22-101959] Loop script templates use positional format verbs — adding a verb requires updating all callers

**Context**: Adding TplLoopNotify to TplLoopScript and TplLoopMaxIter required updating the fmt.Sprintf calls in script.go to pass the new argument at the correct position.

**Lesson**: Go fmt.Sprintf with positional args is fragile — adding a %s to a template shifts all downstream arguments. The comment documenting Args: in the template constant is the only guard.

**Application**: When modifying loop templates, always update the Args: comment and verify script.go passes arguments in the correct order.

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

## [2026-02-22-120003] Journal site rendering pipeline (consolidated)

**Consolidated from**: 8 entries (2026-02-20)

- Python-Markdown ends ALL HTML blocks at blank lines — no CommonMark Type 1 exception for `<pre>`. Only fenced code blocks survive blank lines
- `<details>` is Type 6 HTML block — ends at first blank line; collapsible tool output needs CSS/JS approach
- pymdownx.highlight hijacks `<pre><code>` patterns; disable with `use_pygments=false`
- normalizeContent three-layer fix: always run wrapToolOutputs, pick fence depth exceeding inner fences, stripPreWrapper only unescapes when `<pre>` found
- Tool output boundary detection: pre-scan turn numbers, sort+dedup, find min > N, use LAST positional occurrence (last-match-wins)
- Inline code spans with angle brackets: replace backticks with double-quotes and brackets with HTML entities via RegExInlineCodeAngle
- Title sanitization: strip Claude tags, replace angle brackets, strip backticks/hash, truncate to 75 chars (RecallMaxTitleLen)
- AI normalization poor fit for bulk files; code-level pipeline handles rendering at build time

---

## [2026-02-22-120004] Skills evaluation and design (consolidated)

**Consolidated from**: 6 entries (2026-01-23 to 2026-02-15)

- Skills are markdown files in `.claude/commands/` with YAML frontmatter; `$ARGUMENTS` passes args
- Skills that restate or contradict system prompt create tension — check platform defaults first
- Red flags: urgency tags, reasoning overrides, tables labeling hesitation as wrong — discard entirely
- ~80% of external skill files are redundant; apply E/A/R classification, only keep expert-level knowledge delta
- Prefer skills over CLI commands for judgment-heavy workflows; reserve CLI for deterministic operations
- When a skill would edit files controlling agent behavior (permissions, hooks), use a runbook instead

---

## [2026-02-22-120005] Claude Code session and JSONL format (consolidated)

**Consolidated from**: 4 entries (2026-01-28 to 2026-02-04)

- Claude Code, Cursor, Aider each have different JSONL formats; use tool-agnostic Session type with tool-specific parsers
- JSONL files are append-only, never shrink after compaction; file size overreports post-compaction
- Subagent files in `/subagents/` share parent sessionId; skip when scanning (check `isSidechain:true`)
- `slug` field removed in Claude Code v2.1.29+; parse by `sessionId` + valid `type` instead

---

## [2026-02-22-120006] Permission and settings drift (consolidated)

**Consolidated from**: 4 entries (2026-02-15)

- Permission drift is distinct from code drift — settings.local.json is gitignored, no review catches stale entries
- `Skill()` permissions don't support name prefix globs — list each skill individually
- Wildcard trusted binaries (`Bash(ctx:*)`, `Bash(make:*)`), but keep git commands granular (never `Bash(git:*)`)
- settings.local.json accumulates session debris; run periodic hygiene via `/sanitize-permissions` and `/ctx-drift`

---

## [2026-02-22-120007] Superseded session-related learnings (consolidated)

**Consolidated from**: 4 entries (2026-01-20)

> **Historical**: `.context/sessions/` was removed in v0.4.0. These learnings are superseded but preserved for context.

- SessionEnd hook fires on all exits including Ctrl+C — hook behavior still accurate for Claude Code, but ctx no longer uses it
- Session filenames used YYYY-MM-DD-HHMMSS-topic.md — journal entries now use `ctx recall export` naming
- Two tiers remain: curated (`.context/*.md`) and full dump (`~/.claude/projects/` + `.context/journal/`); middle `.context/sessions/` tier eliminated
- Auto-load via PreToolUse worked; auto-save via SessionEnd removed because Claude Code retains transcripts natively

---

## [2026-02-22-120008] Gitignore and filesystem hygiene (consolidated)

**Consolidated from**: 3 entries (2026-02-11 to 2026-02-15)

- Gitignored directories are invisible to `git status`; stale artifacts persist indefinitely — periodically `ls` gitignored working directories
- Add editor artifacts (*.swp, *.swo, *~) to .gitignore alongside IDE directories from day one
- Gitignore entries for sensitive paths are security controls, not documentation — never remove during cleanup sweeps

---

## [2026-02-21-200036] Zensical section icons require index pages

**Context**: Mobile nav showed icons for Manifesto, Blog, and Recipes but not Reference, Operations, or Security

**Lesson**: Zensical (like MkDocs Material) only renders section-header icons when the section has an index.md listed as its first nav entry, via the navigation.indexes feature. Icons in child page frontmatter don't propagate to the section header.

**Application**: When adding a new top-level nav section, always create a <section>/index.md with an icon: frontmatter field and list it first in zensical.toml

---

## [2026-02-21-195840] zensical serve supports -a flag for dev_addr override

**Context**: Needed a way to override dev_addr without modifying toml files

**Lesson**: zensical serve -a IP:PORT overrides the config file dev_addr. ctx journal site --serve does not pass through extra flags to zensical — it hardcodes zensical serve with no args

**Application**: For any future zensical config that needs per-developer overrides, prefer CLI flags in Make targets over config file changes

---

## [2026-02-20-142442] Default export already preserves enrichment — T2.1 was partially stale

**Context**: Investigated ctx recall export --update and found the default behavior already preserves YAML frontmatter during re-export. The --force flag has a bug where it claims to discard frontmatter but does not.

**Lesson**: Always read the current code before speccing a feature — the need may already be met, and the real work may be a bug fix rather than a new feature.

**Application**: When speccing tasks from the backlog, investigate current state first. Rewrite the task to reflect what is actually needed.

---

## [2026-02-19-215200] Feature can be code-complete but invisible to users

**Context**: ctx pad merge was fully implemented with 19 passing tests and binary support, but had zero coverage in user-facing docs (scratchpad.md, cli-reference.md, scratchpad-sync recipe). Only discoverable via --help.

**Lesson**: Implementation completeness \!= user-facing completeness. A feature without docs is invisible to users who don't explore CLI help.

**Application**: After implementing a new CLI subcommand, always check: feature page, cli-reference.md, relevant recipes, and zensical.toml nav (if new page).

---

## [2026-02-19-214909] GCM authentication makes try-decrypt a reliable format discriminator

**Context**: Needed to auto-detect whether pad merge input files are encrypted or plaintext without relying on file extensions or user flags.

**Lesson**: Authenticated encryption (AES-256-GCM) guarantees that decryption with the wrong key always fails — unlike unauthenticated ciphers that produce silent garbage. This makes 'try decrypt, fall back to plaintext' a safe and simple detection strategy.

**Application**: Use try-decrypt-first as the default pattern for any ctx feature that handles mixed encrypted/plaintext input. No need for format flags or extension-based heuristics.

---

## [2026-02-14-163549] normalizeCodeFences regex splits language specifiers

**Context**: Writing test for normalizeCodeFences, expected inline fence with lang tag to stay joined but the regex matched characters after backticks

**Lesson**: The inline fence regex treats any non-whitespace adjacent to triple-backtick fences as a split point, separating lang tags from the fence

**Application**: When testing normalizeCodeFences, use plain fences without language tags. See internal/cli/recall/fmt_test.go.

---

## [2026-02-06-200000] PROMPT.md deleted — was stale project briefing, not a Ralph loop prompt

**Context**: During consolidation, reviewed PROMPT.md and found it had drifted
into a stale project briefing — duplicating CLAUDE.md (session start/end rituals,
build commands, context file table) and containing outdated Phase 2 monitor
architecture diagrams for work that was already completed differently.

**Lesson**: PROMPT.md's actual purpose is as a Ralph loop iteration prompt: a
focused "what to do next and how to know when done" document consumed by
`ctx loop` between iterations. CLAUDE.md serves a different role: always-loaded
project operating manual for Claude Code. When PROMPT.md drifts into duplicating
CLAUDE.md, it becomes stale weight that misleads future sessions.

**Application**: Re-introduce PROMPT.md only when actively using Ralph loops.
Keep it to: iteration goal + completion signal + current phase focus. Project
context (build commands, file tables, session rituals) belongs in CLAUDE.md and
.context/ files, not PROMPT.md.

---

## [2026-02-03-160000] User input often has inline code fences that break markdown rendering

**Context**: Journal export showed broken code blocks where user typed
`text: ```code` on a single line without proper newlines before/after the
code fence.

**Lesson**: Users naturally type inline code fences like `This is the error:
```Error: foo```. Markdown requires code fences to be on their own lines with
blank lines separating them. You can't force users to format correctly,
but you can normalize on export.

**Application**: Use regex to detect fences preceded/followed by non-whitespace
on same line. Insert `\n\n` to ensure proper spacing. Apply only to user
messages (assistant output is already well-formatted).

---

## [2026-02-03-154500] Claude Code injects system-reminder tags into tool results, breaking markdown export

**Context**: Journal site had rendering errors starting from "Tool Output"
sections. A closing triple-backtick appeared orphaned. Investigation traced
it to `<system-reminder>` tags in the JSONL source - 32 occurrences in one
session file.

**Lesson**: Claude Code injects `<system-reminder>...</system-reminder>` blocks
into tool result content before storing in JSONL. When exported to markdown
and wrapped in code fences, these XML-like tags break rendering - some
markdown parsers treat them as HTML, causing the closing fence to appear as
orphaned literal text instead of terminating the code block.

**Application**: Extract system reminders from tool result content before
wrapping in code fences. Render them as markdown (`**System Reminder**: ...`)
outside the fence. This preserves the information (useful for debugging Claude
Code behavior) while fixing the rendering issue.

---

## [2026-01-28-051426] IDE is already the UI

**Context**: Considering whether to build custom UI for .context/ files

**Lesson**: Discovery, search, and editing of .context/ markdown files works
better in VS Code/IDE than any custom UI we'd build. Full-text search,
git integration, extensions - all free.

**Application**: Don't reinvent the editor. Let users use their preferred IDE.

---

## [2026-01-26-180000] Go json.Marshal Escapes Shell Characters

**Context**: Generated settings.local.json had `2\u003e/dev/null` instead
of `2>/dev/null`.

**Lesson**: Go's json.Marshal escapes `>`, `<`, and `&` as
unicode (`\u003e`, `\u003c`, `\u0026`) by default for HTML safety.
Use `json.Encoder` with `SetEscapeHTML(false)` when generating config files
that contain shell commands.

---

## [2026-01-21-140000] One Templates Directory, Not Two

**Context**: Confusion arose about `templates/` (root) vs
`internal/templates/` (embedded).

**Lesson**: Only `internal/templates/` matters — it's where Go embeds files
into the binary. A root `templates/` directory is spec baggage that serves
no purpose.

**The actual flow:**
```
internal/templates/  ──[ctx init]──>  .context/
     (baked into binary)              (agent's working copy)
```

**Application**: Don't create duplicate template directories. One source of truth.

---

## [2026-01-20-200000] ctx and Ralph Loop Are Separate Systems

**Context**: User asked "How do I use the ctx binary to recreate this project?"

**Lesson**: `ctx` and Ralph Loop are two distinct systems:
- `ctx init` creates `.context/` for context management (decisions, learnings, tasks)
- Ralph Loop uses PROMPT.md, IMPLEMENTATION_PLAN.md, specs/ for iterative AI development
- `ctx` does NOT create Ralph Loop infrastructure

**Application**: To bootstrap a new project with both:
1. Run `ctx init` to create `.context/`
2. Manually copy/adapt PROMPT.md, AGENTS.md, specs/ from a reference project
3. Create IMPLEMENTATION_PLAN.md with your tasks
4. Run `/ralph-loop` to start iterating

---
