//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared scratchpad operations
// used by all pad subcommands.
//
// The scratchpad ("pad") is an encrypted, per-project
// store for short text notes and binary file blobs.
// This core package is the root of the business logic
// tree; it holds sub-packages that each handle one
// concern of the pad lifecycle.
//
// # Sub-Package Overview
//
// The core package itself serves as a namespace. All
// logic lives in sub-packages:
//
//   - add -- append new text or blob entries with
//     stable ID assignment.
//   - blob -- base64 encode/decode binary content
//     within entry strings.
//   - crypto -- decrypt external scratchpad files.
//   - edit -- replace, append, prepend, or update
//     blob content of existing entries.
//   - export -- plan blob extraction to the
//     filesystem with collision avoidance.
//   - imp -- import entries from line-delimited
//     files, stdin, or directory trees.
//   - load -- orchestrate import from files or
//     directories with user output.
//   - merge -- read and merge external scratchpad
//     files with conflict detection.
//   - parse -- split raw bytes into entries and
//     assign stable IDs.
//   - resolve -- convert raw entries to display
//     form for listing.
//   - store -- read/write the default project
//     scratchpad with encryption.
//   - tag -- scan entries for hash-prefixed tags.
//   - validate -- bounds-check entry indexes.
//
// # Data Flow
//
// The cmd/pad layer receives user input, delegates to
// one of these sub-packages for business logic, and
// passes results to the write/pad package for output.
// The store sub-package owns the encrypted file on
// disk; all other sub-packages operate on in-memory
// entry slices.
package core
