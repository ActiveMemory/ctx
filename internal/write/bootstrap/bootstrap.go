//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/spf13/cobra"
)

// Text prints the human-readable bootstrap output to stdout.
//
// Parameters:
//   - cmd: Cobra command whose stdout stream receives the output.
//   - dir: absolute path to the context directory.
//   - fileList: pre-formatted, wrapped file list string.
//   - rules: ordered rule strings (numbered automatically).
//   - nextSteps: ordered next-step strings (numbered automatically).
//   - warning: optional warning string (empty string skips).
func Text(
	cmd *cobra.Command,
	dir string,
	fileList string,
	rules []string,
	nextSteps []string,
	warning string,
) {
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteBootstrapTitle))
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteBootstrapSep))
	cmd.Println()
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.TextDescKeyWriteBootstrapDir), dir))
	cmd.Println()
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteBootstrapFiles))
	cmd.Println(fileList)
	cmd.Println()
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteBootstrapRules))
	for i, r := range rules {
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(text.TextDescKeyWriteBootstrapNumbered), i+1, r,
		))
	}
	cmd.Println()
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteBootstrapNextSteps))
	for i, s := range nextSteps {
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(text.TextDescKeyWriteBootstrapNumbered), i+1, s,
		))
	}
	if warning != "" {
		cmd.Println()
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(text.TextDescKeyWriteBootstrapWarning), warning,
		))
	}
}

// JSON prints the JSON bootstrap output to stdout.
//
// If encoding fails, a structured JSON error object is printed instead.
//
// Parameters:
//   - cmd: Cobra command whose stdout stream receives the output.
//   - dir: absolute path to the context directory.
//   - files: list of context file names.
//   - rules: list of rule strings.
//   - nextSteps: list of next-step strings.
//   - warning: optional warning string (empty string omits warnings).
func JSON(
	cmd *cobra.Command,
	dir string,
	files []string,
	rules []string,
	nextSteps []string,
	warning string,
) {
	out := entity.BootstrapOutput{
		ContextDir: dir,
		Files:      files,
		Rules:      rules,
		NextSteps:  nextSteps,
	}

	if warning != "" {
		out.Warnings = []string{warning}
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if encodeErr := enc.Encode(out); encodeErr != nil {
		cmd.PrintErrln(fmt.Sprintf(
			`{"error": "json encode: %v"}`, encodeErr,
		))
	}
}
