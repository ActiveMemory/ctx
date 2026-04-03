//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for trigger subcommands.
const (
	UseTrigger        = "trigger"
	UseTriggerAdd     = "add <trigger-type> <name>"
	UseTriggerList    = "list"
	UseTriggerTest    = "test <trigger-type>"
	UseTriggerEnable  = "enable <name>"
	UseTriggerDisable = "disable <name>"
)

// DescKeys for trigger subcommands.
const (
	DescKeyTrigger        = "trigger"
	DescKeyTriggerAdd     = "trigger.add"
	DescKeyTriggerList    = "trigger.list"
	DescKeyTriggerTest    = "trigger.test"
	DescKeyTriggerEnable  = "trigger.enable"
	DescKeyTriggerDisable = "trigger.disable"
)
