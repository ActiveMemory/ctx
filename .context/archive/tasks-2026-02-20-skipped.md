# Skipped Tasks — Archived 2026-02-20

Tasks deliberately skipped with documented reasons. Preserved for traceability.

## Phase 0: Ideas

### Improvements (REPORT-6)

- [-] CI coverage enforcement: extend `test-coverage` Makefile target
      beyond just `internal/context` to enforce project-wide and
      per-package minimums. #priority:medium #source:report-6
      Skipped: Replaced with HTML coverage report task (make test-cover → dist/coverage.html).

### User-Facing Documentation (REPORT-7)

- [-] Create migration/adoption guide for existing projects: how ctx
      interacts with existing CLAUDE.md, .cursorrules, the --merge
      flag in workflow context. #priority:medium #source:report-7
      Skipped: No migration needed — ctx is additive markdown files.
      Coexistence with CLAUDE.md/.cursorrules already covered in
      Getting Started and Integrations.

- [-] Create troubleshooting page: consolidate scattered troubleshooting
      into one page (ctx not in PATH, hooks not firing, context not
      loading). #priority:low #source:report-7
      Skipped: Not enough content to warrant a page. The few troubleshooting
      items are already inline where they belong (Integrations, Autonomous Loop).

## Phase 1: Journal Site Improvements

- [-] T1.4: Investigate inconsistent tool output collapsing — tool output
      collapsing removed, no longer part of the pipeline.
- [-] Timeline narrative — dropped, duplicates `/ctx-blog` skill (see DECISIONS.md)
- [-] Outcome filtering — uncertain value, revisit after seeing data
- [-] Stats dashboard — gamification, low ROI
- [-] Technology index — not useful for this project
