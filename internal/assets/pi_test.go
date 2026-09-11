//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"bytes"
	"io/fs"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// TestPiSkillsMirrorOpenCode freezes the intentionally identical Pi
// and OpenCode skill trees: every file under one tree must exist under
// the other with byte-identical content, so an edit to one side fails
// until it is mirrored. See specs/pi-cli-integration.md.
func TestPiSkillsMirrorOpenCode(t *testing.T) {
	opencode := readEmbeddedTree(t, asset.DirIntegrationsOpenCodeSkill)
	pi := readEmbeddedTree(t, asset.DirIntegrationsPiSkill)
	if len(opencode) == 0 {
		t.Fatalf("no files under %s", asset.DirIntegrationsOpenCodeSkill)
	}
	for rel, want := range opencode {
		got, ok := pi[rel]
		switch {
		case !ok:
			t.Errorf("%s: missing from %s", rel, asset.DirIntegrationsPiSkill)
		case !bytes.Equal(got, want):
			t.Errorf(
				"%s: differs between %s and %s",
				rel, asset.DirIntegrationsOpenCodeSkill,
				asset.DirIntegrationsPiSkill,
			)
		}
	}
	for rel := range pi {
		if _, ok := opencode[rel]; !ok {
			t.Errorf(
				"%s: missing from %s", rel, asset.DirIntegrationsOpenCodeSkill,
			)
		}
	}
}

// readEmbeddedTree maps every file under root in the embedded FS to
// its content, keyed by the path relative to root.
func readEmbeddedTree(t *testing.T, root string) map[string][]byte {
	t.Helper()
	sub, subErr := fs.Sub(FS, root)
	if subErr != nil {
		t.Fatalf("sub %s: %v", root, subErr)
	}
	out := make(map[string][]byte)
	walkErr := fs.WalkDir(
		sub, ".",
		func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			data, readErr := fs.ReadFile(sub, p)
			if readErr != nil {
				return readErr
			}
			out[p] = data
			return nil
		},
	)
	if walkErr != nil {
		t.Fatalf("walk %s: %v", root, walkErr)
	}
	return out
}
