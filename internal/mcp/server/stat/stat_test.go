//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stat

import "testing"

func TestTotalAddsEmpty(t *testing.T) {
	if got := TotalAdds(nil); got != 0 {
		t.Errorf("TotalAdds(nil) = %d, want 0", got)
	}
}

func TestTotalAddsMultiple(t *testing.T) {
	m := map[string]int{"decision": 2, "learning": 3, "convention": 1}
	if got := TotalAdds(m); got != 6 {
		t.Errorf("TotalAdds = %d, want 6", got)
	}
}
