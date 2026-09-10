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

### The cluster that could never start — `internal/cli/hub/core/server`

Two defects made the leadership fields unreachable from the CLI
no matter how correctly they were wired.

`RunDaemon` built the re-exec argv from `--port` and
`--data-dir` only. `--peers` was accepted by the flag parser,
never forwarded, and never seen by the process that actually
ran: every `ctx hub start --daemon --peers ...` started a
standalone hub while reporting success, and the HA recipe's
three `--daemon` commands produced three unrelated hubs.

Running the same command in the foreground, where the flag did
arrive, failed anyway. `Run` derived the Raft address as
`fmt.Sprintf(":%d", port+1)` and passed it to
`raft.NewTCPTransport` as both bind and advertise address.
Raft refuses to advertise an unspecified address, so every
cluster start died on:

```
Error: local bind address is not advertisable
```

Cluster mode has never started on any machine. The two defects
hid each other: the daemon path silently dropped the flag that
would have surfaced the crash.

### The three commands that printed and returned — `internal/cli/hub`

`ctx hub peer add`, `ctx hub peer remove` and `ctx hub
stepdown` each called a `write` helper and returned nil.
`Cluster.Stepdown()` had no caller anywhere in the tree, and
raft's `AddVoter` / `RemoveServer` were never called at all.
The HA recipe documented all three as working cluster
operations, so an operator handing off leadership before
maintenance got "Leadership transferred" from a process that
had asked nobody anything.

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

### Give the cluster an address it can advertise

4. `internal/config/flag`, `internal/config/embed/flag`,
   `internal/assets/commands/flags.yaml` — `--raft-bind`, the
   address this node binds its Raft transport to and
   advertises to its peers. `--peers` becomes the list of the
   other nodes' `--raft-bind` addresses, so every node
   bootstraps the same `{ID, Address}` set; each node
   registers under its own Raft address as both.
5. `internal/cli/hub/core/server/setup.go` —
   `validateRaftBind` rejects an empty, host-less or wildcard
   address before Raft does, because Raft's own error names
   neither the flag nor the value. `errHub.RaftBindRequired`
   and `errHub.RaftBindUnroutable` carry the text, keyed in
   `errors.yaml` like every other hub error.
6. `internal/cli/hub/core/server/run.go` — the Raft node
   starts when `--raft-bind` is given, with or without peers:
   a lone node bootstraps a one-server cluster and elects
   itself, which is the cheapest way for an operator to see
   the new Status fields. Asking for `--peers` without
   `--raft-bind` is an error rather than a hub that quietly
   is not in a cluster.
7. `internal/cli/hub/core/server` — `daemonArgs` (extracted
   from `RunDaemon` so it can be tested without forking)
   forwards both cluster flags.

### Make the cluster commands do what they print

8. `internal/hub` — two admin-token-gated RPCs, `Peer` and
   `Stepdown`, gated the way `Register` and `Revoke` are:
   reshaping a cluster is an operator action, not a client
   one. `Cluster.AddPeer` / `RemovePeer` wrap raft's
   `AddVoter` / `RemoveServer` keyed on the node's Raft
   address, which is also its ID.
9. `internal/hub/err_check.go` — `clusterOpErr` maps
   `raft.ErrNotLeader` to `FailedPrecondition` with a message
   naming `ctx hub status` as the way to find the leader. A
   follower's refusal is something the operator can act on;
   an opaque `Internal` is not.
10. `internal/cli/hub/core/{peer,stepdown}` — both dial the
    hub from the saved connection config and call the RPC.
    `internal/cli/hub/core/admin.Token` resolves `--token`
    then `CTX_HUB_ADMIN_TOKEN` for all three admin commands,
    including `revoke`, whose inline copy it replaces.
11. `--join` (`internal/cli/hub`, `internal/hub.ClusterConfig`)
    — a node that bootstraps its own configuration is a second
    cluster of one, not a member of the first, so a node being
    added has to start with no configuration and wait. `--join`
    with `--peers` is an error: a node either bootstraps or
    joins.

### Carry it on the wire

12. `internal/hub/types.go` — `StatusResponse` gains
   `ClusterEnabled`, `IsLeader`, `LeaderAddr` and
   `ClusterPeers`. All four are additive and JSON-omitempty
   where a zero value is meaningless, so an older client
   decoding a newer response is unaffected.
   `ClusterEnabled` is the disambiguator: without it a
   standalone hub is indistinguishable from a clustered node
   that has lost its leader — both report `IsLeader: false`,
   `LeaderAddr: ""`.
13. `internal/hub/handler.go` — `hubStatus` fills those fields
   from `s.cluster` when one is attached, and warns to stderr
   (`cfgWarn.HubClusterPeers`) if the configuration read fails
   rather than failing the RPC: a Status call is a diagnostic,
   and losing the peer count is not a reason to deny the
   operator the rest of it.
14. `internal/config/warn/warn.go` — `HubClusterPeers`, the
   stderr format for that read.

### Render what the hub actually said

15. `internal/config/hub/hub.go` — `RoleLeader` and
   `RoleStandalone` join `RoleFollower`. `RoleActive` is
   deleted: it labelled the listener-count heuristic and has no
   meaning once the role comes from Raft.
16. `internal/write/hub` — `ClusterStatus` takes a
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
17. `internal/cli/hub/core/status/status.go` — the role,
   the leader and the peer count all come from the response.
   Standalone hubs print no leader and no peer count instead
   of inventing both.

### Correct the docs the change falsifies

18. `docs/recipes/hub-cluster.md` — the "expected output" block
    is replaced with what the command prints. It currently
    shows a per-peer table with sync state and uptime that no
    version of this code has ever produced.
19. `docs/cli/hub.md`, `internal/assets/commands/commands.yaml`
    — `ctx hub status` is described as role, leader, entries,
    peers and dropped listeners; "sync state" and "uptime" are
    dropped because neither is reported.
20. `docs/operations/hub.md` — the monitoring section drops
    `ctx hub status --exit-code` (no such flag exists; the
    command exits non-zero only on RPC failure) and the
    per-peer replication-lag claim (the response carries no
    per-peer sequence). Role flaps stay: with a truthful
    `Role:` line they are now actually observable.
21. `internal/hub/doc.go` — the Status paragraph names the
    leadership fields.
22. `docs/cli/hub.md`, `docs/recipes/hub-cluster.md` —
    `--raft-bind` and `--join` in the start reference, the
    cluster recipe's start commands, and its topology diagram,
    which showed one port per node where there are two. The
    membership and maintenance sections document what the
    commands now do: admin-gated, leader-only, addressed by
    Raft address, and — for an addition — preceded by starting
    the new node with `--join`.

## Tests

- `TestCluster_ThreeNodesElectOneLeader` — three real Raft
  nodes on loopback ports, each bootstrapped with the full
  server list: exactly one leads, all three name the same
  leader address, and each counts the other two as peers. This
  is the regression test for the address arithmetic: on `main`
  the nodes cannot start at all.
- `TestValidateRaftBind` — the wildcard, bare-port and
  host-less forms an operator reaches for, each rejected with
  a message naming the flag.
- `TestDaemonArgs_ForwardsClusterFlags` /
  `TestDaemonArgs_OmitsClusterFlags` — the argv the daemon is
  re-executed with. `--peers` used to be missing from it, so a
  daemonized "cluster" node ran standalone while reporting
  success; the second test keeps a standalone daemon's argv
  unchanged.
- `TestHubStatus_StandaloneReportsClusterDisabled` — a Server
  with no `SetCluster` reports `ClusterEnabled == false` and
  zero values for the rest. Pins the disambiguator.
- `TestHubStatus_ClusterReportsLeadershipState` — a real
  single-node Raft node elects itself; Status then reports
  `ClusterEnabled`, `IsLeader`, a non-empty `LeaderAddr` and
  `ClusterPeers == 0` (one server, no peers). Verified by
  mutation: returning early from the `s.cluster != nil` block
  fails it.
- `TestCluster_PeersExcludesSelf` — the count excludes self on
  a single-node configuration.
- `TestCluster_LeaderAddrIsTheAddress` — the node registers
  under an ID that is not an address, so the old
  address-discarding implementation fails here.
- `TestRenderInfo_*` — the CLI mapping: Standalone with no
  leader and no peers, Leader when the hub says so, and
  Follower on a clustered node with no listeners (the case the
  old listener-count heuristic labelled backwards).
- `TestClusterStatus_Standalone` / `TestClusterStatus_Cluster`
  / `TestClusterStatus_LeaderUnknown` — the three rendered
  shapes. `desc.Text` returns `""` for an unknown key, so a
  renamed text key blanks a line silently; asserting on the
  rendered strings catches it. `TestClusterStatus_*Dropped*`
  keep their existing contract under the new struct argument.

- `TestCluster_PeerAddJoinsNode` /
  `TestCluster_PeerRemoveShrinksCluster` — a node started in
  join mode holds no configuration until the leader adds it,
  then follows that leader and counts as a peer; removing it
  shrinks the configuration again. Real Raft nodes, not mocks.
- `TestCluster_StepdownHandsOffLeadership` — after a transfer
  the old leader is not leading and the other node is.
- `TestPeer_RejectsBadAdminToken` /
  `TestStepdown_RejectsBadAdminToken` — the admin gate.
- `TestPeer_NoClusterIsPrecondition` /
  `TestStepdown_NoClusterIsPrecondition` /
  `TestPeer_FollowerIsPrecondition` — the three refusals that
  used to be confirmations: no Raft node at all, and a
  configuration change asked of a follower (the
  `raft.ErrNotLeader` mapping).
- `TestPeer_ValidatesRequest` — unknown action and empty
  address.
- `TestDaemonArgs_ForwardsJoin` — the boolean flag, which
  carries no value and is the easiest one to drop.

## Out of Scope

- **A `Leadership` streaming RPC** and a `ctx hub leader`
  shortcut, both named as non-goals by the issue. Status is
  the minimum useful surface; a subscription and a shortcut
  are sugar until a caller asks for them.
- **Raft term, commit index, log position.** Useful for
  debugging quorum, a different question from "who leads now".
- **Deterministic bootstrap (H-12).** A node with `--peers`
  still bootstraps the full server list, which works because
  every node is given the same list and `ErrCantBootstrap` is
  tolerated on restart. `--join` + `ctx hub peer add` is the
  `AddVoter` half of H-12 and lands here because `peer add`
  is meaningless without it; the single-`--bootstrap`-node
  flow and the persisted bootstrapped flag stay with H-12.
- **Client-side failover to the new leader.** `ctx hub
  stepdown` now moves leadership, and a client whose stream
  was on the old leader still has to be re-run: `Listen` is
  called once with `sinceSequence` hardcoded to `0`, the same
  latent full-refetch `fix-hub-silent-error-suppression.md`
  parked. Reconnect stays manual and the failure-modes doc
  says so.
- **Authenticated Raft transport (H-10/H-11).** `--raft-bind`
  now makes a real cluster possible on a LAN, and the Raft
  transport is still unauthenticated and unencrypted. The
  security docs already say so; nothing here changes it.
