//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgTrigger "github.com/ActiveMemory/ctx/internal/config/trigger"
	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trigger"
	writeTrigger "github.com/ActiveMemory/ctx/internal/write/trigger"
)

// Cmd returns the "ctx hook test" subcommand.
//
// Returns:
//   - *cobra.Command: Configured test subcommand
func Cmd() *cobra.Command {
	var toolName string
	var path string

	short, long := desc.Command(cmd.DescKeyTriggerTest)

	c := &cobra.Command{
		Use:   cmd.UseTriggerTest,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return Run(c, args[0], toolName, path)
		},
	}

	flagbind.StringFlag(c, &toolName, flag.Tool, embedFlag.DescKeyTriggerTestTool)
	flagbind.StringFlag(c, &path, flag.Path, embedFlag.DescKeyTriggerTestPath)

	return c
}

// Run tests hooks for a given hook type by constructing a mock input
// and executing all enabled hooks.
//
// Parameters:
//   - c: The cobra command for output
//   - hookType: The hook type to test
//   - toolName: Optional tool name for mock input
//   - path: Optional file path for mock input
func Run(c *cobra.Command, hookType, toolName, path string) error {
	// Validate hook type.
	ht := hookType
	valid := trigger.ValidTypes()

	found := false
	for _, v := range valid {
		if v == ht {
			found = true
			break
		}
	}

	if !found {
		names := make([]string, len(valid))
		copy(names, valid)
		return errTrigger.InvalidType(hookType, strings.Join(names, token.CommaSpace))
	}

	hooksDir := rc.HooksDir()
	timeout := time.Duration(rc.HookTimeout()) * time.Second

	// Build mock input.
	params := make(map[string]any)
	if path != "" {
		params[flag.Path] = path
	}

	input := &trigger.HookInput{
		TriggerType: hookType,
		Tool:        toolName,
		Parameters:  params,
		Session: trigger.HookSession{
			ID:    cfgTrigger.MockSessionID,
			Model: cfgTrigger.MockModel,
		},
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		CtxVersion: cfgTrigger.MockVersion,
	}

	writeTrigger.TestingHeader(c, hookType)

	inputJSON, _ := json.MarshalIndent(input, "", "  ")
	writeTrigger.TestInput(c, string(inputJSON))

	agg, err := trigger.RunAll(hooksDir, ht, input, timeout)
	if err != nil {
		return err
	}

	if agg.Cancelled {
		writeTrigger.Cancelled(c, agg.Message)
		return nil
	}

	if agg.Context != "" {
		writeTrigger.ContextOutput(c, agg.Context)
	}

	if len(agg.Errors) > 0 {
		writeTrigger.ErrorsHeader(c)
		for _, e := range agg.Errors {
			writeTrigger.ErrorLine(c, e)
		}
		writeTrigger.BlankLine(c)
	}

	if agg.Context == "" && len(agg.Errors) == 0 {
		writeTrigger.NoOutput(c)
	}

	return nil
}
