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

// CargoMetadata runs `cargo metadata --format-version 1` and
// returns the raw output. When noDeps is true, --no-deps is
// appended to skip dependency resolution.
//
// Parameters:
//   - noDeps: if true, skip full dependency resolution
//
// Returns:
//   - []byte: raw cargo stdout
//   - error: non-nil if cargo is not found or the command fails
func CargoMetadata(noDeps bool) ([]byte, error) {
	if _, lookErr := exec.LookPath(cfgDep.CargoBinary); lookErr != nil {
		return nil, lookErr
	}
	args := []string{
		cfgDep.CargoMetadataCmd,
		cfgDep.CargoFlagFormatVersion, cfgDep.CargoFormatVersion1,
	}
	if noDeps {
		args = append(args, cfgDep.CargoFlagNoDeps)
	}
	//nolint:gosec // G204: args are literal constants
	return exec.Command(cfgDep.CargoBinary, args...).Output()
}
