# Tasks

## In Progress

## Next Up

### Enhance `amem init` to Create Claude Hooks `#priority:high` `#area:cli`
- [ ] Embed hook scripts in binary (like templates)
- [ ] Create `.claude/hooks/auto-save-session.sh` during init
- [ ] Create `.claude/settings.local.json` with PreToolUse and SessionEnd hooks
- [ ] Detect platform to set correct binary path in hooks
- [ ] Update `amem init` output to mention Claude Code integration

### Handle CLAUDE.md Creation/Merge `#priority:high` `#area:cli`
- [ ] Create CLAUDE.md if it doesn't exist
- [ ] If CLAUDE.md exists, backup to CLAUDE.md.<unix_timestamp>.bak before any modification
- [ ] Detect existing amem content via marker comment (`<!-- amem:context -->`)
- [ ] If no amem content, offer to merge (output snippet + prompt)
- [ ] Add `--merge` flag to auto-append without prompting
- [ ] Ensure idempotency — running init twice doesn't duplicate content

### Session Management Commands `#priority:high` `#area:cli`
- [ ] Implement `amem session save` — manually dump context to sessions/
- [ ] Implement `amem session list` — list saved sessions with summaries
- [ ] Implement `amem session load <file>` — load/summarize a previous session
- [ ] Implement `amem session parse` — convert .jsonl transcript to readable markdown

### Auto-Save Enhancements `#priority:medium` `#area:cli`
- [ ] Add PreCompact behavior — auto-save before `amem compact` runs
- [ ] Extract key decisions/learnings from transcript automatically
- [ ] Consider `amem watch --auto-save` mode

### Documentation `#priority:medium` `#area:docs`
- [ ] Document Claude Code integration in README
- [ ] Add "Dogfooding Guide" — how to use amem on amem itself
- [ ] Document session auto-save setup for new users

## Completed (Recent)

- [x] Set up PreToolUse hook for auto-load — 2025-01-20
- [x] Set up SessionEnd hook for auto-save — 2025-01-20
- [x] Create `.context/sessions/` directory structure — 2025-01-20
- [x] Create CLAUDE.md for native Claude Code bootstrapping — 2025-01-20
- [x] Document session persistence in AGENT_PLAYBOOK.md — 2025-01-20
- [x] Decide: always create .claude/ hooks (no --claude flag needed) — 2025-01-20

## Blocked
