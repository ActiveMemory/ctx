//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// TestSchemaCoversCtxRC asserts a bijection between the yaml tags of
// the real rc.CtxRC struct and the properties declared in the embedded
// .ctxrc JSON schema. It reflects over the live struct (not a
// hand-maintained copy) so a field added to CtxRC without a matching
// schema property — or vice versa — fails the build's tests.
func TestSchemaCoversCtxRC(t *testing.T) {
	schemaData, readErr := assets.FS.ReadFile(asset.PathCtxrcSchema)
	if readErr != nil {
		t.Fatalf("read schema: %v", readErr)
	}
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if parseErr := json.Unmarshal(schemaData, &schema); parseErr != nil {
		t.Fatalf("parse schema: %v", parseErr)
	}

	structKeys := make(map[string]bool)
	rt := reflect.TypeOf(CtxRC{})
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("yaml")
		if tag == "" || tag == "-" {
			continue
		}
		// Strip options like ",omitempty" to get the bare key.
		name := strings.Split(tag, ",")[0]
		if name == "" {
			continue
		}
		structKeys[name] = true
	}

	if len(structKeys) == 0 {
		t.Fatal("reflected zero yaml-tagged fields from CtxRC")
	}

	// Every struct field must appear in the schema.
	for key := range structKeys {
		if _, ok := schema.Properties[key]; !ok {
			t.Errorf("CtxRC field %q has no schema property", key)
		}
	}
	// Every schema property must map to a struct field.
	for key := range schema.Properties {
		if _, ok := structKeys[key]; !ok {
			t.Errorf("schema property %q has no CtxRC field", key)
		}
	}
}
