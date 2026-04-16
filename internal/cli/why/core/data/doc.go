//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package data provides document alias mappings and
// display ordering for the why command's interactive
// menu. It is the static data layer consumed by the
// menu and show packages.
//
// # Document Aliases
//
// [DocAliases] maps user-facing document names (like
// "manifesto", "about", "invariants") to the embedded
// asset names used by the philosophy asset reader. The
// show package uses this map to resolve a user-provided
// alias to the correct embedded document.
//
// # Display Ordering
//
// [DocOrder] defines the sequence in which documents
// appear in the interactive menu. Each entry is a
// [DocEntry] pairing an alias with its human-readable
// menu label. The menu package iterates this slice to
// render numbered menu items.
//
// # DocEntry Type
//
// [DocEntry] pairs a document lookup key (Alias) with
// its display label (Label). The alias corresponds to a
// key in DocAliases; the label is the text shown to the
// user in the numbered menu.
package data
