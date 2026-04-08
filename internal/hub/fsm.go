//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	stdio "io"

	"github.com/hashicorp/raft"
)

// leaderFSM is a no-op finite state machine for Raft.
//
// Raft requires an FSM, but we only use Raft for leader
// election — data replication is handled separately via
// sequence-based gRPC sync. All FSM methods are no-ops.
type leaderFSM struct{}

// Apply is a no-op. Data is not replicated via Raft.
//
// Parameters:
//   - log: Raft log entry (ignored)
//
// Returns:
//   - any: nil
func (f *leaderFSM) Apply(_ *raft.Log) any {
	return nil
}

// Snapshot returns a no-op snapshot. No state to persist.
//
// Returns:
//   - raft.FSMSnapshot: no-op snapshot
//   - error: always nil
func (f *leaderFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &noopSnapshot{}, nil
}

// Restore is a no-op. No state to restore.
//
// Parameters:
//   - rc: snapshot reader (ignored)
//
// Returns:
//   - error: always nil
func (f *leaderFSM) Restore(_ stdio.ReadCloser) error {
	return nil
}

// noopSnapshot is a no-op snapshot implementation.
type noopSnapshot struct{}

// Persist is a no-op.
//
// Parameters:
//   - sink: snapshot sink (ignored)
//
// Returns:
//   - error: always nil
func (s *noopSnapshot) Persist(_ raft.SnapshotSink) error {
	return nil
}

// Release is a no-op.
func (s *noopSnapshot) Release() {}
