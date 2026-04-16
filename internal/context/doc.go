//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package context loads and manages .context/ files
// with token counting and validation.
//
// This is a parent package that organizes the context
// subsystem into focused subpackages. Each subpackage
// handles a specific concern in the context lifecycle.
//
// # Subpackages
//
// The context subsystem is composed of:
//
//   - load: reads context files from disk into an
//     entity.Context struct.
//   - resolve: locates the context directory and
//     builds paths to subdirectories like journal/.
//   - sanitize: checks file content for emptiness,
//     distinguishing files with only headers from
//     files with real content.
//   - summary: condenses context content for display.
//   - token: estimates LLM token counts from byte
//     content using a characters-per-token heuristic.
//   - validate: checks whether the .context/ directory
//     exists and contains all required files.
//
// # Lifecycle
//
// A typical context operation follows this flow:
//
//  1. validate.Exists checks the directory is present
//  2. validate.Initialized checks required files exist
//  3. load.Do reads files into an entity.Context
//  4. token.Estimate counts tokens for budget limits
//  5. sanitize.EffectivelyEmpty filters empty files
package context
