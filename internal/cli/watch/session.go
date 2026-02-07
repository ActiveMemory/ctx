//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// watchAutoSaveSession saves a session snapshot during watch mode.
//
// Creates a timestamped markdown file in the sessions directory
// containing all updates applied during the watch session. Called
// periodically when --auto-save is enabled.
//
// Parameters:
//   - updates: Slice of ContextUpdate records applied during watch
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func watchAutoSaveSession(updates []ContextUpdate) error {
	sessionsDir := filepath.Join(rc.ContextDir(), config.DirSessions)
	if err := os.MkdirAll(sessionsDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	now := time.Now()
	filename := fmt.Sprintf("%s-watch.md", now.Format("2006-01-02-150405"))
	filePath := filepath.Join(sessionsDir, filename)

	content := buildWatchSession(now, updates)

	if err := os.WriteFile(filePath, []byte(content), config.PermFile); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// buildWatchSession creates a session snapshot from watch mode updates.
//
// Generates a Markdown document with metadata, updates grouped by type,
// and a snapshot of the current TASKS.md content.
//
// Parameters:
//   - timestamp: Time to record as session timestamp
//   - updates: Slice of ContextUpdate records to include
//
// Returns:
//   - string: Formatted Markdown session content
func buildWatchSession(timestamp time.Time, updates []ContextUpdate) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString("# Watch Mode Session" + nl + nl)
	sb.WriteString(fmt.Sprintf("**Date**: %s"+nl, timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s"+nl, timestamp.Format("15:04:05")))
	sb.WriteString("**Type**: watch-auto-save" + nl + nl)
	sb.WriteString(config.Separator + nl + nl)

	sb.WriteString("## Applied Updates" + nl + nl)

	// Group updates by type
	updatesByType := make(map[string][]string)
	for _, u := range updates {
		updatesByType[u.Type] = append(updatesByType[u.Type], u.Content)
	}

	// Write updates by type
	typeOrder := []string{
		config.EntryTask,
		config.EntryDecision,
		config.EntryLearning,
		config.EntryConvention,
		config.EntryComplete,
	}
	for _, t := range typeOrder {
		contents, ok := updatesByType[t]
		if !ok || len(contents) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("### %ss"+nl+nl, strings.ToUpper(t[:1])+t[1:]))
		for _, c := range contents {
			sb.WriteString(fmt.Sprintf("- %s"+nl, c))
		}
		sb.WriteString(nl)
	}

	// Add the current context snapshot
	sb.WriteString(config.Separator + nl + nl)
	sb.WriteString("## Context Snapshot" + nl + nl)

	tasksPath := filepath.Join(rc.ContextDir(), config.FileTask)
	if tasksContent, err := os.ReadFile(tasksPath); err == nil {
		sb.WriteString("### Current Tasks" + nl + nl)
		sb.WriteString("```markdown" + nl)
		sb.WriteString(string(tasksContent))
		sb.WriteString(nl + "```" + nl + nl)
	}

	return sb.String()
}
