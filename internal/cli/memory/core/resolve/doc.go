//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve centralizes path resolution shared by every
// memory-bridge subcommand. Each subcommand needs the declared
// context directory (for the .context/memory/ mirror) and its
// parent (the project root, where MEMORY.md lives).
//
// Before this package existed, every memory Run function repeated
// the rc.RequireContextDir + filepath.Dir + cobra.SilenceUsage
// sequence verbatim. Collapsing those three lines into a single
// ContextAndRoot call makes the Run functions read like the task
// they perform, not like the setup scaffolding every Run shares.
//
// The package does not cover memory.DiscoverPath: each caller
// handles its discovery-failure case differently (some emit
// StatusNotActive output, some a tailored NotFound error), so that
// step stays inline where the differences live.
package resolve
