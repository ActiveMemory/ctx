
# Tasks

<!--
STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently — never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers
-->

### Phase 1: Journal Site Improvements `#priority:high`

**Context**: Enriched journal files now have YAML frontmatter (topics, type, outcome,
technologies, key_files). The site generator should leverage this metadata for
better discovery and navigation.

- [ ] Check for the latest backup of .context (there is a script in ~/WORKSPACE)
  and if it is >2days old, remind the user to run the backup script.
  we can use a hook for this. Alternative could be detecting this,
  exiting early, and using a cron job; but the disadvantage is, if nothing has
  changed, the cron job will still make a redundant backup: If nobody is using
  the project, then there is no need to run the backup anyway.

- [ ] Build minimal ctx monitor dashboard (see specs/monitor-architecture.md). 
  - [ ] MVP: a terminal script that runs in a separate window, finds the active 
    session JSONL, estimates token usage, shows context health (green/yellow/red), 
    and refreshes periodically. Later: pluggable auditors 
    (drift, repetition detection), signal injection via hooks. #priority:medium #added:2026-02-04-223356

- [x] Search for a consolidation/code-hygiene skill that can enforce
  project-specific conventions. Today's codebase scan found concrete
  drift: 5 Is* predicate violations, magic strings in 7+ files, 80+ hardcoded
  file permissions. A skill that periodically checks for ctx's known drift
  patterns (godoc style, predicate naming, file organization, magic strings
  in config) would be the enforcement arm of CONVENTIONS.md. Look for skills
  focused on codebase consistency, convention enforcement, or technical debt
  scanning — not generic code review checklists. #priority:medium #added:2026-02-04-222224
  **Done**: Found generic skill (70% redundant per E/A/R), adapted to
  `/consolidate` with 9 project-specific checks. Added line width +
  duplication conventions. Blog post drafted. #completed:2026-02-05

**Features (priority order):**

- [ ] T1.1: Topics system
      - Single `/topics/index.md` page
      - Popular topics (2+ sessions) get dedicated pages (`/topics/{topic}.md`)
      - Long-tail topics (1 session) listed inline with direct session links
      - All on one page for Ctrl+F discoverability
      #added:2026-02-03

- [ ] T1.2: Key files index
      - Reverse lookup: file → sessions that touched it
      - Uses `key_files` from frontmatter (not parsed from conversation)
      - `/files/index.md` or similar
      #added:2026-02-03

- [ ] T1.3: Session type pages
      - Dedicated page per type: `/types/debugging.md`, `/types/feature.md`, etc.
      - Groups sessions by type (feature, bugfix, refactor, exploration, debugging, documentation)
      #added:2026-02-03

**Deferred:**
- Timeline narrative (nice-to-have, low priority)
- Outcome filtering (uncertain value, revisit after seeing data)
- Stats dashboard (skipped - gamification, low ROI)
- Technology index (skipped - not useful for this project)

**Design note**: Topics and key-files both need content search across journal
entries (not just filename matching). Consider FTS5 or similar full-text
indexing — grep across 100+ journal files won't scale.

**Design status**: Understanding confirmed, ready for design approaches.

### Phase 2: Export Preservation `#priority:medium`

- [ ] T2.1: `ctx recall export --update` mode
      - Preserve YAML frontmatter and summary when re-exporting
      - Update only the raw conversation content
      - Solves: `--force` loses enrichments, no-force can't update
      #added:2026-02-03

## Blocked

## Reference

**Task Status Labels**:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)
