# ctx post-tool-use hook for Copilot CLI (PowerShell)
# Checks for post-commit context and task completion

$Tool = $args[0]

if ($Tool -eq "bash" -or $Tool -eq "powershell") {
    try { ctx system post-commit 2>$null } catch {}
}

if ($Tool -eq "edit" -or $Tool -eq "write") {
    try { ctx system check-task-completion 2>$null } catch {}
}
