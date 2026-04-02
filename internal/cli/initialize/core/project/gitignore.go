//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// EnsureGitignoreEntries appends recommended .gitignore entries that are not
// already present. Creates .gitignore if it does not exist.
//
// Parameters:
//   - cmd: Cobra command for status output
//
// Returns:
//   - error: Non-nil on read or write failure
func EnsureGitignoreEntries(cmd *cobra.Command) error {
	content, readErr := io.SafeReadUserFile(file.FileGitignore)
	if readErr != nil && !os.IsNotExist(readErr) {
		return readErr
	}

	// Build set of existing trimmed lines.
	existing := make(map[string]bool)
	for _, line := range strings.Split(string(content), token.NewlineLF) {
		existing[strings.TrimSpace(line)] = true
	}

	// Collect missing entries.
	var missing []string
	for _, e := range file.Gitignore {
		if !existing[e] {
			missing = append(missing, e)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	// Build block to append.
	var sb strings.Builder
	if len(content) > 0 && !strings.HasSuffix(string(content), token.NewlineLF) {
		sb.WriteString(token.NewlineLF)
	}
	sb.WriteString(token.NewlineLF + file.GitignoreHeader + token.NewlineLF)
	for _, e := range missing {
		sb.WriteString(e + token.NewlineLF)
	}

	if writeErr := io.SafeWriteFile(
		file.FileGitignore, append(content, []byte(sb.String())...),
		fs.PermFile,
	); writeErr != nil {
		return writeErr
	}

	initialize.InfoGitignoreUpdated(cmd, len(missing))
	initialize.InfoGitignoreReview(cmd)
	return nil
}
