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

- [-] Promote 'block-dangerous-commands' to a real ctx system Go subcommand so OpenCode and other non-Claude editor integrations can ship the safety hook #priority:medium #added:2026-04-26-152911 #skipped:2026-04-26-231517 reason: decided not to do — OpenCode's exit-code semantics make a Cobra-based block-command shim too risky, and the safety-net omission in OpenCode is now treated as permanent (see decision 2026-04-26-231517)

- [ ] Task 1
- [ ] Task 2

### Phase 2: [Name] `#priority:medium`
- [ ] Task 1
- [ ] Task 2

### Phase CP: Ceremony Profiles `#priority:medium #added:2026-04-26`

Spec: `specs/ceremony-profiles.md`

- [ ] Add `Ceremony{Remember,WrapUp}` to `internal/rc/types.go`; apply defaults in `internal/rc/rc.go` from `internal/config/ceremony/ceremony.go` constants
- [ ] Thread resolved ceremony names into `ScanJournalsForCeremonies` and `Emit` in `internal/cli/system/core/ceremony/ceremony.go` (replace direct constant reads)
- [ ] Convert `internal/assets/hooks/messages/check-ceremony/{remember,wrapup,both}.txt` to `{REMEMBER}` / `{WRAPUP}` sentinels; audit `internal/config/embed/text` ceremony desc keys for the same
- [ ] Add a single sentinel-substitution helper (extend `internal/cli/system/core/message.Load` or sibling) so substitution happens in one place
- [ ] Show active ceremony profile (one line) in `ctx status` output
- [ ] Tests: default profile renders `/ctx-remember` `/ctx-wrap-up`; project with `ceremony.remember: dp-remember` renders `/dp-remember` and scanner only counts `dp-remember` as fulfilling the open-bookend
- [ ] Document in `docs/recipes/` with the editorial-project (DR knowledgebase) consumer as the worked example

## Blocked
