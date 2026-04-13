//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// UnknownProfile returns an error for an unrecognized config profile name.
//
// Parameters:
//   - name: the profile name that was not recognized
//
// Returns:
//   - error: "unknown profile <name>: must be dev, base, or prod"
func UnknownProfile(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrConfigUnknownProfile), name,
	)
}

// ReadProfile wraps a failure to read a profile file.
//
// Parameters:
//   - name: profile filename
//   - cause: the underlying read error
//
// Returns:
//   - error: "read <name>: <cause>"
func ReadProfile(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrConfigReadProfile), name, cause,
	)
}

// InvalidTool returns an error for an unsupported AI tool name.
//
// Parameters:
//   - tool: the tool name that was not recognized
//
// Returns:
//   - error: "invalid tool <tool>: must be claude, aider, or generic"
func InvalidTool(tool string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrConfigInvalidTool), tool,
	)
}

// UnsupportedTool returns an error for an unrecognized AI tool name.
//
// Parameters:
//   - tool: the tool name that was not recognized
//
// Returns:
//   - error: "unsupported tool: <tool>"
func UnsupportedTool(tool string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrConfigUnsupportedTool), tool,
	)
}

// UnknownUpdateType returns an error for an unrecognized context update type.
//
// Parameters:
//   - typeName: the update type that was not recognized.
//
// Returns:
//   - error: "unknown update type: <typeName>"
func UnknownUpdateType(typeName string) error {
	return fmt.Errorf(desc.Text(
		text.DescKeyErrConfigUnknownUpdateType), typeName,
	)
}

// SettingsNotFound returns an error when settings.local.json is missing.
//
// Returns:
//   - error: "no .claude/settings.local.json found"
func SettingsNotFound() error {
	return errors.New(
		desc.Text(text.DescKeyErrConfigSettingsNotFound),
	)
}

// GoldenNotFound returns an error when settings.golden.json is missing.
//
// Returns:
//   - error: advises the user to run 'ctx permission snapshot' first
func GoldenNotFound() error {
	return errors.New(
		desc.Text(text.DescKeyErrConfigGoldenNotFound),
	)
}

// ReadEmbeddedSchema wraps a failure to read the embedded JSON Schema.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "read embedded schema: <cause>"
func ReadEmbeddedSchema(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrConfigReadEmbeddedSchema), cause,
	)
}

// MarshalSettings wraps a failure to marshal settings JSON.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "failed to marshal settings: <cause>"
func MarshalSettings(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrConfigMarshalSettings), cause,
	)
}

// MarshalPlugins wraps a failure to marshal enabledPlugins JSON.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "failed to marshal enabledPlugins: <cause>"
func MarshalPlugins(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrConfigMarshalPlugins), cause,
	)
}
