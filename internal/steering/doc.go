//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package steering parses, scores, and synchronizes
// **steering files** — the small frontmattered Markdown
// documents under `.context/steering/` that tell each
// configured AI tool *how to behave* when a specific kind of
// prompt arrives.
//
// Steering is the declarative half of ctx's behavior layer
// (the imperative half is [internal/trigger]: scripts that
// *do* things on lifecycle events). A steering file says
// "when the user asks about Y, prepend these rules to the
// prompt"; a trigger says "when X happens, run this code."
//
// # The Steering File
//
// Each `.md` file under the steering directory is a
// [SteeringFile]: a short YAML frontmatter block followed by
// a Markdown body. The schema:
//
//   - **name** — unique identifier; matches the manual
//     selector in `ctx steering preview --names ...`.
//   - **description** — one-line summary; doubles as the
//     match phrase for [cfgSteering.InclusionAuto].
//   - **inclusion** — `always` | `auto` | `manual`
//     ([cfgSteering.InclusionMode]). Default `manual`.
//   - **tools** — list of AI-tool IDs the file applies to;
//     empty/nil means "all tools".
//   - **priority** — injection order; lower priority is
//     injected earlier (default 50).
//
// [Parse] reads bytes + a path and returns a fully populated
// [SteeringFile] with defaults applied; YAML errors are wrapped
// via [internal/err/steering] so the file path is always part
// of the message. [LoadAll] is the bulk variant that walks a
// directory.
//
// # The Inclusion Modes
//
// Three modes determine when a file's body is appended to the
// next prompt:
//
//   - **always** — every prompt, every turn, no questions.
//     Heaviest on context budget; reserve for genuinely
//     foundational rules.
//   - **auto** — included when the lowercased prompt contains
//     the file's lowercased description (substring match —
//     simple, deterministic, fast). The most common mode for
//     project-specific guidance.
//   - **manual** — only when the file's name appears in the
//     `manualNames` argument to [Filter] / [matchInclusion].
//     Used by `ctx steering preview --names ...` and by the
//     MCP `steering_get` tool.
//
// [matchInclusion] does the per-file decision; [matchTool]
// adds tool-scope filtering on top. [Filter] composes the two
// against a list of files for a given (prompt, tool, manual)
// triple.
//
// # Two Tool Families, Two Delivery Paths
//
// Not every AI editor consumes steering the same way; ctx
// handles two families:
//
//   - **Native-rules tools** — Cursor, Cline, Kiro
//     ([syncableTools]) — have a built-in rules primitive
//     (`.cursor/rules/*.mdc`, `.clinerules`,
//     `.kiro/steering/*.md`). [SyncTool] writes
//     ctx-managed `.context/steering/*.md` into each tool's
//     native format. [SyncAll] does this for every supported
//     tool in one call. Idempotent: unchanged content is
//     skipped.
//   - **Hook-driven tools** — Claude Code and Codex use
//     `ctx agent` to assemble the context packet on every
//     prompt; their steering arrives via the agent pipeline
//     (no file sync). They are deliberately **not** in
//     [syncableTools]; calling `SyncTool` for them returns
//     [errSteering.UnsupportedTool].
//
// Mixed setups (project uses both Cursor and Claude Code)
// run `ctx steering sync` for the native-rules tools and let
// the hook+MCP pipeline cover Claude Code automatically. See
// `docs/home/steering.md` for the user-facing summary of this
// split.
//
// # Foundation Files
//
// `ctx init` scaffolds four foundation steering files —
// `product`, `tech`, `structure`, `workflow` — so users have
// real templates to edit instead of an empty directory.
// [FoundationFiles] returns the set; bodies and descriptions
// come from YAML text assets at call time so they stay in sync
// with the embedded copy. Re-running `ctx init` is safe:
// existing files are left alone.
//
// # Format Adapters
//
// Each native tool needs a slightly different frontmatter
// shape:
//
//   - [cursorFrontmatter] — `description`, `globs`,
//     `alwaysApply`.
//   - [kiroFrontmatter] — `name`, `description`, `mode`.
//   - Cline takes plain Markdown with no frontmatter.
//
// [format.go] holds the per-tool serializers; the unexported
// types in [types.go] keep the YAML shape decoupled from the
// canonical [SteeringFile].
//
// # Concurrency and Idempotency
//
// Functions are stateless. [SyncTool] reads from the steering
// directory, computes the desired output for each file,
// compares it to what is on disk, and writes only the
// changed files — so running it twice in a row produces no
// `Written` entries the second time, just `Skipped`. Output
// paths are validated to resolve within `projectRoot` before
// writing.
//
// # Related Packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/steering] —
//     CLI surface: `add`, `list`, `preview`, `init`, `sync`.
//   - [github.com/ActiveMemory/ctx/internal/mcp/handler] —
//     MCP `steering_get` tool that surfaces matched files to
//     Claude Code via JSON-RPC.
//   - [github.com/ActiveMemory/ctx/internal/config/steering]
//     — inclusion modes and foundation file names.
//   - [github.com/ActiveMemory/ctx/internal/config/hook] —
//     supported tool ID constants used by [syncableTools].
//   - [github.com/ActiveMemory/ctx/internal/err/steering] —
//     typed error constructors with file-path context.
//   - [github.com/ActiveMemory/ctx/internal/parse] —
//     [SplitFrontmatter] used by [Parse].
package steering
