//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"time"
)

// replicateInterval is how often a follower retries
// connecting to the master for replication.
const replicateInterval = 5 * time.Second

// StartReplication connects to the master and streams
// entries into the local store. Blocks until the context
// is cancelled. Retries on failure.
//
// Parameters:
//   - ctx: context for cancellation
//   - masterAddr: gRPC address of the master hub
//   - store: local store to write replicated entries
//   - clientToken: bearer token for auth
func StartReplication(
	ctx context.Context,
	masterAddr string,
	store *Store,
	clientToken string,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		replicateOnce(
			ctx, masterAddr, store, clientToken,
		)

		select {
		case <-ctx.Done():
			return
		case <-time.After(replicateInterval):
		}
	}
}
