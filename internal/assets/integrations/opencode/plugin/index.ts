// ctx OpenCode plugin — thin shim to ctx system subcommands.
// All real logic lives in the ctx Go binary; this plugin just
// wires OpenCode lifecycle hooks to ctx system calls.
import type { Plugin } from "@opencode-ai/plugin"

export default ((ctx) => ({
  "shell.env": () => ({
    CTX_DIR: ".context",
  }),
  event: {
    "session.created": async () => {
      await ctx.$`ctx system bootstrap 2>/dev/null || true`
      await ctx.$`ctx agent --budget 4000 2>/dev/null || true`
    },
    "session.idle": async () => {
      await ctx.$`ctx system check-persistence 2>/dev/null || true`
      await ctx.$`ctx system check-task-completion 2>/dev/null || true`
    },
  },
  "tool.execute.before": async ({ tool, input }) => {
    if (tool === "shell" || tool === "bash") {
      const cmd = typeof input === "string" ? input : JSON.stringify(input)
      const result =
        await ctx.$`echo ${cmd} | ctx system block-dangerous-commands --caller opencode 2>/dev/null`
      if (result.exitCode !== 0) {
        return { blocked: true, reason: result.stdout.toString().trim() }
      }
    }
  },
  "tool.execute.after": async ({ tool }) => {
    if (tool === "shell" || tool === "bash") {
      await ctx.$`ctx system post-commit 2>/dev/null || true`
    }
    if (tool === "edit" || tool === "write" || tool === "file_edit") {
      await ctx.$`ctx system check-task-completion 2>/dev/null || true`
    }
  },
})) satisfies Plugin
