# Tasks

<!--
UPDATE WHEN:
- New work is identified → add task with #added timestamp
- Starting work → add #in-progress or #started timestamp
- Work completes → mark [x]
- Work is blocked → add to Blocked section with reason
- Scope changes → update task description inline

DO NOT UPDATE FOR:
- Reorganizing or moving tasks (violates CONSTITUTION)
- Removing completed tasks (use ctx task archive instead)

STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently: never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers

TASK STATUS LABELS:
  `[ ]`: pending
  `[x]`: completed
  `[-]`: skipped (with reason)
  `#in-progress`: currently being worked on (add inline, don't move task)
-->

### Phase 1: [Name] `#priority:high`
- [ ] Add TypeScript type-check step (bunx tsc --noEmit) for embedded editor-plugin assets to CI; nothing currently checks .opencode/plugins/ctx/index.ts before embedding #priority:low #added:2026-04-26-152912

- [ ] Promote 'block-dangerous-commands' to a real ctx system Go subcommand so OpenCode and other non-Claude editor integrations can ship the safety hook #priority:medium #added:2026-04-26-152911

- [ ] Task 1
- [ ] Task 2

### Phase 2: [Name] `#priority:medium`
- [ ] Task 1
- [ ] Task 2

## Blocked
