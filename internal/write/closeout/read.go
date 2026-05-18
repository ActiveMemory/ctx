//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errCloseout "github.com/ActiveMemory/ctx/internal/err/closeout"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Read parses a closeout file at the given path. Returns the
// parsed frontmatter and the body bytes (everything after the
// closing `---`).
//
// Parameters:
//   - path: absolute path to a closeout markdown file.
//
// Returns:
//   - File: parsed frontmatter + body.
//   - error: [errCloseout.ErrMissingFrontmatter] if the file
//     does not open with `---`; [errCloseout.ErrMissingFields]
//     when required fields are absent; wrapped errors for I/O
//     or YAML parse failures.
func Read(path string) (File, error) {
	raw, ioErr := io.SafeReadUserFile(path)
	if ioErr != nil {
		return File{}, errCloseout.ReadFailed(ioErr)
	}
	fm, body, splitErr := splitFrontmatter(string(raw))
	if splitErr != nil {
		return File{}, splitErr
	}
	if fieldsErr := requireFields(fm); fieldsErr != nil {
		return File{}, fieldsErr
	}
	return File{
		Path:        path,
		Frontmatter: fm,
		Body:        body,
	}, nil
}

// List enumerates closeouts in dir, parses each, and returns
// the parsed File slice sorted by GeneratedAt ascending.
//
// Files that fail to parse are returned in the second slice
// (their paths) so the caller can surface a doctor-style
// warning. List returns nil error on file-level parse failures;
// only directory-walk failures bubble up.
//
// Parameters:
//   - closeoutsDir: absolute path to scan.
//
// Returns:
//   - []File: successfully-parsed closeouts (sorted).
//   - []string: paths of files that failed to parse.
//   - error: non-nil only on directory enumeration failure.
func List(closeoutsDir string) ([]File, []string, error) {
	entries, readErr := os.ReadDir(closeoutsDir)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil, nil
		}
		return nil, nil, errCloseout.ReadCloseoutsDir(readErr)
	}

	var ok []File
	var bad []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, cfgKB.CloseoutSuffix) {
			continue
		}
		path := filepath.Join(closeoutsDir, name)
		f, parseErr := Read(path)
		if parseErr != nil {
			bad = append(bad, path)
			continue
		}
		ok = append(ok, f)
	}
	sort.Slice(ok, func(i, j int) bool {
		return ok[i].Frontmatter.GeneratedAt.Before(
			ok[j].Frontmatter.GeneratedAt,
		)
	})
	return ok, bad, nil
}

// PostdatedBy returns the subset of files whose GeneratedAt is
// strictly after cursor. Used by the handover-fold mechanism
// to find closeouts produced since the last handover.
//
// Parameters:
//   - files: pre-parsed closeouts (typically from List).
//   - cursor: handover-fold cursor (latest handover's
//     generated-at, or zero value for "fold everything").
//
// Returns:
//   - []File: files with GeneratedAt > cursor, order preserved.
func PostdatedBy(files []File, cursor time.Time) []File {
	var out []File
	for _, f := range files {
		if f.Frontmatter.GeneratedAt.After(cursor) {
			out = append(out, f)
		}
	}
	return out
}
