//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package crypto defines constants for the AES-256-GCM
// encryption subsystem used by ctx's scratchpad and
// notification features.
//
// ctx encrypts sensitive data at rest: the scratchpad
// (short notes that travel with the project) and webhook
// URLs for notifications. This package provides the
// cryptographic parameters and file names that the
// encryption layer depends on.
//
// # Cryptographic Parameters
//
//   - [KeySize] is 32 bytes (256 bits), matching the
//     AES-256 key length required by the Go standard
//     library's crypto/aes package.
//   - [NonceSize] is 12 bytes, the standard GCM nonce
//     length. Each encryption operation generates a
//     fresh random nonce prepended to the ciphertext.
//
// # Encrypted File Names
//
//   - [NotifyEnc] (".notify.enc") stores the encrypted
//     webhook URL used by the notification system. The
//     URL is encrypted so it can be committed to version
//     control without exposing the endpoint.
//   - [ContextKey] (".ctx.key") is the encryption key
//     file. It lives in .context/ and is excluded from
//     version control via .gitignore.
//
// # Why Centralized
//
// The pad command, the notify command, and the key
// management utilities all need the same key size,
// nonce size, and file names. Centralizing them here
// ensures the encryption parameters are consistent
// and auditable in one place.
package crypto
