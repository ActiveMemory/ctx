//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dir defines directory path constants used
// throughout the ctx application for locating project
// structure, context subdirectories, and platform-specific
// home paths.
//
// ctx organizes project state under .context/ with many
// subdirectories (archive, journal, hooks, skills, etc.)
// and interacts with .claude/ for Claude Code integration.
// This package names every directory so that path
// construction never relies on raw strings.
//
// # Context Subdirectories
//
// The .context/ directory contains:
//
//   - [Archive]: archived task snapshots
//   - [Hooks] / [HooksMessages]: lifecycle hook scripts
//     and message overrides
//   - [Journal]: session journal entries
//   - [JournalObsidian]: Obsidian-format journal export
//   - [JournalSite]: static site output for journals
//   - [Logs]: log files
//   - [Memory] / [MemoryArchive]: memory bridge files
//   - [Reports]: generated reports
//   - [Sessions]: session summaries
//   - [Skills]: skill definitions
//   - [Steering]: steering files for hook behavior
//   - [State]: project-scoped runtime state
//   - [Trace]: commit context trace data
//   - [Templates]: entry scaffolding templates
//
// # Project Root Directories
//
//   - [Claude] (".claude"): Claude Code configuration
//   - [Context] (".context"): the default context root
//   - [Ideas] ("ideas"): early-stage explorations
//   - [Specs] ("specs"): formalized feature specs
//   - [Projects] ("projects"): .claude/projects/
//
// # User-Level Directories
//
//   - [CtxData] ("~/.ctx/"): user-level ctx data
//
// # Default Paths
//
//   - [DefaultSteeringPath]: .context/steering
//   - [DefaultHooksPath]: .context/hooks
//
// # Platform Paths
//
//   - [HomeLinux] ("home") and [HomeMacOS] ("Users")
//     identify home directory parents for cross-platform
//     path resolution.
//
// # Journal Site Output
//
//   - [JournalDocs], [JournTopics], [JournalFiles],
//     [JournalTypes] name subdirectories in the generated
//     journal static site.
//
// # Why Centralized
//
// Nearly every package in ctx constructs paths to context
// subdirectories. A single source of truth prevents path
// typos and makes directory renames safe.
package dir
