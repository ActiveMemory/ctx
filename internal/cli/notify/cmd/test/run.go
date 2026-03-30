//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/spf13/cobra"

	coreTest "github.com/ActiveMemory/ctx/internal/cli/notify/core/test"
	"github.com/ActiveMemory/ctx/internal/config/crypto"
	writeNotify "github.com/ActiveMemory/ctx/internal/write/notify"
)

// Run sends a test notification to the configured webhook.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on webhook load or HTTP failure
func Run(cmd *cobra.Command) error {
	r, sendErr := coreTest.Send()
	if sendErr != nil {
		return sendErr
	}

	if r.NoWebhook {
		writeNotify.TestNoWebhook(cmd)
		return nil
	}

	if r.Filtered {
		writeNotify.TestFiltered(cmd)
	}

	writeNotify.TestResult(cmd, r.StatusCode, coreTest.OK(r), crypto.NotifyEnc)
	return nil
}
