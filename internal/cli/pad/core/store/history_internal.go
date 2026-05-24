//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package store

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	cfgPad "github.com/ActiveMemory/ctx/internal/config/pad"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
)

// mostRecentSnapshot returns the absolute path of the newest
// snapshot, or the empty string when no snapshots exist.
//
// Returns:
//   - string: absolute path of the newest snapshot file, or ""
//   - error: non-nil on history-list failure
func mostRecentSnapshot() (string, error) {
	dir, dirErr := HistoryDir()
	if dirErr != nil {
		return "", dirErr
	}

	entries, readErr := os.ReadDir(dir) //nolint:gosec // dir is rc-derived
	if errors.Is(readErr, os.ErrNotExist) {
		return "", nil
	}
	if readErr != nil {
		return "", errPad.HistoryRead(readErr)
	}

	var candidates []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !validSnapshotName(e.Name()) {
			continue
		}
		candidates = append(candidates, e.Name())
	}
	if len(candidates) == 0 {
		return "", nil
	}
	sort.Strings(candidates)
	return filepath.Join(dir, candidates[len(candidates)-1]), nil
}

// buildSnapshotName composes the filename for a new snapshot.
// Format: `<timestamp>-<op><ext>` where ext mirrors the live
// pad file (e.g. `.enc` or `.md`) so plaintext and encrypted
// modes co-exist cleanly under the same history directory.
//
// Parameters:
//   - cmd: cobra command whose Name() supplies the op label
//   - padPath: absolute path of the live pad (for ext detection)
//   - now: snapshot timestamp; pass [time.Now] in UTC at the
//     call site so tests can inject deterministic values
//
// Returns:
//   - string: snapshot filename (basename, not absolute path)
func buildSnapshotName(
	cmd *cobra.Command, padPath string, now time.Time,
) string {
	ts := now.Format(cfgPad.HistoryTimeFormat)
	op := opLabel(cmd)
	ext := filepath.Ext(padPath)
	return ts + cfgPad.HistorySnapshotSeparator + op + ext
}

// opLabel derives the snapshot op label from a cobra command,
// falling back to [cfgPad.HistoryOpUnknown] when no command
// context is available (nil cmd or empty Name).
//
// Parameters:
//   - cmd: subcommand-scoped cobra command, may be nil
//
// Returns:
//   - string: subcommand name (e.g. "rm", "edit") or the
//     unknown-op fallback
func opLabel(cmd *cobra.Command) string {
	if cmd == nil {
		return cfgPad.HistoryOpUnknown
	}
	name := cmd.Name()
	if name == "" {
		return cfgPad.HistoryOpUnknown
	}
	return name
}

// validSnapshotName checks that a filename starts with a
// parseable [cfgPad.HistoryTimeFormat] timestamp followed by
// the snapshot separator. Defends list / prune passes against
// partial files, leftover editor swap files, and any other
// foreign content in the history directory.
//
// Parameters:
//   - name: basename of a candidate file inside the history dir
//
// Returns:
//   - bool: true when the name matches the snapshot shape
func validSnapshotName(name string) bool {
	tsLen := len(cfgPad.HistoryTimeFormat)
	if len(name) < tsLen+len(cfgPad.HistorySnapshotSeparator) {
		return false
	}
	if !strings.HasPrefix(
		name[tsLen:], cfgPad.HistorySnapshotSeparator,
	) {
		return false
	}
	_, err := time.Parse(
		cfgPad.HistoryTimeFormat, name[:tsLen],
	)
	return err == nil
}

// slotFromName strips the extension from a snapshot filename
// to produce a user-displayable slot identifier (the part the
// `undo --to <slot>` flag will accept in Phase 2).
//
// Parameters:
//   - name: snapshot basename including extension
//
// Returns:
//   - string: the timestamp + op portion, no extension
func slotFromName(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}
	return strings.TrimSuffix(name, ext)
}
