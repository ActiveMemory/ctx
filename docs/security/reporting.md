---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Security Policy
icon: lucide/shield
---

![ctx](../images/ctx-banner.png)

## Reporting Vulnerabilities

At `ctx` we take security very seriously.

If you discover a security vulnerability in `ctx`, please report it responsibly.

**Do NOT open a public issue for security vulnerabilities.**

### Email

Send details to **security@ctx.ist**.

### GitHub Private Reporting

1. Go to the [Security tab](https://github.com/ActiveMemory/ctx/security);
2. Click "*Report a Vulnerability*";
3. Provide a **detailed** description.

### Encrypted Reports (*Optional*)

If your report contains sensitive details (*proof-of-concept exploits,
credentials, or internal system information*), you can encrypt your message
with our PGP key:

* **In-repo**: [`SECURITY_KEY.asc`](https://github.com/ActiveMemory/ctx/blob/main/SECURITY_KEY.asc)
* **Keybase**: [keybase.io/alekhinejose](https://keybase.io/alekhinejose/pgp_keys.asc)

```bash
# Import the key
gpg --import SECURITY_KEY.asc

# Encrypt your report
gpg --armor --encrypt --recipient security@ctx.ist report.txt
```

Encryption is optional. Unencrypted reports to **security@ctx.ist** or via
GitHub Private Reporting are perfectly fine.

### What to Include

* Description of the vulnerability,
* Steps to reproduce,
* Potential impact,
* Suggested fix (*if any*).

## Attribution

We appreciate responsible disclosure and will acknowledge security researchers
who report valid vulnerabilities (*unless they prefer to remain anonymous*).

----

### Response Timeline

!!! note "Open source, Best-Effort Timelines"
    `ctx` is a volunteer-maintained open source project. 

    The timelines below are **guidelines**, not guarantees, and depend on 
    contributor availability.
    
    We will address security reports on a best-effort basis and prioritize
    them by severity.


| Stage              | Timeframe                                |
|--------------------|------------------------------------------|
| Acknowledgment     | Within 48 hours                          |
| Initial assessment | Within 7 days                            |
| Resolution target  | Within 30 days (*depending on severity*) |

----

## Trust Model

`ctx` operates within a single trust boundary: **the local filesystem**.

The person who authors `.context/` files is the same person who runs the
agent that reads them. There is no remote input, no shared state, and no
server component.

This means:

* **`ctx` does not sanitize context files for prompt injection.** This is a
  deliberate design choice, not an oversight. The files are authored by the
  developer who owns the machine: Sanitizing their own instructions back
  to them would be counterproductive.
* **If you place adversarial instructions in your own `.context/` files,
  your agent will follow them.** This is expected behavior. You control the
  context; the agent trusts it.

!!! warning "Shared Repositories"
    In shared repositories, `.context/` files should be reviewed in code
    review (*the same way you would review CI/CD config or Makefiles*). A
    malicious contributor could add harmful instructions to
    `CONSTITUTION.md` or `TASKS.md`.

## Security Design

`ctx` is designed with security in mind:

* **No secrets in context**: The constitution explicitly forbids storing
  secrets, tokens, API keys, or credentials in `.context/` files
* **Local only**: `ctx` runs entirely locally with no external network calls
* **No code execution**: ctx reads and writes Markdown files only; it does
  not execute arbitrary code
* **Git-tracked**: Core context files are meant to be committed, so they should
  never contain sensitive data. Exception: `sessions/` and `journal/` contain
  raw conversation data and should be gitignored

## Permission Hygiene

Claude Code evaluates permissions in deny → ask → allow order. `ctx init`
automatically populates `permissions.deny` with rules that block dangerous
operations before the allow list is ever consulted.

**Default deny rules block:**

* `sudo`, `git push`, `rm -rf /`, `rm -rf ~`, `curl`, `wget`, `chmod 777`
* `Read`/`Edit` of `.env`, credentials, secrets, `.pem`, `.key` files

Even with deny rules in place, the allow list accumulates one-off permissions
over time. Periodically review for:

* **Destructive commands**: `git reset --hard`, `git clean -f`, etc.
* **Config injection vectors**: permissions that allow modifying files
  controlling agent behavior (`CLAUDE.md`, `settings.local.json`)
* **Broad wildcards**: overly permissive patterns that pre-approve
  more than intended

## Temp File Cleanup

Hook state files accumulate in the user-specific temp directory
(`$XDG_RUNTIME_DIR/ctx-<uid>/` or `/tmp/ctx-<uid>/`). A `SessionEnd`
hook automatically removes files older than 15 days. For manual cleanup
or cron scheduling, use:

```bash
.context/tools/cleanup-ctx-tmp.sh        # default: 15-day threshold
.context/tools/cleanup-ctx-tmp.sh 7      # custom: 7-day threshold
```

## Best Practices

1. **Review before committing**: Always review `.context/` files before committing
2. **Use `.gitignore`**: If you must store sensitive notes locally,
   add them to `.gitignore`
3. **Drift detection**: Run `ctx drift` to check for potential issues
4. **Permission audit**: Review `.claude/settings.local.json` after busy sessions
