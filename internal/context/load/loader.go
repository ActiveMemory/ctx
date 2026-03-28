//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/context/sanitize"
	"github.com/ActiveMemory/ctx/internal/context/summary"
	"github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// Do reads all context files from the specified directory.
//
// If dir is empty, it uses the configured context directory from .ctxrc,
// CTX_DIR environment variable, or the default ".context".
//
// Parameters:
//   - dir: Directory path to load from, or empty string for default
//
// Returns:
//   - *Context: Loaded context with files, token counts, and metadata
//   - error: NotFoundError if directory doesn't exist, or other IO errors
func Do(dir string) (*entity.Context, error) {
	if dir == "" {
		dir = rc.ContextDir()
	}

	// Check if the directory exists
	info, statErr := os.Stat(dir)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return nil, errCtx.NotFound(dir)
		}
		return nil, statErr
	}
	if !info.IsDir() {
		return nil, errCtx.NotFound(dir)
	}

	// Reject context directories that contain symlinks (M-2 defense).
	if err := validation.CheckSymlinks(dir); err != nil {
		return nil, err
	}

	ctx := &entity.Context{
		Dir:   dir,
		Files: []entity.FileInfo{},
	}

	// Read all .md files in the directory
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil, readErr
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != file.ExtMarkdown {
			continue
		}

		filePath := filepath.Clean(filepath.Join(dir, name))
		content, readFileErr := os.ReadFile(filePath)
		if readFileErr != nil {
			continue
		}

		fileInfo, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}

		tokens := token.Estimate(content)
		fi := entity.FileInfo{
			Name:    name,
			Path:    filePath,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Content: content,
			IsEmpty: len(content) == 0 || sanitize.EffectivelyEmpty(content),
			Tokens:  tokens,
			Summary: summary.Generate(name, content),
		}

		ctx.Files = append(ctx.Files, fi)
		ctx.TotalTokens += tokens
		ctx.TotalSize += fileInfo.Size()
	}

	return ctx, nil
}
