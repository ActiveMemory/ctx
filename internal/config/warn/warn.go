//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package warn provides format string constants for best-effort
// warning messages routed through log.Warn.
//
// These are Printf-style format strings for common I/O failure
// patterns. Using constants prevents typo drift across 40+ call sites.
// Import as config/warn.
package warn

// Format strings for file I/O warnings. Each takes (path, error).
const (
	// Close is the format for file close failures.
	Close = "close %s: %v"

	// Write is the format for file write failures.
	Write = "write %s: %v"

	// Remove is the format for file remove failures.
	Remove = "remove %s: %v"

	// Mkdir is the format for directory creation failures.
	Mkdir = "mkdir %s: %v"

	// Rename is the format for file rename failures.
	Rename = "rename %s: %v"

	// Walk is the format for directory walk failures.
	Walk = "walk %s: %v"

	// Getwd is the format for working directory resolution failures.
	Getwd = "getwd: %v"

	// Marshal is the format for JSON marshal failures. Takes (error).
	Marshal = "marshal: %v"

	// Readdir is the format for directory read failures.
	Readdir = "readdir %s: %v"

	// CloseResponse is the format for HTTP response body close failures.
	CloseResponse = "close response body: %v"

	// ParseConfig is the format for config file parse failures.
	ParseConfig = "warning: failed to parse %s: %v (using defaults)"

	// CopilotClose is the format for Copilot CLI file close failures.
	CopilotClose = "copilot-cli: close %s: %v"

	// JSONEncode is the JSON-safe error for encoding failures.
	JSONEncode = `{"error": "json encode: %v"}`
)

// Warn context identifiers for index generation.
const (
	// IndexHeader is the context label for index header write errors.
	IndexHeader = "index-header"
	// IndexSeparator is the context label for index separator write
	// errors.
	IndexSeparator = "index-separator"
	// IndexRow is the context label for index row write errors.
	IndexRow = "index-row"
	// ResponseBody is the context label for HTTP response body
	// close errors.
	ResponseBody = "response body"
)
