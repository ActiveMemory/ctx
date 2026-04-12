//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package render

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/io"
)

// typedFileName returns the markdown filename for a type.
//
// Parameters:
//   - entryType: Entry type to build filename for
//
// Returns:
//   - string: Filename in the form "<type>s.md"
func typedFileName(entryType string) string {
	return entryType + cfgHub.SuffixPluralMD
}

// groupByType groups entries by their Type field.
//
// Parameters:
//   - entries: Slice of hub entries to group
//
// Returns:
//   - map[string][]hub.EntryMsg: Entries keyed by type
func groupByType(
	entries []hub.EntryMsg,
) map[string][]hub.EntryMsg {
	result := make(map[string][]hub.EntryMsg)
	for i := range entries {
		t := entries[i].Type
		result[t] = append(result[t], entries[i])
	}
	return result
}

// toMarkdown renders a slice of entries as markdown.
//
// Parameters:
//   - entries: Slice of hub entries to render
//
// Returns:
//   - string: Concatenated markdown for all entries
func toMarkdown(entries []hub.EntryMsg) string {
	var b strings.Builder
	for i := range entries {
		writeEntry(&b, &entries[i])
	}
	return b.String()
}

// writeEntry renders a single entry as markdown with
// origin tag and date header.
//
// Parameters:
//   - b: Builder to append markdown to
//   - e: Entry to render
func writeEntry(b *strings.Builder, e *hub.EntryMsg) {
	ts := time.Unix(e.Timestamp, 0).UTC()
	date := ts.Format(cfgTime.DateFormat)
	if _, err := fmt.Fprintf(b,
		tpl.TplEntryMarkdown,
		date, firstLine(e.Content),
		e.Origin, e.Content,
	); err != nil {
		return
	}
}

// firstLine returns the first line of s for use as a title.
//
// Parameters:
//   - s: Multi-line string to extract from
//
// Returns:
//   - string: First line, or the full string if single-line
func firstLine(s string) string {
	if line, _, ok := strings.Cut(s, token.NewlineLF); ok {
		return line
	}
	return s
}

// appendShared appends content to a file, creating it if
// needed.
//
// Parameters:
//   - path: Target file path
//   - content: Content to append
//
// Returns:
//   - error: Non-nil on I/O failure
func appendShared(path, content string) error {
	existing, readErr := io.SafeReadUserFile(path)
	if readErr != nil && !os.IsNotExist(readErr) {
		return readErr
	}
	return io.SafeWriteFile(
		path,
		append(existing, []byte(content)...),
		fs.PermFile,
	)
}
