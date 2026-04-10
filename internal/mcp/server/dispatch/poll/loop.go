//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package poll

import "time"

// poll checks subscribed resources for mtime changes on a
// fixed interval.
func (p *Poller) poll() {
	ticker := time.NewTicker(defaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.pollStop:
			return
		case <-ticker.C:
			p.CheckChanges()
		}
	}
}
