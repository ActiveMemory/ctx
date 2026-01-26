//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package templates provides embedded template files for initializing .context/ directories.
package templates

import "embed"

//go:embed *.md entry-templates/*.md claude/commands/*.md claude/hooks/*.sh
var FS embed.FS

// GetTemplate reads a template file by name from the embedded filesystem.
func GetTemplate(name string) ([]byte, error) {
	return FS.ReadFile(name)
}

// ListTemplates returns all available template file names.
func ListTemplates() ([]string, error) {
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

// ListEntryTemplates returns available entry template file names.
func ListEntryTemplates() ([]string, error) {
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

// GetEntryTemplate reads an entry template by name.
func GetEntryTemplate(name string) ([]byte, error) {
	return FS.ReadFile("entry-templates/" + name)
}

// ListClaudeCommands returns available Claude Code slash command file names.
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

// GetClaudeCommand reads a Claude Code slash command template by name.
func GetClaudeCommand(name string) ([]byte, error) {
	return FS.ReadFile("claude/commands/" + name)
}

// GetClaudeHook reads a Claude Code hook script by name.
func GetClaudeHook(name string) ([]byte, error) {
	return FS.ReadFile("claude/hooks/" + name)
}
