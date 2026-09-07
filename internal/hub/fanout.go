//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// fanOutBuffer is the channel buffer size for each listener.
const fanOutBuffer = 64

// newFanOut creates a fan-out broadcaster.
//
// Returns:
//   - *fanOut: initialized broadcaster with no subscribers
func newFanOut() *fanOut {
	return &fanOut{
		subs: make(map[chan []Entry]struct{}),
	}
}

// subscribe returns a channel that receives broadcast
// entries. Call unsubscribe when done.
//
// Returns:
//   - chan []Entry: channel delivering broadcast entries
func (f *fanOut) subscribe() chan []Entry {
	f.mu.Lock()
	defer f.mu.Unlock()

	ch := make(chan []Entry, fanOutBuffer)
	f.subs[ch] = struct{}{}
	return ch
}

// unsubscribe removes and closes a listener channel. It is
// idempotent: [fanOut.broadcast] may already have disconnected
// and closed ch, and every Listen stream unsubscribes on the way
// out via defer. Membership in f.subs is the open/closed record,
// so a channel already gone from the map is left alone rather
// than closed a second time — which would panic and, with no
// recovery interceptor on the gRPC server, take the hub daemon
// down.
//
// Parameters:
//   - ch: channel previously returned by subscribe
func (f *fanOut) unsubscribe(ch chan []Entry) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, live := f.subs[ch]; !live {
		return
	}
	delete(f.subs, ch)
	close(ch)
}

// broadcast sends entries to all active listeners.
// Non-blocking: slow listeners get disconnected to prevent
// unbounded buffering. Each disconnect emits a warning so the
// event is visible to operators rather than only bumping a
// counter.
//
// The warnings are written after f.mu is released. Warn writes
// to stderr, and a stalled stderr pipe holding the broadcast
// mutex would freeze subscribe, unsubscribe and the Status RPC
// along with every publisher.
//
// Parameters:
//   - entries: entries to deliver to all subscribers
func (f *fanOut) broadcast(entries []Entry) {
	for _, n := range f.deliver(entries) {
		logWarn.Warn(cfgWarn.HubFanOutSlowListener, n)
	}
}

// deliver is the locked half of [fanOut.broadcast]: it offers
// entries to every subscriber and disconnects the ones that
// cannot take them.
//
// Parameters:
//   - entries: entries to deliver to all subscribers
//
// Returns:
//   - []uint64: cumulative disconnect count after each
//     disconnect this call made, one element per disconnected
//     listener; nil (and unallocated) on the healthy path
func (f *fanOut) deliver(entries []Entry) []uint64 {
	f.mu.Lock()
	defer f.mu.Unlock()

	var counts []uint64
	for ch := range f.subs {
		select {
		case ch <- entries:
		default:
			// Slow listener: disconnect to prevent loss.
			delete(f.subs, ch)
			close(ch)
			counts = append(counts, f.dropped.Add(1))
		}
	}
	return counts
}

// count returns the number of active listeners.
//
// Returns:
//   - uint32: number of active subscriber channels
func (f *fanOut) count() uint32 {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := len(f.subs)
	if n < 0 {
		n = 0
	}
	return uint32(n) //nolint:gosec // len is non-negative
}

// droppedCount returns the cumulative number of listeners
// disconnected for being too slow. The read is atomic rather
// than mutex-guarded so the Status RPC handler never contends
// with an in-flight broadcast.
//
// Returns:
//   - uint64: cumulative slow-listener disconnects
func (f *fanOut) droppedCount() uint64 {
	return f.dropped.Load()
}
