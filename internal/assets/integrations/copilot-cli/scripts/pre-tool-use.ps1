# ctx pre-tool-use hook for Copilot CLI (PowerShell)
# Ensures context is loaded and blocks dangerous commands

$Tool = $args[0]

try { ctx system context-load-gate 2>$null } catch {}

if ($Tool -eq "bash" -or $Tool -eq "powershell") {
    try { ctx system block-non-path-ctx 2>$null } catch {}
    try { ctx system qa-reminder 2>$null } catch {}
}
