//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package context provides functionality for loading and managing .context/ files.
package context

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Load reads all context files from the specified directory.
//
// If dir is empty, it uses the configured context directory from .contextrc,
// CTX_DIR environment variable, or the default ".context".
//
// Parameters:
//   - dir: Directory path to load from, or empty string for default
//
// Returns:
//   - *Context: Loaded context with files, token counts, and metadata
//   - error: NotFoundError if directory doesn't exist, or other IO errors
func Load(dir string) (*Context, error) {
	if dir == "" {
		dir = rc.ContextDir()
	}

	// Check if the directory exists
	info, statErr := os.Stat(dir)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return nil, &NotFoundError{Dir: dir}
		}
		return nil, statErr
	}
	if !info.IsDir() {
		return nil, &NotFoundError{Dir: dir}
	}

	ctx := &Context{
		Dir:   dir,
		Files: []FileInfo{},
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
		if filepath.Ext(name) != config.ExtMarkdown {
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

		tokens := EstimateTokens(content)
		fi := FileInfo{
			Name:    name,
			Path:    filePath,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Content: content,
			IsEmpty: len(content) == 0 || effectivelyEmpty(content),
			Tokens:  tokens,
			Summary: generateSummary(name, content),
		}

		ctx.Files = append(ctx.Files, fi)
		ctx.TotalTokens += tokens
		ctx.TotalSize += fileInfo.Size()
	}

	return ctx, nil
}
