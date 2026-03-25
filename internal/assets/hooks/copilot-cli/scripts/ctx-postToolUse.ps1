# ctx postToolUse hook for GitHub Copilot CLI
# Reads tool result JSON from stdin and appends to audit log.
$ErrorActionPreference = 'SilentlyContinue'

if (Get-Command ctx -ErrorAction SilentlyContinue) {
    $RawInput = $input | Out-String
    $LogDir = Join-Path '.context' 'state'
    $LogFile = Join-Path $LogDir 'copilot-cli-audit.jsonl'

    if (Test-Path '.context') {
        if (-not (Test-Path $LogDir)) {
            New-Item -ItemType Directory -Path $LogDir -Force | Out-Null
        }
        $Timestamp = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
        $Entry = "{`"timestamp`":`"$Timestamp`",`"event`":`"postToolUse`",`"data`":$RawInput}"
        Add-Content -Path $LogFile -Value $Entry -ErrorAction SilentlyContinue
    }
}
