//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pass

import (
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	dreamPaths "github.com/ActiveMemory/ctx/internal/cli/dream/core/paths"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeDream "github.com/ActiveMemory/ctx/internal/write/dream"
)

// Run executes one dream pass. It resolves paths, ensures the
// notebook, computes the delta, gates on an empty delta, takes the
// lock, invokes the executor (fail-loud on a missing binary, leaving a
// failmark), validates the proposals written, persists state for the
// processed sources, and prints a digest.
//
// Parameters:
//   - cmd: cobra command for output and the executor's stream wiring
//   - opts: the resolved run parameters
//
// Returns:
//   - error: a resolution, lock, executor, or persistence failure;
//     nil on the empty-delta and lock-held exit-0 paths
func Run(cmd *cobra.Command, opts Opts) error {
	loc, locErr := dreamPaths.Resolve()
	if locErr != nil {
		return locErr
	}
	if mkErr := ctxIo.SafeMkdirAll(
		loc.Dreams, cfgFs.PermRestrictedDir,
	); mkErr != nil {
		return errDream.Mkdir(loc.Dreams, mkErr)
	}

	if !opts.Force && !dreamDue(loc.Dreams) {
		writeDream.Nothing(cmd)
		return nil
	}

	selected, scanErr := selectSources(loc, opts.Max)
	if scanErr != nil {
		return scanErr
	}
	if len(selected) == 0 {
		writeDream.Nothing(cmd)
		return nil
	}

	lockPath := filepath.Join(loc.Dreams, cfgDream.FileLock)
	acquired, lockErr := ctxIo.SafeTryLock(lockPath, cfgFs.PermSecret)
	if lockErr != nil {
		return errDream.LockAcquire(lockPath, lockErr)
	}
	if !acquired {
		writeDream.Locked(cmd)
		return nil
	}
	defer func() { _ = ctxIo.SafeUnlock(lockPath) }()

	runDir := filepath.Join(
		loc.Dreams, time.Now().UTC().Format(cfgDream.RunTimeLayout),
	)
	if mkErr := ctxIo.SafeMkdirAll(
		runDir, cfgFs.PermRestrictedDir,
	); mkErr != nil {
		return errDream.Mkdir(runDir, mkErr)
	}

	if execErr := invoke(cmd, loc, runDir, opts); execErr != nil {
		writeFailmark(cmd, loc.Dreams)
		return execErr
	}

	valid, validateErr := validateRun(runDir)
	if validateErr != nil {
		return validateErr
	}
	if saveErr := persist(loc, selected); saveErr != nil {
		return saveErr
	}
	writeDream.Digest(cmd, len(selected), valid)
	return nil
}
