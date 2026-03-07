//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
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
	pattern := config.RegExFromAttrName(attrName)
	match := pattern.FindStringSubmatch(tag)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

// ProcessStream reads from a stream and applies context updates.
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
func ProcessStream(cmd *cobra.Command, reader io.Reader, dryRun bool) error {
	scanner := bufio.NewScanner(reader)
	// Use a larger buffer for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	updateCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Check for context-update commands
		matches := config.RegExContextUpdate.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				openingTag := match[1]
				update := ContextUpdate{
					Type:         strings.ToLower(ExtractAttribute(openingTag, "type")),
					Content:      strings.TrimSpace(match[2]),
					Context:      ExtractAttribute(openingTag, "context"),
					Lesson:       ExtractAttribute(openingTag, "lesson"),
					Application:  ExtractAttribute(openingTag, "application"),
					Rationale:    ExtractAttribute(openingTag, "rationale"),
					Consequences: ExtractAttribute(openingTag, "consequences"),
				}

				if dryRun {
					cmd.Println(fmt.Sprintf(
						"%s Would apply: [%s] %s\n", yellow("○"),
						update.Type, update.Content,
					))
				} else {
					err := ApplyUpdate(update)
					if err != nil {
						cmd.Println(fmt.Sprintf(
							"%s Failed to apply [%s]: %v\n", color.RedString("✗"),
							update.Type, err,
						))
					} else {
						cmd.Println(fmt.Sprintf(
							"%s Applied: [%s] %s\n", green("✓"), update.Type, update.Content,
						))
						updateCount++
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}
