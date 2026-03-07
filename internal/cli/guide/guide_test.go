//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// executeGuide runs the guide command with the given args and returns output.
func executeGuide(t *testing.T, args ...string) string {
	t.Helper()

	root := &cobra.Command{Use: "ctx"}
	root.AddCommand(Cmd())

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(append([]string{"guide"}, args...))

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	return buf.String()
}

func TestGuideDefaultOutput(t *testing.T) {
	out := executeGuide(t)

	sections := []string{
		"GETTING STARTED",
		"TRACKING DECISIONS",
		"BROWSING HISTORY",
		"AI CONTEXT",
		"MAINTENANCE",
		"KEY SKILLS",
		"RECIPES",
	}
	for _, section := range sections {
		if !strings.Contains(out, section) {
			t.Errorf("default output missing section %q", section)
		}
	}

	// Check for flag hints.
	if !strings.Contains(out, "--skills") {
		t.Error("default output missing --skills hint")
	}
	if !strings.Contains(out, "--commands") {
		t.Error("default output missing --commands hint")
	}
}

func TestGuideSkillsFlag(t *testing.T) {
	out := executeGuide(t, "--skills")

	if !strings.Contains(out, "Available Skills:") {
		t.Error("--skills output missing header")
	}

	// Should contain skill names starting with /ctx-.
	lines := strings.Split(out, "\n")
	skillCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "/ctx-") {
			skillCount++
		}
	}
	if skillCount < 10 {
		t.Errorf("expected at least 10 skills, got %d", skillCount)
	}
}

func TestGuideCommandsFlag(t *testing.T) {
	root := &cobra.Command{Use: "ctx"}

	visible := &cobra.Command{Use: "status", Short: "Show status"}
	hidden := &cobra.Command{Use: "secret", Short: "Hidden cmd", Hidden: true}
	root.AddCommand(visible, hidden, Cmd())

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"guide", "--commands"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "status") {
		t.Error("--commands output missing visible command 'status'")
	}
	if strings.Contains(out, "secret") {
		t.Error("--commands output should not contain hidden command 'secret'")
	}
}
