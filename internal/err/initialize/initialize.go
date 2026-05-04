//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgInit "github.com/ActiveMemory/ctx/internal/config/initialize"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ErrContextPopulated is returned by ctx init when the target
// .context/ directory contains a populated context (essential
// files present) and the operator did not pass --reset.
//
// Wrap with [Populated] to attach the file list and the
// recovery hint pointing at --reset.
//
// The message lives in config/initialize (not resolved through
// desc.Text) because the sentinel is initialized at package
// load time, before the embedded YAML lookup is populated.
var ErrContextPopulated = errors.New(cfgInit.ErrMsgContextPopulated)

// ErrResetRequiresInteractive is returned by ctx init --reset
// when the invocation is non-interactive (the --caller flag
// is set, indicating an editor or scripted entry point). Reset
// is destructive and must come from a real terminal session.
var ErrResetRequiresInteractive = errors.New(
	cfgInit.ErrMsgResetRequiresInteractive,
)

// Populated wraps [ErrContextPopulated] with the populated
// files (basenames) and a hint about how to proceed.
//
// Parameters:
//   - dir:   absolute path of the .context/ directory
//   - files: basenames of populated essential files (already
//     filtered to those that exist)
//
// Returns:
//   - error: wrapping ErrContextPopulated for errors.Is matches
func Populated(dir string, files []string) error {
	listing := strings.Join(files, token.CommaSpace)
	return fmt.Errorf(
		desc.Text(text.DescKeyErrInitContextPopulated),
		ErrContextPopulated,
		dir,
		listing,
		cfgInit.ResetFlag,
		filepath.Join(dir, cfgInit.BackupPlaceholder),
	)
}

// ResetRequiresInteractive wraps [ErrResetRequiresInteractive]
// with a hint about how to invoke init from a real terminal.
//
// Returns:
//   - error: wrapping ErrResetRequiresInteractive for errors.Is matches
func ResetRequiresInteractive() error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrInitResetRequiresInteractive),
		ErrResetRequiresInteractive,
	)
}

// BackupMkdir wraps a failure to create the backup directory.
//
// Parameters:
//   - path:  absolute path of the backup directory that failed
//   - cause: the underlying mkdir error
//
// Returns:
//   - error: "create backup dir <path>: <cause>"
func BackupMkdir(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrInitBackupMkdir), path, cause,
	)
}

// BackupRead wraps a failure to read a populated source file
// while preparing the snapshot.
//
// Parameters:
//   - path:  absolute path of the source file
//   - cause: the underlying read error
//
// Returns:
//   - error: "read <path> for backup: <cause>"
func BackupRead(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrInitBackupRead), path, cause,
	)
}

// BackupWrite wraps a failure to write the snapshot copy.
//
// Parameters:
//   - path:  absolute path of the backup target file
//   - cause: the underlying write error
//
// Returns:
//   - error: "write backup <path>: <cause>"
func BackupWrite(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrInitBackupWrite), path, cause,
	)
}
