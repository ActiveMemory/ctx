 ---
title: Add workflow tips section to ctx init output
date: 2026-03-22
status: ready
---

# Workflow Tips After ctx init

## Problem

`ctx init` shows generic next steps (edit TASKS.md, run ctx status).
New users don't know about the ceremony loop, journal pipeline, or
key skills that make ctx effective.

## Solution

Add a "Workflow tips" block after the existing "Next steps" in the
init output. Each skill gets a one-line description derived from
reading the actual SKILL.md files.

### Layout

```
Workflow tips:

  Every session:
    /ctx-remember             Recall context and pick up where you left off
    /ctx-wrap-up              Capture learnings, decisions, and tasks before ending

  During work:
    /ctx-status               Check context health, token usage, and file summary
    /ctx-next                 Analyze tasks and suggest 1-3 concrete next actions
    /ctx-commit               Commit code, then prompt for decisions worth persisting
    /ctx-verify               Run verification before claiming a task is done

  Planning and design:
    /ctx-brainstorm           Structured design dialogue before implementation
    /ctx-spec                 Scaffold a feature spec from the project template
    /ctx-implement            Execute a plan step-by-step with checkpointed verification

  Periodic maintenance:
    /ctx-architecture         Build and refresh ARCHITECTURE.md and DETAILED_DESIGN.md
    /ctx-consolidate          Merge overlapping entries in DECISIONS.md and LEARNINGS.md
    /ctx-drift                Detect stale paths, broken references, and outdated context

  Journal pipeline (every few sessions):
    ctx recall export --all   Export session transcripts to .context/journal/
    /ctx-journal-enrich-all   Add frontmatter, tags, and summaries to exported entries

  Run 'ctx guide' for the full command reference.
```

### Implementation

1. Add YAML entry `write.init-workflow-tips` with the block above
2. Add Go constant `DescKeyWriteInitWorkflowTips`
3. Add `InfoWorkflowTips(cmd)` to `write/initialize/`
4. Call it from `initialize/cmd/root/run.go` after `InfoNextSteps`
