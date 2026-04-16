//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook holds the **constants** every package
// touching the hook subsystem references: hook names
// (`check_persistence`, `check_context_size`, …),
// lifecycle stages (sessionStart, preToolUse, …),
// supported AI tool identifiers (`claude`, `cursor`,
// `cline`, `kiro`, `codex`), category tags
// (`Customizable`, `CtxSpecific`), and the per-hook
// throttling thresholds.
//
// The package is a constants registry — no logic. Its
// existence keeps consumers free of magic strings and
// lets the audit suite catch references to non-existent
// hook names at compile time.
//
// # Constant Families
//
//   - **Hook names** — one per `cmd/check_*` and
//     `cmd/block_*` package under
//     [internal/cli/system/cmd]. Used by the
//     `ctx hook event` query layer and the message
//     loader.
//   - **Lifecycle stages** — `sessionStart`,
//     `sessionEnd`, `preToolUse`, `postToolUse`,
//     `userPromptSubmit`, etc. The Claude Code
//     hook config and the trigger dispatcher both
//     speak this vocabulary.
//   - **Tool identifiers** — `ToolClaude`,
//     `ToolCursor`, `ToolCline`, `ToolKiro`,
//     `ToolCodex` — the `tool:` field in `.ctxrc`
//     and the `tools:` filter in steering files
//     reference these.
//   - **Categories** — `CategoryCustomizable`,
//     `CategoryCtxSpecific` — used by
//     [internal/assets/hooks/messages] to label
//     each message in `ctx hook message list`.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package hook
