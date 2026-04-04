//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"io"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// copyDir recursively copies the contents of src into dst.
// Both directories must already exist.
func copyDir(src, dst string) error {
	entries, readErr := os.ReadDir(src)
	if readErr != nil {
		return readErr
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if mkdirErr := ctxIo.SafeMkdirAll(
				dstPath, fs.PermRestrictedDir,
			); mkdirErr != nil {
				return mkdirErr
			}
			if recurseErr := copyDir(srcPath, dstPath); recurseErr != nil {
				return recurseErr
			}
			continue
		}

		if cpErr := copyFile(srcPath, dstPath); cpErr != nil {
			return cpErr
		}
	}
	return nil
}

// copyFile copies a single file from src to dst, preserving permissions.
func copyFile(src, dst string) error {
	info, statErr := ctxIo.SafeStat(src)
	if statErr != nil {
		return statErr
	}

	in, openErr := ctxIo.SafeOpenUserFile(src)
	if openErr != nil {
		return openErr
	}
	defer func() { _ = in.Close() }()

	out, createErr := ctxIo.SafeCreateFile(dst, info.Mode().Perm())
	if createErr != nil {
		return createErr
	}
	defer func() { _ = out.Close() }()

	_, copyErr := io.Copy(out, in)
	return copyErr
}
