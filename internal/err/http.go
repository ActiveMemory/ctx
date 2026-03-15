//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
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
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrHttpUnsafeURLScheme), scheme)
}

// TooManyRedirects returns an error when an HTTP response exceeds the
// redirect limit.
//
// Returns:
//   - error: "too many redirects: limit exceeded"
func TooManyRedirects() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrHttpTooManyRedirects))
}
