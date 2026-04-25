//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"errors"
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
)

// oversizeContent checks for an injection-oversize flag file and
// returns the raw nudge content if present. Deletes the flag after
// reading (one-shot).
//
// ctxDir is supplied by the caller (FullPreamble-equivalent gate
// above) so this helper does not re-resolve it; a second resolution
// would be dead code today and would pair an ambiguous (zero, err)
// return with the legitimate "nothing to do" result.
//
// Parameters:
//   - ctxDir: absolute path to the context directory
//
// Returns:
//   - string: raw oversize nudge content, or empty string when
//     there is no flag to report or the template silences itself.
//   - error: non-nil when the flag file cannot be read (permission,
//     I/O) or cannot be removed after reading. Legitimate "nothing
//     to do" paths return ("", nil): flag file absent
//     (os.ErrNotExist). A remove failure returns ("", err) rather
//     than (content, err): if we cannot clear the one-shot flag we
//     must not emit the nudge either, otherwise the flag re-fires on
//     every subsequent invocation and the operator sees a nudge
//     storm. Log-first principle: don't emit a user-visible nudge
//     whose persistence cleanup we could not verify.
func oversizeContent(ctxDir string) (string, error) {
	baseDir := filepath.Join(ctxDir, dir.State)
	flagPath := filepath.Join(baseDir, stats.ContextSizeInjectionOversizeFlag)
	data, readErr := io.SafeReadFile(
		baseDir, stats.ContextSizeInjectionOversizeFlag,
	)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			// No flag on disk ⇒ nothing to report; legitimate.
			return "", nil
		}
		// Permission denied, I/O failure: surface.
		return "", readErr
	}

	tokenCount := extractOversizeTokens(data)
	fallback := fmt.Sprintf(
		desc.Text(text.DescKeyCheckContextSizeOversizeFallback), tokenCount,
	)
	content := message.Load(hook.CheckContextSize, hook.VariantOversize,
		map[string]any{stats.VarTokenCount: tokenCount}, fallback)

	// One-shot: remove the flag regardless of whether the template
	// silenced itself, so a silenced template does not leave the
	// flag lingering and re-firing every invocation.
	if removeErr := os.Remove(flagPath); removeErr != nil {
		return "", removeErr
	}
	return content, nil
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
