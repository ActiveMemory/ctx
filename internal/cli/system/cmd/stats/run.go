//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func runStats(cmd *cobra.Command) error {
	follow, _ := cmd.Flags().GetBool("follow")
	session, _ := cmd.Flags().GetString("session")
	last, _ := cmd.Flags().GetInt("last")
	jsonOut, _ := cmd.Flags().GetBool("json")

	dir := filepath.Join(rc.ContextDir(), config.DirState)

	entries, readErr := ReadStatsDir(dir, session)
	if readErr != nil {
		return readErr
	}

	if !follow {
		return DumpStats(cmd, entries, last, jsonOut)
	}

	// Dump existing entries first, then stream.
	if dumpErr := DumpStats(cmd, entries, last, jsonOut); dumpErr != nil {
		return dumpErr
	}

	return streamStats(cmd, dir, session, jsonOut)
}

// StatsEntry is a core.SessionStats with the source file for display.
type StatsEntry struct {
	core.SessionStats
	Session string `json:"session"`
}

// ReadStatsDir reads all stats JSONL files, optionally filtered by session prefix.
//
// Parameters:
//   - dir: Path to the state directory
//   - sessionFilter: Session ID prefix to filter by (empty for all)
//
// Returns:
//   - []StatsEntry: Sorted stats entries
//   - error: Non-nil on glob failure
func ReadStatsDir(dir, sessionFilter string) ([]StatsEntry, error) {
	pattern := filepath.Join(dir, "stats-*.jsonl")
	matches, globErr := filepath.Glob(pattern)
	if globErr != nil {
		return nil, fmt.Errorf("globbing stats files: %w", globErr)
	}

	var entries []StatsEntry
	for _, path := range matches {
		sid := ExtractSessionID(filepath.Base(path))
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

// ExtractSessionID gets the session ID from a filename like "stats-abc123.jsonl".
//
// Parameters:
//   - basename: File basename
//
// Returns:
//   - string: Session ID
func ExtractSessionID(basename string) string {
	s := strings.TrimPrefix(basename, "stats-")
	return strings.TrimSuffix(s, ".jsonl")
}

// ParseStatsFile reads all JSONL lines from a stats file.
//
// Parameters:
//   - path: Absolute path to the stats file
//   - sid: Session ID for this file
//
// Returns:
//   - []StatsEntry: Parsed entries
//   - error: Non-nil on read failure
func ParseStatsFile(path, sid string) ([]StatsEntry, error) {
	data, readErr := os.ReadFile(path) //nolint:gosec // project-local state path
	if readErr != nil {
		return nil, readErr
	}

	var entries []StatsEntry
	for _, line := range strings.Split(strings.TrimSpace(string(data)), config.NewlineLF) {
		if line == "" {
			continue
		}
		var s core.SessionStats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, StatsEntry{SessionStats: s, Session: sid})
	}
	return entries, nil
}

// DumpStats outputs the last N entries.
//
// Parameters:
//   - cmd: Cobra command for output
//   - entries: Stats entries to display
//   - last: Number of entries to show (0 for all)
//   - jsonOut: Whether to output as JSONL
//
// Returns:
//   - error: Non-nil on output failure
func DumpStats(cmd *cobra.Command, entries []StatsEntry, last int, jsonOut bool) error {
	if len(entries) == 0 {
		cmd.Println("No stats recorded yet.")
		return nil
	}

	// Tail: take last N entries.
	if last > 0 && len(entries) > last {
		entries = entries[len(entries)-last:]
	}

	if jsonOut {
		return outputStatsJSON(cmd, entries)
	}

	PrintStatsHeader(cmd)
	for i := range entries {
		PrintStatsLine(cmd, &entries[i])
	}
	return nil
}

// streamStats polls for new JSONL lines and prints them as they arrive.
func streamStats(cmd *cobra.Command, dir, sessionFilter string, jsonOut bool) error {
	// Track file sizes to detect new content.
	offsets := make(map[string]int64)
	matches, _ := filepath.Glob(filepath.Join(dir, "stats-*.jsonl"))
	for _, path := range matches {
		info, statErr := os.Stat(path)
		if statErr == nil {
			offsets[path] = info.Size()
		}
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		matches, _ = filepath.Glob(filepath.Join(dir, "stats-*.jsonl"))
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

// ReadNewLines reads bytes from offset to end and parses JSONL lines.
//
// Parameters:
//   - path: Absolute path to the stats file
//   - offset: Byte offset to start reading from
//   - sid: Session ID for this file
//
// Returns:
//   - []StatsEntry: Newly parsed entries
func ReadNewLines(path string, offset int64, sid string) []StatsEntry {
	f, openErr := os.Open(path) //nolint:gosec // project-local state path
	if openErr != nil {
		return nil
	}
	defer func() { _ = f.Close() }()

	if _, seekErr := f.Seek(offset, 0); seekErr != nil {
		return nil
	}

	buf := make([]byte, 8192)
	n, readErr := f.Read(buf)
	if readErr != nil || n == 0 {
		return nil
	}

	var entries []StatsEntry
	for _, line := range strings.Split(strings.TrimSpace(string(buf[:n])), config.NewlineLF) {
		if line == "" {
			continue
		}
		var s core.SessionStats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, StatsEntry{SessionStats: s, Session: sid})
	}
	return entries
}

// outputStatsJSON writes entries as raw JSONL.
func outputStatsJSON(cmd *cobra.Command, entries []StatsEntry) error {
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
	cmd.Println(fmt.Sprintf("%-19s  %-8s  %6s  %8s  %4s  %-12s",
		"TIME", "SESSION", "PROMPT", "TOKENS", "PCT", "EVENT"))
	cmd.Println(fmt.Sprintf("%-19s  %-8s  %6s  %8s  %4s  %-12s",
		"-------------------", "--------", "------", "--------", "----", "------------"))
}

// PrintStatsLine prints a single stats entry in human-readable format.
//
// Parameters:
//   - cmd: Cobra command for output
//   - e: Stats entry to print
func PrintStatsLine(cmd *cobra.Command, e *StatsEntry) {
	ts := formatStatsTimestamp(e.Timestamp)
	sid := e.Session
	if len(sid) > 8 {
		sid = sid[:8]
	}
	tokens := core.FormatTokenCount(e.Tokens)
	cmd.Println(fmt.Sprintf("%-19s  %-8s  %6d  %7s  %3d%%  %-12s",
		ts, sid, e.Prompt, tokens, e.Pct, e.Event))
}

// formatStatsTimestamp converts an RFC3339 timestamp to local time display.
func formatStatsTimestamp(ts string) string {
	t, parseErr := time.Parse(time.RFC3339, ts)
	if parseErr != nil {
		return ts
	}
	return t.Local().Format("2006-01-02 15:04:05")
}
