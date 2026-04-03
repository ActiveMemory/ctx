//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errTrace "github.com/ActiveMemory/ctx/internal/err/trace"
	"github.com/ActiveMemory/ctx/internal/exec/git"
	"github.com/ActiveMemory/ctx/internal/trace"
	writeTrace "github.com/ActiveMemory/ctx/internal/write/trace"
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
	parts := strings.SplitN(suffix, token.Dash, 2)
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
func Trace(
	cmd *cobra.Command,
	filePath string,
	last int,
	traceDir string,
) error {
	out, err := git.Run(
		cfgGit.Log, fmt.Sprintf("-%d", last),
		cfgGit.FormatHashDateSubj,
		cfgGit.FlagPathSep, filePath,
	)
	if err != nil {
		return errTrace.GitLog(err)
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
		refStr := desc.Text(text.DescKeyWriteTraceNoRefs)
		if len(refs) > 0 {
			refStr = desc.Text(text.DescKeyWriteTraceRefsPrefix) +
				strings.Join(refs, token.CommaSpace)
		}

		writeTrace.FileEntry(cmd, trace.ShortHash(hash), date, subject, refStr)
	}

	return nil
}
