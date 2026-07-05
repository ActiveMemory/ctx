//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package statusline

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	cfgStatusline "github.com/ActiveMemory/ctx/internal/config/statusline"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// render feeds the given payload through Run via a pipe and returns
// the rendered line.
func render(t *testing.T, payloadJSON string) string {
	t.Helper()
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("failed to create pipe: %v", pipeErr)
	}
	if _, writeErr := w.WriteString(payloadJSON); writeErr != nil {
		t.Fatalf("failed to write payload: %v", writeErr)
	}
	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("failed to close pipe writer: %v", closeErr)
	}
	defer func() { _ = r.Close() }()

	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	if runErr := Run(cmd, r); runErr != nil {
		t.Fatalf("Run returned error: %v", runErr)
	}
	return strings.TrimRight(buf.String(), "\n")
}

func TestRunRendersAllSegments(t *testing.T) {
	line := render(t, `{
		"model": {"display_name": "Opus"},
		"workspace": {"current_dir": "/tmp/project"},
		"cost": {"total_cost_usd": 1.2345},
		"context_window": {"used_percentage": 42.4}
	}`)

	for _, want := range []string{"Opus", "ctx: 42%", "$1.23"} {
		if !strings.Contains(line, want) {
			t.Errorf("line %q missing %q", line, want)
		}
	}
	if !strings.Contains(line, cfgStatusline.SegmentSeparator) {
		t.Errorf("line %q missing segment separator", line)
	}
}

func TestRunMalformedPayloadDegrades(t *testing.T) {
	line := render(t, `{not json at all`)
	if strings.Contains(line, "$") || strings.Contains(line, "ctx:") {
		t.Errorf("degraded line %q should carry no payload segments", line)
	}
}

func TestRunEmptyPayloadDegrades(t *testing.T) {
	line := render(t, `{}`)
	if strings.Contains(line, "$") || strings.Contains(line, "ctx:") {
		t.Errorf("line %q should drop segments for absent fields", line)
	}
}

func TestRunMissingCostOmitsSegment(t *testing.T) {
	line := render(t, `{
		"model": {"display_name": "Opus"},
		"context_window": {"used_percentage": 10}
	}`)
	if strings.Contains(line, "$") {
		t.Errorf("line %q renders a cost segment without cost data", line)
	}
}

func TestRunNullPercentageOmitsSegment(t *testing.T) {
	line := render(t, `{
		"model": {"display_name": "Opus"},
		"context_window": {"used_percentage": null},
		"cost": {"total_cost_usd": 0.5}
	}`)
	if strings.Contains(line, "ctx:") {
		t.Errorf("line %q renders ctx segment for null percentage", line)
	}
	if !strings.Contains(line, "$0.50") {
		t.Errorf("line %q missing cost segment", line)
	}
}

func TestRunOutOfRangePercentageDropped(t *testing.T) {
	line := render(t, `{
		"context_window": {"used_percentage": 250}
	}`)
	if strings.Contains(line, "ctx:") {
		t.Errorf("line %q renders out-of-range percentage", line)
	}
}

func TestRunSanitizesEscapedControlChars(t *testing.T) {
	// ANSI escape encoded as a legal JSON \u escape: the payload
	// parses, and sanitize must strip the ESC byte from the output.
	line := render(t, `{
		"model": {"display_name": "\u001b[31mOpus\u001b[0m\nInjected"}
	}`)
	if strings.ContainsAny(line, "\x1b\n\r\x07") {
		t.Errorf("line %q contains control bytes", line)
	}
	if !strings.Contains(line, "Opus") {
		t.Errorf("line %q lost the printable model name", line)
	}
}

func TestRunRawControlBytesDegrade(t *testing.T) {
	// A raw ESC byte inside a JSON string is invalid JSON. The
	// payload must fail closed: degraded line, no control bytes.
	line := render(t, "{\"model\": {\"display_name\": \"\x1b[31mOpus\"}}")
	if strings.ContainsAny(line, "\x1b\n\r\x07") {
		t.Errorf("line %q contains control bytes", line)
	}
	if strings.Contains(line, "Opus") {
		t.Errorf("line %q rendered a field from invalid JSON", line)
	}
}

func TestRunCapsLineLength(t *testing.T) {
	line := render(t, `{
		"model": {"display_name": "`+strings.Repeat("m", 500)+`"}
	}`)
	if len(line) > cfgStatusline.MaxLineLen {
		t.Errorf("line length %d exceeds cap %d",
			len(line), cfgStatusline.MaxLineLen)
	}
}

func TestRunDisabledRendersBlankLine(t *testing.T) {
	tmpDir, mkErr := os.MkdirTemp("", "ctx-statusline-off-*")
	if mkErr != nil {
		t.Fatalf("failed to create temp dir: %v", mkErr)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	if mkdirErr := os.Mkdir(filepath.Join(tmpDir, ".context"), 0o750); mkdirErr != nil {
		t.Fatalf("failed to create .context: %v", mkdirErr)
	}
	rcBody := "statusline:\n  enabled: false\n"
	if writeErr := os.WriteFile(
		filepath.Join(tmpDir, ".ctxrc"), []byte(rcBody), 0o600,
	); writeErr != nil {
		t.Fatalf("failed to write .ctxrc: %v", writeErr)
	}

	origDir, _ := os.Getwd()
	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatalf("failed to chdir: %v", chdirErr)
	}
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()
	defer rc.Reset()

	line := render(t, `{"model": {"display_name": "Opus"}}`)
	if line != "" {
		t.Errorf("disabled statusline rendered %q; want empty", line)
	}
}

func TestRunShowCostFalseSuppressesCost(t *testing.T) {
	tmpDir, mkErr := os.MkdirTemp("", "ctx-statusline-*")
	if mkErr != nil {
		t.Fatalf("failed to create temp dir: %v", mkErr)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// .ctxrc is only read when $PWD/.context/ exists (cwd-anchored
	// resolution model).
	if mkdirErr := os.Mkdir(filepath.Join(tmpDir, ".context"), 0o750); mkdirErr != nil {
		t.Fatalf("failed to create .context: %v", mkdirErr)
	}
	rcPath := filepath.Join(tmpDir, ".ctxrc")
	rcBody := "statusline:\n  show_cost: false\n"
	if writeErr := os.WriteFile(rcPath, []byte(rcBody), 0o600); writeErr != nil {
		t.Fatalf("failed to write .ctxrc: %v", writeErr)
	}

	origDir, _ := os.Getwd()
	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatalf("failed to chdir: %v", chdirErr)
	}
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()
	defer rc.Reset()

	line := render(t, `{
		"model": {"display_name": "Opus"},
		"cost": {"total_cost_usd": 9.99}
	}`)
	if strings.Contains(line, "$") {
		t.Errorf("line %q renders cost despite show_cost: false", line)
	}
	if !strings.Contains(line, "Opus") {
		t.Errorf("line %q lost non-cost segments", line)
	}
}
