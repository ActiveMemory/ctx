//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestContextDirLine_Default(t *testing.T) {
	rc.Reset()
	got := contextDirLine()
	want := "Context: .context"
	if got != want {
		t.Errorf("contextDirLine() = %q, want %q", got, want)
	}
}

func TestContextDirLine_Override(t *testing.T) {
	rc.Reset()
	rc.OverrideContextDir("/mnt/nas/.context")
	defer rc.Reset()

	got := contextDirLine()
	want := "Context: /mnt/nas/.context"
	if got != want {
		t.Errorf("contextDirLine() = %q, want %q", got, want)
	}
}
