//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"path/filepath"

	"github.com/spf13/cobra"

	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/sync"
)

// ContextAndRoot resolves the context directory and its parent
// (project root) for a memory subcommand Run.
//
// Silences cobra's usage dump on error: a missing CTX_DIR is a
// declaration problem, not a misuse of the command. Callers return
// the error unchanged so the standard tailored message from
// rc.RequireContextDir reaches the user.
//
// Parameters:
//   - cmd: the cobra command being run (used only for SilenceUsage).
//
// Returns:
//   - string: absolute path to the declared context directory.
//   - string: project root (filepath.Dir of the context directory),
//     where MEMORY.md is expected to live.
//   - error: non-nil when the context directory is not declared.
func ContextAndRoot(cmd *cobra.Command) (string, string, error) {
	contextDir, err := rc.RequireContextDir()
	if err != nil {
		cmd.SilenceUsage = true
		return "", "", err
	}
	return contextDir, filepath.Dir(contextDir), nil
}

// DiscoverSource runs memory.DiscoverPath and applies the standard
// "auto memory not active" treatment: surface the helper notice to
// the Cobra command's output and return errMemory.NotFound. This is
// the shape four of the six memory subcommands (importer, publish,
// sync, unpublish) share. The diff and status commands want a
// different discovery-failure message and keep their handling
// inline.
//
// Parameters:
//   - cmd: the cobra command being run (passed through to the
//     sync.ErrAutoMemoryNotActive helper for user-facing output).
//   - projectRoot: project root previously resolved via
//     [ContextAndRoot].
//
// Returns:
//   - string: absolute path to the MEMORY.md source file when
//     discovered successfully.
//   - error: errMemory.NotFound when DiscoverPath fails; nil on
//     success.
func DiscoverSource(cmd *cobra.Command, projectRoot string) (string, error) {
	sourcePath, err := memory.DiscoverPath(projectRoot)
	if err != nil {
		sync.ErrAutoMemoryNotActive(cmd, err)
		return "", errMemory.NotFound()
	}
	return sourcePath, nil
}

// ReadSource reads the MEMORY.md file at the given path, splitting
// it into the directory + base filename that io.SafeReadFile wants.
// The helper wraps read failures in errMemory.Read so callers get a
// consistent user-facing error message.
//
// Parameters:
//   - path: absolute path to the MEMORY.md source file.
//
// Returns:
//   - []byte: file contents on success.
//   - error: errMemory.Read wrapping the underlying io error.
func ReadSource(path string) ([]byte, error) {
	data, err := io.SafeReadFile(filepath.Dir(path), filepath.Base(path))
	if err != nil {
		return nil, errMemory.Read(err)
	}
	return data, nil
}
