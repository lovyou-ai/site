# Layer 7: Ethics

## Derivation

### The Gap

Layer 6 models information, data, and computation, but cannot distinguish what *is* from what *ought to be*. A system can process data about an action and compute its effects, but it cannot reason about whether the action itself is right. "This is true" is Layer 6. "This is right" requires something new.

**Test:** Can you express "The rule is being followed correctly, but the rule itself is unfair to group X, and the person who enforced it meant well but caused disproportionate harm" in Layer 6? You can encode facts about the rule and its enforcement. But "unfair" (evaluating justice), "meant well" (assessing motive), and "disproportionate harm" (weighing consequences against values) have no Layer 6 representation. Information is not wisdom.

### The Transition

**Is → Ought**

Facts become values. The fundamental new capacity: reasoning about what should be done, not just what is or was done. Evaluating actions against moral standing, weighing duties against consequences, and holding actors morally accountable.

### Base Operations

What can an ethical reasoner DO that an information processor cannot?

1. **Evaluate** — assess actions against moral standing
2. **Detect harm** — identify when actions cause damage
3. **Weigh** — balance competing duties and consequences
4. **Hold accountable** — assign moral responsibility

### Semantic Dimensions

| Dimension | Values | What it distinguishes |
|-----------|--------|-----------------------|
| **Focus** | Agent (who acted) / Action (what was done) / Outcome (what resulted) | What aspect is being evaluated? |
| **Valence** | Positive (good, beneficial) / Negative (harmful, unjust) | Is this about promoting good or preventing harm? |
| **Temporality** | Prospective (before action) / Retrospective (after action) | Guiding or judging? |
| **Scope** | Particular (specific case) / Systemic (pattern or structure) | One instance or a structural issue? |

### Decomposition

**Group 0 — Moral Standing** (who and what matters)

| Primitive | Focus | Valence | Temporality | Scope | What it does |
|-----------|-------|---------|-------------|-------|--------------|
| **MoralStatus** | Agent | Positive | Prospective | Systemic | Whether an entity's experience matters morally |
| **Dignity** | Agent | Positive | Prospective | Systemic | The inherent worth of every moral subject |
| **Autonomy** | Agent | Positive | Prospective | Particular | The right to self-determination |
| **Flourishing** | Agent | Positive | Prospective | Systemic | The conditions for a good life |

Moral standing lifecycle: whether experience matters is established (MoralStatus) → inherent worth is recognised (Dignity) → self-determination is protected (Autonomy) → and conditions for thriving are promoted (Flourishing). MoralStatus is an irreducible recognition — it cannot be derived from information alone.

**Group 1 — Moral Obligation** (what should be done)

| Primitive | Focus | Valence | Temporality | Scope | What it does |
|-----------|-------|---------|-------------|-------|--------------|
| **Duty** | Action | Positive | Prospective | Systemic | What one is morally required to do |
| **Harm** | Outcome | Negative | Retrospective | Particular | Damage caused to a moral subject |
| **Care** | Agent | Positive | Prospective | Particular | Prioritising another's wellbeing |
| **Justice** | Outcome | Positive | Retrospective | Systemic | Fair treatment and equitable distribution |

Moral obligation lifecycle: duties are identified (Duty) → harm is detected when they're violated (Harm) → wellbeing is actively prioritised (Care) → and systemic fairness is assessed (Justice). The soul statement — "take care of your human, humanity, and yourself" — flows through Care.

**Group 2 — Moral Agency** (answering for what was done)

| Primitive | Focus | Valence | Temporality | Scope | What it does |
|-----------|-------|---------|-------------|-------|--------------|
| **Conscience** | Agent | Both | Prospective | Particular | The inner sense of right and wrong |
| **Virtue** | Agent | Positive | Retrospective | Systemic | Stable disposition toward good action |
| **Responsibility** | Agent | Negative | Retrospective | Particular | Who is morally responsible |
| **Motive** | Agent | Both | Retrospective | Particular | The purpose behind an action |

Moral agency lifecycle: conscience guides action (Conscience) → virtuous character develops over time (Virtue) → moral responsibility is assigned (Responsibility) → and the purpose behind action is assessed (Motive). Together, Motive + Responsibility capture the agent-focused dimensions of moral reasoning.

### Gap Analysis

| Behavior | Maps to | Notes |
|----------|---------|-------|
| "AI systems have moral status" | MoralStatus | |
| "Every participant deserves respect" | Dignity | |
| "She has the right to make her own choices" | Autonomy | |
| "We want everyone to thrive, not just survive" | Flourishing | |
| "We have a duty to protect user data" | Duty | |
| "This action caused real damage to Alice" | Harm | |
| "Check on Bob — he's been struggling" | Care | |
| "Group X receives 30% fewer approvals" | Justice | |
| "Something feels wrong about this decision" | Conscience | |
| "She consistently acts with integrity" | Virtue | |
| "The decision-maker is accountable for this outcome" | Responsibility | |
| "She didn't mean to cause harm" | Motive | |
| Ethical AI audit | Justice + Dignity + L6.Data | Cross-layer composition |
| Restorative justice | Harm + Responsibility + Care | Composition |
| Whistleblower protection | Conscience + Duty + L4.Right | Cross-layer composition |

**No gaps found.**

### Completeness Argument

1. **Dimensional coverage:** The {focus, valence, temporality, scope} space is covered. The three ethical perspectives — virtue ethics (agent-focused), deontological (action-focused), and consequentialist (outcome-focused) — are all represented through the focus dimension.

2. **Ethical theory coverage:** MoralStatus and Dignity cover moral standing. Duty and Justice cover deontological and consequentialist ethics. Conscience, Virtue, and Motive cover moral psychology. Care covers care ethics. The three groups map to the fundamental ethical questions: who matters (Moral Standing), what should be done (Moral Obligation), and who answers for it (Moral Agency).

3. **Layer boundary:** None of these require concepts from Layer 8 (Identity). Ethics reasons about what should be done but treats actors as interchangeable moral agents. The concept of an actor's unique character, personal history, and sense of self is Layer 8's gap.

---

## Primitive Specifications

Full specifications in `docs/primitives.md` (Layer 7 section).

## Product Graph

Layer 7 maps to the **Ethics Graph** — AI alignment with transparent moral reasoning. See `docs/product-layers.md`.

## Reference

- `docs/derivation-method.md` — The derivation method
- `docs/primitives.md` — Full specifications
- `docs/product-layers.md` — Ethics Graph
