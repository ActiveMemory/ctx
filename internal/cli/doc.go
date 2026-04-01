//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cli contains the implementation of all ctx subcommands.
//
// Each command lives in its own package following the taxonomy:
// parent.go (Cmd wiring), cmd/root/ or cmd/<sub>/ (implementation),
// core/ (shared helpers). The bootstrap package registers all
// commands into the root cobra.Command tree.
package cli
