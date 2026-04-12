---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: HA cluster
icon: lucide/layers
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: High-availability cluster

Run **multiple** hub nodes with Raft-based leader election for
redundancy. Any follower can take over if the leader dies.

This recipe assumes you've read the
[`ctx` Hub overview](hub-overview.md) and the
[Multi-machine setup](hub-multi-machine.md). HA only makes
sense in the "small trusted team" story — a personal
cross-project brain on one workstation does not need three Raft
peers.

!!! warning "Raft-lite"
    ctx uses Raft **only for leader election**, not for data
    consensus. Entry replication happens via sequence-based gRPC
    sync on the append-only JSONL store. This is simpler than full
    Raft log replication and is possible because the store is
    append-only and clients are idempotent. **The implication**:
    a write accepted by the leader is durable on the leader
    immediately; followers catch up asynchronously. If the leader
    crashes **between** accepting a write and replicating it,
    that write can be lost. Do not use the hub as a bank ledger.

## Topology

A minimum HA cluster is **three** nodes. Two is worse than one —
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
| :9900 |   | :9900 |   | :9900 |
+-------+   +-------+   +-------+
    ^           ^           ^
    +-----------+-----------+
        Raft (leader election)
        gRPC (data sync)
```

## Step 1 — Bootstrap the first node

```bash
ctx hub start --daemon \
  --port 9900 \
  --peers hub-b.lan:9900,hub-c.lan:9900
```

The node starts a Raft election as soon as it sees its peers.

## Step 2 — Start the other nodes

On `hub-b.lan`:

```bash
ctx hub start --daemon \
  --port 9900 \
  --peers hub-a.lan:9900,hub-c.lan:9900
```

On `hub-c.lan`:

```bash
ctx hub start --daemon \
  --port 9900 \
  --peers hub-a.lan:9900,hub-b.lan:9900
```

After a few seconds, one node wins the election and becomes the
**leader**. The other two are followers.

## Step 3 — Verify cluster state

From any node:

```bash
ctx hub status
```

Expected output:

```
role:       leader
peers:      hub-a.lan:9900 (leader)
            hub-b.lan:9900 (follower, in-sync)
            hub-c.lan:9900 (follower, in-sync)
entries:    1248
uptime:     3h42m
```

## Step 4 — Register clients with failover peers

When registering a client, give it the **full peer list**:

```bash
ctx connection register hub-a.lan:9900 \
  --token ctx_adm_... \
  --peers hub-b.lan:9900,hub-c.lan:9900
```

If the leader becomes unreachable, the client reconnects to the
next peer. Followers redirect to the current leader, so writes
always land on the right node.

## Runtime membership changes

Add a new peer without downtime:

```bash
ctx hub peer add hub-d.lan:9900
```

Remove a decommissioned peer:

```bash
ctx hub peer remove hub-c.lan:9900
```

## Planned maintenance

Before taking a leader offline, hand off leadership:

```bash
ssh hub-a.lan 'ctx hub stepdown'
```

`stepdown` triggers a new election among the remaining followers
before the leader goes offline. In-flight clients briefly pause,
then reconnect to the new leader.

## Failure modes at a glance

| Event                       | What happens                                 |
|-----------------------------|----------------------------------------------|
| Leader crashes              | New election; clients reconnect to new leader |
| Follower crashes            | No write impact; catches up on restart        |
| Network partition (majority) | Majority side keeps serving; minority read-only |
| Network partition (split)   | No quorum; all nodes read-only                |
| Disk full on leader         | Writes rejected; read traffic continues      |

For the full list, see
[Hub failure modes](../operations/hub-failure-modes.md).

## See also

- [Multi-machine recipe](hub-multi-machine.md) — single-node
  deployment
- [Hub operations](../operations/hub.md) — backup and
  maintenance
- [Hub security model](../security/hub.md) — TLS, tokens
