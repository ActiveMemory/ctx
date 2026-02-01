//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// preCompactAutoSave saves a session snapshot before compacting.
//
// Creates a timestamped Markdown file in .context/sessions/ containing
// the current state of TASKS.md for reference after compact operations.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func preCompactAutoSave(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Ensure sessions directory exists
	sessionsDir := filepath.Join(rc.GetContextDir(), "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf("%s-pre-compact.md", now.Format("2006-01-02-150405"))
	filePath := filepath.Join(sessionsDir, filename)

	// Build minimal session content
	content := buildPreCompactSession(now)

	// Write the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	cmd.Printf(
		"%s Auto-saved pre-compact snapshot to %s\n\n", green("✓"), filePath,
	)
	return nil
}

// buildPreCompactSession creates a minimal session snapshot before compact.
//
// The output includes a header with timestamp, purpose description, and
// the full content of TASKS.md wrapped in a Markdown code block.
//
// Parameters:
//   - timestamp: Time to include in the session header
//
// Returns:
//   - string: Formatted Markdown content for the session file
func buildPreCompactSession(timestamp time.Time) string {
	var sb strings.Builder
	nl := config.NewlineLF
	sep := config.Separator

	sb.WriteString("# Pre-Compact Snapshot" + nl + nl)
	sb.WriteString(fmt.Sprintf("**Date**: %s"+nl, timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s"+nl, timestamp.Format("15:04:05")))
	sb.WriteString("**Type**: pre-compact" + nl + nl)
	sb.WriteString(sep + nl + nl)

	sb.WriteString("## Purpose" + nl + nl)
	sb.WriteString(
		"This snapshot was automatically created before running `ctx compact`." + nl,
	)
	sb.WriteString(
		"It preserves the state of context files before any cleanup operations." + nl + nl,
	)
	sb.WriteString(sep + nl + nl)

	// Read and include current TASKS.md content
	tasksPath := filepath.Join(rc.GetContextDir(), config.FileTask)
	if tasksContent, err := os.ReadFile(tasksPath); err == nil {
		sb.WriteString("## Tasks (Before Compact)" + nl + nl)
		sb.WriteString("```markdown" + nl)
		sb.WriteString(string(tasksContent))
		sb.WriteString(nl + "```" + nl + nl)
	}

	return sb.String()
}
