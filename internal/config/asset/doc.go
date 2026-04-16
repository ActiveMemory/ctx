//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package asset defines path constants for the embedded
// asset filesystem that ships inside the ctx binary.
//
// ctx embeds a tree of template files, schemas, hook
// messages, skill definitions, and integration configs
// using Go's embed.FS. This package maps every directory
// and leaf file in that tree to named constants so that
// runtime code never uses raw string paths.
//
// # Directory Constants
//
// Each embedded subdirectory has a constant:
//
//   - [DirClaude] -- Claude Code configuration templates
//   - [DirCommands] -- CLI command metadata (YAML)
//   - [DirContext] -- context file templates
//   - [DirEntryTemplates] -- entry scaffolding templates
//   - [DirIntegrations] -- integration configs (Copilot,
//     Copilot CLI)
//   - [DirHooksMessages] -- hook message templates
//   - [DirJournal] -- journal site assets (CSS)
//   - [DirPermissions] -- allow/deny lists
//   - [DirSchema] -- JSON schemas
//
// # File Constants
//
// Leaf file names are exported individually:
//
//   - [FileCLAUDEMd] -- the CLAUDE.md template
//   - [FileCommandsYAML] -- command registry
//   - [FileFlagsYAML] -- flag definitions
//   - [FilePluginJSON] -- Claude Code plugin manifest
//   - [FileSKILLMd] -- skill template
//   - [FileCtxrcSchemaJSON] -- ctxrc JSON schema
//
// # Composed Paths
//
// Package-level variables combine directories and files
// into full embedded paths:
//
//   - [PathCLAUDEMd] -- claude/CLAUDE.md
//   - [PathPluginJSON] -- plugin manifest path
//   - [PathCommandsYAML] -- command registry path
//   - [PathMessageRegistry] -- hook message registry
//
// # Why Centralized
//
// Embedded paths are referenced by the bootstrap command,
// the init command, the setup command, and every hook
// that renders message templates. A single source of
// truth prevents typos and makes asset reorganization
// safe.
package asset
