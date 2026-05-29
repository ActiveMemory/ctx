//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex_test

import (
	"strings"
	"testing"

	reindex "github.com/ActiveMemory/ctx/internal/cli/kb/core/reindex"
	cfgKbCli "github.com/ActiveMemory/ctx/internal/config/kb/cli"
)

func TestRenderBlock_NestedSlug(t *testing.T) {
	block := reindex.RenderBlock([]string{"g/t", "flat"})

	// Delimited by the managed-block markers.
	if !strings.HasPrefix(block, cfgKbCli.ManagedKBTopicsStart) {
		t.Errorf("block missing start marker:\n%s", block)
	}
	if !strings.HasSuffix(block, cfgKbCli.ManagedKBTopicsEnd) {
		t.Errorf("block missing end marker:\n%s", block)
	}
	// A grouped slug renders a working nested link.
	if !strings.Contains(block, "- [`g/t`](topics/g/t/)") {
		t.Errorf("nested slug not rendered as a topics/g/t/ link:\n%s", block)
	}
	// A flat slug still renders.
	if !strings.Contains(block, "- [`flat`](topics/flat/)") {
		t.Errorf("flat slug not rendered:\n%s", block)
	}
}

func TestRenderBlock_Empty(t *testing.T) {
	block := reindex.RenderBlock(nil)
	if !strings.Contains(block, cfgKbCli.ManagedKBTopicsEmpty) {
		t.Errorf("empty block missing placeholder:\n%s", block)
	}
}
