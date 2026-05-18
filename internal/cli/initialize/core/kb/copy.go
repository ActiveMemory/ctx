//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kb

import (
	"io/fs"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errInitKB "github.com/ActiveMemory/ctx/internal/err/initialize/kb"
	"github.com/ActiveMemory/ctx/internal/io"
)

// CopyEmbedTree copies every regular file at the top level of
// the embedded directory at assetPath into dstDir. Subdirs are
// not recursed (callers pass leaf directories).
//
// Parameters:
//   - assetPath: directory under assets.FS (no leading slash).
//   - dstDir: destination directory on disk.
//
// Returns:
//   - error: wrapped errors on read / write failure.
func CopyEmbedTree(assetPath, dstDir string) error {
	entries, readErr := fs.ReadDir(assets.FS, assetPath)
	if readErr != nil {
		return errInitKB.ReadEmbed(assetPath, readErr)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		src := assetPath + token.Slash + e.Name()
		dst := filepath.Join(dstDir, e.Name())
		if copyErr := CopyEmbedFile(src, dst); copyErr != nil {
			return copyErr
		}
	}
	return nil
}

// CopyEmbedFile copies a single embedded file to dst. Skips
// silently when dst already exists (init preserves curated
// content).
//
// Parameters:
//   - assetPath: file under assets.FS.
//   - dst: destination path.
//
// Returns:
//   - error: wrapped errors on read / write failure.
func CopyEmbedFile(assetPath, dst string) error {
	if _, statErr := io.SafeStat(dst); statErr == nil {
		return nil
	}
	raw, readErr := fs.ReadFile(assets.FS, assetPath)
	if readErr != nil {
		return errInitKB.ReadEmbed(assetPath, readErr)
	}
	mkErr := io.SafeMkdirAll(filepath.Dir(dst), cfgFs.PermExec)
	if mkErr != nil {
		return errInitKB.MkdirFor(dst, mkErr)
	}
	writeErr := io.SafeWriteFile(dst, raw, cfgFs.PermSecret)
	if writeErr != nil {
		return errInitKB.WriteFile(dst, writeErr)
	}
	return nil
}
