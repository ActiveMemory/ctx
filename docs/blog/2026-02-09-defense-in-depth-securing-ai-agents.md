---
title: "Defense in Depth: Securing AI Agents"
date: 2026-02-09
author: Jose Alekhinne
reviewed_and_finalized: true
topics:
  - agent security
  - defense in depth
  - prompt injection
  - autonomous loops
  - container isolation
---

# Defense in Depth: Securing AI Agents

![ctx](../images/ctx-banner.png)

## When Markdown Is Not a Security Boundary

*Jose Alekhinne / 2026-02-09*

!!! question "What Happens When Your AI Agent Runs Overnight and Nobody Is Watching?"
    It follows instructions: **That is the problem**.

    Not because it is malicious. Because it is **controllable**.

    It follows instructions from context, and **context can be poisoned**.

I was writing the [autonomous loops recipe][loops-recipe] for `ctx`:
the guide for running an AI agent in a loop overnight, unattended,
working through tasks while you sleep. The original draft had a tip
at the bottom:

> *Use `CONSTITUTION.md` for guardrails. Tell the agent "never delete
> tests" and it usually won't.*

Then I read that sentence back and realized: **that is wishful
thinking.**

[loops-recipe]: ../recipes/autonomous-loops.md

## The Realization

`CONSTITUTION.md` is a Markdown file. The agent reads it at session
start alongside everything else in `.context/`. It is *one* source
of instructions in a context window that also contains system prompts,
project files, conversation history, tool outputs, and whatever the
agent fetched from the internet.

An attacker who can inject content into *any* of those sources can
redirect the agent's behavior. And "*attacker*" does not always mean
a person with malicious intent. It can be:

| Vector           | Example                                                                     |
|------------------|-----------------------------------------------------------------------------|
| A dependency     | A malicious npm package with instructions in its README or error output     |
| A URL            | Documentation page with embedded adversarial instructions                   |
| A project file   | A contributor who adds instructions to `CLAUDE.md` or `.cursorrules`        |
| The agent itself | In an autonomous loop, the agent modifies its own config between iterations |
| A command output | An error message containing instructions the agent interprets and follows   |

**That last vector is the one that kept me up at night** (*literally!*):

In an autonomous loop, the agent modifies files as part of its job. 

If it modifies its *own* configuration files, the next iteration loads the
modified config. 

* No human reviews it. 
* No diff is shown. 
* The agent that starts iteration `N+1` is running with rules written by
iteration `N`.

**The agent can rewrite its own guardrails.**

## Five Layers, Each With a Hole

That's five nested layers of swiss cheese. Alone, each of them has large
holes. Together, they create a **boundary**.

What followed was a week of peeling back assumptions:

Every defenseI examined had a bypass, and the bypass was always the same shape:
the defense was enforced at a level the agent could reach.

### Layer 1: Soft Instructions (Probabilistic)

`CONSTITUTION.md`, the **Agent Playbook**, **system prompts**: 
These tell the agent what to do. 

The agent *usually* follows them.

"**Usually**" is the keyword here.

**The hole**: Prompt injection:

* A sufficiently crafted payload overrides soft instructions. 
* Long context windows **dilute** attention on rules stated early. 
* Edge cases where instructions are ambiguous get resolved in 
  **unpredictable** ways.

**Verdict**: Necessary. Not sufficient. Good for the common case.
Never trust it as a security boundary.

### Layer 2: Application Controls (*Deterministic at Runtime, Mutable Across Iterations*)

Permission allowlists in `.claude/settings.local.json`:

```json
{
  "permissions": {
    "allow": [
      "Bash(make:*)",
      "Bash(go:*)",
      "Bash(git:*)",
      "Read",
      "Write",
      "Edit"
    ]
  }
}
```

If `rm`, `curl`, `sudo`, or `docker` are not in the allowlist, the
agent cannot invoke them. This is **deterministic**. 

The application enforces it regardless of what any prompt says.

**The hole**: The agent can modify the allowlist itself: 

* It has `Write` permission. 
* The allowlist lives in a file. 
* The agent writes to the file. 
* The next iteration loads the modified allowlist.

The application enforces the rules, but the application reads the
rules from files the agent can write.

**Verdict**: Strong first layer. Must be combined with
self-modification prevention.

### Layer 3: OS-Level Isolation (*Unbypassable*)

This is where the defenses stop having holes in the same shape.

The operating system enforces access controls that no
application-level trick can override. An unprivileged user cannot
read files owned by root. A process without `CAP_NET_RAW` cannot
open raw sockets. These are **kernel** boundaries.

| Control                     | What it stops                                      |
|-----------------------------|----------------------------------------------------|
| Dedicated unprivileged user | Privilege escalation, `sudo`, group-based access   |
| Filesystem permissions      | Lateral movement to other projects, system config  |
| Immutable config files      | Self-modification of guardrails between iterations |

Make the agent's instruction files read-only: `CLAUDE.md`,
`.claude/settings.local.json`, `.context/CONSTITUTION.md`. Own them
as a different user, or mark them immutable with `chattr +i` on Linux.

**The hole**: Actions within the agent's legitimate scope: 

* If the agent has write access to source code (*which it needs*), it can
introduce vulnerabilities in the code itself. 
* You cannot prevent this without removing the agent's ability to do its job.

**Verdict**: Essential. This is the layer that makes Layers 1 and 2
**trustworthy**.

OS-level isolation does not make the agent safe; it makes the other
layers **meaningful**.

### Layer 4: Network Controls

An agent that cannot reach the internet cannot exfiltrate data.

It also cannot ingest new instructions mid-loop from external
documents, error pages, or hostile content.

```bash
# Container with no network
docker run --network=none ...

# Or firewall rules allowing only package registries
iptables -A OUTPUT -d registry.npmjs.org -j ACCEPT
iptables -A OUTPUT -d proxy.golang.org -j ACCEPT
iptables -A OUTPUT -j DROP
```

* If the agent genuinely does not need the network, disable it
  entirely. 
* If it needs to fetch dependencies, allow specific
  registries and block everything else.

**The hole**: **None**, if the agent does not need the network. 

Thetradeoff is that many real workloads need dependency resolution,
so a full airgap requires pre-populated caches.

### Layer 5: Infrastructure Isolation

The strongest boundary is a **separate machine**.

**The moment you stop arguing about prompts and start arguing about
kernels, you are finally doing security**.

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

!!! danger "Never Mount the Docker Socket"
    Do not mount `/var/run/docker.sock`, like, **ever**. 

    An agent with socket access can spawn sibling containers with 
    full host access, effectively escaping the sandbox. 

    This is **not** theoretical: the Docker socket grants
    root-equivalent access to the host.

    **Use rootless Docker or Podman to eliminate this escalation path entirely**.

Virtual machines are even stronger: The guest kernel has no
visibility into the host OS. No shared folders, no filesystem
passthrough, no SSH keys to other machines.

## The Pattern

Each layer is straightforward: The strength is in the **combination**:

| Layer                 | Implementation                  | What it stops                                        |
|-----------------------|---------------------------------|------------------------------------------------------|
| Soft instructions     | `CONSTITUTION.md`               | Common mistakes (probabilistic)                      |
| Application allowlist | `.claude/settings.local.json`   | Unauthorized commands (deterministic within runtime) |
| Immutable config      | `chattr +i` on config files     | Self-modification between iterations                 |
| Unprivileged user     | Dedicated user, no sudo         | Privilege escalation                                 |
| Container             | `--cap-drop=ALL --network=none` | Host escape, data exfiltration                       |
| Resource limits       | `--memory=4g --cpus=2`          | Resource exhaustion                                  |

No layer is redundant. Each one catches what the others miss:

* The soft instructions handle the 99% case: "*don't delete tests.*"
* The allowlist prevents the agent from running commands it should
  not.
* The immutable config prevents the agent from modifying the
  allowlist.
* The unprivileged user prevents the agent from removing
  the immutable flag.
* The container prevents the agent from reaching
  anything outside its workspace.
* The resource limits prevent the agent from consuming all system resources.

**Remove any one layer and there is an attack path through the
remaining ones.**

## Common Mistakes I See

These are real patterns, **not** hypotheticals:

**"I'll just use `--dangerously-skip-permissions`."** This disables
Layer 2 entirely. Without Layers 3 through 5, you have no
protection at all. The flag means what it says. If you **ever** need to,
**think thrice**, you probably don't. But, if you ever need to usee this
**only use it inside a properly isolated VM** (*not even a container: a "VM"*).

**"The agent is sandboxed in Docker."** A Docker container with the
Docker socket mounted, running as root, with `--privileged`, and
full network access is not sandboxed. It is a root shell with extra
steps.

**"I reviewed `CLAUDE.md`, it's fine."** You reviewed it **before**
the loop started. The agent modified it during iteration 3. Iteration
4 loaded the modified version. Unless the file is immutable, your
review is futile.

**"The agent only has access to this one project."** Does the
project directory contain `.env` files? SSH keys? API tokens?
A `.git/config` with push access to a remote? Filesystem isolation
means isolating what is **in** the directory too.

## The Connection to Context Engineering

This is the same lesson I keep rediscovering, wearing different
clothes.

In [The Attention Budget][attention-post], I wrote about how every
token competes for the AI's focus. Security instructions in
`CONSTITUTION.md` are subject to the same budget pressure: if the
context window is full of code, error messages, and tool outputs,
the security rules stated at the top get diluted.

In [Skills That Fight the Platform][fight-post], I wrote about how
custom instructions can conflict with the AI's built-in behavior.
Security rules have the same problem: telling an agent "*never run
curl*" in Markdown while giving it unrestricted shell access creates
a contradiction: The agent resolves contradictions unpredictably.
**The agent will often pick the path of least resistance to attain
its objective function**. And, trust me, agents can get far more
creative than the best red-teamer you know.

In [You Can't Import Expertise][import-post], I wrote about how
generic templates fail because they do not encode project-specific
knowledge. Generic security advice fails the same way: "*Don't
exfiltrate data*" is a **category**; blocking outbound network access is
a **control**.

**The pattern across all of these**: Soft instructions are useful
for the common case. **Hard boundaries are required for security**.

Know which is which.

[attention-post]: 2026-02-03-the-attention-budget.md
[fight-post]: 2026-02-04-skills-that-fight-the-platform.md
[import-post]: 2026-02-05-you-cant-import-expertise.md

## The Checklist

Before running an unattended AI agent:

* [ ] Agent runs as a dedicated unprivileged user (*no sudo, no
  docker group*)
* [ ] Agent's config files are immutable or owned by a different
  user
* [ ] Permission allowlist restricts tools to the project's
  toolchain
* [ ] Container drops all capabilities (`--cap-drop=ALL`)
* [ ] Docker socket is **NOT** mounted
* [ ] Network is disabled or restricted to specific domains
* [ ] Resource limits are set (*memory, CPU, disk*)
* [ ] No SSH keys, API tokens, or credentials are accessible
* [ ] Project directory does not contain `.env` or secrets files
* [ ] Iteration cap is set (`--max-iterations`)

This checklist lives in the [Agent Security][security-doc]
reference alongside the full threat model and detailed guidance
for each layer.

[security-doc]: ../agent-security.md

## What Changed in ctx

The autonomous loops recipe now has a full
[permissions and isolation section][loops-permissions] instead of a
one-line tip about `CONSTITUTION.md`. It covers both the explicit
allowlist approach and the `--dangerously-skip-permissions` flag,
with honest guidance about when each is appropriate.

It also has an OS-level isolation table that is not optional:
unprivileged users, filesystem permissions, containers, VMs,
network controls, resource limits, and self-modification prevention.

The [Agent Security][security-doc] page consolidates the threat
model and defense layers into a standalone reference.

These are **not** theoretical improvements. They are the minimum
responsible guidance for a tool that helps people run AI agents
overnight.

[loops-permissions]: ../recipes/autonomous-loops.md#step-4-configure-permissions

---

!!! quote "If you remember one thing from this post..."
    **Markdown is not a security boundary.**

    `CONSTITUTION.md` is a nudge. An allowlist is a gate.

    An unprivileged user in a network-isolated container is a wall.

    Use all three. Trust only the wall.
---

*This post was written during the session that added permissions,
isolation, and self-modification prevention to the autonomous loops
recipe. The security guidance started as a single tip and grew into
two documents. The meta continues.*
