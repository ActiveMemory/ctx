//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// findInsertionPoint finds where to insert ctx content in an existing file.
//
// Logic:
//   - Find the first Markdown heading in the file
//   - If it's an H1 (# Title), return the position after that line
//   - Otherwise return 0 (insert at top)
//
// This ensures ctx content appears prominently near the top, after the
// document title if one exists, rather than being buried at the end.
//
// Parameters:
//   - content: The existing file content
//
// Returns:
//   - int: Byte position where ctx content should be inserted
func findInsertionPoint(content string) int {
	lines := strings.Split(content, "\n")
	pos := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines at the start
		if trimmed == "" {
			pos += len(line) + 1 // +1 for newline
			continue
		}

		// Check if this is a heading
		if strings.HasPrefix(trimmed, "#") {
			// Count the heading level
			level := 0
			for _, ch := range trimmed {
				if ch == '#' {
					level++
				} else {
					break
				}
			}

			if level == 1 {
				// H1 found - insert after this line
				pos += len(line) + 1
				// Skip any blank lines immediately after the H1
				for j := i + 1; j < len(lines); j++ {
					if strings.TrimSpace(lines[j]) == "" {
						pos += len(lines[j]) + 1
					} else {
						break
					}
				}
				return pos
			}
			// Not H1 - insert at top (pos 0)
			return 0
		}

		// Non-empty, non-heading line found first - insert at top
		return 0
	}

	// Empty file or only whitespace - insert at top
	return 0
}

// updateCtxSection replaces the existing ctx section between markers with
// new content.
//
// Locates the ctx markers in the existing content and replaces that section
// with the corresponding section from the template. Creates a timestamped
// backup before modifying.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - existing: Current file content containing ctx markers
//   - newTemplate: Template content with updated ctx section
//
// Returns:
//   - error: Non-nil if the markers are not found or file operations fail
func updateCtxSection(
	cmd *cobra.Command, existing string, newTemplate []byte,
) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Find the start marker
	startIdx := strings.Index(existing, config.CtxMarkerStart)
	if startIdx == -1 {
		return fmt.Errorf("ctx start marker not found")
	}

	// Find the end marker
	endIdx := strings.Index(existing, config.CtxMarkerEnd)
	if endIdx == -1 {
		// No end marker - append from start marker to end
		endIdx = len(existing)
	} else {
		endIdx += len(config.CtxMarkerEnd)
	}

	// Extract the ctx content from the template (between markers)
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, config.CtxMarkerStart)
	templateEnd := strings.Index(templateStr, config.CtxMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return fmt.Errorf("template missing ctx markers")
	}
	ctxContent := templateStr[templateStart : templateEnd+
		len(config.CtxMarkerEnd)]

	// Build new content: before ctx + new ctx content + after ctx
	newContent := existing[:startIdx] + ctxContent + existing[endIdx:]

	// Back up before updating
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FileClaudeMd, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), config.PermFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	cmd.Printf("  %s %s (backup)\n", green("✓"), backupName)

	if err := os.WriteFile(
		config.FileClaudeMd, []byte(newContent), config.PermFile,
	); err != nil {
		return fmt.Errorf("failed to update %s: %w", config.FileClaudeMd, err)
	}
	cmd.Printf(
		"  %s %s (updated ctx section)\n", green("✓"), config.FileClaudeMd,
	)

	return nil
}
