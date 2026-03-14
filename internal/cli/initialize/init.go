//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"github.com/spf13/cobra"

	initroot "github.com/ActiveMemory/ctx/internal/cli/initialize/cmd/root"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core"
)

// PluginInstalled reports whether the ctx plugin is registered in
// ~/.claude/plugins/installed_plugins.json.
// Re-exported from the core subpackage for backward compatibility.
var PluginInstalled = core.PluginInstalled

// PluginEnabledGlobally reports whether the ctx plugin is enabled in
// ~/.claude/settings.json.
// Re-exported from the core subpackage for backward compatibility.
var PluginEnabledGlobally = core.PluginEnabledGlobally

// PluginEnabledLocally reports whether the ctx plugin is enabled in
// .claude/settings.local.json in the current project.
// Re-exported from the core subpackage for backward compatibility.
var PluginEnabledLocally = core.PluginEnabledLocally

// Cmd returns the "ctx init" command for initializing a .context/ directory.
func Cmd() *cobra.Command {
	return initroot.Cmd()
}
