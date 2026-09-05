//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
	cfgVersion "github.com/ActiveMemory/ctx/internal/config/version"
)

// repoRoot is the project root relative to this package directory
// (go test runs with the package directory as cwd).
var repoRoot = filepath.Join("..", "..")

// repoVersion returns the trimmed contents of the project-root
// VERSION file, the single source of truth every plugin manifest
// must agree with (hack/release.sh and `make sync-version` keep
// them in step; this is the guard that they did).
func repoVersion(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, cfgVersion.FileVersion))
	if err != nil {
		t.Fatalf("read %s: %v", cfgVersion.FileVersion, err)
	}
	ver := strings.TrimSpace(string(data))
	if ver == "" {
		t.Fatalf("%s is empty", cfgVersion.FileVersion)
	}
	return ver
}

// manifestVersion reads the top-level "version" string of an
// embedded plugin manifest.
func manifestVersion(t *testing.T, embeddedPath string) string {
	t.Helper()
	data, err := FS.ReadFile(embeddedPath)
	if err != nil {
		t.Fatalf("%s: unexpected error: %v", embeddedPath, err)
	}
	var manifest map[string]json.RawMessage
	if unmarshalErr := json.Unmarshal(data, &manifest); unmarshalErr != nil {
		t.Fatalf("%s: parse error: %v", embeddedPath, unmarshalErr)
	}
	raw, ok := manifest[asset.JSONKeyVersion]
	if !ok {
		t.Fatalf("%s missing 'version' key", embeddedPath)
	}
	var ver string
	if parseErr := json.Unmarshal(raw, &ver); parseErr != nil {
		t.Fatalf("%s: version parse error: %v", embeddedPath, parseErr)
	}
	return ver
}

// TestPluginVersion asserts that every embedded plugin manifest
// (Claude Code and Codex) carries a semver version equal to the
// project-root VERSION file.
func TestPluginVersion(t *testing.T) {
	want := repoVersion(t)
	for _, manifestPath := range []string{
		asset.PathPluginJSON,
		asset.PathCodexPluginJSON,
	} {
		t.Run(manifestPath, func(t *testing.T) {
			ver := manifestVersion(t, manifestPath)
			if ver == "" {
				t.Error("version is empty")
			}
			if !strings.Contains(ver, ".") {
				t.Errorf("version = %q, expected semver format", ver)
			}
			if ver != want {
				t.Errorf(
					"version = %q, want %q (VERSION file); run `make sync-version`",
					ver, want,
				)
			}
		})
	}
}
