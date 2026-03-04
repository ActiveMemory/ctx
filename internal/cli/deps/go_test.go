//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"testing"
)

func TestIsStdlib(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"fmt", true},
		{"os/exec", true},
		{"encoding/json", true},
		{"github.com/foo/bar", false},
		{"golang.org/x/tools", false},
	}
	for _, tt := range tests {
		if got := isStdlib(tt.input); got != tt.want {
			t.Errorf("isStdlib(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestShortPkgName(t *testing.T) {
	mod := "github.com/ActiveMemory/ctx"
	tests := []struct {
		input string
		want  string
	}{
		{"github.com/ActiveMemory/ctx/internal/cli/deps", "internal/cli/deps"},
		{"github.com/ActiveMemory/ctx", "github.com/ActiveMemory/ctx"},
		{"github.com/other/pkg", "github.com/other/pkg"},
	}
	for _, tt := range tests {
		if got := shortPkgName(tt.input, mod); got != tt.want {
			t.Errorf("shortPkgName(%q, %q) = %q, want %q", tt.input, mod, got, tt.want)
		}
	}
}

func TestGoModulePath(t *testing.T) {
	modPath := "github.com/test/project"
	pkgs := []goPackage{
		{ImportPath: "github.com/test/project/pkg", Module: &struct {
			Path string `json:"Path"`
		}{Path: modPath}},
	}
	if got := goModulePath(pkgs); got != modPath {
		t.Errorf("goModulePath() = %q, want %q", got, modPath)
	}

	// No module info.
	if got := goModulePath([]goPackage{{ImportPath: "test"}}); got != "" {
		t.Errorf("goModulePath() with no module = %q, want empty", got)
	}
}
