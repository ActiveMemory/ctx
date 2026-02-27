# Archived Learnings (consolidated 2026-02-27)

Originals replaced by consolidated entries in LEARNINGS.md.

## Group: Context injection and compliance strategy

## [2026-02-26-182747] Verbal summaries with linked diagram files are the effective pattern for lean context injection

**Context**: Auto-injected context via context-load-gate was ~20K tokens. ARCHITECTURE.md alone was 49K chars (~12K tokens) — 61% of the budget — mostly ASCII/Mermaid diagrams.

**Lesson**: Verbal summaries with linked diagram files (architecture-dia-*.md) cut ARCHITECTURE.md from ~12K to ~3.8K tokens without losing information. Diagrams live outside FileReadOrder so they're never auto-injected but remain available on demand. Same pattern applied to AGENT_PLAYBOOK.md by deduplicating sections already present in CLAUDE.md or CONVENTIONS.md.

**Application**: When context files grow heavy, check for diagrams and duplicated content first. Extract visual content to linked files, keep prose summaries inline. The 4-chars-per-token estimator is accurate — optimize the content, not the estimator.

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

## Group: Journal and recall parsing edge cases

## [2026-02-24-022214] /ctx-journal-normalize is dangerous at scale on non-ctx projects

**Context**: Discussed whether to keep normalize in the default journal pipeline

**Lesson**: On projects with large session JSONL files (millions of lines), the normalize skill blows up subagent context windows, consumes excessive tokens, and produces nondeterministic half-baked outputs

**Application**: Keep expensive AI skills out of batch pipelines; offer them as targeted per-file tools instead

---

## [2026-02-14-163549] normalizeCodeFences regex splits language specifiers

**Context**: Writing test for normalizeCodeFences, expected inline fence with lang tag to stay joined but the regex matched characters after backticks

**Lesson**: The inline fence regex treats any non-whitespace adjacent to triple-backtick fences as a split point, separating lang tags from the fence

**Application**: When testing normalizeCodeFences, use plain fences without language tags. See internal/cli/recall/fmt_test.go.

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

## Group: Zensical site builder quirks

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

## Group: Project identity and structure clarifications

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
