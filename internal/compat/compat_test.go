//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compat

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/agent/core/budget"
	cfgTrigger "github.com/ActiveMemory/ctx/internal/config/trigger"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/skill"
	"github.com/ActiveMemory/ctx/internal/steering"
	"github.com/ActiveMemory/ctx/internal/trigger"
)

// TestBackwardCompat_AssemblePacket_NoSteeringNoSkill verifies that
// AssemblePacket with nil steering and empty skill produces a packet
// without Steering or Skill sections — identical to pre-extension
// behaviour.
//
// Validates: Requirements 14.1
func TestBackwardCompat_AssemblePacket_NoSteeringNoSkill(t *testing.T) {
	ctx := &entity.Context{}

	pkt := budget.AssemblePacket(ctx, 8000, nil, "")

	if len(pkt.Steering) != 0 {
		t.Errorf("expected no steering entries, got %d", len(pkt.Steering))
	}
	if pkt.Skill != "" {
		t.Errorf("expected empty skill, got %q", pkt.Skill)
	}

	// Verify the packet still contains the always-present fields.
	if pkt.Budget != 8000 {
		t.Errorf("expected budget 8000, got %d", pkt.Budget)
	}
	if pkt.Instruction == "" {
		t.Error("expected non-empty instruction")
	}

	// Render to markdown and confirm no Steering/Skill sections appear.
	md := budget.RenderMarkdownPacket(pkt)
	if strings.Contains(md, "## Steering") {
		t.Error("rendered markdown should not contain Steering section")
	}
	if strings.Contains(md, "## Skill") {
		t.Error("rendered markdown should not contain Skill section")
	}
}

// TestBackwardCompat_HookRunAll_NonExistentDir verifies that RunAll
// on a non-existent hooks directory returns an empty AggregatedOutput
// without error.
//
// Validates: Requirements 14.2
func TestBackwardCompat_HookRunAll_NonExistentDir(t *testing.T) {
	nonexistent := filepath.Join(t.TempDir(), "no-such-hooks")
	input := &trigger.HookInput{TriggerType: "pre-tool-use", Tool: "test"}

	agg, err := trigger.RunAll(nonexistent, cfgTrigger.PreToolUse, input, 5*time.Second)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if agg == nil {
		t.Fatal("expected non-nil AggregatedOutput")
	}
	if agg.Cancelled {
		t.Error("expected Cancelled to be false")
	}
	if agg.Context != "" {
		t.Errorf("expected empty context, got %q", agg.Context)
	}
	if agg.Message != "" {
		t.Errorf("expected empty message, got %q", agg.Message)
	}
	if len(agg.Errors) != 0 {
		t.Errorf("expected no errors, got %v", agg.Errors)
	}
}

// TestBackwardCompat_HookDiscover_NonExistentDir verifies that Discover
// on a non-existent hooks directory returns an empty map without error.
//
// Validates: Requirements 14.2
func TestBackwardCompat_HookDiscover_NonExistentDir(t *testing.T) {
	nonexistent := filepath.Join(t.TempDir(), "no-such-hooks")

	result, err := trigger.Discover(nonexistent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

// TestBackwardCompat_SkillLoadAll_NonExistentDir verifies that LoadAll
// on a non-existent skills directory returns nil without error.
//
// Validates: Requirements 14.4
func TestBackwardCompat_SkillLoadAll_NonExistentDir(t *testing.T) {
	nonexistent := filepath.Join(t.TempDir(), "no-such-skills")

	skills, err := skill.LoadAll(nonexistent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if skills != nil {
		t.Errorf("expected nil, got %d skills", len(skills))
	}
}

// TestBackwardCompat_SteeringLoadAll_NonExistentDir verifies that
// LoadAll on a non-existent steering directory returns an error, which
// callers (like ctx agent) handle gracefully by skipping steering.
//
// Validates: Requirements 14.1
func TestBackwardCompat_SteeringLoadAll_NonExistentDir(t *testing.T) {
	nonexistent := filepath.Join(t.TempDir(), "no-such-steering")

	files, err := steering.LoadAll(nonexistent)
	if err == nil {
		t.Fatal("expected error for non-existent steering directory")
	}
	if files != nil {
		t.Errorf("expected nil files on error, got %d", len(files))
	}

	// The error should be an os-level "not exist" error that callers
	// can detect with os.IsNotExist or errors.Is.
	if !os.IsNotExist(unwrapPathError(err)) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

// TestBackwardCompat_FullAgentPath_NoExtensions exercises the full
// backward-compatible agent assembly path: no steering directory,
// no skills directory, no hooks directory. The resulting packet should
// be structurally identical to the pre-extension output.
//
// Validates: Requirements 14.1, 14.2, 14.3, 14.4, 14.5
func TestBackwardCompat_FullAgentPath_NoExtensions(t *testing.T) {
	ctx := &entity.Context{}

	// Simulate the agent path: no steering files loaded (directory
	// missing → error → caller passes nil), no skill.
	pkt := budget.AssemblePacket(ctx, 8000, nil, "")

	// Verify core structure is intact.
	if pkt.Budget != 8000 {
		t.Errorf("budget = %d, want 8000", pkt.Budget)
	}
	if pkt.Instruction == "" {
		t.Error("instruction should be populated from embedded assets")
	}

	// Verify no extension sections are present.
	if len(pkt.Steering) != 0 {
		t.Errorf("steering should be empty, got %d", len(pkt.Steering))
	}
	if pkt.Skill != "" {
		t.Errorf("skill should be empty, got %q", pkt.Skill)
	}

	// Render and verify the markdown output has no extension sections.
	md := budget.RenderMarkdownPacket(pkt)
	if !strings.Contains(md, "# Context Packet") {
		t.Error("rendered markdown should contain Context Packet header")
	}
	for _, section := range []string{"## Steering", "## Skill"} {
		if strings.Contains(md, section) {
			t.Errorf("rendered markdown should not contain %q when no extensions are active", section)
		}
	}
}

// unwrapPathError extracts the underlying error from a wrapped path
// error chain for os.IsNotExist checking.
func unwrapPathError(err error) error {
	for {
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}
		err = u.Unwrap()
	}
}
