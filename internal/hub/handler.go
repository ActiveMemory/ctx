//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
)

// register handles the Register RPC.
//
// Parameters:
//   - ctx: request context (unused)
//   - req: registration request with admin token
//
// Returns:
//   - *RegisterResponse: client ID and token
//   - error: non-nil if auth or registration fails
func (s *Server) register(
	_ context.Context, req *RegisterRequest,
) (*RegisterResponse, error) {
	if req.AdminToken != s.adminToken {
		return nil, status.Error(
			codes.PermissionDenied,
			cfgHub.ErrInvalidAdminToken,
		)
	}
	if req.ProjectName == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			cfgHub.ErrProjectNameRequired,
		)
	}

	clientToken, genErr := GenerateClientToken()
	if genErr != nil {
		return nil, errHub.InternalErr(genErr)
	}

	clientID, idErr := generateClientID()
	if idErr != nil {
		return nil, errHub.InternalErr(idErr)
	}

	client := ClientInfo{
		ID:          clientID,
		ProjectName: req.ProjectName,
		Token:       clientToken,
	}
	if regErr := s.store.RegisterClient(client); regErr != nil {
		return nil, errHub.InternalErr(regErr)
	}

	return &RegisterResponse{
		ClientID:    clientID,
		ClientToken: clientToken,
	}, nil
}

// publish handles the Publish RPC.
//
// Parameters:
//   - ctx: request context (unused)
//   - req: publish request with entries
//
// Returns:
//   - *PublishResponse: assigned sequence numbers
//   - error: non-nil if validation or append fails
func (s *Server) publish(
	_ context.Context, req *PublishRequest,
) (*PublishResponse, error) {
	if len(req.Entries) == 0 {
		return &PublishResponse{}, nil
	}

	for _, pe := range req.Entries {
		if valErr := validateEntry(pe); valErr != nil {
			return nil, valErr
		}
	}

	entries := make([]Entry, len(req.Entries))
	for i, pe := range req.Entries {
		entries[i] = Entry{
			ID:        pe.ID,
			Type:      pe.Type,
			Content:   pe.Content,
			Origin:    pe.Origin,
			Meta:      pe.Meta,
			Timestamp: time.Unix(pe.Timestamp, 0),
		}
	}

	seqs, appendErr := s.store.Append(entries)
	if appendErr != nil {
		return nil, errHub.InternalErr(appendErr)
	}

	for i := range entries {
		entries[i].Sequence = seqs[i]
	}
	s.listeners.broadcast(entries)

	return &PublishResponse{Sequences: seqs}, nil
}

// syncEntries handles the Sync RPC (server-streaming).
//
// Parameters:
//   - req: sync request with type filter and sequence
//   - send: callback to send each entry to the client
//
// Returns:
//   - error: non-nil if send fails
func (s *Server) syncEntries(
	req *SyncRequest, send func(*EntryMsg) error,
) error {
	results := s.store.Query(
		req.Types, req.SinceSequence,
	)
	for i := range results {
		if sendErr := send(
			entryToMsg(&results[i]),
		); sendErr != nil {
			return sendErr
		}
	}
	return nil
}

// listenEntries handles the Listen RPC (long-lived stream).
//
// Parameters:
//   - req: listen request with type filter and sequence
//   - send: callback to send each entry to the client
//   - ctx: context for cancellation
//
// Returns:
//   - error: non-nil if send fails
func (s *Server) listenEntries(
	req *ListenRequest,
	send func(*EntryMsg) error,
	ctx context.Context,
) error {
	results := s.store.Query(
		req.Types, req.SinceSequence,
	)
	for i := range results {
		if sendErr := send(
			entryToMsg(&results[i]),
		); sendErr != nil {
			return sendErr
		}
	}

	typeSet := make(map[string]bool, len(req.Types))
	for _, t := range req.Types {
		typeSet[t] = true
	}

	ch := s.listeners.subscribe()
	defer s.listeners.unsubscribe(ch)

	for {
		select {
		case <-ctx.Done():
			return nil
		case entries := <-ch:
			for i := range entries {
				if len(typeSet) > 0 &&
					!typeSet[entries[i].Type] {
					continue
				}
				if sendErr := send(
					entryToMsg(&entries[i]),
				); sendErr != nil {
					return sendErr
				}
			}
		}
	}
}

// hubStatus handles the Status RPC.
//
// Parameters:
//   - ctx: request context (unused)
//
// Returns:
//   - *StatusResponse: hub statistics
//   - error: always nil
func (s *Server) hubStatus(
	_ context.Context,
) (*StatusResponse, error) {
	total, byType, byProject := s.store.Stats()
	return &StatusResponse{
		TotalEntries:     total,
		ConnectedClients: s.listeners.count(),
		EntriesByType:    byType,
		EntriesByProject: byProject,
	}, nil
}
