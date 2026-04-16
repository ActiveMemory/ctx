//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve implements the "ctx pad resolve"
// subcommand for inspecting merge conflicts in the
// encrypted scratchpad.
//
// # Behavior
//
// When a git merge conflict occurs on the encrypted
// pad file, git cannot merge the binary ciphertext.
// Instead, ctx stores both sides as separate files
// (ours and theirs). This command decrypts and
// displays both sides so the user can decide which
// entries to keep.
//
// The command requires scratchpad encryption to be
// enabled. If encryption is off, it returns an error
// immediately. It loads the project encryption key
// and attempts to decrypt each conflict file. If
// both files are missing, it reports that no conflict
// exists.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// For each conflict side that exists (ours, theirs),
// prints a labeled header followed by the decrypted
// entries in display format. Missing sides are
// silently skipped. If neither side exists, returns
// an error indicating no conflict files were found.
//
// # Delegation
//
// Decryption is performed by [padCrypto.DecryptFile].
// Display formatting uses [coreResolve.DisplayAll].
// Key loading goes through [crypto.LoadKey]. Output
// is routed through [writePad.ResolveSide].
package resolve
