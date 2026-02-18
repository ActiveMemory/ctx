# Persistent IRC Presence with ZNC

## The Problem

IRC is stateless. You disconnect, you vanish. You reconnect, you reconstruct.

This is not unique to IRC. Close the browser tab, lose the Slack scrollback,
open a new LLM session, start from zero. Resets externalize reconstruction
cost onto humans.

This runbook builds a persistent IRC presence using [ZNC](https://znc.in) 
as an always-on bouncer:

```
Client → ZNC → IRC
```

Client sessions become disposable. Presence becomes infrastructural.

This is not a nostalgia project. It is a pattern: stateless protocol,
stateful wrapper, persistent presence. The same pattern that ctx applies
to AI sessions.

## Architecture

* Small VPS (IPv4 + IPv6 optional)
* ZNC installed
* TLS listener
* Firewall restricted to your IP
* SASL authentication upstream
* Buffers stored on disk

No screenshots. No click-driven instructions. Only primitives.

## The Workflow

### Step 1: Provision a VPS

Any minimal Linux host works. You need:

* A static public IP
* Inbound TCP port (e.g., 6699)
* SSH access
* Firewall control

Lock SSH before anything else.

### Step 2: Install ZNC

On Fedora:

```bash
sudo dnf install znc
```

On Debian/Ubuntu:

```bash
sudo apt install znc
```

Do not start the systemd service yet.

### Step 3: Generate Configuration

Run configuration as the ZNC service user, not your shell user.

This avoids the classic "*config in wrong directory*" failure.

Check the service file:

```bash
systemctl cat znc
```

Note the `User=` and `--datadir=` values. Usually:

```text
- User: `_znc`
- Datadir: `/var/lib/znc`
```

Switch to that user:

```bash
sudo -u _znc znc --makeconf --datadir=/var/lib/znc
```

> [!WARNING]
> If you previously ran `znc --makeconf` as yourself, you likely
> created `~/.znc/`. That will **not** be used by systemd.
>
> Delete or ignore it.

### Step 4: TLS Certificate

ZNC needs a PEM file. If you see `Unable to locate pem file`, generate one:

```bash
sudo -u _znc znc --datadir=/var/lib/znc --makepem
```

Ensure `znc.conf` points to the correct PEM path inside `/var/lib/znc`.
Do not reference your home directory.

### Step 5: Start ZNC

```bash
sudo systemctl enable znc
sudo systemctl start znc
```

Verify:

```bash
sudo systemctl status znc
ss -lntp | grep 6699
```

### Step 6: Lock Firewall

Do not leave the port open to the world. Restrict inbound TCP `6699` 
**to your IP only**.

Test connectivity from your machine:

```bash
openssl s_client -connect <vps_ip>:6699
```

If this hangs, it is the firewall. Not TLS. Not ZNC.

### Step 7: Connect Your IRC Client to ZNC

Configure your client:

| Setting  | Value                         |
|----------|-------------------------------|
| Server   | `<vps_ip>`                    |
| Port     | `6699`                        |
| SSL      | Enabled                       |
| Username | `znc_user/network_name`       |
| Password | Your ZNC password             |

Example username: `jose/libera`

> [!WARNING]
> Do not use "personal password" for ZNC.
>
> That field is for `NickServ`, not ZNC.
>
> If you see `You need to send your password`,
> you put it in the wrong field.

### Step 8: Enable SASL

Configure SASL **inside ZNC**, not in your client.

In the ZNC web interface: 

```text
User → Networks → libera → Modules → enable "sasl"
```

Configure:

* **Username**: Your IRC nickname
* **Password**: Your NickServ password
* **Mechanism**: `PLAIN`

Reconnect the network through ZNC's `*status` virtual user:

```
/msg *status Disconnect
/msg *status Connect
```

`*status` is ZNC's built-in control interface. You message it like a
user, and it executes commands on ZNC itself. `/msg *status Help` lists
all available commands.

Then verify SASL worked:

```
/msg NickServ STATUS <nick>
```

It should return `"identified"` without manual `IDENTIFY`.

**Authentication belongs in infrastructure**.

### Step 9: Enable Useful Modules

**Network modules** (*per-network, e.g., `Libera.Chat`*):

* `autoattach`: Reattach to channels on client reconnect
* `buffextras`: Log joins, parts, and mode changes in buffer
* `keepnick`: Reclaim your nick automatically
* `log`: Write channel and query logs to disk
* `route_replies`: Route server replies to the requesting client
* `sasl`: Authenticate upstream (*configured in Step 8*)
* `savebuff`: Persist buffers across ZNC restarts
* `simple_away`: Set away status when no clients are connected

**Global modules** (*apply to all networks*):

* `corecaps`: Negotiate IRCv3 capabilities with the server (*highly recommended*)
* `webadmin`: Browser-based admin panel for managing ZNC

**Avoid**:

* `nickserv`: Use `sasl` instead
* `autoop`: Use `ChanServ` instead

### Step 10: Auto-Op Correctly

If you are the founder of a channel:

```
/msg ChanServ FLAGS #channel <nick> +AO
```

If `ChanServ` says "You are not logged in," SASL is not configured correctly.

### Step 11: Optional MOTD

ZNC can print a banner when your client connects:

```
// Context: https://ctx.ist
// do you remember?
```

This is not decoration; it is a **reminder**: **memory is infrastructure**.

## Common Failure Modes

| Problem                         | Cause                              | Fix                                        |
|---------------------------------|------------------------------------|--------------------------------------------|
| systemd service fails           | Config in wrong directory          | Generate under `/var/lib/znc` as `_znc`    |
| TLS timeout                     | Firewall                           | Test with `openssl`, not assumptions       |
| ChanServ says not logged in     | SASL not configured in ZNC         | Move auth upstream                         |
| Duplicate nick on reconnect     | Client handling NickServ, not ZNC  | Move auth upstream                         |

## Why Do We Use Bouncers?

Without a bouncer: disconnect → reset → reconstruct.

With a bouncer: **continuity**.

ZNC is **not** retro: It is a memory layer. 

Before LLM context windows, we had buffers. 

**Stateless protocols require stateful wrappers**: The pattern persists.

## See Also

* [Before Context Windows, We Had Bouncers](https://ctx.ist/blog/2026-02-14-irc-as-context/):
  The blog post exploring this pattern
* [Running an Unattended AI Agent](https://ctx.ist/recipes/autonomous-loops/): Another recipe
  about persistent infrastructure for ephemeral sessions
