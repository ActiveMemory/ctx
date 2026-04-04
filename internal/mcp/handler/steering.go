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

	errMcp "github.com/ActiveMemory/ctx/internal/err/mcp"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
)

// Steering result messages.
const (
	// msgNoSteeringFiles is returned when no steering files exist.
	msgNoSteeringFiles = "No steering files found."
	// msgNoMatchingSteering is returned when no files match.
	msgNoMatchingSteering = "No matching steering files."
)

// SteeringGet returns applicable steering files for the given prompt.
// If prompt is empty, returns only "always" inclusion files.
//
// Parameters:
//   - prompt: optional prompt text for auto-inclusion matching
//
// Returns:
//   - string: formatted list of matching steering files
//   - error: steering load error
func (h *Handler) SteeringGet(prompt string) (string, error) {
	steeringDir := rc.SteeringDir()

	files, loadErr := steering.LoadAll(steeringDir)
	if loadErr != nil {
		if errors.Is(loadErr, os.ErrNotExist) {
			return msgNoSteeringFiles, nil
		}
		return "", loadErr
	}

	if len(files) == 0 {
		return msgNoSteeringFiles, nil
	}

	filtered := steering.Filter(files, prompt, nil, "")

	if len(filtered) == 0 {
		return msgNoMatchingSteering, nil
	}

	var sb strings.Builder
	for _, sf := range filtered {
		fmt.Fprintf(&sb, "## %s\n\n%s\n\n", sf.Name, sf.Body)
	}

	return sb.String(), nil
}

// Search searches across all .context/ files for the given query.
// Returns matching excerpts with file paths and line numbers.
//
// Parameters:
//   - query: search text to find in context files
//
// Returns:
//   - string: formatted search results with paths and line numbers
//   - error: directory read error
func (h *Handler) Search(query string) (string, error) {
	if query == "" {
		return "", errMcp.QueryRequired()
	}

	entries, readErr := os.ReadDir(h.ContextDir)
	if readErr != nil {
		return "", errMcp.SearchRead(h.ContextDir, readErr)
	}

	queryLower := strings.ToLower(query)
	var sb strings.Builder
	matches := 0

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(h.ContextDir, e.Name())
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
				fmt.Fprintf(&sb, "%s:%d: %s\n", e.Name(), lineNum, line)
				matches++
			}
		}
	}

	if matches == 0 {
		return fmt.Sprintf("No matches for %q in %s.", query, h.ContextDir), nil
	}

	return sb.String(), nil
}
