//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package theme declares progressive-disclosure themes from the
// add-path. It is what `ctx <kind> add --section Themes` runs
// instead of the ordinary entry writer.
//
// A theme is the unit of bounded recall: the root keeps a
// one-line gist and a link, and the bodies live in a tier-1 file
// the link points at (see specs/progressive-disclosure.md). The
// digesting pass creates themes as a side effect of folding
// entries into them; this package is the other way in, for
// naming a theme up front.
//
// # The two halves
//
// Declaring a theme writes two things, and both are required:
//
//  1. the gist bullet in the root's `## Themes` region, and
//  2. the theme file the bullet links to.
//
// A bullet without its file is a dangling link, which is exactly
// what the root's gist-to-file pairing invariant refuses — a root
// left in that state makes `ctx disclosure` refuse to operate on
// the file at all. So the file is created before the rewritten
// root is handed back to be written.
//
// # The spec shape
//
// The caller supplies `<name> — <gist>`, split on the same
// em-dash separator the theme parser reads back, so what is
// written round-trips through a parse. The slug is derived from
// the name, which is what keeps the bullet's link and the file's
// path in agreement.
//
// Re-declaring an existing theme revises its gist and leaves the
// theme file's accumulated bodies untouched.
//
// # Concurrency
//
// Stateless; safe for concurrent use. Callers serialize their own
// writes to a given root.
package theme
