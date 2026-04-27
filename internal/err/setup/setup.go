//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// CreateDir wraps a failure to create a setup directory.
//
// Parameters:
//   - dir: the directory path
//   - cause: the underlying OS error
//
// Returns:
//   - error: "create <dir>: <cause>"
func CreateDir(dir string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSetupCreateDir), dir, cause,
	)
}

// MarshalConfig wraps a failure to marshal MCP configuration JSON.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "marshal mcp config: <cause>"
func MarshalConfig(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSetupMarshalConfig), cause,
	)
}

// WriteFile wraps a failure to write a setup file.
//
// Parameters:
//   - path: the file path
//   - cause: the underlying OS error
//
// Returns:
//   - error: "write <path>: <cause>"
func WriteFile(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSetupFileWrite), path, cause,
	)
}

// SyncSteering wraps a failure during steering sync in setup.
//
// Parameters:
//   - cause: the underlying sync error
//
// Returns:
//   - error: "sync steering: <cause>"
func SyncSteering(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSetupSyncSteering), cause,
	)
}

// MissingEmbeddedAsset reports that an asset expected to be
// embedded in the binary is missing — typically a setup-time
// invariant violation rather than a user-facing failure.
//
// Parameters:
//   - name: the asset key that was looked up
//
// Returns:
//   - error: "embedded asset missing: <name>"
func MissingEmbeddedAsset(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSetupMissingEmbeddedAsset), name,
	)
}
