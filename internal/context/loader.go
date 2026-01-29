//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Packages context provides functionality for loading and managing .context/ files.
package context

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// FileInfo represents metadata about a context file.
//
// Fields:
//   - Name: Filename (e.g., "TASKS.md")
//   - Path: Full path to the file
//   - Size: File size in bytes
//   - ModTime: Last modification time
//   - Content: Raw file content
//   - IsEmpty: True if file has no meaningful content (only headers/whitespace)
//   - Tokens: Estimated token count for the content
//   - Summary: Brief description generated from the content
type FileInfo struct {
	Name    string
	Path    string
	Size    int64
	ModTime time.Time
	Content []byte
	IsEmpty bool
	Tokens  int
	Summary string
}

// Context represents the loaded context from a .context/ directory.
//
// Fields:
//   - Dir: Path to the context directory
//   - Files: All loaded context files with their metadata
//   - TotalTokens: Sum of estimated tokens across all files
//   - TotalSize: Sum of file sizes in bytes
type Context struct {
	Dir         string
	Files       []FileInfo
	TotalTokens int
	TotalSize   int64
}

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
		dir = config.GetContextDir()
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
			IsEmpty: len(content) == 0 || isEffectivelyEmpty(content),
			Tokens:  tokens,
			Summary: generateSummary(name, content),
		}

		ctx.Files = append(ctx.Files, fi)
		ctx.TotalTokens += tokens
		ctx.TotalSize += fileInfo.Size()
	}

	return ctx, nil
}

// Exists checks if a context directory exists.
//
// If dir is empty, it uses the configured context directory.
//
// Parameters:
//   - dir: Directory path to check, or empty string for default
//
// Returns:
//   - bool: True if the directory exists and is a directory
func Exists(dir string) bool {
	if dir == "" {
		dir = config.GetContextDir()
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

// NotFoundError is returned when the context directory doesn't exist.
type NotFoundError struct {
	Dir string
}

// Error implements the error interface for NotFoundError.
//
// Returns:
//   - string: Error message including the missing directory path
func (e *NotFoundError) Error() string {
	return "context directory not found: " + e.Dir
}

// isEffectivelyEmpty checks if a file only contains headers and whitespace.
//
// Parameters:
//   - content: File content to check
//
// Returns:
//   - bool: True if file has no meaningful content (only headers, comments, whitespace)
func isEffectivelyEmpty(content []byte) bool {
	// Simple heuristic: if `content` is less than 100 bytes
	// and mostly headers/whitespace
	if len(content) < 20 {
		return true
	}

	// Count non-header, non-whitespace content
	lines := 0
	contentLines := 0
	for _, line := range splitLines(content) {
		lines++
		trimmed := trimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if trimmed[0] == '#' {
			continue
		}
		if len(trimmed) > 0 && trimmed[0] == '-' && len(trimmed) < 5 {
			continue
		}
		// Check for HTML comment markers
		if len(trimmed) >= 4 && string(trimmed[:4]) == "<!--" {
			continue
		}
		if len(trimmed) >= 3 && string(trimmed[len(trimmed)-3:]) == "-->" {
			continue
		}
		contentLines++
	}

	return contentLines == 0
}

// splitLines splits content into lines without using strings package.
//
// Parameters:
//   - content: Byte slice to split
//
// Returns:
//   - [][]byte: Slice of lines (without newline characters)
func splitLines(content []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range content {
		if b == '\n' {
			lines = append(lines, content[start:i])
			start = i + 1
		}
	}
	if start < len(content) {
		lines = append(lines, content[start:])
	}
	return lines
}

// trimSpace trims leading and trailing whitespace from a byte slice.
//
// Parameters:
//   - b: Byte slice to trim
//
// Returns:
//   - []byte: Slice with spaces, tabs, and carriage returns removed from ends
func trimSpace(b []byte) []byte {
	start := 0
	end := len(b)
	for start < end && (b[start] == ' ' || b[start] == '\t' || b[start] == '\r') {
		start++
	}
	for end > start && (b[end-1] == ' ' || b[end-1] == '\t' || b[end-1] == '\r') {
		end--
	}
	return b[start:end]
}
