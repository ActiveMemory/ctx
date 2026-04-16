//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package steering defines inclusion modes, directory
// paths, and file extension constants for the steering
// file subsystem.
//
// Steering files are Markdown documents injected into
// AI prompts to shape agent behavior. They live under
// .context/steering/ and are synced to each tool's
// native rules directory (Cursor .mdc files, Cline
// .clinerules, Kiro steering/). This package provides
// the vocabulary for that sync.
//
// # Inclusion Modes
//
// The [InclusionMode] type controls when a steering
// file is injected:
//
//   - [InclusionAlways]: included in every packet.
//   - [InclusionAuto]: included when the prompt
//     matches the file's globs or keywords.
//   - [InclusionManual]: included only when named
//     explicitly.
//
// # Tool-Native Paths
//
// Directory and extension constants map each editor
// to its native format:
//
//   - [DirCursorDot] + [DirRules] + [ExtMDC]:
//     Cursor rules (.cursor/rules/*.mdc).
//   - [DirClinerules]: Cline rules directory.
//   - [DirKiroDot] + [DirSteering]: Kiro steering.
//
// # Foundation Files
//
// [NameProduct], [NameTech], [NameStructure], and
// [NameWorkflow] are the starter steering files
// scaffolded by ctx init and ctx steering init.
//
// # Other Constants
//
//   - [DefaultPriority] (50): injection priority
//     when omitted from frontmatter.
//   - [LabelAllTools]: display label when a file
//     applies to all tools.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package steering
