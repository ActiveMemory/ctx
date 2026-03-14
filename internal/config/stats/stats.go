//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

// Stats display configuration.
const (
	// FilePrefix is the filename prefix for per-session stats JSONL files.
	FilePrefix = "stats-"
	// ReadBufSize is the byte buffer size for reading new lines
	// from stats files during follow/stream mode.
	ReadBufSize = 8192
	// HeaderTime is the column header label for timestamp.
	HeaderTime = "TIME"
	// HeaderSession is the column header label for session ID.
	HeaderSession = "SESSION"
	// HeaderPrompt is the column header label for prompt count.
	HeaderPrompt = "PROMPT"
	// HeaderTokens is the column header label for token count.
	HeaderTokens = "TOKENS"
	// HeaderPct is the column header label for percentage.
	HeaderPct = "PCT"
	// HeaderEvent is the column header label for the event type.
	HeaderEvent = "EVENT"
	// SepTime is the column separator for the time field.
	SepTime = "-------------------"
	// SepSession is the column separator for the session field.
	SepSession = "--------"
	// SepPrompt is the column separator for the prompt field.
	SepPrompt = "------"
	// SepTokens is the column separator for the tokens field.
	SepTokens = "--------"
	// SepPct is the column separator for the percentage field.
	SepPct = "----"
	// SepEvent is the column separator for the event field.
	SepEvent = "------------"
)
