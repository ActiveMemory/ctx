---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Reporting Vulnerabilities
icon: lucide/shield
---

![ctx](../images/ctx-banner.png)

Disclosure process for security issues in `ctx`. For the broader
security model (trust boundaries, audit trail, permission hygiene),
see [Security Design](design.md).

## Reporting Vulnerabilities

At `ctx` we take security very seriously.

If you discover a security vulnerability in `ctx`, please report it
responsibly.

**Do NOT open a public issue for security vulnerabilities.**

### Email

Send details to **security@ctx.ist**.

### GitHub Private Reporting

1. Go to the [Security tab](https://github.com/ActiveMemory/ctx/security);
2. Click "*Report a Vulnerability*";
3. Provide a **detailed** description.

### Encrypted Reports (*Optional*)

If your report contains sensitive details (*proof-of-concept exploits,
credentials, or internal system information*), you can encrypt your
message with our PGP key:

* **In-repo**: [`SECURITY_KEY.asc`](https://github.com/ActiveMemory/ctx/blob/main/SECURITY_KEY.asc)
* **Keybase**: [keybase.io/alekhinejose](https://keybase.io/alekhinejose/pgp_keys.asc)

```bash
# Import the key
gpg --import SECURITY_KEY.asc

# Encrypt your report
gpg --armor --encrypt --recipient security@ctx.ist report.txt
```

Encryption is optional. Unencrypted reports to **security@ctx.ist** or
via GitHub Private Reporting are perfectly fine.

### What to Include

* Description of the vulnerability,
* Steps to reproduce,
* Potential impact,
* Suggested fix (*if any*).

## Attribution

We appreciate responsible disclosure and will acknowledge security
researchers who report valid vulnerabilities (*unless they prefer to
remain anonymous*).

## Response Timeline

!!! note "Open Source, Best-Effort Timelines"
    `ctx` is a volunteer-maintained open source project.

    The timelines below are **guidelines**, not guarantees, and depend
    on contributor availability.

    We will address security reports on a best-effort basis and
    prioritize them by severity.


| Stage              | Timeframe                                |
|--------------------|------------------------------------------|
| Acknowledgment     | Within 48 hours                          |
| Initial assessment | Within 7 days                            |
| Resolution target  | Within 30 days (*depending on severity*) |
