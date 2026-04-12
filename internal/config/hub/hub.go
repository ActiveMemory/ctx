//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

// gRPC service descriptor.
const (
	// ServiceName is the fully qualified gRPC service name.
	ServiceName = "ctx.hub.v1.CtxHub"
	// ServicePath is the gRPC service path prefix used in
	// method descriptors.
	ServicePath = "/" + ServiceName + "/"
)

// gRPC method names.
const (
	// MethodRegister is the Register RPC method name.
	MethodRegister = "Register"
	// MethodPublish is the Publish RPC method name.
	MethodPublish = "Publish"
	// MethodSync is the Sync RPC method name.
	MethodSync = "Sync"
	// MethodListen is the Listen RPC method name.
	MethodListen = "Listen"
	// MethodStatus is the Status RPC method name.
	MethodStatus = "Status"
)

// Full gRPC method paths (ServicePath + MethodName).
const (
	// PathRegister is the full gRPC path for Register.
	PathRegister = ServicePath + MethodRegister
	// PathPublish is the full gRPC path for Publish.
	PathPublish = ServicePath + MethodPublish
	// PathSync is the full gRPC path for Sync.
	PathSync = ServicePath + MethodSync
	// PathListen is the full gRPC path for Listen.
	PathListen = ServicePath + MethodListen
	// PathStatus is the full gRPC path for Status.
	PathStatus = ServicePath + MethodStatus
)

// Authorization header.
const (
	// HeaderAuthorization is the gRPC metadata key for bearer
	// token authentication.
	HeaderAuthorization = "authorization"
)

// EntryMeta field names used in validation.
const (
	// MetaDisplayName is the JSON field name for display name.
	MetaDisplayName = "display_name"
	// MetaHost is the JSON field name for the host.
	MetaHost = "host"
	// MetaTool is the JSON field name for the tool.
	MetaTool = "tool"
	// MetaVia is the JSON field name for the via field.
	MetaVia = "via"
)

// Context directory layout.
const (
	// DirHub is the subdirectory under .context/ for hub
	// entries and sync state.
	DirHub = "hub"
)

// Replication timing.
const (
	// ReplicateInterval is how often a follower retries
	// connecting to the master for replication.
	ReplicateInterval = 5 // seconds
)

// Token generation.
const (
	// TokenBytes is the number of random bytes in a
	// generated bearer token.
	TokenBytes = 32
)

// Cluster configuration.
const (
	// RaftDir is the subdirectory name for Raft state.
	RaftDir = "raft"
	// RaftTransport is the transport protocol for Raft.
	RaftTransport = "tcp"
	// RaftLogDB is the BoltDB file name for Raft log storage.
	RaftLogDB = "log.db"
)

// gRPC method descriptor metadata.
const (
	// ProtoFile is the virtual proto file name in the service
	// descriptor (no actual .proto file — the service is hand-rolled).
	ProtoFile = "hub.proto"
)

// Peer action names.
const (
	// ActionAdd is the peer add action.
	ActionAdd = "add"
	// ActionRemove is the peer remove action.
	ActionRemove = "remove"
)

// Daemon file names.
const (
	// FilePID is the PID file written by the daemonized hub.
	FilePID = "hub.pid"
	// DirHubData is the subdirectory for hub data files.
	DirHubData = "hub-data"
	// FileAdminToken stores the admin token after first run.
	FileAdminToken = "admin.token"
)

// Status role labels.
const (
	// RoleFollower is the role label for a follower node.
	RoleFollower = "Follower"
	// RoleActive is the role label for an active node.
	RoleActive = "Active"
)

// Address formatting.
const (
	// FmtPort is the format string for a port-only address.
	FmtPort = ":%d"
	// FmtFlagPrefix is the prefix for long-form CLI flags.
	FmtFlagPrefix = "--"
)

// Daemon re-exec argument tokens.
const (
	// ArgHub is the subcommand name when re-execing the
	// hub binary in daemon mode.
	ArgHub = "hub"
	// ArgStart is the start subcommand for daemon re-exec.
	ArgStart = "start"
)

// Throttle identifiers.
const (
	// ThrottleHubSync is the daily throttle marker for hub
	// sync operations.
	ThrottleHubSync = "hub-sync"
)

// Persistence file names.
const (
	// FileEntries is the append-only entries JSONL file.
	FileEntries = "entries.jsonl"
	// FileClients is the client registry JSON file.
	FileClients = "clients.json"
	// FileMeta is the hub metadata JSON file.
	FileMeta = "meta.json"
	// FileSyncState is the sync state JSON file.
	FileSyncState = ".sync-state.json"
	// FileSyncLock is the lock file to prevent concurrent
	// syncs.
	FileSyncLock = ".sync.lock"
	// FileConnect is the encrypted connection config file.
	FileConnect = ".connect.enc"
	// JSONIndent is the indentation string for JSON marshaling.
	JSONIndent = "  "
	// LockSentinel is the content written to lock files.
	LockSentinel = "lock"
	// SuffixPluralMD is the suffix for typed hub markdown
	// filenames (e.g. "decisions.md").
	SuffixPluralMD = "s.md"
)

// Token prefixes.
const (
	// AdminTokenPrefix is the prefix for admin tokens.
	AdminTokenPrefix = "ctx_adm_" //nolint:gosec // prefix, not a credential
	// ClientTokenPrefix is the prefix for client tokens.
	ClientTokenPrefix = "ctx_cli_" //nolint:gosec // prefix, not a credential
)

// Bearer authentication.
const (
	// BearerPrefix is the prefix stripped from the authorization
	// header value.
	BearerPrefix = "Bearer "
)

// Handler error messages.
const (
	// ErrInvalidAdminToken is the gRPC error for invalid admin token.
	ErrInvalidAdminToken = "invalid admin token"
	// ErrProjectNameRequired is the gRPC error for missing project name.
	ErrProjectNameRequired = "project_name required"
	// ErrMissingMetadata is the gRPC error for missing metadata.
	ErrMissingMetadata = "missing metadata"
	// ErrMissingToken is the gRPC error for missing auth token.
	ErrMissingToken = "missing token"
	// ErrInvalidToken is the gRPC error for invalid auth token.
	ErrInvalidToken = "invalid token"
)

// StructTagJSON is the struct tag key used by types.go for
// JSON field name resolution during meta validation.
const StructTagJSON = "json"

// Validation size limits.
const (
	// MaxContentLen is the maximum entry content size (1 MB).
	MaxContentLen = 1 << 20
	// MaxMetaFieldLen is the per-field size cap for
	// EntryMeta fields.
	MaxMetaFieldLen = 256
	// MaxMetaTotalLen caps the sum of all Meta field
	// lengths to prevent abuse via many nearly-full fields.
	MaxMetaTotalLen = 2048
	// MetaControlSpaceLow is the lowest printable ASCII
	// byte; anything below this (except tab) is a control
	// character.
	MetaControlSpaceLow = 0x20
	// MetaControlDelete is the DEL character.
	MetaControlDelete = 0x7f
)

// Client ID size.
const (
	// ClientIDBytes is the byte length of generated client
	// UUIDs (hex-encoded to 32 chars).
	ClientIDBytes = 16
)

// Validation error messages.
const (
	// ErrEntryIDRequired is the gRPC error for missing entry ID.
	ErrEntryIDRequired = "entry ID required"
	// ErrInvalidEntryType is the gRPC error format for invalid
	// entry types.
	ErrInvalidEntryType = "invalid entry type %q"
	// ErrEntryOriginRequired is the gRPC error for missing origin.
	ErrEntryOriginRequired = "entry origin required"
	// ErrEntryContentOversize is the gRPC error for oversized content.
	ErrEntryContentOversize = "entry content exceeds 1MB limit"
	// ErrMetaFieldOversize is the gRPC error format for an
	// oversized meta field.
	ErrMetaFieldOversize = "meta.%s exceeds %d bytes"
	// ErrMetaTotalOversize is the gRPC error format for oversized
	// total meta.
	ErrMetaTotalOversize = "meta total exceeds %d bytes"
	// ErrMetaControlChar is the gRPC error format for control
	// characters in meta fields.
	ErrMetaControlChar = "meta.%s contains control character"
)
