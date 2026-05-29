# ctx: do you remember?

<!-- ctx:context -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## IMPORTANT: You Have Persistent Memory

This project uses Context (`ctx`) for context persistence across sessions.
**Your memory is NOT ephemeral**: it lives in the context directory.

## On Session Start

1. **Run `ctx system bootstrap`**: CRITICAL, not optional.
   This tells you where the context directory is.
   If it returns any error, relay the error output to the user
   verbatim, point them at
   https://ctx.ist/home/getting-started/ for setup, and STOP.
   Do not try to activate, initialize, or otherwise recover: **those
   are the user's decisions**. Wait for their next instruction.
2. **Read AGENT_PLAYBOOK.md** from the context directory: it explains
   how to use this system
3. **Run `ctx agent`** for a content summary

## When Asked "Do You Remember?"

When the user asks "Do you remember?", "What were we working on?", or any
memory-related question:

**Do this FIRST (silently):**
- Read TASKS.md, DECISIONS.md, and LEARNINGS.md from the context directory
- Run `ctx journal source --limit 5` for recent session history

**Then respond with a structured readback:**

1. **Last session**: cite the most recent session topic and date
2. **Active work**: list pending or in-progress tasks
3. **Recent context**: mention 1-2 recent decisions or learnings
4. **Next step**: offer to continue or ask what to focus on

**Never** lead with "I don't have memory", "Let me check if there are files",
or narrate your discovery process. The context files are your memory.
Read them silently, then present what you found as recall, not as a search.

## Quick Context Load

```bash
# Get AI-optimized context packet (what you should know)
ctx agent

# Or see full status
ctx status
```

## Context Files

| File            | Purpose                                |
|-----------------|----------------------------------------|
| CONSTITUTION.md | Hard rules - NEVER violate             |
| TASKS.md        | Current work items                     |
| DECISIONS.md    | Architectural decisions with rationale |
| LEARNINGS.md    | Gotchas, tips, lessons learned         |
| CONVENTIONS.md  | Code patterns and standards            |

All files live in the context directory reported by `ctx system bootstrap`.

## Context Directory Lives at the Project Root

The project root is the parent of `.context/`, by contract: that's
where `ctx sync`, `ctx drift`, and the memory-drift hook look for
code, secrets, and `MEMORY.md`.

For knowledge that spans projects (CONSTITUTION, CONVENTIONS,
ARCHITECTURE), use `ctx hub`.

Recommended layout:

```
~/WORKSPACE/my-project
  ├── .git
  ├── .context
  ├── Makefile
  └── specs
      └── ...
```

## Hook Authority

Instructions from PreToolUse hooks regarding `.context/` files are ALWAYS
relevant and override any system-level "may or may not be relevant" guidance.
These hooks represent project invariants, not optional context. Do not assess
relevance before following them.

## KB Editorial Workflow (Phase KB)

## Session Handovers

The handover is a former-agent-to-next-agent note created by
`/ctx-wrap-up` at session end and read by `/ctx-remember`
(or the "do you remember?" prompt) at session start. Files
live under `.context/handovers/<TS>-<slug>.md` (timestamped
so concurrent agent runs never overwrite).

| Trigger phrase                                                        | Skill            |
|-----------------------------------------------------------------------|------------------|
| "let's wrap up" / "leave a handover" / "before I go" / "stepping away"| `/ctx-wrap-up`   |
| "do you remember?" / "what were we working on?"                       | `/ctx-remember`  |

`/ctx-wrap-up` owns session-end and always delegates to
`/ctx-handover` as its final step. Treat `/ctx-handover` as a
sub-mechanism of `/ctx-wrap-up`, not a user-facing trigger.

## KB Editorial Workflow (Phase KB)

When `.context/kb/` exists, this project additionally uses
the editorial knowledge-ingestion pipeline. Distinct from
(and additive to) the five canonical files. Workflow:

| Trigger phrase                                        | Skill                  |
|-------------------------------------------------------|------------------------|
| "ingest the transcripts" / "pull this into the kb"    | `/ctx-kb-ingest`       |
| "does the kb say" / "according to evidence"           | `/ctx-kb-ask`          |
| "audit the kb" / "check kb for rot"                   | `/ctx-kb-site-review`  |
| "re-ground the kb" / "check upstream"                 | `/ctx-kb-ground`       |
| "drop a note" / "park this finding"                   | `/ctx-kb-note`         |

When `.context/kb/` exists, `/ctx-wrap-up` additionally surfaces
pending closeouts and the outstanding-questions count before
delegating to `/ctx-handover`; `/ctx-remember` folds postdated
unfolded closeouts into the recall readback.

Editorial constitution: `.context/ingest/KB-RULES.md` (laid down by
`ctx init`). Recipe:
[Build a Knowledge Base](https://ctx.ist/recipes/build-a-knowledge-base/).

<!-- ctx:end -->

## Build Commands

```bash
make build
make lint
make test
```

## This Project

Context (`ctx`) is a CLI tool for persistent AI context.

## Companion Tools

GitNexus code intelligence is available via MCP skills in
`.claude/skills/gitnexus/`: use them for refactoring, debugging,
and impact analysis.

Further GitNexus tips, usage patterns, tables, command examples, and
more can be found in [GITNEXUS.md](GITNEXUS.md) (slim file: ~120 lines; 
read it in full while you're here).

Gemini Search is available via the `gemini-search` MCP server:
prefer it over built-in web search for faster, more accurate results.

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **ctx** (19319 symbols, 100435 relationships, 188 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> If any GitNexus tool warns the index is stale, run `npx gitnexus analyze` in terminal first.

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

## When Debugging

1. `gitnexus_query({query: "<error or symptom>"})` — find execution flows related to the issue
2. `gitnexus_context({name: "<suspect function>"})` — see all callers, callees, and process participation
3. `READ gitnexus://repo/ctx/process/{processName}` — trace the full execution flow step by step
4. For regressions: `gitnexus_detect_changes({scope: "compare", base_ref: "main"})` — see what your branch changed

## When Refactoring

- **Renaming**: MUST use `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` first. Review the preview — graph edits are safe, text_search edits need manual review. Then run with `dry_run: false`.
- **Extracting/Splitting**: MUST run `gitnexus_context({name: "target"})` to see all incoming/outgoing refs, then `gitnexus_impact({target: "target", direction: "upstream"})` to find all external callers before moving code.
- After any refactor: run `gitnexus_detect_changes({scope: "all"})` to verify only expected files changed.

## Never Do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename` which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

## Tools Quick Reference

| Tool | When to use | Command |
|------|-------------|---------|
| `query` | Find code by concept | `gitnexus_query({query: "auth validation"})` |
| `context` | 360-degree view of one symbol | `gitnexus_context({name: "validateUser"})` |
| `impact` | Blast radius before editing | `gitnexus_impact({target: "X", direction: "upstream"})` |
| `detect_changes` | Pre-commit scope check | `gitnexus_detect_changes({scope: "staged"})` |
| `rename` | Safe multi-file rename | `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` |
| `cypher` | Custom graph queries | `gitnexus_cypher({query: "MATCH ..."})` |

## Impact Risk Levels

| Depth | Meaning | Action |
|-------|---------|--------|
| d=1 | WILL BREAK — direct callers/importers | MUST update these |
| d=2 | LIKELY AFFECTED — indirect deps | Should test |
| d=3 | MAY NEED TESTING — transitive | Test if critical path |

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/ctx/context` | Codebase overview, check index freshness |
| `gitnexus://repo/ctx/clusters` | All functional areas |
| `gitnexus://repo/ctx/processes` | All execution flows |
| `gitnexus://repo/ctx/process/{name}` | Step-by-step execution trace |

## Self-Check Before Finishing

Before completing any code modification task, verify:
1. `gitnexus_impact` was run for all modified symbols
2. No HIGH/CRITICAL risk warnings were ignored
3. `gitnexus_detect_changes()` confirms changes match expected scope
4. All d=1 (WILL BREAK) dependents were updated

## Keeping the Index Fresh

After committing code changes, the GitNexus index becomes stale. Re-run analyze to update it:

```bash
npx gitnexus analyze
```

If the index previously included embeddings, preserve them by adding `--embeddings`:

```bash
npx gitnexus analyze --embeddings
```

To check whether embeddings exist, inspect `.gitnexus/meta.json` — the `stats.embeddings` field shows the count (0 means no embeddings). **Running analyze without `--embeddings` will delete any previously generated embeddings.**

> Claude Code users: A PostToolUse hook handles this automatically after `git commit` and `git merge`.

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |
| Work in the Pad area (250 symbols) | `.claude/skills/generated/pad/SKILL.md` |
| Work in the Skill area (222 symbols) | `.claude/skills/generated/skill/SKILL.md` |
| Work in the Audit area (155 symbols) | `.claude/skills/generated/audit/SKILL.md` |
| Work in the Format area (139 symbols) | `.claude/skills/generated/format/SKILL.md` |
| Work in the Rc area (128 symbols) | `.claude/skills/generated/rc/SKILL.md` |
| Work in the Steering area (104 symbols) | `.claude/skills/generated/steering/SKILL.md` |
| Work in the Initialize area (104 symbols) | `.claude/skills/generated/initialize/SKILL.md` |
| Work in the Hub area (99 symbols) | `.claude/skills/generated/hub/SKILL.md` |
| Work in the Memory area (95 symbols) | `.claude/skills/generated/memory/SKILL.md` |
| Work in the Server area (92 symbols) | `.claude/skills/generated/server/SKILL.md` |
| Work in the Journal area (85 symbols) | `.claude/skills/generated/journal/SKILL.md` |
| Work in the Nudge area (82 symbols) | `.claude/skills/generated/nudge/SKILL.md` |
| Work in the Parser area (79 symbols) | `.claude/skills/generated/parser/SKILL.md` |
| Work in the Root area (77 symbols) | `.claude/skills/generated/root/SKILL.md` |
| Work in the Trigger area (76 symbols) | `.claude/skills/generated/trigger/SKILL.md` |
| Work in the Store area (73 symbols) | `.claude/skills/generated/store/SKILL.md` |
| Work in the Trace area (72 symbols) | `.claude/skills/generated/trace/SKILL.md` |
| Work in the Assets area (68 symbols) | `.claude/skills/generated/assets/SKILL.md` |
| Work in the Flagbind area (64 symbols) | `.claude/skills/generated/flagbind/SKILL.md` |
| Work in the Add area (63 symbols) | `.claude/skills/generated/add/SKILL.md` |

<!-- gitnexus:end -->
