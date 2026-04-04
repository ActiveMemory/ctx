# Shared Context Hub — Implementation & Federation

Companion to [shared-context-hub.md](shared-context-hub.md). Contains
package layout, compliance, implementation phases, and future
federation design.

---

## Package Layout

```
internal/
├── hub/
│   ├── doc.go            ← package documentation
│   ├── server.go         ← gRPC server implementation
│   ├── client.go         ← gRPC client (used by ctx connect)
│   ├── store.go          ← JSONL append-only storage
│   ├── auth.go           ← token generation, validation, interceptor
│   ├── types.go          ← Entry, Subscription, ConnectionConfig
│   ├── encrypt.go        ← connection config encryption (reuse notify pattern)
│   └── proto/
│       └── hub.proto     ← gRPC service definition
├── cli/
│   ├── serve/
│   │   └── shared.go     ← ctx serve --shared command
│   └── connect/
│       ├── register.go   ← ctx connect register
│       ├── subscribe.go  ← ctx connect subscribe
│       ├── sync.go       ← ctx connect sync
│       ├── listen.go     ← ctx connect listen
│       ├── publish.go    ← ctx connect publish
│       └── status.go     ← ctx connect status
```

---

## Compliance & Invariants

### Design Invariants Preserved

| Invariant | How preserved |
|-----------|---------------|
| Markdown-on-filesystem | Shared entries stored as .md in .context/shared/ |
| Zero runtime deps (core) | gRPC scoped to `internal/hub/` — not in local-only list |
| Deterministic assembly | Shared budget is additive, same files + budget = same output |
| Human authority | `--share` is explicit, shared knowledge is informational |
| Local-first | Core ctx works without hub; shared is opt-in |
| No telemetry | Hub is self-hosted, no external services |

### Compliance Test Update

The existing `TestNoNetworkImportsInCore` checks a curated list of
`localOnlyPackages` that must not import `net` or `net/http`. The
`internal/hub/` package is not added to that list (same approach as
`internal/notify/`). Core packages (`context`, `config`, `drift`,
`task`, `validation`, `crypto`, `assets`, `index`) remain network-free.

### Security

- **TLS** — `grpcs://` for encrypted transport
- **Token auth** — per-client tokens, validated via gRPC interceptor
- **Encrypted config** — connection config stored with AES-256-GCM
  (same pattern as webhook URL via `internal/crypto`)
- **No sensitive data** — entries are architectural knowledge, not
  secrets. CONSTITUTION.md invariant on secrets still applies.

---

## Implementation Phases

### Phase 1: Foundation

- Proto definition and code generation
- Hub server with JSONL storage
- Register, Publish, Sync RPCs
- `ctx serve --shared` and `ctx connect register/sync/publish`
- Token-based auth with encrypted local storage

### Phase 2: Real-Time

- Listen RPC (server-streaming with fan-out)
- `ctx connect listen` with auto-reconnect
- `ctx add --share` flag integration
- Background listener option

### Phase 3: Agent Integration

- `ctx agent --include-shared` with Tier 6 budget
- Shared file rendering in agent packet
- Scoring shared entries (recency + type relevance)

### Phase 4: Operational

- `ctx connect status` with detailed stats
- Hub-side Status RPC
- Connection health monitoring
- Graceful shutdown and reconnection

---

## Future: Distributed Hub (Federation)

The append-only, sequence-based design enables future hub-to-hub
replication:

```
[Hub-EU] ←──gRPC──→ [Hub-US]
   ↑                    ↑
 clients              clients
```

Each hub maintains its own sequence space. Federation maps remote
sequences to local ones. The entry UUID ensures global deduplication.

### Federation Protocol

- Hubs connect to each other as peers (bidirectional gRPC streams)
- Each hub assigns local sequence numbers to replicated entries
- Entry UUID prevents duplicate ingestion
- Conflict-free: append-only means no write conflicts between hubs
- Partition-tolerant: hubs queue entries during disconnection,
  replay on reconnect using the same since-sequence mechanism
  clients use

### Hub Discovery

- Manual configuration: `ctx serve --shared --peer grpcs://hub-us:9900`
- Future: DNS-based discovery or a lightweight registry

### Consistency Model

- **Eventual consistency** — all hubs converge to the same entry set
- No ordering guarantee across hubs (local sequence only)
- Entry UUID + timestamp allows consumers to sort globally if needed
