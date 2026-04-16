//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parser auto-detects and parses AI-coding-assistant session
// transcripts from multiple tool formats into a single normalized
// [github.com/ActiveMemory/ctx/internal/entity.Session] type the rest
// of the journal pipeline can consume uniformly.
//
// # Why a Parser Layer
//
// Every AI tool ctx integrates with stores its session history in a
// different on-disk format:
//
//   - **Claude Code** writes one JSONL file per project under
//     `~/.claude/projects/<slug>/*.jsonl`, with multiple sessions
//     interleaved by `sessionId` field.
//   - **Copilot** (VS Code) keeps a binary-ish chunked store in
//     the workspace state directory.
//   - **Copilot CLI** writes a different, JSON-with-metadata layout
//     under its own home tree.
//   - **MarkdownSession** is the round-trip format ctx itself
//     produces when an enriched journal entry is *re-imported*; it
//     parses the YAML frontmatter + body that
//     `ctx journal import` produced earlier.
//
// Downstream consumers (`ctx journal source`, `ctx journal import`,
// the journal site builder, the obsidian exporter) should never
// have to know which tool wrote a file. They get back
// `[]*entity.Session` and work with that.
//
// # Public Surface
//
// Three entry points cover the common use cases:
//
//   - [ParseFile](path):                parse one file; returns all
//     sessions it contains (a JSONL file may interleave many).
//   - [ScanDirectory](dir):             recursively walk a tree,
//     parse every parseable file, return sessions sorted
//     newest-first; per-file errors are swallowed so one bad file
//     does not abort the scan.
//   - [ScanDirectoryWithErrors](dir):   same walk, but also
//     returns a slice of (path, err) pairs for every parse failure
//     so callers can surface them to the user.
//
// Tool-specific constructors ([NewClaudeCode], [NewCopilot],
// [NewCopilotCLI], [NewMarkdownSession]) are exported for callers
// that need to operate on a known format directly (tests, format
// converters, the schema validator).
//
// # Dispatch Mechanism
//
// All tool implementations satisfy the unexported `Session`
// interface (Tool, Matches, ParseFile, ParseLine). The package-level
// `registeredParsers` slice holds one instance of each. Dispatch is
// first-match-wins: [ParseFile] iterates the slice and asks each
// parser whether it `Matches(path)`. Implementations may check
// extension, directory shape, or peek at the first line; order in
// the slice matters when a file could plausibly match more than one
// (in practice, the four formats are disjoint).
//
// **Adding a new tool**: implement the four interface methods on a
// new type, then append a constructor call to `registeredParsers`
// in `parser.go`. No other changes are required.
//
// # Output Shape
//
// Every parser yields `*entity.Session` values populated with:
//
//   - identity: ID, Slug, Tool, SourceFile
//   - context: CWD, Project (basename of CWD), GitBranch
//   - timing: StartTime, EndTime, Duration
//   - content: a flat []Message in chronological order
//   - rollups: TurnCount, FirstUserMsg (preview, truncated at
//     [config/session.PreviewMaxLen])
//
// [ScanDirectory] sorts the aggregated slice by `StartTime`
// descending so the most recent session lands at index 0, the
// invariant the journal CLI and site generator both rely on.
//
// # Error Handling
//
// Errors fall into three buckets:
//
//   - **No matching parser**: [ParseFile] returns
//     [internal/err/parser.NoMatch] when no registered parser claims
//     the file. Callers should treat this as "skip", not "fail";
//     the directory may legitimately contain unrelated files.
//   - **Per-file parse errors**: malformed JSON, truncated stream,
//     unexpected schema. [ScanDirectory] swallows these silently;
//     [ScanDirectoryWithErrors] surfaces them paired with the path
//     for the caller to log.
//   - **Filesystem errors**: walk-time IO errors (permission,
//     device) are returned directly from the Scan functions and
//     terminate the walk.
//
// # Concurrency
//
// Parsers are stateless; the `registeredParsers` slice is read-only
// after package init. A single parser instance is reused across all
// calls. Concurrent [ParseFile] / [ScanDirectory] calls are safe.
package parser
