//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for the `ctx system` unknown-subcommand verbatim relay.
const (
	// DescKeySystemUnknownRelayPrefix is the text key for the
	// "relay this verbatim" prefix above the unknown-subcommand box.
	DescKeySystemUnknownRelayPrefix = "system-unknown.relay-prefix"
	// DescKeySystemUnknownBoxTitle is the text key for the
	// unknown-subcommand box title.
	DescKeySystemUnknownBoxTitle = "system-unknown.box-title"
	// DescKeySystemUnknownBody is the text key for the
	// unknown-subcommand box body. It is a format string taking the
	// unrecognised subcommand name.
	DescKeySystemUnknownBody = "system-unknown.body"
	// DescKeySystemUnknownRelayMessage is the text key for the
	// human-readable relay-event description recorded in the event
	// log / webhook. Format string taking the subcommand name.
	DescKeySystemUnknownRelayMessage = "system-unknown.relay-message"
)
