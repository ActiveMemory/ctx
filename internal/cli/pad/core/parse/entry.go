//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// EntriesWithIDs parses raw bytes into entries with stable
// IDs. Lines with valid `[N] ` prefixes keep their IDs.
// Lines without prefixes are assigned the next available ID.
// Duplicate IDs are reassigned.
//
// Parameters:
//   - data: Raw scratchpad content
//
// Returns:
//   - []Entry: Parsed entries with stable IDs
func EntriesWithIDs(data []byte) []Entry {
	if len(data) == 0 {
		return nil
	}

	lines := strings.Split(string(data), token.NewlineLF)

	// First pass: parse IDs and find max.
	type raw struct {
		id      int
		content string
		hasID   bool
	}
	var raws []raw
	seen := make(map[int]bool)
	maxID := 0

	for _, line := range lines {
		if line == "" {
			continue
		}
		match := regex.PadEntryID.FindStringSubmatch(line)
		if match != nil {
			id, _ := strconv.Atoi(match[1])
			content := line[len(match[0]):]
			if seen[id] {
				// Duplicate: treat as no ID.
				raws = append(raws, raw{
					content: content, hasID: false,
				})
			} else {
				seen[id] = true
				if id > maxID {
					maxID = id
				}
				raws = append(raws, raw{
					id: id, content: content, hasID: true,
				})
			}
		} else {
			raws = append(raws, raw{
				content: line, hasID: false,
			})
		}
	}

	// Second pass: assign IDs to entries without them.
	entries := make([]Entry, 0, len(raws))
	nextID := maxID + 1
	for _, r := range raws {
		if r.hasID {
			entries = append(entries, Entry{
				ID: r.id, Content: r.content,
			})
		} else {
			entries = append(entries, Entry{
				ID: nextID, Content: r.content,
			})
			nextID++
		}
	}

	return entries
}

// FormatEntriesWithIDs serializes entries with ID prefixes.
//
// Parameters:
//   - entries: Entries to serialize
//
// Returns:
//   - []byte: Serialized content with ID prefixes
func FormatEntriesWithIDs(entries []Entry) []byte {
	if len(entries) == 0 {
		return nil
	}
	lines := make([]string, len(entries))
	for i, e := range entries {
		lines[i] = fmt.Sprintf(pad.FmtPadEntryID, e.ID, e.Content)
	}
	return []byte(
		strings.Join(lines, token.NewlineLF) +
			token.NewlineLF)
}

// NextID returns the next available ID for a new entry.
//
// Parameters:
//   - entries: Current entries
//
// Returns:
//   - int: Next ID (max existing + 1, minimum 1)
func NextID(entries []Entry) int {
	maxID := 0
	for _, e := range entries {
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	return maxID + 1
}

// FindByID returns the index of the entry with the given ID.
// Returns -1 if not found.
//
// Parameters:
//   - entries: Entries to search
//   - id: Stable ID to find
//
// Returns:
//   - int: Index in the slice, or -1 if not found
func FindByID(entries []Entry, id int) int {
	for i, e := range entries {
		if e.ID == id {
			return i
		}
	}
	return -1
}

// Normalize reassigns IDs as 1..N in current order.
//
// Parameters:
//   - entries: Entries to normalize
//
// Returns:
//   - []Entry: Entries with sequential IDs
func Normalize(entries []Entry) []Entry {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		result[i] = Entry{ID: i + 1, Content: e.Content}
	}
	return result
}

// IDs expands string args into a flat list of integer IDs.
// Supports individual IDs ("3") and ranges ("3-5").
//
// Parameters:
//   - args: command arguments to parse
//
// Returns:
//   - []int: expanded IDs
//   - error: non-nil if any arg is invalid
func IDs(args []string) ([]int, error) {
	var ids []int
	for _, arg := range args {
		if strings.Contains(arg, token.Dash) {
			parts := strings.SplitN(arg, token.Dash, 2)
			start, startErr := strconv.Atoi(parts[0])
			if startErr != nil {
				return nil, startErr
			}
			end, endErr := strconv.Atoi(parts[1])
			if endErr != nil {
				return nil, endErr
			}
			for i := start; i <= end; i++ {
				ids = append(ids, i)
			}
		} else {
			n, nErr := strconv.Atoi(arg)
			if nErr != nil {
				return nil, nErr
			}
			ids = append(ids, n)
		}
	}
	return ids, nil
}

// ToStrings extracts content strings from entries.
// Used for backward compatibility with existing code paths.
//
// Parameters:
//   - entries: Entries to extract from
//
// Returns:
//   - []string: Content strings without ID prefixes
func ToStrings(entries []Entry) []string {
	if entries == nil {
		return nil
	}
	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = e.Content
	}
	return result
}
