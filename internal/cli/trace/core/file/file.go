//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/trace"
)

// ParsePathArg strips an optional :line-range suffix from a path argument
// so git log gets the clean file path.
//
// Examples:
//
//	"src/auth.go:42-60"  → "src/auth.go"
//	"src/auth.go:42"     → "src/auth.go"
//	"src/auth.go"        → "src/auth.go"
//	"src/auth.go:latest" → "src/auth.go:latest"
//
// Parameters:
//   - arg: combined path and optional line range argument
//
// Returns:
//   - string: file path with line range stripped
func ParsePathArg(arg string) string {
	idx := strings.LastIndex(arg, ":")
	if idx < 0 {
		return arg
	}
	suffix := arg[idx+1:]
	// Check if suffix looks like a line range (digits or digits-digits)
	parts := strings.SplitN(suffix, "-", 2)
	for _, p := range parts {
		for _, c := range p {
			if c < '0' || c > '9' {
				return arg // not a line range, return as-is
			}
		}
		if p == "" {
			return arg
		}
	}
	return arg[:idx]
}

// TraceFile runs git log to retrieve commits touching the given file and
// prints context refs for each commit.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - filePath: clean file path (no line-range suffix)
//   - last: maximum number of commits to show
//   - traceDir: absolute path to the trace directory
//
// Returns:
//   - error: non-nil on execution failure
func TraceFile(cmd *cobra.Command, filePath string, last int, traceDir string) error {
	gitArgs := []string{"log", fmt.Sprintf("-%d", last), "--format=%H %ci %s", "--", filePath}

	//nolint:gosec // gitArgs built from integer flag + user file path, standard git usage
	out, err := exec.Command("git", gitArgs...).Output()
	if err != nil {
		return fmt.Errorf("git log: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), token.NewlineLF)

	for _, line := range lines {
		if line == "" {
			continue
		}
		// format: <hash> <date> <subject>
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 1 {
			continue
		}
		hash := parts[0]
		date := ""
		if len(parts) > 1 {
			date = parts[1]
		}
		subject := ""
		if len(parts) > 2 {
			subject = parts[2]
		}

		refs := trace.CollectRefsForCommit(hash, traceDir, false)
		refStr := "(none)"
		if len(refs) > 0 {
			refStr = "\u2192 " + strings.Join(refs, ", ")
		}

		cmd.Println(fmt.Sprintf("%s  %s  %s  [%s]", trace.ShortHash(hash), date, subject, refStr))
	}

	return nil
}
