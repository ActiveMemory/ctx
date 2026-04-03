//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stream

import (
	"bufio"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/watch/core"
	"github.com/ActiveMemory/ctx/internal/cli/watch/core/apply"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/watch"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	writeWatch "github.com/ActiveMemory/ctx/internal/write/watch"
)

// ExtractAttribute extracts a named attribute from an XML tag string.
//
// Parameters:
//   - tag: XML tag string to search (e.g., `<context-update type="task">`)
//   - attrName: Attribute name to extract (e.g., "type")
//
// Returns:
//   - string: Attribute value, or empty string if not found
func ExtractAttribute(tag, attrName string) string {
	// Use simple string search — the attribute names are fixed XML
	// attributes, no regex needed.
	prefix := attrName + marker.AttrEquals
	idx := strings.Index(tag, prefix)
	if idx == -1 {
		return ""
	}
	start := idx + len(prefix)
	end := strings.Index(tag[start:], token.DoubleQuote)
	if end == -1 {
		return ""
	}
	return tag[start : start+end]
}

// Process reads from a stream and applies context updates.
//
// Scans input line-by-line looking for <context-update> XML tags.
// When found, parses the type and content, then either displays
// what would happen (dry-run) or applies the update.
//
// Parameters:
//   - cmd: Cobra command for output
//   - reader: Input stream to scan (stdin or log file)
//   - dryRun: If true, show what would happen without applying
//
// Returns:
//   - error: Non-nil if a read error occurs
func Process(cmd *cobra.Command, reader io.Reader, dryRun bool) error {
	scanner := bufio.NewScanner(reader)
	// Use a larger buffer for long lines
	buf := make([]byte, 0, watch.StreamScannerInitCap)
	scanner.Buffer(buf, watch.StreamScannerMaxSize)

	updateCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Check for context-update commands
		matches := regex.SystemContextUpdate.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= watch.ContextUpdateMinGroups {
				openingTag := match[1]
				update := core.ContextUpdate{
					Type:        strings.ToLower(ExtractAttribute(openingTag, cli.AttrType)),
					Content:     strings.TrimSpace(match[2]),
					Context:     ExtractAttribute(openingTag, cli.AttrContext),
					Lesson:      ExtractAttribute(openingTag, cli.AttrLesson),
					Application: ExtractAttribute(openingTag, cli.AttrApplication),
					Rationale:   ExtractAttribute(openingTag, cli.AttrRationale),
					Consequence: ExtractAttribute(openingTag, cli.AttrConsequence),
				}

				if dryRun {
					writeWatch.DryRunPreview(cmd, update.Type, update.Content)
				} else {
					applyErr := apply.Update(update)
					if applyErr != nil {
						writeWatch.ApplyFailed(cmd, update.Type, applyErr)
					} else {
						writeWatch.ApplySuccess(cmd, update.Type, update.Content)
						updateCount++
					}
				}
			}
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return errFs.ReadInputStream(scanErr)
	}

	return nil
}
