//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package messages provides metadata for hook message templates.
//
// The embedded registry.yaml maps each hook+variant pair to a
// category and description. [Registry] returns all entries,
// [Lookup] finds a specific one, and [Variants] enumerates
// the available names. Category constants live in
// config/hook (CategoryCustomizable, CategoryCtxSpecific).
package messages
