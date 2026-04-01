//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ReadTokenInfo finds the current session's JSONL file and returns
// the most recent total input token count and model ID from the last
// assistant message. Returns zero value if the file isn't found or has no
// usage data.
//
// Parameters:
//   - sessionID: The Claude Code session ID
//
// Returns:
//   - SessionTokenInfo: Token count and model from the last assistant message
//   - error: Non-nil only on unexpected I/O errors
func ReadTokenInfo(sessionID string) (entity.TokenInfo, error) {
	if sessionID == "" || sessionID == session.IDUnknown {
		return entity.TokenInfo{}, nil
	}

	path, findErr := FindJSONLPath(sessionID)
	if findErr != nil || path == "" {
		return entity.TokenInfo{}, findErr
	}

	return ParseLastUsageAndModel(path)
}

// FindJSONLPath locates the JSONL file for a session ID.
//
// Uses glob: ~/.claude/projects/*/{sessionID}.jsonl
// Caches the result in StateDir()/jsonl-path-{sessionID} so the glob
// runs once per session.
//
// Parameters:
//   - sessionID: The Claude Code session ID
//
// Returns:
//   - string: Path to the JSONL file, or empty if not found
//   - error: Non-nil only on unexpected errors
func FindJSONLPath(sessionID string) (string, error) {
	// Check cache first
	cacheFile := filepath.Join(state.Dir(), stats.JsonlPathCachePrefix+sessionID)
	if data, readErr := internalIo.SafeReadUserFile(cacheFile); readErr == nil {
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

	pattern := filepath.Join(
		home, dir.Claude, dir.Projects,
		token.GlobStar, sessionID+file.ExtJSONL,
	)
	matches, globErr := filepath.Glob(pattern)
	if globErr != nil {
		return "", globErr
	}

	if len(matches) == 0 {
		return "", nil
	}

	// Cache the result for subsequent calls this session.
	if writeErr := os.WriteFile(
		cacheFile, []byte(matches[0]), fs.PermSecret,
	); writeErr != nil {
		ctxLog.Warn(warn.Write, cacheFile, writeErr)
	}
	return matches[0], nil
}

// ParseLastUsageAndModel reads the tail of a JSONL file and extracts the
// last assistant message's usage data and model ID.
//
// Parameters:
//   - path: Absolute path to the JSONL file
//
// Returns:
//   - SessionTokenInfo: Token count and model, or zero value if not found
//   - error: Non-nil only on I/O errors
func ParseLastUsageAndModel(path string) (entity.TokenInfo, error) {
	f, openErr := internalIo.SafeOpenUserFile(path)
	if openErr != nil {
		return entity.TokenInfo{}, openErr
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn(warn.Close, path, closeErr)
		}
	}()

	info, statErr := f.Stat()
	if statErr != nil {
		return entity.TokenInfo{}, statErr
	}

	// Read the tail of the file
	size := info.Size()
	offset := int64(0)
	if size > claude.MaxTailBytes {
		offset = size - claude.MaxTailBytes
	}

	if _, seekErr := f.Seek(offset, io.SeekStart); seekErr != nil {
		return entity.TokenInfo{}, seekErr
	}

	tail, readErr := io.ReadAll(f)
	if readErr != nil {
		return entity.TokenInfo{}, readErr
	}

	// Scan lines in reverse for the last assistant message with usage
	lines := bytes.Split(tail, []byte(token.NewlineLF))
	for i := len(lines) - 1; i >= 0; i-- {
		line := bytes.TrimSpace(lines[i])
		if len(line) == 0 {
			continue
		}

		// Quick check: skip lines that can't contain usage data
		if !bytes.Contains(line, []byte(claude.FieldUsage)) {
			continue
		}
		if !bytes.Contains(line, []byte(claude.FieldInputTokens)) {
			continue
		}

		var msg jsonlMessage
		if jsonErr := json.Unmarshal(line, &msg); jsonErr != nil {
			continue
		}

		if msg.Message.Role != claude.RoleAssistant {
			continue
		}

		u := msg.Message.Usage
		total := u.InputTokens + u.CacheCreationInputTokens + u.CacheReadInputTokens
		if total > 0 {
			return entity.TokenInfo{
				Tokens: total,
				Model:  msg.Message.Model,
			}, nil
		}
	}

	return entity.TokenInfo{}, nil
}

// ModelContextWindow returns the context window size for a known model ID.
// Returns 0 if the model is not recognized, signaling callers to fall back
// to other detection tiers.
//
// Parameters:
//   - model: Model ID string from the JSONL (e.g., "claude-opus-4-6-20260205")
//
// Returns:
//   - int: Context window size in tokens, or 0 if unknown
func ModelContextWindow(model string) int {
	if model == "" {
		return 0
	}

	if !strings.HasPrefix(model, claude.ModelPrefix) {
		return 0
	}

	lower := strings.ToLower(model)

	// 1M models: explicit [1m] suffix OR Opus 4.6+ (always 1M).
	if strings.Contains(lower, claude.ModelSuffix1M) {
		return claude.ContextWindow1M
	}
	if strings.Contains(lower, claude.ModelOpus) {
		return claude.ContextWindow1M
	}

	return rc.DefaultContextWindow
}

// EffectiveContextWindow returns the context window size using a four-tier
// fallback where ground truth outranks configuration:
//
//  1. JSONL model ID: actual model running the session (ground truth)
//  2. Claude Code ~/.claude/settings.json: configured model selection
//  3. Explicit .ctxrc context_window: manual override / escape hatch
//  4. rc.ContextWindow() default (200k)
//
// Parameters:
//   - model: Model ID string from JSONL (may be empty)
//
// Returns:
//   - int: Effective context window size in tokens
func EffectiveContextWindow(model string) int {
	// Tier 1: model-based detection (ground truth from session JSONL).
	if w := ModelContextWindow(model); w > 0 {
		return w
	}
	// Tier 2: auto-detect from Claude Code settings.
	if ClaudeSettingsHas1M() {
		return claude.ContextWindow1M
	}
	// Tier 3: explicit .ctxrc override (fallback for non-Claude tools).
	if w := rc.RC().ContextWindow; w > 0 && w != rc.DefaultContextWindow {
		return w
	}
	// Tier 4: default.
	return rc.ContextWindow()
}

// ClaudeSettingsHas1M reads ~/.claude/settings.json and returns true if the
// selected model name contains "[1m]", indicating the user has opted into
// the 1M extended context window. Returns false on any error.
//
// Returns:
//   - bool: True if 1M context is enabled
func ClaudeSettingsHas1M() bool {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return false
	}
	data, readErr := internalIo.SafeReadUserFile(
		filepath.Join(home, dir.Claude, claude.GlobalSettings),
	)
	if readErr != nil {
		return false
	}
	var settings struct {
		Model string `json:"model"`
	}
	if jsonErr := json.Unmarshal(data, &settings); jsonErr != nil {
		return false
	}
	return strings.Contains(strings.ToLower(settings.Model), claude.ModelSuffix1M)
}

// FormatTokenCount formats a token count as a human-readable abbreviated
// string: "1.2k", "52k", "164k".
//
// Parameters:
//   - tokens: Token count to format
//
// Returns:
//   - string: Abbreviated token count
func FormatTokenCount(tokens int) string {
	if tokens < cfgFmt.SIThreshold {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatSIInteger), tokens)
	}
	k := float64(tokens) / cfgFmt.SIThreshold
	if k < 10 {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatSIKilo), k)
	}
	return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatSIKiloInt), int(k))
}

// FormatWindowSize formats the context window size as a human-readable
// abbreviated string for display in token usage lines: "200k", "128k".
//
// Parameters:
//   - size: Window size in tokens
//
// Returns:
//   - string: Abbreviated window size
func FormatWindowSize(size int) string {
	if size < cfgFmt.SIThreshold {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatSIInteger), size)
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteFormatSIKiloInt), size/cfgFmt.SIThreshold,
	)
}
