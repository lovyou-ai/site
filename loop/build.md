# Build Report — Fix: Agent persona badge N+1 and empty-initial issues

## Gap
Critic flagged three issues in `cedd52e`:
1. N+1 query: `GetAgentPersonaForConversation` called once per conversation
2. Blank badge circle when `AgentPersona.Display` is empty (`initial("")` → `""`)
3. No test for the persona-lookup path

## Changes

### `graph/store.go`
- Added `GetAgentPersonasForConversations(ctx, convos)` — single SQL query joining `users` and `agent_personas` for all conversation tag IDs at once. Returns `map[string]*AgentPersona` (convo ID → persona). O(1) round-trips regardless of page size.

### `graph/handlers.go`
- Replaced the N+1 loop in `handleConversations` with a single call to `GetAgentPersonasForConversations`.

### `graph/views.templ` + `graph/views_templ.go`
- Fixed `initial()`: returns `"?"` instead of `""` for empty input. Prevents blank badge circle for unnamed or freshly-created agents.

### `graph/store_test.go`
- Added `TestGetAgentPersonasForConversations`: seeds a persona and agent user, constructs two fake ConversationSummary objects (one with agent, one without), verifies the mapping is correct and that empty input doesn't panic. Skips without `DATABASE_URL`.

## Verification
- `templ generate`: 15 updates, no errors
- `go build -buildvcs=false ./...`: clean
- `go test ./...`: all pass (graph tests run in 0.507s)
