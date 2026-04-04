//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// User-facing messages for hook list and test output.
const (
	// msgNoHooksFound is shown when no hooks are discovered.
	msgNoHooksFound = "No hooks found."
	// msgErrors is the section header for hook errors.
	msgErrors = "Errors:"
	// msgNoOutput is shown when hooks produce no output.
	msgNoOutput = "No output from hooks."
)

// Created prints confirmation that a hook script was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: Path to the created hook script
func Created(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerCreated), path,
	))
}

// Disabled prints confirmation that a hook was disabled.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Hook name
//   - path: Path to the hook script
func Disabled(cmd *cobra.Command, name, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerDisabled),
		name, path,
	))
}

// Enabled prints confirmation that a hook was enabled.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Hook name
//   - path: Path to the hook script
func Enabled(cmd *cobra.Command, name, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerEnabled),
		name, path,
	))
}

// TypeHeader prints a hook type section header.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hookType: The hook type name
func TypeHeader(cmd *cobra.Command, hookType string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerTypeHdr), hookType,
	))
}

// Entry prints a single hook entry with name, status, and path.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Hook name
//   - status: "enabled" or "disabled"
//   - path: Path to the hook script
func Entry(cmd *cobra.Command, name, status, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerEntry),
		name, status, path,
	))
}

// BlankLine prints a blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func BlankLine(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
}

// NoHooksFound prints a message indicating no hooks were found.
//
// Parameters:
//   - cmd: Cobra command for output
func NoHooksFound(cmd *cobra.Command) {
	cmd.Println(msgNoHooksFound)
}

// Count prints the total hook count.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of hooks
func Count(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerCount), count,
	))
}

// TestingHeader prints the header for hook testing output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hookType: The hook type being tested
func TestingHeader(cmd *cobra.Command, hookType string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerTestHdr), hookType,
	))
	cmd.Println()
}

// TestInput prints the test input JSON block.
//
// Parameters:
//   - cmd: Cobra command for output
//   - inputJSON: Pretty-printed JSON input
func TestInput(cmd *cobra.Command, inputJSON string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerTestInput), inputJSON,
	))
	cmd.Println()
}

// Cancelled prints a cancellation message from hook output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - message: The cancellation reason
func Cancelled(cmd *cobra.Command, message string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerCancelled), message,
	))
}

// ContextOutput prints context output from hook execution.
//
// Parameters:
//   - cmd: Cobra command for output
//   - context: The context string from hooks
func ContextOutput(cmd *cobra.Command, context string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerContext), context,
	))
	cmd.Println()
}

// ErrorsHeader prints the errors section header.
//
// Parameters:
//   - cmd: Cobra command for output
func ErrorsHeader(cmd *cobra.Command) {
	cmd.Println(msgErrors)
}

// ErrorLine prints a single error line.
//
// Parameters:
//   - cmd: Cobra command for output
//   - errMsg: The error message
func ErrorLine(cmd *cobra.Command, errMsg string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteTriggerErrLine), errMsg,
	))
}

// NoOutput prints a message indicating no output from hooks.
//
// Parameters:
//   - cmd: Cobra command for output
func NoOutput(cmd *cobra.Command) {
	cmd.Println(msgNoOutput)
}
