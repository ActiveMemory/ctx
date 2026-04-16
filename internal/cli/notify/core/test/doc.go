//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package test provides webhook test notification
// logic.
//
// The "ctx notify test" command sends a test payload
// to the configured webhook URL to verify that the
// notification pipeline works end-to-end. This package
// contains the business logic for building, sending,
// and evaluating the test notification.
//
// # Sending a Test Notification
//
// [Send] is the primary function. It performs the
// following steps:
//
//  1. Loads the webhook URL from the project
//     configuration via notify.LoadWebhook. If no URL
//     is configured, it returns a Result with
//     NoWebhook set to true.
//  2. Determines the project name from the current
//     working directory, falling back to a default
//     name if the directory cannot be resolved.
//  3. Builds a NotifyPayload with event type "test",
//     a test message, an RFC3339 timestamp, and the
//     project name.
//  4. Marshals the payload to JSON.
//  5. Checks whether the "test" event type is allowed
//     by the configured event filter, recording the
//     filtered flag.
//  6. Posts the JSON body to the webhook URL via
//     notify.PostJSON.
//  7. Returns a Result containing the HTTP status code
//     and filtered flag.
//
// # Result Evaluation
//
// [OK] checks whether a Result indicates success by
// testing that the status code falls in the 2xx range.
//
// # Result Type
//
// The [Result] struct carries three fields: NoWebhook
// (no URL configured), Filtered (event excluded by
// filter), and StatusCode (the HTTP response code).
package test
