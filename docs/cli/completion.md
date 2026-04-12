---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Completion
icon: lucide/keyboard
---

![ctx](../images/ctx-banner.png)

## `ctx completion`

Generate shell autocompletion scripts.

```bash
ctx completion <shell>
```

### Subcommands

| Shell        | Command                     |
|--------------|-----------------------------|
| `bash`       | `ctx completion bash`       |
| `zsh`        | `ctx completion zsh`        |
| `fish`       | `ctx completion fish`       |
| `powershell` | `ctx completion powershell` |

**Examples**:

```bash
ctx completion bash > /etc/bash_completion.d/ctx
ctx completion zsh  > "${fpath[1]}/_ctx"
ctx completion fish > ~/.config/fish/completions/ctx.fish
ctx completion powershell | Out-String | Invoke-Expression
```

### Installation

=== "Bash"

    ```bash
    # Add to ~/.bashrc
    source <(ctx completion bash)
    ```

=== "Zsh"

    ```bash
    # Add to ~/.zshrc
    source <(ctx completion zsh)
    ```

=== "Fish"

    ```bash
    ctx completion fish | source
    # Or save to completions directory
    ctx completion fish > ~/.config/fish/completions/ctx.fish
    ```

=== "PowerShell"

    ```powershell
    # Add to your PowerShell profile
    ctx completion powershell | Out-String | Invoke-Expression
    ```
