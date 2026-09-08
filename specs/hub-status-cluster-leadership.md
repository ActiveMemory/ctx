# Hub Status: Cluster Leadership

The hub keeps a Raft `Cluster` on `Server` and shuts it down in
`GracefulStop`, but nothing ever reads its leadership state. The
Status RPC does not carry it, so `ctx hub status` cannot answer
"who is the leader?" — and the three lines it does print about
cluster state are fabricated from data that has nothing to do
with Raft. Upstream issue:
[ActiveMemory/ctx#96](https://github.com/ActiveMemory/ctx/issues/96).

## Problem

### The leadership state no RPC carries — `internal/hub/handler.go`

`Cluster.IsLeader()` and `Cluster.LeaderAddr()` exist and work.
`hubStatus` never calls either: it reports store stats and
listener counts and returns. `StatusResponse` has no field to
put leadership in even if it did, so an operator running the
hub in HA mode has no supported way to ask which node leads.

### The three lines that answer from the wrong data — `internal/cli/hub/core/status/status.go`

`ctx hub status` prints a role, a leader and a peer count today.
None of the three is what its label claims:

- **`Role:`** is `Active` when `ConnectedClients > 0` and
  `Follower` otherwise. It reports whether anyone is subscribed,
  not what this node's Raft role is. A leader with no listeners
  reads `Follower`.
- **`Leader:`** is `cfg.HubAddr` — the address the CLI just
  dialed. It says "the hub you asked is the leader" for every
  hub, including a follower and including a standalone node with
  no Raft at all.
- **`Peers:`** is `len(resp.EntriesByProject)`, the number of
  distinct origin projects in the store. A single-node hub
  holding entries from three projects reports `Peers: 3`.

This is worse than silence: the failure-modes runbook tells
operators to run `ctx hub status` on each peer when they see
"no leader" errors, and the answer they get back is
manufactured.

### `LeaderAddr` returns an ID, not an address — `internal/hub/cluster.go`

```go
func (c *Cluster) LeaderAddr() string {
    _, id := c.raftNode.LeaderWithID()
    return string(id)
}
```

`LeaderWithID` returns `(ServerAddress, ServerID)`. The method
named `LeaderAddr`, documented as "leader address", discards the
address and returns the ID. Surfacing that value on the wire as
`leader_addr` would make the mismatch a published contract.

### The bootstrap error nothing checks — `internal/hub/cluster.go`

`NewCluster` calls `r.BootstrapCluster(config)` in both branches
and discards the future. A node whose bootstrap fails returns a
healthy-looking `*Cluster` that never elects anyone — the exact
state the new Status fields are meant to make visible, arriving
with no error and no log line.

## Solution

### Make the cluster answerable

1. `internal/hub/cluster.go` — `LeaderAddr` returns the Raft
   `ServerAddress` of the leader, matching its name and its
   docstring. Empty while no leader is known (election in
   progress, or quorum lost).
2. `internal/hub/cluster.go` — `Peers()` reports the number of
   *other* servers in the committed Raft configuration
   (`GetConfiguration`), returning the future's error rather
   than swallowing it. The count is accumulated in a `uint32`
   instead of converting `len()`, so the wire type needs no
   `gosec` narrowing suppression.
3. `internal/hub/cluster.go` — the two bootstrap branches
   collapse into one server list plus one `BootstrapCluster`
   call whose error is checked. `raft.ErrCantBootstrap` is
   tolerated: it is what a restart against existing Raft state
   returns, and that node is already bootstrapped.

### Carry it on the wire

4. `internal/hub/types.go` — `StatusResponse` gains
   `ClusterEnabled`, `IsLeader`, `LeaderAddr` and
   `ClusterPeers`. All four are additive and JSON-omitempty
   where a zero value is meaningless, so an older client
   decoding a newer response is unaffected.
   `ClusterEnabled` is the disambiguator: without it a
   standalone hub is indistinguishable from a clustered node
   that has lost its leader — both report `IsLeader: false`,
   `LeaderAddr: ""`.
5. `internal/hub/handler.go` — `hubStatus` fills those fields
   from `s.cluster` when one is attached, and warns to stderr
   (`cfgWarn.HubClusterPeers`) if the configuration read fails
   rather than failing the RPC: a Status call is a diagnostic,
   and losing the peer count is not a reason to deny the
   operator the rest of it.
6. `internal/config/warn/warn.go` — `HubClusterPeers`, the
   stderr format for that read.

### Render what the hub actually said

7. `internal/config/hub/hub.go` — `RoleLeader` and
   `RoleStandalone` join `RoleFollower`. `RoleActive` is
   deleted: it labelled the listener-count heuristic and has no
   meaning once the role comes from Raft.
8. `internal/write/hub` — `ClusterStatus` takes a
   `ClusterStatusInfo` struct rather than growing a seventh
   positional parameter, and renders two shapes:

   ```
   Role: Standalone          Role: Leader
   Entries: 1248             Leader: 10.0.0.5:9901
                             Entries: 1248  Peers: 2
   ```

   The leader line reads `Leader: unknown (election in
   progress)` when the cluster is up but `LeaderAddr` is empty.
   `Entries:` gets its own text key for the standalone shape,
   where the combined `Entries: %d  Peers: %d` line would be
   printing a peer count that does not exist. The
   `Dropped listeners:` line keeps its non-zero condition.
9. `internal/cli/hub/core/status/status.go` — the role,
   the leader and the peer count all come from the response.
   Standalone hubs print no leader and no peer count instead
   of inventing both.

### Correct the docs the change falsifies

10. `docs/recipes/hub-cluster.md` — the "expected output" block
    is replaced with what the command prints. It currently
    shows a per-peer table with sync state and uptime that no
    version of this code has ever produced.
11. `docs/cli/hub.md`, `internal/assets/commands/commands.yaml`
    — `ctx hub status` is described as role, leader, entries,
    peers and dropped listeners; "sync state" and "uptime" are
    dropped because neither is reported.
12. `docs/operations/hub.md` — the monitoring section drops
    `ctx hub status --exit-code` (no such flag exists; the
    command exits non-zero only on RPC failure) and the
    per-peer replication-lag claim (the response carries no
    per-peer sequence). Role flaps stay: with a truthful
    `Role:` line they are now actually observable.
13. `internal/hub/doc.go` — the Status paragraph names the
    leadership fields.

## Tests

- `TestHubStatus_StandaloneReportsClusterDisabled` — a Server
  with no `SetCluster` reports `ClusterEnabled == false` and
  zero values for the rest. Pins the disambiguator.
- `TestHubStatus_ClusterReportsLeadershipState` — a real
  single-node Raft node elects itself; Status then reports
  `ClusterEnabled`, `IsLeader`, a non-empty `LeaderAddr` and
  `ClusterPeers == 0` (one server, no peers). Verified by
  mutation: returning early from the `s.cluster != nil` block
  fails it.
- `TestCluster_PeersCountsOtherServers` — the count excludes
  self on a single-node configuration.
- `TestClusterStatus_Standalone` / `TestClusterStatus_Cluster`
  / `TestClusterStatus_LeaderUnknown` — the three rendered
  shapes. `desc.Text` returns `""` for an unknown key, so a
  renamed text key blanks a line silently; asserting on the
  rendered strings catches it. `TestClusterStatus_*Dropped*`
  keep their existing contract under the new struct argument.

## Out of Scope

- **A `Leadership` streaming RPC** and a `ctx hub leader`
  shortcut, both named as non-goals by the issue. Status is
  the minimum useful surface; a subscription and a shortcut
  are sugar until a caller asks for them.
- **Raft term, commit index, log position.** Useful for
  debugging quorum, a different question from "who leads now".
- **The peer addressing the count now exposes.** `ctx hub
  start --peers` registers each peer with `ID == Address ==
  the operator's string` while the local node registers
  `ID = ":<grpc port>"`, `Address = ":<grpc port + 1>"`, so a
  multi-node cluster never agrees on identity or on which port
  carries Raft. That is tracked as **H-12** (deterministic
  bootstrap) and **H-28** (decouple the Raft bind port) in
  TASKS.md and needs a flag surface this change does not have.
  The new `Peers:` line reports the committed configuration
  faithfully, including when that configuration is wrong,
  which is what makes the wrongness visible.
- **`ctx hub stepdown` is still a no-op** that prints
  "Leadership transferred" without calling
  `Cluster.Stepdown()`. Same class of defect as this one but a
  separate RPC to add; filed as a TASKS entry rather than
  bundled in.
