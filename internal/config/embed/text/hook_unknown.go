//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for the `ctx hook` unknown-subcommand verbatim relay.
const (
	// DescKeyHookUnknownRelayPrefix is the text key for the
	// "relay this verbatim" prefix above the unknown-subcommand box.
	DescKeyHookUnknownRelayPrefix = "hook-unknown.relay-prefix"
	// DescKeyHookUnknownBoxTitle is the text key for the
	// unknown-subcommand box title.
	DescKeyHookUnknownBoxTitle = "hook-unknown.box-title"
	// DescKeyHookUnknownBody is the text key for the
	// unknown-subcommand box body. It is a format string taking the
	// unrecognised subcommand name.
	DescKeyHookUnknownBody = "hook-unknown.body"
	// DescKeyHookUnknownRelayMessage is the text key for the
	// human-readable relay-event description recorded in the event
	// log / webhook. Format string taking the subcommand name.
	DescKeyHookUnknownRelayMessage = "hook-unknown.relay-message"
)
