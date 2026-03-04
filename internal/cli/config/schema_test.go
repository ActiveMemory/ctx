//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestSchemaCmd_OutputsJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)

	schema := schemaCmd()
	schema.SetOut(buf)
	schema.SetArgs([]string{})

	if err := schema.Execute(); err != nil {
		t.Fatalf("schema command failed: %v", err)
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
