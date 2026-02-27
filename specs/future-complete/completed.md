### Phase 6: Session Ceremonies `#priority:medium`

**Context**: Sessions have two bookend rituals — init (`/ctx-remember`) and
wrap-up (`/ctx-wrap-up`). Unlike other ctx skills that encourage conversational
use, these should be explicitly invoked as slash commands for precision and
completeness. `/ctx-remember` already exists; `/ctx-wrap-up` is new.
Spec: `specs/session-wrap-up.md`

### Phase 7: Smart Retrieval `#priority:high`

**Context**: `ctx agent --budget` is cosmetic — the budget value is displayed
but never used for content selection. LEARNINGS.md is entirely excluded.
Decisions are title-only. No relevance filtering. This is the highest-impact
improvement for agent context quality. Spec: `specs/smart-retrieval.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 1)

### Phase 8: Drift Nudges `#priority:high`

**Context**: Context files grow without feedback. A project with 47 learnings
gets the same `ctx drift` output as one with 5. Entry count warnings nudge
users to consolidate or archive. Spec: `specs/drift-nudges.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 2)

