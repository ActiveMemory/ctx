# META: Build the Planning Runner System (Claude Code)

You are building a **planning runner system** for Claude Code workflows.

## Goal
Create a reusable, deterministic “execution-first” planning loop that:
- asks the human for missing config details,
- reads a small set of context files,
- generates an execution-grade plan,
- maintains a compact canonical plan state to control token cost,
- produces timestamped run artifacts,
- converges quickly with explicit stop conditions.

You have autonomy to improve the design, templates, and file layout as you see fit.
Prefer out-of-the-box optimizations that reduce token cost and improve convergence.

## Constraints / Invariants
- The runner must be **generic** (usable across projects).
- It must be **execution-first by default** (deterministic plans; fast convergence).
- It must still support **strategy mode** as an optional flag (next-quarter planning).
- It must not require me to write long prose for every run.
- It must be safe against “planner drift” and “plan-of-the-day syndrome”:
    - enforce stop rules, iteration budget defaults, and change discipline.
- It must address the “git context” issue:
    - require repo baseline fields (branch + commit hash) and/or a sync protocol.
- It must keep token cost low:
    - canonical state file and bounded rereads.

## Deliverables (you will create files in this repo)
Create these paths and files (you may add more if useful):

1) `planning/README.md`
    - How to use the planning runner
    - Example usage snippets
    - Explanation of modes: execution vs strategy
    - Explanation of folder layout and artifacts

2) `planning/PLANNER_RUNNER.md`
    - The actual runner “command” prompt I paste into Claude Code for planning runs
    - Must include:
        - CONFIG block
        - INTERVIEW rules (ask only missing fields; <= 6 questions)
        - planning loop steps
        - plan template headings
        - budget + stop rules
    - You may refine/amend the template below

3) `planning/DEV_RUNNER.md` (optional but recommended)
    - A companion prompt that consumes the canonical plan and executes it stepwise
    - Must report back only “delta facts” useful for the planner:
        - interface/API changes
        - schema/config changes
        - failing tests + error summaries
        - commits produced (hashes) and branch
    - Must prefer small commits and validations per milestone

4) `planning/_state/CANONICAL_PLAN.md`
    - Initialize as an empty but structured template with brief guidance
    - Must have a hard size guideline (e.g., <= 300 lines)

5) Create `.gitignore` additions if needed for generated artifacts (only if appropriate)

## Given Runner Template (starting point; you can improve)
Execution mode headings must exist (you can add more):
- Plan ID / Timestamp
- Repo baseline
- Scope boundaries
- Goals / Non-goals
- Constraints / invariants
- Success criteria (binary/observable)
- Assumptions (minimize; list explicitly)
- Milestones (ordered; each includes)
    - Change summary
    - Files likely touched
    - Validation (exact command/test + expected outcome)
    - Rollback / escape hatch
- Risks + mitigations
- Open questions (each has a next action)
- PR/commit strategy
- Changelog since previous iteration

Default loop behaviors:
- Create a timestamped run file under `planning/_runs/`
- Maintain `planning/_state/CANONICAL_PLAN.md` as the concise “truth”
- On iterations after v1, reread only CANONICAL_PLAN + last run file
- Default budget: max_iterations=3
- Default stop rule:
  Stop if:
  (a) milestones all have validations,
  (b) risks/open questions have next actions,
  (c) last iteration made no structural change.

## Autonomy (important)
You may propose a better architecture for the runner system if you can justify it:
- e.g., more compact state representation,
- better convergence signals,
- better separation between planning and execution,
- smarter interview gating,
- improved artifact naming,
- improved strategies to avoid oscillation.

If you change the structure, update `planning/README.md` accordingly.

## Workflow
1) Inspect the repository briefly to understand context:
    - list directories
    - find existing docs or workflows that might integrate well
2) Create/update the deliverable files.
3) Ensure the README includes at least one full “copy/paste” example for:
    - execution run
    - strategy run
    - dev runner usage (if you create it)
4) Keep everything concise. This system should feel lightweight.

## Output requirement
After writing the files, print a short summary:
- what files were created/updated
- how to run the planning loop
- what defaults are used and how to override them
