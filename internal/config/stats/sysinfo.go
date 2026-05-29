//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

// Resources display formatting.
const (
	// ResourcesStatusCol is the column where the status indicator starts
	// in the resources text output.
	ResourcesStatusCol = 52
	// ResourcesLabelWidth is the left-aligned label column width.
	ResourcesLabelWidth = 7
)

// Resources formatting patterns.
const (
	// FormatGiB is the precision format for GiB values.
	FormatGiB = "%.1f"
)

// Stats command defaults.
const (
	// DefaultLast is the default number of stats entries to display.
	DefaultLast = 20
)

// Resource threshold constants for health evaluation.
const (
	// ThresholdMemPressureSomeWarnPct is the Linux PSI
	// "some" avg10 percentage (share of the last 10s in
	// which at least one task stalled on memory) that
	// triggers a warning. PSI is a rate-of-stall pressure
	// signal, not a static-occupancy ratio.
	ThresholdMemPressureSomeWarnPct = 10.0
	// ThresholdMemPressureFullDangerPct is the Linux PSI
	// "full" avg10 percentage (share of the last 10s in
	// which every runnable task stalled on memory) that
	// triggers a danger alert.
	ThresholdMemPressureFullDangerPct = 10.0
	// ThresholdDiskWarnPct is the disk usage percentage
	// that triggers a warning.
	ThresholdDiskWarnPct = 85
	// ThresholdDiskDangerPct is the disk usage percentage
	// that triggers a danger alert.
	ThresholdDiskDangerPct = 95
	// ThresholdLoadWarnRatio is the load-to-CPU ratio
	// that triggers a warning.
	ThresholdLoadWarnRatio = 0.8
	// ThresholdLoadDangerRatio is the load-to-CPU ratio
	// that triggers a danger alert.
	ThresholdLoadDangerRatio = 1.5
	// ThresholdBytesPerGiB is the number of bytes in one gibibyte.
	ThresholdBytesPerGiB = 1 << 30
)
