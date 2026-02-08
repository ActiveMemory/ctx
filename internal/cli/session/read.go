//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// readContextFile reads a context file and returns its content as a string.
//
// Parameters:
//   - fileName: Name of the file in .context/ (e.g., "DECISIONS.md")
//
// Returns:
//   - string: File content
//   - error: Non-nil if the file cannot be read
func readContextFile(fileName string) (string, error) {
	filePath := filepath.Join(rc.ContextDir(), fileName)
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// readContextSection reads a section from a context file between two headers.
//
// Extracts text between a start header and an optional end header from a
// context file in .context/. Useful for reading specific sections like
// "## In Progress" from TASKS.md.
//
// Parameters:
//   - filename: Name of the file in .context/ (e.g., "TASKS.md")
//   - startHeader: Header marking the section start (e.g., "## In Progress")
//   - endHeader: Header marking the section end, or empty string for
//     the end of the file
//
// Returns:
//   - string: Trimmed content between the headers
//   - error: Non-nil if the file cannot be read or the start header is
//     not found
func readContextSection(
	filename, startHeader, endHeader string,
) (string, error) {
	filePath := filepath.Join(rc.ContextDir(), filename)
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Find start
	startIdx := strings.Index(contentStr, startHeader)
	if startIdx == -1 {
		return "", errSectionNotFound(startHeader)
	}
	startIdx += len(startHeader)

	// Find end
	endIdx := len(contentStr)
	if endHeader != "" {
		idx := strings.Index(contentStr[startIdx:], endHeader)
		if idx != -1 {
			endIdx = startIdx + idx
		}
	}

	section := strings.TrimSpace(contentStr[startIdx:endIdx])
	return section, nil
}

// readRecentDecisions extracts the most recent decisions from DECISIONS.md.
//
// Parses DECISIONS.md to find decision headers (## [YYYY-MM-DD] Title) and
// returns the 3 most recent as a formatted list.
//
// Returns:
//   - string: Formatted list of recent decision titles, or empty if none found
//   - error: Non-nil if DECISIONS.md cannot be read
func readRecentDecisions() (string, error) {
	contentStr, err := readContextFile(config.FileDecision)
	if err != nil {
		return "", err
	}

	// Find decision headers (## [YYYY-MM-DD] Title)
	matches := config.RegExDecision.FindAllStringIndex(contentStr, -1)
	if len(matches) == 0 {
		return "", nil
	}

	// Get the last 3 decisions (most recent)
	limit := config.MaxDecisionsToSummarize
	if len(matches) < limit {
		limit = len(matches)
	}

	var decisions []string
	for i := len(matches) - limit; i < len(matches); i++ {
		start := matches[i][0]
		end := len(contentStr)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		decision := strings.TrimSpace(contentStr[start:end])
		// Only include the header for brevity
		headerEnd := strings.Index(decision, config.NewlineLF)
		if headerEnd != -1 {
			decisions = append(decisions, "- "+decision[:headerEnd])
		}
	}

	return strings.Join(decisions, config.NewlineLF), nil
}

// readRecentLearnings extracts the most recent learnings from LEARNINGS.md.
//
// Parses LEARNINGS.md to find learning entries (- **[YYYY-MM-DD]** text) and
// returns the 5 most recent.
//
// Returns:
//   - string: Formatted list of recent learnings, or empty if none found
//   - error: Non-nil if LEARNINGS.md cannot be read
func readRecentLearnings() (string, error) {
	contentStr, err := readContextFile(config.FileLearning)
	if err != nil {
		return "", err
	}

	// Find learning entries (- **[YYYY-MM-DD]** text)
	matches := config.RegExLearning.FindAllString(contentStr, -1)
	if len(matches) == 0 {
		return "", nil
	}

	// Get the last 5 learnings (most recent)
	limit := config.MaxLearningsToSummarize
	if len(matches) < limit {
		limit = len(matches)
	}

	return strings.Join(matches[len(matches)-limit:], config.NewlineLF), nil
}
