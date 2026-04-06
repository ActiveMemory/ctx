//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"sync"
	"time"
)

// Entry is the unit of sharing in the hub.
//
// Every published piece of context is an Entry. Entries are
// append-only — once published, never modified or deleted.
// Each entry gets a monotonically increasing sequence number
// assigned by the hub.
//
// Fields:
//   - ID: UUID, globally unique
//   - Type: entry type (decision, learning, convention, task)
//   - Content: the actual text (markdown)
//   - Origin: project name that published it
//   - Author: optional, who wrote it
//   - Timestamp: when it was published
//   - Sequence: monotonic counter, assigned by hub on publish
type Entry struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Origin    string    `json:"origin"`
	Author    string    `json:"author,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Sequence  uint64    `json:"sequence"`
}

// ClientInfo holds registration data for a connected client.
//
// Fields:
//   - ID: unique client identifier (UUID)
//   - ProjectName: name of the project this client represents
//   - Token: bearer token for authenticating RPCs
type ClientInfo struct {
	ID          string `json:"id"`
	ProjectName string `json:"project_name"`
	Token       string `json:"token"`
}

// Meta holds hub-level metadata persisted alongside the log.
//
// Fields:
//   - SequenceCounter: next sequence number to assign
//   - CreatedAt: when the hub was first started
type Meta struct {
	SequenceCounter uint64    `json:"sequence_counter"`
	CreatedAt       time.Time `json:"created_at"`
}

// Store is an append-only JSONL storage backend for hub entries.
//
// All writes are serialized via a mutex. Entries are appended to
// a single JSONL file. Client registry and metadata are stored
// as separate JSON files.
//
// Fields:
//   - dir: directory where data files live
//   - mu: serializes all reads and writes
//   - meta: hub-level metadata (sequence counter)
//   - clients: registered client tokens
//   - entries: in-memory cache of all entries (append-only)
type Store struct {
	dir     string
	mu      sync.Mutex
	meta    Meta
	clients []ClientInfo
	entries []Entry
}
