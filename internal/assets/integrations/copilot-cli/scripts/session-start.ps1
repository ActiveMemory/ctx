# ctx session start hook for Copilot CLI (PowerShell)
# Bootstraps context and loads the agent packet

try { ctx system bootstrap 2>$null } catch {}
try { ctx agent 2>$null } catch {}
