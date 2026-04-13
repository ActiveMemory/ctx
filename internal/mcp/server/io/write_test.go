//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"bytes"
	"os"
	"testing"
)

func TestWriteJSONSuccess(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteJSON(map[string]int{"a": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "{\"a\":1}\n" {
		t.Errorf("output = %q", got)
	}
}

func TestWriteJSONMarshalError(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteJSON(make(chan int))
	if err == nil {
		t.Fatal("expected marshal error")
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, os.ErrClosed
}

func TestWriteJSONWriteError(t *testing.T) {
	w := NewWriter(errWriter{})
	err := w.WriteJSON("hello")
	if err == nil {
		t.Fatal("expected write error")
	}
}
