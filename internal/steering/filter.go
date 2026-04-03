//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// LoadAll reads all .md files from the steering directory and parses
// them into SteeringFile values. Returns an error if the directory
// cannot be read or any file fails to parse.
func LoadAll(steeringDir string) ([]*SteeringFile, error) {
	entries, readDirErr := os.ReadDir(steeringDir)
	if readDirErr != nil {
		return nil, errSteering.ReadDir(steeringDir, readDirErr)
	}

	var files []*SteeringFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), file.ExtMarkdown) {
			continue
		}
		path := filepath.Join(steeringDir, e.Name())
		data, readErr := ctxIo.SafeReadUserFile(path)
		if readErr != nil {
			return nil, errSteering.ReadFile(path, readErr)
		}
		sf, parseErr := Parse(data, path)
		if parseErr != nil {
			return nil, parseErr
		}
		files = append(files, sf)
	}
	return files, nil
}

// Filter returns steering files applicable for the given context.
//
// Inclusion rules:
//   - always: included unconditionally
//   - auto: included when prompt contains the file's description
//     as a case-insensitive substring
//   - manual: included only when the file's name appears in manualNames
//
// When tool is non-empty, files whose Tools list is non-nil and
// non-empty are excluded if the list does not contain the tool.
// When tool is empty, no tool filtering is applied.
//
// Results are sorted by ascending priority, then alphabetically
// by name on tie.
func Filter(
	files []*SteeringFile, prompt string,
	manualNames []string, tool string,
) []*SteeringFile {
	promptLower := strings.ToLower(prompt)

	var result []*SteeringFile
	for _, sf := range files {
		if !matchInclusion(sf, promptLower, manualNames) {
			continue
		}
		if !matchTool(sf, tool) {
			continue
		}
		result = append(result, sf)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Priority != result[j].Priority {
			return result[i].Priority < result[j].Priority
		}
		return result[i].Name < result[j].Name
	})

	return result
}
