//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// extractAttribute extracts a named attribute from an XML tag string.
//
// Parameters:
//   - tag: XML tag string to search (e.g., `<context-update type="task">`)
//   - attrName: Attribute name to extract (e.g., "type")
//
// Returns:
//   - string: Attribute value, or empty string if not found
func extractAttribute(tag, attrName string) string {
	pattern := regexp.MustCompile(attrName + `="([^"]*)"`)
	match := pattern.FindStringSubmatch(tag)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

// processStream reads from a stream and applies context updates.
//
// Scans input line-by-line looking for <context-update> XML tags.
// When found, parses the type and content, then either displays
// what would happen (--dry-run) or applies the update. Triggers
// auto-save after every WatchAutoSaveInterval updates when enabled.
//
// Parameters:
//   - cmd: Cobra command for output
//   - reader: Input stream to scan (stdin or log file)
//
// Returns:
//   - error: Non-nil if a read error occurs
func processStream(cmd *cobra.Command, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	// Use a larger buffer for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	// Pattern to match context-update tags
	// Captures: 1=full opening tag with attributes, 2=content between tags
	updatePattern := regexp.MustCompile(
		`(<context-update\s+[^>]+)>([^<]+)</context-update>`,
	)

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Track applied updates for auto-save
	updateCount := 0
	var appliedUpdates []ContextUpdate

	for scanner.Scan() {
		line := scanner.Text()

		// Check for context-update commands
		matches := updatePattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				openingTag := match[1]
				update := ContextUpdate{
					Type:         strings.ToLower(extractAttribute(openingTag, "type")),
					Content:      strings.TrimSpace(match[2]),
					Context:      extractAttribute(openingTag, "context"),
					Lesson:       extractAttribute(openingTag, "lesson"),
					Application:  extractAttribute(openingTag, "application"),
					Rationale:    extractAttribute(openingTag, "rationale"),
					Consequences: extractAttribute(openingTag, "consequences"),
				}

				if watchDryRun {
					cmd.Printf(
						"%s Would apply: [%s] %s\n", yellow("○"),
						update.Type, update.Content,
					)
				} else {
					err := applyUpdate(update)
					if err != nil {
						cmd.Printf(
							"%s Failed to apply [%s]: %v\n", color.RedString("✗"),
							update.Type, err,
						)
					} else {
						cmd.Printf(
							"%s Applied: [%s] %s\n", green("✓"), update.Type, update.Content,
						)
						updateCount++
						appliedUpdates = append(appliedUpdates, update)

						// Auto-save every N updates
						if watchAutoSave && updateCount%config.WatchAutoSaveInterval == 0 {
							if err := watchAutoSaveSession(appliedUpdates); err != nil {
								cmd.Printf("%s Auto-save failed: %v\n", yellow("⚠"), err)
							} else {
								cmd.Printf(
									"%s Auto-saved session after %d updates\n", cyan("📸"),
									updateCount,
								)
							}
						}
					}
				}
			}
		}
	}

	// Final auto-save if there are remaining updates
	if watchAutoSave && len(appliedUpdates) > 0 &&
		updateCount%config.WatchAutoSaveInterval != 0 {
		if err := watchAutoSaveSession(appliedUpdates); err != nil {
			cmd.Printf("%s Final auto-save failed: %v\n", yellow("⚠"), err)
		} else {
			cmd.Printf(
				"%s Final auto-save completed (%d total updates)\n",
				cyan("📸"), updateCount,
			)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}
