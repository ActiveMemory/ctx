//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"crypto/rand"
	"encoding/hex"

	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
)

const (
	// tokenBytes is the number of random bytes in a token.
	tokenBytes = 32

	// adminTokenPrefix distinguishes admin tokens.
	adminTokenPrefix = "ctx_adm_" //nolint:gosec // prefix, not a credential

	// clientTokenPrefix distinguishes client tokens.
	clientTokenPrefix = "ctx_cli_" //nolint:gosec // prefix, not a credential
)

// generateToken creates a prefixed hex-encoded random token.
func generateToken(prefix string) (string, error) {
	b := make([]byte, tokenBytes)
	if _, randErr := rand.Read(b); randErr != nil {
		return "", errHub.GenerateToken(randErr)
	}
	return prefix + hex.EncodeToString(b), nil
}
