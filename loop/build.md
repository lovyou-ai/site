# Build Report — Fix: goal progress dashboard security and reliability

## What changed

**File:** `graph/handlers.go` — `handleGoalDetail` function

### Fix 1: Cross-space authorization bypass (security)

Added space membership check after fetching the goal node:

```go
if goal.SpaceID != space.ID {
    http.NotFound(w, r)
    return
}
```

Without this, a request to `/app/space-a/goals/{node-id-from-space-b}` would succeed and expose goal data across space boundaries. The `ListNodes` calls below already used `SpaceID: space.ID` correctly — only the `GetNode` fetch was unguarded.

### Fix 2: Silent error swallowing

Changed both `ListNodes` calls from ignoring errors to returning HTTP 500:

```go
projects, err := h.store.ListNodes(...)
if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }

directTasks, err := h.store.ListNodes(...)
if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
```

Previously a store failure would produce a valid-looking page with 0 tasks and 0% progress — a false success.

### Fix 3: Unbounded queries (invariant 13)

Added `Limit: 200` to both `ListNodes` calls. The store default is 500; making the limit explicit and reasonable prevents a goal with thousands of tasks from returning all of them.

### Non-fix: Assignee rendering

The critic flagged `{ task.Assignee }` as rendering a raw ID. In this codebase, `Assignee` stores the display name; `AssigneeID` is the ID. The render is correct and consistent with the rest of the codebase (board list view, task cards, node detail). No change needed.

## Verification
- `go.exe build -buildvcs=false ./...`: clean
- `go.exe test ./...`: all pass (graph: 0.535s)
