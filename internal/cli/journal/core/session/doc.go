//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package session provides session statistics helpers
// for journal generation.
//
// When the journal site is rendered, the index page
// shows aggregate statistics such as the total number
// of unique sessions across all topics. This package
// provides the counting logic for that purpose.
//
// # Unique Session Counting
//
// [CountUnique] is the sole exported function. It
// accepts a slice of TopicData values, each containing
// a list of journal entries. The function iterates
// every entry across all topics, collecting filenames
// into a set, and returns the set size. Because a
// single session can appear under multiple topics, the
// set-based approach avoids double-counting.
//
// The cmd/journal layer calls CountUnique after
// grouping entries by topic and passes the result to
// the write layer for rendering in the site header.
//
// # Data Flow
//
// TopicData arrives from the grouping and index
// packages. Each entry carries a Filename field that
// uniquely identifies the source session file. The
// count feeds into template data for the journal
// index page.
package session
