//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify provides terminal output for webhook notification
// setup and testing (ctx notify setup, ctx notify test).
//
// [SetupPrompt] displays the webhook URL prompt, [SetupDone]
// confirms successful configuration. [TestResult] reports the HTTP
// response from a test notification, [TestNoWebhook] handles the
// unconfigured case, and [TestFiltered] explains when an event
// type is excluded by the filter.
package notify
