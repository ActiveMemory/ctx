//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// oversizeNudgeContent checks for an injection-oversize flag file and returns
// the raw nudge content if present. Deletes the flag after reading (one-shot).
//
// Returns:
//   - string: raw oversize nudge content, or empty string if no flag
func oversizeNudgeContent() string {
	baseDir := filepath.Join(rc.ContextDir(), dir.State)
	flagPath := filepath.Join(baseDir, stats.ContextSizeInjectionOversizeFlag)
	data, readErr := io.SafeReadFile(
		baseDir, stats.ContextSizeInjectionOversizeFlag,
	)
	if readErr != nil {
		return ""
	}

	tokenCount := extractOversizeTokens(data)
	fallback := fmt.Sprintf(
		desc.Text(text.DescKeyCheckContextSizeOversizeFallback), tokenCount,
	)
	content := message.LoadMessage(hook.CheckContextSize, hook.VariantOversize,
		map[string]any{stats.VarTokenCount: tokenCount}, fallback)
	if content == "" {
		_ = os.Remove(flagPath) // silenced, still consume the flag
		return ""
	}

	_ = os.Remove(flagPath) // one-shot: consumed
	return content
}

// extractOversizeTokens parses the token count from an injection-oversize
// flag file.
//
// Parameters:
//   - data: raw bytes from the flag file
//
// Returns:
//   - int: parsed token count, or 0 if not found
func extractOversizeTokens(data []byte) int {
	m := regex.OversizeTokens.FindSubmatch(data)
	if m == nil {
		return 0
	}
	n, parseErr := strconv.Atoi(string(m[1]))
	if parseErr != nil {
		return 0
	}
	return n
}
