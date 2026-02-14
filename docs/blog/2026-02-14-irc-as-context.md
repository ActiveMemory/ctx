---
title: "Before Context Windows, We Had Bouncers"
date: 2026-02-14
author: Jose Alekhinne
topics:
  - context engineering
  - infrastructure
  - IRC
  - persistence
  - state continuity
---

# Before Context Windows, We Had Bouncers

![ctx](../images/ctx-banner.png)

## The Reset Problem

[IRC][irc] is **stateless**.

[irc]: https://en.wikipedia.org/wiki/Internet_Relay_Chat

* You disconnect, you vanish.
* You reconnect, you begin again.

No buffer.

No memory.

No continuity.

Modern systems are not much different:

* Close the browser tab.
* Lose the Slack scrollback.
* Open a new LLM session.
* Start from zero.

**Resets externalize reconstruction cost onto humans.**

Reconstruction is **tax**.
Tax becomes **entropy**.

---

## Stateless Protocol, Stateful Life

IRC is minimal:

* A TCP connection.
* A nickname.
* A channel.
* A stream of lines.

When the connection drops, you literally ***disappear** from the graph.

The protocol is **stateless**; human systems **are not**.

So you:

* Reconnect.
* Ask what you missed.
* Scroll.
* Reconstruct.

The machine forgets; **you** pay.

---

## The Bouncer Pattern

A `bouncer` is a daemon that remains connected when you do not.

* It holds your seat.
* It buffers what you missed.
* It keeps your identity online.

[ZNC][znc] is one such bouncer.

[znc]: https://en.wikipedia.org/wiki/ZNC

With **ZNC**:

* Your client does not connect to IRC.
* It connects to ZNC.
* ZNC connects upstream.

Client sessions become **ephemeral**.  
Presence becomes **infrastructural**.

!!! tip "ZNC is tmux for IRC"
    * Close your laptop.
      * ZNC remains.

    * Switch devices.
      * ZNC persists.

This is **not** convenience; this is **continuity**.

---

## Presence Without Flapping

With a bouncer:

* Closing your client does not emit `PART`.
* Reopening does not emit `JOIN`.

You do not flap in and out of existence.

From the channel’s perspective, you remain.

From your perspective, history accumulates.

* Buffers persist.
* Identity persists.
* Context persists.

This pattern predates AI.

---

## Before LLM Context Windows

An LLM session without memory is IRC without a bouncer.

* Close the window.
* Start over.
* Re-explain intent.
* Rehydrate context.

That is **friction**.

!!! tip "This Walks and Talks like ctx"
    Context engineering moves memory
    **out of sessions** and **into infrastructure**.

    * `ZNC` does this for IRC.
    * `ctx` does this for agents.

    Same principle:

    * Volatile interface.
    * Persistent substrate.

    Different fabric.

---

## Minimal Architecture

My setup is intentionally boring:

* A $5 small VPS.
* ZNC installed.
* TLS enabled.
* Firewall restricted.

Then:

* **ZNC** connects to `Libera.Chat`.
* `SASL` authentication lives inside **ZNC**.
* Buffers are stored on disk.

My client connects to my VPS, not the network.

The commands do not matter: The **boundaries** do.

* Authentication in **infrastructure**, **not** in the client.
* Memory **server-side**, **not** in scrollback.
* Presence **decoupled** from activity.

Everything else is configuration.

---

## Platform Memory

Yes, it is 2026.

* Discord stores history.
* Slack stores history.
* The dumpster fire on gasoline called X, too, stores history.

**HOWEVER**, they own your **substrate**.

Running a bouncer is **quiet sovereignty**:

* Logs are mine.
* Presence is continuous.
* State does not reset because I closed a tab.

**Small acts compound**.

---

## Signal Density

**Primitive systems select for builders**.

Consistent presence in small rooms compounds reputation.

**Quiet compounding outperforms viral spikes.**

---

## Infrastructure as Cognition

**ZNC** is not interesting because it is retro;
it is interesting because it models a principle:

* Stateless protocols **require** stateful wrappers.
* Volatile interfaces **require** durable memory.
* Human systems **require** continuity.

Distilled:

**Humans require context.**

Before context windows, we had bouncers.  

Before AI memory files, we had buffers.

Continuity is **not** a feature; it is a **design decision**.

---

## Build It

If you want the actual setup (*VPS, ZNC, TLS, SASL, firewall...*) there is
a step-by-step recipe:

**[Persistent IRC Presence with ZNC](https://github.com/ActiveMemory/ctx/blob/main/docs/recipes/persistent-irc.md)**

---

## MOTD

When my client connects to my bouncer, it prints:

```text
//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0
```
