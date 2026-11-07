//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	cfgMcpTool "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

func TestDefsCount(t *testing.T) {
	if len(Defs()) != 15 {
		t.Errorf("tool count = %d, want 15", len(Defs()))
	}
}

func TestDefsNoDuplicateNames(t *testing.T) {
	seen := make(map[string]bool)
	for _, d := range Defs() {
		if seen[d.Name] {
			t.Errorf("duplicate tool name: %s", d.Name)
		}
		seen[d.Name] = true
	}
}

func TestDefsAllNamed(t *testing.T) {
	for i, d := range Defs() {
		if d.Name == "" {
			t.Errorf("tool[%d] has empty name", i)
		}
	}
}

// Note: Description fields are populated by desc.Text() at package
// init time. They are verified as non-empty in the server integration
// tests where lookup.Init() runs before this package is imported.

func TestDefsAllHaveObjectSchema(t *testing.T) {
	for _, d := range Defs() {
		if d.InputSchema.Type != "object" {
			t.Errorf(
				"tool %q schema type = %q, want %q",
				d.Name, d.InputSchema.Type, "object",
			)
		}
	}
}

func TestDefsContainsAllConfigTools(t *testing.T) {
	want := []string{
		cfgMcpTool.Status,
		cfgMcpTool.Add,
		cfgMcpTool.Complete,
		cfgMcpTool.Drift,
		cfgMcpTool.JournalSource,
		cfgMcpTool.WatchUpdate,
		cfgMcpTool.Compact,
		cfgMcpTool.Next,
		cfgMcpTool.CheckTaskCompletion,
		cfgMcpTool.SessionEvent,
		cfgMcpTool.Remind,
		cfgMcpTool.SteeringGet,
		cfgMcpTool.Search,
		cfgMcpTool.SessionStart,
		cfgMcpTool.SessionEnd,
	}
	names := make(map[string]bool)
	for _, d := range Defs() {
		names[d.Name] = true
	}
	for _, w := range want {
		if !names[w] {
			t.Errorf("missing tool: %s", w)
		}
	}
}

func TestDefsAnnotations(t *testing.T) {
	for _, d := range Defs() {
		if d.Annotations == nil {
			t.Errorf(
				"tool %q has nil annotations", d.Name,
			)
		}
	}
}

func TestDefsAddRequiredFields(t *testing.T) {
	for _, d := range Defs() {
		if d.Name != cfgMcpTool.Add {
			continue
		}
		if len(d.InputSchema.Required) < 2 {
			t.Errorf(
				"add tool requires at least 2 fields, got %d",
				len(d.InputSchema.Required),
			)
		}
		return
	}
	t.Error("add tool not found in Defs")
}

// Entry validation (internal/entry/validate.go) can require each of
// these fields, and the handlers read priority alongside them. A
// client can only send what a schema advertises, so both entry tools
// must declare every one.
func TestDefsEntryToolsAdvertiseEntryFields(t *testing.T) {
	missing := map[string]bool{
		cfgMcpTool.Add: true, cfgMcpTool.WatchUpdate: true,
	}
	for _, d := range Defs() {
		if !missing[d.Name] {
			continue
		}
		delete(missing, d.Name)
		for _, key := range []string{
			field.Section, field.Priority,
			field.SessionID, field.Branch, field.Commit,
			cli.AttrContext, cli.AttrRationale, cli.AttrConsequence,
			cli.AttrLesson, cli.AttrApplication,
		} {
			if _, ok := d.InputSchema.Properties[key]; !ok {
				t.Errorf("%s schema missing prop: %s", d.Name, key)
			}
		}
	}
	for name := range missing {
		t.Errorf("entry tool %s not found in Defs", name)
	}
}

func TestDefsMergeProps(t *testing.T) {
	dst := map[string]proto.Property{
		"a": {Type: "string"},
	}
	src := map[string]proto.Property{
		"b": {Type: "number"},
	}
	result := MergeProps(dst, src)
	if len(result) != 2 {
		t.Errorf("merged length = %d, want 2", len(result))
	}
	if result["b"].Type != "number" {
		t.Errorf(
			"result[b].Type = %q, want %q",
			result["b"].Type, "number",
		)
	}
}

func TestDefsEntryAttrProps(t *testing.T) {
	props := EntryAttrProps("test.key")
	expected := []string{
		"context", "rationale", "consequence",
		"lesson", "application",
	}
	for _, key := range expected {
		if _, ok := props[key]; !ok {
			t.Errorf("missing entry attr prop: %s", key)
		}
	}
}
