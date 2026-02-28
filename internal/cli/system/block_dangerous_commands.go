//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// blockDangerousCommandsCmd returns the "ctx system block-dangerous-commands"
// command.
//
// Regex safety net for commands that the deny-list cannot express. The bulk of
// command blocking is handled by permissions.deny in settings.local.json; this
// hook catches only patterns requiring regex matching:
//   - Mid-command sudo/git-push (after &&, ||, ;)
//   - cp/mv to bin directories
//   - cp/install to ~/.local/bin
func blockDangerousCommandsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block-dangerous-commands",
		Short: "Block dangerous command patterns (regex safety net)",
		Long: `Regex safety net for commands that the deny-list cannot express.
Catches mid-command sudo, mid-command git push, and binary installs
to bin directories.

Hook event: PreToolUse (Bash)
Output: {"decision":"block","reason":"..."} or silent
Silent when: command doesn't match any dangerous pattern`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runBlockDangerousCommands(cmd, os.Stdin)
		},
	}
}

// Compiled regex patterns for dangerous command detection.
var (
	// Mid-command sudo after && || ;
	reMidSudo = regexp.MustCompile(`(;|&&|\|\|)\s*sudo\s`)
	// Mid-command git push after && || ;
	reMidGitPush = regexp.MustCompile(`(;|&&|\|\|)\s*git\s+push`)
	// cp/mv to bin directories
	reCpMvToBin = regexp.MustCompile(`(cp|mv)\s+\S+\s+(/usr/local/bin|/usr/bin|~/go/bin|~/.local/bin|/home/\S+/go/bin|/home/\S+/.local/bin)`)
	// cp/install to ~/.local/bin
	reInstallToLocalBin = regexp.MustCompile(`(cp|install)\s.*~/\.local/bin`)
)

func runBlockDangerousCommands(cmd *cobra.Command, stdin *os.File) error {
	input := readInput(stdin)
	command := input.ToolInput.Command

	if command == "" {
		return nil
	}

	var variant, fallback string

	// Mid-command sudo — after && || ; (prefix sudo caught by deny rule)
	if reMidSudo.MatchString(command) {
		variant = "mid-sudo"
		fallback = "Cannot use sudo (no password access). Use 'make build && sudo make install' manually if needed."
	}

	// Mid-command git push — after && || ; (prefix git push caught by deny rule)
	if variant == "" && reMidGitPush.MatchString(command) {
		variant = "mid-git-push"
		fallback = "git push requires explicit user approval."
	}

	// cp/mv to bin directories — agent must never install binaries
	if variant == "" && reCpMvToBin.MatchString(command) {
		variant = "cp-to-bin"
		fallback = "Agent must not copy binaries to bin directories. Ask the user to run 'sudo make install' instead."
	}

	// cp/install to ~/.local/bin — breaks PATH ctx rules
	if variant == "" && reInstallToLocalBin.MatchString(command) {
		variant = "install-to-local-bin"
		fallback = "Do not copy binaries to ~/.local/bin — this overrides the system ctx in /usr/local/bin. Use 'ctx' from PATH."
	}

	var reason string
	if variant != "" {
		reason = loadMessage("block-dangerous-commands", variant, nil, fallback)
	}

	if reason != "" {
		resp := blockResponse{
			Decision: "block",
			Reason:   reason,
		}
		data, _ := json.Marshal(resp)
		cmd.Println(string(data))
		ref := notify.NewTemplateRef("block-dangerous-commands", variant, nil)
		_ = notify.Send("relay", "block-dangerous-commands: "+reason, input.SessionID, ref)
		eventlog.Append("relay", "block-dangerous-commands: "+reason, input.SessionID, ref)
	}

	return nil
}
