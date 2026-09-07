---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hub Failure Modes
icon: lucide/alert-triangle
---

![ctx](../images/ctx-banner.png)

# `ctx` Hub: Failure Modes

What can go wrong, what the system does about it, and what you
should do. Complementary to
[`ctx` Hub Operations](hub.md).

!!! info "Design Posture"
    The hub is **best-effort knowledge sharing**, not a durable
    ledger. Local `.context/` files are the source of truth for
    each project; the hub is a fan-out channel. This framing
    informs every failure-mode decision below.

## Network

### Client Loses Connection Mid-Stream

**What happens:** the stream ends and `ctx connection listen`
exits. There is no automatic reconnect: the command calls Listen
once and returns when the stream ends.

**What you should do:** re-run `ctx connection listen`. Nothing
is lost on the hub — its log is append-only, and the replay
covers every entry newer than the sequence the client asks for.
If disconnects repeat, check firewall state on the hub and
`ctx hub status` output.

!!! warning "Reconnect Is Manual Today"
    Two consequences until automatic reconnect lands:

    - Keeping a listener up across disconnects is a supervisor's
      job (systemd, a shell loop), not the command's.
    - The re-run asks for sequence `0`, not the client's
      last-seen sequence, so entries already written to
      `.context/hub/` are appended a second time.

### Slow Listener Disconnected

**What happens:** each `ctx connection listen` stream gets a
buffered fan-out channel. A client that stops draining it (paused
process, saturated link, a laptop that went to sleep) fills the
buffer. Rather than block every publisher, the hub disconnects
that one listener: it drops the subscription and closes the
channel. The stream then ends with a `ResourceExhausted` error
(`listener disconnected: stream not drained, fan-out buffer
full`), so `ctx connection listen` exits non-zero with that
message instead of hanging on a stream that will never carry
another entry.

Only that one client is affected. Other listeners and every
publisher keep going, and nothing is removed from the hub's log —
the entries the disconnected client missed are still there, and a
fresh `ctx connection listen` picks up from the sequence it asks
for.

Each disconnect writes a warning to the hub's stderr and increments
a cumulative counter reported as `Dropped listeners:` in
`ctx hub status`.

**What you should do:** re-run `ctx connection listen` on the
affected client. As with any lost stream, reconnect is manual
today — see
[Client Loses Connection Mid-Stream](#client-loses-connection-mid-stream)
for the caveats. A count that climbs steadily means listeners
cannot keep up with the publish rate: check the listening
client's health and the link to it before assuming the hub is at
fault.

### Partition: Majority Side Reachable

**What happens:** clients routed to the majority side continue to
publish and listen. The minority nodes step down to followers
that cannot accept writes (Raft quorum lost).

**What you should do:** let it heal. When the partition closes,
followers catch up via sequence-based sync automatically.

### Partition: Split Brain (No Quorum)

**What happens:** no node holds a majority, so no leader is
elected. All nodes become read-only. `ctx connection publish` and
`ctx add --share` fail with a "no leader" error; local writes
still succeed.

**What you should do:** fix the network. If the partition is
permanent (e.g., a data center is gone), bootstrap a new cluster
from the survivors with `ctx hub peer remove` for the dead nodes.

### Hub Unreachable during `ctx add --share`

**What happens:** the local write succeeds; the share step prints
a warning and exits non-zero on the share leg only. `--share` is
best-effort; it never blocks local context updates.

**What you should do:** run `ctx connection publish` later to
backfill. Publish the entry once: the hub's log is append-only
and does not deduplicate by entry ID, so re-sharing the same
entry adds a second copy under a new sequence number.

## Storage

### Disk Full on the Leader

**What happens:** `entries.jsonl` append fails. The hub rejects
writes with an error and stays up for read traffic. Clients
retry; followers keep their in-sync status using whatever the
leader already wrote.

**What you should do:** free disk or grow the volume, then
nothing else; the hub resumes accepting writes on the next
append attempt.

### Corrupt `entries.jsonl`

**What happens:** if the last line is a partial JSON write from a
crash, the hub truncates it on startup and logs a warning. If any
earlier line is malformed, the hub refuses to start.

**What you should do:** inspect with
`jq -c . <data-dir>/entries.jsonl > /dev/null` to find the bad
line. Move the bad region to a `.quarantine` file, then start.
Nothing is ever silently dropped.

### `meta.json` / `entries.jsonl` Sequence Mismatch

**What happens:** the hub refuses to start. This usually means
someone copied one file without the other.

**What you should do:** restore both files from the same backup,
or accept the higher sequence by regenerating `meta.json` from
`entries.jsonl` (manual for now; file a bug).

## Cluster

### Leader Crash, Clean Shutdown

**What happens:** `ctx hub stop` triggers `stepdown` first, so
a new leader is elected before the old one exits. In-flight
writes drain. Clients reconnect to the new leader transparently.

### Leader Crash, Hard Fail (Kill -9, Power Loss)

**What happens:** Raft detects the missing heartbeat and elects
a new leader within a few seconds. Writes the old leader accepted
**but had not yet replicated** can be lost. See the Raft-lite
warning in [the cluster recipe](../recipes/hub-cluster.md).

**What you should do:** if you need stronger durability, run
`ctx connection listen` on a dedicated "collector" project that
persists entries locally as a write-ahead backup.

### Split-Brain After Rejoin

**What happens:** Raft reconciles: the minority side's uncommitted
writes are discarded, and the majority's log is authoritative.

**What you should do:** nothing automatic. If you know the
minority had important writes, grep for them in
`<data-dir>/entries.jsonl.rejected` (written by the reconciliation
pass) and replay them with `ctx connection publish`.

## Auth and Tokens

### Lost Admin Token

**What happens:** you cannot register new projects.

**What you should do:** retrieve it from
`<data-dir>/admin.token`. If that file is also gone, stop the hub
and regenerate. Note that **all existing client tokens keep
working**; only new registrations need the admin token.

### Compromised Admin Token

**What happens:** anyone with the token can register new
projects and publish. They cannot read existing entries without
a client token for a project that subscribes.

**What you should do:** rotate the admin token
(regenerate `<data-dir>/admin.token` and restart), revoke
suspicious client registrations via `clients.json`, and audit
`entries.jsonl` for unexpected origins.

### Compromised Client Token

**What happens:** the attacker can publish as that project and
read anything that project is subscribed to. Because `Origin`
is self-asserted on publish, the attacker can also publish
entries tagged with **any other project's** name, so
attribution in `entries.jsonl` cannot be trusted after a
token compromise.

**What you should do:** remove the client's entry from
`clients.json`, restart the hub, and re-register the legitimate
project with a fresh token. Audit `entries.jsonl` for entries
published after the compromise timestamp and quarantine any
that look suspicious; remember that `Origin` on those entries
proves nothing.

### Compromised Hub Host

**What happens:** `<data-dir>/clients.json` stores client
tokens **verbatim** (not hashed). Anyone with read access to
that file has every client token in hand and can impersonate
any registered project until each one is rotated.

**What you should do:** treat it as a total hub compromise.
Stop the hub, wipe `<data-dir>` (keep a forensic copy first),
regenerate the admin token, and have every client re-register.
See [Security model](../security/hub.md#hub-side-token-storage)
for the mitigations that reduce the blast radius while the
hashing follow-up is pending.

## Clock Skew

Hub entries carry a timestamp assigned **by the publishing
client**. The hub does not rewrite timestamps. Clients with
significant clock skew will publish entries that look out of
order in the shared feed.

**What you should do:** run NTP on all client machines. If you
see entries dated in the future or far past, the publisher's
clock is the culprit.

## The Short List

| Symptom                           | First thing to check              |
|-----------------------------------|-----------------------------------|
| Client can't reach hub            | Firewall, then `ctx hub status`   |
| "No leader" errors                | Cluster quorum; run `ctx hub status` on each peer |
| Hub won't start after crash       | Last line of `entries.jsonl`      |
| Entries missing after restore     | Check `clients.json` sequence vs local `.sync-state.json` |
| Duplicate entries in shared feed  | A client re-published; the hub never dedups by ID |
| Followers lagging                 | Disk or network on the follower, not the leader |

## See Also

- [`ctx` Hub Operations](hub.md)
- [`ctx` Hub security model](../security/hub.md)
- [HA cluster recipe](../recipes/hub-cluster.md)
