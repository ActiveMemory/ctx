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

// NewFailoverClient creates a client that tries peers in
// order until one succeeds. Fails fast on auth errors
// since the token is the same for all peers.
//
// Parameters:
//   - peers: ordered list of hub addresses
//   - bearerToken: token for authenticated RPCs
//
// Returns:
//   - *Client: connected client to the first reachable peer
//   - error: non-nil if no peer is reachable
func NewFailoverClient(
	peers []string, bearerToken string,
) (*Client, error) {
	var lastErr error
	for _, addr := range peers {
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
			lastErr = dialErr
			continue
		}

		resp := &StatusResponse{}
		callErr := conn.Invoke(
			addBearerMD(
				context.Background(), bearerToken,
			),
			"/ctx.hub.v1.CtxHub/Status",
			&struct{}{},
			resp,
		)
		if callErr != nil {
			_ = conn.Close()

			// Fail fast on auth errors — same token
			// won't work on other peers either.
			if isAuthErr(callErr) {
				return nil, callErr
			}

			lastErr = callErr
			continue
		}

		return &Client{
			conn: conn, token: bearerToken,
		}, nil
	}
	return nil, lastErr
}
