//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package why defines constants for the ctx why
// command, which displays embedded philosophy and
// design documents.
//
// ctx ships three embedded documents that explain
// the project's motivation and design principles.
// The ctx why command presents them in a menu for
// interactive reading. This package maps user-facing
// aliases to embedded asset names so the command can
// resolve both CLI arguments and menu selections.
//
// # User-Facing Aliases
//
// These constants are used as CLI arguments and
// interactive menu keys:
//
//   - [DocManifesto]: the project manifesto.
//   - [DocAbout]: the about page.
//   - [DocInvariants]: design invariants.
//
// # Embedded Asset Names
//
// Each alias maps to an asset file stem under
// internal/assets/why/:
//
//   - [DocAliasManifesto]: "manifesto".
//   - [DocAliasAbout]: "about".
//   - [DocAliasInvariants]: "design-invariants".
//
// The alias and asset name differ only for
// invariants, where the user-facing key is shorter
// than the asset filename.
//
// # Why Centralized
//
// The command, the menu renderer, and the asset
// loader all need to agree on these strings. A
// single package prevents the CLI from accepting
// an alias the loader cannot resolve.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package why
