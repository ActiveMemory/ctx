# ctx v0.6.0: Claude Code Plugin Conversion

## Context

ctx (v0.4.0) distributes Claude Code integration via `ctx init`, which scaffolds
9 shell hook scripts into `.claude/hooks/`, writes hook config + permissions into
`settings.local.json`, and deploys 25+ skills into `.claude/skills/`. This approach
requires `jq` as a runtime dependency, couples ctx to per-project scaffolding, and
does not leverage Claude Code's plugin system (available since v1.0.33).

v0.6.0 converts ctx into a first-class Claude Code plugin. Shell hook scripts become
Go subcommands (`ctx system *`). Skills and hooks ship as a plugin directory. `ctx init`
becomes tool-agnostic. `ctx hook claude-code` is removed (replaced by the plugin).
Non-Claude tools (Cursor, Aider, Copilot, Windsurf) keep their `ctx hook <tool>` doc
generators. Version jumps from 0.4.0 to 0.6.0 to signal the magnitude of the change.

---

## 1. New `ctx system` Subcommands

**Package**: `internal/cli/system/`

Six subcommands, all hidden (`cmd.Hidden = true`), documented in user-facing docs
under "Advanced: Hook Internals".

| Subcommand | Replaces | Hook Event | Reads stdin |
|---|---|---|---|
| `check-context-size` | `check-context-size.sh` | UserPromptSubmit | `session_id` |
| `check-persistence` | `check-persistence.sh` | UserPromptSubmit | `session_id` |
| `check-journal` | `check-journal.sh` | UserPromptSubmit | (minimal) |
| `block-non-path-ctx` | `block-non-path-ctx.sh` | PreToolUse:Bash | `tool_input.command` |
| `post-commit` | `post-commit.sh` | PostToolUse:Bash | `tool_input.command` |
| `cleanup-tmp` | `cleanup-tmp.sh` | SessionEnd | (none) |

**NOT converted** (project-specific, not in embedded templates):
- `block-git-push.sh` — stays in ctx project's `.claude/hooks/`
- `block-dangerous-commands.sh` — stays in ctx project's `.claude/hooks/`
- `check-backup-age.sh` — depends on SMB config, stays project-local

### Shared infrastructure (`internal/cli/system/`)

- `input.go` — `HookInput` / `ToolInput` structs, `readInput()` helper
- `state.go` — `secureTempDir()` (duplicate of `internal/cli/agent/cooldown.go:26`
  pattern), counter read/write, log helper, daily throttle check
- `system.go` — parent Cobra command registration

### stdin JSON contract

```go
type HookInput struct {
    SessionID string    `json:"session_id"`
    ToolInput ToolInput `json:"tool_input"`
}
type ToolInput struct {
    Command string `json:"command"`
}
```

All subcommands exit 0. Block commands output `{"decision":"block","reason":"..."}`.

---

## 2. Plugin Directory

**Location**: `plugin/ctx-plugin/` (distributable artifact, outside `internal/`)

```
plugin/ctx-plugin/
├── .claude-plugin/
│   └── plugin.json
├── skills/                  # 27 generic ctx skills (ctx-* prefixed + borrow/worktree)
│   ├── ctx-status/SKILL.md
│   ├── ctx-add-decision/SKILL.md
│   ├── ctx-add-learning/SKILL.md
│   ├── ctx-add-task/SKILL.md
│   ├── ctx-add-convention/SKILL.md
│   ├── ctx-agent/SKILL.md
│   ├── ctx-alignment-audit/SKILL.md
│   ├── ctx-archive/SKILL.md
│   ├── ctx-blog/SKILL.md
│   ├── ctx-blog-changelog/SKILL.md
│   ├── ctx-borrow/SKILL.md
│   ├── ctx-commit/SKILL.md
│   ├── ctx-context-monitor/SKILL.md
│   ├── ctx-drift/SKILL.md
│   ├── ctx-implement/SKILL.md
│   ├── ctx-journal-enrich/SKILL.md
│   ├── ctx-journal-enrich-all/SKILL.md
│   ├── ctx-journal-normalize/SKILL.md
│   ├── ctx-loop/SKILL.md
│   ├── ctx-next/SKILL.md
│   ├── ctx-pad/SKILL.md
│   ├── ctx-prompt-audit/SKILL.md
│   ├── ctx-recall/SKILL.md
│   ├── ctx-reflect/SKILL.md
│   ├── ctx-remember/SKILL.md
│   └── ctx-worktree/SKILL.md
└── hooks/
    └── hooks.json
```

**Skills NOT in plugin** (project-specific, stay in `.claude/skills/`):
brainstorm, qa, verify, skill-creator, update-docs, release-notes, release,
backup, sanitize-permissions, check-links, consolidate

### plugin.json

```json
{
  "name": "ctx",
  "version": "0.6.0",
  "description": "Persistent context for AI coding assistants",
  "author": {"name": "Context contributors"},
  "homepage": "https://ctx.ist",
  "repository": "https://github.com/ActiveMemory/ctx",
  "license": "Apache-2.0",
  "keywords": ["context", "memory", "persistence", "decisions", "learnings"]
}
```

### hooks.json

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [{"type": "command", "command": "ctx system block-non-path-ctx"}]
      },
      {
        "matcher": ".*",
        "hooks": [{"type": "command", "command": "ctx agent --budget 4000 2>/dev/null || true"}]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Bash",
        "hooks": [{"type": "command", "command": "ctx system post-commit"}]
      }
    ],
    "UserPromptSubmit": [
      {
        "hooks": [
          {"type": "command", "command": "ctx system check-context-size"},
          {"type": "command", "command": "ctx system check-persistence"},
          {"type": "command", "command": "ctx system check-journal"}
        ]
      }
    ],
    "SessionEnd": [
      {
        "hooks": [{"type": "command", "command": "ctx system cleanup-tmp"}]
      }
    ]
  }
}
```

---

## 3. Marketplace

**Location**: repo root `.claude-plugin/marketplace.json`

```json
{
  "name": "activememory-ctx",
  "owner": {"name": "Context contributors"},
  "metadata": {
    "description": "Official ctx plugins for Claude Code",
    "version": "1.0.0"
  },
  "plugins": [
    {
      "name": "ctx",
      "source": "./plugin/ctx-plugin",
      "description": "Persistent context for AI coding assistants",
      "version": "0.6.0"
    }
  ]
}
```

Users install with:
```
/plugin marketplace add ActiveMemory/ctx
/plugin install ctx@activememory-ctx
```

---

## 4. Changes to `ctx init`

### Removes
- `createClaudeHooks()` call (line 151 of `internal/cli/initialize/run.go`)
- `createClaudeSkills()` call (inside `createClaudeHooks` in `hook.go`)
- Hook entries from `mergeSettingsHooks()` — **keep permissions merge only**
- Shell script deployment via `deployHookScript()`

### Keeps
- `.context/` scaffolding (unchanged)
- CLAUDE.md creation/merge (unchanged)
- PROMPT.md, IMPLEMENTATION_PLAN.md (unchanged)
- Makefile.ctx (unchanged)
- .gitignore entries (unchanged)
- `Bash(ctx:*)` permission seeding in settings.local.json (still useful)

### Adds
- End-of-init message: plugin install guidance for Claude Code users

---

## 5. Changes to `ctx hook`

### Remove
- `case "claude-code", "claude":` branch in `internal/cli/hook/run.go`
- Replace with redirect: "Claude Code integration is now provided via the ctx plugin."

### Keep
- `cursor`, `aider`, `copilot`, `windsurf` cases (unchanged)
- Help text updated to remove claude-code from supported tools list

---

## 6. Code Deletion

### Embedded shell scripts (dead code)
- `internal/tpl/claude/hooks/block-non-path-ctx.sh`
- `internal/tpl/claude/hooks/check-context-size.sh`
- `internal/tpl/claude/hooks/check-persistence.sh`
- `internal/tpl/claude/hooks/check-journal.sh`
- `internal/tpl/claude/hooks/post-commit.sh`
- `internal/tpl/claude/hooks/cleanup-tmp.sh`

### Script loaders (`internal/claude/script.go`)
- All 6 functions: `BlockNonPathCtxScript`, `CheckContextSizeScript`,
  `CleanupTmpScript`, `CheckPersistenceScript`, `PostCommitScript`,
  `CheckJournalScript`

### Hook matchers (`internal/claude/matcher.go`)
- `preToolUserHookMatcher`, `postToolUseHookMatcher`,
  `sessionEndHookMatcher`, `userPromptSubmitHookMatcher`

### Hook builder (`internal/claude/hook.go`)
- `DefaultHooks()` function

### Init scaffolding (`internal/cli/initialize/hook.go`)
- `createClaudeHooks()`, `deployHookScript()` — delete
- `mergeSettingsHooks()` — reduce to permissions-only merge

### Constants (`internal/config/file.go`)
- `FileBlockNonPathScript`, `FileCheckContextSize`, `FileCheckPersistence`,
  `FileCheckJournal`, `FilePostCommit`, `FileCleanupTmp`, `CmdAutoloadContext`

### Embed directive (`internal/tpl/embed.go`)
- Remove `claude/hooks/*.sh` from `//go:embed`
- Remove `ClaudeHookByFileName()` function

---

## 7. New File Structure

```
internal/cli/system/
├── doc.go
├── system.go                  # Parent command (Hidden=true)
├── input.go                   # HookInput/ToolInput types, readInput()
├── state.go                   # secureTempDir, counter, log, throttle helpers
├── check_context_size.go      # Adaptive prompt counter
├── check_context_size_test.go
├── check_persistence.go       # Context file mtime watcher
├── check_persistence_test.go
├── check_journal.go           # Unexported sessions + unenriched entries
├── check_journal_test.go
├── block_non_path_ctx.go      # Command pattern blocker
├── block_non_path_ctx_test.go
├── post_commit.go             # Post-commit context capture nudge
├── post_commit_test.go
├── cleanup_tmp.go             # Old temp file removal
└── cleanup_tmp_test.go
```

Registration in `internal/bootstrap/bootstrap.go`: add `system.Cmd` to the command list.

---

## 8. Testing Strategy

### Unit tests per subcommand
- Input parsing: known JSON -> correct field extraction
- Behavior parity: same inputs as shell scripts -> identical outputs
- Edge cases: empty stdin, missing session_id, missing dirs, permission errors

### Key test scenarios
- `check-context-size`: silent at count 5, checkpoint at count 18, checkpoint at count 33
- `check-persistence`: mtime reset when .context/ modified, nudge at prompt 20
- `check-journal`: both stages fire, daily throttle works, no journal dir = silent
- `block-non-path-ctx`: blocks `./ctx`, `go run ./cmd/ctx`, `/home/.../ctx`;
  allows `ctx status`, `git -C ./ctx/path`
- `post-commit`: triggers on `git commit`, skips `git commit --amend`, skips `ls`
- `cleanup-tmp`: removes 16-day-old file, keeps 14-day-old file

### Integration test
Build binary, pipe JSON via `echo '...' | ctx system <cmd>`, verify stdout/exit code.

### Plugin smoke test
`claude --plugin-dir ./plugin/ctx-plugin` — verify hooks fire, skills appear.

---

## 9. Task Breakdown for TASKS.md

### Phase 5: Plugin Conversion (ctx v0.6.0) `#priority:high`

**Context**: Convert ctx from shell-script hook scaffolding to a Claude Code plugin.
Spec: this plan document.

```
P5.0:  Write spec to specs/plugin-conversion.md
P5.1:  Create internal/cli/system/ package scaffold (doc.go, system.go, input.go, state.go)
       Register system.Cmd in bootstrap.go.
       Done: `ctx system --help` works.

P5.2:  Implement ctx system check-context-size
       Done: unit tests pass, output matches shell script behavior.

P5.3:  Implement ctx system check-persistence
       Done: unit tests pass, mtime reset + nudge timing verified.

P5.4:  Implement ctx system check-journal
       Done: unit tests pass, both stages + daily throttle verified.

P5.5:  Implement ctx system block-non-path-ctx
       Done: unit tests pass, all 3 patterns blocked, JSON output correct.

P5.6:  Implement ctx system post-commit
       Done: unit tests pass, amend skip verified.

P5.7:  Implement ctx system cleanup-tmp
       Done: unit tests pass, age threshold verified.

P5.8:  Create plugin/ctx-plugin/ directory
       Copy 27 skills from internal/tpl/claude/skills/.
       Create .claude-plugin/plugin.json and hooks/hooks.json.
       Done: valid plugin structure, `claude plugin validate` passes.

P5.9:  Create .claude-plugin/marketplace.json at repo root
       Done: marketplace validates.

P5.10: Update ctx init — remove Claude scaffolding
       Remove createClaudeHooks, createClaudeSkills.
       Keep permissions merge. Add plugin guidance message.
       Done: ctx init creates .context/ only, no .claude/hooks/*.sh.

P5.11: Update ctx hook — remove claude-code case
       Replace with plugin redirect message.
       Done: ctx hook claude-code prints install instructions.

P5.12: Delete dead code
       Shell scripts, script loaders, matchers, DefaultHooks, constants, embed.
       Done: go build succeeds, go test passes, no dead code.

P5.13: Integration tests for ctx system subcommands
       Done: all pass in CI.

P5.14: Bump version to 0.6.0
       Done: ctx --version shows 0.6.0.

P5.15: Plugin smoke test
       Done: manual verification with claude --plugin-dir.

P5.16: Documentation updates
       New plugin install page, hook internals page, updated getting-started.
       Remove shell script references.
       Done: docs reflect v0.6.0.
```

### Dependency graph

```
P5.0 -> P5.1 -> P5.2-P5.7 (parallel)
                    |
        P5.8 + P5.9 (parallel with P5.2-P5.7)
                    |
        P5.10 + P5.11 (parallel)
                    |
               P5.12 -> P5.13 -> P5.14 -> P5.15 -> P5.16
```

---

## Verification

1. **Build**: `CGO_ENABLED=0 go build -o ctx ./cmd/ctx` — succeeds
2. **Tests**: `CGO_ENABLED=0 go test ./...` — all pass including new system/ tests
3. **Version**: `./ctx --version` -> `0.6.0`
4. **Init**: `./ctx init` in a temp dir -> creates `.context/` only, no `.claude/hooks/*.sh`
5. **System**: `echo '{"session_id":"test"}' | ./ctx system check-context-size` -> exits 0
6. **Plugin**: `claude --plugin-dir ./plugin/ctx-plugin` -> hooks fire, `/ctx:status` works
7. **Hook redirect**: `./ctx hook claude-code` -> prints plugin install message
8. **Other hooks**: `./ctx hook cursor` -> still prints Cursor instructions
