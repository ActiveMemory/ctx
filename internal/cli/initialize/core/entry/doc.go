//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry handles template creation and insertion
// point detection during the ctx init pipeline.
//
// # Overview
//
// This package provides two capabilities needed when
// initialising a project's context directory: deploying
// entry template files and finding the right place to
// insert ctx content into existing markdown files.
//
// # Behavior
//
// [CreateTemplates] deploys entry template files (TASKS.md,
// DECISIONS.md, LEARNINGS.md, etc.) into .context/templates/.
// [FindInsertionPoint] parses existing markdown to locate
// the position where ctx content should be inserted, placing
// it after a level-1 heading or at the top of the file.
//
// # Algorithms
//
// FindInsertionPoint works by parsing lines top-down:
//
//  1. Skips leading blank lines.
//  2. If the first non-blank line is a level-1 heading,
//     returns the position after the heading and any
//     trailing blank lines.
//  3. If the first non-blank line is a deeper heading
//     or non-heading text, returns position 0 (insert
//     at the top).
//
// CreateTemplates delegates to the tpl sub-package's
// DeployTemplates function, passing the embedded entry
// template list and reader functions.
package entry
