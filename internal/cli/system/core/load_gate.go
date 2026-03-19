//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ExtractIndex returns the content between INDEX:START and INDEX:END
// markers within a context file.
//
// Parameters:
//   - content: full file content to search
//
// Returns:
//   - string: trimmed index content, or empty string if markers are
//     not found or improperly ordered
func ExtractIndex(content string) string {
	start := strings.Index(content, marker.IndexStart)
	end := strings.Index(content, marker.IndexEnd)
	if start < 0 || end < 0 || end <= start {
		return ""
	}
	startPos := start + len(marker.IndexStart)
	return strings.TrimSpace(content[startPos:end])
}

// WriteOversizeFlag writes an injection-oversize flag file when the total
// injected tokens exceed the configured threshold. The flag file is read
// by check-context-size to emit an oversize warning.
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//   - totalTokens: total injected token count
//   - perFile: per-file token breakdown for diagnostics
func WriteOversizeFlag(
	contextDir string, totalTokens int, perFile []FileTokenEntry,
) {
	threshold := rc.InjectionTokenWarn()
	if threshold == 0 || totalTokens <= threshold {
		return
	}

	sd := filepath.Join(contextDir, dir.State)
	_ = os.MkdirAll(sd, fs.PermRestrictedDir)

	var flag strings.Builder
	flag.WriteString(desc.TextDesc(text.DescKeyContextLoadGateOversizeHeader))
	flag.WriteString(strings.Repeat("=", stats.ContextSizeOversizeSepLen) + token.NewlineLF)
	flag.WriteString(fmt.Sprintf(
		desc.TextDesc(text.DescKeyContextLoadGateOversizeTimestamp),
		time.Now().UTC().Format(time.RFC3339)))
	flag.WriteString(fmt.Sprintf(
		desc.TextDesc(text.DescKeyContextLoadGateOversizeInjected),
		totalTokens, threshold))
	flag.WriteString(desc.TextDesc(text.DescKeyContextLoadGateOversizeBreakdown))
	for _, entry := range perFile {
		flag.WriteString(fmt.Sprintf(
			desc.TextDesc(text.DescKeyContextLoadGateOversizeFileEntry),
			entry.Name, entry.Tokens))
	}
	flag.WriteString(token.NewlineLF)
	flag.WriteString(desc.TextDesc(text.DescKeyContextLoadGateOversizeAction))

	_ = os.WriteFile(
		filepath.Join(sd, stats.ContextSizeInjectionOversizeFlag),
		[]byte(flag.String()), fs.PermSecret)
}
