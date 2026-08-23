//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pi

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deployExtension writes the embedded extension to
// .pi/extensions/ctx.ts. Pi auto-loads flat top-level files under
// .pi/extensions/ (project-local; loaded only after project trust).
// If an installed ctx extension differs from the embedded one, it
// is refreshed in place.
//
// The extension is a static embedded asset: a thin TypeScript shim
// that wires Pi lifecycle events to ctx system subcommands, with a
// type-only import of @earendil-works/pi-coding-agent erased at
// compile time. No runtime npm dependency tree is needed.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func deployExtension(cmd *cobra.Command) error {
	extensionDir := filepath.Join(
		cfgHook.DirPi, cfgHook.DirPiExtensions,
	)
	target := filepath.Join(
		extensionDir, cfgHook.FilePiExtensionDeploy,
	)
	if _, validateErr := validateManagedTarget(target); validateErr != nil {
		return validateErr
	}

	files, readErr := agent.PiExtension()
	if readErr != nil {
		return readErr
	}
	content, ok := files[cfgHook.FilePiExtensionDeploy]
	if !ok {
		return errSetup.MissingEmbeddedAsset(cfgHook.FilePiExtensionDeploy)
	}

	if existing, statErr := ctxIo.SafeReadUserFile(target); statErr == nil {
		if bytes.Equal(existing, content) {
			writeSetup.InfoPiSkipped(cmd, target)
			return nil
		}
	} else if !os.IsNotExist(statErr) {
		return errFs.FileRead(target, statErr)
	}

	if mkErr := ctxIo.SafeMkdirAll(
		extensionDir, fs.PermExec,
	); mkErr != nil {
		return errFs.Mkdir(extensionDir, mkErr)
	}

	if wErr := ctxIo.SafeWriteFile(
		target, content, fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(target, wErr)
	}
	writeSetup.InfoPiCreated(cmd, target)

	return nil
}
