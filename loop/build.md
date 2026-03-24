# Build Report — Checklist Progress Badge

## Gap
The "Getting Started" checklist on the Board had no progress indicator, and step 3 ("See their response") was tracked via localStorage click events rather than the server-side `first_completion_at` field that already existed.

## Changes

### `graph/views.templ`
- **`gettingStartedChecklist`**: Added `hasCompletion bool` parameter.
  - Added progress badge `"X/4"` next to the "Getting Started" heading. Server-side count from `checklistDoneCount(hasTask, hasAgentTask, hasCompletion)` rendered via `data-server-done` attribute; JS updates to `+1` if invite localStorage key is set.
  - Step 3 now uses `hasCompletion` (server-side) instead of localStorage click tracking — checks off "See their response" when `space.FirstCompletionAt != nil` (i.e., an agent completed a task).
  - Step 4 (Invite) still uses localStorage; JS marks it done and updates the badge to `(serverDone+1)/4` on click or on load if already visited.
  - Removed the old nav-count auto-dismiss logic and the separate `checklist_chat_*` localStorage key.
- **`BoardView`**: Added `hasCompletion bool` parameter; passes it to `gettingStartedChecklist`.

### `graph/handlers.go`
- Added `checklistDoneCount(hasTask, hasAgentTask, hasCompletion bool) int` helper (used by the templ template at render time).
- In `handleBoard`: added `hasCompletion := space.FirstCompletionAt != nil` and passed it to `BoardView`.

## Verification
- `templ generate` ✓
- `go.exe build -buildvcs=false ./...` ✓
- `go.exe test ./...` ✓ (all pass)
