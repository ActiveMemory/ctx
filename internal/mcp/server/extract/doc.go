//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package extract converts raw MCP tool arguments
// into typed Go values for use by tool handlers.
//
// # Entry Arguments
//
// EntryArgs extracts the required "type" and "content"
// fields from an MCP tool argument map. It returns an
// error if either field is missing or empty.
//
//	entryType, content, err := extract.EntryArgs(args)
//
// # Options
//
// Opts builds an entity.EntryOpts struct from an MCP
// argument map by extracting optional fields such as
// priority, section, context, rationale, consequence,
// lesson, application, session ID, branch, and commit.
// Missing fields are left as zero values.
//
//	opts := extract.Opts(args)
//	// opts.Priority, opts.Section, etc.
//
// # Design
//
// Both functions perform safe type assertions on the
// interface{} values in the argument map, returning
// zero values for fields that are absent or have
// unexpected types. This avoids panics from malformed
// MCP requests.
package extract
