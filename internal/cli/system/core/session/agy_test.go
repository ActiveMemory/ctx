//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/agy"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func writeInitializedProject(t *testing.T, tempDir string) string {
	t.Helper()
	ctxDir := filepath.Join(tempDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0o750); mkErr != nil {
		t.Fatalf("MkdirAll(%s): %v", ctxDir, mkErr)
	}
	for _, f := range cfgCtx.FilesRequired {
		p := filepath.Join(ctxDir, f)
		if wErr := os.WriteFile(p, []byte("# "+f+"\n"), 0o600); wErr != nil {
			t.Fatalf("WriteFile(%s): %v", p, wErr)
		}
	}
	return ctxDir
}

func TestAgyModeState(t *testing.T) {
	// Ensure clean initial state
	SetAgyMode(false, "", "", "")
	if Agy() {
		t.Fatal("expected Agy() to be false initially")
	}

	SetAgyMode(true, "sess-123", "/path/to/transcript.jsonl", "gemini-pro")
	if !Agy() {
		t.Error("expected Agy() to be true")
	}
	if got := AgyTranscriptPath(); got != "/path/to/transcript.jsonl" {
		t.Errorf("AgyTranscriptPath() = %q, want %q", got, "/path/to/transcript.jsonl")
	}
	if got := AgyModelName(); got != "gemini-pro" {
		t.Errorf("AgyModelName() = %q, want %q", got, "gemini-pro")
	}

	if AgyOutputWritten() {
		t.Error("expected AgyOutputWritten() to be false")
	}
	SetAgyOutputWritten(true)
	if !AgyOutputWritten() {
		t.Error("expected AgyOutputWritten() to be true after set")
	}

	// Reset
	SetAgyMode(false, "", "", "")
	if Agy() {
		t.Error("expected Agy() to be false after reset")
	}
}

func TestFormatAgyInjectSteps(t *testing.T) {
	// Empty slice
	if got := FormatAgyInjectSteps(nil); got != agy.EmptyJSON {
		t.Errorf("FormatAgyInjectSteps(nil) = %q, want %q", got, agy.EmptyJSON)
	}
	if got := FormatAgyInjectSteps([]string{}); got != agy.EmptyJSON {
		t.Errorf("FormatAgyInjectSteps([]) = %q, want %q", got, agy.EmptyJSON)
	}

	// Single step
	res := FormatAgyInjectSteps([]string{"hello world"})
	var parsed agyInjectStepsResponse
	if err := json.Unmarshal([]byte(res), &parsed); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if len(parsed.InjectSteps) != 1 {
		t.Fatalf("len(InjectSteps) = %d, want 1", len(parsed.InjectSteps))
	}
	if got := parsed.InjectSteps[0].SystemMessage.SystemMessage; got != "hello world" {
		t.Errorf("SystemMessage = %q, want %q", got, "hello world")
	}

	// Multiple steps
	resMulti := FormatAgyInjectSteps([]string{"step 1", "step 2"})
	var parsedMulti agyInjectStepsResponse
	if err := json.Unmarshal([]byte(resMulti), &parsedMulti); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if len(parsedMulti.InjectSteps) != 2 {
		t.Fatalf("len(InjectSteps) = %d, want 2", len(parsedMulti.InjectSteps))
	}
}

func TestFormatAgyDecision(t *testing.T) {
	allow := FormatAgyDecision(agy.DecisionAllow, "")
	var parsedAllow agyDecisionResponse
	if err := json.Unmarshal([]byte(allow), &parsedAllow); err != nil {
		t.Fatalf("unmarshal allow JSON: %v", err)
	}
	if parsedAllow.Decision != "allow" || parsedAllow.Reason != "" {
		t.Errorf("unexpected allow response: %+v", parsedAllow)
	}

	deny := FormatAgyDecision(agy.DecisionDeny, "access denied")
	var parsedDeny agyDecisionResponse
	if err := json.Unmarshal([]byte(deny), &parsedDeny); err != nil {
		t.Fatalf("unmarshal deny JSON: %v", err)
	}
	if parsedDeny.Decision != "deny" || parsedDeny.Reason != "access denied" {
		t.Errorf("unexpected deny response: %+v", parsedDeny)
	}
}

func TestQueueAndDrainPendingAgy(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origWd)
		rc.Reset()
	}()
	_ = os.Chdir(tmpDir)
	writeInitializedProject(t, tmpDir)
	rc.Reset()

	sessionID := "test-session-42"
	SetAgyMode(true, sessionID, "", "")
	defer SetAgyMode(false, "", "", "")

	// Initially empty
	drained := DrainPendingAgy()
	if len(drained) != 0 {
		t.Errorf("expected empty drain, got: %v", drained)
	}

	// Queue two nudges
	QueuePendingAgy("first nudge")
	QueuePendingAgy("second nudge\nmultiline")

	// Drain
	drained = DrainPendingAgy()
	if len(drained) != 2 {
		t.Fatalf("DrainPendingAgy len = %d, want 2 (%+v)", len(drained), drained)
	}
	if drained[0] != "first nudge" {
		t.Errorf("drained[0] = %q, want %q", drained[0], "first nudge")
	}
	if drained[1] != "second nudge\nmultiline" {
		t.Errorf("drained[1] = %q, want %q", drained[1], "second nudge\nmultiline")
	}

	// Next drain should be empty
	drainedAgain := DrainPendingAgy()
	if len(drainedAgain) != 0 {
		t.Errorf("expected empty second drain, got: %v", drainedAgain)
	}
}

func TestReadInputAgyNormalization(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origWd)
		SetAgyMode(false, "", "", "")
	}()

	agyJSON := `{
		"conversationId": "conv-999",
		"workspacePaths": ["` + filepath.ToSlash(tmpDir) + `"],
		"transcriptPath": "/brain/conv-999/transcript.jsonl",
		"modelName": "gemini-2.0-flash",
		"lastUserInput": "ctx task list",
		"toolCall": {
			"name": "run_command",
			"args": {
				"CommandLine": "git commit -m 'test'"
			}
		}
	}`

	input := ReadInput(strings.NewReader(agyJSON))

	if !input.Agy() {
		t.Error("expected input.Agy() to be true")
	}
	if input.SessionID != "conv-999" {
		t.Errorf("SessionID = %q, want %q", input.SessionID, "conv-999")
	}
	if input.Prompt != "ctx task list" {
		t.Errorf("Prompt = %q, want %q", input.Prompt, "ctx task list")
	}
	if input.ToolName != "run_command" {
		t.Errorf("ToolName = %q, want %q", input.ToolName, "run_command")
	}
	if input.ToolInput.Command != "git commit -m 'test'" {
		t.Errorf("Command = %q, want %q", input.ToolInput.Command, "git commit -m 'test'")
	}
	if input.CWD != filepath.ToSlash(tmpDir) {
		t.Errorf("CWD = %q, want %q", input.CWD, filepath.ToSlash(tmpDir))
	}
	if !Agy() {
		t.Error("expected session.Agy() to be true")
	}
	if input.SessionID != "conv-999" {
		t.Errorf("input.SessionID = %q, want %q", input.SessionID, "conv-999")
	}

	// Now parse non-Agy input -> should reset Agy mode
	claudeJSON := `{"session_id": "claude-1", "tool_input": {"command": "ls"}}`
	nonAgyInput := ReadInput(strings.NewReader(claudeJSON))
	if nonAgyInput.Agy() {
		t.Error("expected nonAgyInput.Agy() to be false")
	}
	if Agy() {
		t.Error("expected session.Agy() to be reset to false")
	}
}

func TestAgyStepWithBufferedProvenance(t *testing.T) {
	defer SetAgyMode(false, "", "", "")
	SetAgyMode(true, "sess-prov", "", "")

	SetAgyProvenance("[ctx: sess branch 1234567]")
	AddAgyStep("reminder box content")

	steps := DrainAgySteps()
	if len(steps) != 1 {
		t.Fatalf("len(steps) = %d, want 1", len(steps))
	}
	if !strings.Contains(steps[0], "[ctx: sess branch 1234567]") {
		t.Errorf("step should contain buffered provenance, got: %q", steps[0])
	}
	if !strings.Contains(steps[0], "reminder box content") {
		t.Errorf("step should contain reminder content, got: %q", steps[0])
	}
}

func TestHookInputAgyEnvOverride(t *testing.T) {
	h := entity.HookInput{}
	if h.Agy() {
		t.Error("expected h.Agy() to be false without env or conversationId")
	}

	t.Setenv(agy.EnvHookTool, hook.ToolAgy)
	if !h.Agy() {
		t.Error("expected h.Agy() to be true with CTX_HOOK_TOOL=agy")
	}

	t.Setenv(agy.EnvHookTool, hook.ToolAntigravity)
	if !h.Agy() {
		t.Error("expected h.Agy() to be true with CTX_HOOK_TOOL=antigravity")
	}
}

func TestAgyConstants(t *testing.T) {
	if agy.ToolRunCommand != "run_command" {
		t.Errorf("unexpected ToolRunCommand: %s", agy.ToolRunCommand)
	}
	if agy.ToolEditFile != "edit_file" {
		t.Errorf("unexpected ToolEditFile: %s", agy.ToolEditFile)
	}
	if agy.ToolWriteFile != "write_file" {
		t.Errorf("unexpected ToolWriteFile: %s", agy.ToolWriteFile)
	}
	if agy.ToolReadFile != "read_file" {
		t.Errorf("unexpected ToolReadFile: %s", agy.ToolReadFile)
	}
	if agy.EnvConfigDir != "AGY_CONFIG_DIR" {
		t.Errorf("unexpected EnvConfigDir: %s", agy.EnvConfigDir)
	}
	if agy.EventSessionStart != "SessionStart" {
		t.Errorf("unexpected EventSessionStart: %s", agy.EventSessionStart)
	}
	if agy.EventPreInvocation != "PreInvocation" {
		t.Errorf("unexpected EventPreInvocation: %s", agy.EventPreInvocation)
	}
	if agy.EventPostInvocation != "PostInvocation" {
		t.Errorf("unexpected EventPostInvocation: %s", agy.EventPostInvocation)
	}
	if agy.EventStop != "Stop" {
		t.Errorf("unexpected EventStop: %s", agy.EventStop)
	}
}

func TestParseLastUsageAndModel_Agy(t *testing.T) {
	tempDir := t.TempDir()
	transcriptPath := filepath.Join(tempDir, agy.FileTranscriptJSONL)

	line := `{"step_index":1,"source":"MODEL","type":"PLANNER_RESPONSE",` +
		`"status":"DONE","created_at":"2026-11-06T23:30:04Z",` +
		`"input_tokens":7306,"cache_read_tokens":120,"output_tokens":225}` + "\n"
	if err := os.WriteFile(transcriptPath, []byte(line), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	SetAgyMode(true, "sess-test", transcriptPath, "gemini-2.5-pro")
	defer SetAgyMode(false, "", "", "")

	info, err := ParseLastUsageAndModel(transcriptPath)
	if err != nil {
		t.Fatalf("ParseLastUsageAndModel: %v", err)
	}

	if info.Tokens != 7426 {
		t.Errorf("Tokens = %d, want 7426 (7306+120)", info.Tokens)
	}
	if info.Model != "gemini-2.5-pro" {
		t.Errorf("Model = %q, want %q", info.Model, "gemini-2.5-pro")
	}
}
