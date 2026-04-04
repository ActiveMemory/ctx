//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for config subcommands.
const (
	// UseConfig is the cobra Use string for the config command.
	UseConfig = "config"
	// UseConfigSchema is the cobra Use string for the config schema command.
	UseConfigSchema = "schema"
	// UseConfigStatus is the cobra Use string for the config status command.
	UseConfigStatus = "status"
	// UseConfigSwitch is the cobra Use string for the config switch command.
	UseConfigSwitch = "switch [dev|base]"
)

// DescKeys for config subcommands.
const (
	// DescKeyConfig is the description key for the config command.
	DescKeyConfig = "config"
	// DescKeyConfigSchema is the description key for the config schema command.
	DescKeyConfigSchema = "config.schema"
	// DescKeyConfigStatus is the description key for the config status command.
	DescKeyConfigStatus = "config.status"
	// DescKeyConfigSwitch is the description key for the config switch command.
	DescKeyConfigSwitch = "config.switch"
)
