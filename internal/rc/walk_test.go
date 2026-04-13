//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
)

func TestWalkForContextDir_GitAnchor(t *testing.T) {
	// Parent workspace has .context, child project has .git but no .context.
	// Walk should discard parent's .context and anchor to child's git root.
	//
	//   workspace/
	//     .context/          ← parent's context (should be ignored)
	//     child-project/
	//       .git/            ← child's git root
	//       src/             ← CWD
	tmp := t.TempDir()
	workspace := filepath.Join(tmp, "workspace")
	parentCtx := filepath.Join(workspace, dir.Context)
	childProject := filepath.Join(workspace, "child-project")
	childGit := filepath.Join(childProject, ".git")
	childSrc := filepath.Join(childProject, "src")

	for _, d := range []string{parentCtx, childGit, childSrc} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	origDir, _ := os.Getwd()
	_ = os.Chdir(childSrc)
	defer func() { _ = os.Chdir(origDir) }()

	got := walkForContextDir(dir.Context)

	// Should anchor to child-project, not use parent's .context.
	wantResolved, _ := filepath.EvalSymlinks(filepath.Join(childProject, dir.Context))
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("walkForContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestWalkForContextDir_NoGit(t *testing.T) {
	// No .git anywhere, parent has .context.
	// Walk should fall through to cwd/.context.
	//
	//   workspace/
	//     .context/          ← parent's context (no git to confirm)
	//     child/             ← CWD
	tmp := t.TempDir()
	workspace := filepath.Join(tmp, "workspace")
	parentCtx := filepath.Join(workspace, dir.Context)
	child := filepath.Join(workspace, "child")

	for _, d := range []string{parentCtx, child} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	origDir, _ := os.Getwd()
	_ = os.Chdir(child)
	defer func() { _ = os.Chdir(origDir) }()

	got := walkForContextDir(dir.Context)

	wantResolved, _ := filepath.EvalSymlinks(filepath.Join(child, dir.Context))
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("walkForContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestWalkForContextDir_SameGitRoot(t *testing.T) {
	// .context and CWD share the same git root.
	// Walk should return the found .context.
	//
	//   project/
	//     .git/
	//     .context/
	//     src/deep/          ← CWD
	tmp := t.TempDir()
	project := filepath.Join(tmp, "project")
	projectGit := filepath.Join(project, ".git")
	projectCtx := filepath.Join(project, dir.Context)
	deep := filepath.Join(project, "src", "deep")

	for _, d := range []string{projectGit, projectCtx, deep} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	origDir, _ := os.Getwd()
	_ = os.Chdir(deep)
	defer func() { _ = os.Chdir(origDir) }()

	got := walkForContextDir(dir.Context)

	wantResolved, _ := filepath.EvalSymlinks(projectCtx)
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("walkForContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestWalkForContextDir_GitWorktreeFile(t *testing.T) {
	// .git is a file (worktree), not a directory.
	// Should still be detected as git root.
	//
	//   project/
	//     .git              ← file (worktree marker)
	//     .context/
	//     src/              ← CWD
	tmp := t.TempDir()
	project := filepath.Join(tmp, "project")
	projectCtx := filepath.Join(project, dir.Context)
	src := filepath.Join(project, "src")

	for _, d := range []string{projectCtx, src} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	// Create .git as a file (like git worktrees do).
	gitFile := filepath.Join(project, ".git")
	if err := os.WriteFile(gitFile, []byte("gitdir: /some/other/path\n"), 0600); err != nil {
		t.Fatalf("write .git file: %v", err)
	}

	origDir, _ := os.Getwd()
	_ = os.Chdir(src)
	defer func() { _ = os.Chdir(origDir) }()

	got := walkForContextDir(dir.Context)

	wantResolved, _ := filepath.EvalSymlinks(projectCtx)
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("walkForContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestWalkForContextDir_NothingFound_GitRoot(t *testing.T) {
	// No .context anywhere, but .git exists.
	// Walk should anchor to git root.
	//
	//   project/
	//     .git/
	//     src/              ← CWD
	tmp := t.TempDir()
	project := filepath.Join(tmp, "project")
	projectGit := filepath.Join(project, ".git")
	src := filepath.Join(project, "src")

	for _, d := range []string{projectGit, src} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	origDir, _ := os.Getwd()
	_ = os.Chdir(src)
	defer func() { _ = os.Chdir(origDir) }()

	got := walkForContextDir(dir.Context)

	wantResolved, _ := filepath.EvalSymlinks(filepath.Join(project, dir.Context))
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("walkForContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestWalkForContextDir_CWDHasContext(t *testing.T) {
	// .context exists in CWD — should always use it regardless of git.
	//
	//   workspace/
	//     .context/
	//     child/
	//       .context/       ← CWD has its own
	tmp := t.TempDir()
	workspace := filepath.Join(tmp, "workspace")
	parentCtx := filepath.Join(workspace, dir.Context)
	child := filepath.Join(workspace, "child")
	childCtx := filepath.Join(child, dir.Context)

	for _, d := range []string{parentCtx, childCtx} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	origDir, _ := os.Getwd()
	_ = os.Chdir(child)
	defer func() { _ = os.Chdir(origDir) }()

	got := walkForContextDir(dir.Context)

	wantResolved, _ := filepath.EvalSymlinks(childCtx)
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("walkForContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestWalkForContextDir_NestedGitRepos(t *testing.T) {
	// Inner git repo (like a submodule) should use inner git root,
	// rejecting outer project's .context.
	//
	//   outer/
	//     .git/
	//     .context/
	//     vendor/inner/
	//       .git/           ← inner git root
	//       src/            ← CWD
	tmp := t.TempDir()
	outer := filepath.Join(tmp, "outer")
	outerGit := filepath.Join(outer, ".git")
	outerCtx := filepath.Join(outer, dir.Context)
	inner := filepath.Join(outer, "vendor", "inner")
	innerGit := filepath.Join(inner, ".git")
	innerSrc := filepath.Join(inner, "src")

	for _, d := range []string{outerGit, outerCtx, innerGit, innerSrc} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	origDir, _ := os.Getwd()
	_ = os.Chdir(innerSrc)
	defer func() { _ = os.Chdir(origDir) }()

	got := walkForContextDir(dir.Context)

	// Should anchor to inner git root, not use outer's .context.
	wantResolved, _ := filepath.EvalSymlinks(filepath.Join(inner, dir.Context))
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("walkForContextDir() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestFindGitRoot_Found(t *testing.T) {
	tmp := t.TempDir()
	project := filepath.Join(tmp, "project")
	gitDir := filepath.Join(project, ".git")
	deep := filepath.Join(project, "a", "b", "c")

	for _, d := range []string{gitDir, deep} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	got := findGitRoot(deep)
	wantResolved, _ := filepath.EvalSymlinks(project)
	gotResolved, _ := filepath.EvalSymlinks(got)

	if gotResolved != wantResolved {
		t.Errorf("findGitRoot() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestFindGitRoot_NotFound(t *testing.T) {
	tmp := t.TempDir()
	noGit := filepath.Join(tmp, "no-git", "deep")
	if err := os.MkdirAll(noGit, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	got := findGitRoot(noGit)
	if got != "" {
		t.Errorf("findGitRoot() = %q, want empty", got)
	}
}
