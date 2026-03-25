# Build: Grounded-in indicator on agent chat messages

## Gap
Agent replies that used document grounding had no visual indicator. Users couldn't tell whether an agent response was grounded in space documents or drawn purely from the model's training.

## Changes

### `graph/mind.go` — store doc count as tag on reply node
In `replyTo`, when `len(docs) > 0`, a `grounded:N` tag is added to the `CreateNodeParams.Tags` for the reply node. When there are no docs (human-only spaces), no tag is added.

### `graph/handlers.go` — `groundedLabel` helper
New package-level function `groundedLabel(tags []string) string`:
- Scans node tags for the `grounded:N` format
- Returns `"grounded in 1 doc"` or `"grounded in N docs"` for N > 0
- Returns `""` when no grounding tag is present

### `graph/views.templ` — render label in chat templates
Both `chatMessage` and `chatMessageCompact` now call `groundedLabel(msg.Tags)` after the message body. When non-empty, a `<div class="text-[9px] text-violet-400/50 mt-1">` label is rendered inside the bubble below the body text. Styled violet/low-opacity to match the agent aesthetic without competing with message content.

## Verification
- `templ generate` — 15 updates, no errors
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (graph: 0.608s)
