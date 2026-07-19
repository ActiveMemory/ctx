# git-and-signing

## [2026-07-15-141749] Carving one feature out of a working tree that also holds an unrelated in-progress feature needs hunk-level staging + an isolation-gate

**Context**: The tree mixed a computed-index feature, a stale local copy of a contributor's hub PR (#134), and a gosec-hardening effort — interleaved even within shared files (commands.yaml, errors.yaml, .context/*). Needed to commit ONE of them cleanly. Interactive 'git add -p' isn't available headless.

**Lesson**: Whole-file 'git add' bundles the other feature; naive per-file staging silently drops shared-file edits and forgets coupled files. The reliable method: (1) safety-checkpoint the whole tree to a throwaway branch first; (2) stage pure-mine files + apply only your hunks to mixed files via 'git diff HEAD -- f | filter-hunks | git apply --cached'; (3) ISOLATION-GATE: 'git stash --keep-index -u' then run full build+lint+test on the staged-only tree.

**Application**: Always run the isolation-gate before committing a carve — it caught 5 unstaged files (orphaned DescKeys would've broken the build alone) and a half-done env.go fix that a whole-tree build masked. A green full-tree build does NOT prove the carved subset is self-consistent; only building the stashed-down index does.

---

## [2026-07-06-214523] git filter-branch leaves the originals in refs/original and the reflog

**Context**: Rewrote the branch to strip agent Co-Authored-By trailers and called it clean after checking only the branch tip; the originals were still recoverable and visible in git log --all.

**Lesson**: filter-branch rewrites the ref but keeps originals in refs/original/* and every reflog; clean requires update-ref -d on refs/original, git reflog expire --expire=now --all, and git gc --prune=now, then a re-check across ALL refs and dangling objects. Rewriting an annotated/signed tag strips its GPG signature.

**Application**: After any history rewrite: purge refs/original + reflog + gc, then verify agent trailers are zero across the whole object store, not just main..HEAD.

---

## [2026-06-07-170004] git CLI wrapping quirks (consolidated)

**Consolidated from**: 4 entries (2026-03-24 to 2026-05-22)

- `git rev-parse` exits 0 on an unknown long-flag and echoes the literal arg
  back as its only stdout line (treats it as a candidate revision name). A
  non-zero-exit guard never trips, so `--show-current` shipped verbatim into
  handover frontmatter. Validate the OUTPUT shape (length, no `--` prefix,
  hex-ness for SHAs) when wrapping rev-parse, not just the exit code.
  (`--show-current` is a `git branch` flag, not rev-parse.)
- Group git flag constants by the subcommand whose argv they're valid in (//
  Branch subcommand flags, // Rev-parse flags), not by "loose CLI flags" — the
  group comment is informal type info; mis-grouping enables wrong-subcommand
  bugs. Genuinely-spanning flags (-C, --) go under an explicit Cross-subcommand
  group.
- `git describe --tags --abbrev=0` follows reachability from HEAD, not the
  global tag list (diffed against v0.3.0 instead of v0.6.0 on a diverged release
  branch). For "latest release globally" use `git tag --sort=-v:refname | head
  -1`.
- A trailing regex word boundary \b does NOT exclude hyphenated continuations
  (\bgit commit\b matches `git commit-tree`). For porcelain with hyphenated
  cousins (commit-tree, commit-graph, for-each-ref) append a (?!-) negative
  lookahead.

---

## [2026-04-13-153618] GPG signing from non-TTY contexts requires pinentry-mac (or equivalent)

**Context**: git commit failed from Claude Code's shell with 'gpg: signing
failed: No such file or directory' — the default pinentry-curses cannot open a
TTY in agent-invoked shells. Manual commits from a real terminal worked fine.

**Lesson**: GPG's default curses pinentry requires an interactive TTY. In
non-TTY contexts (Claude Code, CI, scripts, cron), signing fails silently-ish.
The fix is to configure a GUI pinentry that uses the OS keychain: brew install
pinentry-mac; echo 'pinentry-program $(brew --prefix)/bin/pinentry-mac' >>
~/.gnupg/gpg-agent.conf; gpgconf --kill gpg-agent. Once the passphrase is saved
in Keychain, signing works from any context.

**Application**: If agents or CI need to sign commits, configure pinentry-mac
(macOS) or pinentry-gtk/pinentry-qt (Linux) with the OS keychain, not
pinentry-curses. This is a one-time setup per machine.

---

