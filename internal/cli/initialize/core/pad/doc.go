//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad provides scratchpad initialisation logic
// for the ctx init pipeline.
//
// # Overview
//
// The scratchpad is a per-project note area that can
// operate in plaintext or encrypted mode. During init,
// this package sets up the appropriate backing storage
// based on the user's runtime configuration.
//
// # Behavior
//
// [Setup] provisions the scratchpad backing storage:
// in encrypted mode it generates a 256-bit AES key,
// in plaintext mode it creates an empty scratchpad.md.
//
// # Modes
//
// Setup checks the rc.ScratchpadEncrypt setting and
// branches accordingly:
//
// Encrypted mode (default):
//
//  1. If the encryption key already exists at the
//     configured key path, skips with an info message.
//  2. If an .enc file exists but no key is found,
//     warns the user about the orphaned ciphertext.
//  3. Otherwise, creates the key directory with
//     restricted permissions, generates a 256-bit
//     AES key via the crypto package, and saves it.
//
// Plaintext mode:
//
//  1. If scratchpad.md already exists, skips with
//     an info message.
//  2. Otherwise, creates an empty scratchpad.md in
//     the .context/ directory.
package pad
