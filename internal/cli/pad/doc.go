//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad implements the "ctx pad" command for managing an encrypted
// scratchpad.
//
// The scratchpad stores short, sensitive one-liners that travel with the
// project via git but remain opaque at rest. Entries are encrypted with
// AES-256-GCM using a symmetric key at .context/.scratchpad.key.
//
// A plaintext fallback (.context/scratchpad.md) is available via the
// scratchpad_encrypt config option in .contextrc.
package pad
