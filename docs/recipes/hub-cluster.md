---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: HA Cluster
icon: lucide/layers
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: High-Availability Cluster

Run **multiple** hub nodes with Raft-based leader election for
redundancy. Any follower can take over if the leader dies.

This recipe assumes you've read the
[`ctx` Hub overview](hub-overview.md) and the
[Multi-machine setup](hub-multi-machine.md). HA only makes
sense in the "small trusted team" story; a personal
cross-project brain on one workstation does not need three Raft
peers.

!!! warning "Raft-Lite"
    `ctx` uses Raft **only for leader election**, not for data
    consensus. Entry replication happens via sequence-based gRPC
    sync on the append-only JSONL store. This is simpler than full
    Raft log replication and is possible because the store is
    append-only and clients are idempotent. **The implication**:
    a write accepted by the leader is durable on the leader
    immediately; followers catch up asynchronously. If the leader
    crashes **between** accepting a write and replicating it,
    that write can be lost. Do not use the hub as a bank ledger.

## Topology

A minimum HA cluster is **three** nodes. Two is worse than one:
it doubles failure probability without providing quorum.

```
         +-------------+
         |  client(s)  |
         +------+------+
                |
    +-----------+-----------+
    |           |           |
+---v---+   +---v---+   +---v---+
| hub A |   | hub B |   | hub C |
| :9900 |   | :9900 |   | :9900 |  gRPC (clients, data sync)
| :9901 |   | :9901 |   | :9901 |  Raft (leader election)
+-------+   +-------+   +-------+
    ^           ^           ^
    +-----------+-----------+
        Raft (leader election)
        gRPC (data sync)
```

Each node runs two listeners: the hub's gRPC port that clients
dial (`--port`), and the Raft port the other nodes dial
(`--raft-bind`). They are separate addresses; the peer list is
made of **Raft** addresses.

## Step 1: Bootstrap the First Node

```bash
ctx hub start --daemon \
  --port 9900 \
  --raft-bind hub-a.lan:9901 \
  --peers hub-b.lan:9901,hub-c.lan:9901
```

`--raft-bind` is the address this node advertises to the other
two, so it has to be a host they can dial: a bare port
(`:9901`) or a wildcard (`0.0.0.0:9901`) is rejected at
startup. Every node's `--raft-bind` appears in the other nodes'
`--peers` lists, and each node bootstraps that same set.

The node starts a Raft election as soon as it sees its peers.
Until a quorum answers, `ctx hub status` reports
`Leader: unknown (election in progress)` — expected while the
other nodes are still coming up.

## Step 2: Start the Other Nodes

On `hub-b.lan`:

```bash
ctx hub start --daemon \
  --port 9900 \
  --raft-bind hub-b.lan:9901 \
  --peers hub-a.lan:9901,hub-c.lan:9901
```

On `hub-c.lan`:

```bash
ctx hub start --daemon \
  --port 9900 \
  --raft-bind hub-c.lan:9901 \
  --peers hub-a.lan:9901,hub-b.lan:9901
```

After a few seconds, one node wins the election and becomes the
**leader**. The other two are followers.

## Step 3: Verify Cluster State

From any node:

```bash
ctx hub status
```

Expected output on the node that won the election:

```
Role: Leader
Leader: hub-a.lan:9901
Entries: 1248  Peers: 2
```

and on either of the others:

```
Role: Follower
Leader: hub-a.lan:9901
Entries: 1248  Peers: 2
```

The leader is named by its Raft address, which is what the
cluster agrees on. Clients still dial the hub port.

`Peers:` counts the servers in the committed Raft configuration
other than the one answering, so a three-node cluster reports
two from every node. If the line reads
`Leader: unknown (election in progress)`, Raft has no leader for
the current term: either the election is still running, or the
node you asked cannot see a quorum.

## Step 4: Register Clients with Failover Peers

The `ctx hub *` commands above run on the hub nodes themselves and
don't need a project. The `ctx connection *` commands below are
different: they live inside a project (the encrypted hub config is
stored at `.context/.connect.enc`), so you have to tell `ctx` which
project first.

When registering a client, give it the **full peer list**:

```bash
# In the project directory on the client:
ctx connection register hub-a.lan:9900 \
  --token ctx_adm_... \
  --peers hub-b.lan:9900,hub-c.lan:9900
```

If the leader becomes unreachable, the client reconnects to the
next peer. Followers redirect to the current leader, so writes
always land on the right node.

## Runtime Membership Changes

!!! warning "Not Wired Yet"
    `ctx hub peer add`, `ctx hub peer remove` and
    `ctx hub stepdown` print a confirmation and return: none of
    them reaches the Raft node. Membership changes today mean
    restarting the affected nodes with new `--peers` lists, and
    a leader handoff means stopping the leader and letting the
    remaining nodes elect. `ctx hub status` reports the outcome
    either way, so you can see what the cluster actually did.

Add a new peer without downtime:

```bash
ctx hub peer add hub-d.lan:9901
```

Remove a decommissioned peer:

```bash
ctx hub peer remove hub-c.lan:9901
```

## Planned Maintenance

Before taking a leader offline, hand off leadership:

```bash
ssh hub-a.lan 'ctx hub stepdown'
```

`stepdown` is meant to trigger a new election among the
remaining followers before the leader goes offline. Until it is
wired to the cluster (see the warning above), stop the leader
and let the survivors elect: with a quorum still up, the new
leader appears in `ctx hub status` within a couple of seconds.

## Failure Modes at a Glance

| Event                       | What happens                                 |
|-----------------------------|----------------------------------------------|
| Leader crashes              | New election; clients reconnect to new leader |
| Follower crashes            | No write impact; catches up on restart        |
| Network partition (majority) | Majority side keeps serving; minority read-only |
| Network partition (split)   | No quorum; all nodes read-only                |
| Disk full on leader         | Writes rejected; read traffic continues      |

For the full list, see
[Hub failure modes](../operations/hub-failure-modes.md).

## See Also

- [Multi-machine recipe](hub-multi-machine.md): single-node
  deployment
- [Hub operations](../operations/hub.md): backup and
  maintenance
- [Hub security model](../security/hub.md): TLS, tokens
