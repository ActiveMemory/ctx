//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashRender returns a stable hex digest of a rendered journal-entry
// body, used to prove a file is still exactly what ctx last wrote.
//
// Callers hash the entry BODY (frontmatter stripped), so agent-side
// enrichment — which only touches frontmatter — leaves the digest
// unchanged, while a hand edit to the transcript body changes it. A
// mismatch on a later growth sweep means the file was edited outside
// ctx and must not be clobbered.
//
// Parameters:
//   - body: the rendered entry content with any YAML frontmatter removed
//
// Returns:
//   - string: lowercase hex SHA-256 of the body
func HashRender(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}
