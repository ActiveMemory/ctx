//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap is the **CLI assembly layer** for ctx:
// the place where every cobra command in the binary gets
// constructed, grouped, and wired into the root command
// before `cmd.Execute()` runs.
//
// `cmd/ctx/main.go` is intentionally tiny:
//
//	cmd := bootstrap.Initialize(bootstrap.RootCmd())
//	if err := cmd.Execute(); err != nil { ... }
//
// All the actual command registration happens here so the
// command tree is in one auditable place and so the audit
// suite (`cli_cmd_structure_test`) can verify invariants
// like "every command has a non-empty Use", "every command
// has a Short", and "every group has at least one
// command".
//
// # The Root Command
//
// [RootCmd] returns the bare root cobra command with the
// banner, version flag, the `--tool` global flag, and
// the persistent error formatter. It is intentionally
// devoid of subcommands; [Initialize] adds them.
//
// # Group-Based Registration
//
// [Initialize] does the wiring through small grouped
// helpers ([gettingStarted], [contextCmds], [artifacts],
// [sessions], [runtimeCmds], [integrations],
// [diagnostics], [hiddenCmds]), each of which returns a
// `[]registration` that pairs a constructor with a
// [Group] tag. The result is the cobra command tree the
// user sees in `ctx --help`, organized into the same
// sections documented in `docs/cli/index.md`.
//
// New commands plug in by:
//
//  1. Implementing a `Cmd() *cobra.Command` factory in
//     `internal/cli/<command>`.
//  2. Adding the constructor to the right group helper
//     in [group.go] under the matching `embedCmd.Group*`
//     constant.
//  3. Adding the `Use` and `DescKey` constants to
//     [internal/config/embed/cmd] and the matching YAML
//     entries to [internal/assets/commands].
//
// # Hidden Commands
//
// [hiddenCmds] keeps `ctx site` and `ctx system` out of
// `ctx --help` because they are agent-/automation-facing
// rather than user-facing. They still execute when
// invoked directly. The criterion for "hidden" is "no
// human is expected to type this".
//
// # Version Stamping
//
// The build embeds the version string into the package
// at link time via `-ldflags` (see Makefile `build`
// target); the value is exposed through [Version] and
// surfaced by `ctx --version`.
//
// # Concurrency
//
// Bootstrap runs once at process start. Concurrent
// execution is not a concern; cobra serializes
// subcommand dispatch.
package bootstrap
