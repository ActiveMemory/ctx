//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/skill"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// vscodeSurface mirrors editors/vscode/src/ctx-cli-surface.json: every
// ctx invocation and skill the VS Code extension dispatches, recorded by
// the extension's own test suite (editors/vscode/src/commandParity.test.ts;
// refresh with `npx vitest run -u` in editors/vscode).
type vscodeSurface struct {
	Skills      []string   `json:"skills"`
	Invocations [][]string `json:"invocations"`
}

// entryNouns are the commands whose `add` validates required fields
// (provenance, --section, decision and learning structure) inside RunE,
// where parsing alone cannot see them.
var entryNouns = []string{"convention", "decision", "learning", "task"}

func loadVSCodeSurface(t *testing.T) vscodeSurface {
	t.Helper()
	path := filepath.Join(
		"..", "..", "editors", "vscode", "src", "ctx-cli-surface.json",
	)
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read %s: %v", path, readErr)
	}
	var surface vscodeSurface
	if jsonErr := json.Unmarshal(data, &surface); jsonErr != nil {
		t.Fatalf("parse %s: %v", path, jsonErr)
	}
	if len(surface.Invocations) == 0 || len(surface.Skills) == 0 {
		t.Fatalf("%s lists no invocations or no skills", path)
	}
	return surface
}

// resolveInvocation finds argv's command in a fresh ctx command tree and
// parses the rest of argv exactly as cobra does before running it.
//
// Parameters:
//   - argv: ctx arguments, without the binary name
//
// Returns:
//   - *cobra.Command: The command argv runs
//   - []string: Its positional arguments
//   - error: Why the CLI would reject argv, if it would
func resolveInvocation(argv []string) (*cobra.Command, []string, error) {
	root := Initialize(RootCmd())
	cmd, rest, findErr := root.Find(argv)
	switch {
	case findErr != nil:
		return nil, nil, findErr
	case cmd == root:
		return nil, nil, fmt.Errorf("no ctx command in %q", argv)
	case cmd.Deprecated != "":
		return nil, nil, fmt.Errorf("%q is deprecated", cmd.CommandPath())
	case !cmd.Runnable():
		return nil, nil, fmt.Errorf("%q only prints help", cmd.CommandPath())
	}
	if parseErr := cmd.ParseFlags(rest); parseErr != nil {
		return nil, nil, parseErr
	}
	args := cmd.Flags().Args()
	// Without an Args validator cobra accepts any positional arguments.
	// A command group then reads a leftover word as a subcommand it does
	// not have (`ctx system stats` prints an unknown-subcommand notice
	// and still exits 0); a command whose Use line names no arguments
	// silently ignores them.
	if cmd.Args == nil && len(args) > 0 {
		if cmd.HasSubCommands() {
			return nil, nil, fmt.Errorf(
				"%q has no subcommand %q", cmd.CommandPath(), args[0],
			)
		}
		if !strings.Contains(cmd.Use, " ") {
			return nil, nil, fmt.Errorf(
				"%q takes no arguments, got %q", cmd.CommandPath(), args,
			)
		}
	}
	if argsErr := cmd.ValidateArgs(args); argsErr != nil {
		return nil, nil, argsErr
	}
	if reqErr := cmd.ValidateRequiredFlags(); reqErr != nil {
		return nil, nil, reqErr
	}
	if groupErr := cmd.ValidateFlagGroups(); groupErr != nil {
		return nil, nil, groupErr
	}
	return cmd, args, nil
}

// TestVSCodeExtensionSurface fails when the VS Code extension dispatches
// a ctx invocation this CLI would reject, so renaming or removing a
// command cannot silently strand an @ctx chat command.
func TestVSCodeExtensionSurface(t *testing.T) {
	for _, argv := range loadVSCodeSurface(t).Invocations {
		t.Run(strings.Join(argv, " "), func(t *testing.T) {
			if _, _, err := resolveInvocation(argv); err != nil {
				t.Errorf("ctx %s: %v", strings.Join(argv, " "), err)
			}
		})
	}
}

// TestVSCodeExtensionEntryAdds runs the extension's `ctx <noun> add`
// invocations in a scratch project, because their required fields are
// validated at run time rather than by flag parsing.
func TestVSCodeExtensionEntryAdds(t *testing.T) {
	surface := loadVSCodeSurface(t)
	testctx.Declare(t, t.TempDir())
	initCmd := initialize.Cmd()
	initCmd.SetArgs([]string{})
	initCmd.SetOut(io.Discard)
	initCmd.SetErr(io.Discard)
	if initErr := initCmd.Execute(); initErr != nil {
		t.Fatalf("init: %v", initErr)
	}

	ran := 0
	for _, argv := range surface.Invocations {
		if len(argv) < 2 || argv[1] != "add" ||
			!slices.Contains(entryNouns, argv[0]) {
			continue
		}
		ran++
		t.Run(strings.Join(argv, " "), func(t *testing.T) {
			cmd, args, resolveErr := resolveInvocation(argv)
			if resolveErr != nil {
				t.Fatalf("resolve: %v", resolveErr)
			}
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			if runErr := cmd.RunE(cmd, args); runErr != nil {
				t.Errorf("ctx %s: %v", strings.Join(argv, " "), runErr)
			}
		})
	}
	if ran == 0 {
		t.Fatal("surface lists no entry add invocations")
	}
}

// TestVSCodeExtensionSkills fails when a skill-backed @ctx command names
// a skill that no longer ships.
func TestVSCodeExtensionSkills(t *testing.T) {
	names, listErr := skill.List()
	if listErr != nil {
		t.Fatalf("list skills: %v", listErr)
	}
	for _, name := range loadVSCodeSurface(t).Skills {
		if !slices.Contains(names, name) {
			t.Errorf("skill %q is not in internal/assets/claude/skills", name)
		}
	}
}
