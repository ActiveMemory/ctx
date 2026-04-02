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
)

// Validate performs strict YAML decoding of .ctxrc content.
//
// Unknown fields are returned as warnings (not errors) so callers can
// distinguish typos from genuinely broken YAML.
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
			return te.Errors, nil
		}

		// Genuinely broken YAML.
		return nil, decErr
	}

	return nil, nil
}
