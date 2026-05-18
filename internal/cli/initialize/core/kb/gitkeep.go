//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kb

import (
	"os"
	"path/filepath"
	"strings"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgInitKB "github.com/ActiveMemory/ctx/internal/config/initialize/kb"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errInitKB "github.com/ActiveMemory/ctx/internal/err/initialize/kb"
	"github.com/ActiveMemory/ctx/internal/io"
)

// PlaceGitkeep adds a .gitkeep stub to dir when the directory
// is empty (so an empty kb-state subdirectory survives git
// add).
//
// Parameters:
//   - dir: directory to stub.
//
// Returns:
//   - error: wrapped error on read / write failure.
func PlaceGitkeep(dir string) error {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return errInitKB.ReadDir(dir, readErr)
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), token.Dot) {
			return nil
		}
		if e.Name() != cfgInitKB.GitkeepName && !e.IsDir() {
			return nil
		}
	}
	path := filepath.Join(dir, cfgInitKB.GitkeepName)
	if _, statErr := io.SafeStat(path); statErr == nil {
		return nil
	}
	if writeErr := io.SafeWriteFile(
		path, nil, cfgFs.PermSecret,
	); writeErr != nil {
		return errInitKB.WriteFile(path, writeErr)
	}
	return nil
}
