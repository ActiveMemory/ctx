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
