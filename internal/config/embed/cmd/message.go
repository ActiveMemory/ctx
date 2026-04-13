//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for the top-level message command and its subcommands.
const (
	// UseMessage is the cobra Use string for the message command.
	UseMessage = "message"
	// UseMessageEdit is the cobra Use string for the message edit command.
	UseMessageEdit = "edit <hook> <variant>"
	// UseMessageList is the cobra Use string for the message list command.
	UseMessageList = "list"
	// UseMessageReset is the cobra Use string for the message reset command.
	UseMessageReset = "reset <hook> <variant>"
	// UseMessageShow is the cobra Use string for the message show command.
	UseMessageShow = "show <hook> <variant>"
)

// DescKeys for the top-level message command and its subcommands.
const (
	// DescKeyMessage is the description key for the message command.
	DescKeyMessage = "message"
	// DescKeyMessageEdit is the description key for the message edit command.
	DescKeyMessageEdit = "message.edit"
	// DescKeyMessageList is the description key for the message list command.
	DescKeyMessageList = "message.list"
	// DescKeyMessageReset is the description key for the message reset command.
	DescKeyMessageReset = "message.reset"
	// DescKeyMessageShow is the description key for the message show command.
	DescKeyMessageShow = "message.show"
)
