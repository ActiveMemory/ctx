//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger provides terminal output for the
// trigger hook commands (ctx trigger create, list,
// enable, disable, test).
//
// Triggers are user-defined hook scripts that run
// at specific lifecycle points. The output functions
// cover management, listing, and testing workflows.
//
// # Management
//
// [Created] confirms a hook script was created at
// a given path. [Disabled] and [Enabled] confirm
// status changes, printing hook name and path.
//
// # Listing
//
// [TypeHeader] prints a section header for each
// hook type. [Entry] prints a single hook with
// its name, enabled/disabled status, and path.
// [Count] prints the total hook count.
// [NoHooksFound] handles the empty-list case.
// [BlankLine] separates sections visually.
//
// # Testing
//
// [TestingHeader] prints the header for a hook
// test run. [TestInput] prints the JSON input
// block sent to the hook. [ContextOutput] prints
// context output from hook execution.
// [Cancelled] prints a cancellation message.
// [ErrorsHeader] and [ErrorLine] render the
// errors section. [NoOutput] reports when hooks
// produced no output.
package trigger
