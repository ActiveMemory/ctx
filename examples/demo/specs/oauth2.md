# OAuth2 Authentication Spec

## Overview

Implement OAuth2 authentication supporting Google and GitHub providers.

## Requirements

### Functional

1. **Provider Support**
   - Google OAuth2
   - GitHub OAuth2
   - Extensible provider interface for future additions

2. **Flow**
   - User clicks "Sign in with Google/GitHub"
   - Redirect to provider's authorization page
   - Provider redirects back with authorization code
   - Exchange code for access token
   - Fetch user profile from provider
   - Create or update local user record
   - Issue JWT session token

3. **User Linking**
   - If email already exists, link OAuth identity to existing account
   - If new email, create new user account
   - Store provider ID for future logins

### Non-Functional

- Token exchange must complete in < 2 seconds
- Handle provider downtime gracefully (show user-friendly error)
- Log all OAuth events for security auditing

## API Endpoints

```
GET  /auth/oauth/{provider}          # Initiate OAuth flow
GET  /auth/oauth/{provider}/callback # Handle OAuth callback
POST /auth/logout                    # Revoke session
```

## Data Model

```go
type OAuthIdentity struct {
    ID         string    `json:"id"`
    UserID     string    `json:"user_id"`
    Provider   string    `json:"provider"`   // "google", "github"
    ProviderID string    `json:"provider_id"`
    Email      string    `json:"email"`
    CreatedAt  time.Time `json:"created_at"`
}
```

## Configuration

```yaml
oauth:
  google:
    client_id: ${GOOGLE_CLIENT_ID}
    client_secret: ${GOOGLE_CLIENT_SECRET}
    redirect_url: https://example.com/auth/oauth/google/callback
  github:
    client_id: ${GITHUB_CLIENT_ID}
    client_secret: ${GITHUB_CLIENT_SECRET}
    redirect_url: https://example.com/auth/oauth/github/callback
```

## Security Considerations

- Use `state` parameter to prevent CSRF attacks
- Validate redirect URLs against allowlist
- Never log access tokens or client secrets
- Store only necessary user data from provider

## Testing

- Unit tests for token exchange logic
- Integration tests with mock OAuth provider
- E2E test with real providers in staging environment

## Tasks

These map to `.context/TASKS.md` Phase 2:

1. [ ] Create OAuth provider interface
2. [ ] Implement Google OAuth provider
3. [ ] Implement GitHub OAuth provider
4. [ ] Add OAuth callback handler
5. [ ] Implement user linking logic
6. [ ] Add OAuth configuration loading
7. [ ] Write integration tests
