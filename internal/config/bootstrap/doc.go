//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap defines display, version, and parsing
// constants for the ctx system bootstrap command.
//
// When an AI agent starts a session, it runs
// "ctx system bootstrap" to discover the context directory
// and display a summary of available context files. This
// package provides the formatting and parsing rules used
// by that bootstrap output.
//
// # Version Fallback
//
// [DefaultVersion] is the version string ("dev") used when
// the binary is built without ldflags. Production builds
// inject the real version at compile time; this constant
// ensures the bootstrap banner always has a value.
//
// # File List Formatting
//
// The bootstrap output lists context files in a wrapped,
// indented format:
//
//   - [FileListWidth] sets the wrap width at 55 characters
//     to fit comfortably in terminal sidebars.
//   - [FileListIndent] adds a two-space indent prefix to
//     each wrapped line for visual hierarchy.
//
// # Numbered List Parsing
//
// Some context files use numbered lists (e.g. "1. item").
// The bootstrap parser strips these prefixes:
//
//   - [NumberedListSep] is the ". " separator between the
//     number and the text.
//   - [NumberedListMaxDigits] limits prefix detection to
//     two-digit numbers, avoiding false positives on lines
//     that happen to start with digits.
//
// # Why Centralized
//
// The bootstrap command, the agent packet builder, and
// the version banner all reference these values. Keeping
// them together ensures consistent formatting across all
// entry points.
package bootstrap
