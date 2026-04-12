# Copilot Feature Parity Kit

## Problem

Claude Code integration has 41 skills, 20+ hooks, a governance/ceremony
system, and deep code intelligence (GitNexus). Copilot CLI has 5 skills
and 4 lifecycle hooks. Copilot VS Code has 45+ slash commands but most
delegate to CLI without agent-level intelligence (workflow orchestration,
proactive nudges, code-aware operations).

This creates a two-tier experience: Claude users get a full context-aware
AI partner; Copilot users get a context reader. The gap is not in data
access (MCP tools are shared) but in **agent behavior** — the skills,
hooks, and governance that turn raw context into workflow.

## Approach

A phased spec kit that brings Copilot CLI and VS Code to parity with
Claude Code. Each phase is independently shippable. Work is organized
into three layers:

1. **Skills** — Copilot CLI `.github/skills/` and VS Code slash commands
2. **Hooks** — Copilot CLI hook scripts and VS Code event handlers
3. **Governance** — Proactive nudges, ceremony checks, session health

Architecture principle: **single source of truth**. Skills are authored
as markdown SKILL.md files in `internal/assets/integrations/copilot-cli/skills/`
and deployed by `ctx setup copilot-cli --write`. VS Code commands call
`ctx` CLI or MCP tools — the extension does not duplicate skill logic.

### Cross-references

- `specs/copilot-cli-integration.md` — existing feature matrix (context injection, hooks, MCP, recall)
- `specs/vscode-feature-parity.md` — existing layer-by-layer VS Code mapping
- This spec supersedes neither; it fills the **skill + governance gap** they identify.

---

## Phase 1 — Core Workflow Skills (Copilot CLI + VS Code)

Port the skills that encode the work cycle: pick → implement → commit → reflect.

### 1.1 Skills to Port

| # | Claude Skill | Copilot CLI Skill | VS Code Command | Priority |
|---|-------------|-------------------|-----------------|----------|
| 1 | `ctx-next` | `ctx-next/SKILL.md` ✅ exists | `/next` | P0 — already done |
| 2 | `ctx-implement` | `ctx-implement/SKILL.md` | `/implement` | P0 |
| 3 | `ctx-commit` | `ctx-commit/SKILL.md` | `/commit` | P0 |
| 4 | `ctx-reflect` | `ctx-reflect/SKILL.md` | `/reflect` | P0 |
| 5 | `ctx-remember` | `ctx-remember/SKILL.md` | `/remember` | P0 |
| 6 | `ctx-wrap-up` | `ctx-wrap-up/SKILL.md` | `/wrapup` | P0 |
| 7 | `ctx-code-review` | `ctx-code-review/SKILL.md` | `/review` | P1 |
| 8 | `ctx-refactor` | `ctx-refactor/SKILL.md` | `/refactor` | P1 |
| 9 | `ctx-explain` | `ctx-explain/SKILL.md` | `/explain` | P1 |
| 10 | `ctx-brainstorm` | `ctx-brainstorm/SKILL.md` | `/brainstorm` | P1 |
| 11 | `ctx-spec` | `ctx-spec/SKILL.md` | `/spec` | P1 |

### 1.2 Skill File Format (Copilot CLI)

Copilot CLI skills live in `.github/skills/<name>/SKILL.md`. Format:

```markdown
---
name: ctx-implement
description: Execute implementation plan step-by-step with verification
tools: [bash, read, write, edit, glob, grep]
---

# ctx-implement

## When to Use
- User says "implement this", "build it", "start coding"
- A task from TASKS.md is selected for implementation

## When NOT to Use
- No spec or plan exists (use ctx-spec first)
- Task is ambiguous (use ctx-brainstorm first)

## Workflow
1. Read the referenced spec from `specs/`
2. Read CONVENTIONS.md for code patterns
3. Break work into chunks, commit after each
4. Run `make lint && make test` after each chunk
5. Mark task done in TASKS.md when complete
6. Offer to record learnings/decisions discovered

## Quality Gates
- [ ] Spec exists and was read
- [ ] Tests pass after each chunk
- [ ] Lint passes
- [ ] TASKS.md updated
```

### 1.3 VS Code Slash Command Wiring

For each new skill, add to `editors/vscode/package.json` contributes
and handle in `extension.ts`:

```typescript
case '/implement':
  return runCtxCommand(stream, 'implement', request.prompt);
```

VS Code commands delegate to `ctx` CLI. The extension provides UI
(progress, follow-ups, markdown rendering) but not logic.

### 1.4 Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/integrations/copilot-cli/skills/ctx-implement/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-commit/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-reflect/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-remember/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-wrap-up/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-code-review/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-refactor/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-explain/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-brainstorm/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-spec/SKILL.md` | New skill |
| `internal/cli/setup/core/copilot_cli/copilot_cli.go` | Deploy new skills |
| `editors/vscode/package.json` | Add slash commands |
| `editors/vscode/src/extension.ts` | Handle new commands |

---

## Phase 2 — Architecture & Design Skills

### 2.1 Skills to Port

| # | Claude Skill | Copilot CLI Skill | VS Code Command | Priority |
|---|-------------|-------------------|-----------------|----------|
| 1 | `ctx-architecture` | `ctx-architecture/SKILL.md` | `/architecture` | P1 |
| 2 | `ctx-architecture-enrich` | `ctx-architecture-enrich/SKILL.md` | — (CLI only) | P2 |
| 3 | `ctx-architecture-failure-analysis` | `ctx-architecture-failure-analysis/SKILL.md` | — (CLI only) | P2 |
| 4 | `ctx-doctor` | `ctx-doctor/SKILL.md` | `/system doctor` ✅ | P1 |

### 2.2 GitNexus Dependency

`ctx-architecture-enrich` and `ctx-architecture-failure-analysis` use
GitNexus MCP tools (`mcp__gitnexus__*`). For Copilot CLI:

- GitNexus MCP server must be registered in `~/.copilot/mcp-config.json`
- Skills should gracefully degrade if GitNexus is unavailable
- Fallback: use `grep`/`go doc` for basic code intelligence

### 2.3 Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/integrations/copilot-cli/skills/ctx-architecture/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-architecture-enrich/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-architecture-failure-analysis/SKILL.md` | New skill |

---

## Phase 3 — Governance & Proactive Hooks

The biggest behavioral gap. Claude's hook system fires on every tool
use and user prompt, surfacing nudges for persistence, ceremonies,
and health. Copilot has no equivalent.

### 3.1 Copilot CLI Hooks to Add

Extend `.github/hooks/ctx-hooks.json` with richer behavior in existing
hook scripts:

| Hook | Trigger | Behavior | Claude Equivalent |
|------|---------|----------|-------------------|
| `preToolUse` (enhanced) | Every tool call | Dangerous cmd block + context load gate | `block-non-path-ctx` + `context-load-gate` |
| `postToolUse` (enhanced) | After edit/write | Task completion check + learning nudge | `check-task-completion` + `post-commit` |
| `sessionStart` (enhanced) | Session begin | Load context + version check + reminder relay | `budget-agent` + `check-version` + `check-reminders` |
| `sessionEnd` (enhanced) | Session close | Persistence ceremony + journal capture | `check-ceremonies` + `check-persistence` |

### 3.2 New Hook Scripts

```
.github/hooks/scripts/
├── ctx-preToolUse.sh        # Enhanced: dangerous cmd + context gate
├── ctx-preToolUse.ps1       # PowerShell mirror
├── ctx-postToolUse.sh       # Enhanced: task check + learning nudge
├── ctx-postToolUse.ps1      # PowerShell mirror
├── ctx-sessionStart.sh      # Enhanced: bootstrap + reminders
├── ctx-sessionStart.ps1     # PowerShell mirror
├── ctx-sessionEnd.sh        # Enhanced: ceremony + journal
└── ctx-sessionEnd.ps1       # PowerShell mirror
```

### 3.3 Governance Messages

Port the hook message registry. Each message has:
- **Condition**: when to fire (e.g., "uncompleted tasks > 5")
- **Message**: what to show (e.g., "⚠ 5 tasks pending — run `/next`")
- **Cooldown**: don't repeat within N minutes

| Message ID | Condition | Copilot Surface |
|------------|-----------|-----------------|
| `ceremony-remember` | Session > 30 min, no recall done | `sessionStart` script output |
| `ceremony-wrapup` | Session > 2 hours, no persist | `sessionEnd` script output |
| `persistence-nudge` | Decision made but not recorded | `postToolUse` script output |
| `task-completion` | File edited matching task description | `postToolUse` script output |
| `version-drift` | `ctx --version` != expected | `sessionStart` script output |
| `reminder-relay` | Pending reminders exist | `sessionStart` script output |

### 3.4 VS Code Governance

Map governance to VS Code extension events:

| Governance | VS Code Mechanism |
|------------|-------------------|
| Ceremony check | `vscode.window.onDidChangeWindowState` (focus loss → prompt) |
| Persistence nudge | `vscode.workspace.onDidSaveTextDocument` (`.context/` watch) |
| Task completion | `onDidSaveTextDocument` → `ctx system check-task-completion` |
| Reminder relay | Status bar item + 5-min timer (already partial) |
| Version check | `ensureCtxAvailable()` version comparison |
| Session ceremony | Extension `deactivate()` → wrap-up prompt |

### 3.5 Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/integrations/copilot-cli/scripts/ctx-preToolUse.sh` | Enhance with context gate |
| `internal/assets/integrations/copilot-cli/scripts/ctx-preToolUse.ps1` | PowerShell mirror |
| `internal/assets/integrations/copilot-cli/scripts/ctx-postToolUse.sh` | Enhance with task check + nudge |
| `internal/assets/integrations/copilot-cli/scripts/ctx-postToolUse.ps1` | PowerShell mirror |
| `internal/assets/integrations/copilot-cli/scripts/ctx-sessionStart.sh` | Enhance with bootstrap |
| `internal/assets/integrations/copilot-cli/scripts/ctx-sessionStart.ps1` | PowerShell mirror |
| `internal/assets/integrations/copilot-cli/scripts/ctx-sessionEnd.sh` | Enhance with ceremony |
| `internal/assets/integrations/copilot-cli/scripts/ctx-sessionEnd.ps1` | PowerShell mirror |
| `internal/assets/integrations/copilot-cli/ctx-hooks.json` | Updated hook config |
| `editors/vscode/src/extension.ts` | Add governance event handlers |

---

## Phase 4 — Context Health & Maintenance Skills

### 4.1 Skills to Port

| # | Claude Skill | Copilot CLI Skill | VS Code Command | Priority |
|---|-------------|-------------------|-----------------|----------|
| 1 | `ctx-consolidate` | `ctx-consolidate/SKILL.md` | `/consolidate` | P2 |
| 2 | `ctx-permission-sanitize` | `ctx-permission-sanitize/SKILL.md` | — | P2 |
| 3 | `ctx-prompt-audit` | `ctx-prompt-audit/SKILL.md` | — | P3 |
| 4 | `ctx-skill-audit` | `ctx-skill-audit/SKILL.md` | — | P3 |
| 5 | `ctx-skill-create` | `ctx-skill-create/SKILL.md` | — | P3 |
| 6 | `ctx-link-check` | `ctx-link-check/SKILL.md` | `/check-links` ✅ | Done |
| 7 | `ctx-pad` | `ctx-pad/SKILL.md` | `/pad` ✅ | Done |

### 4.2 Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/integrations/copilot-cli/skills/ctx-consolidate/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-permission-sanitize/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-prompt-audit/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-skill-audit/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-skill-create/SKILL.md` | New skill |

---

## Phase 5 — Journal & Documentation Skills

### 5.1 Skills to Port

| # | Claude Skill | Copilot CLI Skill | VS Code Command | Priority |
|---|-------------|-------------------|-----------------|----------|
| 1 | `ctx-journal-enrich` | `ctx-journal-enrich/SKILL.md` | `/journal enrich` | P2 |
| 2 | `ctx-journal-enrich-all` | `ctx-journal-enrich-all/SKILL.md` | — | P3 |
| 3 | `ctx-blog` | `ctx-blog/SKILL.md` | `/blog` | P3 |
| 4 | `ctx-blog-changelog` | `ctx-blog-changelog/SKILL.md` | `/changelog` ✅ | Done |
| 5 | `ctx-plan-import` | `ctx-plan-import/SKILL.md` | — | P3 |

### 5.2 Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/integrations/copilot-cli/skills/ctx-journal-enrich/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-journal-enrich-all/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-blog/SKILL.md` | New skill |
| `internal/assets/integrations/copilot-cli/skills/ctx-plan-import/SKILL.md` | New skill |

---

## Phase 6 — Advanced / Infrastructure Skills

### 6.1 Skills to Port

| # | Claude Skill | Copilot CLI Skill | VS Code Command | Priority |
|---|-------------|-------------------|-----------------|----------|
| 1 | `ctx-loop` | `ctx-loop/SKILL.md` | — (N/A for chat UI) | P3 |
| 2 | `ctx-worktree` | `ctx-worktree/SKILL.md` | `/worktree` ✅ | Done |
| 3 | `ctx-pause` / `ctx-resume` | `ctx-pause/SKILL.md` | `/pause` | P3 |

---

## Summary: Parity Scorecard

### Current State

| Surface | Skills | Hooks | Governance | Total Score |
|---------|--------|-------|------------|-------------|
| Claude Code | 41 | 20+ | Full | 100% |
| Copilot CLI | 5 | 4 | None | ~15% |
| Copilot VS Code | 45 cmds (CLI delegate) | 6 watchers | Partial | ~40% |

### After This Spec Kit

| Surface | Skills | Hooks | Governance | Total Score |
|---------|--------|-------|------------|-------------|
| Claude Code | 41 | 20+ | Full | 100% |
| Copilot CLI | 36 | 4 (enhanced) | Hook-based | ~85% |
| Copilot VS Code | 50+ cmds | 10+ handlers | Event-based | ~90% |

### Remaining Gap (Intentional)

These features are Claude Code specific and **not ported**:

| Feature | Reason |
|---------|--------|
| `UserPromptSubmit` hooks (12 types) | Copilot CLI has no equivalent trigger point |
| `check-context-size` (token budget) | Copilot does not expose token counts |
| `check-knowledge` | Claude Code knowledge graph specific |
| `heartbeat` telemetry | Different telemetry model |
| Plugin system (`.claude-plugin/`) | Claude Code specific packaging |

---

## Edge Cases

| Case | Expected Behavior |
|------|-------------------|
| Skill references MCP tool not available | Graceful degradation: use CLI fallback |
| GitNexus not registered | Architecture-enrich falls back to grep/ast |
| Hook script timeout (>5s) | Script returns empty, no block |
| Concurrent skill invocation | Each invocation is independent |
| `ctx` binary not on PATH | `sessionStart` hook warns; VS Code auto-downloads |
| `.context/` doesn't exist | Skills prompt `ctx init` |
| Windows vs Unix line endings | Scripts use native endings per platform |

## Validation Rules

- Every skill SKILL.md must have: name, description, tools, workflow, quality gates
- Every hook script must have: bash + PowerShell variants
- Every VS Code command must: delegate to CLI, show progress, offer follow-ups
- Skill names must match between Claude and Copilot (e.g., `ctx-implement` in both)

## Testing

- **Unit**: Each skill SKILL.md passes `ctx skill audit` (format, completeness)
- **Integration**: `ctx setup copilot-cli --write` deploys all skills + hooks
- **E2E**: Run Copilot CLI session with skills, verify workflow cycle
- **VS Code**: Extension test suite covers new slash commands
- **Cross-platform**: Hook scripts tested on bash (Linux/macOS) and PowerShell (Windows)

## Non-Goals

- Replacing Copilot's built-in features (code completion, inline suggestions)
- Porting Claude Code's plugin packaging system
- Real-time token budget monitoring (Copilot doesn't expose this)
- Bidirectional memory sync (separate spec: Copilot memory bridge)
- ACP server mode (separate spec: `specs/copilot-cli-integration.md` Phase 4)

## Open Questions

1. **Copilot CLI skill discovery**: Does Copilot CLI auto-discover `.github/skills/`
   or do skills need explicit registration? Need to verify with latest CLI docs.
2. **Hook script output rendering**: How does Copilot CLI render hook script
   stdout? Markdown? Plain text? This affects governance message formatting.
3. **VS Code command registration limit**: Is there a practical limit on slash
   commands in a chat participant? Current 45 → 55+ after this spec.
4. **Skill frontmatter schema**: Does Copilot CLI enforce a specific YAML
   frontmatter schema for SKILL.md, or is it freeform markdown?
