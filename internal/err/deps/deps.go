//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// CargoNotFound returns an error when cargo is not in PATH.
//
// Returns:
//   - error: advises installing the Rust toolchain
func CargoNotFound() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrDepsCargoNotFound))
}

// CargoMetadataFailed wraps a cargo metadata command failure.
//
// Parameters:
//   - cause: the underlying command error
//
// Returns:
//   - error: "cargo metadata failed: <cause>"
func CargoMetadataFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrDepsCargoMetadataFailed), cause)
}

// ParseCargoMetadata wraps a cargo metadata parse failure.
//
// Parameters:
//   - cause: the underlying unmarshal error
//
// Returns:
//   - error: "parsing cargo metadata: <cause>"
func ParseCargoMetadata(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrDepsParseCargoMetadata), cause)
}
