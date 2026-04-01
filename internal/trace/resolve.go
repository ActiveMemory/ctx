//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/task"
)

// Reference type constants used by parseRef and Resolve.
const (
	refTypeNote     = "note"
	refTypeSession  = "session"
	refTypeDecision = "decision"
	refTypeLearning = "learning"
)

// parseRef breaks a raw reference string into its components.
//
// Formats:
//   - "decision:12"  → ("decision", 12, "")
//   - "session:abc"  → ("session", 0, "abc")
//   - `"Some note"`  → ("note", 0, "Some note")
//   - unknown        → ("note", 0, ref)
//
// Parameters:
//   - ref: raw reference string
//
// Returns:
//   - refType: type keyword (decision, learning, convention, task, session, note)
//   - number: numeric value, 0 when not applicable
//   - text: text value, empty when not applicable
func parseRef(ref string) (refType string, number int, text string) {
	// Quoted strings are free-form notes.
	if strings.HasPrefix(ref, `"`) && strings.HasSuffix(ref, `"`) {
		return refTypeNote, 0, strings.Trim(ref, `"`)
	}

	parts := strings.SplitN(ref, ":", 2)
	if len(parts) != 2 {
		return refTypeNote, 0, ref
	}

	kind := parts[0]
	value := parts[1]

	switch kind {
	case refTypeDecision, refTypeLearning, "convention", "task":
		n, err := strconv.Atoi(value)
		if err != nil {
			return refTypeNote, 0, ref
		}
		return kind, n, ""
	case refTypeSession:
		return refTypeSession, 0, value
	default:
		return refTypeNote, 0, ref
	}
}

// Resolve looks up a raw reference and returns its full details.
//
// Parameters:
//   - ref: raw reference string (e.g. "decision:12", "task:8", `"Some note"`)
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - ResolvedRef: resolved reference with title, detail, and found status
func Resolve(ref, contextDir string) ResolvedRef {
	refType, number, text := parseRef(ref)

	resolved := ResolvedRef{
		Raw:    ref,
		Type:   refType,
		Number: number,
	}

	switch refType {
	case refTypeDecision:
		return resolveEntry(resolved, contextDir, cfgCtx.Decision, number)
	case refTypeLearning:
		return resolveEntry(resolved, contextDir, cfgCtx.Learning, number)
	case "convention":
		return resolveEntry(resolved, contextDir, cfgCtx.Convention, number)
	case "task":
		return resolveTask(resolved, contextDir, number)
	case refTypeSession:
		resolved.Title = text
		resolved.Found = true
		return resolved
	default: // refTypeNote
		resolved.Title = text
		resolved.Found = true
		return resolved
	}
}

// resolveEntry reads the specified context file, parses entry headers,
// and finds the entry at the given 1-based position.
//
// Parameters:
//   - resolved: partially populated ResolvedRef (Raw, Type, Number already set)
//   - contextDir: absolute path to the .context/ directory
//   - fileName: context file name (e.g. "DECISIONS.md")
//   - number: 1-based entry number
//
// Returns:
//   - ResolvedRef: populated with Title and Detail if found
func resolveEntry(resolved ResolvedRef, contextDir, fileName string, number int) ResolvedRef {
	path := filepath.Clean(filepath.Join(contextDir, fileName))

	//nolint:gosec // path built from trusted contextDir + constant filename
	content, err := os.ReadFile(path)
	if err != nil {
		return resolved
	}

	entries := index.ParseHeaders(string(content))

	// Entries are 1-based; index into slice using number-1.
	if number < 1 || number > len(entries) {
		return resolved
	}

	entry := entries[number-1]
	resolved.Title = entry.Title
	resolved.Detail = "Date: " + entry.Date
	resolved.Found = true

	return resolved
}

// resolveTask reads TASKS.md and finds the nth top-level task (1-based),
// counting both pending and completed tasks sequentially.
//
// Parameters:
//   - resolved: partially populated ResolvedRef (Raw, Type, Number already set)
//   - contextDir: absolute path to the .context/ directory
//   - number: 1-based task number
//
// Returns:
//   - ResolvedRef: populated with Title and Detail if found
func resolveTask(resolved ResolvedRef, contextDir string, number int) ResolvedRef {
	path := filepath.Clean(filepath.Join(contextDir, cfgCtx.Task))

	//nolint:gosec // path built from trusted contextDir + constant filename
	f, err := os.Open(path)
	if err != nil {
		return resolved
	}
	defer func() { _ = f.Close() }()

	count := 0
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		m := regex.Task.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// Skip subtasks (indented).
		if task.Sub(m) {
			continue
		}

		count++
		if count == number {
			status := "pending"
			if task.Completed(m) {
				status = "completed"
			}
			resolved.Title = task.Content(m)
			resolved.Detail = "Status: " + status
			resolved.Found = true
			return resolved
		}
	}

	return resolved
}
