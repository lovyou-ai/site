# The Friction

*Everything that could stop this from working. Honestly, without flinching.*

Matt Searles (+Claude) · March 2026

---

The Weight named the suffering. The Transition mapped the construction plan. If you read those and thought "yes, but can this actually work?" — good. That's the right question. This post is the honest answer.

Every serious objection encountered — from AI researchers, engineers, and Google's Gemini — is catalogued here. Some have answers. Some have partial answers. Some remain genuinely unsolved. The post sorts them honestly into those three categories.

## The Oracle Problem

The graph records events but cannot verify their truth at entry. A builder logging "high-quality concrete poured" creates an immutable record, whether accurate or false.

The framework doesn't verify truth immediately but makes lies discoverable over time. When concrete fails inspection, that event gets recorded. Contradictions accumulate into visible patterns. A builder cannot relocate and start fresh; their track record follows them on the chain.

**Weakness:** This only works if contradicting events eventually surface. If inspections are corrupt or buildings don't fail for decades, lies persist.

**Verdict:** Solved enough to ship, though safety-critical domains need additional verification layers beyond what the framework currently specifies.

## The Panopticon

Total causal traceability equals total surveillance. Every action by every agent becomes traceable back to authorization, making the event graph "the most comprehensive surveillance apparatus ever built."

Three architectural safeguards exist: intent sanctuary (recording actions, not thoughts), causality boundaries (determining where responsibility ends), and access control (layered, permission-gated visibility). However, governments with sufficient power can demand access regardless of architectural constraints.

**Weakness:** Access control depends entirely on governance enforcement. The data's existence creates incentives for state-level demands that wouldn't otherwise exist.

**Verdict:** Solved in principle, needs concrete work. Mechanisms for enforcing boundaries against state demands remain unspecified.

## Goodhart's Law

When reputation becomes a measurable target, people optimize for the measurement. They perform ethical acts strategically, gaming the system like websites gaming search rankings.

However, causal chain analysis distinguishes genuine commitment from strategic performance over time. The strategic volunteer disappears when visibility does; the genuine volunteer persists unseen. Patterns emerge at scale.

**Weakness:** "Over time, at scale" demands much. Early adoption will include system gamers. Convergence toward truth versus exploitation speed is empirically unknowable pre-deployment.

**Verdict:** Solved enough to ship. The causal structure provides better detection tools than existing reputation systems.

## The Memory Problem

Immutable records preserve mistakes permanently, contradicting human psychology's need for forgetting and fresh starts.

The framework includes forgiveness as a primitive event, making repair visible alongside wounds. Human memory still degrades naturally. The graph supplements rather than replaces human memory.

**Weakness:** Archives have audiences. Future employers and partners can query old conflicts. The "right to be forgotten" (enshrined in European law) conflicts with append-only chains.

**Verdict:** Solved enough to ship, with caveats. The tension between permanent records and human healing may require policy solutions beyond architecture.

## The Scalability Wall

Billions of AI agents generating events continuously create unprecedented data volumes. Graph databases currently struggle at millions of nodes. Traversing causal chains across billions in real time may exceed current computational capacity.

The possible solution involves neural networks learning graph structure to predict chains rather than traversing them. This remains speculative research.

**Weakness:** "A research direction" isn't an answer. If unsolvable at civilizational scale, the framework works locally but not globally.

**Verdict:** Genuinely unsolved. Company and community deployment works today; global scalability remains an open research question.

## The Bootstrapping Paradox

Network value depends on network effects, but early adoption provides minimal value when no one else participates.

The framework provides standalone value even at single-company scale. Internal accountability, AI agent management, and institutional memory justify initial deployment independent of network formation.

**Weakness:** The leap from "useful tool" to "global accountability infrastructure" is enormous. Local value doesn't guarantee network scaling.

**Verdict:** Solved enough to ship. Local value is real and demonstrable.

## Governance of the Graph Itself

The framework governs everything except itself. Who decides primitives? Who updates the schema? Who resolves disputes about the graph's own rules?

Current stewardship exists, but transition to community governance needs formalization before the network outgrows its current structure. Open-source models (like Linux and TCP/IP) offer a path but face challenges unique to accountability infrastructure.

**Weakness:** Open-source governance works for technical protocols less well for systems with direct political implications. Powerful actors have incentives to capture, corrupt, or undermine accountability systems.

**Verdict:** Solved in principle, needs work. Political dimensions differ qualitatively from existing open-source projects.

## The AI Reliability Gap

The framework assumes AI agents can meaningfully consent, authorize, and maintain transparency. Current AI systems hallucinate, confabulate, and produce confident falsehoods.

Building accountability infrastructure now, before AI scales, creates space for failures to surface at manageable scale. Reliability gaps narrow as AI improves while infrastructure develops.

**Verdict:** Solved in principle, needs work. Early deployments reveal gaps between assumptions and reality; convergence depends on parallel AI improvement trajectories.

## Physical World Interop

Every digital-to-physical interface presents vulnerability. Broken sensors, angled cameras, drifting GPS coordinates undermine precision. Contaminated data creates false confidence.

Imperfect data on the graph remains superior to no data anywhere. Sensor accuracy improves over time. The framework handles this better than most systems by making sensor reliability itself auditable.

**Verdict:** Solved enough to ship. Imperfect inputs improve over time.

## The Energy Problem

Billions of agents generating events continuously requires immense compute and energy. If ecological damage from the infrastructure exceeds damage it prevents, net impact becomes negative.

Solar energy is already cheapest electricity in most markets. AI-accelerated research into storage, fusion, and grid efficiency compresses previously decade-spanning timelines. The framework should prioritize energy efficiency, but current trajectories suggest this won't be fatal.

**Verdict:** Solved in principle by trajectory. Energy transition moves faster than projections.

## The Adversarial Environment

Sybil attacks (deploying millions of fake agents) prove detectable because identity derives from meaningful behavior over time. Reputation is earned slowly; fake accounts lack depth and community participation.

Compromised high-reputation actors present an unsolvable insider threat. A state-cultivated sleeper agent with genuine credentials and relationships betrays unpredictably. The system must be designed for resilience rather than prevention.

**Verdict:** Sybil attacks solved. Insider threats genuinely unsolved and likely unsolvable generally. Framework needs resilience design.

## The Cultural Adoption Gap

Privacy represents genuine human need, not merely bad-behavior concealment. Desiring untracked spaces and private messiness reflects healthy psychology, not pathology.

The framework demands transparency from systems, not individuals. Personal life remains private; access controls protect it. Institutional transparency differs from individual surveillance.

**Weakness:** The individual-institution boundary is practically blurry. Work events become institutional; community participation becomes social.

**Verdict:** Solved in principle, needs work. Practical details matter enormously for adoption.

## The Wealth Transition

Trillions in existing wealth depend on current infrastructure failures. Mortgage, insurance, legal, and financial intermediary industries extract value from system opacity. These won't dismantle themselves.

Coexistence mechanics for old and new systems need detailed economic modeling currently lacking.

**Verdict:** Genuinely needs work. Policy-level solutions required beyond technical specification.

## Mono-Infrastructure Risk

Fragmented current infrastructure offers resilience through redundancy; unified infrastructure offers coherence at cost of vulnerability. If one system fails, everything fails.

The distributed event graph with inter-system protocol enables sovereign systems with independent graphs to communicate through signed envelopes, cross-graph references, and bilateral treaties. No single point of failure exists.

**Weakness:** Distributed systems have their own failure modes. Social dependency risks remain: if society depends on event graphs for accountability and the protocol is compromised, the accountability vacuum becomes total.

**Verdict:** Solved in architecture, unsolved in practice.

## Compound Friction

Every friction above is manageable alone. Whether the system survives all of them simultaneously during messy transition remains empirical.

Scalability limits deployment speed. Bootstrapping slows adoption. Cultural gaps reduce the user base. Wealth transition creates enemies. AI reliability produces failures. Physical interop undermines confidence. Adversarial environments test systems before readiness. Governance remains unresolved.

No historical framework faced all its frictions simultaneously and survived. The only strategy involves starting narrow, proving local value, and expanding as each friction gets addressed.

**Verdict:** Genuinely unsolved. This is a condition to survive, not a problem to solve. Full thirteen-layer deployment would fail immediately.

---

## The Honest Tally

**Solved enough to ship:** Oracle problem, Goodhart's Law, memory problem, Sybil attacks, physical interop, bootstrapping paradox, energy problem.

These friction points are real but match those every ambitious system faces. They don't require fundamental breakthroughs — just engineering, iteration, and deployment experience.

**Solved in principle, needs work:** Panopticon, governance, AI reliability gap, cultural adoption gap.

These have architectural answers but lack implementation detail. The panopticon tension remains most critical; get it wrong and the framework becomes the surveillance system it aimed to prevent.

**Genuinely unsolved:** Scalability at civilizational scale, insider threats, wealth transition, compound friction.

These could kill the project. Scalability remains a research problem possibly unsolvable with current technology. Insider threats are likely unsolvable generally. Wealth transition demands political solutions the framework cannot provide.

---

## Why Build Anyway

Four genuinely unsolved problems. Several needing real work. Compound friction risk multiplying everything.

Success isn't guaranteed. However, the alternative guarantees failure. The suffering described in "The Weight" continues: 138 million children in labour, 176,000 deaths of despair, 180 children in Minab, 4% conviction rates, persistent drug economies, widespread loneliness, two billion people governed algorithmically without consent.

The asymmetry: building and failing costs effort. Not building costs the weight.

Every friction here is real. None minimized or dismissed. The honest accounting: this might not work. The honest response: we must try anyway because not trying guarantees worse outcomes.
