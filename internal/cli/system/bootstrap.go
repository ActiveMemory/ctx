//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// bootstrapCmd returns the "ctx system bootstrap" subcommand.
func bootstrapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Print context location for AI agents",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runBootstrap(cmd)
		},
	}
	cmd.Flags().Bool("json", false, "Output in JSON format")
	cmd.Flags().BoolP("quiet", "q", false, "Output only the context directory path")
	return cmd
}

// bootstrapRules are the standard rules emitted by the bootstrap command.
var bootstrapRules = []string{
	"Use context_dir above for ALL file reads/writes",
	"Never say \"I don't have memory\" — context IS your memory",
	"Read files silently, present as recall (not search)",
	"Persist learnings/decisions before session ends",
	"Run `ctx agent` for content summaries",
	"Run `ctx status` for context health",
}

// bootstrapNextSteps tells the agent what to do immediately after bootstrap.
var bootstrapNextSteps = []string{
	"Read AGENT_PLAYBOOK.md from the context directory",
	"Run `ctx agent --budget 4000` for a content summary",
}

func runBootstrap(cmd *cobra.Command) error {
	dir := rc.ContextDir()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("context directory not found: %s — run 'ctx init'", dir)
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	if quiet {
		cmd.Println(dir)
		return nil
	}

	files := listContextFiles(dir)

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		return outputBootstrapJSON(cmd, dir, files)
	}
	outputBootstrapText(cmd, dir, files)
	return nil
}

func outputBootstrapText(cmd *cobra.Command, dir string, files []string) {
	cmd.Println("ctx bootstrap")
	cmd.Println("=============")
	cmd.Println()
	cmd.Println("context_dir: " + dir)
	cmd.Println()
	cmd.Println("Files:")
	cmd.Println(wrapFileList(files, 55, "  "))
	cmd.Println()
	cmd.Println("Rules:")
	for i, r := range bootstrapRules {
		cmd.Println(fmt.Sprintf("  %d. %s", i+1, r))
	}
	cmd.Println()
	cmd.Println("Next steps:")
	for i, s := range bootstrapNextSteps {
		cmd.Println(fmt.Sprintf("  %d. %s", i+1, s))
	}

	if w := pluginWarning(); w != "" {
		cmd.Println()
		cmd.Println("Warning: " + w)
	}
}

func outputBootstrapJSON(cmd *cobra.Command, dir string, files []string) error {
	type jsonOutput struct {
		ContextDir string   `json:"context_dir"`
		Files      []string `json:"files"`
		Rules      []string `json:"rules"`
		NextSteps  []string `json:"next_steps"`
		Warnings   []string `json:"warnings,omitempty"`
	}

	out := jsonOutput{
		ContextDir: dir,
		Files:      files,
		Rules:      bootstrapRules,
		NextSteps:  bootstrapNextSteps,
	}

	if w := pluginWarning(); w != "" {
		out.Warnings = []string{w}
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// pluginWarning returns a warning string if the ctx plugin is installed
// but not enabled in either global or local settings. Returns empty string
// if no warning is needed.
func pluginWarning() string {
	if !initialize.PluginInstalled() {
		return ""
	}
	if initialize.PluginEnabledGlobally() || initialize.PluginEnabledLocally() {
		return ""
	}
	return "ctx plugin is installed but not enabled. " +
		"Run 'ctx init' to auto-enable, or add " +
		"{\"enabledPlugins\": {\"ctx@activememory-ctx\": true}} " +
		"to ~/.claude/settings.json"
}

// listContextFiles reads the given directory and returns sorted .md filenames.
func listContextFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), ".md") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files
}

// wrapFileList formats file names as a comma-separated list, wrapping lines
// at approximately maxWidth characters. Continuation lines are prefixed with
// the given indent string.
func wrapFileList(files []string, maxWidth int, indent string) string {
	if len(files) == 0 {
		return indent + "(none)"
	}

	var lines []string
	current := indent

	for i, f := range files {
		entry := f
		if i < len(files)-1 {
			entry += ","
		}

		switch {
		case current == indent:
			// First entry on this line — always add it.
			current += entry
		case len(current)+1+len(entry) > maxWidth:
			// Would exceed width — start a new line.
			lines = append(lines, current)
			current = indent + entry
		default:
			current += " " + entry
		}
	}
	lines = append(lines, current)
	return strings.Join(lines, "\n")
}
