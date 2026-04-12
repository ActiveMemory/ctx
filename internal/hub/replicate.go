//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"time"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Reserved for cluster mode: startReplication will be
// called from Server.Start when a follower peer is
// configured.
var _ = startReplication

// startReplication connects to the master and streams
// entries into the local store. Blocks until the context
// is canceled. Retries on failure.
//
// Parameters:
//   - ctx: context for cancellation
//   - masterAddr: gRPC address of the master hub
//   - store: local store to write replicated entries
//   - clientToken: bearer token for auth
func startReplication(
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
		case <-time.After(cfgHub.ReplicateInterval * time.Second):
		}
	}
}

// replicateOnce connects to the master, syncs all entries
// since the local store's last sequence, and appends them.
//
// Parameters:
//   - ctx: context for cancellation
//   - masterAddr: gRPC address of the master hub
//   - store: local store to write replicated entries
//   - clientToken: bearer token for auth
func replicateOnce(
	ctx context.Context,
	masterAddr string,
	store *Store,
	clientToken string,
) {
	conn, dialErr := grpc.NewClient(
		masterAddr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithDefaultCallOptions(
			grpc.CallContentSubtype(codecName),
		),
	)
	if dialErr != nil {
		return
	}
	defer func() { _ = conn.Close() }()

	_, lastSeq := store.lastSequence()
	authed := addBearerMD(ctx, clientToken)

	stream, streamErr := conn.NewStream(
		authed,
		&grpc.StreamDesc{ServerStreams: true},
		cfgHub.PathSync,
	)
	if streamErr != nil {
		return
	}

	if sendErr := stream.SendMsg(&SyncRequest{
		SinceSequence: lastSeq,
	}); sendErr != nil {
		return
	}
	if closeErr := stream.CloseSend(); closeErr != nil {
		return
	}

	for {
		msg := &EntryMsg{}
		if recvErr := stream.RecvMsg(msg); recvErr != nil {
			return
		}
		entry := Entry{
			ID:        msg.ID,
			Type:      msg.Type,
			Content:   msg.Content,
			Origin:    msg.Origin,
			Meta:      msg.Meta,
			Timestamp: time.Unix(msg.Timestamp, 0),
			Sequence:  msg.Sequence,
		}
		_, _ = store.Append([]Entry{entry})
	}
}
