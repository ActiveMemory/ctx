//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "testing"

func TestGitPush(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		// Positive: bare
		{"bare", "git push", true},
		{"bare_with_args", "git push origin main", true},
		{"bare_with_force", "git push --force origin main", true},

		// Positive: after statement separators
		{"after_semicolon", "cd foo; git push", true},
		{"after_and", "make && git push", true},
		{"after_and_no_space", "make &&git push", true},
		{"after_or", "foo || git push", true},
		{"after_pipe", "echo x | git push", true},
		{"after_bg", "sleep 1 & git push", true},

		// Positive: subshells and command substitution
		{"subshell", "(git push)", true},
		{"command_sub_dollar", "$(git push)", true},
		{"command_sub_backtick", "`git push`", true},

		// Positive: newline-separated multi-line
		{"newline", "git status\ngit push origin main", true},

		// Positive: env var prefix
		{"env_var", "GIT_DIR=/tmp/foo git push", true},
		{"multi_env", "GIT_DIR=/x GIT_SSH_COMMAND=ssh git push", true},

		// Positive: command wrappers
		{"time_wrapper", "time git push", true},
		{"nice_wrapper", "nice git push", true},
		{"nohup_wrapper", "nohup git push", true},

		// Positive: git top-level flags
		{"dash_c_path", "git -C /path push", true},
		{"dash_c_config", "git -c push.default=simple push", true},
		{"long_git_dir", "git --git-dir=/path push", true},
		{"long_work_tree", "git --work-tree=/path push", true},
		{"long_no_pager", "git --no-pager push", true},
		{"long_bare", "git --bare push", true},
		{"short_paginate", "git -p push", true},
		{"short_no_pager", "git -P push", true},
		{"mixed_flags", "git -C /path --no-pager push origin", true},
		{"flags_other_order", "nice git --no-pager -C /path push", true},

		// Negative: not a push
		{"empty", "", false},
		{"no_git", "echo hello", false},
		{"other_subcommand", "git status", false},
		{"git_pull", "git pull origin main", false},
		{"git_log", "git log --oneline", false},
		{"git_log_with_grep_push", "git log --grep=push", false},

		// Negative: ref-name starting with push (tail anchor rejects)
		{"push_hyphen", "git push-to-remote", false},
		{"push_underscore", "git push_branch", false},
		{"push_slash", "git push/foo", false},
		{"push_dot", "git push.default", false},

		// Negative: not the `git` program
		{"mygit", "mygit push", false},
		{"gitpush_joined", "gitpush", false},
		{"git_push_joined", "gitpush origin", false},

		// Accepted false positives: `push` as a literal arg after another
		// subcommand. Documented trade-off — over-blocking is preferred
		// to under-blocking for a push guard.
		{"false_positive_log_push", "git log push", true},
		{"false_positive_commit_msg_push", "git commit -m push", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := GitPush.MatchString(c.input)
			if got != c.want {
				t.Errorf("GitPush.MatchString(%q) = %v, want %v", c.input, got, c.want)
			}
		})
	}
}
