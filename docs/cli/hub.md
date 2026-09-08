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

Operator commands for a **`ctx` Hub**: the gRPC server that
fans out decisions, learnings, conventions, and tasks across
projects. Use `ctx hub` to start and stop the server, inspect
cluster state, add or remove peers at runtime, and hand off
leadership before maintenance.

!!! tip "Who Needs This Page"
    You only need `ctx hub` if you are **running** a hub
    server or cluster. For client-side operations (register,
    subscribe, sync, publish, listen), see
    [`ctx connection`](connection.md). For the mental model behind
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
stdout. Save this token; it's required for
[`ctx connection register`](connection.md#ctx-connection-register) in
client projects. Subsequent runs reuse the stored token from
`<data-dir>/admin.token`.

**Default data directory**: `~/.ctx/hub-data/`

#### Daemon Mode

Run the hub as a detached background process:

```bash
ctx hub start --daemon          # Fork to background
ctx hub stop                    # Graceful shutdown
```

The daemon writes a PID file to `<data-dir>/hub.pid`. Stop
the daemon with `ctx hub stop` (see below).

#### Cluster Mode

For high availability, run multiple hubs with Raft-based
leader election. `--raft-bind` is the address this node binds
its Raft transport to and advertises to the others, and
`--peers` lists the `--raft-bind` addresses of the other
nodes — Raft addresses, not hub ports:

```bash
ctx hub start --port 9900 \
  --raft-bind host1:9901 \
  --peers host2:9901,host3:9901
```

`--raft-bind` must name a host a peer can dial. A bare port
(`:9901`) or a wildcard (`0.0.0.0:9901`) is rejected at
startup, because Raft refuses to advertise an address that
does not identify this node to anyone else.

`--raft-bind` on its own — with no `--peers` — runs a
single-node Raft cluster that elects itself. That is the
cheapest way to see the leadership fields of
[`ctx hub status`](#ctx-hub-status) before adding nodes.

To add a node to a cluster that is already running, start it
with `--join` instead of `--peers`: it brings up its Raft
transport, bootstraps nothing, and waits for
[`ctx hub peer add`](#ctx-hub-peer) on the leader to hand it a
configuration. `--join` and `--peers` together are an error —
a node either bootstraps a cluster or joins one.

Raft is used **only** for leader election. Data replication
uses sequence-based gRPC sync on the append-only JSONL log;
there is no multi-node consensus on writes. See the
[HA cluster recipe](../recipes/hub-cluster.md) for the full
setup and the Raft-lite durability caveat.

#### Flags

| Flag          | Description                                       | Default          |
|---------------|---------------------------------------------------|------------------|
| `--port`      | Hub listen port                                   | `9900`           |
| `--data-dir`  | Hub data directory                                | `~/.ctx/hub-data/` |
| `--daemon`    | Run the hub server in the background              | `false`          |
| `--raft-bind` | Raft address this node binds and advertises       | *(none)*         |
| `--peers`     | Comma-separated peer Raft addresses               | *(none)*         |
| `--join`      | Wait to be added by a leader (no bootstrap)       | `false`          |

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
Safe to rerun: if no daemon is running, returns a
"no running hub" error without side effects.

### `ctx hub status`

Show what the hub reports about itself: its role, the current
leader, the entry count and the peer count.

A hub started without `--peers` runs no Raft node, so it has no
leader and no peers to name:

```
Role: Standalone
Entries: 1248
```

A hub started with peers answers from its Raft node. The role is
`Leader` or `Follower`, the leader is the address Raft holds for
the current term, and the peer count is the committed cluster
configuration minus the node answering:

```
Role: Leader
Leader: 10.0.0.5:9901
Entries: 1248  Peers: 2
```

While an election is in progress — or after quorum is lost —
Raft knows no leader, and the line says so:

```
Leader: unknown (election in progress)
```

When the hub has disconnected any slow listeners, the output
gains a `Dropped listeners:` line with the cumulative count.
The line is omitted while that count is zero, so a healthy hub
looks exactly as it did before. See
[Slow Listener Disconnected](../operations/hub-failure-modes.md#slow-listener-disconnected).

**Examples**:

```bash
ctx hub status
```

### `ctx hub peer`

Add or remove peers in the cluster's Raft configuration at
runtime. Useful for scaling up or replacing a decommissioned
node without restarting the leader.

The address is the peer's **Raft** address (its `--raft-bind`),
not its hub port. Membership changes are admin-gated, like
[`ctx hub revoke`](#ctx-hub-revoke): pass `--token` or set
`CTX_HUB_ADMIN_TOKEN`.

Only the leader can change the configuration. Run the command
against the leader — [`ctx hub status`](#ctx-hub-status) names
it — or the hub answers `not the leader`.

A node being added must already be running with
`--raft-bind <its address> --join`, so that it is waiting for a
configuration instead of bootstrapping one of its own.

**Examples**:

```bash
# On the new node:
ctx hub start --daemon --port 9900 \
  --raft-bind host4:9901 --join

# On the leader:
ctx hub peer add host4:9901 --token ctx_adm_...
ctx hub peer remove host3:9901 --token ctx_adm_...
```

### `ctx hub stepdown`

Transfer leadership to another node gracefully. Triggers a
new election among the remaining followers before the current
leader steps down. Use before taking the leader offline for
maintenance.

Admin-gated like [`ctx hub peer`](#ctx-hub-peer), and
leader-only: a follower answers `not the leader` instead of
reporting a transfer that did not happen. The confirmation
prints after the transfer returns;
[`ctx hub status`](#ctx-hub-status) names the node that won.

**Examples**:

```bash
ctx hub stepdown --token ctx_adm_...
CTX_HUB_ADMIN_TOKEN=ctx_adm_... ctx hub stepdown
```

### See Also

- [`ctx connection`](connection.md): client-side commands
  (register, subscribe, sync, publish, listen)
- [`ctx` Hub overview](../recipes/hub-overview.md): mental
  model and user stories
- [`ctx` Hub: Getting Started](../recipes/hub-getting-started.md)
- [Hub operations](../operations/hub.md): production
  deployment, backup, monitoring
- [Hub failure modes](../operations/hub-failure-modes.md)
- [Hub security model](../security/hub.md)
