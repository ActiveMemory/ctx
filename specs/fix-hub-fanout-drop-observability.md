# Fix Hub Fanout Drop Observability

The fan-out broadcaster disconnected a listener whose buffer
filled, but the only record of it was `f.dropped` — a counter
incremented in `internal/hub/fanout.go` and read nowhere. The
disconnect itself was also broken in two ways that made it a
pending daemon crash. Upstream issue:
[ActiveMemory/ctx#94](https://github.com/ActiveMemory/ctx/issues/94).

## Problem

### The counter nothing read — `internal/hub/fanout.go`

`broadcast` bumped `f.dropped` on the `default` branch of its
non-blocking send and moved on. No RPC exposed it, no warning
fired, no test asserted it. An operator whose listeners were
being cut loose had no surface that said so, and a regression
that deleted the `delete(f.subs, ch); close(ch)` block would
have passed the whole suite.

### The stream never ended — `internal/hub/handler.go`

`listenEntries` received with `case entries := <-ch` and never
checked the closed state. Once `broadcast` closed the channel,
a closed channel is always receivable: the handler drained the
buffered slices and then spun on `nil` forever at full CPU. The
RPC never returned, so the client never saw the stream end and
silently missed every entry published afterwards.

### The double close — `internal/hub/fanout.go`

Every Listen stream runs `defer s.listeners.unsubscribe(ch)`,
and `unsubscribe` closed unconditionally. When the stream
finally ended (client TCP drop, shutdown), that second close
landed on a channel `broadcast` had already closed:
`panic: close of closed channel`. `grpc.NewServer()` is
constructed with no recovery interceptor
(`internal/hub/server.go`), so grpc-go does not catch it —
every slow-listener disconnect was a pending hub-daemon crash.

## Solution

### Make the disconnect a real, terminating event

1. `internal/hub/fanout.go` — `unsubscribe` is idempotent.
   Membership in `f.subs` is the open/closed record for each
   channel: a channel already gone from the map has already
   been closed and is left alone. Both closers (broadcast's
   disconnect, the stream's deferred unsubscribe) can now run
   in either order.
2. `internal/config/hub/hub.go` — `ErrSlowListener`, the
   wire-visible message, alongside the existing handler error
   constants.
3. `internal/hub/err_check.go` — `errSlowListener`, a
   package-level `status.Error(codes.ResourceExhausted, ...)`
   sentinel, so the handler returns the same value every time,
   tests can match it with `errors.Is`, and the code travels to
   the client.
4. `internal/hub/handler.go` — `listenEntries` receives with
   `entries, live := <-ch` and returns `errSlowListener` on
   `!live`. The stream ends with a reason instead of spinning,
   and its deferred unsubscribe is now safe.

### Make the counter observable

5. `internal/config/warn/warn.go` — `HubFanOutSlowListener`
   format (the `HubReplicate*` family's pattern), carrying the
   cumulative count so aggregators can rate it.
6. `internal/hub/fanout.go` — `broadcast` splits into a locked
   `deliver` that returns the post-increment count per
   disconnect, and an unlocked warn loop. `logWarn.Warn` writes
   to stderr; doing that under `f.mu` would let a stalled
   stderr pipe freeze subscribe, unsubscribe, every publisher,
   and the Status RPC (which takes `f.mu` via `count()`).
7. `internal/hub/types.go` — `dropped` is an `atomic.Uint64`.
   Writes are already serialized under `f.mu`; atomic makes
   the Status RPC's read safe without taking the broadcast
   mutex, and the typed field carries no 32-bit alignment
   caveat.
8. `StatusResponse.DroppedListeners` (additive, wire-
   compatible) is populated from a new `droppedCount()` and
   rendered by `internal/write/hub.ClusterStatus` through
   `DescKeyWriteHubDroppedListeners`. The line is printed only
   when the count is non-zero, so a healthy hub's output is
   byte-for-byte what it was.

### Correct the docs the mechanism falsified

9. `docs/operations/hub-failure-modes.md` — the Slow Listener
   entry states what actually happens: one client's stream ends
   with `ResourceExhausted`, other listeners and publishers are
   unaffected, nothing leaves the hub's log. Client reconnect
   is documented as manual, because it is.
10. `docs/operations/hub.md`, `hub-failure-modes.md`,
    `internal/cli/connection/core/render/doc.go` — three claims
    that the hub or store "deduplicates by entry ID" were
    false; `Store.Append` assigns a new sequence
    unconditionally and `render.appendShared` appends without
    inspecting existing content. Corrected rather than left to
    contradict the new text.

## Tests

`logWarn.SetSink(io.Discard)` (existing seam) keeps the new
warning out of test output.

- `TestFanOut_DisconnectsSlowListener` — the channel closes,
  the subscriber leaves `f.subs`, the counter increments.
  Verified by mutation: deleting the `delete`/`close` block
  fails it while the original three fan-out tests pass.
- `TestFanOut_UnsubscribeAfterDisconnect` — unsubscribe after a
  disconnect, twice. Verified by mutation: removing the
  membership guard panics with `close of closed channel`.
- `TestFanOut_DroppedCountStartsAtZero` — pins the healthy path
  so the counter cannot drift upward.
- `TestFanOut_DroppedCountRaceWithBroadcast` — four goroutines
  read `droppedCount()` (the Status RPC's read path) while
  `broadcast` disconnects listeners. Verified by mutation: a
  plain `f.dropped++` makes `-race` report a data race.
- `TestListenEntries_SlowListenerEndsStream` — drives the real
  handler with a stalled `send`, then asserts it returns
  `errSlowListener` with code `ResourceExhausted`, that the
  already-buffered entries were delivered first, and that its
  deferred unsubscribe does not panic. Verified by mutation:
  restoring `case entries := <-ch` hangs it to the deadline.
- `TestListenEntries_ContextCancelEndsStream` — the healthy
  shutdown still returns nil and still closes a live channel.
- `TestIntegration_SlowListenerReachesClient` — the same
  contract over a real gRPC stream: a client whose handler
  stalls is disconnected server-side and `Client.Listen`
  returns `ResourceExhausted`, which is what makes
  `ctx connection listen` exit non-zero instead of reporting
  success on a stream it no longer receives. Verified by the
  same mutation, which hangs it.
- `TestClusterStatus_DroppedListeners` /
  `TestClusterStatus_NoDroppedListeners` —
  `desc.Text` returns `""` for an unknown key, so a renamed
  text key would blank the new line silently. The first
  asserts the rendered count, the second the omission at zero
  and the survival of the existing stats line.

## Out of Scope

- **Configurable `fanOutBuffer`** (issue item #3). Tuning the
  constant before the counter is observable would be tuning
  blind.
- **Automatic client reconnect with backoff.** No backoff
  implementation exists anywhere in the repo, and
  `internal/cli/connection/core/listen` calls `Listen` once
  with `sinceSequence` hardcoded to `0` — the same shape of
  latent full-refetch that
  `fix-hub-silent-error-suppression.md` parked for the hubsync
  hook. Until reconnect lands, the failure-modes doc says
  reconnect is manual and names the duplicate-append
  consequence, rather than describing a recovery the code does
  not perform.
- **A recovery interceptor on `grpc.NewServer()`.** It would
  have contained the double-close panic, but containing a
  panic is not the same as not having one; the idempotent
  `unsubscribe` removes the cause. A general interceptor is a
  server-wide decision, not a fan-out fix.
