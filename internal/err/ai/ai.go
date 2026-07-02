//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \\    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ai

import (
	"errors"

	cfgAI "github.com/ActiveMemory/ctx/internal/config/ai"
)

// EmitRequired reports that ctx ai propose requires a non-empty --emit value.
//
// Returns:
//   - error: "emit is required"
func EmitRequired() error {
	return errors.New(cfgAI.ErrEmitRequired)
}

// InvalidArtifact reports that a proposal artifact is structurally invalid.
//
// Returns:
//   - error: "invalid proposal artifact"
func InvalidArtifact() error {
	return errors.New(cfgAI.ErrInvalidArtifact)
}

// InvalidArtifactResponse reports that a proposal response payload is invalid.
//
// Returns:
//   - error: "invalid proposal artifact response"
func InvalidArtifactResponse() error {
	return errors.New(cfgAI.ErrInvalidArtifactResponse)
}
