//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dep

import (
	"os/exec"

	cfgDep "github.com/ActiveMemory/ctx/internal/config/dep"
)

// GoListPackages runs `go list -json ./...` and returns the raw
// output. The caller is responsible for JSON decoding (go list
// emits concatenated JSON objects, not an array).
//
// Returns:
//   - []byte: raw go list stdout
//   - error: non-nil if go is not found or the command fails
func GoListPackages() ([]byte, error) {
	//nolint:gosec // G204: all args are package constants
	return exec.Command(
		cfgDep.GoBinary, cfgDep.GoList,
		cfgDep.GoFlagJSON, cfgDep.GoAllPackages,
	).Output()
}
