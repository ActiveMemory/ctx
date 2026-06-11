//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pass

// Opts carries the run-pass parameters resolved from flags and rc.
//
// Fields:
//   - Mode: execution mode (discipline in v1)
//   - Max: ceiling on ideas files processed this pass
//   - Budget: step/token budget for the pass (executor turn bound)
//   - Force: bypass the opt-in/cadence trigger gate for a manual run
type Opts struct {
	Mode   string
	Max    int
	Budget int
	Force  bool
}
