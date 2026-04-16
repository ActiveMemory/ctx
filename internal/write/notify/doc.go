//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify provides terminal output for webhook
// notification setup and testing (ctx hook notify setup,
// ctx hook notify test).
//
// # Setup Flow
//
// [SetupPrompt] displays the interactive webhook URL
// prompt where the user enters their endpoint.
// [SetupDone] prints the success block after saving
// a webhook, showing the masked URL and the encrypted
// file path where credentials are stored.
//
// # Test Flow
//
// [TestResult] reports the HTTP response from a test
// notification including the status code and status
// text. When the response indicates success (2xx),
// an additional confirmation line is printed.
//
// [TestNoWebhook] handles the case when no webhook is
// configured. [TestFiltered] explains when a test
// event type is excluded by the user's event filter
// configuration.
//
// # Message Categories
//
//   - Info: setup confirmation, test results
//   - Warning: unconfigured or filtered states
//
// # Nil Safety
//
// All functions treat a nil *cobra.Command as a no-op.
package notify
