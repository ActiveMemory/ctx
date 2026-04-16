//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package crypto provides the **AES-256-GCM** encryption
// primitives ctx uses for the two pieces of state that
// must never land on disk in plaintext: the **scratchpad**
// ([internal/pad]) and the **webhook URL** ([internal/notify]).
//
// The package is deliberately small. Heavy lifting (key
// management policies, rotation cadence, file paths) lives
// in the consumer packages; this package owns only the
// "given a key and bytes, produce ciphertext / produce
// plaintext" primitives plus the on-disk key file format.
//
// # Public Surface
//
//   - **[GenerateKey]** — returns a fresh 32-byte
//     (256-bit) key from `crypto/rand`. The caller
//     persists it via [SaveKey].
//   - **[SaveKey](path, key)** — writes the key to
//     `path` with `0o600` permissions. Refuses to
//     overwrite an existing file (the caller must
//     remove first; `ctx pad rotate` does this
//     intentionally).
//   - **[LoadKey](path)** — reads the key back. Returns
//     a typed error from [internal/err/crypto] when
//     the file is missing, world-readable, or the
//     wrong size.
//   - **[Encrypt](key, plaintext)** — produces nonce
//     prepended to ciphertext: a fresh random
//     12-byte nonce per call concatenated with the
//     AES-256-GCM ciphertext.
//   - **[Decrypt](key, payload)** — splits nonce from
//     ciphertext and decrypts. Returns a typed error
//     on auth-tag mismatch, short payload, or
//     missing key.
//
// # File Format
//
// Both encrypted blobs (`.notify.enc`,
// `.context/.scratchpad.enc`) are the raw output of
// [Encrypt] — no header, no version, no JSON wrapper.
// The format is purely
// `nonce(12) || ciphertext(...) || tag(16)`.
//
// # Per-Machine Key
//
// The key lives at `~/.ctx/.ctx.key` (one key per user,
// shared by every project on that machine). Cross-machine
// scratchpad sync requires copying that key — see
// `docs/recipes/scratchpad-sync.md` for the user-facing
// procedure.
//
// # Concurrency
//
// All functions are stateless. AES-256-GCM is safe for
// concurrent calls — each [Encrypt] generates a fresh
// nonce internally, and the underlying cipher is
// reentrant.
package crypto
