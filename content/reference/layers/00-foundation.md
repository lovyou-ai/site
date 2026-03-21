# Layer 0: Foundation

## Overview

45 primitives in 11 groups (Group 0 has 5, Groups 1-10 have 4 each). The irreducible computational foundations — everything above is derived from gaps in this layer.

Every primitive below declares:
- **Subscribes to** — event type prefixes it listens for
- **Emits** — event types it produces (with content schemas)
- **Depends on** — interfaces it requires
- **State** — what it stores between ticks
- **Mechanical/Intelligent** — whether it starts with deterministic logic or needs LLM

See `docs/interfaces.md` for all referenced types (`EventID`, `ActorID`, `EdgeType`, etc.).

---

## Group 0 — Core

The event graph itself.

### Event

Creates and validates events. The gatekeeper for everything entering the graph.

| | |
|---|---|
| **Subscribes to** | `*` (all — validates every event) |
| **Emits** | `event.invalid` { eventID: EventID, reason: string } |
| **Depends on** | Hash, Signature |
| **State** | `lastEventID`: EventID, `eventCount`: int |
| **Mechanical** | Yes |

Validates: required fields present, UUID v7 format, causes exist in store, hash computes correctly, signature verifies against source actor's public key.

### EventStore

Persistence. Wraps the Store interface for the tick engine.

| | |
|---|---|
| **Subscribes to** | `store.*` |
| **Emits** | `store.appended` { eventID: EventID }, `store.error` { operation: string, error: string } |
| **Depends on** | Store |
| **State** | `eventCount`: int, `lastHash`: Hash |
| **Mechanical** | Yes |

### Clock

Temporal ordering. Tick counting. Timestamps.

| | |
|---|---|
| **Subscribes to** | `tick.*` |
| **Emits** | `clock.tick` { tick: int, timestamp: time, elapsed: duration } |
| **Depends on** | (none — reads system clock) |
| **State** | `currentTick`: int, `startTime`: time, `lastTickTime`: time |
| **Mechanical** | Yes |

### Hash

Cryptographic hashing. Chain integrity.

| | |
|---|---|
| **Subscribes to** | `event.created` |
| **Emits** | `hash.computed` { eventID: EventID, hash: Hash, prevHash: Hash } |
| **Depends on** | (none — pure computation) |
| **State** | `chainHead`: Hash |
| **Mechanical** | Yes |

Input: canonical string (`prev_hash|id|type|source|conversation_id|timestamp_nanos|content_json`). Output: SHA-256 `Hash`.

### Self

System identity and message routing.

| | |
|---|---|
| **Subscribes to** | `message.*`, `system.*` |
| **Emits** | `self.routed` { messageID: EventID, target: PrimitiveID }, `self.identity` { actorID: ActorID, systemURI: SystemURI } |
| **Depends on** | IIdentity, IActorStore, Primitive Registry |
| **State** | `systemActorID`: ActorID, `systemURI`: SystemURI |
| **Intelligent** | Yes — routes messages to relevant domain primitives |

---

## Group 1 — Causality

Why things happened.

### CausalLink

Establishing and validating causal edges between events.

| | |
|---|---|
| **Subscribes to** | `event.created` |
| **Emits** | `causal.linked` { from: EventID, to: EventID }, `causal.invalid` { eventID: EventID, missingCause: EventID } |
| **Depends on** | Store |
| **State** | (stateless — validates on each event) |
| **Mechanical** | Yes |

Validates that all declared causes exist. Rejects events with missing causal predecessors (except Bootstrap).

### Ancestry

Traversing causal chains upward.

| | |
|---|---|
| **Subscribes to** | `query.ancestors` |
| **Emits** | `query.ancestors.result` { eventID: EventID, ancestors: []EventID, depth: int } |
| **Depends on** | Store (Ancestors method) |
| **State** | (stateless) |
| **Mechanical** | Yes |

Input: `{ eventID: EventID, maxDepth: int }`. Output: ordered list of ancestor events.

### Descendancy

Traversing causal chains downward.

| | |
|---|---|
| **Subscribes to** | `query.descendants` |
| **Emits** | `query.descendants.result` { eventID: EventID, descendants: []EventID, depth: int } |
| **Depends on** | Store (Descendants method) |
| **State** | (stateless) |
| **Mechanical** | Yes |

Input: `{ eventID: EventID, maxDepth: int }`. Output: ordered list of descendant events.

### FirstCause

Finding root causes. Identifying bootstrap events.

| | |
|---|---|
| **Subscribes to** | `query.firstcause` |
| **Emits** | `query.firstcause.result` { eventID: EventID, firstCause: EventID, chainLength: int } |
| **Depends on** | Store (Ancestors method), IIntelligence (for complex causal analysis) |
| **State** | (stateless) |
| **Mostly mechanical** | Intelligent for complex multi-branch causal analysis |

Walks ancestors to the root. For simple linear chains, mechanical. For complex causal DAGs with multiple roots, may use IIntelligence to determine the most significant root cause.

---

## Group 2 — Identity

Who did things.

### ActorID

Actor identity management. Keypair association.

| | |
|---|---|
| **Subscribes to** | `actor.register`, `actor.update` |
| **Emits** | `actor.registered` { actorID: ActorID, publicKey: PublicKey, type: ActorType }, `actor.updated` { actorID: ActorID, changes: ActorUpdate } |
| **Depends on** | IActorStore |
| **State** | `actorCount`: int |
| **Mechanical** | Yes |

Creates `ActorID` from public key. Ensures uniqueness. Manages the link between keypair and identity.

### ActorRegistry

Actor registration, lookup, lifecycle.

| | |
|---|---|
| **Subscribes to** | `actor.*` |
| **Emits** | `actor.suspended` { actorID: ActorID, reason: EventID }, `actor.memorial` { actorID: ActorID, reason: EventID }, `actor.reactivated` { actorID: ActorID } |
| **Depends on** | IActorStore |
| **State** | `activeCount`: int, `suspendedCount`: int, `memorialCount`: int |
| **Mechanical** | Yes |

Manages actor lifecycle transitions. Suspension and memorial are irreversible without explicit reactivation. Memorial preserves the actor's graph forever.

### Signature

Ed25519 signing of events and messages.

| | |
|---|---|
| **Subscribes to** | `event.created` (pre-append) |
| **Emits** | `signature.signed` { eventID: EventID, signer: ActorID } |
| **Depends on** | IIdentity |
| **State** | (stateless — signing is pure computation) |
| **Mechanical** | Yes |

Input: canonical event bytes, signer's private key. Output: Ed25519 signature bytes.

### Verify

Signature verification. Authentication.

| | |
|---|---|
| **Subscribes to** | `event.created`, `egip.envelope.received` |
| **Emits** | `signature.verified` { eventID: EventID, valid: bool, signer: ActorID }, `signature.failed` { eventID: EventID, reason: string } |
| **Depends on** | IActorStore (to look up public keys) |
| **State** | `verifiedCount`: int, `failedCount`: int |
| **Mechanical** | Yes |

Input: data, public key, signature. Output: bool. Rejects events with invalid signatures.

---

## Group 3 — Expectations

What should happen.

### Expectation

Defining what the system expects to happen after an event.

| | |
|---|---|
| **Subscribes to** | `*` (analyses events to create expectations) |
| **Emits** | `expectation.created` { expectation: Expectation }, `expectation.met` { expectationID: EventID, metBy: EventID } |
| **Depends on** | IIntelligence (to reason about what should happen next) |
| **State** | `pending`: map[EventID]Expectation |
| **Intelligent** | Yes |

Given an event, determines what should happen next and by when. E.g., after `authority.requested`, expects `authority.resolved` within the timeout period.

### Timeout

Tracking when expectations expire.

| | |
|---|---|
| **Subscribes to** | `expectation.created`, `clock.tick` |
| **Emits** | `expectation.expired` { expectationID: EventID, deadline: time } |
| **Depends on** | Clock |
| **State** | `deadlines`: map[EventID]time |
| **Mechanical** | Yes |

Watches the clock. When an expectation's deadline passes without being met, emits expiry event.

### Violation

Detecting when expectations are not met.

| | |
|---|---|
| **Subscribes to** | `expectation.expired` |
| **Emits** | `violation.detected` { violation: ViolationRecord } |
| **Depends on** | IIntelligence (to assess whether the violation is genuine or excusable) |
| **State** | `violations`: map[EventID]ViolationRecord |
| **Intelligent** | Yes |

Analyses expired expectations. Determines if the expectation was truly violated or if circumstances changed. Creates ViolationRecord with responsible actor and evidence.

### Severity

Assessing how serious a violation is.

| | |
|---|---|
| **Subscribes to** | `violation.detected` |
| **Emits** | `violation.assessed` { violationID: EventID, severity: SeverityLevel, impact: string } |
| **Depends on** | IIntelligence, ITrustModel |
| **State** | (stateless per violation) |
| **Intelligent** | Yes |

Input: ViolationRecord. Output: SeverityLevel + impact description. May trigger trust updates for the responsible actor.

---

## Group 4 — Trust

Who to believe.

### TrustScore

Maintaining trust scores for actors and systems.

| | |
|---|---|
| **Subscribes to** | `query.trust`, `trust.updated` |
| **Emits** | `trust.score` { actor: ActorID, metrics: TrustMetrics } |
| **Depends on** | ITrustModel, IActorStore |
| **State** | `scores`: map[ActorID]TrustMetrics |
| **Both** | Mechanical for queries, intelligent for complex scoring |

Input: IActor (via `ActorID`). Output: TrustMetrics (overall, by-domain, confidence, trend, decay rate).

### TrustUpdate

Processing events that affect trust.

| | |
|---|---|
| **Subscribes to** | `violation.assessed`, `task.completed`, `authority.resolved`, `corroboration.found`, `contradiction.found` |
| **Emits** | `trust.updated` { actor: ActorID, previous: TrustMetrics, current: TrustMetrics, cause: EventID } |
| **Depends on** | ITrustModel, IActorStore |
| **State** | (delegates to ITrustModel) |
| **Intelligent** | Yes — determines how much trust changes |

Edges created: `AddEdge { type: EdgeType.Trust, from: system, to: actor, weight: newScore }`.

### Corroboration

Detecting when multiple sources agree.

| | |
|---|---|
| **Subscribes to** | `*` (looks for agreement patterns) |
| **Emits** | `corroboration.found` { claim: string, sources: []ActorID, confidence: Score, evidence: []EventID } |
| **Depends on** | IIntelligence |
| **State** | `claims`: map[string][]EventID |
| **Intelligent** | Yes |

When multiple independent actors make consistent claims, corroboration increases confidence and trust.

### Contradiction

Detecting when sources disagree.

| | |
|---|---|
| **Subscribes to** | `*` (looks for disagreement patterns) |
| **Emits** | `contradiction.found` { claim: string, sideA: []ActorID, sideB: []ActorID, evidence: []EventID } |
| **Depends on** | IIntelligence |
| **State** | `claims`: map[string][]EventID |
| **Intelligent** | Yes |

When actors make conflicting claims, contradiction triggers investigation and potential trust updates.

---

## Group 5 — Confidence

How sure we are.

### Confidence

Maintaining confidence levels for beliefs and states.

| | |
|---|---|
| **Subscribes to** | `query.confidence`, `evidence.*`, `corroboration.*`, `contradiction.*` |
| **Emits** | `confidence.level` { proposition: string, confidence: Score, evidence: []EventID } |
| **Depends on** | Evidence primitive |
| **State** | `levels`: map[string]Score |
| **Both** | Mechanical for queries, intelligent for assessment |

### Evidence

Tracking evidence for and against propositions.

| | |
|---|---|
| **Subscribes to** | `*` (identifies events that constitute evidence) |
| **Emits** | `evidence.for` { proposition: string, eventID: EventID, weight: Weight }, `evidence.against` { proposition: string, eventID: EventID, weight: Weight } |
| **Depends on** | IIntelligence |
| **State** | `propositions`: map[string]{ for: []EventID, against: []EventID } |
| **Intelligent** | Yes |

### Revision

Updating beliefs when new evidence arrives.

| | |
|---|---|
| **Subscribes to** | `evidence.for`, `evidence.against` |
| **Emits** | `belief.revised` { proposition: string, previous: Score, current: Score, cause: EventID } |
| **Depends on** | IIntelligence |
| **State** | `beliefs`: map[string]Score |
| **Intelligent** | Yes — Bayesian-style revision |

### Uncertainty

Quantifying and communicating what the system doesn't know.

| | |
|---|---|
| **Subscribes to** | `confidence.*`, `query.uncertainty` |
| **Emits** | `uncertainty.report` { unknowns: []string, confidenceGaps: map[string]Score } |
| **Depends on** | Confidence, Evidence |
| **State** | `gaps`: map[string]Score |
| **Intelligent** | Yes |

---

## Group 6 — Instrumentation

What we're watching.

### InstrumentationSpec

Defining what should be monitored.

| | |
|---|---|
| **Subscribes to** | `system.*`, `instrumentation.*` |
| **Emits** | `instrumentation.defined` { spec: { eventTypes: []string, actors: []ActorID, primitives: []PrimitiveID } } |
| **Depends on** | IIntelligence |
| **State** | `specs`: []InstrumentationSpec |
| **Intelligent** | Yes |

### CoverageCheck

Verifying that instrumentation covers all critical paths.

| | |
|---|---|
| **Subscribes to** | `instrumentation.defined`, `clock.tick` (periodic check) |
| **Emits** | `instrumentation.coverage` { covered: Score, uncovered: []string } |
| **Depends on** | InstrumentationSpec, Primitive Registry |
| **State** | `coveragePercentage`: Score |
| **Both** | |

### Gap

Identifying what's not being monitored.

| | |
|---|---|
| **Subscribes to** | `instrumentation.coverage` |
| **Emits** | `instrumentation.gap` { description: string, severity: SeverityLevel, recommendation: string } |
| **Depends on** | IIntelligence |
| **State** | `gaps`: []string |
| **Intelligent** | Yes |

### Blind

Detecting blind spots — things the system can't see.

| | |
|---|---|
| **Subscribes to** | `instrumentation.gap`, `violation.detected` (violations from ungapped areas suggest blind spots) |
| **Emits** | `instrumentation.blind` { description: string, evidence: []EventID } |
| **Depends on** | IIntelligence |
| **State** | `blindSpots`: []string |
| **Intelligent** | Yes |

---

## Group 7 — Query

How we find things.

### PathQuery

Finding paths between events or actors in the graph.

| | |
|---|---|
| **Subscribes to** | `query.path` |
| **Emits** | `query.path.result` { from: EventID, to: EventID, path: []EventID, length: int } |
| **Depends on** | Store |
| **State** | (stateless) |
| **Mechanical** | Yes |

Input: `{ from: EventID, to: EventID, maxDepth: int }`. Output: shortest causal path between two events.

### SubgraphExtract

Extracting relevant subgraphs for analysis.

| | |
|---|---|
| **Subscribes to** | `query.subgraph` |
| **Emits** | `query.subgraph.result` { root: EventID, events: []EventID, edges: []Edge } |
| **Depends on** | Store, IIntelligence (for relevance determination) |
| **State** | (stateless) |
| **Both** | |

Input: `{ root: EventID, criteria: SubgraphCriteria }`. Output: subgraph of related events and edges.

### Annotate

Adding metadata to existing events.

| | |
|---|---|
| **Subscribes to** | `annotation.create` |
| **Emits** | `annotation.created` { target: EventID, key: string, value: any, annotator: ActorID } |
| **Depends on** | Store |
| **State** | (stateless — annotations are events) |
| **Both** | |

Creates an Annotation edge from the annotation event to the target event. Edges created: `AddEdge { type: EdgeType.Annotation, from: annotationEventID, to: targetEventID }`.

### Timeline

Constructing temporal views of event sequences.

| | |
|---|---|
| **Subscribes to** | `query.timeline` |
| **Emits** | `query.timeline.result` { events: []EventID, start: time, end: time } |
| **Depends on** | Store |
| **State** | (stateless) |
| **Mechanical** | Yes |

Input: `{ start: time, end: time, filter: { type: Option<EventType>, source: Option<ActorID> } }`. Output: ordered events within the time range.

---

## Group 8 — Integrity

Is the graph intact.

### HashChain

Managing the linear hash chain.

| | |
|---|---|
| **Subscribes to** | `event.created` |
| **Emits** | `chain.extended` { eventID: EventID, hash: Hash, position: int } |
| **Depends on** | Hash primitive, Store |
| **State** | `chainHead`: Hash, `chainLength`: int |
| **Mechanical** | Yes |

Maintains the append-only linear chain. Every event gets a PrevHash linking to the previous event's hash.

### ChainVerify

Verifying chain integrity.

| | |
|---|---|
| **Subscribes to** | `query.verify`, `clock.tick` (periodic verification) |
| **Emits** | `chain.verified` { valid: bool, length: int, duration: duration }, `chain.broken` { position: int, expected: Hash, actual: Hash } |
| **Depends on** | Store (VerifyChain method), Hash |
| **State** | `lastVerified`: time, `lastResult`: bool |
| **Mechanical** | Yes |

Walks the chain from genesis to head, recomputing each hash. Any tampering breaks the chain at the modification point.

### Witness

External witnessing of chain state (for cross-system trust).

| | |
|---|---|
| **Subscribes to** | `egip.proof.request`, `clock.tick` (periodic witness generation) |
| **Emits** | `witness.created` { chainHead: Hash, chainLength: int, timestamp: time, signature: Signature } |
| **Depends on** | IIdentity, HashChain |
| **State** | `witnesses`: []{ hash: Hash, time: time, signature: Signature } |
| **Mechanical** | Yes |

Creates signed attestations of chain state. Used by EGIP for cross-system integrity verification.

### IntegrityViolation

Detecting and responding to integrity breaches.

| | |
|---|---|
| **Subscribes to** | `chain.broken`, `signature.failed` |
| **Emits** | `integrity.violated` { type: IntegrityViolationType, description: string, severity: SeverityLevel, evidence: []EventID } |
| **Depends on** | IAuthorityChain (may escalate to human) |
| **State** | `violations`: []{ time: time, type: IntegrityViolationType, resolved: bool } |
| **Both** | |

Integrity violations are always SeverityLevel.Critical. May trigger actor suspension or system quarantine.

---

## Group 9 — Deception

Is someone lying.

### Pattern

Recognising behavioural patterns across events.

| | |
|---|---|
| **Subscribes to** | `*` (analyses all events for patterns) |
| **Emits** | `pattern.detected` { actor: ActorID, pattern: string, confidence: Score, evidence: []EventID } |
| **Depends on** | IIntelligence, Store |
| **State** | `patterns`: map[ActorID][]{ pattern: string, frequency: int } |
| **Intelligent** | Yes |

### DeceptionIndicator

Identifying potential deception signals.

| | |
|---|---|
| **Subscribes to** | `pattern.detected`, `contradiction.found`, `signature.failed` |
| **Emits** | `deception.indicator` { actor: ActorID, indicator: string, confidence: Score, evidence: []EventID } |
| **Depends on** | IIntelligence, Pattern |
| **State** | `indicators`: map[ActorID][]string |
| **Intelligent** | Yes |

Indicators: inconsistent timestamps, contradictory claims, replay patterns, abnormal event velocity.

### Suspicion

Maintaining and updating suspicion levels.

| | |
|---|---|
| **Subscribes to** | `deception.indicator` |
| **Emits** | `suspicion.updated` { actor: ActorID, previous: Score, current: Score, indicators: []string } |
| **Depends on** | ITrustModel |
| **State** | `levels`: map[ActorID]Score |
| **Intelligent** | Yes |

Suspicion is separate from trust — trust is what you've earned, suspicion is active investigation. High suspicion may trigger quarantine.

### Quarantine

Isolating suspected-compromised events or actors.

| | |
|---|---|
| **Subscribes to** | `suspicion.updated`, `integrity.violated` |
| **Emits** | `quarantine.actor` { actor: ActorID, reason: EventID }, `quarantine.event` { eventID: EventID, reason: EventID }, `quarantine.actor.lifted` { actor: ActorID, reason: EventID }, `quarantine.event.lifted` { eventID: EventID, reason: EventID } |
| **Depends on** | IAuthorityChain (requires authority to quarantine), IActorStore |
| **State** | `quarantined`: { actors: []ActorID, events: []EventID } |
| **Both** | |

Quarantined actors can't emit events. Quarantined events are excluded from queries but not deleted (the graph is append-only). Lifting quarantine requires authority approval.

---

## Group 10 — Health

Is the system well.

### GraphHealth

Overall health assessment of the event graph.

| | |
|---|---|
| **Subscribes to** | `clock.tick` (periodic health check) |
| **Emits** | `health.report` { overall: Score, chainIntegrity: bool, primitiveHealth: map[PrimitiveID]Score, activeActors: int, eventRate: float64 } |
| **Depends on** | All other health-related primitives, Store |
| **State** | `lastReport`: time, `metrics`: HealthReportContent |
| **Both** | |

### Invariant

Defining system invariants that must hold.

| | |
|---|---|
| **Subscribes to** | `system.start`, `invariant.define` |
| **Emits** | `invariant.defined` { name: InvariantName, description: string, checker: PrimitiveID } |
| **Depends on** | IIntelligence |
| **State** | `invariants`: map[InvariantName]{ description: string, checker: PrimitiveID } |
| **Intelligent** | Yes |

The 10 system invariants (CAUSALITY, INTEGRITY, OBSERVABLE, SELF-EVOLVE, DIGNITY, TRANSPARENT, CONSENT, AUTHORITY, VERIFY, RECORD) are loaded at bootstrap.

### InvariantCheck

Verifying invariants are maintained.

| | |
|---|---|
| **Subscribes to** | `clock.tick` (periodic), `event.created` (spot checks) |
| **Emits** | `invariant.passed` { name: InvariantName }, `invariant.failed` { name: InvariantName, description: string, evidence: []EventID, severity: SeverityLevel } |
| **Depends on** | Invariant, Store |
| **State** | `results`: map[string]{ lastChecked: time, passed: bool } |
| **Both** | |

Invariant failures are always SeverityLevel.Critical. May trigger system halt if CAUSALITY or INTEGRITY fail.

### Bootstrap

System initialisation and recovery.

| | |
|---|---|
| **Subscribes to** | `system.start` |
| **Emits** | `system.bootstrapped` { timestamp: time, actorID: ActorID, chainGenesis: Hash } |
| **Depends on** | Store, IActorStore, IIdentity |
| **State** | `bootstrapped`: bool, `genesisEvent`: EventID |
| **Mechanical** | Yes |

The only event that has no causes. Creates the system actor, genesis event, and initialises the hash chain. If recovering from a crash, verifies chain integrity before activating Layer 0 primitives.

---

## Layer Activation

Layer 0 primitives activate on system start. They must be stable before any Layer 1 primitive activates. "Stable" means: in `LifecycleState.Active` with consistent outputs for at least one full tick cycle.

## Implementation Notes

- Mechanical primitives should start with mostly deterministic decision tree branches
- Intelligent primitives should start with `NeedsLLM = true` leaves
- All primitives must emit events for their significant actions
- All primitives must declare event subscriptions
- Default cadence for Layer 0: 1 (every tick)
- Every `ActorID` passed must resolve via `IActorStore.Get()` — no orphan references
- Every `EventID` referenced in causes must exist in Store — no dangling pointers
