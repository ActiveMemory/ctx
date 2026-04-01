//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

// EmDash is the em dash character used as a metadata separator.
const EmDash = "—"

const enDash = "-"

// MiddleDot is the middle dot character used as a metadata item separator.
const MiddleDot = "·"

// MetaSeparator is the em dash separator with surrounding spaces.
const MetaSeparator = " " + EmDash + " "

// MetaJoin is the middle dot separator with surrounding spaces.
const MetaJoin = " " + MiddleDot + " "

const (
	// Colon is the colon character used as a key-value separator.
	Colon = ":"
	// Comma is the comma character.
	Comma = ","
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
	// CloseBracket is the closing square bracket.
	CloseBracket = "]"
	// PeriodSpace is a period-space separator for joining sentences.
	PeriodSpace = ". "
	// Quotes is the set of quote characters to trim from TOML/JSON values.
	Quotes = `"'`
	// SemicolonSpace is a semicolon-space separator for joining clauses.
	SemicolonSpace = "; "
	// Underscore is the underscore character used as a word separator.
	Underscore = "_"
)

// TopicSeparators are the delimiters between a date and topic in session
// headers (e.g., "2026-01-15 - Fix API").
var TopicSeparators = []string{MetaSeparator, " " + enDash + " "}
