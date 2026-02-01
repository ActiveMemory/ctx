//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

func extractContent(args []string, flags addConfig) (string, error) {
	if flags.fromFile != "" {
		// Read from the file
		fileContent, err := os.ReadFile(flags.fromFile)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", flags.fromFile, err)
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
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}
		return strings.TrimSpace(strings.Join(lines, config.NewlineLF)), nil
	}
	return "", fmt.Errorf("no content provided")
}
