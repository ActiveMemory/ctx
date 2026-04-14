//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub implements the ctx Hub: a gRPC server that fans
// structured entries (decisions, learnings, conventions, tasks)
// across multiple ctx projects and the client primitives those
// projects use to talk to it.
//
// The hub is deliberately small in scope. It is not a wiki, an
// audit log, or a multi-tenant service. It is a fan-out channel
// for the entries that should travel between trusted projects
// without dragging the rest of `.context/` along.
//
// # Architecture at a Glance
//
// The package layers four concerns:
//
//   - Storage     [Store]         append-only JSONL + sequence
//     numbers + per-client tokens
//   - Transport   [Server]        gRPC Register / Publish / Sync
//     / Listen / Status RPCs
//   - Cluster     [Cluster]       Raft *for leader election only*
//     (see "Raft-Lite" below)
//   - Client      [Client]        connection registration, sync
//     catch-up, push streaming, and
//     ordered-peer failover
//
// Auth, validation, and fan-out broadcast support the four pillars:
//
//   - Auth        [GenerateAdminToken], [GenerateClientToken],
//     bearer-token authentication on every RPC.
//   - Validate    [ValidateEntry] enforces the entry schema and
//     normalizes provenance fields.
//   - Fan-out     internal broadcaster delivers each new entry to
//     all live Listen subscribers without coupling
//     them.
//
// # Storage Model
//
// The store is **append-only JSONL** under a hub data directory:
//
//   - entries.jsonl      one [Entry] per line, monotonically
//     sequence-numbered
//   - clients.json       registered client tokens + per-client
//     subscription filters
//   - meta.json          schema version + admin token hash
//
// Sequence numbers make replication and resume strictly
// idempotent: a follower or a returning client only asks for
// "entries after seq N" and the leader streams the tail. Because
// the log is append-only, there is no "edit" operation to
// reconcile and no conflict resolution to perform.
//
// # Raft-Lite
//
// The package embeds HashiCorp Raft (see [Cluster] and the no-op
// [leaderFSM]) for **leader election only** — never for data
// consensus. Entry replication is performed independently via the
// sequence-based gRPC sync described above.
//
// The trade-off is explicit and documented in
// docs/recipes/hub-cluster.md: writes are durable on the leader
// at the moment they are accepted, and followers catch up
// asynchronously. If the leader crashes between accepting a
// write and replicating it, that write may be lost. We take that
// risk in exchange for a much simpler implementation than full
// Raft log replication, and it is sound because the store is
// append-only and clients are idempotent.
//
// # Trust Model
//
// The hub assumes every holder of a client token is friendly.
// Origin is **self-asserted** by the publishing client; there is
// no per-user attribution and no read ACL beyond subscription
// filters. This is by design — the hub serves "Story 1" (single
// developer, multiple projects) and "Story 2" (small trusted
// team) shapes, not public multi-tenant deployments.
//
// Hostile clients, untrusted networks, and compliance-grade audit
// trails are explicitly out of scope. See docs/security/hub.md
// for the threat-model write-up.
//
// # Concurrency
//
// [Store] guards its in-memory indexes and the appender with a
// single mutex; gRPC handlers serialize through it. Listen
// streams subscribe to a fan-out channel and receive entries in
// publish order; slow subscribers are dropped rather than
// blocking publishers (see [fanOutBuffer]).
//
// # Encryption
//
// Connection state on the client side (the hub address and
// per-client token) is encrypted at rest via the package
// [github.com/ActiveMemory/ctx/internal/crypto] using AES-256-GCM
// with the same per-machine key that protects [internal/pad].
//
// # Related Packages
//
//   - [internal/cli/hub]         server-side CLI
//     (start/stop/status/peer/stepdown)
//   - [internal/cli/connection]  client-side CLI
//     (register/subscribe/sync/listen/
//     publish/status)
//   - [internal/err/hub]         typed error constructors used by
//     this package
//   - [internal/config/hub]      protocol/runtime constants
//     (ports, tokens prefixes, paths)
//
// # Key Exports
//
// Server side: [Store], [Server], [NewServer], [Cluster],
// [NewCluster], [GenerateAdminToken], [GenerateClientToken].
//
// Client side: [Client], [Connect], [Publish], [Sync], [Listen],
// [Status], [Register].
//
// Domain types: [Entry], [EntryMsg], [ClientRecord].
package hub
