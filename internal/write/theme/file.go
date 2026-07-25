//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package theme

import (
	"os"
	"path/filepath"

	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// ensureFile creates <contextDir>/<noun>/<stem>.md with an H1 when it
// does not already exist. An existing file is left untouched — a theme
// may be re-declared to revise its gist, and its accumulated bodies must
// survive that.
//
// Parameters:
//   - contextDir: the context directory
//   - noun: the kind's theme-file subdirectory
//   - stem: the theme file's basename stem
//   - name: the theme's human-readable name, used as the H1
//
// Returns:
//   - error: an IO error, or nil
func ensureFile(contextDir, noun, stem, name string) error {
	dir := filepath.Join(contextDir, noun)
	if mkErr := internalIo.SafeMkdirAll(dir, cfgFs.PermExec); mkErr != nil {
		return mkErr
	}
	path := filepath.Join(dir, stem+cfgFile.ExtMarkdown)
	if _, statErr := os.Stat(path); statErr == nil {
		return nil
	}
	header := token.HeadingLevelOneStart + name + token.NewlineLF
	return internalIo.SafeWriteFile(path, []byte(header), cfgFs.PermFile)
}
