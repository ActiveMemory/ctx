//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package http

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// UnsafeURLScheme returns an error when a URL uses a scheme other than
// http or https.
//
// Parameters:
//   - scheme: the rejected URL scheme
//
// Returns:
//   - error: "unsafe URL scheme: <scheme>: only http and https are allowed"
func UnsafeURLScheme(scheme string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrHttpUnsafeURLScheme), scheme,
	)
}

// ParseURL wraps a failure to parse a URL.
//
// Parameters:
//   - cause: the underlying parse error
//
// Returns:
//   - error: "parse URL: <cause>"
func ParseURL(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrHttpParseURL), cause,
	)
}

// TooManyRedirects returns an error when an HTTP response exceeds the
// redirect limit.
//
// Returns:
//   - error: "too many redirects: limit exceeded"
func TooManyRedirects() error {
	return errors.New(
		desc.TextDesc(text.DescKeyErrHttpTooManyRedirects),
	)
}
