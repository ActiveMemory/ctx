//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package crypto defines the **typed error constructors**
// returned by [internal/crypto] and its consumers
// ([internal/pad], [internal/notify]). Every encryption,
// decryption, and key-management failure flows through
// one of these constructors.
//
// # Why Typed Errors
//
//   - **Stability** — error categories are part of
//     the public API.
//   - **Routing** — write-side packages map error
//     types to localized text via
//     [internal/assets/read/desc].
//   - **Wrapping** — constructors wrap the
//     underlying cause via `%w` so callers can
//     `errors.Is` against system errors
//     (`io.EOF`, `os.ErrNotExist`) when needed.
//
// # Public Surface
//
// Constructors (one per failure mode):
// [LoadKey], [EncryptFailed], [DecryptFailed],
// [NoKeyAt], [SaveKey], [MkdirKeyDir].
//
// # Why "NoKeyAt" Is Distinct from "LoadKey"
//
// "Key file does not exist yet" is the *normal*
// state on first use; consumers ([pad], [notify])
// treat it as "generate one" rather than "fail".
// Other load failures (permission denied, wrong
// size) are real errors and surface through
// [LoadKey].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package crypto
