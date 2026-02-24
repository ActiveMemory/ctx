//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

// helper: build a turn header line.
func turnHeader(n int, role, ts string) string {
	return fmt.Sprintf("### %d. %s (%s)", n, role, ts)
}

// helper: generate N numbered body lines.
func bodyLines(n int) string {
	var lines []string
	for i := 1; i <= n; i++ {
		lines = append(lines, fmt.Sprintf("line %d", i))
	}
	return strings.Join(lines, config.NewlineLF)
}

func TestCollapseToolOutputs_LongOutputWrapped(t *testing.T) {
	header := turnHeader(1, config.LabelToolOutput, "10:00:00")
	body := bodyLines(12)
	input := header + "\n\n" + body + "\n"

	got := collapseToolOutputs(input)

	if !strings.Contains(got, "<details>") {
		t.Error("expected <details> tag for long tool output")
	}
	if !strings.Contains(got, "</details>") {
		t.Error("expected </details> closing tag")
	}
	if !strings.Contains(got, "<summary>12 lines</summary>") {
		t.Errorf("expected summary with 12 lines, got:\n%s", got)
	}
	// Body content must still be present
	if !strings.Contains(got, "line 1") {
		t.Error("body content missing after collapse")
	}
	if !strings.Contains(got, "line 12") {
		t.Error("body content missing after collapse")
	}
}

func TestCollapseToolOutputs_ShortOutputUnchanged(t *testing.T) {
	header := turnHeader(1, config.LabelToolOutput, "10:00:00")
	body := bodyLines(5)
	input := header + "\n\n" + body + "\n"

	got := collapseToolOutputs(input)

	if strings.Contains(got, "<details>") {
		t.Error("short tool output should not be wrapped in <details>")
	}
	if !strings.Contains(got, "line 5") {
		t.Error("body content missing")
	}
}

func TestCollapseToolOutputs_ExactThresholdUnchanged(t *testing.T) {
	header := turnHeader(1, config.LabelToolOutput, "10:00:00")
	body := bodyLines(config.RecallDetailsThreshold)
	input := header + "\n\n" + body + "\n"

	got := collapseToolOutputs(input)

	if strings.Contains(got, "<details>") {
		t.Error("output at exactly the threshold should not be wrapped")
	}
}

func TestCollapseToolOutputs_AlreadyWrappedNotDoubled(t *testing.T) {
	header := turnHeader(1, config.LabelToolOutput, "10:00:00")
	body := "<details>\n<summary>15 lines</summary>\n\n" +
		bodyLines(15) + "\n</details>"
	input := header + "\n\n" + body + "\n"

	got := collapseToolOutputs(input)

	count := strings.Count(got, "<details>")
	if count != 1 {
		t.Errorf("expected 1 <details> block, got %d", count)
	}
}

func TestCollapseToolOutputs_NonToolTurnsUntouched(t *testing.T) {
	user := turnHeader(1, "User", "10:00:00") + "\n\n" + bodyLines(15) + "\n"
	assistant := turnHeader(2, "Assistant", "10:01:00") + "\n\n" +
		bodyLines(15) + "\n"
	input := user + "\n" + assistant

	got := collapseToolOutputs(input)

	if strings.Contains(got, "<details>") {
		t.Error("non-tool-output turns should never be collapsed")
	}
}

func TestCollapseToolOutputs_MixedTurns(t *testing.T) {
	short := turnHeader(1, config.LabelToolOutput, "10:00:00") +
		"\n\n" + bodyLines(3) + "\n"
	long := turnHeader(2, config.LabelToolOutput, "10:01:00") +
		"\n\n" + bodyLines(15) + "\n"
	user := turnHeader(3, "User", "10:02:00") +
		"\n\n" + bodyLines(20) + "\n"

	input := short + "\n" + long + "\n" + user

	got := collapseToolOutputs(input)

	count := strings.Count(got, "<details>")
	if count != 1 {
		t.Errorf("expected exactly 1 <details> block, got %d", count)
	}
	if !strings.Contains(got, "<summary>15 lines</summary>") {
		t.Error("the long tool output should show 15 lines in summary")
	}
}

func TestCollapseToolOutputs_NoHeaders(t *testing.T) {
	input := "Just some text\nwithout any turn headers\n"

	got := collapseToolOutputs(input)

	if got != input {
		t.Errorf("content without headers should pass through unchanged\ngot:\n%s", got)
	}
}
