//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package config implements the ctx config command group
// for managing runtime configuration profiles.
//
// Runtime configuration lives in .ctxrc files at the
// project root. Multiple profiles can coexist as
// .ctxrc.<name> files and be swapped with ctx config
// switch. The active profile controls tool selection,
// hook behavior, notification routing, and other
// runtime knobs without modifying code or context files.
//
// # Subcommands
//
//   - schema: outputs the JSON Schema for .ctxrc,
//     useful for editor validation and documentation
//   - status: prints the active profile name and key
//     configuration values
//   - switch: activates a named profile by symlinking
//     or copying .ctxrc.<name> to .ctxrc
//
// # Subpackages
//
//	cmd/schema: JSON Schema output implementation
//	cmd/status: active config display
//	cmd/switchcmd: profile switching logic
//	core: shared config helpers
package config
