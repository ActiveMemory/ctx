//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

// fanOutBuffer is the channel buffer size for each listener.
const fanOutBuffer = 64

// newFanOut creates a fan-out broadcaster.
func newFanOut() *fanOut {
	return &fanOut{
		subs: make(map[chan []Entry]struct{}),
	}
}

// subscribe returns a channel that receives broadcast
// entries. Call unsubscribe when done.
func (f *fanOut) subscribe() chan []Entry {
	f.mu.Lock()
	defer f.mu.Unlock()

	ch := make(chan []Entry, fanOutBuffer)
	f.subs[ch] = struct{}{}
	return ch
}

// unsubscribe removes and closes a listener channel.
func (f *fanOut) unsubscribe(ch chan []Entry) {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.subs, ch)
	close(ch)
}

// broadcast sends entries to all active listeners.
// Non-blocking: slow listeners may miss entries.
func (f *fanOut) broadcast(entries []Entry) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for ch := range f.subs {
		select {
		case ch <- entries:
		default:
		}
	}
}

// count returns the number of active listeners.
func (f *fanOut) count() uint32 {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := len(f.subs)
	if n < 0 {
		n = 0
	}
	return uint32(n) //nolint:gosec // len is non-negative
}
