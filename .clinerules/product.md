# product


# Product Context

`ctx` is **persistent context for AI coding sessions**. It gives
the AI memory across sessions by writing project state to
git-versioned Markdown files in `.context/` and feeding that
state back to the AI on every turn.

## Target users

Developers using AI coding tools (Claude Code, Cursor, OpenCode,
Copilot CLI, Aider, Cline, Kiro, Codex) who want their AI to
remember decisions, conventions, and learnings across sessions
without re-explaining the project every time.

## Load-bearing constraints

These shape every design decision; treat them as invariants when
proposing features:

- **Local-first.** All state lives in the user's filesystem. No
  hosted service, no cloud account, no network call required for
  normal operation.
- **Single statically-linked binary.** No runtime dependency
  tree, no package manager, no install step beyond "drop the
  binary on PATH."
- **Git-friendly.** Context is plain Markdown with stable
  ordering; diffs are human-readable. Designed so context
  history lives in the same repo as the code it describes.
- **Tool-agnostic.** ctx integrates with multiple AI tools as
  symmetric peers. No tool is the "primary"; new tools land via
  the same `ctx setup <tool>` and `ctx steering sync` paths.
- **No telemetry, no anonymous data collection.** Period.

## Out of scope

- Cloud-hosted state, SaaS sync, or any solution that requires a
  network round-trip during normal use. If a proposal needs a
  server, it's the wrong proposal for ctx.
- Embedding an LLM into ctx. ctx is the persistence layer; the
  LLM lives in the user's chosen AI tool.
- AI-tool lock-in. Features must work across at least two of the
  supported tool families (hook-based + native-rules), not be
  Claude-Code-only or Cursor-only by design.
