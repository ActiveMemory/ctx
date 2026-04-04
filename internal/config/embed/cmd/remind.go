//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for remind subcommands.
const (
	// UseRemindAdd is the cobra Use string for the remind add command.
	UseRemindAdd = "add TEXT"
	// UseRemindDismiss is the cobra Use string for the remind dismiss command.
	UseRemindDismiss = "dismiss [ID]"
	// UseRemindDismissAlias is the cobra Use string for the remind dismiss alias
	// command.
	UseRemindDismissAlias = "rm"
	// UseRemindList is the cobra Use string for the remind list command.
	UseRemindList = "list"
	// UseRemindListAlias is the cobra Use string for the remind list alias
	// command.
	UseRemindListAlias = "ls"
)

// DescKeys for remind subcommands.
const (
	// DescKeyRemind is the description key for the remind command.
	DescKeyRemind = "remind"
	// DescKeyRemindAdd is the description key for the remind add command.
	DescKeyRemindAdd = "remind.add"
	// DescKeyRemindDismiss is the description key for the remind dismiss command.
	DescKeyRemindDismiss = "remind.dismiss"
	// DescKeyRemindList is the description key for the remind list command.
	DescKeyRemindList = "remind.list"
)
