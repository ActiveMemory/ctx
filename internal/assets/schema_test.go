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

	"gopkg.in/yaml.v3"

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

func TestSchemaCoversCtxRC(t *testing.T) {
	// Parse the schema to get its property keys.
	schemaData, readErr := FS.ReadFile(asset.PathCtxrcSchema)
	if readErr != nil {
		t.Fatalf("read schema: %v", readErr)
	}
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if parseErr := json.Unmarshal(schemaData, &schema); parseErr != nil {
		t.Fatalf("parse schema: %v", parseErr)
	}

	// Parse a zero-value CtxRC to YAML then back to a map to get yaml tags.
	// We marshal a struct with all fields set to get every key emitted.
	type ctxRC struct {
		Profile             string `yaml:"profile"`
		TokenBudget         int    `yaml:"token_budget"`
		PriorityOrder       []int  `yaml:"priority_order"`
		AutoArchive         bool   `yaml:"auto_archive"`
		ArchiveAfterDays    int    `yaml:"archive_after_days"`
		ScratchpadEncrypt   *bool  `yaml:"scratchpad_encrypt"`
		EntryCountLearnings int    `yaml:"entry_count_learnings"`
		EntryCountDecisions int    `yaml:"entry_count_decisions"`
		ConventionLineCount int    `yaml:"convention_line_count"`
		InjectionTokenWarn  int    `yaml:"injection_token_warn"`
		ContextWindow       int    `yaml:"context_window"`
		BillingTokenWarn    int    `yaml:"billing_token_warn"`
		EventLog            bool   `yaml:"event_log"`
		KeyRotationDays     int    `yaml:"key_rotation_days"`
		TaskNudgeInterval   int    `yaml:"task_nudge_interval"`
		KeyPathOverride     string `yaml:"key_path"`
		StaleAgeDays        int    `yaml:"stale_age_days"`
		SessionPrefixes     []int  `yaml:"session_prefixes"`
		CompanionCheck      *bool  `yaml:"companion_check"`
		ClassifyRules       []int  `yaml:"classify_rules"`
		SpecSignalWords     []int  `yaml:"spec_signal_words"`
		SpecNudgeMinLen     int    `yaml:"spec_nudge_min_len"`
		Placeholders        []int  `yaml:"placeholders"`
		Notify              *int   `yaml:"notify"`
		FreshnessFiles      []int  `yaml:"freshness_files"`
		Tool                string `yaml:"tool"`
		Steering            *int   `yaml:"steering"`
		Hooks               *int   `yaml:"hooks"`
		ProvenanceRequired  *int   `yaml:"provenance_required"`
	}
	yamlBytes, marshalErr := yaml.Marshal(ctxRC{})
	if marshalErr != nil {
		t.Fatalf("marshal: %v", marshalErr)
	}
	var structKeys map[string]any
	if unmarshalErr := yaml.Unmarshal(yamlBytes, &structKeys); unmarshalErr != nil {
		t.Fatalf("unmarshal: %v", unmarshalErr)
	}

	// Every struct field must appear in schema.
	for key := range structKeys {
		if _, ok := schema.Properties[key]; !ok {
			t.Errorf("CtxRC field %q has no schema property", key)
		}
	}
	// Every schema property must appear in struct.
	for key := range schema.Properties {
		if _, ok := structKeys[key]; !ok {
			t.Errorf("schema property %q has no CtxRC field", key)
		}
	}
}
