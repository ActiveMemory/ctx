# Decision: Federation Strategy for Shared Context Hub

Reference: [ctx-wg#30](https://github.com/ActiveMemory/ctx-wg/discussions/30)

---

## Context

The shared context hub ([shared-context-hub.md](shared-context-hub.md))
starts as a single hub-and-spoke topology. Multiple users work across
projects, each publishing knowledge to the hub. We need:

1. A **master** hub that collects from all ctx clients
2. Master pushes updates to all connected clients
3. If master dies, **another hub takes over automatically**
4. Clients redirect to the new master

This is a low-volume, append-only workload — tens of entries per day.
Entries are immutable once published, each with a UUID.

## Candidates Evaluated

### 1. Raft (Full Consensus)

- Leader election, log compaction, snapshotting
- Strong total ordering across all nodes
- Go libraries: `hashicorp/raft`, `etcd/raft`
- **Overkill for data consensus** — we don't need linearizable
  writes for append-only entries

### 2. CRDT (G-Set / Grow-Only Log)

- Entries are a natural G-Set — set-union merge, no conflicts
- Mathematically proven convergence with zero coordination
- ~50 lines of Go, no library needed
- **No leader election** — doesn't solve the "who is master" problem

### 3. Simple Log Replication (Primary-Secondary)

- Primary assigns sequences, secondaries pull
- **No automatic failover** — requires manual promotion

### 4. Gossip Protocol

- Anti-entropy, membership management, failure detection
- **Overkill for low volume** — designed for epidemic dissemination

### 5. No Consensus — Append-Only + UUIDs

- Simple peer sync, every hub is autonomous
- **No master** — doesn't match the master/follower model users need

## Decision

**Raft-lite for leader election + G-Set sync for data replication.**

Use `hashicorp/raft` **only** for agreeing on who the master is.
Data replication stays simple — append-only entries with UUID keys,
set-union merge, no conflicts by construction.

This gives us:
- Automatic master election and failover (Raft)
- Simple, conflict-free data sync (G-Set)
- No full Raft log — entries don't go through Raft consensus

---

## Architecture

### Master/Follower Model

```
Normal operation:              After master failure:

 Master (Hub-1)                Hub-1 X (down)
  ↑  ↑  ↑                     
ctx-A ctx-B ctx-C              Hub-2 elected master (Raft)
  ↓  ↓                         ↑  ↑  ↑
Hub-2  Hub-3 (followers)       ctx-A ctx-B ctx-C
                                ↓
                               Hub-3 (follower)
```

- **Master** accepts writes from all clients, replicates to followers
- **Followers** receive replicated entries, serve read-only queries
- **Clients** know the peer list; on master failure, they reconnect
  to the new master automatically

### What Raft Does (and Doesn't Do)

| Raft handles | Raft does NOT handle |
|-------------|---------------------|
| Leader election | Data consensus |
| Failure detection | Entry ordering |
| Leader step-down | Entry storage |
| Cluster membership | Client connections |

Raft is a thin layer for one purpose: **everyone agrees who the
master is**. The actual entry data flows through simple gRPC
replication, not through Raft's log.

### Data Replication

Entries are still append-only with UUIDs — no conflicts possible.

```
Client publishes → Master hub
  Master assigns sequence number
  Master appends to local entries.jsonl
  Master pushes to followers via gRPC stream
  Followers append to their local entries.jsonl
```

If a follower misses entries (was offline), it catches up via
sequence-based sync — same as the client Sync RPC.

---

## Daemon Mode

Hubs must run as **long-lived background processes** to participate
in leader election and maintain cluster state.

`ctx serve --shared` becomes a daemon:

```bash
# Start hub daemon (foreground, for dev/testing)
ctx serve --shared --port 9900 \
  --peers hub-2:9900,hub-3:9900

# Start as background daemon
ctx serve --shared --port 9900 \
  --peers hub-2:9900,hub-3:9900 \
  --daemon

# Stop daemon
ctx serve --stop

# Check daemon and cluster status
ctx hub status
# Role:    master
# Peers:
#   hub-2  follower  connected  last_sync: 2s ago
#   hub-3  follower  connected  last_sync: 5s ago
# Entries: 312
# Uptime:  4h 23m
```

### Deployment

Users deploy the daemon however they prefer:

- `--daemon` flag (backgrounds the process, writes PID file)
- systemd / launchd service
- Docker container
- Screen / tmux session

The CLI (`ctx connect`, `ctx fleet`) talks to the local or remote
daemon via gRPC — same as the existing design.

---

## Failover Flow

```
1. Hub-1 (master) goes down
2. Raft detects missing heartbeat (configurable timeout)
3. Hub-2 and Hub-3 hold election
4. Hub-2 wins (highest priority or Raft term)
5. Hub-2 becomes master, starts accepting writes
6. Clients detect connection failure to Hub-1
7. Clients try next peer in list → connect to Hub-2
8. Hub-2 continues serving — no data loss (entries replicated)
9. Hub-1 comes back online → joins as follower
10. Hub-1 syncs missed entries from Hub-2
```

### Client Reconnection

Clients maintain an ordered peer list from registration:

```yaml
# .context/.connect.enc (decrypted view)
server: grpcs://hub-1:9900
peers:
  - grpcs://hub-1:9900
  - grpcs://hub-2:9900
  - grpcs://hub-3:9900
```

On connection failure, the client tries peers in order. The connected
hub responds with the current master address if the client connected
to a follower.

---

## CLI

```bash
# Start a 3-node cluster
# On hub-1:
ctx serve --shared --port 9900 --peers hub-2:9900,hub-3:9900 --daemon
# On hub-2:
ctx serve --shared --port 9900 --peers hub-1:9900,hub-3:9900 --daemon
# On hub-3:
ctx serve --shared --port 9900 --peers hub-1:9900,hub-2:9900 --daemon

# Check cluster state from any hub
ctx hub status

# Add a peer at runtime
ctx hub peer add grpcs://hub-4:9900

# Remove a peer
ctx hub peer remove hub-4

# Force leadership transfer (graceful)
ctx hub stepdown

# Client registration (unchanged from single-hub)
ctx connect register grpcs://hub-1:9900 --token ctx_adm_xxx
# Client automatically learns peer list from hub
```

---

## Entry Identity

Each entry gets:
- `id` — UUID (globally unique)
- `origin` — which project published it
- `sequence` — assigned by master (total order)
- `created_at` — wall-clock time

Since only the master assigns sequence numbers, total ordering is
preserved across the cluster. Followers mirror the master's sequence.

---

## Consequences

- Federation uses Raft for leader election only, not data consensus
- Data replication is simple sequence-based sync (append-only, no
  conflicts)
- `ctx serve --shared` gains `--daemon` and `--peers` flags
- `hashicorp/raft` becomes a dependency (scoped to `internal/hub/`)
- Clients auto-reconnect to new master on failover
- Single-hub deployments still work — Raft with one node elects
  itself immediately
- Total ordering preserved (master assigns all sequences)
