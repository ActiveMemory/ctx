//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import "testing"

func TestContextFile(t *testing.T) {
	ctx := &Context{
		Files: []FileInfo{
			{Name: "TASKS.md", Path: "/tmp/TASKS.md"},
			{Name: "DECISIONS.md", Path: "/tmp/DECISIONS.md"},
		},
	}

	t.Run("found", func(t *testing.T) {
		f := ctx.File("TASKS.md")
		if f == nil {
			t.Fatal("expected non-nil FileInfo")
		}
		if f.Name != "TASKS.md" {
			t.Fatalf("got Name=%q, want TASKS.md", f.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		f := ctx.File("NOPE.md")
		if f != nil {
			t.Fatalf("expected nil, got %+v", f)
		}
	})

	t.Run("empty files", func(t *testing.T) {
		empty := &Context{}
		f := empty.File("TASKS.md")
		if f != nil {
			t.Fatalf("expected nil, got %+v", f)
		}
	})
}
