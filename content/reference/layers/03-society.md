# Layer 3: Society

## Derivation

### The Gap

Layer 2 models dyadic exchange — two actors establishing terms, making offers, and forming agreements. But when a third actor joins, new phenomena appear that dyads can't express. What's acceptable isn't just what two people agree on — it's what the GROUP considers normal. Roles, norms, reputation within a community, inclusion, and exclusion don't exist at the dyad level.

**Test:** Can you express "The community considers this behavior unacceptable" in Layer 2? You can have individual Breaches and bilateral Agreements, but the concept of a norm — a shared expectation held by a group, not derived from any specific agreement — has no Layer 2 representation. Norms emerge from groups, not pairs.

### The Transition

**Dyad → Group**

Two actors in relation become N actors forming a collective. The fundamental new capacity: state that belongs to a group rather than to individuals or pairs.

### Base Operations

What can a group DO that a dyad cannot?

1. **Norm** — establish shared expectations not derived from bilateral agreement
2. **Role** — assign social positions with specific expectations
3. **Include/Exclude** — control group membership
4. **Govern** — make collective decisions

### Semantic Dimensions

| Dimension | Values | What it distinguishes |
|-----------|--------|-----------------------|
| **Scope of effect** | Individual (about one member) / Collective (about the group) | Does this affect one person or everyone? |
| **Formality** | Informal (implicit) / Formal (explicit, recorded) | Is this a tacit understanding or a stated rule? |
| **Direction** | Centripetal (toward group) / Centrifugal (toward individual) | Does this bring people together or push them apart? |
| **Temporality** | Emergent (from patterns) / Declared (from authority) | Did this arise naturally or was it established by decree? |

### Decomposition

**Group 0 — Collective Identity** (who we are together)

| Primitive | Effect | Formality | Direction | Temporality | What it does |
|-----------|--------|-----------|-----------|-------------|--------------|
| **Group** | Collective | Formal | Centripetal | Declared | A named collective with boundaries |
| **Membership** | Individual | Formal | Centripetal | Declared | Joining, belonging to, and leaving a group |
| **Role** | Individual | Formal | Centripetal | Declared | Named position with expectations |
| **Consent** | Individual | Formal | Centripetal | Declared | Voluntary agreement to participate or be governed |

Collective identity lifecycle: a group is formed (Group) → members join (Membership) → positions are assigned (Role) → and participation is voluntary (Consent). This covers how collectives constitute themselves.

**Group 1 — Social Order** (how we live together)

| Primitive | Effect | Formality | Direction | Temporality | What it does |
|-----------|--------|-----------|-----------|-------------|--------------|
| **Norm** | Collective | Formal | Centripetal | Declared | Shared behavioral expectation |
| **Reputation** | Individual | Informal | Centripetal | Emergent | Community standing from track record |
| **Sanction** | Individual | Formal | Centrifugal | Declared | Consequence for norm violation |
| **Authority** | Individual | Formal | Centripetal | Declared | Legitimate power to act on behalf of the group |

Social order lifecycle: expectations are established (Norm) → standing develops (Reputation) → violations have consequences (Sanction) → and some members have legitimate power (Authority). This is how societies regulate behavior.

**Group 2 — Collective Agency** (what we do together)

| Primitive | Effect | Formality | Direction | Temporality | What it does |
|-----------|--------|-----------|-----------|-------------|--------------|
| **Property** | Individual | Formal | Centrifugal | Declared | Recognised claim over resources |
| **Commons** | Collective | Formal | Centripetal | Declared | Shared resources that belong to the group |
| **Governance** | Collective | Formal | Centripetal | Declared | How collective decisions are made |
| **CollectiveAct** | Collective | Formal | Centripetal | Declared | Action taken by the group as a whole |

Collective agency lifecycle: individuals have recognised claims (Property) → some resources are shared (Commons) → decision-making structures exist (Governance) → and the group acts as one (CollectiveAct). This covers how groups exercise agency beyond individual action.

### Gap Analysis

| Behavior | Maps to | Notes |
|----------|---------|-------|
| "This is the backend team" | Group | |
| "Alice joined the project" | Membership | |
| "Alice is a moderator" | Role | |
| "Everyone agreed to the code of conduct" | Consent | |
| "We don't do that here" | Norm | |
| "Bob has high standing in the community" | Reputation | |
| "Violating the code of conduct results in suspension" | Sanction | |
| "The lead has authority to approve PRs" | Authority | |
| "This codebase belongs to the team" | Property | |
| "The shared library is a community resource" | Commons | |
| "Decisions are made by majority vote" | Governance | |
| "The community voted to adopt the new policy" | CollectiveAct | |
| Community votes on a proposal | Governance + CollectiveAct | Composition |
| Online community moderation | Norm + Sanction + Authority | Composition |
| Open source project governance | Group + Role + Governance + Commons | Composition |

**No gaps found.**

### Completeness Argument

1. **Dimensional coverage:** All meaningful combinations covered. The key distinction is formal/informal crossed with individual/collective — this generates the essential structures of social life.

2. **Sociological coverage:** The four pillars of social structure from sociology — norms, roles, membership, and stratification — are all represented. Collective agency (Property, Commons, Governance, CollectiveAct) captures how groups act beyond individual action.

3. **Layer boundary:** None of these require concepts from Layer 4 (Legal). Society has informal norms and sanctions, but not formal adjudication, precedent, or jurisdiction. Those are Layer 4's gap.

---

## Primitive Specifications

Full specifications in `docs/primitives.md` (Layer 3 section).

## Product Graph

Layer 3 maps to the **Social Graph** — user-owned social platform where communities set norms. See `docs/product-layers.md`.

## Reference

- `docs/derivation-method.md` — The derivation method
- `docs/primitives.md` — Full specifications
- `docs/product-layers.md` — Social Graph
- `docs/tests/primitives/04-community-governance.md` — Related integration test scenario
