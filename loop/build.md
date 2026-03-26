# Build Report — Invite Management UI in Space Settings

## Gap
Space settings had a minimal invite link section (show/generate one token) with no ability to see active invite codes, their usage, or revoke them.

## Files Changed

### `graph/store.go`
- Added `ListInvites(ctx, spaceID)` — queries all invite codes for a space, newest first; scans token, space_id, created_by, created_at, expires_at, max_uses, use_count
- Added `RevokeInvite(ctx, token)` — deletes an invite code by token

### `graph/handlers.go`
- Added route `POST /app/{slug}/invites` → `handleCreateInviteHTMX` — creates a new invite and returns an `InviteCodeRow` HTML fragment for HTMX inline insertion
- Added route `DELETE /app/{slug}/invites/{token}` → `handleRevokeInvite` — verifies the token belongs to the space (owner-only), then deletes it
- Updated `handleSpaceSettings` to call `ListInvites` and pass `[]InviteCode` to `SettingsView` (removed legacy `inviteToken` query-param path)
- Fixed two error-path `SettingsView` calls in `handleUpdateSpace` and `handleDeleteSpace` to use `nil, nil` for members/invites

### `graph/views.templ` (+ generated `views_templ.go`)
- Added `InviteCodeRow(inv InviteCode, slug string)` component — displays invite URL (readonly, click-to-select), use_count/max_uses and expiry (when set), and a Revoke button using `hx-delete` + `hx-swap="delete"`
- Updated `SettingsView` signature: replaced `inviteToken string` with `members []SpaceMember, invites []InviteCode`
- Replaced "Invite people" section with "Invitations" section: header + "Generate invite link" button (HTMX POST to `/app/{slug}/invites`, target `#invites-list afterbegin`) + `#invites-list` div rendering existing codes via `InviteCodeRow`

## Verification
- `templ generate` — ✓ 16 updates, no errors
- `go.exe build -buildvcs=false ./...` — ✓ clean
- `go.exe test ./...` — ✓ all pass (graph: 0.604s)
