//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema implements the "ctx config schema"
// subcommand that prints the embedded JSON Schema for
// the .ctxrc configuration file.
//
// # What It Does
//
// Reads the JSON Schema from the embedded asset
// bundle and writes it to stdout. The schema
// defines all valid .ctxrc fields, their types,
// defaults, and constraints.
//
// # Flags
//
// None. The command accepts no arguments.
//
// # Output
//
// Raw JSON written to stdout. The output can be
// piped to a file or used for editor integration:
//
//	ctx config schema > ctxrc-schema.json
//
// Editors that support JSON Schema (VS Code,
// JetBrains, Neovim with LSP) can use this file
// to provide autocompletion and validation when
// editing .ctxrc.
//
// # Delegation
//
// [Cmd] builds the cobra.Command. The RunE handler
// reads the schema via [schema.Schema] from the
// embedded assets and writes it through
// [writeConfig.Schema].
package schema
