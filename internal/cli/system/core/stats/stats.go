//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errRecall "github.com/ActiveMemory/ctx/internal/err/recall"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// ReadDir reads all stats JSONL files, optionally filtered by session prefix.
//
// Parameters:
//   - dir: path to the state directory
//   - sessionFilter: session ID prefix to filter by (empty for all)
//
// Returns:
//   - []Entry: sorted stats entries
//   - error: non-nil on glob failure
func ReadDir(dir, sessionFilter string) ([]Entry, error) {
	pattern := filepath.Join(dir, stats.FilePrefix+"*"+file.ExtJSONL)
	matches, globErr := filepath.Glob(pattern)
	if globErr != nil {
		return nil, errRecall.StatsGlob(globErr)
	}

	var entries []Entry
	for _, path := range matches {
		sid := ExtractSessionID(filepath.Base(path))
		if sessionFilter != "" && !strings.HasPrefix(sid, sessionFilter) {
			continue
		}
		fileEntries, parseErr := ParseFile(path, sid)
		if parseErr != nil {
			continue
		}
		entries = append(entries, fileEntries...)
	}

	sort.Slice(entries, func(i, j int) bool {
		ti, ei := time.Parse(time.RFC3339, entries[i].Timestamp)
		tj, ej := time.Parse(time.RFC3339, entries[j].Timestamp)
		if ei != nil || ej != nil {
			return entries[i].Timestamp < entries[j].Timestamp
		}
		return ti.Before(tj)
	})

	return entries, nil
}

// ExtractSessionID gets the session ID from a filename like
// "stats-abc123.jsonl".
//
// Parameters:
//   - basename: file basename
//
// Returns:
//   - string: session ID
func ExtractSessionID(basename string) string {
	s := strings.TrimPrefix(basename, stats.FilePrefix)
	return strings.TrimSuffix(s, file.ExtJSONL)
}

// ParseFile reads all JSONL lines from a stats file.
//
// Parameters:
//   - path: absolute path to the stats file
//   - sid: session ID for this file
//
// Returns:
//   - []Entry: parsed entries
//   - error: non-nil on read failure
func ParseFile(path, sid string) ([]Entry, error) {
	data, readErr := internalIo.SafeReadUserFile(path)
	if readErr != nil {
		return nil, readErr
	}

	var entries []Entry
	for _, line := range strings.Split(strings.TrimSpace(string(data)), token.NewlineLF) {
		if line == "" {
			continue
		}
		var s entity.Stats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, Entry{Stats: s, Session: sid})
	}
	return entries, nil
}

// FormatDump formats the last N entries in either JSON or human-readable format.
//
// Parameters:
//   - entries: stats entries to display
//   - last: number of entries to show (0 for all)
//   - jsonOut: whether to output as JSONL
//
// Returns:
//   - []string: formatted output lines
func FormatDump(entries []Entry, last int, jsonOut bool) []string {
	if len(entries) == 0 {
		return []string{desc.Text(text.DescKeyStatsEmpty)}
	}

	// Tail: take last N entries.
	if last > 0 && len(entries) > last {
		entries = entries[len(entries)-last:]
	}

	if jsonOut {
		return FormatJSON(entries)
	}

	h1, h2 := FormatHeader()
	lines := []string{h1, h2}
	for i := range entries {
		lines = append(lines, FormatLine(&entries[i]))
	}
	return lines
}

// FormatJSON formats entries as raw JSONL lines.
//
// Parameters:
//   - entries: stats entries to serialize
//
// Returns:
//   - []string: JSON lines (marshal errors are silently skipped)
func FormatJSON(entries []Entry) []string {
	var lines []string
	for _, e := range entries {
		line, marshalErr := json.Marshal(e)
		if marshalErr != nil {
			continue
		}
		lines = append(lines, string(line))
	}
	return lines
}

// FormatHeader formats the column header lines for human output.
//
// Returns:
//   - string: header line
//   - string: separator line
func FormatHeader() (string, string) {
	fmtStr := desc.Text(text.DescKeyStatsHeaderFormat)
	header := fmt.Sprintf(fmtStr,
		stats.HeaderTime, stats.HeaderSession,
		stats.HeaderPrompt, stats.HeaderTokens,
		stats.HeaderPct, stats.HeaderEvent)
	separator := fmt.Sprintf(fmtStr,
		stats.SepTime, stats.SepSession,
		stats.SepPrompt, stats.SepTokens,
		stats.SepPct, stats.SepEvent)
	return header, separator
}

// FormatLine formats a single stats entry in human-readable format.
//
// Parameters:
//   - e: stats entry to format
//
// Returns:
//   - string: formatted stats line
func FormatLine(e *Entry) string {
	ts := FormatTimestamp(e.Timestamp)
	sid := e.Session
	if len(sid) > journal.SessionIDShortLen {
		sid = sid[:journal.SessionIDShortLen]
	}
	tokens := session.FormatTokenCount(e.Tokens)
	return fmt.Sprintf(desc.Text(text.DescKeyStatsLineFormat),
		ts, sid, e.Prompt, tokens, e.Pct, e.Event)
}

// FormatTimestamp converts an RFC3339 timestamp to local time display
// using the DateTimePreciseFmt layout.
//
// Parameters:
//   - ts: RFC3339-formatted timestamp string
//
// Returns:
//   - string: local time formatted as "2006-01-02 15:04:05", or the
//     original string on parse failure
func FormatTimestamp(ts string) string {
	t, parseErr := time.Parse(time.RFC3339, ts)
	if parseErr != nil {
		return ts
	}
	return t.Local().Format(cfgTime.DateTimePreciseFmt)
}

// ReadNewLines reads bytes from offset to end and parses JSONL lines.
//
// Parameters:
//   - path: absolute path to the stats file
//   - offset: byte offset to start reading from
//   - sid: session ID for this file
//
// Returns:
//   - []Entry: newly parsed entries
func ReadNewLines(path string, offset int64, sid string) []Entry {
	f, openErr := internalIo.SafeOpenUserFile(path)
	if openErr != nil {
		return nil
	}
	defer func() { _ = f.Close() }()

	if _, seekErr := f.Seek(offset, 0); seekErr != nil {
		return nil
	}

	buf := make([]byte, stats.ReadBufSize)
	n, readErr := f.Read(buf)
	if readErr != nil || n == 0 {
		return nil
	}

	var entries []Entry
	for _, line := range strings.Split(strings.TrimSpace(string(buf[:n])), token.NewlineLF) {
		if line == "" {
			continue
		}
		var s entity.Stats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, Entry{Stats: s, Session: sid})
	}
	return entries
}

// Stream polls for new JSONL lines and writes them as they arrive.
//
// Parameters:
//   - w: output writer
//   - dir: path to the state directory
//   - sessionFilter: session ID prefix to filter by (empty for all)
//   - jsonOut: whether to output as JSONL
//
// Returns:
//   - error: Always nil
func Stream(w io.Writer, dir, sessionFilter string, jsonOut bool) error {
	// Track file sizes to detect new content.
	offsets := make(map[string]int64)
	matches, _ := filepath.Glob(filepath.Join(dir, stats.FilePrefix+"*"+file.ExtJSONL))
	for _, path := range matches {
		info, statErr := os.Stat(path)
		if statErr == nil {
			offsets[path] = info.Size()
		}
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		matches, _ = filepath.Glob(filepath.Join(dir, stats.FilePrefix+"*"+file.ExtJSONL))
		for _, path := range matches {
			sid := ExtractSessionID(filepath.Base(path))
			if sessionFilter != "" && !strings.HasPrefix(sid, sessionFilter) {
				continue
			}

			info, statErr := os.Stat(path)
			if statErr != nil {
				continue
			}
			prev := offsets[path]
			if info.Size() <= prev {
				continue
			}

			newEntries := ReadNewLines(path, prev, sid)
			for i := range newEntries {
				if jsonOut {
					line, marshalErr := json.Marshal(newEntries[i])
					if marshalErr == nil {
						_, _ = fmt.Fprintln(w, string(line))
					}
				} else {
					_, _ = fmt.Fprintln(w, FormatLine(&newEntries[i]))
				}
			}
			offsets[path] = info.Size()
		}
	}

	return nil
}
