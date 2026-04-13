//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"testing"

	cfgPrompt "github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
)

func TestDefsCount(t *testing.T) {
	if len(Defs) != 5 {
		t.Errorf("prompt count = %d, want 5", len(Defs))
	}
}

func TestDefsNoDuplicateNames(t *testing.T) {
	seen := make(map[string]bool)
	for _, d := range Defs {
		if seen[d.Name] {
			t.Errorf("duplicate prompt name: %s", d.Name)
		}
		seen[d.Name] = true
	}
}

func TestDefsAllNamed(t *testing.T) {
	for i, d := range Defs {
		if d.Name == "" {
			t.Errorf("prompt[%d] has empty name", i)
		}
	}
}

func TestDefsContainsAllConfigPrompts(t *testing.T) {
	want := []string{
		cfgPrompt.SessionStart,
		cfgPrompt.AddDecision,
		cfgPrompt.AddLearning,
		cfgPrompt.Reflect,
		cfgPrompt.Checkpoint,
	}
	names := make(map[string]bool)
	for _, d := range Defs {
		names[d.Name] = true
	}
	for _, w := range want {
		if !names[w] {
			t.Errorf("missing prompt: %s", w)
		}
	}
}

func TestDefsAddDecisionArgs(t *testing.T) {
	for _, d := range Defs {
		if d.Name != cfgPrompt.AddDecision {
			continue
		}
		if len(d.Arguments) != 4 {
			t.Errorf(
				"add-decision argument count = %d, want 4",
				len(d.Arguments),
			)
		}
		for _, a := range d.Arguments {
			if !a.Required {
				t.Errorf(
					"argument %q should be required", a.Name,
				)
			}
		}
		return
	}
	t.Error("add-decision prompt not found")
}

func TestDefsAddLearningArgs(t *testing.T) {
	for _, d := range Defs {
		if d.Name != cfgPrompt.AddLearning {
			continue
		}
		if len(d.Arguments) != 4 {
			t.Errorf(
				"add-learning argument count = %d, want 4",
				len(d.Arguments),
			)
		}
		for _, a := range d.Arguments {
			if !a.Required {
				t.Errorf(
					"argument %q should be required", a.Name,
				)
			}
		}
		return
	}
	t.Error("add-learning prompt not found")
}

func TestDefsSessionStartNoArgs(t *testing.T) {
	for _, d := range Defs {
		if d.Name != cfgPrompt.SessionStart {
			continue
		}
		if len(d.Arguments) != 0 {
			t.Errorf(
				"session-start should have 0 args, got %d",
				len(d.Arguments),
			)
		}
		return
	}
	t.Error("session-start prompt not found")
}

func TestDefsReflectNoArgs(t *testing.T) {
	for _, d := range Defs {
		if d.Name != cfgPrompt.Reflect {
			continue
		}
		if len(d.Arguments) != 0 {
			t.Errorf(
				"reflect should have 0 args, got %d",
				len(d.Arguments),
			)
		}
		return
	}
	t.Error("reflect prompt not found")
}

func TestDefsCheckpointNoArgs(t *testing.T) {
	for _, d := range Defs {
		if d.Name != cfgPrompt.Checkpoint {
			continue
		}
		if len(d.Arguments) != 0 {
			t.Errorf(
				"checkpoint should have 0 args, got %d",
				len(d.Arguments),
			)
		}
		return
	}
	t.Error("checkpoint prompt not found")
}
