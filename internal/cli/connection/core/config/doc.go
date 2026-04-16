//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package config manages the encrypted connection
// configuration for ctx Hub client operations.
//
// # Config Type
//
// [Config] is the persisted hub connection state. It
// holds three fields:
//
//   - HubAddr: the gRPC address (host:port) of the hub.
//   - Token: the client bearer token received during
//     registration.
//   - Types: an optional list of subscribed entry types
//     for filtered listening.
//
// Config is serialized as JSON and encrypted at rest.
//
// # Save
//
// [Save] persists a Config to disk. It marshals the
// struct to JSON, encrypts the bytes using AES-GCM via
// crypto.Encrypt, and writes the ciphertext to
// .context/.connect.enc with restricted permissions
// (fs.PermSecret). The encryption key is loaded from
// the global key path via crypto.GlobalKeyPath.
//
// # Load
//
// [Load] reads and decrypts the stored configuration.
// It reads the ciphertext from .connect.enc, loads the
// encryption key, decrypts via crypto.Decrypt, and
// unmarshals the JSON into a Config struct. Returns an
// error when the file is missing, the key is
// unreadable, or decryption fails.
//
// # Key Management
//
// The unexported loadKey helper reads the encryption
// key from crypto.GlobalKeyPath(). The unexported
// filePath helper resolves the absolute path to
// .connect.enc within the context directory.
//
// # Data Flow
//
// The register subpackage calls Save after receiving a
// client token from the hub. The listen, publish, and
// status subpackages call Load to retrieve credentials
// before making gRPC calls.
package config
