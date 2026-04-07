//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/parser"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/journal/schema"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Check runs schema validation across JSONL files in the
// directories determined by opts.
//
// Parameters:
//   - opts: flags controlling which directories to scan
//
// Returns:
//   - *schema.Collector: accumulated validation findings
//   - error: non-nil if scanning fails entirely
func Check(opts CheckOpts) (*schema.Collector, error) {
	s := schema.Default()
	c := schema.NewCollector(s.Version)

	dirs := ScanDirs(opts)
	if len(dirs) == 0 {
		return c, nil
	}

	for _, scanDir := range dirs {
		walkErr := filepath.Walk(scanDir, func(
			path string, info os.FileInfo,
			walkItemErr error,
		) error {
			if walkItemErr != nil {
				return walkItemErr
			}
			if info.IsDir() {
				if info.Name() == parser.DirSubagents {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, file.ExtJSONL) {
				return nil
			}
			sep := string(filepath.Separator)
			subPath := sep + parser.DirSubagents + sep
			if strings.Contains(path, subPath) {
				return nil
			}
			validateErr := schema.ValidateFile(
				s, path, c,
			)
			if validateErr != nil {
				ctxLog.Warn(
					warn.Walk, path, validateErr,
				)
			}
			return nil
		})
		if walkErr != nil {
			ctxLog.Warn(warn.Walk, scanDir, walkErr)
		}
	}

	return c, nil
}

// CheckSessions validates the source files of the given sessions.
//
// Parameters:
//   - sessions: sessions whose source files to validate
//
// Returns:
//   - *schema.Collector: accumulated validation findings
func CheckSessions(
	sessions []*entity.Session,
) *schema.Collector {
	s := schema.Default()
	c := schema.NewCollector(s.Version)

	seen := make(map[string]bool)
	for _, sess := range sessions {
		if sess.SourceFile == "" || seen[sess.SourceFile] {
			continue
		}
		seen[sess.SourceFile] = true
		validateErr := schema.ValidateFile(
			s, sess.SourceFile, c,
		)
		if validateErr != nil {
			ctxLog.Warn(
				warn.Walk, sess.SourceFile, validateErr,
			)
		}
	}

	return c
}

// ScanDirs resolves the directories to scan based on flags.
//
// Parameters:
//   - opts: flags controlling directory resolution
//
// Returns:
//   - []string: directories to scan for JSONL files
func ScanDirs(opts CheckOpts) []string {
	if opts.Dir != "" {
		return []string{opts.Dir}
	}

	var dirs []string
	home, homeErr := os.UserHomeDir()
	if homeErr == nil {
		if !opts.AllProjects {
			cwd, cwdErr := os.Getwd()
			if cwdErr == nil {
				pd := ClaudeProjectDir(home, cwd)
				if pd != "" {
					dirs = append(dirs, pd)
				}
			}
		} else {
			dirs = append(dirs,
				filepath.Join(home, dir.Claude, dir.Projects))
		}
	}

	return dirs
}

// ClaudeProjectDir returns the Claude Code project directory
// for the given cwd.
//
// Parameters:
//   - home: user home directory
//   - cwd: current working directory
//
// Returns:
//   - string: project directory path, or empty if not found
func ClaudeProjectDir(home, cwd string) string {
	base := filepath.Join(home, dir.Claude, dir.Projects)
	entries, readErr := os.ReadDir(base)
	if readErr != nil {
		return ""
	}

	sep := string(filepath.Separator)
	cwdDashed := strings.ReplaceAll(cwd, sep, token.Dash)
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == cwdDashed {
			return filepath.Join(base, entry.Name())
		}
	}

	return ""
}

// SortedRecordTypes returns record type keys sorted.
//
// Parameters:
//   - m: record type map from the schema
//
// Returns:
//   - []string: sorted type names
func SortedRecordTypes(
	m map[string]schema.RecordSchema,
) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedBlockTypes returns block type keys sorted.
//
// Parameters:
//   - m: block type map from the schema
//
// Returns:
//   - []string: sorted block type names
func SortedBlockTypes(
	m map[string]schema.BlockKind,
) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// WriteReport writes or removes the drift report file in
// .context/reports/schema-drift.md.
//
// Parameters:
//   - c: collector with accumulated findings
//
// Returns:
//   - error: non-nil if the report cannot be written
func WriteReport(c *schema.Collector) error {
	contextDir := rc.ContextDir()
	if contextDir == "" {
		return nil
	}

	reportsDir := filepath.Join(contextDir, dir.Reports)
	reportPath := filepath.Join(reportsDir, file.SchemaDrift)

	if !c.Drift() {
		if _, statErr := os.Stat(reportPath); statErr == nil {
			return os.Remove(reportPath)
		}
		return nil
	}

	mkErr := ctxIo.SafeMkdirAll(reportsDir, fs.PermExec)
	if mkErr != nil {
		return mkErr
	}

	report := schema.Report(c)
	return ctxIo.SafeWriteFile(
		reportPath, []byte(report), fs.PermFile,
	)
}
