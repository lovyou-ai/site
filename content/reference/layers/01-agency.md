# Layer 1: Agency

## Derivation

### The Gap

Layer 0 can record events from actors, hash-chain them, verify integrity, track trust, and enforce authority. But every actor looks the same — the graph records what happened, not whether the actor CHOSE for it to happen. A logging system that mechanically records sensor readings is indistinguishable from an AI agent that deliberates and decides. The difference between "something happened" and "someone did something on purpose" doesn't exist.

**Test:** Can you express "Actor A chose to do X rather than Y, based on goal G" as a sequence of Layer 0 operations? You can record that X happened, link it causally, and track trust. But the concepts of "chose," "rather than Y," and "based on goal G" have no Layer 0 representation. These are genuine structural gaps.

### The Transition

**Observer → Participant**

An actor that merely records becomes an actor that acts with purpose. Three new capabilities emerge that Layer 0 cannot express:
1. **Volition** — acting toward what matters, not just reacting to events
2. **Action** — producing effects in the world, not just recording them
3. **Communication** — exchanging signals with other actors

### Base Operations

What can an agent DO that a recorder cannot?

1. **Value** — identify what matters and why
2. **Intend** — declare a desired future state
3. **Choose** — select one action over alternatives
4. **Act** — produce effects in the world
5. **Signal** — send structured information to others

These are the irreducible operations of agency. Each requires representing something that doesn't yet exist (an intent, a choice, a commitment) — which is precisely what Layer 0 lacks.

### Semantic Dimensions

| Dimension | Values | What it distinguishes |
|-----------|--------|-----------------------|
| **Temporal orientation** | Past (review) / Present (act) / Future (plan) | Is this about what happened, what's happening, or what should happen? |
| **Initiative** | Reactive (responding to events) / Proactive (self-initiated) | Did something trigger this, or did the agent decide on its own? |
| **Scope** | Internal (self) / External (others) | Is this about the agent's own state or about other actors? |
| **Binding** | Uncommitted (assessable) / Committed (binding) | Is this an assessment or a commitment? |

**Justification:** These four dimensions capture the essential axes of agency as studied in philosophy of action (Bratman's planning theory), cognitive science (attention models), and organisational theory (delegation/accountability). Each varies independently:
- An intent is future + proactive + internal + committed
- A resource is present + reactive + internal + uncommitted
- A signal is present + proactive + external + committed
- A consequence is past + reactive + external + uncommitted

### Decomposition

Applying dimensions to base operations:

**Group 0 — Volition** (what matters, what to pursue)

| Primitive | Temporal | Initiative | Scope | Binding | What it does |
|-----------|----------|------------|-------|---------|--------------|
| **Value** | Future | Proactive | Internal | Committed | Identifies what matters and why |
| **Intent** | Future | Proactive | Internal | Committed | Declares a desired future state |
| **Choice** | Present | Proactive | Internal | Committed | Selects one action over alternatives |
| **Risk** | Future | Reactive | Internal | Uncommitted | Assesses potential negative consequences of action |

Gap check: "This matters" (Value), "I want X" (Intent), "I pick this over that" (Choice), "what could go wrong?" (Risk). This covers the volition lifecycle.

**Group 1 — Action** (producing effects)

| Primitive | Temporal | Initiative | Scope | Binding | What it does |
|-----------|----------|------------|-------|---------|--------------|
| **Act** | Present | Proactive | External | Committed | Produces an effect in the world |
| **Consequence** | Past | Reactive | External | Uncommitted | The effects that result from action |
| **Capacity** | Present | Reactive | Internal | Uncommitted | What an actor is able to do |
| **Resource** | Present | Reactive | Internal | Uncommitted | What an actor has available to use |

Gap check: "Do this" (Act), "this happened because of that" (Consequence), "I can do X" (Capacity), "I have Y available" (Resource). This covers the action lifecycle.

**Group 2 — Communication** (exchanging signals)

| Primitive | Temporal | Initiative | Scope | Binding | What it does |
|-----------|----------|------------|-------|---------|--------------|
| **Signal** | Present | Proactive | External | Committed | Sends structured information to others |
| **Reception** | Present | Reactive | External | Uncommitted | Receives and processes incoming signals |
| **Acknowledgment** | Present | Reactive | External | Committed | Confirms receipt and understanding |
| **Commitment** | Past→Present | Reactive | Internal | Committed | Binds oneself to a future action |

Gap check: "Here's what I'm saying" (Signal), "I received that" (Reception), "I understand" (Acknowledgment), "I will do this" (Commitment). This covers the communication lifecycle.

### Gap Analysis

**Behaviors tested against the 12 primitives:**

| Behavior | Maps to | Notes |
|----------|---------|-------|
| AI agent identifies what matters | Value | |
| Agent sets a task objective | Intent | |
| Agent picks approach A over approach B | Choice | |
| Agent assesses risk of deployment | Risk | |
| Agent executes a build step | Act | |
| Agent traces what happened from an action | Consequence | |
| Agent reports its available skills | Capacity | |
| Agent checks available compute budget | Resource | |
| Agent sends a status update | Signal | |
| Agent processes incoming event | Reception | |
| Agent confirms receipt of a directive | Acknowledgment | |
| Agent promises to complete a task by deadline | Commitment | |
| Agent changes its mind about an intent | Intent (intent.abandoned) | |
| Human checks what AI agent is doing | Signal + Consequence | Composition |
| Agent reports why it took an action | Choice + Consequence | Composition |

**No gaps found.** All tested agency behaviors map to a single primitive or a composition of two primitives.

### Completeness Argument

1. **Dimensional coverage:** All meaningful combinations of {temporal, initiative, scope, binding} have a primitive or are degenerate (e.g., past + proactive + internal + committed = a completed commitment, which is Commitment tracking completion, not a separate primitive).

2. **Lifecycle coverage:** The volition lifecycle (value → intend → choose → assess risk) maps to Value → Intent → Choice → Risk. The action lifecycle (act → observe consequence → check capacity → allocate resource) maps to Act → Consequence → Capacity → Resource. The communication lifecycle (signal → receive → acknowledge → commit) maps to Signal → Reception → Acknowledgment → Commitment.

3. **Layer boundary:** None of these primitives require concepts from Layer 2 (Exchange). Agency is about individual actors — it doesn't model interaction between actors (that's Exchange's gap). An agent can value, intend, choose, and act without needing reciprocity, negotiation, or agreement.

---

## Primitive Specifications

Full specifications for all 12 primitives are in `docs/primitives.md` (Layer 1 section). Each primitive declares:
- Subscriptions, emitted events, dependencies, state, mechanical/intelligent flag

## Product Graph

Layer 1 maps to the **Work Graph** — task management where AI agents and humans are on the same graph. See `docs/product-layers.md`.

## Reference

- `docs/derivation-method.md` — The method used for this derivation
- `docs/primitives.md` — Full primitive specifications
- `docs/product-layers.md` — Work Graph product layer
- `docs/tests/primitives/01-agent-audit-trail.md` — Integration test scenario exercising Agency patterns
