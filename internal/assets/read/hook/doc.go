//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook provides access to hook message
// templates and the hook registry from embedded
// assets.
//
// # Message Templates
//
// Message reads a specific template file by hook
// name and filename. Templates are stored under
// hooks/messages/<hook>/<filename> in the embedded
// filesystem.
//
//	data, err := hook.Message("qa-reminder", "gate.txt")
//
// # Registry
//
// MessageRegistry returns the raw registry.yaml that
// describes all hook message templates, their
// categories, and template variables. The registry is
// parsed by the hooks/messages package for structured
// access.
//
//	raw, err := hook.MessageRegistry()
//
// # Trace Scripts
//
// TraceScript reads an embedded git hook script by
// filename. These scripts are installed into
// .git/hooks/ by ctx init to enable commit tracing.
//
//	script, err := hook.TraceScript("prepare-commit-msg.sh")
package hook
