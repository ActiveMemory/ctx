//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for trigger subcommands.
const (
	// UseTrigger is the cobra Use string for the trigger command.
	UseTrigger = "trigger"
	// UseTriggerAdd is the cobra Use string for the trigger add command.
	UseTriggerAdd = "add <trigger-type> <name>"
	// UseTriggerList is the cobra Use string for the trigger list command.
	UseTriggerList = "list"
	// UseTriggerTest is the cobra Use string for the trigger test command.
	UseTriggerTest = "test <trigger-type>"
	// UseTriggerEnable is the cobra Use string for the trigger enable command.
	UseTriggerEnable = "enable <name>"
	// UseTriggerDisable is the cobra Use string for the trigger disable command.
	UseTriggerDisable = "disable <name>"
)

// DescKeys for trigger subcommands.
const (
	// DescKeyTrigger is the description key for the trigger command.
	DescKeyTrigger = "trigger"
	// DescKeyTriggerAdd is the description key for the trigger add command.
	DescKeyTriggerAdd = "trigger.add"
	// DescKeyTriggerList is the description key for the trigger list command.
	DescKeyTriggerList = "trigger.list"
	// DescKeyTriggerTest is the description key for the trigger test command.
	DescKeyTriggerTest = "trigger.test"
	// DescKeyTriggerEnable is the description key for the trigger enable command.
	DescKeyTriggerEnable = "trigger.enable"
	// DescKeyTriggerDisable is the description key for the trigger disable
	// command.
	DescKeyTriggerDisable = "trigger.disable"
)
