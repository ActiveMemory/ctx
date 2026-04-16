//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_resource implements the
// **`ctx system check-resource`** hidden hook, which
// warns when system resources reach dangerous levels.
//
// # What It Does
//
// The hook collects a snapshot of system resource
// metrics (memory usage, swap pressure, disk space,
// CPU load average) and evaluates each against
// configured thresholds. When any metric reaches
// "danger" severity, the hook emits a nudge box
// listing the affected resources and recommending
// actions such as persisting work and ending the
// session before the system becomes unresponsive.
//
// Metrics below the danger threshold are silently
// ignored.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On danger-level resource: a nudge box listing each
// danger alert with an error icon. On all resources
// healthy or throttled: no output.
//
// # Throttling
//
// The hook respects the session pause state but has
// no daily throttle; resource warnings fire on every
// invocation when danger levels are detected.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], collects metrics
// through [sysinfo.Collect], evaluates thresholds with
// [sysinfo.Evaluate], and emits the warning via
// [core/nudge.LoadAndEmit].
package check_resource
