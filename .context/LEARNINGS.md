# Learnings

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-02-24 | CLI tools don't benefit from in-memory caching of context files |
| 2026-02-24 | Worktree agents lack key-dependent features by design |
| 2026-02-24 | /ctx-journal-normalize is dangerous at scale on non-ctx projects |
| 2026-02-24 | ARCHITECTURE.md had significant drift — 4 core packages and 4 CLI commands missing |
| 2026-02-24 | All runCmd() returns must be consumed in tests |
| 2026-02-24 | url.Parse works for SMB URLs |
| 2026-02-22 | Interaction pattern capture risks softening agent rigor |
| 2026-02-22 | No reliable agent-side before-session-end event exists |
| 2026-02-22 | Plain-text hook output is silently ignored by the agent |
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
| 2026-02-21 | Parallel agents work cleanly on disjoint file sets |
| 2026-02-21 | Zensical section icons require index pages |
| 2026-02-21 | zensical serve supports -a flag for dev_addr override |
| 2026-02-21 | Multi-agent work requires commit auditing before batching |
| 2026-02-20 | Default export already preserves enrichment — T2.1 was partially stale |
| 2026-02-19 | Trust the binary output over source code analysis |
| 2026-02-19 | Feature can be code-complete but invisible to users |
| 2026-02-19 | GCM authentication makes try-decrypt a reliable format discriminator |
| 2026-02-18 | Blog posts are living documents |
| 2026-02-17 | rsync between worktrees can clobber permissions and gitignored files |
| 2026-02-16 | Security docs are most vulnerable to stale paths after architecture migrations |
| 2026-02-16 | Duplicate skills appear with namespace prefix in Claude Code |
| 2026-02-16 | Local marketplace plugin enables live skill editing |
| 2026-02-15 | Dead link checking is consolidation check 12, not a standalone concern |
| 2026-02-15 | Cross-repo links to published docs should use ctx.ist |
| 2026-02-14 | ctx add learning/decision requires structured flags, not just a string |
| 2026-02-14 | ctx init is non-destructive toward tool-specific configs |
| 2026-02-14 | merge insertion is position-aware, not append |
| 2026-02-14 | ctx init CLAUDE.md handling is a 3-state machine |
| 2026-02-14 | color.NoColor in init for CLI test files |
| 2026-02-14 | Recall CLI tests isolate via HOME env var |
| 2026-02-14 | formatDuration accepts interface not time.Duration |
| 2026-02-14 | normalizeCodeFences regex splits language specifiers |
| 2026-02-13 | Specs get lost without cross-references from TASKS.md |
| 2026-02-11 | Chain-of-thought prompting improves agent reasoning accuracy |
| 2026-02-06 | PROMPT.md deleted — was stale project briefing, not a Ralph loop prompt |
| 2026-02-03 | User input often has inline code fences that break markdown rendering |
| 2026-02-03 | Claude Code injects system-reminder tags into tool results, breaking markdown export |
| 2026-01-30 | Say 'project conventions' not 'idiomatic X' |
| 2026-01-29 | Documentation audits require verification against actual standards |
| 2026-01-28 | Required flags now enforced for learnings |
| 2026-01-28 | Changelogs vs Blogs serve different audiences |
| 2026-01-28 | IDE is already the UI |
| 2026-01-28 | Subtasks complete does not mean parent task complete |
| 2026-01-27 | Always Complete Decision Record Sections |
| 2026-01-27 | Slash Commands Require Matching Permissions |
| 2026-01-26 | Go json.Marshal Escapes Shell Characters |
| 2026-01-25 | CI Tests Need CTX_SKIP_PATH_CHECK |
| 2026-01-25 | AGENTS.md Is Not Auto-Loaded |
| 2026-01-25 | Autonomous Mode Creates Technical Debt |
| 2026-01-23 | ctx agent vs Manual File Reading Trade-offs |
| 2026-01-23 | Infer Intent on "Do You Remember?" Questions |
| 2026-01-23 | Always Use ctx from PATH |
| 2026-01-21 | Exit Criteria Must Include Verification |
| 2026-01-21 | Orchestrator vs Agent Tasks Must Be Separate |
| 2026-01-21 | One Templates Directory, Not Two |
| 2026-01-21 | Hooks Should Use PATH, Not Hardcoded Paths |
| 2026-01-20 | ctx and Ralph Loop Are Separate Systems |
| 2026-01-20 | .context/ Is NOT a Claude Code Primitive |
| 2026-01-20 | Always Backup Before Modifying User Files |
| 2026-01-19 | CGO Must Be Disabled for ARM64 Linux |
<!-- INDEX:END -->

---

## [2026-02-24-032945] CLI tools don't benefit from in-memory caching of context files

**Context**: Discussed whether ctx should read and cache LEARNINGS.md, DECISIONS.md etc. in memory

**Lesson**: ctx is a short-lived CLI process, not a daemon. Context files are tiny (few KB), sub-millisecond to read. Cache invalidation complexity exceeds the read cost. Caching only makes sense if ctx becomes a long-lived process (MCP server, watch daemon).

**Application**: Don't add caching layers to ctx's file reads. If an MCP server mode is ever added, revisit then.

---

## [2026-02-24-030252] Worktree agents lack key-dependent features by design

**Context**: Investigated whether hooks/notify/pad would break in git worktrees when .context.key is gitignored

**Lesson**: ctx pad fails gracefully (no key), ctx notify silently no-ops, journal enrichment writes to the worktree and is orphaned on teardown. All path resolution is cwd-relative with no git-root awareness. The only real gap is notify — pad and journal are naturally avoided by workflow.

**Application**: Document worktree limitations in skill docs, recipes, and reference pages. File a task to enable notify in worktrees (the one feature that would genuinely help autonomous opaque agents).

---

## [2026-02-24-022214] /ctx-journal-normalize is dangerous at scale on non-ctx projects

**Context**: Discussed whether to keep normalize in the default journal pipeline

**Lesson**: On projects with large session JSONL files (millions of lines), the normalize skill blows up subagent context windows, consumes excessive tokens, and produces nondeterministic half-baked outputs

**Application**: Keep expensive AI skills out of batch pipelines; offer them as targeted per-file tools instead

---

## [2026-02-24-015941] ARCHITECTURE.md had significant drift — 4 core packages and 4 CLI commands missing

**Context**: During the first /ctx-map run, analysis revealed crypto, sysinfo, notify, journal/state were missing from the core packages table, and notify, pad, permissions, system CLI commands were absent. The doc claimed 19 commands (actually 22) and 16 skill templates (actually 28).

**Lesson**: ARCHITECTURE.md drifts silently when new packages are added without updating the doc. The existing /ctx-drift skill catches stale paths but not missing packages — it cannot detect what is absent from a table.

**Application**: Run /ctx-map after adding new packages or CLI commands. The tracking file staleness detection catches modules with new commits, but new modules need the first-run survey to be discovered.

---

## [2026-02-24-015827] All runCmd() returns must be consumed in tests

**Context**: golangci-lint errcheck flagged every runCmd(...) call in remind_test.go that didn't capture the return

**Lesson**: Even setup calls in tests that run commands as preconditions need '_, _ = runCmd(...)' to satisfy errcheck

**Application**: When writing test helpers that call cobra commands as setup, always capture both returns

---

## [2026-02-24-015713] url.Parse works for SMB URLs

**Context**: Shell script used sed to extract host/share from smb://host/share. Needed a Go equivalent.

**Lesson**: Go net/url.Parse handles smb:// scheme correctly — u.Host gives hostname, u.Path gives share path with leading slash.

**Application**: Use url.Parse for any custom-scheme URL parsing instead of hand-rolled regex.

---

## [2026-02-22-212342] Interaction pattern capture risks softening agent rigor

**Context**: Considered a skill to analyze session history and extract user interaction patterns into .context/

**Lesson**: Automated pattern capture from sessions risks training the agent to please rather than push back. Scoping to 'process only' is insufficient — the agent doing the analysis has an incentive to learn the wrong lessons. Existing mechanisms (learnings, hooks, constitution) already capture process preferences explicitly and with human review.

**Application**: Do not build implicit user-modeling from session history. Rely on explicit, human-reviewed context (learnings, conventions, hooks) for behavioral shaping. If interaction patterns are ever revisited, require a CONSTITUTION clause preventing compliance-oriented learning.

---

## [2026-02-22-194443] No reliable agent-side before-session-end event exists

**Context**: Investigated whether SessionEnd hooks could trigger context persistence automatically

**Lesson**: SessionEnd fires after the agent is gone — it can only run fire-and-forget shell commands (cleanup), not LLM reasoning. CTRL+C terminates immediately. The playbook already encodes the correct strategy: persist as you go.

**Application**: Mid-session nudges (check-persistence) and explicit /ctx-wrap-up are the only reliable persistence mechanisms; do not design hooks that assume the agent gets a final turn

---

## [2026-02-22-194441] Plain-text hook output is silently ignored by the agent

**Context**: The qa-reminder PreToolUse:Edit hook was firing successfully (confirmed in debug log) but producing no visible effect on agent behavior

**Lesson**: Claude Code parses hook stdout that starts with { as JSON directives; plain text is injected as context but the agent treats it as disposable. Structured JSON with hookSpecificOutput.additionalContext is reliably processed.

**Application**: All non-blocking PreToolUse/PostToolUse hooks should return JSON via printHookContext(), not cmd.Println() plain text

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

## [2026-02-21-200039] Parallel agents work cleanly on disjoint file sets

**Context**: The .contextrc → .ctxrc rename touched ~40 files across Go source, docs, specs, and context files

**Lesson**: Splitting work across 4 parallel agents by file-type (Go source, new files, docs, specs/context) with zero file overlap completed the entire rename in one pass with no conflicts

**Application**: For broad renames or refactors, partition by file domain and run agents in parallel rather than sequentially

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

## [2026-02-21-195820] Multi-agent work requires commit auditing before batching

**Context**: Found 122 uncommitted files from 3+ agent sessions spanning different features (knowledge health, recall safety, system ergonomics)

**Lesson**: When multiple agents work in parallel, their changes intermingle in the working tree with no commit boundaries. An audit pass is needed to group changes logically before committing.

**Application**: After multi-agent sessions, run git diff --stat and group by logical feature before committing. Don't commit everything in one giant batch.

---

## [2026-02-20-142442] Default export already preserves enrichment — T2.1 was partially stale

**Context**: Investigated ctx recall export --update and found the default behavior already preserves YAML frontmatter during re-export. The --force flag has a bug where it claims to discard frontmatter but does not.

**Lesson**: Always read the current code before speccing a feature — the need may already be met, and the real work may be a bug fix rather than a new feature.

**Application**: When speccing tasks from the backlog, investigate current state first. Rewrite the task to reflect what is actually needed.

---

## [2026-02-19-215204] Trust the binary output over source code analysis

**Context**: Wrongly concluded ctx decisions archive was missing from the installed binary based on a single CLI test that showed parent help instead of subcommand help. The user's own terminal showed it working fine.

**Lesson**: A single ambiguous CLI output is not proof of absence. Re-run the exact command before claiming something is missing. When the user contradicts your finding, they are probably right.

**Application**: When checking if a subcommand exists, run the subcommand directly (e.g., ctx decisions archive --help) and if results are ambiguous, retry before drawing conclusions.

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

## [2026-02-18-071508] Blog posts are living documents

**Context**: Session spent enriching two blog posts with cross-links, update admonitions, citations, and contextual notes. Every post had 3-6 places where a link or admonition improved reader experience.

**Lesson**: Blog posts benefit from periodic enrichment passes: cross-linking to newer content, adding update admonitions for superseded features, citing sources, and adding contextual admonitions that connect ideas across posts.

**Application**: Schedule blog enrichment as part of consolidation sessions. When a new feature supersedes something described in a blog post, add an update admonition immediately rather than waiting.

---

## [2026-02-17] Blog publishing from ideas/ requires a consistent checklist

**Context**: Published 4 blog posts from ideas/ drafts in one session. Each required the same steps: date update, path fixes, cross-links, Arc section, blog index, See also in companions. Missing any step left broken links or orphaned posts.

**Lesson**: Blog publishing is a repeatable workflow with 7 steps: (1) update date and frontmatter, (2) fix relative paths from ideas/ to docs/blog/, (3) add cross-links to/from companion posts, (4) add "The Arc" section connecting to the series narrative, (5) update blog index, (6) add "See also" in related posts, (7) verify all link targets exist.

**Application**: Follow this checklist for every ideas/ → docs/blog/ promotion. Consider making it a recipe in hack/runbooks/ if the pattern continues.

---

## [2026-02-17] Reports graduate to ideas/done/ only after all items are tracked or resolved

**Context**: Moving REPORT-6 and REPORT-7 to ideas/done/. Each had a mix of completed, skipped, and untracked items. Moving before tracking would lose the untracked items.

**Lesson**: Before graduating a report: (1) cross-reference every item against TASKS.md and the codebase, (2) add trackers for undone items, (3) create specs for items that need design, (4) put remaining low-priority items in a future-considerations document, (5) update TASKS.md path references, (6) then move.

**Application**: Always do the full cross-reference before moving reports to done/. The report is the source of truth until every item has a home elsewhere.

---

## [2026-02-17] Agent must never place binaries — nudge the user to install

**Context**: Agent removed ~/go/bin/ctx and discussed copying to /usr/local/bin. Both actions bypass the proper installation path (make install with elevated privileges) which the agent cannot run.

**Lesson**: The agent must never place binaries in any bin directory — not via cp, mv, go install, or any other mechanism. When a rebuild is needed, the agent builds with `make build` and asks the user to run the privileged install step themselves.

**Application**: When ctx binary is stale or missing: (1) run `make build`, (2) ask the user to install it (requires privileges), (3) wait for confirmation before continuing. Hooks in block-dangerous-commands.sh now block cp/mv to bin dirs and `go install` as a command.

---

## [2026-02-17-183937] rsync between worktrees can clobber permissions and gitignored files

**Context**: Used rsync -av to borrow upstream changes; it overwrote .claude/hooks/*.sh with non-executable copies and clobbered gitignored settings.local.json

**Lesson**: rsync -av preserves source permissions, not destination. Gitignored files have no git safety net. Use --no-perms or --chmod=+x for scripts, and --exclude gitignored paths explicitly.

**Application**: When borrowing between worktrees: 1) exclude gitignored paths (.claude/settings.local.json, ideas/, .context/logs/) 2) restore +x on hook scripts after sync 3) consider the /absorb skill which handles these edge cases

---

## [2026-02-16-164547] Security docs are most vulnerable to stale paths after architecture migrations

**Context**: Migrated from per-project .claude/hooks/ and .claude/skills/ to plugin model; found 5 security docs still referencing the old paths

**Lesson**: When moving infrastructure from per-project files to a plugin/external model, audit security docs first — stale paths in security guidance give users a false sense of protection (e.g. 'make .claude/hooks/ immutable' for a directory that no longer exists)

**Application**: After any file-layout migration, grep security and agent-security docs for old paths before anything else

---

## [2026-02-16-164521] Duplicate skills appear with namespace prefix in Claude Code

**Context**: Had both .claude/skills/ctx-status and the marketplace plugin providing the same skill

**Lesson**: When a repo-local .claude/skills/ directory and a marketplace plugin both define the same skill name, Claude Code lists both: the local version unprefixed and the plugin version with a ctx: namespace prefix (e.g. ctx-status and ctx:ctx-status)

**Application**: To avoid confusing duplicates, ensure distributed skills live only in the plugin source (internal/assets/claude/skills/) and not also in .claude/skills/. Dev-only skills that aren't in the plugin won't collide.

---

## [2026-02-16-164518] Local marketplace plugin enables live skill editing

**Context**: Setting up the contributor workflow for ctx development

**Lesson**: Claude Code marketplace plugins source from the repo root where `.claude-plugin/marketplace.json` lives (e.g. ~/WORKSPACE/ctx). The marketplace.json points to the actual plugin in `internal/assets/claude`. Edits to skills and hooks under that path take effect on the next Claude Code load — no reinstall needed

**Application**: The contributor docs instruct devs to add their local clone as a marketplace source rather than using the GitHub URL. This gives them live feedback on skill changes without a rebuild cycle.

---

## [2026-02-15-231022] Dead link checking is consolidation check 12, not a standalone concern

**Context**: User identified dead links in rendered site as a problem. Initial instinct was a standalone task or drift extension.

**Lesson**: Doc link rot is code-level drift — same category as magic strings or stale architecture diagrams. It belongs in /consolidate's check list, with a standalone /check-links skill that consolidate invokes.

**Application**: When a new audit concern emerges, check if it fits an existing audit skill before creating an isolated one. Consolidate is the natural home for anything that drifts silently between sessions.

---

## [2026-02-15-040313] Cross-repo links to published docs should use ctx.ist

**Context**: hack/runbooks/persistent-irc.md linked to docs/ via relative paths, getting-started.md linked to MANIFESTO.md via GitHub — both bypass ctx.ist rendering (admonitions, nav, search)

**Lesson**: When content is published on ctx.ist, always link to the site URL, not the GitHub blob or a relative file path. GitHub won't render zensical admonitions and readers lose navigation context.

**Application**: When adding See Also or cross-references in runbooks or docs, use https://ctx.ist/... URLs for anything the site publishes. Reserve GitHub links for repo-only content (issues, releases, security tab, source files not on the site).

---

## [2026-02-14-164053] ctx add learning/decision requires structured flags, not just a string

**Context**: Repeatedly suggested bare ctx add learning '...' in session endings despite this being wrong

**Lesson**: Learnings require --context, --lesson, --application. Decisions require --context, --rationale, --consequences. A bare string only sets the title — the command will fail without the required flags.

**Application**: Never suggest ctx add learning 'text' as a one-liner. Always show the full flag form. The CLAUDE.md template and session-end prompts should model the correct syntax.

---

## [2026-02-14-164029] ctx init is non-destructive toward tool-specific configs

**Context**: Verified by reading run.go — no code paths touch .cursorrules, .aider.conf.yml, or copilot instructions

**Lesson**: ctx init only creates .context/, CLAUDE.md, .claude/, PROMPT.md, and IMPLEMENTATION_PLAN.md. It has zero awareness of other tools' config files.

**Application**: State this definitively in docs rather than hedging — it's confirmed by the code

---

## [2026-02-14-164013] merge insertion is position-aware, not append

**Context**: Reading fs.go findInsertionPoint() to document --merge behavior

**Lesson**: The --merge flag finds the first H1 heading, skips trailing blank lines, and inserts the ctx block there. If no H1 is found, it inserts at the top. Content is never appended to the end.

**Application**: Document the insertion position clearly — users care about where their content ends up in the merged file

---

## [2026-02-14-164011] ctx init CLAUDE.md handling is a 3-state machine

**Context**: Reading claude.go to write the migration guide

**Lesson**: ctx init checks for: no file (create), file without ctx markers (merge/prompt), file with markers (skip or force-replace). The markers <!-- ctx:context --> / <!-- ctx:end --> are the pivot.

**Application**: When documenting merge behavior, describe all three states explicitly rather than just the happy path

---

## [2026-02-14-163552] color.NoColor in init for CLI test files

**Context**: Recall CLI tests had ANSI escape codes in output making string assertions unreliable

**Lesson**: Setting color.NoColor = true in a package-level init function disables ANSI codes for all tests in the package

**Application**: Add init with color.NoColor = true in test files for CLI packages that use fatih/color. Cleaner than per-test setup.

---

## [2026-02-14-163551] Recall CLI tests isolate via HOME env var

**Context**: Needed integration tests for recall list/show/export without touching real session data

**Lesson**: parser.FindSessions reads os.UserHomeDir which uses HOME env var. Setting t.Setenv HOME to tmpDir with .claude/projects/ structure gives full isolation.

**Application**: For recall integration tests: t.Setenv HOME to tmpDir, create .claude/projects/dir/ with JSONL fixtures. See internal/cli/recall/run_test.go.

---

## [2026-02-14-163550] formatDuration accepts interface not time.Duration

**Context**: Writing unit tests for formatDuration in recall/fmt.go

**Lesson**: formatDuration takes interface with Minutes method, not time.Duration directly. A stub type is needed for testing.

**Application**: Use a stubDuration struct with a mins field and Minutes method when testing formatDuration. See internal/cli/recall/fmt_test.go.

---

## [2026-02-14-163549] normalizeCodeFences regex splits language specifiers

**Context**: Writing test for normalizeCodeFences, expected inline fence with lang tag to stay joined but the regex matched characters after backticks

**Lesson**: The inline fence regex treats any non-whitespace adjacent to triple-backtick fences as a split point, separating lang tags from the fence

**Application**: When testing normalizeCodeFences, use plain fences without language tags. See internal/cli/recall/fmt_test.go.

---

## [2026-02-13-133314] Specs get lost without cross-references from TASKS.md

**Context**: Designed encrypted scratchpad feature, wrote spec in specs/scratchpad.md, tasked it out in TASKS.md. Realized a new session picking up the tasks might never find the spec.

**Lesson**: Agents read TASKS.md early but may never discover specs/ on their own. Single-layer instructions get skipped under pressure; redundancy across layers is the only reliable mitigation for probabilistic instruction-following.

**Application**: Three-layer defense for every spec: (1) playbook instruction for the general pattern, (2) spec reference in the Phase header, (3) bold breadcrumb in the first task of the phase. Added 'Planning Non-Trivial Work' section to AGENT_PLAYBOOK.md to codify this.

---

## [2026-02-12] Git worktrees for parallel agent development

**Context**: Explored using git worktrees for parallel agent development across a large task backlog

**Lesson**: Git worktrees enable parallel Claude Code agent sessions without file conflicts. Create worktrees OUTSIDE the project as sibling directories (`git worktree add ../ctx-docs -b work/docs`). Each worktree gets its own branch, staging area, and working files but shares the same `.git` object database. Group tasks by blast radius (files touched) to minimize merge conflicts. 3-4 parallel worktrees is the practical limit before merge complexity outweighs productivity gains.

**Application**: When tackling many independent tasks: (1) group by file overlap, (2) create worktrees as siblings with `git worktree add ../name -b work/name`, (3) launch separate claude sessions in each, (4) merge back to main as tracks complete, (5) cleanup with `git worktree remove`. Don't run `ctx init` in worktrees — `.context/` is already tracked.

---

## [2026-02-11-124635] Chain-of-thought prompting improves agent reasoning accuracy

**Context**: Research shows accuracy on reasoning tasks jumps from 17.7% to 78.7% by adding think step-by-step to prompts. Applied this across agent guidelines.

**Lesson**: Explicit think step-by-step instructions in agent prompts dramatically improve reasoning accuracy at negligible token cost. This applies to skill files, playbooks, and autonomous loop prompts — anywhere the agent makes decisions before acting.

**Application**: Added Reason Before Acting section to AGENT_PLAYBOOK.md and reasoning nudges to 7 skills (ctx-implement, brainstorm, ctx-reflect, ctx-loop, qa, verify, consolidate). For autonomous loops, include reasoning instructions in PROMPT.md.

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

## [2026-01-30-120009] Say 'project conventions' not 'idiomatic X'

**Context**: When asking Claude to follow documentation style, saying 
'idiomatic Go' triggered training priors (stdlib conventions) instead of 
project-specific standards.

**Lesson**: Use 'follow project conventions' or 'check AGENT_PLAYBOOK' rather 
than 'idiomatic [language]' to ensure Claude looks at project files first.

**Application**: In prompts requesting style alignment, reference project 
files explicitly rather than language-wide conventions.

---

## [2026-01-29-164322] Documentation audits require verification against actual standards

**Context**: Agent claimed 'no Go docstring issues found' but manual inspection 
revealed many functions missing Parameters/Returns sections. The agent only 
checked if comments existed, not if they followed the standard format.

**Lesson**: When auditing documentation, compare against a known-good example 
first. Pattern-match for the COMPLETE standard (e.g., '// Parameters:' 
AND '// Returns:' sections), not just presence of any comment.

**Application**: Before declaring 'no issues', manually verify at least 5 
random samples match the documented standard. Use grep patterns that detect 
missing sections, not just missing comments.

---

## [2026-01-28-191951] Required flags now enforced for learnings

**Context**: Implemented ctx add learning flags to match decision's ADR 
(Architectural Decision Record) pattern

**Lesson**: Structured entries with Context/Lesson/Application are more useful
than one-liners

**Application**: Always use ctx add learning with all three flags; agents
guided via AGENT_PLAYBOOK.md

---

## [2026-01-28-072838] Changelogs vs Blogs serve different audiences

**Context**: Synthesizing session history into documentation

**Lesson**: Changelogs document WHAT; blogs explain WHY. Same information,
different engagement. Changelogs are for machines (audits, dependency trackers).
Blogs are for humans (narrative, context, lessons).

**Application**: When synthesizing session history, output both: changelog for
completeness, blog for readability.

---

## [2026-01-28-051426] IDE is already the UI

**Context**: Considering whether to build custom UI for .context/ files

**Lesson**: Discovery, search, and editing of .context/ markdown files works
better in VS Code/IDE than any custom UI we'd build. Full-text search,
git integration, extensions - all free.

**Application**: Don't reinvent the editor. Let users use their preferred IDE.

---

## [2026-01-28-040915] Subtasks complete does not mean parent task complete

**Context**: AI marked parent task done after finishing subtasks but missing
actual deliverable

**Lesson**: Subtask completion is implementation progress, not delivery.
The parent task defines what the user gets.

**Application**: Parent tasks should have explicit deliverables; don't close
until deliverable is verified.

---

## [2026-01-27-180000] Always Complete Decision Record Sections

**Context**: Decisions added via `ctx add decision` were left with placeholder
text like "[Add context here]".

**Lesson**: When recording decisions, always fill in Context
(what prompted this), Rationale (why this choice over alternatives), and
Consequences (what changes as a result). Placeholder text is a code smell -
decisions without rationale lose their value over time.

**Application**: After using `ctx add decision`, immediately edit the file to
complete all sections. Future: use `--context`, `--rationale`, `--consequences`
flags when available.

---

## [2026-01-27-160000] Slash Commands Require Matching Permissions

**Context**: Claude Code slash commands using `!` bash syntax require matching
permissions in settings.local.json.

**Lesson**: When adding new /ctx-* commands, ensure ctx init pre-seeds the
required `Bash(ctx <subcommand>:*)` permissions. Use additive merging for user
config - never remove existing permissions.

---

## [2026-01-26-180000] Go json.Marshal Escapes Shell Characters

**Context**: Generated settings.local.json had `2\u003e/dev/null` instead
of `2>/dev/null`.

**Lesson**: Go's json.Marshal escapes `>`, `<`, and `&` as
unicode (`\u003e`, `\u003c`, `\u0026`) by default for HTML safety.
Use `json.Encoder` with `SetEscapeHTML(false)` when generating config files
that contain shell commands.

---

## [2026-01-25-180000] CI Tests Need CTX_SKIP_PATH_CHECK

**Context**: CI tests were failing because ctx binary isn't installed on CI runners.

**Lesson**: Tests that call `ctx init` will fail without `CTX_SKIP_PATH_CHECK=1`
env var, because init checks if ctx is in PATH.

---

## [2026-01-25-170000] AGENTS.md Is Not Auto-Loaded

**Context**: Had both AGENTS.md and CLAUDE.md in project root, causing confusion.

**Lesson**: Only CLAUDE.md is read automatically by Claude Code. Projects
using ctx should rely on the CLAUDE.md → AGENT_PLAYBOOK.md chain, not AGENTS.md.

---

## [2026-01-25-140000] Autonomous Mode Creates Technical Debt

**Context**: Compared commits from autonomous "YOLO mode" (auto-accept,
agent-driven) vs human-guided refactoring sessions.

**Lesson**: YOLO mode is effective for feature velocity but accumulates technical debt:

| YOLO Pattern                           | Human-Guided Fix                      |
|----------------------------------------|---------------------------------------|
| `"TASKS.md"` scattered in 10 files     | `config.FilenameTask` constant        |
| `dir + "/" + file`                     | `filepath.Join(dir, file)`            |
| `{"task": "TASKS.md"}`                 | `{UpdateTypeTask: FilenameTask}`      |
| Monolithic `cli_test.go` (1500+ lines) | Colocated `package/package_test.go`   |
| `package initcmd` in `init/` folder    | `package initialize` in `initialize/` |

**Application**:
1. Schedule periodic consolidation sessions (not just feature sprints)
2. When same literal appears 3+ times, extract to constant
3. Constants should reference constants (self-referential maps)
4. Tests belong next to implementations, not in monoliths

---

## [2026-01-23-180000] ctx agent vs Manual File Reading Trade-offs

**Context**: User asked "Do you remember?" and agent used parallel file reads
instead of `ctx agent`. Compared outputs to understand the delta.

**Lesson**: `ctx agent` is optimized for task execution:
- Filters to pending tasks only
- Surfaces constitution rules inline
- Provides prioritized read order
- Token-budget aware

Manual file reading is better for exploratory/memory questions:
- Session history access
- Timestamps ("modified 8 min ago")
- Completed task context
- Parallel reads for speed

**Application**: No need to mandate one approach. Agents naturally pick appropriately:
- "Do you remember?" → parallel file reads (need history)
- "What should I work on?" → `ctx agent` (need tasks)

---

## [2026-01-23-140000] Infer Intent on "Do You Remember?" Questions

**Context**: User asked "Do you remember?" at session start. Agent asked for
clarification instead of proactively checking context files.

**Lesson**: In a ctx-enabled project, "do you remember?" has an obvious
meaning: check the `.context/` files and report what you know from previous
sessions. Don't ask for clarification - just do it.

**Application**: When user asks memory-related questions ("do you remember?",
"what were we working on?", "where did we leave off?"), immediately:
1. Read `.context/TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`
2. Run `ctx recall list --limit 5` for recent session history
3. Summarize what you find

Don't ask "would you like me to check the context files?" - that's the
obvious intent.

---

## [2026-01-23-120000] Always Use ctx from PATH

**Context**: Agent used `./dist/ctx-linux-arm64` and `go run ./cmd/ctx`
instead of just `ctx`, even though the binary was installed to PATH.

**Lesson**: When working on a ctx-enabled project, always use `ctx` directly:
```bash
ctx status        # correct
ctx agent         # correct
./dist/ctx        # avoid hardcoded paths
go run ./cmd/ctx  # avoid unless developing ctx itself
```

**Application**: Check `which ctx` if unsure. The binary is installed during
setup (`sudo make install` or `sudo cp ./ctx /usr/local/bin/`).

---

## [2026-01-21-180000] Exit Criteria Must Include Verification

**Context**: Dogfooding experiment had another Claude rebuild `ctx` from specs.
All tasks were marked complete, Ralph Loop exited successfully. But the built
binary didn't work — commands just printed help text instead of executing.

**Lesson**: "All tasks checked off" ≠ "Implementation works." This applies to
US too, not just the dogfooding clone. Our own verification is based on manual
testing, not automated proof. Blind spots exist in both projects.

Exit criteria must include:
- **Integration tests**: Binary executes commands correctly (not just unit tests)
- **Coverage targets**: Quantifiable proof that code paths are tested
- **Smoke tests**: Basic "does it run" verification in CI

**Application**:
1. Add integration test suite that invokes the actual binary
2. Set coverage targets (e.g., 70% for core packages)
3. Add verification tasks to TASKS.md — we have the same blind spot
4. Being proud of our achievement doesn't prove its validity

---

## [2026-01-21-160000] Orchestrator vs Agent Tasks Must Be Separate

**Context**: Ralph Loop checked `IMPLEMENTATION_PLAN.md`, found all tasks
done, exited — ignoring `.context/TASKS.md`.

**Lesson**: Separate concerns:
- **`IMPLEMENTATION_PLAN.md`** = Orchestrator directive ("check your tasks")
- **`.context/TASKS.md`** = Agent's mind (actual task list)

The orchestrator shouldn't maintain a parallel ledger. It just says
"check your mind."

**Application**: For new projects, `IMPLEMENTATION_PLAN.md` has ONE task:
"Check `.context/TASKS.md`"

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

## [2026-01-21-120000] Hooks Should Use PATH, Not Hardcoded Paths

**Context**: Original hooks used hardcoded absolute paths like
`/home/user/project/dist/ctx-linux-arm64`. This caused issues when dogfooding
or sharing configs.

**Lesson**: Hooks should assume `ctx` is in the user's PATH:
- More portable across machines/users
- Standard Unix practice
- `ctx init` now checks if `ctx` is in PATH before proceeding
- Hooks use `ctx agent` instead of `/full/path/to/ctx-linux-arm64 agent`

**Application**:
1. Users must install ctx to PATH: `sudo make install` or `sudo cp ./ctx /usr/local/bin/`
2. `ctx init` will fail with clear instructions if ctx is not in PATH
3. Tests can skip this check with `CTX_SKIP_PATH_CHECK=1`

**Supersedes**: Previous learning "Binary Path Must Be Absolute" (2026-01-20)

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

## [2026-01-20-180000] .context/ Is NOT a Claude Code Primitive

**Context**: User asked if Claude Code natively understands `.context/`.

**Lesson**: Claude Code only natively reads:
- `CLAUDE.md` (auto-loaded at session start)
- `.claude/settings.json` (hooks and permissions)

The `.context/` directory is a ctx convention. Claude won't know about it unless:
1. A hook runs `ctx agent` to inject context
2. CLAUDE.md explicitly instructs reading `.context/`

**Application**: Always create CLAUDE.md as the bootstrap entry point.

---

## [2026-01-20-080000] Always Backup Before Modifying User Files

**Context**: `ctx init` needs to create/modify CLAUDE.md, but user may have existing customizations.

**Lesson**: When modifying user files (especially config files like CLAUDE.md):
1. **Always backup first** — `file.bak` before any modification
2. **Check for existing content** — use marker comments for idempotency
3. **Offer merge, don't overwrite** — respect user's customizations
4. **Provide escape hatch** — `--merge` flag for automation, manual merge for control

**Application**: Any `ctx` command that modifies user files should follow this pattern.

---

## [2026-01-19-120000] CGO Must Be Disabled for ARM64 Linux

**Context**: `go test` failed with `gcc: error: unrecognized command-line option '-m64'`

**Lesson**: On ARM64 Linux, CGO causes cross-compilation issues. Always use `CGO_ENABLED=0`.

**Application**:
```bash
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
CGO_ENABLED=0 go test ./...
```

---
