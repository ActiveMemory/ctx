//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flag

// DescKeys for init command flags.
const (
	// DescKeyInitializeReset is the description key for the initialize reset
	// flag (replaces the retired --force flag; reset enumerates the existing
	// populated files, backs them up, and only proceeds on interactive
	// y/N confirmation).
	DescKeyInitializeReset = "initialize.reset"
	// DescKeyInitializeMerge is the description key for the initialize merge flag.
	DescKeyInitializeMerge = "initialize.merge"
	// DescKeyInitializeMinimal is the description key for the initialize minimal
	// flag.
	DescKeyInitializeMinimal = "initialize.minimal"
	// DescKeyInitializeNoPluginEnable is the description key for the initialize
	// no plugin enable flag.
	DescKeyInitializeNoPluginEnable = "initialize.no-plugin-enable"
	// DescKeyInitializeNoSteeringInit is the description key for the initialize
	// no steering init flag.
	DescKeyInitializeNoSteeringInit = "initialize.no-steering-init"
)
