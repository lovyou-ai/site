# Build Report — invite_codes schema + store methods

## Gap

Task requested: add `invite_codes` table (id, space_id, created_by, code, expires_at, max_uses, use_count, created_at) and `InviteCode` struct + `CreateInviteCode`, `GetInviteCode`, `UseInviteCode` methods.

## Findings

This task is already implemented. Prior iterations built the full invite infrastructure in `graph/store.go`:

### Schema (lines 349–431)
- `invites` table: `token` (primary key, serves as both id and code), `space_id`, `created_by`, `created_at`
- `ALTER TABLE invites ADD COLUMN IF NOT EXISTS expires_at`, `max_uses`, `use_count`
- `invite_uses` table: `(token, user_id)` primary key for idempotent use tracking

### `InviteCode` struct (line 1765)
```go
type InviteCode struct {
    Token     string
    SpaceID   string
    CreatedBy string
    CreatedAt time.Time
    ExpiresAt *time.Time
    MaxUses   int
    UseCount  int
}
```

### Methods (lines 1775–1855)
- `CreateInviteCode(ctx, spaceID, createdBy, expiresAt, maxUses) (string, error)` — generates token, inserts
- `GetInviteCode(ctx, token) (*InviteCode, error)` — validates expiry + exhaustion, returns nil/nil if invalid
- `UseInviteCode(ctx, token, userID) error` — idempotent per (token, userID), increments use_count
- `ListInvites` and `RevokeInvite` also present

### Note on task path mismatch
Task referenced `internal/store/migrations/`, `internal/store/store.go`, `internal/store/pg.go` — these paths don't exist. The actual store is `graph/store.go`. The task was written against a hypothetical architecture. Functionality is equivalent.

## Verification
- `go.exe build -buildvcs=false ./...` — clean (no output)
- `go.exe test ./...` — all pass
