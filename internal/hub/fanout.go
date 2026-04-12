//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

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

// unsubscribe removes and closes a listener channel.
//
// Parameters:
//   - ch: channel previously returned by subscribe
func (f *fanOut) unsubscribe(ch chan []Entry) {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.subs, ch)
	close(ch)
}

// broadcast sends entries to all active listeners.
// Non-blocking: slow listeners get disconnected to prevent
// unbounded buffering.
//
// Parameters:
//   - entries: entries to deliver to all subscribers
func (f *fanOut) broadcast(entries []Entry) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for ch := range f.subs {
		select {
		case ch <- entries:
		default:
			// Slow listener — disconnect to prevent loss.
			delete(f.subs, ch)
			close(ch)
			f.dropped++
		}
	}
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
