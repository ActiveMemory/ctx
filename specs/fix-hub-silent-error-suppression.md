# Fix Hub Silent Error Suppression

The session-start hubsync hook and the cluster replication
loop swallowed errors with no logging surface: operators could
not tell whether a sync succeeded, partially failed, or never
reached the network. Upstream issue:
[ActiveMemory/ctx#100](https://github.com/ActiveMemory/ctx/issues/100).

## Problem

### Hubsync hook — `internal/cli/system/core/hubsync/sync.go`

`Sync` returned `""` silently on config-load failure, dial
failure, sync-RPC failure, and entry-write failure. Worse, the
sync-error check was conflated with the empty-result check:

```go
entries, syncErr := client.Sync(
    context.Background(), cfg.Types, 0,
)
if syncErr != nil || len(entries) == 0 {
    return ""
}
```

A real network error was indistinguishable from "nothing new."
The package doc codified the behavior ("Every error is
silently swallowed so the hook never blocks the session
start"). The never-block constraint is correct; the silence is
the bug.

### Replication loop — `internal/hub/replicate.go`

`replicateOnce` returned silently on dial, stream-open, send,
and close-send failures, and on every receive error — including
real transport failures. (The `conn.Close` defer and the
`store.Append` failure path already warn, and append already
keeps consuming the stream; those two sub-items of #100 landed
upstream before this change.)

## Solution

Wire every silent return through the established
`internal/log/warn` sink with format constants in
`internal/config/warn`, preserving both functions' signatures
and non-blocking contracts. Logging is the only behavior
change, plus one un-conflation:

1. `internal/config/warn/warn.go` — nine new format
   constants: `HubSyncLoadConfig`, `HubSyncDial`,
   `HubSyncPull`, `HubSyncWrite` (hubsync hook; `hubsync:`
   prefix per the `notify:` precedent) and
   `HubReplicateDial`, `HubReplicateStream`,
   `HubReplicateSend`, `HubReplicateCloseSend`,
   `HubReplicateRecv` (extending the existing
   `HubReplicateAppend` family).
2. `internal/cli/system/core/hubsync/sync.go` — warn at all
   four silent sites; split `syncErr` from the
   `len(entries) == 0` check so only the error case warns. A
   genuine empty result stays silent.
3. `internal/cli/system/core/hubsync/doc.go` — the contract
   sentence becomes "Every error is surfaced as a stderr
   warning via the warn sink, but never propagates: the hook
   must not block the session start."
4. `internal/hub/replicate.go` — warn at dial, stream-open,
   send, and close-send failures. The receive site
   distinguishes three cases: `io.EOF` is the normal end of
   every sync stream (returns silently — warning here would
   spam stderr once per `ReplicateInterval`); a done caller
   context is routine shutdown noise (silent); anything else
   is a transport failure and warns. Issue #100's proposed
   code warns on every receive error and would have made
   clean replication cycles noisy; this is the one deliberate
   deviation.

## Tests

`warn.SetSink` (existing test seam) captures output in all of
them.

- `internal/cli/system/core/hubsync/sync_test.go` (new; the
  package had no tests):
  - `TestSync_WarnsOnLoadError` — no connect config present;
    warns `hubsync: load connection config:`.
  - `TestSync_WarnsOnDialError` — `HubAddr` containing a
    control character. Empirically the only eager
    `grpc.NewClient` failure mode: almost every malformed
    target (`://invalid`, `unix://not-abs`) is deferred to
    first use by the lazy resolver, but a control character
    fails URL parsing at construction. Warns `hubsync: dial`.
  - `TestSync_WarnsOnPullError` — well-formed but closed
    address; `grpc.NewClient` is lazy, so the failure
    surfaces at the Sync RPC; warns `hubsync: sync from`.
  - `TestSync_NoWarnOnEmptyResult` — real in-process hub
    with zero entries; no warning, empty return (pins the
    un-conflation).
- `internal/hub/replicate_test.go` (new; `replicateOnce` had
  no direct coverage):
  - `TestReplicateOnce_WarnsOnDialError` — master target
    with a control character (same eager-failure rationale
    as the hubsync dial test).
  - `TestReplicateOnce_WarnsOnTransportError` — closed port;
    asserts a `hub replicate` warning from whichever lazy
    stage surfaces the failure.
  - `TestReplicateOnce_CleanReplicationDoesNotWarn` — real
    master with two entries, writable follower; entries
    replicate, `io.EOF` ends the cycle, zero warnings (pins
    the EOF deviation).
  - `TestReplicateOnce_KeepsConsumingAfterAppendError` —
    read-only follower store directory; both appends fail,
    two `hub replicate append` warnings, loop reaches EOF
    (pins continue-on-append-failure).

## Out of Scope

- Structured (JSON) event-log emission; stderr via
  `warn.Warn` is the established pattern (issue's own
  out-of-scope list).
- Making `Sync` return an error — the non-blocking hook
  contract is hard.
- Wiring `startReplication` into `Server.Start` (#96
  territory) and the hubsync hook's hardcoded
  `sinceSequence=0` full-refetch (a separate latent issue,
  noted during review of #93).
