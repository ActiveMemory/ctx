# Consolidated Learnings Archive — 2026-02-22

Original entries removed during consolidation.

---

## [2026-02-20-142441] Own the codebase: fix pre-existing lint issues

**Context**: make lint failed on pre-existing SA9003 and goconst violations. We initially committed without running make lint.

**Lesson**: A good agent owns the full codebase — if lint is broken, fix it regardless of who introduced it. Always run make lint, not just targeted linting, before committing.

**Application**: Run make lint as part of every pre-commit check and fix any issues found even if pre-existing.

---

## [2026-02-20-121941] details tag cannot wrap pre blocks with blank lines

**Context**: Considered making tool output collapsible using <details><summary>...</summary><pre><code>...</code></pre></details>. The content inside contains blank lines from tool output.

**Lesson**: <details> is a CommonMark Type 6 HTML block — it ends at the first blank line, regardless of what is inside. A <pre> (Type 1) inside <details> does not override this because CommonMark does not nest HTML block types. Any blank line in the content terminates the <details> block, orphaning the closing tags.

**Application**: Collapsible tool output in the journal site requires a CSS/JS approach (e.g., a toggle class on the <pre> block) rather than native <details> elements. This is a future feature, not a current blocker.

---

## [2026-02-20-121939] pymdownx.highlight hijacks pre/code blocks

**Context**: After switching journal site Tool Output and User turns to <pre><code> with HTML escaping, the rendered output still swallowed subsequent turns. The pymdownx.highlight extension was intercepting <pre><code> patterns and transforming them into fancy code blocks with line numbers and copy buttons, changing block boundaries.

**Lesson**: MkDocs/zensical's pymdownx.highlight extension pattern-matches on <pre><code> and transforms it into a widget with spans, nav buttons, and line numbers. This transformation changes the HTML block structure. Disable with use_pygments=false in markdown_extensions config. Bare <pre> (without <code>) also avoids the hijacking but loses semantic markup.

**Application**: When generating <pre><code> blocks for any zensical/MkDocs site, always configure pymdownx.highlight with use_pygments=false. When overriding markdown_extensions in zensical, the entire default set must be replicated since providing the key replaces all defaults.

---

## [2026-02-20-110254] Journal site normalizeContent: three-layer tool output fix

**Context**: Journal site rendering broke for fencesVerified files: (1) <details>/<pre> wrappers from old export passed through unchanged, (2) inner fence markers from embedded content caused nesting conflicts, (3) HTML unescaping ran incorrectly on collapseToolOutputs format which never HTML-escapes

**Lesson**: Three fixes work together: (a) always run wrapToolOutputs regardless of fencesVerified — use isAlreadyFenced to detect content already in a fenced block and skip re-wrapping; (b) fenceForContent picks a fence depth exceeding any inner fences (4 backticks if content has 3); (c) stripPreWrapper only unescapes HTML when <pre> was specifically found, not just when <details> wrapper was stripped. normalizeContent fence tracking upgraded to CommonMark-compliant variable-length matching.

**Application**: When generating site copies, tool output always needs conversion from source format (details/pre or raw) to fenced code blocks. The fencesVerified flag only affects stripFences, never wrapToolOutputs. If adding new wrapper formats, ensure stripPreWrapper and isAlreadyFenced handle them.

---

## [2026-02-20-110243] Tool output boundary detection: pre-scan + last-match-wins

**Context**: wrapToolOutputs boundary detection was fooled by embedded turn headers from other journal files inside tool output content (e.g., reading Part 5 of another session showed ### 802. Assistant inside turn 41's body, causing the real ### 42. to be swallowed)

**Lesson**: Pre-scan all turn numbers, sort+dedup, find min > N as the expected next turn, then use the LAST positional occurrence of that number as the boundary. Last-match-wins handles duplicate turn numbers (embedded copy appears first, real one appears after closing tags). No magic constants or wrapper tracking needed — just the sorted sequence and positional ordering.

**Application**: Apply this pattern whenever parsing structured content that may contain recursive/embedded copies of itself. The boundary target should be derived from the document's global structure (sorted turn set), not local heuristics (gap limits, wrapper tracking). collapseToolOutputs in collapse.go has similar boundary detection that may benefit from the same approach.

---

## [2026-02-20-070902] Pre-commit gate: build + lint + test, every time

**Context**: Repeatedly shipped code that passed tests but would have bounced from CI due to linter violations. User had to explicitly remind me to run golangci-lint.

**Lesson**: Tests verify behavior; the linter enforces style, security (gosec), and static analysis that CI gates on. Running only 'go test' is an incomplete check. The full pre-commit gate is: go build, golangci-lint run, go test — all three, every time, before calling code done.

**Application**: Before any code is considered commit-ready, always run: (1) CGO_ENABLED=0 go build ./cmd/ctx (2) golangci-lint run ./path/to/package/ (3) CGO_ENABLED=0 go test ./path/to/package/. Never skip the linter. Never assume tests alone are sufficient.

---

## [2026-02-20-061907] AI normalization at scale hits context limits

**Context**: Attempted AI-driven normalization of 290 journal files (~1M lines total). Agents hit context limits on files over ~500 lines (median file was 3.4K lines).

**Lesson**: AI normalization is a poor fit for large files. The source files were already well-formatted — agents found almost no issues to fix. Code-level pipeline (normalizeContent) handles rendering concerns more reliably and consistently.

**Application**: For bulk journal processing, use code-level transforms at render time. Reserve AI normalization for specific files with known issues. Mark all as normalized since the pipeline handles it.

---

## [2026-02-20-044412] Inline code spans with angle brackets break markdown rendering

**Context**: Journal entry body content discussing XML fragments like backtick-less-than-slash-com introduced broken HTML into the rendered page because the angle brackets inside backticks were interpreted as raw HTML tags.

**Lesson**: Single-line backtick spans containing angle brackets need special handling: replace backticks with double-quotes (preserves visual signal) and replace angle brackets with HTML entities. This is done via RegExInlineCodeAngle regex in the normalizeContent line-by-line pass. Multi-line or angle-bracket-free spans are left untouched.

**Application**: The regex pattern matches single backtick on one line containing < or >. Applied after fence stripping but in the line-by-line pass, not in wrapToolOutputs (which handles entire Tool Output sections).

---

## [2026-02-20-044403] Journal title sanitization requires multiple passes

**Context**: Link text in journal index.md broke rendering when titles contained angle brackets (from truncated XML tags like command-message), backticks, or hash characters. Titles over 75 chars wrapped to a second line and lost heading formatting.

**Lesson**: Title sanitization pipeline: 1) Strip Claude Code XML tags (command-message, command-name, local-command-caveat) via RegExClaudeTag. 2) Replace angle brackets with HTML entities. 3) Strip backticks and hash (meaningless in link text). 4) Truncate to 75 chars on word boundary. This applies at both parse time (parseJournalEntry in parse.go) and export time (cleanTitle in recall/slug.go). The H1 heading in normalizeContent also strips Claude tags and truncates to 75 chars.

**Application**: When adding new title sources or display contexts, ensure the full sanitization chain applies. RecallMaxTitleLen (75) is the single source of truth for title length. RegExClaudeTag lives in config/regex.go for sharing between journal and recall packages.

---

## [2026-02-20-044352] Python-Markdown HTML blocks end at blank lines unlike CommonMark

**Context**: Debugging journal site rendering: tool output content with blank lines, headings, thematic breaks, and lists was being interpreted as markdown even inside pre/code and details/pre wrappers.

**Lesson**: Python-Markdown (used by mkdocs/zensical) ends ALL HTML blocks at blank lines, regardless of tag type. CommonMark has Type 1 blocks (pre) that survive blank lines, but Python-Markdown does not. html.EscapeString only handles angle brackets, ampersand, quotes — markdown syntax (hash, dashes, asterisk, numbered lists) passes through untouched. The only reliable way to prevent markdown interpretation of arbitrary content is fenced code blocks, which survive blank lines and block all markdown/HTML parsing.

**Application**: For journal site tool output wrapping: always use fenced code blocks. Run stripFences before wrapToolOutputs so content has no fence lines, making triple-backtick safe as a wrapper. For overflow control (replacing details collapsibility), use CSS max-height + overflow-y: auto on pre elements.

---

## [2026-02-17] Hook grep patterns match inside quoted arguments — use specific anchors

**Context**: Added `(cp|install|mv)\s.*/bin` to block-dangerous-commands.sh. It matched "install" inside `ctx add learning "...install...*/bin/..."` quoted text, blocking legitimate commands.

**Lesson**: Shell hook grep patterns operate on the full command string and cannot distinguish between command names and text inside quoted arguments. Generic patterns like `install\s.*/bin` are too broad. Use specific directory lists and anchor to command-start positions to reduce false positives.

**Application**: When writing hook patterns: (1) list specific dangerous destinations instead of generic `/bin`, (2) anchor with `(^|;|&&|\|\|)\s*` to match command position, (3) test with `ctx add learning` containing the blocked words to verify no false positives.

---

## [2026-02-16-100442] gosec G301/G306: use 0o750 for dirs, 0o600 for files in test code too

**Context**: Plugin conversion: test files used 0o755 and 0o644 which triggered gosec warnings

**Lesson**: gosec checks ALL code including tests. Test helper MkdirAll and WriteFile calls need the same restrictive permissions as production code.

**Application**: Use 0o750 for os.MkdirAll and 0o600 for os.WriteFile everywhere, including test setup code.

---

## [2026-02-16-100438] golangci-lint errcheck: use cmd.Printf not fmt.Fprintf in Cobra commands

**Context**: Plugin conversion: permissions/run.go had 7 errcheck failures from fmt.Fprintf(cmd.OutOrStdout(), ...)

**Lesson**: Cobra's cmd.Printf/cmd.Println write to OutOrStdout() without returning errors, avoiding errcheck lint. fmt.Fprintf returns (int, error) that must be handled.

**Application**: Always use cmd.Printf/cmd.Println for Cobra command output. Reserve fmt.Fprintf for non-Cobra io.Writer contexts.

---

## [2026-02-15-194827] Hook scripts can lose execute permission without warning

**Context**: Every Bash call showed PreToolUse:Bash hook error. Commands succeeded but UX was degraded.

**Lesson**: block-non-path-ctx.sh had -rw-r--r-- instead of -rwxr-xr-x. Claude Code reports non-executable hooks as 'hook error' but still runs the command.

**Application**: After editing or regenerating hook scripts, verify permissions with ls -la .claude/hooks/*.sh. Consider adding a chmod +x step to ctx init or a drift check for hook permissions.

---

## [2026-02-15-170015] Two-tier hook output is sufficient — don't over-engineer severity levels

**Context**: Evaluated whether ctx hooks need a formal INFO/WARN/CRITICAL severity protocol (Pattern 8 in hook-output-patterns.md). Reviewed all shipped hooks: block-non-path-ctx (hard gate), check-context-size (VERBATIM relay), check-persistence (unprefixed nudge), check-journal (VERBATIM + suggested action), check-backup-age (VERBATIM + suggested action), cleanup-tmp (silent side-effect).

**Lesson**: ctx already has a working two-tier system: unprefixed output (agent absorbs as context, mentions if relevant — e.g. check-persistence.sh) and 'IMPORTANT: Relay VERBATIM' prefixed output (agent interrupts immediately — e.g. check-context-size.sh). A three-tier system adds protocol complexity (agent training in CLAUDE.md, consistent prefix usage) without covering cases the two tiers don't already handle.

**Application**: When writing new hooks, choose between silent (no output), unprefixed (agent context — may or may not relay), and VERBATIM (guaranteed relay). Don't introduce new severity prefixes unless the two-tier model demonstrably fails for a specific hook.

---

## [2026-02-15-105918] Gitignored folders accumulate stale artifacts

**Context**: Found a published blog draft still sitting in the gitignored ideas/ folder

**Lesson**: Gitignored directories are invisible to git status, so stale files persist indefinitely. Published drafts, old reports, and resolved spikes linger because nothing flags them.

**Application**: Periodically ls gitignored working directories (ideas/, dist/, etc.) and clean up artifacts that have been promoted or are no longer relevant.

---

## [2026-02-15-105914] Editor artifacts need gitignore coverage from day one

**Context**: Found .swp files showing as untracked — vim swap files were not in .gitignore

**Lesson**: The default Go .gitignore template covers .idea/ and .vscode/ but not vim artifacts (*.swp, *.swo, *~). These accumulate silently.

**Application**: When setting up a new project, add *.swp, *.swo, *~ to .gitignore alongside IDE directories.

---

## [2026-02-15-044503] Permission drift needs auditing like code drift

**Context**: settings.local.json is gitignored so it drifts independently — no PR review, no CI check catches stale or missing permissions

**Lesson**: Permission drift is a distinct category from code or context drift. Skills get added/removed but their Skill() entries in settings.local.json lag behind. The /ctx-drift skill now checks for this.

**Application**: Run /ctx-drift periodically to catch: missing Bash(ctx:*), missing Skill(ctx-*) for installed skills, stale Skill(ctx-*) for removed skills, granular entries that should be consolidated.

---

## [2026-02-15-044500] Skill() permissions do not support name prefix globs

**Context**: Tried to use Skill(ctx-*) to cover all ctx skills in settings.local.json

**Lesson**: Claude Code Skill() permission wildcards only match arguments (e.g., Skill(commit *)), not skill name prefixes. Skill(ctx-*) will not match ctx-add-learning, ctx-agent, etc.

**Application**: List each Skill(ctx-*) entry individually in DefaultClaudePermissions and settings.local.json. When adding a new ctx-* skill, add its Skill() entry to both places.

---

## [2026-02-15-044457] Wildcard trusted binaries, keep git granular

**Context**: Consolidated 22 ctx entries into Bash(ctx:*) and 6 make entries into Bash(make:*), but kept git commands individual

**Lesson**: Trusted binaries (your own CLI, make) should use a single Bash(cmd:*) wildcard. Git needs per-command entries because safe (git log) and destructive (git reset --hard) commands share the same binary and hooks don't block all destructive git operations.

**Application**: Use Bash(ctx:*) and Bash(make:*) wildcards. List git commands individually: git add, git branch, git commit, git diff, git log, git remote, git restore, git show, git stash, git status, git tag. Never wildcard Bash(git:*).

---

## [2026-02-15-044453] settings.local.json accumulates session debris

**Context**: Audited settings.local.json and found 24 removable entries out of 90 — garbage, one-offs, subsumed patterns, stale references

**Lesson**: Every Allow click appends an entry. Over time: hardcoded paths, literal arguments, duplicate intent (env var ordering), garbage entries, and stale skill references accumulate. Invisible drift because the file is gitignored.

**Application**: Run periodic permission hygiene using hack/runbooks/sanitize-permissions.md runbook. Use /ctx-drift to detect permission drift (missing skills, stale entries, consolidation opportunities).

---

## [2026-02-15-044450] Skill vs runbook for agent self-modification

**Context**: Considered building a skill to clean up settings.local.json permissions

**Lesson**: When a skill would edit files that control agent behavior (permissions, hooks, instructions), a runbook is safer. Auto-accept makes self-modifying skills an escalation vector.

**Application**: Use runbooks (human edits, agent advises) for operations on .claude/settings.local.json, CLAUDE.md, hooks, and CONSTITUTION.md. Reserve skills for operations where agent autonomy is safe.

---

## [2026-02-15-034225] G304 gosec false positives in test files are safe to suppress

**Context**: gosec flags os.ReadFile with variable paths as G304 (potential file inclusion), even in test files where paths come from t.TempDir() and compile-time constants

**Lesson**: G304 requires user-controlled input to be exploitable. Test files using t.TempDir() and constants have no attack vector. Suppress with //nolint:gosec // test file path

**Application**: When gosec raises G304 in test files, verify paths aren't from external input, then suppress with nolint comment rather than restructuring the code

---

## [2026-02-14-163855] Skills can replace CLI commands for interactive workflows

**Context**: Evaluating whether /absorb (formerly /ctx-borrow) needed a full CLI command or if the skill was sufficient

**Lesson**: A well-structured skill recipe is a guide, not a rigid script. The agent improvises beyond literal instructions and adapts to edge cases using its available tools.

**Application**: Prefer skills over CLI commands when the workflow requires judgment calls (conflict resolution, selective application, strategy selection). Reserve CLI commands for deterministic, non-interactive operations.

---

## [2026-02-12-005911] Claude Code UserPromptSubmit hooks: stderr with exit 0 is swallowed (only visible in verbose mode Ctrl+O). stdout with exit 0 is prepended as context for the AI. For user-visible warnings use systemMessage JSON on stdout. For AI-facing nudges use plain text on stdout. There is no non-blocking stderr channel for this hook type.

**Context**: All three UserPromptSubmit hooks (check-context-size, check-persistence, prompt-coach) were outputting to stderr, making their output invisible to both user and AI

**Lesson**: stderr from UserPromptSubmit hooks is invisible. Use stdout for AI context, systemMessage JSON for user-visible warnings.

**Application**: AI-facing hooks: drop >&2 redirects. User-facing hooks: output {"systemMessage": "..."} JSON to stdout.

---

## [2026-02-12-005510] Prompt-coach hook outputs to stdout (UserPromptSubmit) which is prepended as AI context, not shown to the user. stderr with exit 0 is swallowed entirely. The only user-visible options are systemMessage JSON (warning banner) or exit 2 (blocks the prompt). There is no non-blocking user-visible output channel for UserPromptSubmit hooks.

**Context**: Debugging why prompt-coach tips were invisible to the user despite firing correctly

**Lesson**: UserPromptSubmit hook stdout goes to the AI as context, not the user terminal. stderr with exit 0 is invisible. No non-blocking user-facing output channel exists for this hook type.

**Application**: Design hooks for their actual audience: AI-facing hooks use stdout, user-facing feedback needs systemMessage or a different mechanism entirely.

---

## [2026-02-11-195405] Gitignore rules for sensitive directories must survive cleanup sweeps

> **Superseded** (2026-02-17): `.context/sessions/` was fully removed in v0.4.0. The directory is no longer created, referenced, or used by any code path. The gitignore entry was removed as dead weight during the issue #7 cleanup. The general principle (audit before removing security controls) remains sound, but no longer applies to sessions.

**Context**: During a stale-reference sweep, the .context/sessions/ gitignore rule was removed because sessions were consolidated into journals. But the gitignore rule exists to prevent sensitive data from being committed, not to document architecture. The directory may still exist locally.

**Lesson**: Gitignore entries for sensitive paths are security controls, not documentation. Never remove them during doc/reference cleanups even if the feature they relate to was removed.

**Application**: Before removing any gitignore entry, ask: does this entry exist for security/privacy or for architecture? Security entries stay permanently.

---

## [2026-02-07-014920] Agent ignores repeated hook output (repetition fatigue)

**Context**: PreToolUse hook ran ctx agent on every tool use, injecting the same
context packet repeatedly. Agent tuned it out and didn't follow conventions.

**Lesson**: Repeated injection causes the agent to ignore the output. A cooldown 
tombstone (--session $PPID --cooldown 10m) emits once per window. A readback 
instruction (confirm to user you read context) creates a behavioral gate harder 
to skip than silent injection.

**Application**: Use --session $PPID in hook commands to enable cooldown. Pair 
context injection with a readback instruction so the agent must acknowledge 
before starting work.

---

## [2026-02-05-174304] Use $CLAUDE_PROJECT_DIR in hook paths

**Context**: Migrating hooks after username rename (parallels→jose) broke all 
absolute paths in settings.local.json

**Lesson**: Claude Code provides $CLAUDE_PROJECT_DIR env var for hook commands — 
resolves to project root at runtime, survives renames

**Application**: Always use "$CLAUDE_PROJECT_DIR"/.claude/hooks/... in 
settings.local.json, never hardcode /home/user/...

---

## [2026-02-04-230943] JSONL session files are append-only

**Context**: Built context-watch.sh monitor; it showed 90% after compaction 
while /context showed 16%

**Lesson**: Claude Code JSONL files never shrink after compaction. Any monitoring 
tool based on file size will overreport post-compaction. The /context command 
shows actual tokens sent to the model.

**Application**: Per ctx workflow, sessions should end before compaction fires — 
so JSONL size is a valid time-to-wrap-up signal. Don't try to make 
context-watch.sh compaction-aware.

---

## [2026-02-04-230941] Most external skill files are redundant with Claude's system prompt

**Context**: Reviewed ~30 external skill/prompt files during systematic skill audit

**Lesson**: Only ~20% had salvageable content — and even those yielded just a few 
heuristics each. The signal is in the knowledge delta, not the word count.

**Application**: When evaluating new skills, apply E/A/R classification ruthlessly. 
Default to delete. Only keep content an expert would say took years to learn.

---

## [2026-02-04-193920] Skills that restate or contradict Claude Code's built-in system prompt create tension, not clarity

**Context**: Reviewing entropy.txt skill that duplicated system prompt guidance 
about code minimalism

**Lesson**: Skills that conflict with system prompts cause unpredictable behavior — 
the AI has to reconcile contradictory instructions. The system prompt already 
covers: avoid over-engineering, don't add unnecessary features, prefer 
simplicity. Skills should complement the system prompt, not compete with it.

**Application**: When evaluating or writing skills, first check Claude Code's 
system prompt defaults. Only create skills for guidance the platform does NOT 
already provide.

---

## [2026-02-04-192812] Skill files that suppress AI judgment are jailbreak patterns, not productivity tools

**Context**: Reviewing power.txt skill that forced skill invocation on every message

**Lesson**: Red flags: <EXTREMELY-IMPORTANT> urgency tags, 'you cannot rationalize' 
overrides, tables that label hesitation as wrong, absurdly low thresholds (1%). 
The fix for 'AI forgets skills' is better skill descriptions, not overriding 
reasoning. Discard these entirely — nothing is salvageable.

**Application**: When evaluating skills, check for judgment-suppression 
patterns before assessing content.

---

## [2026-02-03-064236] Claude Code subagent sessions share parent sessionId

**Context**: After fixing the slug issue, sessions still showed wrong content 
(SUGGESTION MODE instead of actual conversation). Investigation revealed 
subagent files in /subagents/ directories use the same sessionId as the parent.

**Lesson**: Subagent files (e.g., prompt_suggestion, compact) share the parent 
sessionId. When scanning directories, subagent sessions can appear 'newer' 
(later timestamp) and win during deduplication, causing main session content 
to be lost.

**Application**: Skip /subagents/ directories when scanning for sessions. 
Use filepath.SkipDir for efficiency. Subagent sessions have isSidechain:true 
and an agentId field.

---

## [2026-02-03-063337] Claude Code JSONL format changed: slug field removed in v2.1.29+

**Context**: ctx recall export --all --force was skipping February 2026 sessions. 
Investigation revealed sessions like c9f12373 had 0 slug fields but 19 
sessionId fields.

**Lesson**: Claude Code removed the 'slug' field from message records in newer 
versions. The parser's CanParse function required both sessionId AND slug, 
causing it to reject valid session files.

**Application**: When parsing Claude Code sessions, check for sessionId and 
valid type (user/assistant) instead of requiring slug. The slug may be 
available in sessions-index.json if needed.

---

## [2026-01-28-194113] Claude Code Hooks Receive JSON via Stdin

**Context**: Debugging Claude Code PreToolUse hooks - they were not receiving
command data when using environment variables like CLAUDE_TOOL_INPUT

**Lesson**: Claude Code hooks receive input as JSON via stdin, not environment
variables. Use HOOK_INPUT=$(cat) then parse with
jq: COMMAND=$(echo "$HOOK_INPUT" | jq -r ".tool_input.command // empty")

**Application**: All hook scripts should read stdin for input. The JSON
structure includes .tool_input.command for Bash commands. Test hooks with
debug logging to /tmp/ to verify they receive expected data.

---

## [2026-01-28-040251] AI session JSONL formats are not standardized

**Context**: Building recall feature to parse session history from multiple
AI tools

**Lesson**: Claude Code, Cursor, Aider each have different JSONL formats
or may not export sessions at all.

**Application**: Use tool-agnostic Session type with tool-specific parsers.

---

## [2026-01-26-160000] Claude Code Hook Key Names

**Context**: Hooks weren't working, getting "Invalid key in record" errors.

**Lesson**: Claude Code settings.local.json hook keys are `PreToolUse` and
`SessionEnd` (not `PreToolUseHooks`/`SessionEndHooks`). The `Hooks` suffix
causes validation errors.

---

## [2026-01-25-200000] defer os.Chdir Fails errcheck Linter

**Context**: `defer os.Chdir(originalDir)` fails golangci-lint errcheck.

**Lesson**: Use `defer func() { _ = os.Chdir(x) }()` to explicitly ignore the
error return value.

---

## [2026-01-25-190000] golangci-lint Go Version Mismatch in CI

**Context**: CI was failing with Go version mismatches between golangci-lint
and the project.

**Lesson**: When golangci-lint is built with an older Go version than the
project targets, use `install-mode: goinstall` in CI to build the linter from
source using the project's Go version.

---

## [2026-01-25-160000] Hook Regex Can Overfit

**Context**: `.claude/hooks/block-non-path-ctx.sh` was blocking legitimate sed
commands because the regex `ctx[^ ]*` matched paths containing "ctx" as a
directory component (e.g., `/home/user/ctx/internal/...`).

**Lesson**: When writing shell hook regexes:
- Test against paths that contain the target string as a substring
- `ctx` as binary vs `ctx` as directory name are different
- Original: `(/home/|/tmp/|/var/)[^ ]*ctx[^ ]* ` — overfits
- Fixed: `(/home/|/tmp/|/var/)[^ ]*/ctx( |$)` — matches binary only

**Application**: Always test hooks with edge cases before deploying.

---

## [2026-01-23-160000] Claude Code Skills Format

**Context**: Needed to understand how to create custom slash commands.

**Lesson**: Claude Code skills are markdown files in `.claude/commands/` with
YAML frontmatter (`description`, `argument-hint`, `allowed-tools`). Body is
the prompt. Use code blocks with `!` prefix for shell execution. `$ARGUMENTS`
passes command args.

---

## [2026-01-20-160000] SessionEnd Hook Catches Ctrl+C

> **Note**: `.context/sessions/` removed in v0.4.0. The auto-save hook was eliminated. The SessionEnd hook behavior documented here is still accurate for Claude Code, but ctx no longer uses it.

**Context**: Needed to auto-save context even when user force-quits with Ctrl+C.

**Lesson**: Claude Code's `SessionEnd` hook fires on ALL exits including Ctrl+C. It provides:
- `transcript_path` - full session transcript (.jsonl)
- `reason` - why session ended (exit, clear, logout, etc.)
- `session_id` - unique session identifier

**Application**: SessionEnd hook is available for custom workflows but ctx no longer uses it for auto-save.

---

## [2026-01-20-140000] Session Filename Must Include Time

> **Note**: `.context/sessions/` removed in v0.4.0. This naming convention is no longer used by ctx.

**Context**: Using just date (`2026-01-20-topic.md`) would overwrite multiple sessions per day.

**Lesson**: Use `YYYY-MM-DD-HHMM-<topic>.md` format to prevent overwrites.

**Application**: Historical reference only. Journal entries now use `ctx recall export` naming.

---

## [2026-01-20-120000] Two Tiers of Persistence

> **Note**: `.context/sessions/` removed in v0.4.0. Two tiers remain but the full-dump tier is now `~/.claude/projects/` (raw JSONL) + `.context/journal/` (enriched markdown via `ctx recall export`).

**Context**: User wanted to ensure nothing is lost when session ends.

**Lesson**: Two levels serve different needs:

| Tier      | Content                         | Purpose                       | Location                      |
|-----------|---------------------------------|-------------------------------|-------------------------------|
| Curated   | Key learnings, decisions, tasks | Quick reload, token-efficient | `.context/*.md`               |
| Full dump | Entire conversation             | Safety net, deep dive         | `~/.claude/projects/` + `.context/journal/` |

**Application**: Before session ends, persist learnings and decisions via `/ctx-reflect`. Full transcripts are retained automatically by Claude Code.

---

## [2026-01-20-100000] Auto-Load Works, Auto-Save Was Missing

> **Note**: `.context/sessions/` removed in v0.4.0. The auto-save hook was eliminated. Claude Code retains transcripts in `~/.claude/projects/` automatically.

**Context**: Explored how to persist context across Claude Code sessions.

**Lesson**: Initial state was asymmetric:
- **Auto-load**: Works via `PreToolUse` hook running `ctx agent`
- **Auto-save**: Did NOT exist

**Original solution**: `SessionEnd` hook that copies transcript to `.context/sessions/`. Removed in v0.4.0 because Claude Code already retains transcripts and `ctx recall export` reads them directly.

---

