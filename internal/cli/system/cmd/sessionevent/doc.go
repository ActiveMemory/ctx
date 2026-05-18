//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sessionevent implements the hidden
// "ctx system session-event" cobra subcommand.
//
// Editor integrations call this command to signal
// session lifecycle boundaries (start or end). The
// event is recorded in the event log and broadcast
// as a notification so downstream consumers can
// react to session transitions.
//
// # Usage
//
//	ctx system session-event \
//	    --type <start|end> --caller <editor>
//
// # Flags
//
//	--type     Required. The event type: "start" or
//	           "end". Any other value returns an error.
//	--caller   Required. Identifier for the calling
//	           editor or integration, for example
//	           "vscode" or "cursor".
//
// # Behavior
//
// The command validates the event type, then:
//
//   - Formats a human-readable message containing
//     the event type and caller.
//   - Appends the message to the event log under the
//     "session" category.
//   - Sends a notification on the session channel
//     with a template ref keyed to the event type.
//   - Prints a confirmation line to stdout.
//
// The command is a no-op when the context directory
// is not initialized.
//
// # Output
//
// Prints a one-line confirmation of the recorded
// event including the type and caller.
//
// # Delegation
//
// Event logging uses log/event. Notifications use
// the notify package. Output formatting uses
// write/session.
package sessionevent
