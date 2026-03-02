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

	"github.com/ActiveMemory/ctx/internal/rc"
)

// maxTailBytes is the maximum number of bytes to read from the end of a
// JSONL file when scanning for the last usage block.
const maxTailBytes = 32768

// contextWindow1M is the context window size for 1M-capable models.
const contextWindow1M = 1_000_000

// sessionTokenInfo holds token usage and model information extracted from a
// session's JSONL file.
type sessionTokenInfo struct {
	Tokens int    // Total input tokens (input + cache_creation + cache_read)
	Model  string // Model ID from the last assistant message, or ""
}

// readSessionTokenInfo finds the current session's JSONL file and returns
// the most recent total input token count and model ID from the last
// assistant message. Returns zero value if the file isn't found or has no
// usage data.
//
// Parameters:
//   - sessionID: The Claude Code session ID
//
// Returns:
//   - sessionTokenInfo: Token count and model from the last assistant message
//   - error: Non-nil only on unexpected I/O errors
func readSessionTokenInfo(sessionID string) (sessionTokenInfo, error) {
	if sessionID == "" || sessionID == sessionUnknown {
		return sessionTokenInfo{}, nil
	}

	path, findErr := findJSONLPath(sessionID)
	if findErr != nil || path == "" {
		return sessionTokenInfo{}, findErr
	}

	return parseLastUsageAndModel(path)
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
// needed to extract usage and model data from assistant messages.
type jsonlMessage struct {
	Type    string `json:"type"`
	Message struct {
		Role  string    `json:"role"`
		Model string    `json:"model"`
		Usage usageData `json:"usage"`
	} `json:"message"`
}

// parseLastUsageAndModel reads the tail of a JSONL file and extracts the
// last assistant message's usage data and model ID.
//
// Parameters:
//   - path: Absolute path to the JSONL file
//
// Returns:
//   - sessionTokenInfo: Token count and model, or zero value if not found
//   - error: Non-nil only on I/O errors
func parseLastUsageAndModel(path string) (sessionTokenInfo, error) {
	f, openErr := os.Open(path) //nolint:gosec // path from glob result
	if openErr != nil {
		return sessionTokenInfo{}, openErr
	}
	defer func() { _ = f.Close() }()

	info, statErr := f.Stat()
	if statErr != nil {
		return sessionTokenInfo{}, statErr
	}

	// Read the tail of the file
	size := info.Size()
	offset := int64(0)
	if size > maxTailBytes {
		offset = size - maxTailBytes
	}

	if _, seekErr := f.Seek(offset, io.SeekStart); seekErr != nil {
		return sessionTokenInfo{}, seekErr
	}

	tail, readErr := io.ReadAll(f)
	if readErr != nil {
		return sessionTokenInfo{}, readErr
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
			return sessionTokenInfo{
				Tokens: total,
				Model:  msg.Message.Model,
			}, nil
		}
	}

	return sessionTokenInfo{}, nil
}

// modelContextWindow returns the context window size for a known model ID.
// Returns 0 if the model is not recognized, signaling callers to fall back
// to rc.ContextWindow() or the default.
//
// Claude Code enables the 1M beta header for supported models, so models
// in the 1M-capable set are mapped to 1,000,000 tokens.
//
// Parameters:
//   - model: Model ID string from the JSONL (e.g., "claude-opus-4-6-20260205")
//
// Returns:
//   - int: Context window size in tokens, or 0 if unknown
func modelContextWindow(model string) int {
	if model == "" {
		return 0
	}

	// 1M-capable models. Claude Code enables the beta header for these.
	// Check most specific prefixes first to avoid ambiguity.
	switch {
	case strings.HasPrefix(model, "claude-opus-4-6"):
		return contextWindow1M
	case strings.HasPrefix(model, "claude-sonnet-4-6"):
		return contextWindow1M
	case strings.HasPrefix(model, "claude-sonnet-4-5"):
		return contextWindow1M
	case strings.HasPrefix(model, "claude-sonnet-4-2"):
		// Matches dated snapshots like "claude-sonnet-4-20250514"
		return contextWindow1M
	case model == "claude-sonnet-4" || model == "claude-sonnet-4-0":
		return contextWindow1M
	}

	// All other recognized Claude models default to 200k.
	if strings.HasPrefix(model, "claude-") {
		return rc.DefaultContextWindow
	}

	return 0
}

// effectiveContextWindow returns the context window size using a three-tier
// fallback: JSONL-detected model > .ctxrc setting > 200k default.
//
// Parameters:
//   - model: Model ID string from JSONL (may be empty)
//
// Returns:
//   - int: Effective context window size in tokens
func effectiveContextWindow(model string) int {
	if w := modelContextWindow(model); w > 0 {
		return w
	}
	return rc.ContextWindow()
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
