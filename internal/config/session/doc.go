//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package session defines constants for parsing,
// displaying, and exporting AI coding session data
// across multiple tools.
//
// ctx supports session transcripts from Claude Code,
// VS Code Copilot Chat, GitHub Copilot CLI, and raw
// Markdown files. This package provides the shared
// vocabulary: tool identifiers, YAML frontmatter
// field names, tool display keys, and session ID
// conventions.
//
// # Tool Identifiers
//
//   - [ToolClaudeCode]  — Claude Code JSONL sessions.
//   - [ToolCopilot]     — VS Code Copilot Chat.
//   - [ToolCopilotCLI]  — GitHub Copilot CLI.
//   - [ToolMarkdown]    — plain Markdown transcripts.
//
// # Claude Code Tool Names
//
// Constants like [ToolRead], [ToolWrite], [ToolEdit],
// [ToolBash], [ToolGrep], and [ToolGlob] match the
// tool names found in Claude Code session transcripts
// and are used for display formatting and filtering.
//
// # Frontmatter Fields
//
// YAML frontmatter keys for journal entries:
// [FrontmatterTitle], [FrontmatterDate],
// [FrontmatterType], [FrontmatterOutcome],
// [FrontmatterTopics], [FrontmatterTechnologies],
// [FrontmatterKeyFiles], [FrontmatterLocked].
// Export-specific keys like [FmKeyTime],
// [FmKeyProject], [FmKeyModel] support the journal
// export pipeline.
//
// # Session Conventions
//
//   - [IDUnknown] — fallback when a session has no ID.
//   - [IDSuffixSummary], [IDSuffixTopic] — suffixes
//     for derived session IDs.
//   - [PreviewMaxLen] — truncation limit for first-
//     message previews.
//   - [DefaultFilename] — fallback for sanitized
//     filenames.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package session
