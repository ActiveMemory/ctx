# ctx session end hook for Copilot CLI (PowerShell)
# Checks for unsaved context and records heartbeat

try { ctx system check-context-size 2>$null } catch {}
try { ctx system check-persistence 2>$null } catch {}
try { ctx system check-journal 2>$null } catch {}
try { ctx system check-version 2>$null } catch {}
try { ctx system heartbeat 2>$null } catch {}
