//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/msg"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// TemplateVars prints a formatted template variables line.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - vars: pre-formatted template variables text
func TemplateVars(cmd *cobra.Command, vars string) {
	if cmd == nil {
		return
	}
	cmd.Println(vars)
}

// CtxSpecificWarning prints the ctx-specific category warning
// followed by an empty line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func CtxSpecificWarning(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyMessageCtxSpecificWarning))
	cmd.Println()
}

// OverrideCreated prints the override file creation confirmation.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: absolute path to the created override file
func OverrideCreated(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyMessageOverrideCreated), path))
}

// EditHint prints the edit hint after override creation. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func EditHint(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyMessageEditHint))
}

// SourceOverride prints the override source header with the file path.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: path to the override file
func SourceOverride(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyMessageSourceOverride), path))
}

// SourceDefault prints the default source header. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func SourceDefault(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyMessageSourceDefault))
}

// ContentBlock prints raw content with a leading blank line. If the
// content does not end with a newline, an extra newline is appended.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - data: raw content bytes to display
func ContentBlock(cmd *cobra.Command, data []byte) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Print(string(data))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		cmd.Println()
	}
}

// NoOverride prints a message indicating no override exists.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hook: hook name
//   - variant: variant name
func NoOverride(cmd *cobra.Command, hook, variant string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyMessageNoOverride),
		hook, variant))
}

// OverrideRemoved prints the override removal confirmation.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hook: hook name
//   - variant: variant name
func OverrideRemoved(cmd *cobra.Command, hook, variant string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyMessageOverrideRemoved),
		hook, variant))
}

// ListHeader prints the message list table header and separator.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func ListHeader(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(msg.MessageListFormat,
		desc.Text(text.DescKeyMessageListHeaderHook),
		desc.Text(text.DescKeyMessageListHeaderVariant),
		desc.Text(text.DescKeyMessageListHeaderCategory),
		desc.Text(text.DescKeyMessageListHeaderOverride)))
	cmd.Println(fmt.Sprintf(msg.MessageListFormat,
		strings.Repeat(token.LineHorizontal, msg.MessageSepHook),
		strings.Repeat(token.LineHorizontal, msg.MessageSepVariant),
		strings.Repeat(token.LineHorizontal, msg.MessageSepCategory),
		strings.Repeat(token.LineHorizontal, msg.MessageSepOverride)))
}

// ListRow prints a single message list table row. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hook: hook name
//   - variant: variant name
//   - category: message category
//   - hasOverride: whether a custom override exists
func ListRow(
	cmd *cobra.Command,
	hook, variant, category string,
	hasOverride bool,
) {
	if cmd == nil {
		return
	}
	override := ""
	if hasOverride {
		override = desc.Text(text.DescKeyMessageOverrideLabel)
	}
	cmd.Println(fmt.Sprintf(
		msg.MessageListFormat,
		hook, variant, category, override))
}
