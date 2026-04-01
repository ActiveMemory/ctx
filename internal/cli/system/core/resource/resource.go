//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// Snapshot collects system resource metrics and evaluates alert thresholds.
//
// Returns:
//   - sysinfo.Snapshot: Current system metrics
//   - []sysinfo.ResourceAlert: Alerts for metrics exceeding thresholds
func Snapshot() (sysinfo.Snapshot, []sysinfo.ResourceAlert) {
	snap := sysinfo.Collect()
	alerts := sysinfo.Evaluate(snap)
	return snap, alerts
}
