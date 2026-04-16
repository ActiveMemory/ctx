//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify defines the typed error constructors
// for the webhook notification subsystem. These
// errors fire when configuring, persisting, or
// sending webhook notifications triggered by context
// changes.
//
// # Domain
//
// Errors fall into three categories:
//
//   - **Validation** -- the webhook URL is blank.
//     Constructor: [WebhookEmpty].
//   - **Persistence** -- saving or loading the
//     encrypted webhook configuration failed.
//     Constructors: [SaveWebhook], [LoadWebhook].
//   - **Delivery** -- marshaling the JSON payload
//     or sending the HTTP request failed.
//     Constructors: [MarshalPayload],
//     [SendNotification].
//
// # Wrapping Strategy
//
// IO and delivery constructors wrap their cause
// with fmt.Errorf %w so callers can inspect the
// underlying error. [WebhookEmpty] returns a plain
// errors.New value. All user-facing text is
// resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package notify
