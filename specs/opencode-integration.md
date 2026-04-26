# Spec: OpenCode Integration for ctx

## Context

OpenCode (opencode.ai) is a terminal-first AI coding agent with 140K+ GitHub stars.
It reads `AGENTS.md` natively, supports MCP servers via `opencode.json`, and has a
plugin system (`@opencode-ai/plugin`) with lifecycle hooks. ctx already mentions
OpenCode as an `AGENTS.md`-compatible tool (`hooks.yaml:373`) but has no dedicated
integration.

**Goal:** Add `ctx setup opencode --write` following the Copilot CLI blueprint —
MCP registration, `AGENTS.md` generation, a thin TypeScript plugin (embedded asset,
not an npm dependency) that shims lifecycle hooks to `ctx system` subcommands, and
OpenCode-native skills.

**Why this shape:** Every ctx integration is a Go package that deploys config +
assets. The TypeScript plugin is a static embedded asset (like Copilot CLI's `.sh`
scripts) — not a build dependency. All real logic stays in Go via `ctx system`.

---

## Files to Create

### 1. Embedded Assets (`internal/assets/integrations/opencode/`)

```
internal/assets/integrations/opencode/
├── plugin/
│   ├── index.ts          # Thin shim plugin (~50 lines)
│   └── package.json      # Minimal: name, version, @opencode-ai/plugin dep
└── skills/               # ctx skills in OpenCode format
    ├── ctx-agent/SKILL.md
    ├── ctx-remember/SKILL.md
    ├── ctx-status/SKILL.md
    └── ctx-wrap-up/SKILL.md
```

OpenCode reads `AGENTS.md` natively, so we deploy the shared
`agent.AgentsMd()` template at the project root rather than shipping
an OpenCode-specific instructions file.

**`plugin/index.ts`** — the core deliverable. Wires `session.created`
and `session.idle` to `ctx system` nudges, runs `post-commit` after
shell commands that contain `git commit`, and runs
`check-task-completion` after edit/write tool calls. Tool name strings
target `@opencode-ai/plugin` v1.4.x; unrecognized tools silently
no-op.

We deliberately do **not** ship a `tool.execute.before` hook here:
the natural fit (block-dangerous-commands) is currently a Claude Code
plugin-local hook, not a `ctx system` subcommand, so a shim that
shells out to it would block every shell command on installs that
don't have the wrapper. Add this back when block-dangerous-commands
is promoted to the ctx Go binary.

**`plugin/package.json`**:
```json
{
  "name": "ctx-opencode-plugin",
  "version": "0.1.0",
  "type": "module",
  "main": "index.ts",
  "dependencies": {
    "@opencode-ai/plugin": "^1.4.0"
  }
}
```

**`skills/`** — Subset of portable skills (ctx-agent, ctx-remember, ctx-status,
ctx-wrap-up). Format: YAML frontmatter + markdown body, same as Copilot CLI skills.

### 2. Asset Reader (`internal/assets/read/agent/agent.go`)

Add functions (following `CopilotCLI*` pattern):

```go
// OpenCodePlugin reads the embedded OpenCode plugin directory.
func OpenCodePlugin() (map[string][]byte, error)  // filename -> content

// OpenCodeSkills reads embedded OpenCode skill templates.
func OpenCodeSkills() (map[string][]byte, error)  // skill name -> SKILL.md
```

### 3. Asset Path Constants (`internal/config/asset/asset.go`)

Add:
```go
DirIntegrationsOpenCode       = "integrations/opencode"
DirIntegrationsOpenCodePlugin = "integrations/opencode/plugin"
DirIntegrationsOpenCodeSkill  = "integrations/opencode/skills"
```

### 4. Hook Constants (`internal/config/hook/hook.go`)

Add to tool constants:
```go
ToolOpenCode = "opencode"
```

Add OpenCode integration path constants:
```go
// OpenCode integration paths.
DirOpenCode        = ".opencode"
DirOpenCodePlugins = "plugins"
DirOpenCodeSkills  = "skills"
FileOpenCodeJSON   = "opencode.json"
```

### 5. Setup Path Constants (`internal/config/setup/setup.go`)

Add:
```go
DisplayOpenCode      = "OpenCode"
MCPConfigPathOpenCode = "opencode.json"
PluginPathOpenCode    = ".opencode/plugins/ctx/"
SkillsPathOpenCode    = ".opencode/skills/"
```

### 6. Text Description Keys (`internal/config/embed/text/hook.go`)

Add:
```go
DescKeyHookOpenCode              = "hook.opencode"
DescKeyWriteHookOpenCodeCreated  = "write.hook-opencode-created"
DescKeyWriteHookOpenCodeSkipped  = "write.hook-opencode-skipped"
DescKeyWriteHookOpenCodeSummary  = "write.hook-opencode-summary"
```

### 7. YAML Text Templates

**`hooks.yaml`** — add `hook.opencode`:
```yaml
hook.opencode:
  short: |
    OpenCode Integration
    ====================

    Generate .opencode/plugins/ctx/ with ctx lifecycle hooks
    and register the ctx MCP server in opencode.json.

    This creates:
      .opencode/plugins/ctx/index.ts      Plugin shim
      .opencode/plugins/ctx/package.json   Dependencies
      .opencode/skills/ctx-*/SKILL.md      ctx skills
      opencode.json                        MCP server registration

    Run with --write to generate all files:

      ctx setup opencode --write
```

Update `hook.supported-tools` to include `opencode`.

**`write.yaml`** — add:
```yaml
write.hook-opencode-created:
  short: '  ✓ %s'
write.hook-opencode-skipped:
  short: '  ○ %s (ctx plugin exists, skipped)'
write.hook-opencode-summary:
  short: |-
    OpenCode will now:
      1. Bootstrap ctx context on session start
      2. Nudge persistence on session idle
      3. Track task completion after edits
      4. Run post-commit capture after `git commit`
```

### 8. Setup Core Package (`internal/cli/setup/core/opencode/`)

```
internal/cli/setup/core/opencode/
├── doc.go           # Package documentation
├── opencode.go      # Deploy() entry point
├── plugin.go        # deployPlugin() — writes .opencode/plugins/ctx/
├── mcp.go           # ensureMCPConfig() — merges opencode.json
├── skill.go         # deploySkills() — writes .opencode/skills/
└── agents.go        # deployAgents() — writes AGENTS.md (shared template)
```

**`opencode.go` — Deploy()**:
```go
func Deploy(cmd *cobra.Command) error {
    // 1. Deploy plugin files (.opencode/plugins/ctx/)
    if pluginErr := deployPlugin(cmd); pluginErr != nil {
        return pluginErr
    }
    // 2. Register MCP server in opencode.json
    if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
        writeErr.WarnFile(cmd, cfgSetup.MCPConfigPathOpenCode, mcpErr)
    }
    // 3. Deploy AGENTS.md (shared template, idempotent)
    if agentsErr := deployAgents(cmd); agentsErr != nil {
        writeErr.WarnFile(cmd, cfgHook.FileAgentsMd, agentsErr)
    }
    // 4. Deploy skills to .opencode/skills/
    if skillErr := deploySkills(cmd); skillErr != nil {
        writeErr.WarnFile(cmd, cfgSetup.SkillsPathOpenCode, skillErr)
    }
    writeSetup.InfoOpenCodeSummary(cmd)
    return nil
}
```

**`mcp.go` — ensureMCPConfig()**:

OpenCode MCP config lives in `opencode.json` at project root:
```json
{
  "mcp": {
    "ctx": {
      "type": "local",
      "command": "ctx",
      "args": ["mcp", "serve"]
    }
  }
}
```

Read-merge-write pattern: read existing `opencode.json`, add/update `mcp.ctx`
entry, write back. Preserve all other config keys.

**`plugin.go` — deployPlugin()**:

Extract embedded `index.ts` and `package.json` to `.opencode/plugins/ctx/`.
Skip if `index.ts` already exists (idempotent). OpenCode auto-runs
`bun install` in plugin directories at startup.

### 9. Write Setup Functions (`internal/write/setup/hook.go`)

Add:
```go
func InfoOpenCodeCreated(cmd *cobra.Command, targetFile string)
func InfoOpenCodeSkipped(cmd *cobra.Command, targetFile string)
func InfoOpenCodeSummary(cmd *cobra.Command)
```

### 10. Tool Dispatcher (`internal/cli/setup/cmd/root/run.go`)

Add import and case:
```go
coreOpenCode "github.com/ActiveMemory/ctx/internal/cli/setup/core/opencode"

case cfgHook.ToolOpenCode:
    if writeFile {
        return coreOpenCode.Deploy(cmd)
    }
    writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookOpenCode))
```

---

## Files to Modify

| File | Change |
|------|--------|
| `internal/config/hook/hook.go` | Add `ToolOpenCode` + OpenCode path constants |
| `internal/config/setup/setup.go` | Add `DisplayOpenCode` + path constants |
| `internal/config/asset/asset.go` | Add `DirIntegrationsOpenCode*` constants |
| `internal/config/embed/text/hook.go` | Add `DescKeyHookOpenCode` + write keys |
| `internal/assets/commands/text/hooks.yaml` | Add `hook.opencode` + update supported-tools |
| `internal/assets/commands/text/write.yaml` | Add `write.hook-opencode-*` entries |
| `internal/assets/read/agent/agent.go` | Add `OpenCode*()` reader functions |
| `internal/write/setup/hook.go` | Add `InfoOpenCode*()` output functions |
| `internal/cli/setup/cmd/root/run.go` | Add `opencode` case to switch |
| `docs/operations/integrations.md` | Add OpenCode section |

---

## What We're NOT Doing

- No steering sync for OpenCode (OpenCode doesn't have a native rules format
  like `.cursor/rules/`; it uses `AGENTS.md` + skills instead)
- No `ctx init` changes (OpenCode reads `AGENTS.md` which `ctx setup agents`
  already generates; `ctx setup opencode` handles the rest)
- No npm publish (plugin is embedded in the Go binary, deployed by setup)
- No session parser (OpenCode session format is TBD; add later when stable)

---

## Implementation Order

1. **Constants** — `hook.go`, `setup.go`, `asset.go`, `text/hook.go`
2. **Embedded assets** — `integrations/opencode/` directory with plugin, instructions, skills
3. **Asset readers** — `agent.go` OpenCode functions
4. **Setup core** — `internal/cli/setup/core/opencode/` package (5 files)
5. **Write functions** — `write/setup/hook.go` additions
6. **Dispatcher** — `run.go` case addition
7. **YAML text** — `hooks.yaml` + `write.yaml` entries
8. **Docs** — `integrations.md` OpenCode section

---

## Verification

1. **Build**: `make build` — verify compilation with new package
2. **Dry run**: `ctx setup opencode` — should print integration instructions
3. **Write**: `ctx setup opencode --write` in a test project — verify:
   - `.opencode/plugins/ctx/index.ts` created
   - `.opencode/plugins/ctx/package.json` created
   - `opencode.json` has `mcp.ctx` entry (merged, not overwritten)
   - `AGENTS.md` created (or skipped if exists with markers)
   - `.opencode/skills/ctx-*/SKILL.md` created
4. **Idempotency**: Run `ctx setup opencode --write` twice — second run skips existing
5. **Lint**: `make lint`
6. **Test**: `make test`
7. **Smoke**: `make smoke`
