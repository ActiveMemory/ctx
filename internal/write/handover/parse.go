//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errHandover "github.com/ActiveMemory/ctx/internal/err/handover"
	"github.com/ActiveMemory/ctx/internal/io"
)

// listHandovers enumerates handover files in handoversDir,
// parses each, and returns the parsed slice sorted by
// GeneratedAt ascending. Files that fail to parse are
// silently skipped; the caller can detect that case by
// comparing len(returned) against the dir listing.
//
// Parameters:
//   - handoversDir: absolute path.
//
// Returns:
//   - []File: parsed handovers, ascending by GeneratedAt.
//   - error: non-nil only on directory enumeration failure.
func listHandovers(handoversDir string) ([]File, error) {
	entries, readErr := os.ReadDir(handoversDir)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errHandover.ReadHandoversDir(readErr)
	}
	var out []File
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), cfgFile.ExtMarkdown) {
			continue
		}
		path := filepath.Join(handoversDir, e.Name())
		f, parseErr := readFile(path)
		if parseErr != nil {
			continue
		}
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Frontmatter.GeneratedAt.Before(
			out[j].Frontmatter.GeneratedAt,
		)
	})
	return out, nil
}

// readFile reads + parses a handover file at path.
//
// Parameters:
//   - path: absolute path to a handover markdown file.
//
// Returns:
//   - File: parsed frontmatter + body.
//   - error: wrapped errors on I/O or YAML parse failures.
func readFile(path string) (File, error) {
	raw, ioErr := io.SafeReadUserFile(path)
	if ioErr != nil {
		return File{}, errHandover.ReadFailed(ioErr)
	}
	lines := strings.SplitN(string(raw), token.NewlineLF, 2)
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != token.Separator {
		return File{}, errHandover.ErrMissingFrontmatter
	}
	openClose := token.NewlineLF + token.Separator + token.NewlineLF
	idx := strings.Index(lines[1], openClose)
	if idx < 0 {
		idx = strings.Index(lines[1], token.NewlineLF+token.Separator)
		if idx < 0 {
			return File{}, errHandover.ErrMissingClosingDelim
		}
	}
	header := lines[1][:idx]
	bodyStart := idx + len(openClose)
	if bodyStart > len(lines[1]) {
		bodyStart = len(lines[1])
	}
	body := strings.TrimLeft(lines[1][bodyStart:], token.NewlineLF)

	var fm Frontmatter
	if yamlErr := yaml.Unmarshal([]byte(header), &fm); yamlErr != nil {
		return File{}, errHandover.ParseFrontmatter(yamlErr)
	}
	if fm.GeneratedAt.IsZero() {
		return File{}, errHandover.ErrMissingGeneratedAt
	}
	return File{Path: path, Frontmatter: fm, Body: body}, nil
}
