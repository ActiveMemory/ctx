//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// maxTailBytes is the maximum number of bytes to read from the end of a
// JSONL file when scanning for the last usage block.
const maxTailBytes = 32768

// readSessionTokenUsage finds the current session's JSONL file and returns
// the most recent total input token count (input_tokens + cache_creation +
// cache_read). Returns 0, nil if the file isn't found or has no usage data.
//
// Parameters:
//   - sessionID: The Claude Code session ID
//
// Returns:
//   - int: Total input tokens from the last assistant message, or 0
//   - error: Non-nil only on unexpected I/O errors
func readSessionTokenUsage(sessionID string) (int, error) {
	if sessionID == "" || sessionID == sessionUnknown {
		return 0, nil
	}

	path, findErr := findJSONLPath(sessionID)
	if findErr != nil || path == "" {
		return 0, findErr
	}

	return parseLastUsage(path)
}

// findJSONLPath locates the JSONL file for a session ID.
//
// Uses glob: ~/.claude/projects/*/{sessionID}.jsonl
// Caches the result in secureTempDir()/jsonl-path-{sessionID} so the glob
// runs once per session.
//
// Parameters:
//   - sessionID: The Claude Code session ID
//
// Returns:
//   - string: Path to the JSONL file, or empty if not found
//   - error: Non-nil only on unexpected errors
func findJSONLPath(sessionID string) (string, error) {
	// Check cache first
	cacheFile := filepath.Join(secureTempDir(), "jsonl-path-"+sessionID)
	if data, readErr := os.ReadFile(cacheFile); readErr == nil { //nolint:gosec // temp dir path
		cached := strings.TrimSpace(string(data))
		if cached != "" {
			if _, statErr := os.Stat(cached); statErr == nil {
				return cached, nil
			}
		}
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", nil
	}

	pattern := filepath.Join(home, ".claude", "projects", "*", sessionID+".jsonl")
	matches, globErr := filepath.Glob(pattern)
	if globErr != nil {
		return "", globErr
	}

	if len(matches) == 0 {
		return "", nil
	}

	// Cache the result for subsequent calls this session
	_ = os.WriteFile(cacheFile, []byte(matches[0]), 0o600)
	return matches[0], nil
}

// usageData represents the minimal usage fields from a Claude Code JSONL
// assistant message. Only the fields needed for token counting are included.
type usageData struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// jsonlMessage represents the minimal structure of a Claude Code JSONL line
// needed to extract usage data from assistant messages.
type jsonlMessage struct {
	Type    string `json:"type"`
	Message struct {
		Role  string    `json:"role"`
		Usage usageData `json:"usage"`
	} `json:"message"`
}

// parseLastUsage reads the tail of a JSONL file and extracts the last
// assistant message's usage data. Returns the sum of input_tokens,
// cache_creation_input_tokens, and cache_read_input_tokens.
//
// Parameters:
//   - path: Absolute path to the JSONL file
//
// Returns:
//   - int: Total input tokens, or 0 if no usage data found
//   - error: Non-nil only on I/O errors
func parseLastUsage(path string) (int, error) {
	f, openErr := os.Open(path) //nolint:gosec // path from glob result
	if openErr != nil {
		return 0, openErr
	}
	defer func() { _ = f.Close() }()

	info, statErr := f.Stat()
	if statErr != nil {
		return 0, statErr
	}

	// Read the tail of the file
	size := info.Size()
	offset := int64(0)
	if size > maxTailBytes {
		offset = size - maxTailBytes
	}

	if _, seekErr := f.Seek(offset, io.SeekStart); seekErr != nil {
		return 0, seekErr
	}

	tail, readErr := io.ReadAll(f)
	if readErr != nil {
		return 0, readErr
	}

	// Scan lines in reverse for the last assistant message with usage
	lines := bytes.Split(tail, []byte("\n"))
	for i := len(lines) - 1; i >= 0; i-- {
		line := bytes.TrimSpace(lines[i])
		if len(line) == 0 {
			continue
		}

		// Quick check: skip lines that can't contain usage data
		if !bytes.Contains(line, []byte(`"usage"`)) {
			continue
		}
		if !bytes.Contains(line, []byte(`"input_tokens"`)) {
			continue
		}

		var msg jsonlMessage
		if jsonErr := json.Unmarshal(line, &msg); jsonErr != nil {
			continue
		}

		if msg.Message.Role != "assistant" {
			continue
		}

		u := msg.Message.Usage
		total := u.InputTokens + u.CacheCreationInputTokens + u.CacheReadInputTokens
		if total > 0 {
			return total, nil
		}
	}

	return 0, nil
}

// formatTokenCount formats a token count as a human-readable abbreviated
// string: "1.2k", "52k", "164k".
//
// Parameters:
//   - tokens: Token count to format
//
// Returns:
//   - string: Abbreviated token count
func formatTokenCount(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	k := float64(tokens) / 1000
	if k < 10 {
		return fmt.Sprintf("%.1fk", k)
	}
	return fmt.Sprintf("%dk", int(k))
}

// formatWindowSize formats the context window size as a human-readable
// abbreviated string for display in token usage lines: "200k", "128k".
//
// Parameters:
//   - size: Window size in tokens
//
// Returns:
//   - string: Abbreviated window size
func formatWindowSize(size int) string {
	if size < 1000 {
		return fmt.Sprintf("%d", size)
	}
	return fmt.Sprintf("%dk", size/1000)
}
