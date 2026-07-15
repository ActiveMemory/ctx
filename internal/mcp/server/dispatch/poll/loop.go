//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package poll

import "time"

// poll checks subscribed resources for mtime changes on a
// fixed interval.
//
// The stop channel is passed by value (captured under the lock
// at goroutine start) instead of read from p.pollStop each
// iteration: Stop and Unsubscribe nil that field out
// concurrently, so reading it here would be a data race.
// Closing the captured channel still unblocks the select.
//
// Parameters:
//   - stop: channel closed to signal the goroutine to exit
func (p *Poller) poll(stop <-chan struct{}) {
	ticker := time.NewTicker(defaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			p.CheckChanges()
		}
	}
}
