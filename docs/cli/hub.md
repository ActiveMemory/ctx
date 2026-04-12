---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hub
icon: lucide/network
---

![ctx](../images/ctx-banner.png)

## `ctx hub`

Operator commands for a **`ctx` Hub** — the gRPC server that
fans out decisions, learnings, conventions, and tasks across
projects. Use `ctx hub` to start and stop the server, inspect
cluster state, add or remove peers at runtime, and hand off
leadership before maintenance.

!!! tip "Who needs this page"
    You only need `ctx hub` if you are **running** a hub
    server or cluster. For client-side operations (register,
    subscribe, sync, publish, listen), see
    [`ctx connect`](connection.md). For the mental model behind
    the hub as a whole, read the
    [`ctx` Hub overview](../recipes/hub-overview.md).

### `ctx hub start`

Start the hub gRPC server.

**Examples**:

```bash
ctx hub start                           # Foreground, default port 9900
ctx hub start --port 8080               # Custom port
ctx hub start --data-dir /srv/ctx-hub   # Custom data directory
```

On first run, generates an **admin token** and prints it to
stdout. Save this token — it's required for
[`ctx connection register`](connection.md#ctx-connect-register) in
client projects. Subsequent runs reuse the stored token from
`<data-dir>/admin.token`.

**Default data directory**: `~/.ctx/hub-data/`

#### Daemon mode

Run the hub as a detached background process:

```bash
ctx hub start --daemon          # Fork to background
ctx hub stop                    # Graceful shutdown
```

The daemon writes a PID file to `<data-dir>/hub.pid`. Stop
the daemon with `ctx hub stop` (see below).

#### Cluster mode

For high availability, run multiple hubs with Raft-based
leader election:

```bash
ctx hub start --port 9900 \
  --peers host2:9901,host3:9901
```

Raft is used **only** for leader election. Data replication
uses sequence-based gRPC sync on the append-only JSONL log —
there is no multi-node consensus on writes. See the
[HA cluster recipe](../recipes/hub-cluster.md) for the full
setup and the Raft-lite durability caveat.

#### Flags

| Flag         | Description                                      | Default          |
|--------------|--------------------------------------------------|------------------|
| `--port`     | Hub listen port                                  | `9900`           |
| `--data-dir` | Hub data directory                               | `~/.ctx/hub-data/` |
| `--daemon`   | Run the hub server in the background             | `false`          |
| `--peers`    | Comma-separated peer addresses for cluster mode  | *(none)*         |

#### Validation

The hub validates every published entry before accepting it:

- **Type** must be one of `decision`, `learning`, `convention`, `task`
- **ID** and **Origin** are required and non-empty
- **Content** size capped at **1 MB** (text-only)
- **Duplicate project registration** is rejected (one token per project)

### `ctx hub stop`

Stop a running hub daemon.

**Examples**:

```bash
ctx hub stop                            # Stop using default data dir
ctx hub stop --data-dir /srv/ctx-hub    # Custom data directory
```

Sends `SIGTERM` to the PID recorded in `<data-dir>/hub.pid`,
waits for in-flight RPCs to drain, and removes the PID file.
Safe to rerun — if no daemon is running, returns a
"no running hub" error without side effects.

### `ctx hub status`

Show cluster status: role, peers, sync state, entry count,
and uptime.

**Examples**:

```bash
ctx hub status
```

### `ctx hub peer`

Add or remove peers from the cluster at runtime. Useful for
scaling up or replacing a decommissioned node without
restarting the leader.

**Examples**:

```bash
ctx hub peer add host2:9901
ctx hub peer remove host2:9901
```

### `ctx hub stepdown`

Transfer leadership to another node gracefully. Triggers a
new election among the remaining followers before the current
leader steps down. Use before taking the leader offline for
maintenance.

**Examples**:

```bash
ctx hub stepdown
```

### See also

- [`ctx connect`](connection.md) — client-side commands
  (register, subscribe, sync, publish, listen)
- [`ctx` Hub overview](../recipes/hub-overview.md) — mental
  model and user stories
- [`ctx` Hub: Getting Started](../recipes/hub-getting-started.md)
- [Hub operations](../operations/hub.md) — production
  deployment, backup, monitoring
- [Hub failure modes](../operations/hub-failure-modes.md)
- [Hub security model](../security/hub.md)
