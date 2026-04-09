//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates a hub client connected to the given address.
//
// Parameters:
//   - addr: hub gRPC address (host:port)
//   - token: bearer token for authenticated RPCs
//
// Returns:
//   - *Client: connected client
//   - error: non-nil if connection fails
func NewClient(
	addr string, token string,
) (*Client, error) {
	conn, dialErr := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithDefaultCallOptions(
			grpc.CallContentSubtype(codecName),
		),
	)
	if dialErr != nil {
		return nil, dialErr
	}
	return &Client{conn: conn, token: token}, nil
}

// Register calls the Register RPC with the admin token.
//
// Parameters:
//   - ctx: context for the call
//   - adminToken: admin token from hub startup
//   - projectName: name of the project to register
//
// Returns:
//   - *RegisterResponse: client ID and token
//   - error: non-nil if registration fails
func (c *Client) Register(
	ctx context.Context,
	adminToken string,
	projectName string,
) (*RegisterResponse, error) {
	resp := &RegisterResponse{}
	callErr := c.conn.Invoke(
		ctx,
		"/ctx.hub.v1.CtxHub/Register",
		&RegisterRequest{
			AdminToken:  adminToken,
			ProjectName: projectName,
		},
		resp,
	)
	return resp, callErr
}

// Publish calls the Publish RPC.
//
// Parameters:
//   - ctx: context for the call
//   - entries: entries to publish
//
// Returns:
//   - *PublishResponse: assigned sequence numbers
//   - error: non-nil if publish fails
func (c *Client) Publish(
	ctx context.Context,
	entries []PublishEntry,
) (*PublishResponse, error) {
	resp := &PublishResponse{}
	callErr := c.conn.Invoke(
		c.authedCtx(ctx),
		"/ctx.hub.v1.CtxHub/Publish",
		&PublishRequest{Entries: entries},
		resp,
	)
	return resp, callErr
}

// Sync calls the Sync RPC and collects all entries.
//
// Parameters:
//   - ctx: context for the call
//   - types: entry types to sync (empty = all)
//   - sinceSequence: return entries after this sequence
//
// Returns:
//   - []EntryMsg: matching entries
//   - error: non-nil if sync fails
func (c *Client) Sync(
	ctx context.Context,
	types []string,
	sinceSequence uint64,
) ([]EntryMsg, error) {
	stream, streamErr := c.conn.NewStream(
		c.authedCtx(ctx),
		&grpc.StreamDesc{ServerStreams: true},
		"/ctx.hub.v1.CtxHub/Sync",
	)
	if streamErr != nil {
		return nil, streamErr
	}

	if sendErr := stream.SendMsg(&SyncRequest{
		Types:         types,
		SinceSequence: sinceSequence,
	}); sendErr != nil {
		return nil, sendErr
	}
	if closeErr := stream.CloseSend(); closeErr != nil {
		return nil, closeErr
	}

	var entries []EntryMsg
	for {
		msg := &EntryMsg{}
		if recvErr := stream.RecvMsg(msg); recvErr != nil {
			if isEOF(recvErr) {
				break
			}
			return nil, recvErr
		}
		entries = append(entries, *msg)
	}
	return entries, nil
}

// Listen opens a server-streaming Listen RPC and calls
// the handler for each entry received. Blocks until the
// context is cancelled or the stream ends.
//
// Parameters:
//   - ctx: context for cancellation
//   - types: entry types to receive (empty = all)
//   - sinceSequence: start from this sequence
//   - handler: called for each received entry
//
// Returns:
//   - error: non-nil if stream setup or recv fails
func (c *Client) Listen(
	ctx context.Context,
	types []string,
	sinceSequence uint64,
	handler func(EntryMsg) error,
) error {
	stream, streamErr := c.conn.NewStream(
		c.authedCtx(ctx),
		&grpc.StreamDesc{ServerStreams: true},
		"/ctx.hub.v1.CtxHub/Listen",
	)
	if streamErr != nil {
		return streamErr
	}

	if sendErr := stream.SendMsg(&ListenRequest{
		Types:         types,
		SinceSequence: sinceSequence,
	}); sendErr != nil {
		return sendErr
	}
	if closeErr := stream.CloseSend(); closeErr != nil {
		return closeErr
	}

	for {
		msg := &EntryMsg{}
		if recvErr := stream.RecvMsg(msg); recvErr != nil {
			if isEOF(recvErr) {
				return nil
			}
			return recvErr
		}
		if handleErr := handler(*msg); handleErr != nil {
			return handleErr
		}
	}
}

// Status calls the Status RPC.
//
// Parameters:
//   - ctx: context for the call
//
// Returns:
//   - *StatusResponse: hub statistics
//   - error: non-nil if call fails
func (c *Client) Status(
	ctx context.Context,
) (*StatusResponse, error) {
	resp := &StatusResponse{}
	callErr := c.conn.Invoke(
		c.authedCtx(ctx),
		"/ctx.hub.v1.CtxHub/Status",
		&struct{}{},
		resp,
	)
	return resp, callErr
}

// Close closes the underlying gRPC connection.
//
// Returns:
//   - error: non-nil if close fails
func (c *Client) Close() error {
	return c.conn.Close()
}
