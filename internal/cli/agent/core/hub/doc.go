//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub loads shared knowledge files from
// .context/hub/ for inclusion in agent context packets.
//
// # LoadBodies
//
// [LoadBodies] reads all Markdown files from the hub
// directory and returns their contents as a string
// slice. It silently skips directories, non-Markdown
// files, empty files, and files that fail to read.
//
// The hub directory is resolved as a subdirectory of
// rc.ContextDir() using the path defined in
// cfgHub.DirHub. When the directory does not exist,
// LoadBodies returns nil, making shared knowledge
// entirely opt-in.
//
// # File Reading
//
// Each file is read through io.SafeReadUserFile, which
// enforces size limits and symlink safety. Only files
// with the .md extension (file.ExtMarkdown) are
// considered.
//
// # Data Flow
//
// The budget subpackage calls LoadBodies during context
// assembly to append shared knowledge sections to the
// agent packet. Hub content is additive and does not
// displace project-local context files.
package hub
