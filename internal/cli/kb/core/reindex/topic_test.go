//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	reindex "github.com/ActiveMemory/ctx/internal/cli/kb/core/reindex"
)

// mkTopic creates <topicsDir>/<slug>/index.md (slug may contain
// slashes for grouped layouts).
func mkTopic(t *testing.T, topicsDir, slug string) {
	t.Helper()
	dir := filepath.Join(topicsDir, filepath.FromSlash(slug))
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(dir, "index.md"),
		[]byte("# topic\n"), 0o600,
	); err != nil {
		t.Fatal(err)
	}
}

// mkDir creates <topicsDir>/<rel>/ with no index.md.
func mkDir(t *testing.T, topicsDir, rel string) {
	t.Helper()
	if err := os.MkdirAll(
		filepath.Join(topicsDir, filepath.FromSlash(rel)), 0o750,
	); err != nil {
		t.Fatal(err)
	}
}

func TestListTopics(t *testing.T) {
	tests := []struct {
		name  string
		flat  []string // slugs to create with an index.md
		bare  []string // dirs to create WITHOUT an index.md
		want  []string
		setup func(t *testing.T, dir string) // optional extra setup
	}{
		{
			name: "flat topics",
			flat: []string{"beta", "alpha"},
			want: []string{"alpha", "beta"},
		},
		{
			name: "grouped topics",
			flat: []string{"g1/t2", "g1/t1", "g2/t3"},
			want: []string{"g1/t1", "g1/t2", "g2/t3"},
		},
		{
			name: "mixed flat and grouped",
			flat: []string{"flat", "grp/nested"},
			want: []string{"flat", "grp/nested"},
		},
		{
			// topics/g/index.md is a group-landing (g also holds
			// topics/g/t/index.md), so g is excluded; only g/t is a
			// topic.
			name: "group-landing excluded",
			flat: []string{"g", "g/t"},
			want: []string{"g/t"},
		},
		{
			name: "dir without index.md is not a topic",
			flat: []string{"real"},
			bare: []string{"empty"},
			want: []string{"real"},
		},
		{
			name: "deep nesting without intermediate index",
			flat: []string{"a/b/c"},
			want: []string{"a/b/c"},
		},
		{
			name: "no topics",
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			topicsDir := t.TempDir()
			for _, s := range tc.flat {
				mkTopic(t, topicsDir, s)
			}
			for _, d := range tc.bare {
				mkDir(t, topicsDir, d)
			}
			got, err := reindex.ListTopics(topicsDir)
			if err != nil {
				t.Fatalf("ListTopics: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestListTopics_NonexistentDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	got, err := reindex.ListTopics(missing)
	if err != nil {
		t.Fatalf("ListTopics on missing dir: %v", err)
	}
	if got != nil {
		t.Errorf("got %v, want nil for a missing topics dir", got)
	}
}
