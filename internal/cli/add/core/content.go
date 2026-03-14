//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bufio"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/write/add"
)

// ExtractContent retrieves content from various sources for adding entries.
//
// Content is extracted in priority order:
//  1. From the file specified by --file flag
//  2. From command line arguments (after the entry type)
//  3. From stdin (if piped)
//
// Parameters:
//   - args: Command arguments where args[1:] may contain inline content
//   - flags: Configuration flags including FromFile path
//
// Returns:
//   - string: Extracted and trimmed content
//   - error: Non-nil if no content source is available or reading fails
func ExtractContent(args []string, flags Config) (string, error) {
	if flags.FromFile != "" {
		// Read from the file
		fileContent, err := os.ReadFile(flags.FromFile)
		if err != nil {
			return "", add.ErrFileRead(flags.FromFile, err)
		}
		return strings.TrimSpace(string(fileContent)), nil
	}

	if len(args) > 1 {
		// Content from arguments
		return strings.Join(args[1:], " "), nil
	}

	// Try reading from stdin (check if it's a pipe)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// stdin is a pipe, read from it
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", add.ErrStdinRead(err)
		}
		return strings.TrimSpace(strings.Join(lines, token.NewlineLF)), nil
	}
	return "", add.ErrNoContent()
}
