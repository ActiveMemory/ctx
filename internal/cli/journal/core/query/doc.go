//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package query provides session discovery for journal
// operations.
//
// Before the journal can import or list sessions, it
// needs to locate Claude Code session files on disk.
// This package wraps the lower-level session parser to
// provide project-scoped or global session discovery.
//
// # Session Discovery
//
// [FindSessions] is the sole exported function. When
// allProjects is false, it resolves the current working
// directory and delegates to parser.FindSessionsForCWD,
// returning only sessions whose project path matches.
// When allProjects is true, it calls
// parser.FindSessions to scan all known session
// directories regardless of project.
//
// The returned sessions are sorted by start time by the
// underlying parser. The cmd/journal layer passes them
// to the plan package for import planning or to the
// write layer for listing.
//
// # Error Handling
//
// If the working directory cannot be determined (e.g.
// the directory was deleted), FindSessions returns an
// fs.WorkingDirectory error. Parser-level errors from
// scanning session directories propagate unchanged.
package query
