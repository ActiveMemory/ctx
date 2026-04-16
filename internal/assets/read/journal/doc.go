//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal provides access to journal site
// assets embedded in the binary.
//
// # Extra CSS
//
// ExtraCSS returns the CSS content injected into the
// generated journal site for styling beyond what the
// zensical static site generator provides by default.
// This CSS customizes the appearance of session
// entries, timestamps, and other journal-specific
// elements.
//
//	css, err := journal.ExtraCSS()
//
// # Journal Site Generation
//
// The journal site is built by the zensical package
// from Markdown session entries. The extra CSS is
// written alongside the generated site to override
// default theme styles.
package journal
