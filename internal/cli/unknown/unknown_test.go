//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unknown

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// stubRelay swaps the package relay seam for the duration of a test
// and returns a restore func.
func stubRelay(
	fn func(string, string, *entity.TemplateRef) error,
) func() {
	prev := relay
	relay = fn
	return func() { relay = prev }
}

// tempStdin writes content to a temp file and rewinds it, yielding an
// *os.File usable as injected stdin. Empty content models a stream
// with no hook JSON (session resolves to IDUnknown).
func tempStdin(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "stdin-*.json")
	if err != nil {
		t.Fatalf("create temp stdin: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })
	if content != "" {
		if _, err := f.WriteString(content); err != nil {
			t.Fatalf("write temp stdin: %v", err)
		}
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("rewind temp stdin: %v", err)
	}
	return f
}

func TestHandleSystemUnknownEmitsBoxSilencesUsageAndErrors(t *testing.T) {
	cmd := &cobra.Command{Use: "system"}
	var out bytes.Buffer
	cmd.SetOut(&out)

	called := false
	defer stubRelay(func(string, string, *entity.TemplateRef) error {
		called = true
		return nil
	})()

	// Empty stdin → no session → relay leg must be skipped.
	err := handle(cmd, []string{"check-anchor-drift"}, tempStdin(t, ""), SystemConfig)
	if err == nil {
		t.Fatal("want non-nil error for an unknown subcommand")
	}
	if !strings.Contains(err.Error(), "check-anchor-drift") {
		t.Errorf("error should name the verb; got %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"check-anchor-drift", "version skew", "VERBATIM",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("relay box missing %q\n---\n%s", want, got)
		}
	}
	if !cmd.SilenceUsage {
		t.Error("want SilenceUsage=true so cobra does not re-dump help")
	}
	if called {
		t.Error("relay must be skipped when stdin carries no session")
	}
}

func TestHandleSystemUnknownFiresRelayWithSession(t *testing.T) {
	cmd := &cobra.Command{Use: "system"}
	cmd.SetOut(&bytes.Buffer{})

	var gotMsg, gotSID string
	var gotRef *entity.TemplateRef
	defer stubRelay(
		func(msg, sid string, ref *entity.TemplateRef) error {
			gotMsg, gotSID, gotRef = msg, sid, ref
			return nil
		},
	)()

	err := handle(
		cmd, []string{"check-anchor-drift"},
		tempStdin(t, `{"session_id":"sess-9"}`), SystemConfig,
	)
	if err == nil {
		t.Fatal("want non-nil error")
	}
	if gotSID != "sess-9" {
		t.Errorf("relay session = %q, want sess-9", gotSID)
	}
	if !strings.Contains(gotMsg, "check-anchor-drift") {
		t.Errorf("relay msg %q should name the verb", gotMsg)
	}
	if gotRef == nil ||
		gotRef.Variant != hook.VariantUnknownSubcommand ||
		gotRef.Hook != hook.System {
		t.Errorf("relay ref = %+v, want hook=%q variant=%q",
			gotRef, hook.System, hook.VariantUnknownSubcommand)
	}
}

func TestHandleRelayFailureDoesNotMaskError(t *testing.T) {
	cmd := &cobra.Command{Use: "system"}
	cmd.SetOut(&bytes.Buffer{})
	defer stubRelay(func(string, string, *entity.TemplateRef) error {
		return errors.New("relay boom")
	})()

	err := handle(cmd, []string{"x"}, tempStdin(t, `{"session_id":"s"}`), SystemConfig)
	if err == nil || !strings.Contains(err.Error(), "x") {
		t.Fatalf("want unknown-subcommand error naming x, got %v", err)
	}
}

func TestHandleBareReturnsNilNoRelay(t *testing.T) {
	cmd := &cobra.Command{Use: "hook"}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	called := false
	defer stubRelay(func(string, string, *entity.TemplateRef) error {
		called = true
		return nil
	})()

	// Bare group invocation (no leftover args) prints help and exits 0;
	// the help-output itself is covered by integration tests against the
	// real command. Here: no error, no relay, no SilenceUsage flip.
	if err := handle(cmd, nil, tempStdin(t, `{"session_id":"s"}`), HookConfig); err != nil {
		t.Fatalf("bare group: want nil error, got %v", err)
	}
	if called {
		t.Error("bare group must not fire the relay")
	}
	if cmd.SilenceUsage {
		t.Error("bare group must not set SilenceUsage")
	}
}

func TestHandleHookUnknownEmitsHookCopy(t *testing.T) {
	cmd := &cobra.Command{Use: "hook"}
	var out bytes.Buffer
	cmd.SetOut(&out)

	defer stubRelay(func(string, string, *entity.TemplateRef) error {
		return nil
	})()

	err := handle(cmd, []string{"notifyy"}, tempStdin(t, ""), HookConfig)
	if err == nil || !strings.Contains(err.Error(), "notifyy") {
		t.Fatalf("want unknown-subcommand error naming notifyy, got %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"notifyy", "Unknown Hook Subcommand", "CLI drift", "VERBATIM",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("hook relay box missing %q\n---\n%s", want, got)
		}
	}
	if !cmd.SilenceUsage {
		t.Error("want SilenceUsage=true")
	}
}

func TestHandleHookUnknownRelayRefUsesHookLabel(t *testing.T) {
	cmd := &cobra.Command{Use: "hook"}
	cmd.SetOut(&bytes.Buffer{})

	var gotRef *entity.TemplateRef
	defer stubRelay(
		func(_, _ string, ref *entity.TemplateRef) error {
			gotRef = ref
			return nil
		},
	)()

	if err := handle(
		cmd, []string{"pausse"},
		tempStdin(t, `{"session_id":"sess-7"}`), HookConfig,
	); err == nil {
		t.Fatal("want non-nil error")
	}
	if gotRef == nil ||
		gotRef.Hook != hook.Hook ||
		gotRef.Variant != hook.VariantUnknownSubcommand {
		t.Errorf("relay ref = %+v, want hook=%q variant=%q",
			gotRef, hook.Hook, hook.VariantUnknownSubcommand)
	}
}
