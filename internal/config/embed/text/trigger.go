//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for trigger (hook runner) output.
const (
	DescKeyTriggerWarn      = "trigger.warn"
	DescKeyTriggerErrorItem = "trigger.error-item"
	DescKeyTriggerSkipWarn  = "trigger.skip-warn"
)

// DescKeys for write/trigger display output.
const (
	DescKeyWriteTriggerCreated   = "write.trigger-created"
	DescKeyWriteTriggerDisabled  = "write.trigger-disabled"
	DescKeyWriteTriggerEnabled   = "write.trigger-enabled"
	DescKeyWriteTriggerTypeHdr   = "write.trigger-type-hdr"
	DescKeyWriteTriggerEntry     = "write.trigger-entry"
	DescKeyWriteTriggerCount     = "write.trigger-count"
	DescKeyWriteTriggerTestHdr   = "write.trigger-test-hdr"
	DescKeyWriteTriggerTestInput = "write.trigger-test-input"
	DescKeyWriteTriggerCancelled = "write.trigger-cancelled"
	DescKeyWriteTriggerContext   = "write.trigger-context"
	DescKeyWriteTriggerErrLine   = "write.trigger-err-line"
)
