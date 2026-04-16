//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package desc provides command, flag, text, and
// example description lookups backed by embedded YAML.
//
// All user-facing strings are externalized to YAML
// files loaded at init time via lookup.Init. The four
// accessors resolve DescKey constants to their
// localized values. Missing keys return the key
// itself as a fallback, making gaps visible without
// crashing.
//
// # Command Descriptions
//
// Command returns Short and Long descriptions for a
// CLI command by dot-notation key.
//
//	short, long := desc.Command("pad.show")
//
// # Flag Descriptions
//
// Flag returns the description for a CLI flag by
// dot-notation key.
//
//	d := desc.Flag("add.file")
//
// # Example Text
//
// Example returns usage example text for an entry
// type key (decision, learning, task, convention).
//
//	ex := desc.Example("decision")
//
// # General Text
//
// Text returns a user-facing text string by
// dot-notation key. This is used throughout ctx for
// error messages, labels, and formatted output.
//
//	msg := desc.Text("backup.run-hint")
package desc
