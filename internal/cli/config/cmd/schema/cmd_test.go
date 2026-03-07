//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"bytes"
	"strings"
	"testing"
)

func TestSchemaCmd_OutputsJSON(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := Cmd()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	if execErr := cmd.Execute(); execErr != nil {
		t.Fatalf("schema command failed: %v", execErr)
	}

	out := buf.String()
	if !strings.Contains(out, "$schema") {
		t.Error("output should contain $schema")
	}
	if !strings.Contains(out, "additionalProperties") {
		t.Error("output should contain additionalProperties")
	}
	if !strings.Contains(out, "ctx.ist") {
		t.Error("output should contain ctx.ist $id")
	}
}
