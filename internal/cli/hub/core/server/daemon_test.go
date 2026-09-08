//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"slices"
	"testing"
)

// TestDaemonArgs_ForwardsClusterFlags pins the flags that turn
// Raft on. The daemon is a re-executed process, so a flag left
// out of this argv is a flag the hub never sees: --peers used to
// be dropped, which started every daemonized hub standalone
// while ctx hub start reported success and the cluster recipe
// said otherwise.
func TestDaemonArgs_ForwardsClusterFlags(t *testing.T) {
	args := daemonArgs(
		9900, "/tmp/hub", "10.0.0.5:9901", "h2:9901,h3:9901",
	)

	for flag, want := range map[string]string{
		"--raft-bind": "10.0.0.5:9901",
		"--peers":     "h2:9901,h3:9901",
	} {
		idx := slices.Index(args, flag)
		if idx < 0 {
			t.Fatalf("%s not forwarded: %v", flag, args)
		}
		if got := args[idx+1]; got != want {
			t.Errorf("%s = %q, want %q", flag, got, want)
		}
	}
}

// TestDaemonArgs_OmitsClusterFlags keeps a standalone daemon's
// argv as it was: an empty --raft-bind would put the hub into
// cluster mode with no address to advertise.
func TestDaemonArgs_OmitsClusterFlags(t *testing.T) {
	args := daemonArgs(9900, "/tmp/hub", "", "")

	want := []string{
		"hub", "start", "--port", "9900",
		"--data-dir", "/tmp/hub",
	}
	if !slices.Equal(args, want) {
		t.Errorf("args = %v, want %v", args, want)
	}
}

// TestValidateRaftBind pins the addresses Raft cannot advertise.
// raft.NewTCPTransport rejects an unspecified address with an
// error naming neither the flag nor the value, and the wildcard
// forms are exactly what an operator reaches for when told to
// give a bind address.
func TestValidateRaftBind(t *testing.T) {
	for _, tc := range []struct {
		addr    string
		wantErr bool
	}{
		{"10.0.0.5:9901", false},
		{"hub-a.lan:9901", false},
		{"127.0.0.1:9901", false},
		{"", true},
		{":9901", true},
		{"0.0.0.0:9901", true},
		{"[::]:9901", true},
		{"10.0.0.5", true},
	} {
		bindErr := validateRaftBind(tc.addr)
		if tc.wantErr && bindErr == nil {
			t.Errorf("%q accepted, want rejected", tc.addr)
		}
		if !tc.wantErr && bindErr != nil {
			t.Errorf("%q rejected: %v", tc.addr, bindErr)
		}
	}
}
