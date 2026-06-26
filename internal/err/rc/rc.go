//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \\    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import "fmt"

import cfgRC "github.com/ActiveMemory/ctx/internal/config/rc"

// InvalidBackendEndpointScheme reports a backend endpoint with a non-http(s)
// scheme.
//
// Parameters:
//   - name: backend name under backends.*
//
// Returns:
//   - error: invalid endpoint scheme message for that backend
func InvalidBackendEndpointScheme(name string) error {
	return fmt.Errorf(cfgRC.ErrBackendsEndpointScheme, name)
}
