//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package file centralizes file extension, filename,
// gitignore, profile, and limit constants used across
// the ctx CLI for file I/O operations.
//
// Any code that creates, reads, or validates files
// imports these constants instead of hard-coding
// strings. This makes renames, extension changes, and
// new file types single-point edits.
//
// # File Extensions
//
// Extension constants for every file format ctx
// handles:
//
//   - ExtMarkdown -- context files, specs
//   - ExtTxt (".txt") -- plain text output
//   - ExtGo (".go") -- Go source files
//   - ExtJSONL (".jsonl") -- event logs, hub entries
//   - ExtYAML (".yaml") -- steering files
//   - ExtSh (".sh") -- Unix hook scripts
//   - ExtPs1 (".ps1") -- Windows hook scripts
//   - ExtTmp (".tmp") -- atomic-write temporaries
//   - ExtExample, ExtSample -- safe suffixes that
//     exempt files from secret detection
//
// BackupFormat is the format string for timestamped
// backup filenames (original.1234567890.bak).
//
// # Common Filenames
//
//   - Readme ("README.md") -- standard readme
//   - Index ("index.md") -- generated site index
//   - SchemaDrift -- schema drift report filename
//   - Violations -- governance violations JSON file
//
// # Gitignore Management
//
// FileGitignore is the .gitignore filename.
// GitignoreHeader is the section comment ctx prepends.
// The [Gitignore] variable lists all paths that ctx
// init adds to .gitignore (journal/, logs/, state/,
// keys, and local settings).
//
// # Runtime Configuration Profiles
//
// Profile constants for .ctxrc management:
//
//   - CtxRC (".ctxrc") -- the active config file
//   - CtxRCBase, CtxRCDev -- profile-specific files
//   - ProfileDev, ProfileBase, ProfileProd -- profile
//     identifiers (prod is an alias for base)
//
// # Filename Limits
//
// MaxNameLen (50) caps sanitized filename components
// to prevent filesystem issues with long names.
//
// # Why Centralized
//
// File names and extensions are referenced by init,
// backup, drift detection, gitignore management,
// profile switching, and hook scripts. A single source
// of truth prevents typos and makes auditing file
// access straightforward.
package file
