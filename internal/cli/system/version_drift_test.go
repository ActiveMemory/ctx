//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
)

func TestCheckVersionDrift_AllMatch(t *testing.T) {
	tmp := t.TempDir()
	orig, dirErr := os.Getwd()
	if dirErr != nil {
		t.Fatal(dirErr)
	}
	if chErr := os.Chdir(tmp); chErr != nil {
		t.Fatal(chErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	// Write VERSION file matching the embedded plugin version.
	pluginVer := readPluginVersionOrSkip(t)

	if writeErr := os.WriteFile("VERSION", []byte(pluginVer+"\n"), 0644); writeErr != nil {
		t.Fatal(writeErr)
	}
	writeMarketplace(t, tmp, pluginVer)

	// Capture output — should be silent when all match.
	cmd := postCommitCmd()
	var buf strings.Builder
	cmd.SetOut(&buf)

	checkVersionDrift(cmd, "test-session")

	if out := buf.String(); out != "" {
		t.Errorf("expected silent output, got %q", out)
	}
}

func TestCheckVersionDrift_Mismatch(t *testing.T) {
	tmp := t.TempDir()
	orig, dirErr := os.Getwd()
	if dirErr != nil {
		t.Fatal(dirErr)
	}
	if chErr := os.Chdir(tmp); chErr != nil {
		t.Fatal(chErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	// Create .context so isInitialized-related helpers don't interfere.
	pluginVer := readPluginVersionOrSkip(t)
	mismatchVer := "99.99.99"

	if writeErr := os.WriteFile("VERSION", []byte(mismatchVer+"\n"), 0644); writeErr != nil {
		t.Fatal(writeErr)
	}
	writeMarketplace(t, tmp, pluginVer)

	cmd := postCommitCmd()
	var buf strings.Builder
	cmd.SetOut(&buf)

	checkVersionDrift(cmd, "test-session")

	out := buf.String()
	if out == "" {
		t.Fatal("expected drift output, got empty")
	}
	if !strings.Contains(out, mismatchVer) {
		t.Errorf("output missing VERSION value %q: %s", mismatchVer, out)
	}
	if !strings.Contains(out, pluginVer) {
		t.Errorf("output missing plugin version %q: %s", pluginVer, out)
	}
}

func TestReadMarketplaceVersion(t *testing.T) {
	tmp := t.TempDir()
	orig, dirErr := os.Getwd()
	if dirErr != nil {
		t.Fatal(dirErr)
	}
	if chErr := os.Chdir(tmp); chErr != nil {
		t.Fatal(chErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	writeMarketplace(t, tmp, "1.2.3")

	got := readMarketplaceVersion()
	if got != "1.2.3" {
		t.Errorf("readMarketplaceVersion() = %q, want %q", got, "1.2.3")
	}
}

func TestReadMarketplaceVersion_Missing(t *testing.T) {
	tmp := t.TempDir()
	orig, dirErr := os.Getwd()
	if dirErr != nil {
		t.Fatal(dirErr)
	}
	if chErr := os.Chdir(tmp); chErr != nil {
		t.Fatal(chErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	got := readMarketplaceVersion()
	if got != "" {
		t.Errorf("readMarketplaceVersion() = %q, want empty", got)
	}
}

// --- helpers ---

func readPluginVersionOrSkip(t *testing.T) string {
	t.Helper()
	ver, verErr := assets.PluginVersion()
	if verErr != nil {
		t.Skipf("cannot read embedded plugin version: %v", verErr)
	}
	return ver
}

func writeMarketplace(t *testing.T, root, version string) {
	t.Helper()
	dir := filepath.Join(root, ".claude-plugin")
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		t.Fatal(mkErr)
	}
	manifest := marketplaceManifest{
		Plugins: []struct {
			Version string `json:"version"`
		}{{Version: version}},
	}
	data, marshalErr := json.Marshal(manifest)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	if writeErr := os.WriteFile(filepath.Join(dir, "marketplace.json"), data, 0644); writeErr != nil {
		t.Fatal(writeErr)
	}
}
