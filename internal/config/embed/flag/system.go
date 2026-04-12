//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flag

// DescKeys for hidden ctx system hook-plumbing command flags.
//
// User-facing command flags (backup, bootstrap, event, message,
// prune, resource, stats) live in their own per-command files in
// this package since those commands were promoted to top-level.
const (
	// DescKeySystemMarkJournalCheck is the description key for the system mark
	// journal check flag.
	DescKeySystemMarkJournalCheck = "system.markjournal.check"
	// DescKeySystemPauseSessionId is the description key for the system pause
	// session id flag.
	DescKeySystemPauseSessionId = "system.pause.session-id"
	// DescKeySystemResumeSessionId is the description key for the system resume
	// session id flag.
	DescKeySystemResumeSessionId = "system.resume.session-id"
	// DescKeySystemSessionEventCaller is the description key for the system
	// session event caller flag.
	DescKeySystemSessionEventCaller = "system.sessionevent.caller"
	// DescKeySystemSessionEventType is the description key for the system session
	// event type flag.
	DescKeySystemSessionEventType = "system.sessionevent.type"
)
