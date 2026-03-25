# ctx preToolUse hook for GitHub Copilot CLI
# Reads tool invocation JSON from stdin and blocks dangerous commands.
$ErrorActionPreference = 'SilentlyContinue'

$RawInput = $input | Out-String
if (-not $RawInput) { exit 0 }

try {
    $Data = $RawInput | ConvertFrom-Json
} catch {
    exit 0
}

$ToolName = if ($Data.tool_name) { $Data.tool_name } elseif ($Data.tool) { $Data.tool } else { '' }

# Block dangerous shell commands matching known patterns.
if ($ToolName -eq 'shell' -or $ToolName -eq 'powershell') {
    $Command = ''
    if ($Data.input -and $Data.input.command) {
        $Command = $Data.input.command
    }

    $DangerousPatterns = @(
        'sudo ',
        'rm -rf /',
        'rm -rf ~',
        'Remove-Item -Recurse -Force C:\',
        'Remove-Item -Recurse -Force $env:USERPROFILE',
        'chmod 777',
        'Format-Volume'
    )
    foreach ($Pattern in $DangerousPatterns) {
        if ($Command -like "*$Pattern*") {
            Write-Error 'ctx: blocked dangerous command'
            exit 1
        }
    }

    $IrreversiblePatterns = @('git push', 'git reset --hard')
    foreach ($Pattern in $IrreversiblePatterns) {
        if ($Command -like "*$Pattern*") {
            Write-Error 'ctx: blocked irreversible git operation — review first'
            exit 1
        }
    }
}
