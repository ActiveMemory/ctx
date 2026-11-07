//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"bytes"
	"encoding/json"
)

// encodeJSON marshals a value with the given indent without HTML
// escaping and with a trailing newline.
//
// Parameters:
//   - v: value to encode
//   - indent: indentation prefix
//
// Returns:
//   - []byte: encoded JSON bytes
//   - error: non-nil if encoding fails
func encodeJSON(v any, indent string) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	if encodeErr := encoder.Encode(v); encodeErr != nil {
		return nil, encodeErr
	}
	return buf.Bytes(), nil
}

// equivalent reports whether two JSON documents carry equivalent data
// ignoring formatting and key order.
//
// Parameters:
//   - a: first JSON byte slice
//   - b: second JSON byte slice
//
// Returns:
//   - bool: true when both parse to identical canonical structures
//   - error: non-nil if either slice fails to parse
func equivalent(a, b []byte) (bool, error) {
	var av, bv any
	if aErr := json.Unmarshal(a, &av); aErr != nil {
		return false, aErr
	}
	if bErr := json.Unmarshal(b, &bv); bErr != nil {
		return false, bErr
	}
	ac, acErr := json.Marshal(av)
	if acErr != nil {
		return false, acErr
	}
	bc, bcErr := json.Marshal(bv)
	if bcErr != nil {
		return false, bcErr
	}
	return bytes.Equal(ac, bc), nil
}
