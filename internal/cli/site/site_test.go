//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestCmd_HasFeedSubcommand verifies the site command includes "feed".
func TestCmd_HasFeedSubcommand(t *testing.T) {
	cmd := Cmd()

	found := false
	for _, sub := range cmd.Commands() {
		if sub.Name() == "feed" {
			found = true
			break
		}
	}

	if !found {
		t.Error("site command should have a 'feed' subcommand")
	}
}

// newSiteTestCmd creates a cobra command with a captured output buffer.
func newSiteTestCmd() *cobra.Command {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

// siteTestOutput returns the captured output from a test command.
func siteTestOutput(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}

// TestPrintReport_NoSkipped verifies clean output with no issues.
func TestPrintReport_NoSkipped(t *testing.T) {
	cmd := newSiteTestCmd()
	report := feedReport{
		included: 3,
	}

	printReport(cmd, "site/feed.xml", report)
	out := siteTestOutput(cmd)

	if !strings.Contains(out, "3 entries") {
		t.Errorf("expected '3 entries' in output, got: %s", out)
	}
	if strings.Contains(out, "Skipped:") {
		t.Error("output should not contain 'Skipped:' section")
	}
	if strings.Contains(out, "Warnings:") {
		t.Error("output should not contain 'Warnings:' section")
	}
}

// TestPrintReport_WithWarnings verifies warnings section appears.
func TestPrintReport_WithWarnings(t *testing.T) {
	cmd := newSiteTestCmd()
	report := feedReport{
		included: 2,
		warnings: []string{"post.md \u2014 no summary paragraph found"},
	}

	printReport(cmd, "site/feed.xml", report)
	out := siteTestOutput(cmd)

	if !strings.Contains(out, "Warnings:") {
		t.Errorf("expected 'Warnings:' section in output, got: %s", out)
	}
	if !strings.Contains(out, "no summary") {
		t.Errorf("expected warning message in output, got: %s", out)
	}
}

// TestRunFeed_NoBlogDir verifies error when blog directory is missing.
func TestRunFeed_NoBlogDir(t *testing.T) {
	cmd := newSiteTestCmd()
	runErr := runFeed(cmd, "/nonexistent/blog/dir", "out.xml", "https://example.com")
	if runErr == nil {
		t.Fatal("expected error for nonexistent blog directory")
	}
	if !strings.Contains(runErr.Error(), "no blog directory") {
		t.Errorf("expected 'no blog directory' error, got: %v", runErr)
	}
}
