//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

// DefaultLastFile is the default number of commits shown by ctx trace file.
const DefaultLastFile = 20

// DefaultLastShow is the default number of commits shown by ctx trace
// when invoked without arguments.
const DefaultLastShow = 10

// ShortHashLen is the number of characters used when abbreviating
// a git commit hash for display.
const ShortHashLen = 7

// Hook action argument values for ctx trace hook <enable|disable>.
const (
	ActionEnable  = "enable"
	ActionDisable = "disable"
)

// TrailerKey is the git trailer key used to embed context
// refs in commit messages.
const TrailerKey = "ctx-context"

// TrailerFormat is the format string for the git trailer line: "key: refs".
const TrailerFormat = TrailerKey + ": %s"

// Reference type identifiers used in ctx-context trailers.
const (
	RefTypeNote       = "note"
	RefTypeSession    = "session"
	RefTypeDecision   = "decision"
	RefTypeLearning   = "learning"
	RefTypeConvention = "convention"
	RefTypeTask       = "task"
)

// Task status labels for resolved refs.
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
)

// RefFirstEntry is the suffix for the first entry in a context file.
const RefFirstEntry = ":1"

// RefFormat is the format string for numbered refs (e.g. "decision:1").
const RefFormat = "%s:%d"

// SessionRefFormat is the format string for session refs
// (e.g. "session:abc123").
const SessionRefFormat = "session:%s"

// Diff line prefix constants for parsing git diff output.
const (
	DiffAddedPrefix  = "+"
	DiffHeaderPrefix = "++"
)

// CtxTraceMarker is the string used to identify ctx-installed git hooks.
const CtxTraceMarker = "ctx trace"

// JSONL storage filenames within the trace and state directories.
const (
	FileHistory   = "history.jsonl"
	FileOverrides = "overrides.jsonl"
	FilePending   = "pending-context.jsonl"
)

// Embedded hook script filenames.
const (
	ScriptPrepareCommitMsg = "prepare-commit-msg.sh"
	ScriptPostCommit       = "post-commit.sh"
)
