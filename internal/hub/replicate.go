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
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"

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
		logWarn.Warn(cfgWarn.HubReplicateDial, masterAddr, dialErr)
		return
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			logWarn.Warn(cfgWarn.Close, masterAddr, cerr)
		}
	}()

	_, lastSeq := store.lastSequence()
	authed := addBearerMD(ctx, clientToken)

	stream, streamErr := conn.NewStream(
		authed,
		&grpc.StreamDesc{ServerStreams: true},
		cfgHub.PathSync,
	)
	if streamErr != nil {
		logWarn.Warn(cfgWarn.HubReplicateStream, masterAddr, streamErr)
		return
	}

	if sendErr := stream.SendMsg(&SyncRequest{
		SinceSequence: lastSeq,
	}); sendErr != nil {
		logWarn.Warn(cfgWarn.HubReplicateSend, masterAddr, sendErr)
		return
	}
	if closeErr := stream.CloseSend(); closeErr != nil {
		logWarn.Warn(cfgWarn.HubReplicateCloseSend, masterAddr, closeErr)
		return
	}

	for {
		msg := &EntryMsg{}
		if recvErr := stream.RecvMsg(msg); recvErr != nil {
			// io.EOF is the normal end of every sync stream
			// and a done caller context is routine shutdown;
			// warning on either would spam stderr once per
			// replication cycle. Anything else is a transport
			// failure worth surfacing.
			if !eof(recvErr) && ctx.Err() == nil {
				logWarn.Warn(
					cfgWarn.HubReplicateRecv, masterAddr, recvErr,
				)
			}
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
		if _, appendErr := store.Append([]Entry{entry}); appendErr != nil {
			logWarn.Warn(cfgWarn.HubReplicateAppend, appendErr)
		}
	}
}
