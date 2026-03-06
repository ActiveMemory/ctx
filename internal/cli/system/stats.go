//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// statsCmd returns the "ctx system stats" command.
//
// Streams or dumps per-session token usage stats from
// .context/state/stats-*.jsonl files.
func statsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show session token usage stats",
		Long: `Display per-session token usage statistics from stats JSONL files.

By default, shows the last 20 entries across all sessions. Use --follow
to stream new entries as they arrive (like tail -f).

Flags:
  --follow, -f   Stream new entries as they arrive
  --session, -s  Filter by session ID (prefix match)
  --last, -n     Show last N entries (default 20)
  --json, -j     Output raw JSONL`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStats(cmd)
		},
	}

	cmd.Flags().BoolP("follow", "f", false, "Stream new entries as they arrive")
	cmd.Flags().StringP("session", "s", "", "Filter by session ID (prefix match)")
	cmd.Flags().IntP("last", "n", 20, "Show last N entries")
	cmd.Flags().BoolP("json", "j", false, "Output raw JSONL")

	return cmd
}

func runStats(cmd *cobra.Command) error {
	follow, _ := cmd.Flags().GetBool("follow")
	session, _ := cmd.Flags().GetString("session")
	last, _ := cmd.Flags().GetInt("last")
	jsonOut, _ := cmd.Flags().GetBool("json")

	dir := filepath.Join(rc.ContextDir(), config.DirState)

	entries, readErr := readStatsDir(dir, session)
	if readErr != nil {
		return readErr
	}

	if !follow {
		return dumpStats(cmd, entries, last, jsonOut)
	}

	// Dump existing entries first, then stream.
	if dumpErr := dumpStats(cmd, entries, last, jsonOut); dumpErr != nil {
		return dumpErr
	}

	return streamStats(cmd, dir, session, jsonOut)
}

// statsEntry is a sessionStats with the source file for display.
type statsEntry struct {
	sessionStats
	Session string `json:"session"`
}

// readStatsDir reads all stats JSONL files, optionally filtered by session prefix.
func readStatsDir(dir, sessionFilter string) ([]statsEntry, error) {
	pattern := filepath.Join(dir, "stats-*.jsonl")
	matches, globErr := filepath.Glob(pattern)
	if globErr != nil {
		return nil, fmt.Errorf("globbing stats files: %w", globErr)
	}

	var entries []statsEntry
	for _, path := range matches {
		sid := extractSessionID(filepath.Base(path))
		if sessionFilter != "" && !strings.HasPrefix(sid, sessionFilter) {
			continue
		}
		fileEntries, parseErr := parseStatsFile(path, sid)
		if parseErr != nil {
			continue
		}
		entries = append(entries, fileEntries...)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp < entries[j].Timestamp
	})

	return entries, nil
}

// extractSessionID gets the session ID from a filename like "stats-abc123.jsonl".
func extractSessionID(basename string) string {
	s := strings.TrimPrefix(basename, "stats-")
	return strings.TrimSuffix(s, ".jsonl")
}

// parseStatsFile reads all JSONL lines from a stats file.
func parseStatsFile(path, sid string) ([]statsEntry, error) {
	data, readErr := os.ReadFile(path) //nolint:gosec // project-local state path
	if readErr != nil {
		return nil, readErr
	}

	var entries []statsEntry
	for _, line := range strings.Split(strings.TrimSpace(string(data)), config.NewlineLF) {
		if line == "" {
			continue
		}
		var s sessionStats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, statsEntry{sessionStats: s, Session: sid})
	}
	return entries, nil
}

// dumpStats outputs the last N entries.
func dumpStats(cmd *cobra.Command, entries []statsEntry, last int, jsonOut bool) error {
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

	printStatsHeader(cmd)
	for i := range entries {
		printStatsLine(cmd, &entries[i])
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
			sid := extractSessionID(filepath.Base(path))
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

			newEntries := readNewLines(path, prev, sid)
			for i := range newEntries {
				if jsonOut {
					line, marshalErr := json.Marshal(newEntries[i])
					if marshalErr == nil {
						cmd.Println(string(line))
					}
				} else {
					printStatsLine(cmd, &newEntries[i])
				}
			}
			offsets[path] = info.Size()
		}
	}

	return nil
}

// readNewLines reads bytes from offset to end and parses JSONL lines.
func readNewLines(path string, offset int64, sid string) []statsEntry {
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

	var entries []statsEntry
	for _, line := range strings.Split(strings.TrimSpace(string(buf[:n])), config.NewlineLF) {
		if line == "" {
			continue
		}
		var s sessionStats
		if jsonErr := json.Unmarshal([]byte(line), &s); jsonErr != nil {
			continue
		}
		entries = append(entries, statsEntry{sessionStats: s, Session: sid})
	}
	return entries
}

// outputStatsJSON writes entries as raw JSONL.
func outputStatsJSON(cmd *cobra.Command, entries []statsEntry) error {
	for _, e := range entries {
		line, marshalErr := json.Marshal(e)
		if marshalErr != nil {
			continue
		}
		cmd.Println(string(line))
	}
	return nil
}

// printStatsHeader prints the column header for human output.
func printStatsHeader(cmd *cobra.Command) {
	cmd.Println(fmt.Sprintf("%-19s  %-8s  %6s  %8s  %4s  %-12s",
		"TIME", "SESSION", "PROMPT", "TOKENS", "PCT", "EVENT"))
	cmd.Println(fmt.Sprintf("%-19s  %-8s  %6s  %8s  %4s  %-12s",
		"-------------------", "--------", "------", "--------", "----", "------------"))
}

// printStatsLine prints a single stats entry in human-readable format.
func printStatsLine(cmd *cobra.Command, e *statsEntry) {
	ts := formatStatsTimestamp(e.Timestamp)
	sid := e.Session
	if len(sid) > 8 {
		sid = sid[:8]
	}
	tokens := formatTokenCount(e.Tokens)
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
