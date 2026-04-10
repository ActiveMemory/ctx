//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
	errMcp "github.com/ActiveMemory/ctx/internal/err/mcp"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
)

// SteeringGet returns applicable steering files for the given prompt.
// If prompt is empty, returns only "always" inclusion files.
//
// Parameters:
//   - d: runtime dependencies (unused, kept for signature uniformity)
//   - prompt: optional prompt text for auto-inclusion matching
//
// Returns:
//   - string: formatted list of matching steering files
//   - error: steering load error
func SteeringGet(_ *entity.MCPDeps, prompt string) (string, error) {
	steeringDir := rc.SteeringDir()

	files, loadErr := steering.LoadAll(steeringDir)
	if loadErr != nil {
		if errors.Is(loadErr, os.ErrNotExist) {
			return desc.Text(text.DescKeyMCPSteeringNoFiles), nil
		}
		return "", loadErr
	}

	if len(files) == 0 {
		return desc.Text(text.DescKeyMCPSteeringNoFiles), nil
	}

	filtered := steering.Filter(files, prompt, nil, "")

	if len(filtered) == 0 {
		return desc.Text(text.DescKeyMCPSteeringNoMatch), nil
	}

	var sb strings.Builder
	for _, sf := range filtered {
		ctxIo.SafeFprintf(&sb,
			desc.Text(text.DescKeyMCPSteeringSection),
			sf.Name, sf.Body)
	}

	return sb.String(), nil
}

// Search searches across all .context/ files for the given query.
// Returns matching excerpts with file paths and line numbers.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//   - query: search text to find in context files
//
// Returns:
//   - string: formatted search results with paths and line numbers
//   - error: directory read error
func Search(d *entity.MCPDeps, query string) (string, error) {
	if query == "" {
		return "", errMcp.QueryRequired()
	}

	entries, readErr := os.ReadDir(d.ContextDir)
	if readErr != nil {
		return "", errMcp.SearchRead(d.ContextDir, readErr)
	}

	queryLower := strings.ToLower(query)
	var sb strings.Builder
	matches := 0

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(d.ContextDir, e.Name())
		data, err := ctxIo.SafeReadUserFile(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), queryLower) {
				ctxIo.SafeFprintf(&sb,
					desc.Text(text.DescKeyMCPSearchHitLine),
					e.Name(), lineNum, line)
				matches++
			}
		}
	}

	if matches == 0 {
		return fmt.Sprintf(
			desc.Text(text.DescKeyMCPSearchNoMatch),
			query, d.ContextDir), nil
	}

	return sb.String(), nil
}
