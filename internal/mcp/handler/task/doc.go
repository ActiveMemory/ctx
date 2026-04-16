//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task iterates pending tasks from TASKS.md
// for MCP tool responses.
//
// This package parses TASKS.md lines and extracts
// top-level pending items, skipping completed sections
// and subtasks. It also provides word-overlap detection
// to check whether a recent action matches an existing
// task.
//
// # Iteration
//
// ForEachPending walks lines from TASKS.md, skipping
// the "## Completed" section and any subtask lines.
// Each top-level pending task is delivered to a visitor
// function. The visitor can return true to stop early.
//
//	task.ForEachPending(lines, func(p Pending) bool {
//	    fmt.Println(p.Index, p.Content)
//	    return false // continue
//	})
//
// # Overlap Detection
//
// ContainsOverlap uses word-set intersection to check
// if a recent action description shares meaningful
// words with a task. It requires at least two
// significant words (length >= MinWordLen) to overlap
// before returning true.
//
//	matched := task.ContainsOverlap(
//	    "added auth validation",
//	    "implement auth validation logic",
//	)
//
// # Types
//
// Pending holds the one-based index and content text
// of a pending top-level task discovered during
// iteration.
package task
