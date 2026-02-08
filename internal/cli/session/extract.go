//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
)

// extractInsights parses a JSONL transcript and extracts potential decisions
// and learnings.
//
// Scans assistant messages for patterns indicating decisions
// (e.g., "decided to", "we'll use", "chose X over Y") and learnings
// (e.g., "learned that", "gotcha", "TIL"). Results are deduplicated.
//
// Parameters:
//   - path: Path to the JSONL transcript file
//
// Returns:
//   - []string: Extracted decision insights
//   - []string: Extracted learning insights
//   - error: Non-nil if the file cannot be opened or read
func extractInsights(path string) ([]string, []string, error) {
	file, openErr := os.Open(filepath.Clean(path))
	if openErr != nil {
		return nil, nil, openErr
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}(file)

	var decisions []string
	var learnings []string

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	seen := make(map[string]bool) // Deduplicate

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry transcriptEntry
		if unmarshalErr := json.Unmarshal(
			[]byte(line), &entry,
		); unmarshalErr != nil {
			continue
		}

		// Only look at assistant messages
		if entry.Type != config.RoleAssistant {
			continue
		}

		// Extract text content
		texts := extractTextContent(entry)

		for _, text := range texts {
			// Check for decisions
			for _, pattern := range config.RegExDecisionPatterns {
				matches := pattern.FindAllStringSubmatch(text, -1)
				for _, match := range matches {
					if len(match) > 1 {
						insight := cleanInsight(match[1])
						if insight != "" && !seen[insight] {
							seen[insight] = true
							decisions = append(decisions, insight)
						}
					}
				}
			}

			// Check for learnings
			for _, pattern := range config.RegExLearningPatterns {
				matches := pattern.FindAllStringSubmatch(text, -1)
				for _, match := range matches {
					if len(match) > 1 {
						insight := cleanInsight(match[len(match)-1])
						if insight != "" && !seen[insight] {
							seen[insight] = true
							learnings = append(learnings, insight)
						}
					}
				}
			}
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, nil, scanErr
	}

	return decisions, learnings, nil
}

// extractTextContent extracts all text content from a transcript entry.
//
// Handles both string content and array content (with text and thinking blocks).
//
// Parameters:
//   - entry: Transcript entry to extract text from
//
// Returns:
//   - []string: All text content found in the entry
func extractTextContent(entry transcriptEntry) []string {
	var texts []string

	switch content := entry.Message.Content.(type) {
	case string:
		texts = append(texts, content)
	case []any:
		for _, block := range content {
			blockMap, isMap := block.(contentBlock)
			if !isMap {
				continue
			}

			text, textOk := blockMap[config.ClaudeFieldText].(string)
			thinking, thinkOk := blockMap[config.ClaudeFieldThinking].(string)

			if textOk {
				texts = append(texts, text)
			}
			if thinkOk {
				texts = append(texts, thinking)
			}
		}
	}

	return texts
}
