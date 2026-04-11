//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// replicateOnce connects to the master, syncs all entries
// since the local store's last sequence, and appends them.
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

	_, lastSeq := lastSequence(store)
	authed := addBearerMD(ctx, clientToken)

	stream, streamErr := conn.NewStream(
		authed,
		&grpc.StreamDesc{ServerStreams: true},
		"/ctx.hub.v1.CtxHub/Sync",
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
			Author:    msg.Author,
			Timestamp: time.Unix(msg.Timestamp, 0),
			Sequence:  msg.Sequence,
		}
		_, _ = store.Append([]Entry{entry})
	}
}

// lastSequence returns the highest sequence in the store.
func lastSequence(store *Store) (bool, uint64) {
	all := store.Query(nil, 0)
	if len(all) == 0 {
		return false, 0
	}
	return true, all[len(all)-1].Sequence
}

// addBearerMD adds a bearer token to outgoing gRPC
// metadata.
func addBearerMD(
	ctx context.Context, tok string,
) context.Context {
	if tok == "" {
		return ctx
	}
	return metadata.NewOutgoingContext(
		ctx, metadata.Pairs(
			"authorization", bearerPrefix+tok,
		),
	)
}
