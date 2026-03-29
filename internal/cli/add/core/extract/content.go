//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"bufio"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
)

// Content retrieves content from various sources for adding entries.
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
func Content(args []string, flags entity.AddConfig) (string, error) {
	if flags.FromFile != "" {
		// Read from the file
		fileContent, readErr := os.ReadFile(flags.FromFile)
		if readErr != nil {
			return "", errFs.FileRead(flags.FromFile, readErr)
		}
		return strings.TrimSpace(string(fileContent)), nil
	}

	if len(args) > 1 {
		// Content from arguments
		return strings.Join(args[1:], token.Space), nil
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
		if scanErr := scanner.Err(); scanErr != nil {
			return "", errFs.StdinRead(scanErr)
		}
		return strings.TrimSpace(strings.Join(lines, token.NewlineLF)), nil
	}
	return "", errAdd.NoContent()
}
