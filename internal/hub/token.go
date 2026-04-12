//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"crypto/rand"
	"encoding/hex"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
)

// generateToken creates a prefixed hex-encoded random
// token.
//
// Parameters:
//   - prefix: token prefix (admin or client)
//
// Returns:
//   - string: prefixed hex-encoded token
//   - error: non-nil if random generation fails
func generateToken(prefix string) (string, error) {
	b := make([]byte, cfgHub.TokenBytes)
	if _, randErr := rand.Read(b); randErr != nil {
		return "", errHub.GenerateToken(randErr)
	}
	return prefix + hex.EncodeToString(b), nil
}
