//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"
)

// BootstrapJSONOutput is the JSON output structure for the bootstrap command.
type BootstrapJSONOutput struct {
	ContextDir string   `json:"context_dir"`
	Files      []string `json:"files"`
	Rules      []string `json:"rules"`
	NextSteps  []string `json:"next_steps"`
	Warnings   []string `json:"warnings,omitempty"`
}

// BootstrapText prints the human-readable bootstrap output to stdout.
//
// Parameters:
//   - cmd: Cobra command whose stdout stream receives the output.
//   - dir: absolute path to the context directory.
//   - fileList: pre-formatted, wrapped file list string.
//   - rules: ordered rule strings (numbered automatically).
//   - nextSteps: ordered next-step strings (numbered automatically).
//   - warning: optional warning string (empty string skips).
func BootstrapText(cmd *cobra.Command, dir string, fileList string, rules []string, nextSteps []string, warning string) {
	cmd.Println(config.TplBootstrapTitle)
	cmd.Println(config.TplBootstrapSep)
	cmd.Println()
	cmd.Println(fmt.Sprintf(config.TplBootstrapDir, dir))
	cmd.Println()
	cmd.Println(config.TplBootstrapFiles)
	cmd.Println(fileList)
	cmd.Println()
	cmd.Println(config.TplBootstrapRules)
	for i, r := range rules {
		cmd.Println(fmt.Sprintf(config.TplBootstrapNumbered, i+1, r))
	}
	cmd.Println()
	cmd.Println(config.TplBootstrapNextSteps)
	for i, s := range nextSteps {
		cmd.Println(fmt.Sprintf(config.TplBootstrapNumbered, i+1, s))
	}

	if warning != "" {
		cmd.Println()
		cmd.Println(fmt.Sprintf(config.TplBootstrapWarning, warning))
	}
}

// BootstrapJSON prints the JSON bootstrap output to stdout.
//
// Parameters:
//   - cmd: Cobra command whose stdout stream receives the output.
//   - dir: absolute path to the context directory.
//   - files: list of context file names.
//   - rules: list of rule strings.
//   - nextSteps: list of next-step strings.
//   - warning: optional warning string (empty string omits warnings).
//
// Returns:
//   - error: non-nil if JSON encoding fails.
func BootstrapJSON(cmd *cobra.Command, dir string, files []string, rules []string, nextSteps []string, warning string) error {
	out := BootstrapJSONOutput{
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
	return enc.Encode(out)
}
