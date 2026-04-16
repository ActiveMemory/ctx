//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stream parses an input stream for XML
// context-update tags and applies them to the project
// context. It powers the "ctx watch" command, which
// monitors agent output for structured context updates.
//
// # Attribute Extraction
//
// [ExtractAttribute] performs simple string-based
// extraction of named attributes from XML tag strings.
// It avoids regex overhead by searching for the
// attribute name followed by an equals sign and double
// quotes.
//
// # Stream Processing
//
// [Process] reads from an io.Reader line by line using
// a buffered scanner with a large buffer for long lines.
// Each line is matched against the context-update regex.
// When a match is found, the function extracts the
// update type, content, and optional attributes
// (section, context, lesson, application, rationale,
// consequence) from the XML tag.
//
// In dry-run mode, Process prints what would happen
// without writing any files. In normal mode, it
// delegates to the apply package to write the update
// to the appropriate context file.
//
// # Supported Update Types
//
// The type attribute determines which context file
// receives the update (e.g., task, decision, learning,
// convention). The apply package handles routing.
package stream
