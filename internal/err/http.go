//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// UnsafeURLScheme returns an error when a URL uses a scheme other than
// http or https.
//
// Parameters:
//   - scheme: the rejected URL scheme
//
// Returns:
//   - error: "unsafe URL scheme: <scheme>: only http and https are allowed"
func UnsafeURLScheme(scheme string) error {
	return fmt.Errorf("unsafe URL scheme: %s: only http and https are allowed", scheme)
}

// TooManyRedirects returns an error when an HTTP response exceeds the
// redirect limit.
//
// Returns:
//   - error: "too many redirects: limit exceeded"
func TooManyRedirects() error {
	return fmt.Errorf("too many redirects: limit exceeded")
}
