// ctx OpenCode plugin — thin shim to ctx system subcommands.
// All real logic lives in the ctx Go binary; this plugin just
// wires OpenCode lifecycle hooks to ctx system calls.
//
// Hook signatures match @opencode-ai/plugin v1.4.x:
//   - shell.env and tool.execute.after take (input, output)
//     and mutate output rather than returning a value.
//   - event is a single dispatcher keyed on input.event.type;
//     it is NOT an object of named per-event handlers.
// ctx subprocess calls go through a CTX_DIR-aware BunShell built
// from ctx.directory — shell.env only injects CTX_DIR into the
// agent's shell tool, not into the plugin's own ctx.$ calls.
// If the upstream renames a hook or changes a signature, the
// corresponding branch silently no-ops; verify against the
// OpenCode plugin SDK type definitions when bumping.
import type { Plugin } from "@opencode-ai/plugin"

const SHELL_TOOLS = new Set(["shell", "bash"])
const EDIT_TOOLS = new Set(["edit", "write", "file_edit"])
// Match `git commit` but not `git commit-tree` / `git commit-graph`.
// The negative lookahead rejects `-` immediately after the boundary.
const GIT_COMMIT_RE = /\bgit\s+commit\b(?!-)/

function extractCommand(input: unknown): string {
  if (typeof input === "string") return input
  if (input && typeof input === "object") {
    const cmd = (input as { command?: unknown }).command
    if (typeof cmd === "string") return cmd
  }
  return ""
}

export default (async (ctx) => {
  const ctxDir = `${ctx.directory}/.context`
  const $ = ctx.$.env({ ...process.env, CTX_DIR: ctxDir })
  return {
    "shell.env": async (input, output) => {
      output.env.CTX_DIR = `${input.cwd}/.context`
    },
    event: async ({ event }) => {
      if (event.type === "session.created") {
        await $`ctx system bootstrap 2>/dev/null || true`
        await $`ctx agent --budget 4000 2>/dev/null || true`
      } else if (event.type === "session.idle") {
        await $`ctx system check-persistence 2>/dev/null || true`
        await $`ctx system check-task-completion 2>/dev/null || true`
      }
    },
    "tool.execute.after": async (input, _output) => {
      if (SHELL_TOOLS.has(input.tool)) {
        const cmd = extractCommand(input.args)
        if (GIT_COMMIT_RE.test(cmd)) {
          await $`ctx system post-commit 2>/dev/null || true`
        }
      }
      if (EDIT_TOOLS.has(input.tool)) {
        await $`ctx system check-task-completion 2>/dev/null || true`
      }
    },
  }
}) satisfies Plugin
