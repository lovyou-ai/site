# Layer 2: Exchange

## Derivation

### The Gap

Layer 1 gives actors values, intents, choices, and communication. An actor can decide, act, and commit. But all of this is unilateral — one actor acting on the world. When two actors interact, Layer 1 has signals ("I'm telling you") but no concept of reciprocity ("I'll do this IF you do that"). Offer, acceptance, obligation, and fairness don't exist.

**Test:** Can you express "Alice offers Bob $100 for a code review, Bob accepts, Alice is now obligated to pay" in Layer 1? Alice can set an Intent, act, and commit. But the concepts of "offer" (conditional on acceptance), "acceptance" (binding bilateral agreement), and "obligation" (created by the agreement, not by unilateral commitment) have no Layer 1 representation. Commitment is unilateral; exchange is bilateral.

### The Transition

**Individual → Dyad**

A single actor becomes two actors in relation. The fundamental new capacity: operations that require two participants and create mutual state.

### Base Operations

What can a dyad DO that an individual cannot?

1. **Establish terms** — define shared language and protocols for interaction
2. **Propose** — make a conditional offer
3. **Agree** — create mutual binding state
4. **Exchange** — transfer value between parties

### Semantic Dimensions

| Dimension | Values | What it distinguishes |
|-----------|--------|-----------------------|
| **Symmetry** | Symmetric (both parties equal) / Asymmetric (one gives, one receives) | Are both parties doing the same thing? |
| **Binding force** | Informational (no obligation) / Conditional (if accepted) / Binding (obligation created) | What happens after this operation? |
| **Completeness** | Partial (in progress) / Complete (resolved) | Is the exchange still open? |
| **Valence** | Positive (value added) / Neutral (information) / Negative (cost, debt) | Does this create value, information, or obligation? |

### Decomposition

**Group 0 — Common Ground** (establishing shared basis for exchange)

| Primitive | Symmetry | Binding | Completeness | Valence | What it does |
|-----------|----------|---------|-------------|---------|--------------|
| **Term** | Symmetric | Informational | Partial | Neutral | Defines shared vocabulary and meaning |
| **Protocol** | Symmetric | Informational | Partial | Neutral | Establishes rules of interaction |
| **Offer** | Asymmetric | Conditional | Partial | Positive | Proposes something to another |
| **Acceptance** | Asymmetric | Binding | Complete | Positive | Accepts or rejects a proposal |

Common ground lifecycle: shared language is established (Term) → rules of interaction are set (Protocol) → proposals are made (Offer) → and accepted or rejected (Acceptance). You can't negotiate if you don't share terms.

**Group 1 — Mutual Binding** (creating and tracking obligations)

| Primitive | Symmetry | Binding | Completeness | Valence | What it does |
|-----------|----------|---------|-------------|---------|--------------|
| **Agreement** | Symmetric | Binding | Complete | Positive | Formalised mutual commitment |
| **Obligation** | Asymmetric | Binding | Partial | Negative | Tracks what actors owe each other |
| **Fulfillment** | Asymmetric | Binding | Complete | Positive | Satisfying an obligation |
| **Breach** | Asymmetric | Binding | Partial | Negative | Detecting when obligations are violated |

Mutual binding lifecycle: parties reach agreement (Agreement) → obligations are created (Obligation) → obligations are satisfied (Fulfillment) → or violated (Breach). These primitives handle the formal side of exchange.

**Group 2 — Value Transfer** (the flow of value)

| Primitive | Symmetry | Binding | Completeness | Valence | What it does |
|-----------|----------|---------|-------------|---------|--------------|
| **Exchange** | Symmetric | Binding | Complete | Positive | Transfer of value between parties |
| **Accountability** | Asymmetric | Informational | Partial | Neutral | Tracing responsibility for exchange outcomes |
| **Debt** | Asymmetric | Binding | Partial | Negative | Outstanding obligation not yet fulfilled |
| **Reciprocity** | Symmetric | Informational | Complete | Positive | Balance of give and take over time |

Value transfer lifecycle: value is exchanged (Exchange) → responsibility is traced (Accountability) → outstanding obligations are tracked (Debt) → and balance is assessed over time (Reciprocity). Reciprocity closes the loop by recognising whether exchanges have been fair.

### Gap Analysis

| Behavior | Maps to | Notes |
|----------|---------|-------|
| Alice and Bob agree on what "done" means | Term | |
| They establish a review process | Protocol | |
| Alice offers to sell her car for $5000 | Offer | |
| Bob says "I'll take it" | Acceptance | |
| They shake hands on the deal | Agreement | |
| Alice owes Bob the car; Bob owes Alice $5000 | Obligation (two instances) | |
| Bob pays Alice $5000 | Fulfillment | |
| Bob doesn't pay | Breach | |
| The car changes hands for the money | Exchange | |
| Who's responsible for the outcome? | Accountability | |
| Bob still owes $2000 from last time | Debt | |
| Over time, they've traded fairly | Reciprocity | |
| Bartering goods | Offer + Acceptance + Exchange | Composition |
| Escrow | Obligation + Agreement + L1.Commitment | Cross-layer composition |
| Reputation from exchange | trust.updated (from L0) via Reciprocity | L0 primitives |

**No gaps found.** Escrow and complex financial instruments are compositions, not missing primitives.

### Completeness Argument

1. **Dimensional coverage:** The {symmetry, binding, completeness, valence} space is covered. Degenerate combinations (e.g., symmetric + binding + complete + negative = mutual breach, which is just two Breach events) don't require new primitives.

2. **Lifecycle coverage:** Common ground: define terms → set protocol → offer → accept. Mutual binding: agree → obligate → fulfil or breach. Value transfer: exchange → account → track debt → assess reciprocity. All three lifecycles are complete.

3. **Layer boundary:** None of these primitives require concepts from Layer 3 (Society). Exchange is about dyads — two actors. Group norms, roles, and collective decision-making are Layer 3's gap.

---

## Primitive Specifications

Full specifications in `docs/primitives.md` (Layer 2 section).

## Product Graph

Layer 2 maps to the **Market Graph** — trust-based marketplace eliminating platform tolls. See `docs/product-layers.md`.

## Reference

- `docs/derivation-method.md` — The derivation method
- `docs/primitives.md` — Full specifications
- `docs/product-layers.md` — Market Graph
- `docs/tests/primitives/02-freelancer-reputation.md` — Integration test scenario
