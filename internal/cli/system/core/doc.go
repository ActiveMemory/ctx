//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared business logic used by
// all system subcommands and hook handlers.
//
// This is the largest core namespace in the project.
// It groups many subpackages that handle bootstrap,
// session management, event logging, hook checks,
// nudge delivery, persistence, and system health.
//
// # Key Subpackages
//
// Bootstrap: session initialization and context
// directory inventory for agent startup.
//
// Ceremony: detection of missing context ceremonies
// (remember and wrap-up) in recent journal entries.
//
// Check: shared hook preamble logic, daily throttling,
// and wrap-up recency detection.
//
// Counter: integer counter persistence in state files.
//
// Event: event log formatting in JSON and human-
// readable column layouts.
//
// Heartbeat: mtime value persistence for session
// activity tracking.
//
// Message: hook message loading and nudge box
// rendering from templates and overrides.
//
// Nudge: pause state, token counting, context size
// detection, and nudge relay for hook notifications.
//
// Session: hook input parsing and session token
// management.
//
// State, Stats, Log, Persistence, Resource, Health,
// Drift, Version, and others provide lower-level
// utilities for the system command surface.
package core
