# ctx sessionEnd hook for GitHub Copilot CLI
# Records session end event for recall and context persistence.
$ErrorActionPreference = 'SilentlyContinue'

if (Get-Command ctx -ErrorAction SilentlyContinue) {
    ctx system session-event --type end --caller copilot-cli 2>$null
}
