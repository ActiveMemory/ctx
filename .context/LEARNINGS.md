# Learnings

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-02-15 | Gitignored folders accumulate stale artifacts |
| 2026-02-15 | Editor artifacts need gitignore coverage from day one |
| 2026-02-15 | Permission drift needs auditing like code drift |
| 2026-02-15 | Skill() permissions do not support name prefix globs |
| 2026-02-15 | Wildcard trusted binaries, keep git granular |
| 2026-02-15 | settings.local.json accumulates session debris |
| 2026-02-15 | Skill vs runbook for agent self-modification |
| 2026-02-15 | Cross-repo links to published docs should use ctx.ist |
| 2026-02-15 | G304 gosec false positives in test files are safe to suppress |
| 2026-02-14 | ctx add learning/decision requires structured flags, not just a string |
| 2026-02-14 | ctx init is non-destructive toward tool-specific configs |
| 2026-02-14 | merge insertion is position-aware, not append |
| 2026-02-14 | ctx init CLAUDE.md handling is a 3-state machine |
| 2026-02-14 | Skills can replace CLI commands for interactive workflows |
| 2026-02-14 | color.NoColor in init for CLI test files |
| 2026-02-14 | Recall CLI tests isolate via HOME env var |
| 2026-02-14 | formatDuration accepts interface not time.Duration |
| 2026-02-14 | normalizeCodeFences regex splits language specifiers |
| 2026-02-13 | Specs get lost without cross-references from TASKS.md |
| 2026-02-12 | Claude Code UserPromptSubmit hooks: stderr with exit 0 is swallowed (only visible in verbose mode Ctrl+O). stdout with exit 0 is prepended as context for the AI. For user-visible warnings use systemMessage JSON on stdout. For AI-facing nudges use plain text on stdout. There is no non-blocking stderr channel for this hook type. |
| 2026-02-12 | Prompt-coach hook outputs to stdout (UserPromptSubmit) which is prepended as AI context, not shown to the user. stderr with exit 0 is swallowed entirely. The only user-visible options are systemMessage JSON (warning banner) or exit 2 (blocks the prompt). There is no non-blocking user-visible output channel for UserPromptSubmit hooks. |
| 2026-02-11 | Gitignore rules for sensitive directories must survive cleanup sweeps |
| 2026-02-11 | Chain-of-thought prompting improves agent reasoning accuracy |
| 2026-02-07 | Agent ignores repeated hook output (repetition fatigue) |
| 2026-02-06 | PROMPT.md deleted — was stale project briefing, not a Ralph loop prompt |
| 2026-02-05 | Use $CLAUDE_PROJECT_DIR in hook paths |
| 2026-02-04 | JSONL session files are append-only |
| 2026-02-04 | Most external skill files are redundant with Claude's system prompt |
| 2026-02-04 | Skills that restate or contradict Claude Code's built-in system prompt create tension, not clarity |
| 2026-02-04 | Skill files that suppress AI judgment are jailbreak patterns, not productivity tools |
| 2026-02-03 | User input often has inline code fences that break markdown rendering |
| 2026-02-03 | Claude Code injects system-reminder tags into tool results, breaking markdown export |
| 2026-02-03 | Claude Code subagent sessions share parent sessionId |
| 2026-02-03 | Claude Code JSONL format changed: slug field removed in v2.1.29+ |
| 2026-01-30 | Say 'project conventions' not 'idiomatic X' |
| 2026-01-29 | Documentation audits require verification against actual standards |
| 2026-01-28 | Required flags now enforced for learnings |
| 2026-01-28 | Claude Code Hooks Receive JSON via Stdin |
| 2026-01-28 | Changelogs vs Blogs serve different audiences |
| 2026-01-28 | IDE is already the UI |
| 2026-01-28 | Subtasks complete does not mean parent task complete |
| 2026-01-28 | AI session JSONL formats are not standardized |
| 2026-01-27 | Always Complete Decision Record Sections |
| 2026-01-27 | Slash Commands Require Matching Permissions |
| 2026-01-26 | Go json.Marshal Escapes Shell Characters |
| 2026-01-26 | Claude Code Hook Key Names |
| 2026-01-25 | defer os.Chdir Fails errcheck Linter |
| 2026-01-25 | golangci-lint Go Version Mismatch in CI |
| 2026-01-25 | CI Tests Need CTX_SKIP_PATH_CHECK |
| 2026-01-25 | AGENTS.md Is Not Auto-Loaded |
| 2026-01-25 | Hook Regex Can Overfit |
| 2026-01-25 | Autonomous Mode Creates Technical Debt |
| 2026-01-23 | ctx agent vs Manual File Reading Trade-offs |
| 2026-01-23 | Claude Code Skills Format |
| 2026-01-23 | Infer Intent on "Do You Remember?" Questions |
| 2026-01-23 | Always Use ctx from PATH |
| 2026-01-21 | Exit Criteria Must Include Verification |
| 2026-01-21 | Orchestrator vs Agent Tasks Must Be Separate |
| 2026-01-21 | One Templates Directory, Not Two |
| 2026-01-21 | Hooks Should Use PATH, Not Hardcoded Paths |
| 2026-01-20 | ctx and Ralph Loop Are Separate Systems |
| 2026-01-20 | .context/ Is NOT a Claude Code Primitive |
| 2026-01-20 | SessionEnd Hook Catches Ctrl+C |
| 2026-01-20 | Session Filename Must Include Time |
| 2026-01-20 | Two Tiers of Persistence |
| 2026-01-20 | Auto-Load Works, Auto-Save Was Missing |
| 2026-01-20 | Always Backup Before Modifying User Files |
| 2026-01-19 | CGO Must Be Disabled for ARM64 Linux |
<!-- INDEX:END -->

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

**Application**: Run periodic permission hygiene using hack/sanitize-permissions.md runbook. Use /ctx-drift to detect permission drift (missing skills, stale entries, consolidation opportunities).

---

## [2026-02-15-044450] Skill vs runbook for agent self-modification

**Context**: Considered building a skill to clean up settings.local.json permissions

**Lesson**: When a skill would edit files that control agent behavior (permissions, hooks, instructions), a runbook is safer. Auto-accept makes self-modifying skills an escalation vector.

**Application**: Use runbooks (human edits, agent advises) for operations on .claude/settings.local.json, CLAUDE.md, hooks, and CONSTITUTION.md. Reserve skills for operations where agent autonomy is safe.

---

## [2026-02-15-040313] Cross-repo links to published docs should use ctx.ist

**Context**: hack/persistent-irc.md linked to docs/ via relative paths, getting-started.md linked to MANIFESTO.md via GitHub — both bypass ctx.ist rendering (admonitions, nav, search)

**Lesson**: When content is published on ctx.ist, always link to the site URL, not the GitHub blob or a relative file path. GitHub won't render zensical admonitions and readers lose navigation context.

**Application**: When adding See Also or cross-references in runbooks or docs, use https://ctx.ist/... URLs for anything the site publishes. Reserve GitHub links for repo-only content (issues, releases, security tab, source files not on the site).

---

## [2026-02-15-034225] G304 gosec false positives in test files are safe to suppress

**Context**: gosec flags os.ReadFile with variable paths as G304 (potential file inclusion), even in test files where paths come from t.TempDir() and compile-time constants

**Lesson**: G304 requires user-controlled input to be exploitable. Test files using t.TempDir() and constants have no attack vector. Suppress with //nolint:gosec // test file path

**Application**: When gosec raises G304 in test files, verify paths aren't from external input, then suppress with nolint comment rather than restructuring the code

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

## [2026-02-14-163855] Skills can replace CLI commands for interactive workflows

**Context**: Evaluating whether /ctx-borrow needed a full CLI command or if the skill was sufficient

**Lesson**: A well-structured skill recipe is a guide, not a rigid script. The agent improvises beyond literal instructions and adapts to edge cases using its available tools.

**Application**: Prefer skills over CLI commands when the workflow requires judgment calls (conflict resolution, selective application, strategy selection). Reserve CLI commands for deterministic, non-interactive operations.

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

**Context**: During a stale-reference sweep, the .context/sessions/ gitignore rule was removed because sessions were consolidated into journals. But the gitignore rule exists to prevent sensitive data from being committed, not to document architecture. The directory may still exist locally.

**Lesson**: Gitignore entries for sensitive paths are security controls, not documentation. Never remove them during doc/reference cleanups even if the feature they relate to was removed.

**Application**: Before removing any gitignore entry, ask: does this entry exist for security/privacy or for architecture? Security entries stay permanently.

---

## [2026-02-11-124635] Chain-of-thought prompting improves agent reasoning accuracy

**Context**: Research shows accuracy on reasoning tasks jumps from 17.7% to 78.7% by adding think step-by-step to prompts. Applied this across agent guidelines.

**Lesson**: Explicit think step-by-step instructions in agent prompts dramatically improve reasoning accuracy at negligible token cost. This applies to skill files, playbooks, and autonomous loop prompts — anywhere the agent makes decisions before acting.

**Application**: Added Reason Before Acting section to AGENT_PLAYBOOK.md and reasoning nudges to 7 skills (ctx-implement, brainstorm, ctx-reflect, ctx-loop, qa, verify, consolidate). For autonomous loops, include reasoning instructions in PROMPT.md.

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

## [2026-01-28-194113] Claude Code Hooks Receive JSON via Stdin

**Context**: Debugging Claude Code PreToolUse hooks - they were not receiving
command data when using environment variables like CLAUDE_TOOL_INPUT

**Lesson**: Claude Code hooks receive input as JSON via stdin, not environment
variables. Use HOOK_INPUT=$(cat) then parse with
jq: COMMAND=$(echo "$HOOK_INPUT" | jq -r ".tool_input.command // empty")

**Application**: All hook scripts should read stdin for input. The JSON
structure includes .tool_input.command for Bash commands. Test hooks with
debug logging to /tmp/ to verify they receive expected data.

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

## [2026-01-28-040251] AI session JSONL formats are not standardized

**Context**: Building recall feature to parse session history from multiple
AI tools

**Lesson**: Claude Code, Cursor, Aider each have different JSONL formats
or may not export sessions at all.

**Application**: Use tool-agnostic Session type with tool-specific parsers.

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

## [2026-01-23-160000] Claude Code Skills Format

**Context**: Needed to understand how to create custom slash commands.

**Lesson**: Claude Code skills are markdown files in `.claude/commands/` with
YAML frontmatter (`description`, `argument-hint`, `allowed-tools`). Body is
the prompt. Use code blocks with `!` prefix for shell execution. `$ARGUMENTS`
passes command args.

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
