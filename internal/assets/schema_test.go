//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
)

func TestSchema(t *testing.T) {
	data, err := FS.ReadFile(asset.PathCtxrcSchema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "$schema") {
		t.Error("does not contain $schema")
	}
	if !strings.Contains(content, "ctx.ist") {
		t.Error("does not contain ctx.ist $id")
	}
}

// The .ctxrc schema-vs-struct bijection guard lives in
// internal/rc/schema_test.go (TestSchemaCoversCtxRC): it reflects over
// the real rc.CtxRC so it cannot silently drift from a hand-maintained
// copy. It cannot live here — internal/assets is imported (via
// assets/read/placeholders) by internal/rc, so a package-assets test
// importing rc would form an import cycle.
