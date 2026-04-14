//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cmd holds the **lookup keys** for every CLI command's
// `Use`, short-description, long-description, and example
// strings, plus the cobra `Group` identifiers that organize
// `ctx --help`.
//
// Two kinds of constants live here:
//
//   - **`UseXxx`** — the cobra `Use` field. The literal
//     command word as users type it. `UseBackup = "backup"`,
//     `UsePause = "pause"`, etc.
//   - **`DescKeyXxx`** — the lookup key for that command's
//     short / long / example text. Resolved at run-time via
//     [internal/assets/read/desc.Command](key) and
//     [.Example](key), which read from the embedded YAML
//     under [internal/assets/commands].
//
// The same two-step indirection [internal/config/embed/text]
// uses for general display text applies here: copy lives in
// YAML so it can be edited without a Go toolchain, and every
// reference is a typed Go constant so a typo fails to
// compile.
//
// # File Layout — One Command per File
//
// Each file in this package corresponds to one command in
// the cobra tree (`backup.go`, `bootstrap.go`, `connect.go`,
// `system.go`, …) and owns that command's `Use`, `DescKey`,
// and any subcommand `Use`/`DescKey` constants. Adding a new
// command means: add the file here, add the matching YAML
// entry in [internal/assets/commands], wire the cobra
// command in [internal/cli/<command>], and register it in
// [internal/bootstrap].
//
// # Naming Convention
//
// Constants follow `Use<CommandPath>` and
// `DescKey<CommandPath>`. The dotted form of the path is
// the YAML key (`backup`, `system.bootstrap`,
// `hub.peer.add`). The audit suite enforces this both ways
// so a constant without a YAML entry — or a YAML entry
// without a constant — fails CI.
//
// # Group Constants
//
// `Group<Section>` constants (e.g. `GroupGettingStarted`,
// `GroupContext`, `GroupRuntime`) name the cobra command
// groups that organize `ctx --help` output. The grouping is
// applied at registration time in [internal/bootstrap].
//
// # Related Packages
//
//   - [internal/assets/read/desc]      — `desc.Command(key)`
//     and `desc.Example(key)` resolve at run-time.
//   - [internal/assets/commands]       — the YAML store of
//     short / long / example text.
//   - [internal/bootstrap]             — wires every cobra
//     command and assigns each to a `Group<Section>`.
//   - [internal/config/embed/text],
//     [internal/config/embed/flag]     — sister key
//     registries for general text and flag help.
package cmd
