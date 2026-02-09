# Glossary

## Domain Terms

| Term                   | Definition                                                                                                                                                                                                  |
|------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Context                | The set of `.context/*.md` files that give AI agents persistent project knowledge across sessions. Not a generic word; when capitalized, refers specifically to the ctx system.                             |
| Context packet         | The token-budgeted markdown blob assembled by `ctx agent`. Contains prioritized excerpts from context files, sized to fit the agent's context window.                                                       |
| Context file           | Any `.md` file in `.context/` that ctx manages (CONSTITUTION, TASKS, DECISIONS, etc.). Each has a defined purpose and priority.                                                                             |
| Constitution           | The set of inviolable rules in `CONSTITUTION.md`. Distinct from conventions: constitution rules cannot be bent; violating one means the task is wrong.                                                      |
| Convention             | A project pattern or standard in `CONVENTIONS.md`. Conventions are strong recommendations that can be bent with justification; contrast with constitution rules.                                            |
| Drift                  | When context files diverge from the actual codebase state. Types: dead path references, stale task counts, missing required files, potential secrets. Detected by `ctx drift`.                              |
| Dead path              | A backtick-enclosed file path in ARCHITECTURE.md or CONVENTIONS.md that references a file no longer on disk. A drift warning type.                                                                          |
| Staleness              | When context files have not been updated to reflect recent code changes. Specific indicator: >10 completed tasks in TASKS.md signals the file needs compaction.                                             |
| Read order             | The priority sequence in which context files are loaded and presented to agents. Defined by `config.FileReadOrder`. Higher priority files are loaded first and survive token budget cuts.                   |
| Token budget           | Maximum estimated token count for assembled context. Default 8000. Configurable via `CTX_TOKEN_BUDGET`, `.contextrc`, or `--budget` flag. Uses 4-chars-per-token heuristic.                                 |
| Curated tier           | The `.context/*.md` files: manually maintained, token-budgeted, loaded by `ctx agent`. Contrast with full-dump tier.                                                                                        |
| Full-dump tier         | The `.context/sessions/` directory: raw session transcripts and snapshots. Not auto-loaded; used for archaeology when curated context is insufficient.                                                      |
| Session snapshot       | A point-in-time capture of session state, saved to `.context/sessions/`. Can be manual (`.md`) or auto-saved (`.jsonl`).                                                                                    |
| Compaction             | The process of archiving completed tasks and cleaning up context files. Run via `ctx compact`. Moves completed tasks to archive; preserves phase structure.                                                 |
| Entry header           | The timestamped heading format used in DECISIONS.md and LEARNINGS.md: `## [YYYY-MM-DD-HHMMSS] Title`. Parsed by `config.RegExEntryHeader`.                                                                  |
| Index table            | The auto-generated markdown table at the top of DECISIONS.md and LEARNINGS.md (between `<!-- INDEX:START -->` and `<!-- INDEX:END -->` markers). Updated by `ctx add` and `ctx decision/learnings reindex`. |
| Readback               | A structured summary where the agent plays back what it knows (last session, active tasks, recent decisions) so the user can confirm correct context was loaded. From aviation: pilots repeat ATC instructions back to confirm they heard correctly. In ctx, triggered by "do you remember?" or `/ctx-remember`. |
| Ralph Loop             | An iterative autonomous AI development workflow that uses PROMPT.md as a directive. Separate from ctx but complementary: Ralph drives the loop, ctx provides the memory.                                    |
| IMPLEMENTATION_PLAN.md | The orchestrator's directive file. Contains the meta-task ("check your tasks"), not the tasks themselves. Lives in project root, not `.context/`.                                                           |
| Skill                  | A Claude Code Agent Skill: a markdown file in `.claude/skills/` that teaches the agent a specialized workflow. Invoked via `/skill-name`.                                                                   |
| Live skill             | The project-local copy of a skill in `.claude/skills/`. Can be edited by the user or agent. Contrast with template skill.                                                                                   |
| Template skill         | The embedded copy of a skill in `internal/tpl/claude/skills/`. Deployed on `ctx init`. Source of truth for the default version.                                                                             |
| Hook                   | A Claude Code lifecycle script in `.claude/hooks/`. Fires on events: PreToolUse, UserPromptSubmit, SessionEnd. Generated by `ctx init`.                                                                     |
| Consolidation          | A code-quality sweep checking for convention drift: magic strings, predicate naming, file size, dead exports, etc. Run via `/consolidate` skill. Distinct from compaction (which is context-level).         |
| 3:1 ratio              | Heuristic for consolidation timing: consolidate after every 3 feature/bugfix sessions. Prevents convention drift from compounding.                                                                          |
| E/A/R classification   | Expert/Activation/Redundant taxonomy for evaluating skill quality. Good skill = >70% Expert knowledge, <10% Redundant with what the model already knows.                                                    |

## Abbreviations

| Abbreviation | Expansion                                                                                                   |
|--------------|-------------------------------------------------------------------------------------------------------------|
| ctx          | Context (the CLI tool and the system it manages)                                                            |
| rc           | Runtime configuration (from Unix `.xxxrc` convention); refers to `.contextrc` and the `internal/rc` package |
| tpl          | Template; the `internal/tpl` package containing embedded go:embed templates                                 |
| CWD          | Current working directory; used in session matching to correlate sessions with projects                     |
| JSONL        | JSON Lines; the format Claude Code uses for session transcripts (one JSON object per line)                  |
