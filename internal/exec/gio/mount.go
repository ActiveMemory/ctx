//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gio

import (
	"os/exec"

	"github.com/ActiveMemory/ctx/internal/config/archive"
)

// Mount runs `gio mount` with the given URL.
//
// Parameters:
//   - url: mount target (e.g. smb://host/share)
//
// Returns:
//   - error: non-nil if gio is not found or the mount fails
func Mount(url string) error {
	//nolint:gosec // G204: url is from user config
	return exec.Command(
		archive.GioBinary, archive.GioMount, url,
	).Run()
}
