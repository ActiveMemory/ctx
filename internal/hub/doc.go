//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub implements the ctx Hub: a gRPC server
// that fans structured entries (decisions, learnings,
// conventions, tasks) across multiple ctx projects,
// plus the client primitives those projects use to
// talk to it.
//
// # Architecture
//
// The package layers four concerns:
//
//   - Storage ([Store]) -- append-only JSONL with
//     sequence numbers and per-client tokens.
//   - Transport ([Server]) -- gRPC Register / Publish
//     / Sync / Listen / Status RPCs.
//   - Cluster ([Cluster]) -- HashiCorp Raft for leader
//     election only (see Raft-Lite below).
//   - Client ([Client]) -- connection registration,
//     sync catch-up, push streaming, and ordered-peer
//     failover.
//
// Supporting pillars:
//
//   - Auth ([GenerateAdminToken],
//     [GenerateClientToken]) -- bearer-token
//     authentication on every RPC.
//   - Validate ([ValidateEntry]) -- entry schema
//     enforcement and provenance normalization.
//   - Fan-out -- internal broadcaster delivers each
//     new entry to all live Listen subscribers.
//
// # Storage Model
//
// The store is append-only JSONL under a hub data
// directory: entries.jsonl (one [Entry] per line),
// clients.json (registered tokens and filters), and
// meta.json (schema version and admin token hash).
// Sequence numbers make replication and resume
// strictly idempotent.
//
// # Raft-Lite
//
// The package embeds HashiCorp Raft for leader
// election only -- never for data consensus. Entry
// replication uses the sequence-based gRPC sync.
// Writes are durable on the leader at acceptance;
// followers catch up asynchronously.
//
// # Trust Model
//
// Every holder of a client token is trusted. Origin
// is self-asserted; there is no per-user attribution.
// The hub serves single-developer and small-team
// shapes, not public multi-tenant deployments.
//
// # Concurrency
//
// [Store] guards its indexes and appender with a
// single mutex. Listen streams subscribe to a
// fan-out channel; slow subscribers are dropped
// rather than blocking publishers.
//
// # Encryption
//
// Client-side connection state is encrypted at rest
// via AES-256-GCM using the same per-machine key
// that protects [internal/pad].
package hub
