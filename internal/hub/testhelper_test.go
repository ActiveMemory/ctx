//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"net"
	"testing"
)

// testCtx returns a background context for test RPCs.
func testCtx() context.Context {
	return context.Background()
}

// listenRandom returns a TCP listener on a random port.
func listenRandom(t *testing.T) net.Listener {
	t.Helper()
	lis, lisErr := net.Listen("tcp", "127.0.0.1:0")
	if lisErr != nil {
		t.Fatal(lisErr)
	}
	return lis
}
