//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tpl provides embedded template files for initializing
// .context/ directories.
package tpl

import "embed"

//go:embed *.md entry-templates/*.md claude/commands/*.md claude/hooks/*.sh
var FS embed.FS

// Template reads a template file by name from the embedded filesystem.
//
// Parameters:
//   - name: Template filename (e.g., "TASKS.md")
//
// Returns:
//   - []byte: Template content
//   - error: Non-nil if the file is not found or read fails
func Template(name string) ([]byte, error) {
	return FS.ReadFile(name)
}

// List returns all available template file names.
//
// Returns:
//   - []string: List of template filenames in the root templates directory
//   - error: Non-nil if directory read fails
func List() ([]string, error) {
	entries, err := FS.ReadDir(".")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// ListEntry returns available entry template file names.
//
// Returns:
//   - []string: List of template filenames in entry-templates/
//   - error: Non-nil if directory read fails
func ListEntry() ([]string, error) {
	entries, err := FS.ReadDir("entry-templates")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// Entry reads an entry template by name.
//
// Parameters:
//   - name: Template filename (e.g., "decision.md")
//
// Returns:
//   - []byte: Template content from entry-templates/
//   - error: Non-nil if the file is not found or read fails
func Entry(name string) ([]byte, error) {
	return FS.ReadFile("entry-templates/" + name)
}

// ListClaudeCommands returns available Claude Code slash command file names.
//
// Returns:
//   - []string: List of command filenames in claude/commands/
//   - error: Non-nil if directory read fails
func ListClaudeCommands() ([]string, error) {
	entries, err := FS.ReadDir("claude/commands")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// ClaudeCommandByName reads a Claude Code slash command template by name.
//
// Parameters:
//   - name: Command filename (e.g., "ctx-status.md")
//
// Returns:
//   - []byte: Command template content from claude/commands/
//   - error: Non-nil if the file not found or read fails
func ClaudeCommandByName(name string) ([]byte, error) {
	return FS.ReadFile("claude/commands/" + name)
}

// ClaudeHookByFileName reads a Claude Code hook script by name.
//
// Parameters:
//   - name: Hook script filename (e.g., "session-end-auto-save.sh")
//
// Returns:
//   - []byte: Hook script content from claude/hooks/
//   - error: Non-nil if the file is not found or read fails
func ClaudeHookByFileName(name string) ([]byte, error) {
	return FS.ReadFile("claude/hooks/" + name)
}
