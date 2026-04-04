//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for trace subcommands.
const (
	// UseTrace is the cobra Use string for the trace command.
	UseTrace = "trace [commit]"
	// UseTraceFile is the cobra Use string for the trace file command.
	UseTraceFile = "file <path[:line-range]>"
	// UseTraceTag is the cobra Use string for the trace tag command.
	UseTraceTag = "tag <commit>"
	// UseTraceCollect is the cobra Use string for the trace collect command.
	UseTraceCollect = "collect"
	// UseTraceHook is the cobra Use string for the trace hook command.
	UseTraceHook = "hook <enable|disable>"
)

// DescKeys for trace subcommands.
const (
	// DescKeyTrace is the description key for the trace command.
	DescKeyTrace = "trace"
	// DescKeyTraceFile is the description key for the trace file command.
	DescKeyTraceFile = "trace.file"
	// DescKeyTraceTag is the description key for the trace tag command.
	DescKeyTraceTag = "trace.tag"
	// DescKeyTraceCollect is the description key for the trace collect command.
	DescKeyTraceCollect = "trace.collect"
	// DescKeyTraceHook is the description key for the trace hook command.
	DescKeyTraceHook = "trace.hook"
)
