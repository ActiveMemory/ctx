//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// runSessionLoad loads and displays a saved session file.
//
// Finds a session file matching the query (by filename, partial match, or index)
// and displays its contents. The query can be a full filename, a substring match,
// or a numeric index from the session list (1 = most recent).
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments where args[0] is the search query
//
// Returns:
//   - error: Non-nil if the sessions directory doesn't exist,
//     the file is not found, or read fails
func runSessionLoad(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Check if the sessions directory exists
	if _, statErr := os.Stat(sessionsDirPath()); os.IsNotExist(statErr) {
		return errNoSessionsDir()
	}

	// Find the matching session file
	filePath, findErr := findSessionFile(query)
	if findErr != nil {
		return findErr
	}

	// Read and display
	content, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return errReadSession(readErr)
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	cmd.Println(fmt.Sprintf("%s Loading: %s", cyan("●"), filepath.Base(filePath)))
	cmd.Println()
	cmd.Println(string(content))

	return nil
}

// runSessionParse parses a JSONL transcript file and outputs formatted content.
//
// Converts a Claude Code JSONL transcript to readable Markdown. Can optionally
// extract potential decisions and learnings from the conversation using pattern
// matching. Output goes to stdout or a specified file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments where args[0] is the input JSONL file path
//   - output: Output file path (empty string for stdout)
//   - extract: If true, extract decisions/learnings instead of full transcript
//
// Returns:
//   - error: Non-nil if the file is not found, parse fails, or write fails
func runSessionParse(
	cmd *cobra.Command, args []string, output string, extract bool,
) error {
	inputPath := args[0]

	// Check if the file exists
	if _, statErr := os.Stat(inputPath); os.IsNotExist(statErr) {
		return errFileNotFound(inputPath)
	}

	green := color.New(color.FgGreen).SprintFunc()

	if extract {
		// Extract decisions and learnings
		decisions, learnings, extractErr := extractInsights(inputPath)
		if extractErr != nil {
			return errExtractInsights(extractErr)
		}

		// Display extracted insights
		cmd.Println(config.SessionHeadingExtractedInsights)
		cmd.Println()
		cmd.Println(fmt.Sprintf(config.MetadataSource+" %s", filepath.Base(inputPath)))
		cmd.Println()

		cmd.Println(config.SessionHeadingPotentialDecisions)
		cmd.Println()
		if len(decisions) == 0 {
			cmd.Println("No decisions detected.")
			cmd.Println()
		} else {
			for _, d := range decisions {
				cmd.Println(fmt.Sprintf("- %s", d))
			}
			cmd.Println()
		}

		cmd.Println(config.SessionHeadingPotentialLearnings)
		cmd.Println()
		if len(learnings) == 0 {
			cmd.Println("No learnings detected.")
			cmd.Println()
		} else {
			for _, l := range learnings {
				cmd.Println(fmt.Sprintf("- %s", l))
			}
			cmd.Println()
		}

		cmd.Println()
		cmd.Println(fmt.Sprintf(
			config.TplSessionInsightsSummary,
			len(decisions), len(learnings),
		))
		return nil
	}

	// Parse the jsonl file
	content, parseErr := parseJsonlTranscript(inputPath)
	if parseErr != nil {
		return errParseTranscript(parseErr)
	}

	// Output
	if output != "" {
		if writeErr := os.WriteFile(
			output, []byte(content), config.PermFile,
		); writeErr != nil {
			return errWriteOutput(writeErr)
		}
		cmd.Println(fmt.Sprintf("%s Parsed transcript saved to %s", green("✓"), output))
	} else {
		cmd.Println(content)
	}

	return nil
}

// runSessionSave saves the current context state to a session file.
//
// Creates a Markdown file in .context/sessions/ containing the current state
// of tasks, decisions, and learnings. The filename includes a timestamp and
// sanitized topic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments where args[0] is the optional topic
//   - sessionType: Type of session (feature, bugfix, refactor, session)
//
// Returns:
//   - error: Non-nil if directory creation, content building, or file write fails
func runSessionSave(
	cmd *cobra.Command, args []string, sessionType string,
) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Get topic from args or use default
	topic := config.DefaultSessionTopic
	if len(args) > 0 {
		topic = args[0]
	}

	// Sanitize the topic for filename
	topic = validation.SanitizeFilename(topic)

	// Ensure sessions directory exists
	if mkdirErr := os.MkdirAll(
		sessionsDirPath(), config.PermExec,
	); mkdirErr != nil {
		return errCreateSessionsDir(mkdirErr)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf(
		config.TplSessionFilename, now.Format("2006-01-02-150405"), topic,
	)
	filePath := filepath.Join(sessionsDirPath(), filename)

	// Build session content
	content, buildErr := buildSessionContent(topic, sessionType, now)
	if buildErr != nil {
		return errBuildContent(buildErr)
	}

	// Write the file
	if writeErr := os.WriteFile(
		filePath, []byte(content), config.PermFile,
	); writeErr != nil {
		return errWriteSession(writeErr)
	}

	cmd.Println(fmt.Sprintf("%s Session saved to %s", green("✓"), filePath))
	return nil
}

// runSessionList lists saved sessions with summaries.
//
// Reads all session files from .context/sessions/, parses their metadata,
// and displays them sorted by date (the newest first). Output includes
// the topic, date, type, summary, and filename for each session.
//
// Parameters:
//   - cmd: Cobra command for output
//   - limit: Maximum number of sessions to display (0 for unlimited)
//
// Returns:
//   - error: Non-nil if reading sessions directory fails
func runSessionList(cmd *cobra.Command, limit int) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if the `sessions` directory exists
	if _, statErr := os.Stat(sessionsDirPath()); os.IsNotExist(statErr) {
		cmd.Println("No sessions found. Use 'ctx session save' to create one.")
		return nil
	}

	// Read directory
	entries, readErr := os.ReadDir(sessionsDirPath())
	if readErr != nil {
		return errReadSessionsDir(readErr)
	}

	// Filter and collect session files
	var sessions []sessionInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Only show .md files (not .jsonl transcripts)
		if !strings.HasSuffix(name, config.ExtMarkdown) {
			continue
		}
		// Skip summary files that accompany jsonl files
		if strings.HasSuffix(name, config.SuffixSummary) {
			continue
		}

		info, parseErr := parseSessionFile(
			filepath.Join(sessionsDirPath(), name),
		)
		if parseErr != nil {
			// Skip files that can't be parsed
			continue
		}
		info.Filename = name
		sessions = append(sessions, info)
	}

	if len(sessions) == 0 {
		cmd.Println("No sessions found. Use 'ctx session save' to create one.")
		return nil
	}

	// Sort by date (newest first) - filenames are date-prefixed
	// so the reverse sort works
	for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	}

	// Limit output
	if limit > 0 && len(sessions) > limit {
		sessions = sessions[:limit]
	}

	// Display
	cmd.Println(fmt.Sprintf("Sessions in %s:", sessionsDirPath()))
	cmd.Println()
	for _, s := range sessions {
		cmd.Println(fmt.Sprintf("%s %s", cyan("●"), s.Topic))
		cmd.Println(fmt.Sprintf("  %s %s | %s %s",
			gray("Date:"), s.Date,
			gray("Type:"), s.Type))
		if s.Summary != "" {
			cmd.Println(fmt.Sprintf("  %s %s",
				gray("Summary:"), truncate(s.Summary, config.MaxPreviewLen)))
		}
		cmd.Println(fmt.Sprintf("  %s %s", yellow("File:"), s.Filename))
		cmd.Println()
	}

	cmd.Println(fmt.Sprintf("Total: %d session(s)", len(sessions)))
	return nil
}
