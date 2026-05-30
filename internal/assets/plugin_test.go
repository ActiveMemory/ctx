//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
)

func TestPluginVersion(t *testing.T) {
	data, err := FS.ReadFile(asset.PathPluginJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var manifest map[string]json.RawMessage
	if unmarshalErr := json.Unmarshal(data, &manifest); unmarshalErr != nil {
		t.Fatalf("parse error: %v", unmarshalErr)
	}
	raw, ok := manifest[asset.JSONKeyVersion]
	if !ok {
		t.Fatal("plugin.json missing 'version' key")
	}
	var ver string
	if parseErr := json.Unmarshal(raw, &ver); parseErr != nil {
		t.Fatalf("version parse error: %v", parseErr)
	}
	if ver == "" {
		t.Error("version is empty")
	}
	if !strings.Contains(ver, ".") {
		t.Errorf("version = %q, expected semver format", ver)
	}
}
