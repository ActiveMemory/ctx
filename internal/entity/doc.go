//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entity is the **shared domain-types layer** for ctx —
// the structs that carry information across package boundaries
// without coupling the producer to the consumer.
//
// # The Two-Sentence Rule
//
// A type belongs in `entity` if and only if **at least two
// packages need it AND it carries no I/O or business logic**.
// Single-package types live with their consumer; types with
// methods that touch the filesystem, the network, or the
// process environment live with the package that performs that
// I/O. Entity types are pure data carriers.
//
// This rule keeps the package free of import cycles: every
// package may depend on `entity`, but `entity` depends on
// nothing except the standard library and a handful of typed
// configuration constants.
//
// # File Layout — One Domain per File
//
// Types are grouped by **the subsystem that owns the data**,
// not by their Go shape. A non-exhaustive tour:
//
//   - **`context.go`**    — [Context], the assembled
//     `.context/` snapshot every reader sees: file list, token
//     stats, drift signals.
//   - **`add.go`**        — [EntryParams], [AddConfig],
//     [EntryOpts] for the `ctx add` family.
//   - **`change.go`**     — [ContextChange], [CodeSummary] for
//     `ctx change`.
//   - **`message.go`**    — [Message], [ToolUse], [ToolResult]
//     — the normalized session-message shape produced by
//     [internal/journal/parser] and consumed everywhere
//     downstream.
//   - **`session.go`**    — the [Session] aggregate (start /
//     end / duration / project / branch / messages / rollups)
//     that flows from parser → journal pipeline → site /
//     obsidian renderers.
//   - **`journal.go`**    — [JournalEntry], [JournalFrontmatter]
//     — the on-disk shape of an enriched journal entry.
//   - **`import.go`**     — [ImportPlan], [ImportResult],
//     [FileAction], [RenameOp] used by the journal-import
//     pipeline.
//   - **`index.go`**      — [IndexEntry], [GroupedIndex],
//     [TopicData], [KeyFileData], [TypeData] — index-table
//     primitives consumed by `internal/index`.
//   - **`hook.go`**       — [HookInput], [ToolInput],
//     [BlockResponse] — the payload shapes for ctx-system
//     hook plumbing.
//   - **`trigger.go`**    — [TriggerSession], [TriggerInput]
//     — payloads for project-authored lifecycle scripts (see
//     [internal/trigger]).
//   - **`system.go`**     — system-hook input/output types.
//   - **`event.go`**      — [EventQueryOpts] and event log
//     types used by `ctx hook event`.
//   - **`notify.go`**     — [NotifyPayload], [TemplateRef] —
//     the webhook delivery payloads.
//   - **`task.go`**       — task-related domain types
//     (priority, completion state, snapshot shapes).
//   - **`mcp_session.go`**, **`mcp_deps.go`**,
//     **`mcp_prompt.go`** — the per-session state, runtime
//     dependency container, and prompt-spec types passed
//     between the MCP server and its handler package.
//   - **`bootstrap.go`**  — [BootstrapOutput] — the JSON
//     emitted by `ctx system bootstrap` for AI agents at
//     session start.
//   - **`deploy.go`**, **`merge.go`** — pipeline params for
//     deploy/merge orchestration.
//   - **`meta.go`**       — [Stats], [TokenInfo] — rollup
//     metadata attached to many other types.
//
// New types should slot into the file whose subsystem owns
// the data; create a new file only when a genuinely new
// subsystem appears.
//
// # The "No Behavior" Constraint
//
// Methods on entity types are limited to:
//
//   - **Pure predicates** (e.g. `Message.BelongsToUser()`,
//     `Task.Done()`) — they read fields and return derived
//     facts.
//   - **Pure derivations** (e.g. `Session.Duration()`).
//   - **Display helpers** (e.g. `String()` overrides for
//     debug logging).
//
// Anything that needs to read a file, hit the network, run
// `git`, or call `time.Now` belongs in the consumer package.
// `entity` types are safe to construct in tests with literal
// field assignment; they need no constructor and have no
// hidden state.
//
// # Concurrency
//
// All entity types are **plain data**. Mutation is the
// responsibility of whoever owns the value; concurrent
// readers of an immutable value are safe by definition.
// Several types (notably [MCPSession]) are documented as
// owned by a single goroutine at a time and must not be
// shared across requests; see their per-type doc comments
// for the full contract.
//
// # Related Packages
//
// Producers and consumers (non-exhaustive):
//
//   - [internal/context/load]      — assembles [Context].
//   - [internal/journal/parser]    — produces [Session],
//     [Message], [ToolUse], [ToolResult].
//   - [internal/index]             — produces [IndexEntry]
//     and the grouped variants.
//   - [internal/drift]             — consumes [Context].
//   - [internal/mcp/handler]       — consumes [MCPDeps],
//     [MCPSession], produces [PromptEntrySpec].
//   - [internal/trigger]           — consumes [TriggerInput]
//     and emits trigger output.
//   - [internal/notify]            — consumes [NotifyPayload].
package entity
