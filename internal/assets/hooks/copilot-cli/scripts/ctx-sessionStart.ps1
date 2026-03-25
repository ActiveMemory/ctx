# ctx sessionStart hook for GitHub Copilot CLI
# Records session start and loads context status.
$ErrorActionPreference = 'SilentlyContinue'

if (Get-Command ctx -ErrorAction SilentlyContinue) {
    ctx system session-event --type start --caller copilot-cli 2>$null
}
