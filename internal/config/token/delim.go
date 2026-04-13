//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

// EmDash is the em dash character used as a metadata separator.
const EmDash = "—"

// enDash is the en dash character used in topic separators.
const enDash = "-"

// MiddleDot is the middle dot character used as a metadata item separator.
const MiddleDot = "·"

// MetaSeparator is the em dash separator with surrounding spaces.
const MetaSeparator = " " + EmDash + " "

// MetaJoin is the middle dot separator with surrounding spaces.
const MetaJoin = " " + MiddleDot + " "

// Constants for delimiter and separator tokens.
const (
	// Colon is the colon character used as a key-value separator.
	Colon = ":"
	// Comma is the comma character.
	Comma = ","
	// ColonSpace is a colon-space separator for key-value display.
	ColonSpace = ": "
	// CommaSpace is a comma-space separator for joining lists.
	CommaSpace = ", "
	// Dash is a hyphen used as a timestamp segment separator.
	Dash = "-"
	// KeyValueSep is the equals sign used as a key-value separator in state files.
	KeyValueSep = "="
	// Separator is a Markdown horizontal rule used between sections.
	Separator = "---"
	// Ellipsis is a Markdown ellipsis.
	Ellipsis = "..."
	// DoubleQuote is the ASCII double-quote character.
	DoubleQuote = `"`
	// EscapedDoubleQuote is a backslash-escaped double quote for TOML/JSON.
	EscapedDoubleQuote = `\"`
	// PeriodSpace is a period-space separator for joining sentences.
	PeriodSpace = ". "
	// Quotes is the set of quote characters to trim from TOML/JSON values.
	Quotes = `"'`
	// SemicolonSpace is a semicolon-space separator for joining clauses.
	SemicolonSpace = "; "
	// Underscore is the underscore character used as a word separator.
	Underscore = "_"
	// Slash is the forward slash character.
	Slash = "/"
	// Dot is the period character.
	Dot = "."
	// Plus is the plus sign character.
	Plus = "+"
	// Hash is the hash/pound character.
	Hash = "#"
	// ParentDir is the relative parent directory component.
	ParentDir = ".."
	// FrontmatterDelimiter is the YAML frontmatter
	// boundary marker.
	FrontmatterDelimiter = "---"
)

// TopicSeparators are the delimiters between a date and topic in session
// headers (e.g., "2026-01-15 - Fix API").
var TopicSeparators = []string{MetaSeparator, " " + enDash + " "}
