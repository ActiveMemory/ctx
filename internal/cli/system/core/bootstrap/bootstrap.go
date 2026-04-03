//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	initCore "github.com/ActiveMemory/ctx/internal/cli/initialize/core/plugin"
	"github.com/ActiveMemory/ctx/internal/config/bootstrap"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// PluginWarning returns a warning string if the ctx plugin is installed
// but not enabled in either global or local settings.
//
// Returns:
//   - string: warning message, or empty string if no warning is needed.
func PluginWarning() string {
	if !initCore.Installed() {
		return ""
	}
	if initCore.EnabledGlobally() || initCore.EnabledLocally() {
		return ""
	}
	return desc.Text(text.DescKeyBootstrapPluginWarning)
}

// ListContextFiles reads the given directory and returns sorted .md filenames.
//
// Parameters:
//   - dir: absolute path to the context directory.
//
// Returns:
//   - []string: sorted list of Markdown filenames, or nil on read error.
func ListContextFiles(dir string) []string {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), file.ExtMarkdown) {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files
}

// WrapFileList formats file names as a comma-separated list, wrapping lines
// at approximately maxWidth characters. Continuation lines are prefixed with
// the given indent string. Returns the "none" label from assets when the
// list is empty.
//
// Parameters:
//   - files: list of filenames to format.
//   - maxWidth: approximate character width before wrapping.
//   - indent: prefix string for each line.
//
// Returns:
//   - string: formatted, wrapped file list.
func WrapFileList(files []string, maxWidth int, indent string) string {
	if len(files) == 0 {
		return indent + desc.Text(text.DescKeyBootstrapNone)
	}

	var lines []string
	current := indent

	for i, f := range files {
		entry := f
		if i < len(files)-1 {
			entry += token.Comma
		}

		switch {
		case current == indent:
			// First entry on this line - always add it.
			current += entry
		case len(current)+1+len(entry) > maxWidth:
			// Would exceed width - start a new line.
			lines = append(lines, current)
			current = indent + entry
		default:
			current += token.Space + entry
		}
	}
	lines = append(lines, current)
	return strings.Join(lines, token.NewlineLF)
}

// ParseNumberedLines splits a numbered multiline string into individual
// items, stripping the leading "N. " prefix from each line. Empty lines
// are skipped.
//
// Parameters:
//   - body: multiline string with "1. ...\n2. ..." formatting
//
// Returns:
//   - []string: list of items with number prefixes removed
func ParseNumberedLines(body string) []string {
	lines := strings.Split(body, token.NewlineLF)
	items := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if idx := strings.Index(line, bootstrap.NumberedListSep); idx >= 0 &&
			idx <= bootstrap.NumberedListMaxDigits {
			line = line[idx+len(bootstrap.NumberedListSep):]
		}
		items = append(items, line)
	}
	return items
}
