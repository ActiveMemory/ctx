//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"testing"

	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// T06(a): verifyContains is the byte-presence guard — present spans pass,
// an absent span returns ErrVerifyFailed (white-box: the predicate is
// unexported and only reachable via a genuine IO race in Apply).
func TestVerifyContains(t *testing.T) {
	const file = "# hooks\n\n## [2026-01-01-000000] Alpha\n\nbody\n"
	if err := verifyContains(file, "## [2026-01-01-000000] Alpha\n\nbody\n"); err != nil {
		t.Errorf("present span: err = %v, want nil", err)
	}
	if err := verifyContains(file, "## [2026-09-09-090909] Ghost\n"); err != errDisc.ErrVerifyFailed {
		t.Errorf("absent span: err = %v, want ErrVerifyFailed", err)
	}
}
