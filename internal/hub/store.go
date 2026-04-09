//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"crypto/subtle"
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/config/token"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
	"github.com/ActiveMemory/ctx/internal/io"
)

// NewStore creates or opens a Store in the given directory.
//
// On first run, creates the directory and initializes empty
// data files. On subsequent runs, loads existing entries,
// clients, and metadata.
//
// Parameters:
//   - dir: directory path for data files
//
// Returns:
//   - *Store: initialized store
//   - error: non-nil if directory or file loading fails
func NewStore(dir string) (*Store, error) {
	if mkErr := io.SafeMkdirAll(dir, dirPerm); mkErr != nil {
		return nil, mkErr
	}

	s := &Store{
		dir:      dir,
		tokenIdx: make(map[string]int),
	}

	if loadErr := loadJSON(metaPath(dir), &s.meta); loadErr != nil {
		return nil, loadErr
	}
	if loadErr := loadJSON(
		clientsPath(dir), &s.clients,
	); loadErr != nil {
		return nil, loadErr
	}
	if loadErr := loadEntries(dir, &s.entries); loadErr != nil {
		return nil, loadErr
	}

	// Build token index from loaded clients.
	for i := range s.clients {
		s.tokenIdx[s.clients[i].Token] = i
	}

	return s, nil
}

// Append adds entries to the store, assigning sequence numbers.
//
// Each entry gets the next monotonic sequence number. Entries
// are appended to the JSONL file and metadata is updated.
//
// Parameters:
//   - entries: entries to append (Sequence is overwritten)
//
// Returns:
//   - []uint64: assigned sequence numbers
//   - error: non-nil if file operations fail
func (s *Store) Append(entries []Entry) ([]uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sequences := make([]uint64, len(entries))
	var lines []byte

	for i := range entries {
		s.meta.SequenceCounter++
		entries[i].Sequence = s.meta.SequenceCounter
		sequences[i] = s.meta.SequenceCounter

		b, marshalErr := json.Marshal(entries[i])
		if marshalErr != nil {
			return nil, marshalErr
		}
		lines = append(lines, b...)
		lines = append(lines, token.NewlineLF...)
		s.entries = append(s.entries, entries[i])
	}

	if appendErr := appendFile(
		entriesPath(s.dir), lines,
	); appendErr != nil {
		return nil, appendErr
	}

	if saveErr := saveJSON(
		metaPath(s.dir), s.meta,
	); saveErr != nil {
		return nil, saveErr
	}

	return sequences, nil
}

// Query returns entries matching types since a sequence.
//
// Parameters:
//   - types: entry types to include (empty = all types)
//   - sinceSequence: entries with sequence > this value
//
// Returns:
//   - []Entry: matching entries in sequence order
func (s *Store) Query(
	types []string, sinceSequence uint64,
) []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	typeSet := make(map[string]bool, len(types))
	for _, t := range types {
		typeSet[t] = true
	}

	var result []Entry
	for _, e := range s.entries {
		if e.Sequence <= sinceSequence {
			continue
		}
		if len(typeSet) > 0 && !typeSet[e.Type] {
			continue
		}
		result = append(result, e)
	}
	return result
}

// RegisterClient adds a client to the registry.
//
// Parameters:
//   - client: client info to register
//
// Returns:
//   - error: non-nil if persistence fails
func (s *Store) RegisterClient(client ClientInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Reject duplicate project names.
	for i := range s.clients {
		if s.clients[i].ProjectName == client.ProjectName {
			return errHub.DuplicateProject(
				client.ProjectName,
			)
		}
	}

	idx := len(s.clients)
	s.clients = append(s.clients, client)
	s.tokenIdx[client.Token] = idx
	return saveJSON(clientsPath(s.dir), s.clients)
}

// ValidateToken checks if a token matches a registered
// client using constant-time comparison.
//
// Parameters:
//   - bearerToken: bearer token to validate
//
// Returns:
//   - *ClientInfo: matching client, or nil if not found
func (s *Store) ValidateToken(bearerToken string) *ClientInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx, ok := s.tokenIdx[bearerToken]
	if !ok || idx >= len(s.clients) {
		return nil
	}
	stored := s.clients[idx].Token
	if subtle.ConstantTimeCompare(
		[]byte(stored), []byte(bearerToken),
	) != 1 {
		return nil
	}
	return &s.clients[idx]
}

// Stats returns current hub statistics.
//
// Returns:
//   - totalEntries: total number of entries
//   - byType: entry count per type
//   - byProject: entry count per origin project
func (s *Store) Stats() (
	uint64, map[string]uint64, map[string]uint64,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	byType := make(map[string]uint64)
	byProject := make(map[string]uint64)

	for _, e := range s.entries {
		byType[e.Type]++
		byProject[e.Origin]++
	}

	return uint64(len(s.entries)), byType, byProject
}
