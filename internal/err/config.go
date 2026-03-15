//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"
)

// UnknownProfile returns an error for an unrecognized config profile name.
//
// Parameters:
//   - name: the profile name that was not recognized
//
// Returns:
//   - error: "unknown profile <name>: must be dev, base, or prod"
func UnknownProfile(name string) error {
	return fmt.Errorf("unknown profile %q: must be dev, base, or prod", name)
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
	return fmt.Errorf("read %s: %w", name, cause)
}

// UnknownFormat returns an error for an unsupported output format.
//
// Parameters:
//   - format: the format string that was not recognized
//   - supported: list of valid formats
//
// Returns:
//   - error: "unknown format <format> (supported: <list>)"
func UnknownFormat(format, supported string) error {
	return fmt.Errorf("unknown format %q (supported: %s)", format, supported)
}

// UnknownProjectType returns an error for an unsupported project type.
//
// Parameters:
//   - projType: the type string that was not recognized
//   - supported: list of valid types
//
// Returns:
//   - error: "unknown project type <type> (supported: <list>)"
func UnknownProjectType(projType, supported string) error {
	return fmt.Errorf("unknown project type %q (supported: %s)", projType, supported)
}

// InvalidTool returns an error for an unsupported AI tool name.
//
// Parameters:
//   - tool: the tool name that was not recognized
//
// Returns:
//   - error: "invalid tool <tool>: must be claude, aider, or generic"
func InvalidTool(tool string) error {
	return fmt.Errorf("invalid tool %q: must be claude, aider, or generic", tool)
}

// UnsupportedTool returns an error for an unrecognized AI tool name.
//
// Parameters:
//   - tool: the tool name that was not recognized
//
// Returns:
//   - error: "unsupported tool: <tool>"
func UnsupportedTool(tool string) error {
	return fmt.Errorf("unsupported tool: %s", tool)
}

// UnknownUpdateType returns an error for an unrecognized context update type.
//
// Parameters:
//   - typeName: the update type that was not recognized.
//
// Returns:
//   - error: "unknown update type: <typeName>"
func UnknownUpdateType(typeName string) error {
	return fmt.Errorf("unknown update type: %s", typeName)
}

// SettingsNotFound returns an error when settings.local.json is missing.
//
// Returns:
//   - error: "no .claude/settings.local.json found"
func SettingsNotFound() error {
	return errors.New("no .claude/settings.local.json found")
}

// GoldenNotFound returns an error when settings.golden.json is missing.
//
// Returns:
//   - error: advises the user to run 'ctx permissions snapshot' first
func GoldenNotFound() error {
	return errors.New(
		"no .claude/settings.golden.json found — run 'ctx permissions snapshot' first",
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
	return fmt.Errorf("read embedded schema: %w", cause)
}

// MarshalSettings wraps a failure to marshal settings JSON.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "failed to marshal settings: <cause>"
func MarshalSettings(cause error) error {
	return fmt.Errorf("failed to marshal settings: %w", cause)
}

// MarshalPlugins wraps a failure to marshal enabledPlugins JSON.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "failed to marshal enabledPlugins: <cause>"
func MarshalPlugins(cause error) error {
	return fmt.Errorf("failed to marshal enabledPlugins: %w", cause)
}
