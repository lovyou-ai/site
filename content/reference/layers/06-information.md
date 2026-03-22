# Layer 6: Information

## Derivation

### The Gap

Layer 5 builds artefacts and processes, but everything is concrete — specific events, specific actions, specific tools. There's no concept of symbolic representation, language, or meaning independent of the physical events that carry it. A system can build a database but cannot reason about what "database" means as a category, cannot generalise from examples to principles, and cannot track whether its knowledge is accurate.

**Test:** Can you express "The pattern we've seen in three projects generalises to a principle about distributed systems, and our earlier belief about consistency was wrong" in Layer 5? You can create artefacts, measure them, and improve efficiency. But "generalises to a principle" (abstraction beyond technology), "belief" (knowledge claim), and "was wrong" (correction of knowledge) have no Layer 5 representation. Building is not knowing.

### The Transition

**Physical → Symbolic**

Concrete artefacts become abstract information. The fundamental new capacity: creating symbols that stand for things, encoding meaning in language, transmitting through channels, and transforming data through computation.

### Base Operations

What can an information processor DO that a builder cannot?

1. **Represent** — create symbols that stand for things
2. **Encode** — transform information between representations
3. **Transmit** — move information through channels
4. **Compute** — transform data algorithmically

### Semantic Dimensions

| Dimension | Values | What it distinguishes |
|-----------|--------|-----------------------|
| **Level** | Concrete (specific instances) / Abstract (general principles) | Particular or general? |
| **Validity** | Unverified (claimed) / Verified (confirmed) / Retracted (corrected) | What's the epistemic status? |
| **Direction** | Constructive (building information) / Critical (challenging information) | Adding to or questioning what's known? |
| **Persistence** | Transient (in motion) / Durable (stored) | Is this information in transit or at rest? |

### Decomposition

**Group 0 — Representation** (how meaning is encoded)

| Primitive | Level | Validity | Direction | Persistence | What it does |
|-----------|-------|----------|-----------|-------------|--------------|
| **Symbol** | Concrete | Unverified | Constructive | Durable | A token that stands for something else |
| **Language** | Abstract | Unverified | Constructive | Durable | A system of symbols with grammar and semantics |
| **Encoding** | Concrete | Unverified | Constructive | Transient | Transforming information between representations |
| **Record** | Concrete | Verified | Constructive | Durable | A durable capture of information |

Representation lifecycle: symbols are created (Symbol) → organised into languages (Language) → transformed between formats (Encoding) → and durably captured (Record). This is how raw events become structured information.

**Group 1 — Dynamics** (how information moves)

| Primitive | Level | Validity | Direction | Persistence | What it does |
|-----------|-------|----------|-----------|-------------|--------------|
| **Channel** | Concrete | Unverified | Constructive | Transient | A path through which information flows |
| **Copy** | Concrete | Verified | Constructive | Durable | A faithful reproduction of information |
| **Noise** | Concrete | Unverified | Critical | Transient | Unwanted distortion in information |
| **Redundancy** | Concrete | Verified | Constructive | Durable | Duplication that protects against loss |

Dynamics lifecycle: information flows through channels (Channel) → is faithfully reproduced (Copy) → distortion is detected (Noise) → and duplication protects against loss (Redundancy). The critical direction (Noise) is what prevents information degradation.

**Group 2 — Transformation** (how information becomes something new)

| Primitive | Level | Validity | Direction | Persistence | What it does |
|-----------|-------|----------|-----------|-------------|--------------|
| **Data** | Concrete | Verified | Constructive | Durable | Structured information ready for processing |
| **Computation** | Abstract | Unverified | Constructive | Transient | Processing data to produce new information |
| **Algorithm** | Abstract | Verified | Constructive | Durable | A repeatable procedure for computation |
| **Entropy** | Abstract | Verified | Critical | Durable | The measure of disorder and information content |

Transformation lifecycle: data is structured (Data) → processed to produce new information (Computation) → via repeatable procedures (Algorithm) → and disorder is measured (Entropy). Entropy is the fundamental limit — it measures what can and cannot be known from a given dataset.

### Gap Analysis

| Behavior | Maps to | Notes |
|----------|---------|-------|
| "Let X represent the user's trust level" | Symbol | |
| "We use JSON as our data language" | Language | |
| "Convert the report to CSV format" | Encoding | |
| "Log this event for audit" | Record | |
| "Send the message over the event bus" | Channel | |
| "Replicate the database to the backup" | Copy | |
| "The signal is corrupted" | Noise | |
| "We keep three replicas for safety" | Redundancy | |
| "The dataset contains 10M rows" | Data | |
| "Calculate the aggregate statistics" | Computation | |
| "Use binary search to find the record" | Algorithm | |
| "The system's entropy is increasing" | Entropy | |
| Knowledge graph construction | Symbol + Language + Record | Composition |
| Scientific method | Data + Computation + Algorithm + Record | Composition |
| Error detection and correction | Noise + Redundancy + Copy | Composition |

**No gaps found.**

### Completeness Argument

1. **Dimensional coverage:** The {level, validity, direction, persistence} space is covered. The constructive/critical distinction ensures the system can both build information and detect its degradation — without the critical dimension (Noise, Entropy), information quality cannot be assessed.

2. **Information theory coverage:** Shannon's information theory (Channel, Noise, Redundancy, Entropy), computation theory (Data, Computation, Algorithm), and semiotics (Symbol, Language, Encoding) are all represented.

3. **Layer boundary:** None of these require concepts from Layer 7 (Ethics). Information models what *is* — symbols, data, computation. Whether those facts are *good* or *just* is Layer 7's gap. A system can process data about an action without reasoning about whether the action itself is right.

---

## Primitive Specifications

Full specifications in `docs/primitives.md` (Layer 6 section).

## Product Graph

Layer 6 maps to the **Knowledge Graph** — verified, provenanced knowledge with integrity guarantees. See `docs/product-layers.md`.

## Reference

- `docs/derivation-method.md` — The derivation method
- `docs/primitives.md` — Full specifications
- `docs/product-layers.md` — Knowledge Graph
- `docs/tests/primitives/06-research-integrity.md` — Related integration test scenario
