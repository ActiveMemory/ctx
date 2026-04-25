//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package start

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/cli"
)

// TestHubStart_AnnotationSkipInit verifies the hub start subcommand
// carries the AnnotationSkipInit annotation. Hub uses
// ~/.ctx/hub-data/, never reads .context/, and must bypass the
// require-context-dir gate so AWS/EKS hub users hit no broken
// windows on first contact.
//
// Spec: specs/single-source-context-anchor.md.
func TestHubStart_AnnotationSkipInit(t *testing.T) {
	c := Cmd()
	if got, ok := c.Annotations[cli.AnnotationSkipInit]; !ok {
		t.Errorf("hub start: missing AnnotationSkipInit annotation")
	} else if got != cli.AnnotationTrue {
		t.Errorf("hub start: AnnotationSkipInit = %q, want %q", got, cli.AnnotationTrue)
	}
}
