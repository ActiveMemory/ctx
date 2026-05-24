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
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgPad "github.com/ActiveMemory/ctx/internal/config/pad"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// HistoryDir returns the absolute path of the scratchpad
// snapshot directory, a sibling of the live pad file inside
// `.context/`.
//
// Returns:
//   - string: history directory path
//   - error: propagated from [rc.ContextDir]
func HistoryDir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, cfgPad.HistoryDirName), nil
}

// SnapshotBefore copies the current pad blob to the history
// directory under a timestamped filename. Intended to be called
// from [WriteEntriesWithIDs] just before the live pad file is
// overwritten.
//
// First-write case (pad file does not exist yet) is a no-op:
// nothing to preserve.
//
// Parameters:
//   - cmd: Cobra command for op-label derivation. Nil yields
//     [cfgPad.HistoryOpUnknown].
//
// Returns:
//   - error: non-nil on read, mkdir, or write failure
func SnapshotBefore(cmd *cobra.Command) error {
	padPath, pathErr := ScratchpadPath()
	if pathErr != nil {
		return pathErr
	}

	info, statErr := os.Stat(padPath)
	if errors.Is(statErr, os.ErrNotExist) {
		return nil
	}
	if statErr != nil {
		return errPad.HistoryWrite(statErr)
	}

	dir, dirErr := HistoryDir()
	if dirErr != nil {
		return dirErr
	}
	if mkErr := io.SafeMkdirAll(
		dir, fs.PermRestrictedDir,
	); mkErr != nil {
		return errPad.HistoryWrite(mkErr)
	}

	data, readErr := io.SafeReadFile(
		filepath.Dir(padPath), filepath.Base(padPath),
	)
	if readErr != nil {
		return errPad.HistoryWrite(readErr)
	}

	snapName := buildSnapshotName(cmd, padPath, time.Now().UTC())
	snapPath := filepath.Join(dir, snapName)
	if writeErr := io.SafeWriteFile(
		snapPath, data, info.Mode().Perm(),
	); writeErr != nil {
		return errPad.HistoryWrite(writeErr)
	}
	return nil
}

// Restore promotes the most recent snapshot back to the live
// pad. The current pad is itself snapshotted first so a
// subsequent `undo` yields a redo.
//
// Parameters:
//   - cmd: Cobra command for op-label derivation on the
//     pre-restore snapshot. Nil yields [cfgPad.HistoryOpUnknown].
//
// Returns:
//   - string: slot identifier (basename without extension) of
//     the snapshot that was restored, or "" when no history
//     exists.
//   - error: non-nil on history-list, snapshot-take, or
//     restore-copy failure.
func Restore(cmd *cobra.Command) (string, error) {
	snapPath, listErr := mostRecentSnapshot()
	if listErr != nil {
		return "", listErr
	}
	if snapPath == "" {
		return "", nil
	}

	if snapErr := SnapshotBefore(cmd); snapErr != nil {
		return "", snapErr
	}

	padPath, pathErr := ScratchpadPath()
	if pathErr != nil {
		return "", pathErr
	}

	data, readErr := io.SafeReadFile(
		filepath.Dir(snapPath), filepath.Base(snapPath),
	)
	if readErr != nil {
		return "", errPad.HistoryRestore(readErr)
	}

	if writeErr := io.SafeWriteFile(
		padPath, data, fs.PermFile,
	); writeErr != nil {
		return "", errPad.HistoryRestore(writeErr)
	}

	return slotFromName(filepath.Base(snapPath)), nil
}

// Prune evicts snapshots that exceed the retention caps. Both
// the count cap ([cfgPad.HistoryMaxSnapshots]) and the age cap
// ([cfgPad.HistoryMaxAge]) apply; either alone is sufficient
// to evict. A cap of zero disables that cap independently.
//
// Returns:
//   - error: nil on success or when the history directory
//     does not exist yet. Per-file remove failures are
//     warned via [logWarn] and otherwise tolerated.
func Prune() error {
	dir, dirErr := HistoryDir()
	if dirErr != nil {
		return dirErr
	}

	entries, readErr := os.ReadDir(dir) //nolint:gosec // dir is rc-derived
	if errors.Is(readErr, os.ErrNotExist) {
		return nil
	}
	if readErr != nil {
		return errPad.HistoryRead(readErr)
	}

	type snap struct {
		name  string
		mtime time.Time
	}
	var snaps []snap
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !validSnapshotName(e.Name()) {
			continue
		}
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		snaps = append(snaps, snap{e.Name(), info.ModTime()})
	}

	// Newest first by lexical (timestamp-prefixed) name.
	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].name > snaps[j].name
	})

	cutoff := time.Now().UTC().Add(-cfgPad.HistoryMaxAge)
	for i, s := range snaps {
		evict := false
		if cfgPad.HistoryMaxSnapshots > 0 &&
			i >= cfgPad.HistoryMaxSnapshots {
			evict = true
		}
		if cfgPad.HistoryMaxAge > 0 && s.mtime.Before(cutoff) {
			evict = true
		}
		if !evict {
			continue
		}
		target := filepath.Join(dir, s.name)
		if rmErr := os.Remove(target); rmErr != nil {
			logWarn.Warn(
				cfgWarn.PadHistoryPruneFile, s.name, rmErr,
			)
		}
	}
	return nil
}
