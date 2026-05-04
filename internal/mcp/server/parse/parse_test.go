//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import (
	"os"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

func TestRequestValid(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	req, errResp := Request(data)
	switch {
	case errResp != nil:
		t.Fatal("unexpected error response")
	case req == nil:
		t.Fatal("expected non-nil request")
	case req.Method != "ping":
		t.Errorf("method = %q, want ping", req.Method)
	}
}

func TestRequestMalformed(t *testing.T) {
	req, errResp := Request([]byte(`not-json`))
	if req != nil {
		t.Fatal("expected nil request")
	}
	if errResp == nil || errResp.Error == nil {
		t.Fatal("expected error response")
	}
	if errResp.Error.Code != -32700 {
		t.Errorf("code = %d, want -32700", errResp.Error.Code)
	}
}

func TestRequestNotification(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","method":"notify"}`)
	req, errResp := Request(data)
	if req != nil {
		t.Error("expected nil request for notification")
	}
	if errResp != nil {
		t.Error("expected nil error for notification")
	}
}
