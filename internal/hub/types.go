//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
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

// Store is an append-only JSONL storage backend for entries.
//
// All writes are serialized via a mutex. Entries are appended
// to a single JSONL file. Client registry and metadata are
// stored as separate JSON files.
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

// Server is the shared context hub gRPC server.
//
// It implements Register, Publish, Sync, Listen, and Status
// RPCs backed by an append-only [Store].
//
// Fields:
//   - store: append-only storage backend
//   - adminToken: token required for Register RPC
//   - grpc: underlying gRPC server
//   - listeners: fan-out broadcaster for Listen streams
type Server struct {
	store      *Store
	adminToken string
	grpc       *grpc.Server
	listeners  *fanOut
}

// fanOut manages real-time entry broadcast to listeners.
//
// Fields:
//   - mu: serializes subscribe/unsubscribe/broadcast
//   - subs: active listener channels
type fanOut struct {
	mu   sync.Mutex
	subs map[chan []Entry]struct{}
}

// RegisterRequest is the input for the Register RPC.
//
// Fields:
//   - AdminToken: admin token from server startup
//   - ProjectName: this project's identifier
type RegisterRequest struct {
	AdminToken  string `json:"admin_token"`
	ProjectName string `json:"project_name"`
}

// RegisterResponse is the output of the Register RPC.
//
// Fields:
//   - ClientID: assigned client identifier
//   - ClientToken: token for future RPCs
type RegisterResponse struct {
	ClientID    string `json:"client_id"`
	ClientToken string `json:"client_token"`
}

// PublishRequest is the input for the Publish RPC.
//
// Fields:
//   - Entries: entries to publish
type PublishRequest struct {
	Entries []PublishEntry `json:"entries"`
}

// PublishEntry is a single entry in a PublishRequest.
//
// Fields:
//   - ID: entry UUID
//   - Type: entry type
//   - Content: markdown text
//   - Origin: source project
//   - Author: optional author
//   - Timestamp: Unix epoch seconds
type PublishEntry struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Origin    string `json:"origin"`
	Author    string `json:"author,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// PublishResponse is the output of the Publish RPC.
//
// Fields:
//   - Sequences: assigned sequence numbers
type PublishResponse struct {
	Sequences []uint64 `json:"sequences"`
}

// SyncRequest is the input for the Sync RPC.
//
// Fields:
//   - Types: entry types to sync (empty = all)
//   - SinceSequence: return entries after this sequence
type SyncRequest struct {
	Types         []string `json:"types"`
	SinceSequence uint64   `json:"since_sequence"`
}

// ListenRequest is the input for the Listen RPC.
//
// Fields:
//   - Types: entry types to receive (empty = all)
//   - SinceSequence: start from this sequence
type ListenRequest struct {
	Types         []string `json:"types"`
	SinceSequence uint64   `json:"since_sequence"`
}

// EntryMsg is a wire-format entry for streaming RPCs.
//
// Fields:
//   - ID: entry UUID
//   - Type: entry type
//   - Content: markdown text
//   - Origin: source project
//   - Author: optional author
//   - Timestamp: Unix epoch seconds
//   - Sequence: hub-assigned sequence
type EntryMsg struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Origin    string `json:"origin"`
	Author    string `json:"author,omitempty"`
	Timestamp int64  `json:"timestamp"`
	Sequence  uint64 `json:"sequence"`
}

// StatusResponse is the output of the Status RPC.
//
// Fields:
//   - TotalEntries: total number of entries
//   - ConnectedClients: active listener count
//   - EntriesByType: entry count per type
//   - EntriesByProject: entry count per origin project
type StatusResponse struct {
	TotalEntries     uint64            `json:"total_entries"`
	ConnectedClients uint32            `json:"connected_clients"`
	EntriesByType    map[string]uint64 `json:"entries_by_type"`
	EntriesByProject map[string]uint64 `json:"entries_by_project"`
}

// Client is a gRPC client for the shared context hub.
//
// Fields:
//   - conn: underlying gRPC connection
//   - token: bearer token for authenticated RPCs
type Client struct {
	conn  *grpc.ClientConn
	token string
}

// Cluster wraps a Raft node for leader election only.
//
// Fields:
//   - raftNode: the underlying Raft instance
//   - transport: network transport for Raft communication
type Cluster struct {
	raftNode  *raft.Raft
	transport *raft.NetworkTransport
}

// jsonCodec is a gRPC codec using JSON encoding instead of
// protobuf. This allows plain Go structs as RPC messages
// without generated protobuf code.
type jsonCodec struct{}

// codecName is the gRPC content-subtype for the JSON codec.
const codecName = "json"

// Marshal encodes v as JSON.
//
// Parameters:
//   - v: value to encode
//
// Returns:
//   - []byte: JSON-encoded bytes
//   - error: non-nil if encoding fails
func (jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes JSON data into v.
//
// Parameters:
//   - data: JSON bytes to decode
//   - v: target to decode into
//
// Returns:
//   - error: non-nil if decoding fails
func (jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Name returns the codec content-subtype.
//
// Returns:
//   - string: "json"
func (jsonCodec) Name() string { return codecName }
