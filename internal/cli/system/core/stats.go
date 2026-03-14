//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	time2 "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// ReadStatsDir reads all stats JSONL files, optionally filtered by session prefix.
//
// Parameters:
//   - dir: path to the state directory
//   - sessionFilter: session ID prefix to filter by (empty for all)
//
// Returns:
//   - []StatsEntry: sorted stats entries
//   - error: non-nil on glob failure
func ReadStatsDir(dir, sessionFilter string) ([]StatsEntry, error) {
	pattern := filepath.Join(dir, stats.FilePrefix+"*"+file.ExtJSONL)
	matches, globErr := filepath.Glob(pattern)
	if globErr != nil {
		return nil, ctxerr.StatsGlob(globErr)
	}

	var entries []StatsEntry
	for _, path := range matches {
		sid := ExtractStatsSessionID(filepath.Base(path))
		if sessionFilter != "" && !strings.HasPrefix(sid, sessionFilter) {
			continue
		}
		fileEntries, parseErr := ParseStatsFile(path, sid)
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

// ExtractStatsSessionID gets the session ID from a filename like
// "stats-abc123.jsonl".
//
// Parameters:
//   - basename: file basename
//
// Returns:
//   - string: session ID
func ExtractStatsSessionID(basename string) string {
	s := strings.TrimPrefix(basename, stats.FilePrefix)
	return strings.TrimSuffix(s, file.ExtJSONL)
}

// ParseStatsFile reads all JSONL lines from a stats file.
//
// Parameters:
//   - path: absolute path to the stats file
//   - sid: session ID for this file
//
// Returns:
//   - []StatsEntry: parsed entries
//   - error: non-nil on read failure
func ParseStatsFile(path, sid string) ([]StatsEntry, error) {
	data, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		return nil, readErr
	}

	var entries []StatsEntry
	for _, line := range strings.Split(strings.TrimSpace(string(data)), token.NewlineLF) {
		if line == "" {
			continue
		}
		var s SessionStats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, StatsEntry{SessionStats: s, Session: sid})
	}
	return entries, nil
}

// DumpStats outputs the last N entries in either JSON or human-readable format.
//
// Parameters:
//   - cmd: Cobra command for output
//   - entries: stats entries to display
//   - last: number of entries to show (0 for all)
//   - jsonOut: whether to output as JSONL
//
// Returns:
//   - error: non-nil on output failure
func DumpStats(cmd *cobra.Command, entries []StatsEntry, last int, jsonOut bool) error {
	if len(entries) == 0 {
		cmd.Println(assets.TextDesc(assets.TextDescKeyStatsEmpty))
		return nil
	}

	// Tail: take last N entries.
	if last > 0 && len(entries) > last {
		entries = entries[len(entries)-last:]
	}

	if jsonOut {
		return OutputStatsJSON(cmd, entries)
	}

	PrintStatsHeader(cmd)
	for i := range entries {
		PrintStatsLine(cmd, &entries[i])
	}
	return nil
}

// OutputStatsJSON writes entries as raw JSONL.
//
// Parameters:
//   - cmd: Cobra command for output
//   - entries: stats entries to serialize
//
// Returns:
//   - error: Always nil (marshal errors are silently skipped)
func OutputStatsJSON(cmd *cobra.Command, entries []StatsEntry) error {
	for _, e := range entries {
		line, marshalErr := json.Marshal(e)
		if marshalErr != nil {
			continue
		}
		cmd.Println(string(line))
	}
	return nil
}

// PrintStatsHeader prints the column header for human output.
//
// Parameters:
//   - cmd: Cobra command for output
func PrintStatsHeader(cmd *cobra.Command) {
	fmtStr := assets.TextDesc(assets.TextDescKeyStatsHeaderFormat)
	cmd.Println(fmt.Sprintf(fmtStr,
		stats.HeaderTime, stats.HeaderSession,
		stats.HeaderPrompt, stats.HeaderTokens,
		stats.HeaderPct, stats.HeaderEvent))
	cmd.Println(fmt.Sprintf(fmtStr,
		stats.SepTime, stats.SepSession,
		stats.SepPrompt, stats.SepTokens,
		stats.SepPct, stats.SepEvent))
}

// PrintStatsLine prints a single stats entry in human-readable format.
//
// Parameters:
//   - cmd: Cobra command for output
//   - e: stats entry to print
func PrintStatsLine(cmd *cobra.Command, e *StatsEntry) {
	ts := FormatStatsTimestamp(e.Timestamp)
	sid := e.Session
	if len(sid) > journal.SessionIDShortLen {
		sid = sid[:journal.SessionIDShortLen]
	}
	tokens := FormatTokenCount(e.Tokens)
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyStatsLineFormat),
		ts, sid, e.Prompt, tokens, e.Pct, e.Event))
}

// FormatStatsTimestamp converts an RFC3339 timestamp to local time display
// using the DateTimePreciseFormat layout.
//
// Parameters:
//   - ts: RFC3339-formatted timestamp string
//
// Returns:
//   - string: local time formatted as "2006-01-02 15:04:05", or the
//     original string on parse failure
func FormatStatsTimestamp(ts string) string {
	t, parseErr := time.Parse(time.RFC3339, ts)
	if parseErr != nil {
		return ts
	}
	return t.Local().Format(time2.DateTimePreciseFormat)
}

// ReadNewLines reads bytes from offset to end and parses JSONL lines.
//
// Parameters:
//   - path: absolute path to the stats file
//   - offset: byte offset to start reading from
//   - sid: session ID for this file
//
// Returns:
//   - []StatsEntry: newly parsed entries
func ReadNewLines(path string, offset int64, sid string) []StatsEntry {
	f, openErr := io.SafeOpenUserFile(path)
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

	var entries []StatsEntry
	for _, line := range strings.Split(strings.TrimSpace(string(buf[:n])), token.NewlineLF) {
		if line == "" {
			continue
		}
		var s SessionStats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, StatsEntry{SessionStats: s, Session: sid})
	}
	return entries
}

// StreamStats polls for new JSONL lines and prints them as they arrive.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: path to the state directory
//   - sessionFilter: session ID prefix to filter by (empty for all)
//   - jsonOut: whether to output as JSONL
//
// Returns:
//   - error: Always nil
func StreamStats(cmd *cobra.Command, dir, sessionFilter string, jsonOut bool) error {
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
			sid := ExtractStatsSessionID(filepath.Base(path))
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
						cmd.Println(string(line))
					}
				} else {
					PrintStatsLine(cmd, &newEntries[i])
				}
			}
			offsets[path] = info.Size()
		}
	}

	return nil
}
