//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hubsync pulls new entries from a remote hub
// and writes them to the local .context/hub/ directory.
// It is triggered automatically during session start by
// the check-hub-sync hook.
//
// # Connection Check
//
// [Connected] reports whether a hub connection config
// exists by checking for the .context/.connect.enc file.
// Hooks call this before attempting a sync to avoid
// unnecessary work when no hub is configured.
//
// # Sync Flow
//
// [Sync] loads the connection config, dials the hub,
// pulls entries matching the configured types, and
// writes them to disk via the connection render layer.
// It returns a formatted status message with the count
// of synced entries, or an empty string when nothing
// was fetched. Every error is silently swallowed so the
// hook never blocks the session start.
//
// The data flow is:
//
//  1. Load connection config from .context/.connect.enc
//  2. Dial the hub using the configured address and token
//  3. Pull entries filtered by configured content types
//  4. Write entries to .context/hub/ via render layer
//  5. Return a formatted summary for the nudge box
package hubsync
