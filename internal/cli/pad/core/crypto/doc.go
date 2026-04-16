//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package crypto reads and decrypts encrypted
// scratchpad files.
//
// The scratchpad stores entries in AES-256 encrypted
// files. When another command needs to read entries
// from a specific encrypted file (as opposed to the
// default project scratchpad), this package provides
// the decryption helper.
//
// # Decryption
//
// [DecryptFile] is the sole exported function. It
// performs three steps:
//
//  1. Reads the encrypted file from baseDir/filename
//     using io.SafeReadFile.
//  2. Decrypts the raw bytes with the provided AES-256
//     key via crypto.Decrypt. If decryption fails, it
//     returns an errCrypto.DecryptFailed error without
//     leaking the ciphertext.
//  3. Parses the decrypted plaintext into individual
//     entries using parse.Entries, which splits on the
//     scratchpad entry separator.
//
// The caller supplies the encryption key; this package
// does not handle key loading. The merge package uses
// DecryptFile when merging entries from an external
// scratchpad file into the current project.
//
// # Error Handling
//
// Read errors propagate from io.SafeReadFile. Decryption
// failures return a generic DecryptFailed error to avoid
// exposing internals. Parse errors are not possible
// because Entries always returns a valid slice.
package crypto
