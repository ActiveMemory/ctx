//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/bootstrap"
)

// ctxBinaryName is the command word a shipped hook uses to invoke
// the ctx binary from PATH.
const ctxBinaryName = "ctx"

// subcommandToken matches a cobra subcommand name: a lowercase
// letter followed by lowercase letters, digits, and hyphens. Used
// to peel the leading command path off a hook's shell command,
// stopping at the first flag, redirection, or shell operator.
var subcommandToken = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// shippedHookFile mirrors the structure of
// internal/assets/claude/hooks/hooks.json: a top-level "hooks"
// object keyed by Claude Code event name, each mapping to a list
// of matcher groups that each carry a list of command hooks.
type shippedHookFile struct {
	Hooks map[string][]struct {
		Hooks []struct {
			Command string `json:"command"`
		} `json:"hooks"`
	} `json:"hooks"`
}

// TestShippedHooksResolveToRegisteredCommands asserts that every
// `ctx <…>` invocation wired into the shipped hooks.json resolves
// to a registered subcommand on the assembled command tree.
//
// This is the recurrence guard for the version-skew bug recorded
// in specs/hooks-wiring-guard.md: a published plugin whose
// hooks.json wired `ctx system check-anchor-drift` after the
// binary had deleted that command, so cobra dumped ~51 lines of
// `system` help (exit 0) into the agent's context on every prompt.
// A half-migrated package now fails here instead of in a session.
func TestShippedHooksResolveToRegisteredCommands(t *testing.T) {
	root := projectRoot(t)
	hooksPath := filepath.Join(
		root, "internal", "assets", "claude", "hooks", "hooks.json",
	)

	data, err := os.ReadFile(filepath.Clean(hooksPath))
	if err != nil {
		t.Fatalf("read shipped hooks.json: %v", err)
	}

	var hf shippedHookFile
	if err := json.Unmarshal(data, &hf); err != nil {
		t.Fatalf("decode shipped hooks.json: %v", err)
	}

	tree := bootstrap.Initialize(bootstrap.RootCmd())

	checked := 0
	for event, groups := range hf.Hooks {
		for _, group := range groups {
			for _, h := range group.Hooks {
				for _, path := range ctxInvocationPaths(h.Command) {
					checked++
					if token, ok := pathResolved(tree, path); !ok {
						t.Errorf(
							"%s hook %q wires `ctx %s`, but %q is not a "+
								"registered subcommand; shipped hooks must "+
								"match the binary's command tree "+
								"(see specs/hooks-wiring-guard.md)",
							event, h.Command,
							strings.Join(path, " "), token,
						)
					}
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal(
			"no `ctx` invocations found in shipped hooks.json; the " +
				"guard parsed nothing — check the asset path and format",
		)
	}
}

// ctxInvocationPaths extracts every ctx command path from a hook's
// shell command string. For each standalone `ctx` token it
// collects the following run of subcommand-shaped tokens, stopping
// at the first flag, redirection, or shell operator. A bare `ctx`
// with no subcommand yields nothing.
//
// strings.Fields does not honour shell quoting, so a literal "ctx"
// inside a quoted message (e.g. the cd guard's "cannot anchor ctx")
// could surface as a token. Such occurrences are either glued to
// punctuation (ctx}") and fail the equality check, or followed by a
// non-subcommand token and yield an empty path — neither produces a
// false invocation.
func ctxInvocationPaths(command string) [][]string {
	tokens := strings.Fields(command)
	var paths [][]string
	for i, tok := range tokens {
		if tok != ctxBinaryName {
			continue
		}
		var path []string
		for _, next := range tokens[i+1:] {
			if !subcommandToken.MatchString(next) {
				break
			}
			path = append(path, next)
		}
		if len(path) > 0 {
			paths = append(paths, path)
		}
	}
	return paths
}

// pathResolved descends the command tree token by token. It returns
// (token, false) at the first token with no matching child, or
// ("", true) when the whole path resolves. Hidden commands are
// traversed — Hidden affects only help display, not lookup.
func pathResolved(root *cobra.Command, path []string) (string, bool) {
	cur := root
	for _, name := range path {
		child := childNamed(cur, name)
		if child == nil {
			return name, false
		}
		cur = child
	}
	return "", true
}

// childNamed returns the immediate subcommand of cur whose name or
// alias equals name, or nil if none matches.
func childNamed(cur *cobra.Command, name string) *cobra.Command {
	for _, c := range cur.Commands() {
		if c.Name() == name {
			return c
		}
		for _, alias := range c.Aliases {
			if alias == name {
				return c
			}
		}
	}
	return nil
}
