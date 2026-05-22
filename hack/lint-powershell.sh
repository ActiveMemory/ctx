#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# lint-powershell.sh — run PSScriptAnalyzer against embedded
# PowerShell scripts.
#
# Scope: PowerShell scripts that ship inside the binary as embedded
# assets and run on user machines (the `*.ps1` halves of the
# Copilot CLI hook pairs under
# `internal/assets/integrations/copilot-cli/scripts/`). Same
# stakes as the shell-side: bug here hits every Windows / pwsh
# user.
#
# Requires pwsh (PowerShell Core) with PSScriptAnalyzer installed.
# Install hint:
#   pwsh -NoProfile -Command 'Install-Module -Name PSScriptAnalyzer -Force -Scope CurrentUser'
#
# Severity: fails on `Warning` and above (matches PSScriptAnalyzer's
# canonical band; equivalent to the `warning` threshold used in
# lint-shellcheck.sh).
#
# Exit code:
#   0 — no findings
#   1 — findings or pwsh / module not available

set -euo pipefail

if ! command -v pwsh >/dev/null 2>&1; then
  echo "pwsh (PowerShell Core) not installed. Install via:" >&2
  echo "  macOS:        brew install powershell/tap/powershell" >&2
  echo "  Debian/Ubuntu: see https://learn.microsoft.com/powershell/scripting/install/install-debian" >&2
  exit 1
fi

SEVERITY="${SEVERITY:-Warning}"

# Targets: embedded scripts only.
TARGETS=()
while IFS= read -r -d '' f; do
  TARGETS+=("$f")
done < <(
  find internal/assets/integrations/copilot-cli/scripts \
    -type f -name "*.ps1" -print0 | sort -z
)

if [[ ${#TARGETS[@]} -eq 0 ]]; then
  echo "No embedded PowerShell scripts found." >&2
  exit 0
fi

echo "Running PSScriptAnalyzer (severity=$SEVERITY) on ${#TARGETS[@]} script(s)..."

# Write the analyzer driver to a temp file and invoke via `pwsh
# -File`. `-Command` would run the inline body but does NOT bind
# subsequent positional args (-Severity / -Paths) to the script's
# `param()` block — they are interpreted as additional pwsh
# commands and the param values fall back to empty. `-File`
# semantics pass remaining args into the script's parameter
# binder, which is what we need for $Severity / $Paths.
PS_FILE="$(mktemp).ps1"
trap 'rm -f "$PS_FILE"' EXIT

cat > "$PS_FILE" <<'PSEOF'
param([string]$Severity, [string[]]$Paths)

if (-not (Get-Module -ListAvailable -Name PSScriptAnalyzer)) {
  Write-Error "PSScriptAnalyzer not installed. Install via: Install-Module -Name PSScriptAnalyzer -Force -Scope CurrentUser"
  exit 1
}
Import-Module PSScriptAnalyzer

$severities = @{
  "Information" = @("Information","Warning","Error","ParseError")
  "Warning"     = @("Warning","Error","ParseError")
  "Error"       = @("Error","ParseError")
}
if (-not $severities.ContainsKey($Severity)) {
  Write-Error "Unknown severity: $Severity (allowed: Information, Warning, Error)"
  exit 1
}
$allowed = $severities[$Severity]

$findings = @()
foreach ($p in $Paths) {
  $r = Invoke-ScriptAnalyzer -Path $p -Severity $allowed -ErrorAction Stop
  if ($r) { $findings += $r }
}
if ($findings.Count -gt 0) {
  $findings | Format-Table -AutoSize | Out-String -Width 200 | Write-Output
  Write-Error "PSScriptAnalyzer: $($findings.Count) finding(s) at severity >= $Severity"
  exit 1
}
PSEOF

pwsh -NoProfile -File "$PS_FILE" -Severity "$SEVERITY" -Paths "${TARGETS[@]}"
echo "PSScriptAnalyzer: no findings at severity >= $SEVERITY."
