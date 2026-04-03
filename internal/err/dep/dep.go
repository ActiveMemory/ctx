//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dep

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// CargoMetadataFailed wraps a cargo metadata command failure.
//
// Parameters:
//   - cause: the underlying command error
//
// Returns:
//   - error: "cargo metadata failed: <cause>"
func CargoMetadataFailed(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrDepsCargoMetadataFailed), cause)
}

// ParseCargoMetadata wraps a cargo metadata parse failure.
//
// Parameters:
//   - cause: the underlying unmarshal error
//
// Returns:
//   - error: "parsing cargo metadata: <cause>"
func ParseCargoMetadata(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrDepsParseCargoMetadata), cause)
}
