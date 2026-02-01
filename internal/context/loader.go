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
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &NotFoundError{Dir: dir}
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, &NotFoundError{Dir: dir}
	}

	ctx := &Context{
		Dir:   dir,
		Files: []FileInfo{},
	}

	// Read all .md files in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".md" {
			continue
		}

		filePath := filepath.Join(dir, name)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
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
