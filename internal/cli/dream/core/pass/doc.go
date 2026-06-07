//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pass orchestrates one executor-agnostic ctx dream run:
// ensure the gitignored dreams/ notebook exists, compute the delta
// against saved state, gate on an empty delta, serialize with a lock,
// invoke the configured executor (fail-loud on a missing binary),
// validate the proposals the executor wrote, persist state for the
// processed sources, and print a short digest.
package pass
