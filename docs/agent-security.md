---
title: "Defense in Depth: Securing AI Agents"
icon: lucide/shield-check
---

![ctx](images/ctx-banner.png)

## The Problem

**An unattended AI agent with unrestricted access to your machine is an
unattended shell with unrestricted access to your machine**.

This is not a theoretical concern. AI coding agents execute shell commands,
write files, make network requests, and modify project configuration. When
running autonomously (*overnight, in a loop, without a human watching*) the
attack surface is the full capability set of the operating system user
account.

The risk is not that the AI is malicious. The risk is that the AI is
**controllable**: it follows instructions from context, and context can be
poisoned.

## Threat Model

### How Agents Get Compromised

AI agents follow instructions from multiple sources: system prompts,
project files, conversation history, and tool outputs. An attacker who can
inject content into any of these sources can redirect the agent's behavior.

| Vector                                   | How it works                                                                                                                                                  |
|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Prompt injection via dependencies**    | A malicious package includes instructions in its README, changelog, or error output. The agent reads these during installation or debugging and follows them. |
| **Prompt injection via fetched content** | The agent fetches a URL (documentation, API response, Stack Overflow answer) containing embedded instructions.                                                |
| **Poisoned project files**               | A contributor adds adversarial instructions to `CLAUDE.md`, `.cursorrules`, or `.context/` files. The agent loads these at session start.                     |
| **Self-modification between iterations** | In an autonomous loop, the agent modifies its own configuration files. The next iteration loads the modified config with no human review.                     |
| **Tool output injection**                | A command's output (error messages, log lines, file contents) contains instructions the agent interprets and follows.                                         |

### What a Compromised Agent Can Do

Depends entirely on what permissions and access the agent has:

| Access level               | Potential impact                                                                |
|----------------------------|---------------------------------------------------------------------------------|
| Unrestricted shell         | Execute any command, install software, modify system files                      |
| Network access             | Exfiltrate source code, credentials, or context files to external servers       |
| Docker socket              | Escape container isolation by spawning privileged sibling containers            |
| SSH keys                   | Pivot to other machines, push to remote repositories, access production systems |
| Write access to own config | Disable its own guardrails for the next iteration                               |

## The Defense Layers

No single layer is sufficient. Each layer catches what the others miss.

```text
Layer 1: Soft instructions     (CONSTITUTION.md, playbook)
Layer 2: Application controls  (permission allowlist, tool restrictions)
Layer 3: OS-level isolation    (user accounts, filesystem, containers)
Layer 4: Network controls      (firewall rules, airgap)
Layer 5: Infrastructure        (VM isolation, resource limits)
```

### Layer 1: Soft Instructions (*Probabilistic*)

Markdown files like `CONSTITUTION.md` and the Agent Playbook tell the
agent what to do and what not to do. These are probabilistic: the agent
*usually* follows them, but there is no enforcement mechanism.

**What it catches**: Most common mistakes. An agent that has been told
"never delete production data" will usually not delete production data.

**What it misses**: Prompt injection. A sufficiently crafted injection
can override soft instructions. Long context windows dilute attention on
rules stated early. Edge cases where instructions are ambiguous.

**Verdict**: Necessary but not sufficient. Good for the common case.
Do not rely on it for security boundaries.

### Layer 2: Application Controls (*Deterministic at Runtime, Mutable Across Iterations*)

AI tool runtimes (Claude Code, Cursor, etc.) provide permission systems:
tool allowlists, command restrictions, confirmation prompts.

For Claude Code, an explicit allowlist in `.claude/settings.local.json`:

```json
{
  "permissions": {
    "allow": [
      "Bash(make:*)",
      "Bash(go:*)",
      "Bash(git:*)",
      "Bash(ctx:*)",
      "Read",
      "Write",
      "Edit"
    ]
  }
}
```

**What it catches**: The agent cannot run commands outside the allowlist.
If `rm`, `curl`, `sudo`, or `docker` are not listed, the agent cannot
invoke them regardless of what any prompt says.

**What it misses**: The agent can modify the allowlist itself. In an
autonomous loop, the agent writes to `.claude/settings.local.json`, and
the next iteration loads the modified config. The application enforces
the rules, but the application reads the rules from files the agent can
write.

**Verdict**: Strong first layer. Must be combined with self-modification
prevention (Layer 3).

### Layer 3: OS-Level Isolation (*Deterministic and Unbypassable*)

The operating system enforces access controls that no application-level
trick can override. An unprivileged user cannot read files owned by root.
A process without `CAP_NET_RAW` cannot open raw sockets. These are kernel
boundaries.

| Control                    | Purpose                                                                                                                                                                                                        |
|----------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Dedicated user account** | No `sudo`, no privileged group membership (`docker`, `wheel`, `adm`). The agent cannot escalate privileges.                                                                                                    |
| **Filesystem permissions** | Project directory writable; everything else read-only or inaccessible. Agent cannot reach other projects, home directories, or system config.                                                                  |
| **Immutable config files** | `CLAUDE.md`, `.claude/settings.local.json`, and `.context/CONSTITUTION.md` owned by a different user or marked immutable (`chattr +i` on Linux). The agent cannot modify its own guardrails. |

**What it catches**: Privilege escalation, self-modification, lateral
movement to other projects or users.

**What it misses**: Actions within the agent's legitimate scope. If the
agent has write access to source code (which it needs to do its job), it
can introduce vulnerabilities in the code itself.

**Verdict**: Essential. This is the layer that makes the other layers
trustworthy.

OS-level isolation does not make the agent safe; it makes the other
layers meaningful.

### Layer 4: Network Controls

An agent that cannot reach the internet cannot exfiltrate data.
It also cannot ingest new instructions mid-loop from external
documents, API responses, or hostile content.

| Scenario                          | Recommended control                                                                                          |
|-----------------------------------|--------------------------------------------------------------------------------------------------------------|
| Agent does not need the internet  | `--network=none` (container) or outbound firewall drop-all                                                   |
| Agent needs to fetch dependencies | Allow specific registries (npmjs.com, proxy.golang.org, pypi.org) via firewall rules. Block everything else. |
| Agent needs API access            | Allow specific API endpoints only. Use an HTTP proxy with allowlisting.                                      |

**What it catches**: Data exfiltration, phone-home payloads, downloading
additional tools, and instruction injection via fetched content.

**What it misses**: Nothing, if the agent genuinely does not need the
network. The tradeoff is that many real workloads need dependency
resolution, so a full airgap requires pre-populated caches.

### Layer 5: Infrastructure Isolation

The strongest boundary is a separate machine — or something that behaves
like one.

The moment you stop arguing about prompts and start arguing about
kernels, you are finally doing security.

**Containers** (Docker, Podman):

```bash
docker run --rm \
  --network=none \
  --cap-drop=ALL \
  --memory=4g \
  --cpus=2 \
  -v /path/to/project:/workspace \
  -w /workspace \
  your-dev-image \
  ./loop.sh
```

!!! danger "Docker Socket is sudo Access"
    Critical: **never mount the Docker socket** (`/var/run/docker.sock`).

    An agent with socket access can spawn sibling containers with full host
    access, effectively escaping the sandbox. 

    Use **rootless Docker** or Podman to eliminate this escalation path.

**Virtual machines**: The strongest isolation. The guest kernel has no
visibility into the host OS. No shared folders, no filesystem passthrough,
no SSH keys to other machines.

**Resource limits**: CPU, memory, and disk quotas prevent a runaway agent
from consuming all resources. Use `ulimit`, cgroup limits, or container
resource constraints.

## Putting It All Together

A defense-in-depth setup for overnight autonomous runs:

| Layer                 | Implementation                                                                    | Stops                                                |
|-----------------------|-----------------------------------------------------------------------------------|------------------------------------------------------|
| Soft instructions     | `CONSTITUTION.md` with "never delete tests", "always run tests before committing" | Common mistakes (probabilistic)                      |
| Application allowlist | `.claude/settings.local.json` with explicit tool permissions                      | Unauthorized commands (deterministic within runtime) |
| Immutable config      | `chattr +i` on `CLAUDE.md`, `.claude/`, `CONSTITUTION.md`                         | Self-modification between iterations                 |
| Unprivileged user     | Dedicated user, no sudo, no docker group                                          | Privilege escalation                                 |
| Container             | `--cap-drop=ALL --network=none`, rootless, no socket mount                        | Host escape, network exfiltration                    |
| Resource limits       | `--memory=4g --cpus=2`, disk quotas                                               | Resource exhaustion                                  |

Each layer is simple. The strength is in the *combination*.

## Common Mistakes

**"I'll just use `--dangerously-skip-permissions`"**: This disables Layer 2
entirely. Without Layers 3-5, you have no protection at all. Only use this
flag inside a properly isolated container or VM.

**"The agent is sandboxed in Docker"**: A Docker container with the Docker
socket mounted, running as root, with `--privileged`, and full network
access is not sandboxed. It is a root shell with extra steps.

**"CONSTITUTION.md says not to do that"**: Markdown is a suggestion. It
works most of the time. It is not a security boundary. Do not use it as
one.

**"I reviewed the CLAUDE.md, it's fine"**: The agent can modify `CLAUDE.md`
during iteration N. Iteration N+1 loads the modified version. Unless the
file is immutable, your review is stale.

**"The agent only has access to this one project"**: Does the project
directory contain `.env` files, SSH keys, API tokens, or credentials? Does
it have a `.git/config` with push access to a remote? Filesystem isolation
means isolating what is *in* the directory too.

## Checklist

Before running an unattended AI agent:

* [ ] Agent runs as a dedicated unprivileged user (no sudo, no docker group)
* [ ] Agent's config files are immutable or owned by a different user
* [ ] Permission allowlist restricts tools to the project's toolchain
* [ ] Container drops all capabilities (`--cap-drop=ALL`)
* [ ] Docker socket is NOT mounted
* [ ] Network is disabled or restricted to specific domains
* [ ] Resource limits are set (memory, CPU, disk)
* [ ] No SSH keys, API tokens, or credentials are accessible to the agent
* [ ] Project directory does not contain `.env` or secrets files
* [ ] Iteration cap is set (`--max-iterations`)

## Further Reading

* [Running an Unattended AI Agent](recipes/autonomous-loops.md): the
  ctx recipe for autonomous loops, including step-by-step permissions
  and isolation setup
* [Security](security.md): `ctx`'s own trust model and vulnerability
  reporting
* [Autonomous Loops](autonomous-loop.md): full documentation of the
  loop pattern, `PROMPT.md` templates, and troubleshooting
