//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package knowledge

// signalKind distinguishes the two knowledge-health signals (M5): a
// foldable root (staging has accreted; the remedy is /ctx-digest) versus
// a heavy page (bytes over the ceiling; the remedy is split or
// extract-to-tooling). They answer different questions and carry
// different remedies, so the formatter routes on this.
type signalKind int

const (
	// foldable is a staging-count finding on a canonical root.
	foldable signalKind = iota
	// heavy is a byte-weight finding on a page (a root or a theme file).
	heavy
)

// finding describes one knowledge-health signal on one page.
type finding struct {
	// Kind selects the signal (foldable | heavy) and thus the remedy.
	Kind signalKind
	// File is the display name of the page: a canonical basename
	// (e.g. DECISIONS.md) for a root, or the context-relative theme-file
	// path (e.g. conventions/code-style.md) for a heavy theme page.
	File string
	// Count is the measured value — staged entries for foldable, bytes
	// for heavy.
	Count int
	// Threshold is the configured limit the measure exceeded.
	Threshold int
	// Unit is the measurement unit ("entries" | "sections" | "bytes").
	Unit string
}

// Thresholds carries the four M5 knowledge-health limits. A zero in any
// field disables that check (the existing convention).
//
// Fields:
//   - Learnings, Decisions, Conventions: staging-count ceilings per kind
//   - PageBytes: the heavy-page byte ceiling (root and theme files)
type Thresholds struct {
	Learnings   int
	Decisions   int
	Conventions int
	PageBytes   int
}
