//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/agy"
)

func TestEmitAgyHelpers(t *testing.T) {
	defer coreSession.SetAgyMode(false, "", "", "")

	// 1. EmitAgyAllow
	coreSession.SetAgyMode(true, "sess-1", "", "")
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	EmitAgyAllow(cmd)
	out := strings.TrimSpace(buf.String())
	if out != `{"decision":"allow"}` {
		t.Errorf("EmitAgyAllow output = %q, want %q", out, `{"decision":"allow"}`)
	}
	// Calling again should be no-op because agyOutputWritten is true
	buf.Reset()
	EmitAgyAllow(cmd)
	if buf.Len() != 0 {
		t.Errorf("expected no-op on second EmitAgyAllow, got: %q", buf.String())
	}

	// 2. EmitAgyEmpty
	coreSession.SetAgyMode(true, "sess-2", "", "")
	buf.Reset()
	EmitAgyEmpty(cmd)
	out = strings.TrimSpace(buf.String())
	if out != "{}" {
		t.Errorf("EmitAgyEmpty output = %q, want %q", out, "{}")
	}

	// 3. FlushAgySteps with steps
	coreSession.SetAgyMode(true, "sess-3", "", "")
	buf.Reset()
	coreSession.AddAgyStep("step one")
	coreSession.AddAgyStep("step two")
	FlushAgySteps(cmd)
	out = strings.TrimSpace(buf.String())
	if !strings.Contains(out, "step one") ||
		!strings.Contains(out, "step two") {
		t.Errorf("FlushAgySteps output = %q, expected steps", out)
	}

	// 4. FlushAgySteps with no steps -> emits "{}"
	coreSession.SetAgyMode(true, "sess-4", "", "")
	buf.Reset()
	FlushAgySteps(cmd)
	out = strings.TrimSpace(buf.String())
	if out != "{}" {
		t.Errorf("FlushAgySteps empty output = %q, want %q", out, "{}")
	}

	// 5. EmitAgyDecision
	coreSession.SetAgyMode(true, "sess-5", "", "")
	buf.Reset()
	EmitAgyDecision(cmd, agy.DecisionDeny, "blocked reason")
	out = strings.TrimSpace(buf.String())
	if !strings.Contains(out, `"decision":"deny"`) ||
		!strings.Contains(out, `"reason":"blocked reason"`) {
		t.Errorf("EmitAgyDecision output = %q", out)
	}

	// 6. EmitAgyInjectSteps
	coreSession.SetAgyMode(true, "sess-6", "", "")
	buf.Reset()
	EmitAgyInjectSteps(cmd, []string{"injected content"})
	out = strings.TrimSpace(buf.String())
	if !strings.Contains(out, "injected content") {
		t.Errorf("EmitAgyInjectSteps output = %q", out)
	}
}
