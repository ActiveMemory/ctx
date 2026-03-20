//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// WebhookEmpty returns an error for blank webhook URL input.
//
// Returns:
//   - error: "webhook URL cannot be empty"
func WebhookEmpty() error {
	return errors.New(desc.TextDesc(text.DescKeyErrNotifyWebhookEmpty))
}

// SaveWebhook wraps a webhook save failure.
//
// Parameters:
//   - cause: the underlying error from the save operation.
//
// Returns:
//   - error: "save webhook: <cause>"
func SaveWebhook(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrNotifySaveWebhook), cause,
	)
}

// LoadWebhook wraps a webhook load failure.
//
// Parameters:
//   - cause: the underlying error from the load operation.
//
// Returns:
//   - error: "load webhook: <cause>"
func LoadWebhook(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrNotifyLoadWebhook), cause,
	)
}

// MarshalPayload wraps a JSON marshal failure.
//
// Parameters:
//   - cause: the underlying marshal error.
//
// Returns:
//   - error: "marshal payload: <cause>"
func MarshalPayload(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrNotifyMarshalPayload), cause,
	)
}

// SendNotification wraps a notification send failure.
//
// Parameters:
//   - cause: the underlying HTTP error.
//
// Returns:
//   - error: "send test notification: <cause>"
func SendNotification(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrNotifySendNotification), cause,
	)
}
