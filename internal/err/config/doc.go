//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package config defines the typed error constructors
// for configuration and settings operations. Every
// failure related to runtime profiles, tool selection,
// settings files, and schema embedding flows through
// one of these constructors.
//
// # Domain
//
// Errors fall into three categories:
//
//   - **Profile / tool validation**: the user
//     specified an unknown profile or unsupported
//     AI tool. Constructors: [UnknownProfile],
//     [InvalidTool], [UnsupportedTool],
//     [UnknownUpdateType].
//   - **Settings files**: the local or golden
//     settings file is missing, or marshaling
//     failed. Constructors: [SettingsNotFound],
//     [GoldenNotFound], [MarshalSettings],
//     [MarshalPlugins].
//   - **Schema / profile IO**: reading a profile
//     file or the embedded JSON schema failed.
//     Constructors: [ReadProfile],
//     [ReadEmbeddedSchema].
//
// # Wrapping Strategy
//
// IO constructors wrap their cause with fmt.Errorf
// %w so callers can errors.Is against system errors.
// Validation constructors return plain errors since
// there is no underlying cause to chain.
//
// All user-facing text is resolved through
// [internal/assets/read/desc] at construction time.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package config
