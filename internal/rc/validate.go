//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/yaml.v3"

	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Validate performs strict YAML decoding of .ctxrc content.
//
// Unknown fields are returned as warnings (not errors) so callers can
// distinguish typos from genuinely broken YAML. Backend configuration is
// stricter: malformed `backends:` keys fail validation.
//
// Parameters:
//   - data: Raw YAML content from a .ctxrc file
//
// Returns:
//   - warnings: Human-readable messages for each unknown field
//   - err: Non-nil only for genuinely malformed YAML
func Validate(data []byte) (warnings []string, err error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)

	var cfg CtxRC
	if decErr := dec.Decode(&cfg); decErr != nil {
		// Empty document: not an error.
		if decErr == io.EOF {
			return nil, nil
		}

		// yaml.v3 returns *yaml.TypeError for unknown fields.
		if te, ok := errors.AsType[*yaml.TypeError](decErr); ok {
			if backendsShapeError(te.Errors) {
				return nil, decErr
			}
			return te.Errors, nil
		}

		// Genuinely broken YAML.
		return nil, decErr
	}

	if cfgErr := cfg.validateBackends(); cfgErr != nil {
		return nil, cfgErr
	}

	return nil, nil
}

// ValidateCurrent validates the cwd-anchored `.ctxrc` file, when one is
// available.
//
// Returns:
//   - error: validation failure for the current `.ctxrc`; nil when the file is
//     absent, unreadable, or valid.
func ValidateCurrent() error {
	rcPath, pathErr := ctxrcPath()
	if pathErr != nil {
		if errors.Is(pathErr, errCtx.ErrNoCtxHere) {
			return nil
		}
		return nil
	}
	data, readErr := ctxIo.SafeReadUserFile(rcPath)
	if readErr != nil {
		return nil
	}
	_, validateErr := Validate(data)
	return validateErr
}
