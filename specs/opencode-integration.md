# Spec: OpenCode Integration for ctx

## Context

OpenCode (opencode.ai) is a terminal-first AI coding agent with 140K+ GitHub stars.
It reads `AGENTS.md` natively, supports MCP servers via its global
`~/.config/opencode/opencode.json` config (or `$OPENCODE_HOME/opencode.json`), and has a
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
│   └── index.ts          # Thin shim plugin (~80 lines)
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
shell commands that contain `git commit`, runs `check-task-completion`
after edit/write tool calls, and injects `ctx system bootstrap` output
into the compaction prompt via `experimental.session.compacting` so
ctx context survives session compaction. The compaction hook pushes
to `output.context` (additive) rather than replacing `output.prompt`,
so it composes with other compaction-aware plugins like oh-my-openagent.
`session.created` does not visibly inject the `ctx agent` packet into chat
because the event hook has no output channel; it prepares ctx in the
background so tools and compaction hooks can use it on demand. Tool name
strings target `@opencode-ai/plugin` v1.4.x; unrecognized tools silently
no-op.

We deliberately do **not** ship a `tool.execute.before` hook here:
the natural fit (block-dangerous-commands) is a Claude Code
plugin-local hook, not a `ctx system` subcommand, so a shim that
shells out to it would block every shell command on installs that
don't have the wrapper. Promoting it to a ctx Go subcommand was
considered and declined — see `.context/DECISIONS.md` entry
`2026-04-26-231517`. The omission is permanent, not deferred.

**Deployment layout**: OpenCode auto-loads top-level `.ts`/`.js`
files under `.opencode/plugins/`; subdirectories are NOT scanned.
The setup deploys a single flat file at `.opencode/plugins/ctx.ts`.
No `package.json` is needed — the plugin uses a type-only import
of `@opencode-ai/plugin` (erased at compile time) and the host
runtime provides the plugin context, so there's no runtime
dependency tree to install.

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
MCPConfigPathOpenCode = "~/.config/opencode/opencode.json"
PluginPathOpenCode    = ".opencode/plugins/ctx.ts"
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

    Generate .opencode/plugins/ctx.ts with ctx lifecycle hooks
    and register the ctx MCP server in the global OpenCode config (`~/.config/opencode/opencode.json` or `$OPENCODE_HOME/opencode.json`).

    This creates:
      .opencode/plugins/ctx.ts             Plugin shim
      .opencode/skills/ctx-*/SKILL.md      ctx skills
      ~/.config/opencode/opencode.json    MCP server registration (or $OPENCODE_HOME/opencode.json)

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
      1. Bootstrap ctx in the background on session start
      2. Nudge persistence on session idle
      3. Track task completion after edits
      4. Run post-commit capture after `git commit`
```

### 8. Setup Core Package (`internal/cli/setup/core/opencode/`)

```
internal/cli/setup/core/opencode/
├── doc.go           # Package documentation
├── opencode.go      # Deploy() entry point (delegates AGENTS to core/agents)
├── plugin.go        # deployPlugin() — writes .opencode/plugins/ctx.ts
├── mcp.go           # ensureMCPConfig() — merges global OpenCode config
└── skill.go         # deploySkills() — writes .opencode/skills/
```

**`opencode.go` — Deploy()**:
```go
func Deploy(cmd *cobra.Command) error {
    // 1. Deploy the plugin file (.opencode/plugins/ctx.ts)
    if pluginErr := deployPlugin(cmd); pluginErr != nil {
        return pluginErr
    }
    // 2. Register MCP server in the global OpenCode config
    if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
        writeErr.WarnFile(cmd, cfgSetup.MCPConfigPathOpenCode, mcpErr)
    }
    // 3. Deploy AGENTS.md via the shared agents package
    if agentsErr := coreAgents.Deploy(cmd); agentsErr != nil {
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

OpenCode MCP config lives in the global config file `~/.config/opencode/opencode.json` (or `$OPENCODE_HOME/opencode.json`). Per
the `@opencode-ai/sdk` `McpLocalConfig` schema, `command` is a
single string array holding the binary and its arguments (no
separate `args` field) and `enabled` is required:
```json
{
  "mcp": {
    "ctx": {
      "type": "local",
      "command": ["sh", "-c", "exec env CTX_DIR=\"$PWD/.context\" /abs/path/to/ctx mcp serve"],
      "enabled": true
    }
  }
}
```

Read-merge-write pattern: read the existing global OpenCode config, add/update `mcp.ctx`
entry, write back. Preserve all other config keys. The setup resolves `ctx` to an
absolute binary path via `exec.LookPath` so OpenCode can spawn it from non-interactive shells.

**`plugin.go` — deployPlugin()**:

Write the embedded `index.ts` content to a flat
`.opencode/plugins/ctx.ts` file. Skip when the installed ctx-managed
plugin already matches the embedded content; refresh it in place when
stale. OpenCode only auto-loads top-level files under
`.opencode/plugins/`; subdirectories are NOT scanned, which is why
the deployment is a single flat file rather than a directory.
No `package.json` is shipped — the plugin uses a type-only import
of `@opencode-ai/plugin` and the host runtime provides the plugin
context, so there's no runtime dependency tree to install.

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
   - `.opencode/plugins/ctx.ts` created (flat file; no subdirectory)
   - `~/.config/opencode/opencode.json` (or `$OPENCODE_HOME/opencode.json`) has `mcp.ctx` entry (merged, not overwritten),
     with `command` as a string array and `enabled: true`
   - `AGENTS.md` created, merged if it exists without ctx markers, or skipped if ctx markers already exist
   - `.opencode/skills/ctx-*/SKILL.md` created
   - Confirm OpenCode actually loads the plugin: launch
     `opencode --print-logs --log-level DEBUG` in the test project,
     ask the agent to make an edit and run `git commit`, and verify
     the plugin's `tool.execute.after` fires the `ctx system
     post-commit` and `ctx system check-task-completion` nudges
4. **Idempotency**: Run `ctx setup opencode --write` twice — second run skips already-matching managed files and refreshes stale plugin/skill/MCP installs
5. **Lint**: `make lint`
6. **Test**: `make test`
7. **Smoke**: `make smoke`
