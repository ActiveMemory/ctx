//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub centralizes configuration constants for
// the ctx Hub, a gRPC service that synchronizes context
// entries across machines and team members.
//
// The hub is a distributed system with client-server
// registration, entry replication, Raft-based
// clustering, bearer-token authentication, and
// persistent storage. This package defines every
// protocol string, file name, error message, and
// validation limit so that the hub implementation and
// its CLI consumers share a single source of truth.
//
// # gRPC Service Descriptor
//
//   - ServiceName ("ctx.hub.v1.CtxHub") -- the fully
//     qualified gRPC service name
//   - ServicePath -- the service path prefix for
//     method descriptors
//   - MethodRegister, MethodPublish, MethodSync,
//     MethodListen, MethodStatus -- RPC method names
//   - PathRegister, PathPublish, PathSync,
//     PathListen, PathStatus -- full method paths
//   - ProtoFile ("hub.proto") -- virtual proto file
//     name in the service descriptor
//
// # Authentication
//
//   - HeaderAuthorization -- gRPC metadata key for
//     bearer tokens
//   - BearerPrefix ("Bearer ") -- prefix stripped
//     from authorization header values
//   - AdminTokenPrefix ("ctx_adm_") -- prefix for
//     admin tokens
//   - ClientTokenPrefix ("ctx_cli_") -- prefix for
//     client tokens
//   - TokenBytes (32) -- random bytes in generated
//     bearer tokens
//
// # Entry Metadata Fields
//
//   - MetaDisplayName, MetaHost, MetaTool, MetaVia
//     -- JSON field names validated during publish
//   - StructTagJSON -- struct tag key for field name
//     resolution
//
// # Persistence Files
//
//   - FileEntries ("entries.jsonl") -- append-only
//     entry store
//   - FileClients, FileMeta, FileSyncState,
//     FileSyncLock, FileConnect -- client registry,
//     metadata, sync state, lock, and encrypted
//     connection config files
//   - FilePID ("hub.pid"), FileAdminToken,
//     DirHubData -- daemon management files
//   - JSONIndent, LockSentinel, SuffixPluralMD
//     -- formatting and naming helpers
//
// # Raft Cluster Configuration
//
//   - RaftDir ("raft") -- subdirectory for Raft state
//   - RaftTransport ("tcp") -- transport protocol
//   - RaftLogDB ("log.db") -- BoltDB log file
//
// # Validation Limits
//
//   - MaxContentLen (1 MB) -- maximum entry content
//   - MaxMetaFieldLen (256) -- per-field meta cap
//   - MaxMetaTotalLen (2048) -- total meta cap
//   - MetaControlSpaceLow, MetaControlDelete --
//     control character boundaries
//   - ClientIDBytes (16) -- UUID byte length
//
// # Error Messages
//
// Handler and validation error strings are defined as
// constants (ErrInvalidAdminToken, ErrEntryIDRequired,
// ErrMetaControlChar, etc.) so that both server and
// client code can match on them reliably.
//
// # Daemon and Peer Management
//
//   - ArgHub, ArgStart -- re-exec argument tokens
//   - ActionAdd, ActionRemove -- peer action names
//   - RoleFollower, RoleActive -- status role labels
//   - ReplicateInterval (5s) -- follower retry
//     interval
//   - ThrottleHubSync -- daily sync throttle marker
//
// # Why Centralized
//
// Hub constants span the gRPC handler, client library,
// CLI commands, Raft integration, and daemon lifecycle.
// Centralizing them prevents protocol drift between
// server and client and makes the wire format easy to
// audit in one place.
package hub
