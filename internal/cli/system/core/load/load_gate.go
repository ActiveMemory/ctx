//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/load_gate"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// WriteOversizeFlag writes an injection-oversize flag file when the total
// injected tokens exceed the configured threshold. The flag file is read
// by check-context-size to emit an oversize warning.
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//   - totalTokens: total injected token count
//   - perFile: per-file token breakdown for diagnostics
func WriteOversizeFlag(
	contextDir string, totalTokens int, perFile []entity.FileTokenEntry,
) {
	threshold := rc.InjectionTokenWarn()
	if threshold == 0 || totalTokens <= threshold {
		return
	}

	sd := filepath.Join(contextDir, dir.State)
	if mkdirErr := os.MkdirAll(sd, fs.PermRestrictedDir); mkdirErr != nil {
		ctxLog.Warn(warn.Mkdir, sd, mkdirErr)
	}

	var flag strings.Builder
	flag.WriteString(desc.Text(text.DescKeyContextLoadGateOversizeHeader))
	sep := strings.Repeat(load_gate.ContextLoadSeparatorChar, stats.ContextSizeOversizeSepLen)
	flag.WriteString(sep + token.NewlineLF)
	io.SafeFprintf(&flag,
		desc.Text(text.DescKeyContextLoadGateOversizeTimestamp),
		time.Now().UTC().Format(time.RFC3339))
	io.SafeFprintf(&flag,
		desc.Text(text.DescKeyContextLoadGateOversizeInjected),
		totalTokens, threshold)
	flag.WriteString(desc.Text(text.DescKeyContextLoadGateOversizeBreakdown))
	for _, entry := range perFile {
		io.SafeFprintf(&flag,
			desc.Text(text.DescKeyContextLoadGateOversizeFileEntry),
			entry.Name, entry.Tokens)
	}
	flag.WriteString(token.NewlineLF)
	flag.WriteString(desc.Text(text.DescKeyContextLoadGateOversizeAction))

	fp := filepath.Join(
		sd, stats.ContextSizeInjectionOversizeFlag,
	)
	if writeErr := os.WriteFile(
		fp, []byte(flag.String()), fs.PermSecret,
	); writeErr != nil {
		ctxLog.Warn(warn.Write, fp, writeErr)
	}
}
