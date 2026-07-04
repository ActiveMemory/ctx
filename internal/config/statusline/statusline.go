//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package statusline

// Rendering limits.
const (
	// MaxSegmentLen bounds each rendered segment so a hostile or
	// degenerate payload field cannot flood the status line.
	MaxSegmentLen = 80

	// MaxLineLen bounds the full rendered line. Claude Code displays
	// a single line; anything longer is terminal noise.
	MaxLineLen = 200

	// MaxDirLen bounds the directory portion of the location segment;
	// longer paths collapse to TruncatedDirPrefix + base.
	MaxDirLen = 32

	// MaxPayloadBytes bounds how much stdin is read; the status line
	// payload is a few kilobytes, so anything near this limit is
	// garbage.
	MaxPayloadBytes = 1 << 20
)

// Segment formats and separators.
const (
	// SegmentCapacity is the number of segments a full line renders
	// (location, model, context usage, cost).
	SegmentCapacity = 4

	// PercentMax is the upper bound of a plausible context-usage
	// percentage; payload values outside [0, PercentMax] are dropped.
	PercentMax = 100

	// SegmentSeparator joins the rendered segments.
	SegmentSeparator = " | "

	// ContextFormat renders the context-usage percentage segment.
	ContextFormat = "ctx: %.0f%%"

	// CostFormat renders the session-cost segment as a plain dollar
	// figure. Informational only: there is deliberately no
	// gating or model-switch nudging attached to it.
	CostFormat = "$%.2f"

	// UserHostSeparator joins username and hostname.
	UserHostSeparator = "@"
)

// Path display tokens.
const (
	// HostDomainSeparator splits a fully-qualified hostname; only the
	// short name is rendered at status-line widths.
	HostDomainSeparator = "."

	// RelSelf is filepath.Rel's answer when the path is the home
	// directory itself.
	RelSelf = "."

	// RelParentPrefix marks a filepath.Rel result that escapes the
	// home directory (no home abbreviation applies).
	RelParentPrefix = ".."

	// HomeAbbrev is the home-directory abbreviation.
	HomeAbbrev = "~"

	// HomeAbbrevPrefix prefixes paths under the home directory.
	HomeAbbrevPrefix = "~/"

	// TruncatedDirPrefix replaces overlong directory paths ahead of
	// the final path element.
	TruncatedDirPrefix = ".../"
)

// Bounds of the byte range kept by output sanitization: printable
// ASCII only, per the behavior contract in specs/statusline.md.
// Control bytes, ANSI escapes, and multi-byte runes are stripped so
// payload content cannot corrupt the terminal.
const (
	// ASCIIPrintableMin is the lowest byte kept.
	ASCIIPrintableMin = 0x20

	// ASCIIPrintableMax is the highest byte kept.
	ASCIIPrintableMax = 0x7E
)
