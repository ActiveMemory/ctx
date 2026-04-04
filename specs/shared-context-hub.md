# Spec: Shared Context Hub (`ctx serve --shared` + `ctx connect`)

Share knowledge between ctx instances across projects via a centralized
gRPC hub. Projects publish decisions, learnings, and conventions to a
shared knowledge base; other projects receive them in real-time and use
them alongside local context for better-informed agent work.

---

## Problem

Each ctx instance is an island. A team running five microservices has
five separate `.context/` directories. When project-alpha learns "gRPC
deadline must be set on client side," project-beta discovers the same
lesson independently — days or weeks later, often the hard way.

There is no mechanism to share knowledge across projects. Copy-pasting
entries between `.context/` directories is manual, error-prone, and
doesn't scale.

## Solution

A **hub-and-spoke** architecture where a shared server aggregates
published entries from multiple ctx instances and streams them to
subscribers in real-time.

```
[project-alpha ctx] ──gRPC──→ [ctx hub :9900] ←──gRPC── [project-beta ctx]
[project-gamma ctx] ──gRPC──┘                  └──gRPC── [project-delta ctx]
```

**Key principles:**

1. **Append-only** — published entries are never modified or deleted
2. **Curated sharing** — each project chooses what to publish (`--share`)
3. **Local authority** — shared knowledge is informational, not imposed
4. **Explicit action** — nothing enters or leaves local context without
   the user's explicit intent

---

## Core Concepts

### Entry — the unit of sharing

Every published piece of context is an Entry:

```go
type Entry struct {
    ID        string    // UUID, globally unique
    Type      string    // entry.Decision, entry.Learning, entry.Convention, entry.Task
    Content   string    // the actual text (markdown)
    Origin    string    // project name that published it
    Author    string    // optional, who wrote it
    Timestamp time.Time // when it was published
    Sequence  uint64    // monotonic, assigned by hub
}
```

- Entries are **append-only** — once published, never modified or deleted
- Each entry gets a **sequence number** from the hub (monotonically
  increasing global counter)
- Clients track their last-seen sequence to know where to resume

### Subscription — what a client cares about

```go
type Subscription struct {
    Types []string // e.g., ["decision", "learning"]
}
```

Type-based filter. The server only streams entries matching the
client's subscription.


---

## gRPC Service Definition

### Proto Service

```protobuf
syntax = "proto3";
package ctx.hub.v1;

option go_package = "github.com/ActiveMemory/ctx/internal/hub/hubpb";

service CtxHub {
  // Auth — one-time registration
  rpc Register(RegisterRequest) returns (RegisterResponse);

  // Publish entries to the hub
  rpc Publish(PublishRequest) returns (PublishResponse);

  // Initial sync — pull all entries matching subscription since a sequence
  rpc Sync(SyncRequest) returns (stream Entry);

  // Incremental updates — long-lived server stream
  rpc Listen(ListenRequest) returns (stream Entry);

  // Query hub state
  rpc Status(StatusRequest) returns (StatusResponse);
}

message RegisterRequest {
  string admin_token = 1;   // admin token from server startup
  string project_name = 2;  // this project's identifier
}

message RegisterResponse {
  string client_id = 1;     // assigned client identifier
  string client_token = 2;  // token for future RPCs
}

message PublishRequest {
  repeated Entry entries = 1;
}

message PublishResponse {
  repeated uint64 sequences = 1; // assigned sequence numbers
}

message SyncRequest {
  repeated string types = 1;  // entry types to sync
  uint64 since_sequence = 2;  // 0 for full sync
}

message ListenRequest {
  repeated string types = 1;
  uint64 since_sequence = 2;
}

message StatusRequest {}

message StatusResponse {
  uint64 total_entries = 1;
  uint32 connected_clients = 2;
  map<string, uint64> entries_by_type = 3;
  map<string, uint64> entries_by_project = 4;
}

message Entry {
  string id = 1;
  string type = 2;
  string content = 3;
  string origin = 4;
  string author = 5;
  int64 timestamp = 6;      // Unix epoch seconds
  uint64 sequence = 7;
}
```

### Message Flow

```
CLIENT                              HUB
  │                                  │
  │── Register(token, project) ────→ │  one-time setup
  │←── RegisterResponse(client_id) ──│
  │                                  │
  │── Sync(types, since_seq=0) ────→ │  initial pull
  │←── Entry stream ─────────────────│  all matching entries
  │←── Entry stream ─────────────────│
  │←── (stream closes) ─────────────│
  │                                  │
  │── Listen(types, since_seq=N) ──→ │  long-lived stream
  │←── Entry (when available) ───────│  real-time updates
  │←── Entry (when available) ───────│
  │         ...                      │
  │                                  │
  │── Publish(entries) ────────────→ │  push local entries
  │←── PublishResponse(sequences) ──│
  │                                  │  hub fans out to
  │                                  │  other listeners
```

### Authentication

- **Register** — client presents admin token + project name. Hub returns
  a client-specific token. One-time operation.
- **All other RPCs** — client token as gRPC metadata:
  `authorization: Bearer <client-token>`
- Server validates via unary/stream interceptor before any handler runs.
- TLS encryption via `--tls-cert` and `--tls-key` flags.

---

## CLI Commands

### Server Side

```bash
# Start the shared hub
ctx serve --shared --port 9900
ctx serve --shared --port 9900 --tls-cert cert.pem --tls-key key.pem

# First run generates an admin token, printed to stderr:
#   Hub started on :9900
#   Admin token: ctx_adm_7f3a...  (save this, shown only once)
```

The hub stores entries in an append-only log. Storage options:

- **v1**: Single JSONL file (`hub-data/entries.jsonl`) — simple, good
  for small-to-medium deployments
- **Future**: SQLite for indexed queries and better concurrent access

### Client Side

```bash
# 1. Register with the hub (one-time)
ctx connect register grpcs://hub.example.com:9900 --token ctx_adm_7f3a...
#   → stores encrypted config in .context/.connect.enc
#   → registers this project with the hub

# 2. Set what you want to receive
ctx connect subscribe decisions learnings
#   → updates local subscription config

# 3. Initial sync (pull all matching entries from hub)
ctx connect sync
#   → streams all matching entries
#   → writes to .context/shared/
#   → records last-seen sequence

# 4. Listen for real-time updates (long-lived)
ctx connect listen
#   → gRPC server-stream, writes new entries to .context/shared/
#   → reconnects automatically on disconnect
#   → ctrl-c to stop

# 5. Publish a local entry to the hub
ctx add decision "Use UTC timestamps everywhere" --share
#   → adds to local DECISIONS.md AND pushes to hub

ctx connect publish --entry decision:5
#   → pushes existing local entry #5 to hub

ctx connect publish --new
#   → pushes all entries created since last publish

# 6. Check connection status
ctx connect status
#   → server, connected, last sync, subscription, entry counts
```

---

## Local File Layout

Shared entries land in a **separate directory**, never mixed with
local context:

```
.context/
├── DECISIONS.md              ← local (this project's decisions)
├── LEARNINGS.md              ← local (this project's learnings)
├── CONVENTIONS.md            ← local (this project's conventions)
├── shared/                   ← from the hub (read-only)
│   ├── decisions.md          ← shared decisions, append-only
│   ├── learnings.md          ← shared learnings, append-only
│   ├── conventions.md        ← shared conventions, append-only
│   └── .sync-state.json     ← last sequence, subscription config
```

### Shared File Format

Each shared file uses the same markdown format as local files, with
origin tags:

```markdown
## [2026-03-14] Use UTC timestamps everywhere

**Origin**: project-alpha

All timestamps in APIs, databases, and logs must use UTC.
Timezone conversion happens only at the UI layer.

---

## [2026-03-14] Never mock the database in integration tests

**Origin**: project-beta

Mocked tests passed but prod migration failed. Integration tests
must hit a real database instance.

---
```

---

## Agent Integration


### Usage

```bash
ctx agent                     # local only (default, unchanged)
ctx agent --include-shared    # local + shared knowledge
ctx agent --include-shared --shared-budget 3000  # custom shared budget
```

### How Agents Use Shared Knowledge

No special logic needed. Context files already influence agent behavior
by being present in the context window. Shared knowledge works the
same way — it's additional context the agent reads and weighs
alongside local context.

When the agent loads shared knowledge, it sees:

```markdown
## Shared Knowledge (from hub)
- [decision] project-alpha: "Use UTC timestamps everywhere"
- [learning] project-beta: "Never mock the database in integration tests"
- [convention] project-alpha: "All APIs return JSON envelope format"
```

The agent uses these naturally when:

- **Making decisions** → shared decisions inform cross-project consistency
- **Writing code** → shared conventions guide patterns
- **Avoiding mistakes** → shared learnings prevent repeating others' bugs

---

## Connection Lifecycle

```
1. ctx connect register grpcs://hub:9900 --token xxx
   └── one-time: registers project, stores encrypted config

2. ctx connect subscribe decisions learnings
   └── declares what entry types to receive

3. ctx connect sync
   └── pulls all matching entries from hub → .context/shared/

4. ctx connect listen
   └── long-lived stream for incremental updates
   └── auto-reconnects on disconnect
   └── writes new entries to .context/shared/ as they arrive

5. ctx add decision "..." --share
   └── normal local add + publish to hub in one step
```

---

## Hub Storage

### v1: JSONL Append-Only Log

```
hub-data/
├── entries.jsonl       ← one JSON object per line
├── clients.json        ← registered clients and tokens
└── meta.json           ← sequence counter, hub metadata
```

Each line in `entries.jsonl`:

```json
{"id":"uuid","type":"decision","content":"...","origin":"project-alpha","author":"","timestamp":1710422400,"sequence":1}
```
