# Layer 4: Legal

## Derivation

### The Gap

Layer 3 models groups with norms, roles, and sanctions. But norms are informal — they emerge from behaviour and consensus. When someone asks "what are the rules?", Layer 3 can point to conventions and norms, but cannot produce a codified law, a jurisdictional boundary, a formal precedent, or a right of appeal. The concept of "this rule applies here but not there" doesn't exist.

**Test:** Can you express "This rule was enacted by the governing body, applies within this jurisdiction, and overrides the conflicting norm from the neighbouring community" in Layer 3? You can have Norms, Governance, and Sanctions. But "enacted by authority" (codification), "applies here but not there" (jurisdiction), and "overrides" (precedent and interpretation) have no Layer 3 representation. Informal governance cannot express formal law.

### The Transition

**Informal → Formal**

Emergent norms become codified rules. The fundamental new capacity: explicit, authoritative, jurisdictionally scoped rules with formal processes for interpretation, adjudication, and reform.

### Base Operations

What can a formal legal system DO that informal norms cannot?

1. **Codify** — write laws with explicit conditions and consequences
2. **Adjudicate** — resolve disputes through formal process
3. **Enforce** — take action when laws are broken
4. **Reform** — change laws based on experience and precedent

### Semantic Dimensions

| Dimension | Values | What it distinguishes |
|-----------|--------|-----------------------|
| **Function** | Prescriptive (what the law says) / Procedural (how disputes are handled) / Corrective (what happens after) | What phase of the legal lifecycle? |
| **Authority** | Constitutive (creates the framework) / Operative (works within it) | Does this create the system or operate within it? |
| **Temporality** | Prospective (forward-looking) / Retrospective (backward-looking) | Does this look forward to prevent or backward to assess? |
| **Scope** | Universal (applies to all) / Particular (applies to specific case) | General rule or specific application? |

### Decomposition

**Group 0 — Codification** (prescriptive, constitutive)

| Primitive | Function | Authority | Temporality | Scope | What it does |
|-----------|----------|-----------|-------------|-------|--------------|
| **Law** | Prescriptive | Constitutive | Prospective | Universal | Formal, codified rule with conditions and consequences |
| **Right** | Prescriptive | Constitutive | Prospective | Universal | Fundamental protections that override other rules |
| **Contract** | Prescriptive | Operative | Prospective | Particular | Formalised bilateral agreement with legal force |
| **Liability** | Prescriptive | Operative | Retrospective | Particular | Legal responsibility for consequences |

Codification lifecycle: laws are enacted (Law) → fundamental protections are established (Right) → specific agreements are formalised (Contract) → and responsibility for consequences is assigned (Liability). This is how informal norms become formal law.

**Group 1 — Process** (procedural)

| Primitive | Function | Authority | Temporality | Scope | What it does |
|-----------|----------|-----------|-------------|-------|--------------|
| **DueProcess** | Procedural | Constitutive | Prospective | Universal | Ensuring procedural fairness |
| **Adjudication** | Procedural | Operative | Retrospective | Particular | Formal dispute resolution |
| **Remedy** | Corrective | Operative | Retrospective | Particular | Making things right after a legal wrong |
| **Precedent** | Prescriptive | Operative | Retrospective | Universal | Past decisions that inform future ones |

Process lifecycle: all proceedings must be fair (DueProcess) → disputes are formally resolved (Adjudication) → wrongs are corrected (Remedy) → and past decisions guide future ones (Precedent). DueProcess violations trigger immediate authority escalation.

**Group 2 — Sovereign Structure** (the framework itself)

| Primitive | Function | Authority | Temporality | Scope | What it does |
|-----------|----------|-----------|-------------|-------|--------------|
| **Jurisdiction** | Prescriptive | Constitutive | Prospective | Particular | Which laws apply where |
| **Sovereignty** | Prescriptive | Constitutive | Prospective | Universal | Ultimate authority within a domain |
| **Legitimacy** | Prescriptive | Constitutive | Retrospective | Universal | Whether authority is rightfully exercised |
| **Treaty** | Prescriptive | Constitutive | Prospective | Particular | Agreement between sovereign entities |

Sovereign structure lifecycle: boundaries of authority are defined (Jurisdiction) → ultimate authority is established (Sovereignty) → that authority is assessed for rightfulness (Legitimacy) → and sovereign entities relate to each other (Treaty). Treaty enables multi-system legal interaction.

### Gap Analysis

| Behavior | Maps to | Notes |
|----------|---------|-------|
| "Code of conduct, Section 3.2" | Law | |
| "Every member has the right to be heard" | Right | |
| "The service-level agreement between A and B" | Contract | |
| "You are liable for the data breach" | Liability | |
| "The accused was not given a chance to respond" | DueProcess (violated) | |
| "The moderator hears the dispute" | Adjudication | |
| "We owe Alice compensation for the wrongful ban" | Remedy | |
| "In the last similar case, we decided X" | Precedent | |
| "EU regulations apply to EU users" | Jurisdiction | |
| "This community governs itself" | Sovereignty | |
| "The election was fair and the leader is legitimate" | Legitimacy | |
| "The two communities agreed to mutual recognition" | Treaty | |
| Constitutional amendment | Law + Right + Sovereignty | Composition |
| Plea bargain | Adjudication + Remedy + Contract | Composition |
| Judicial review | Adjudication + Legitimacy + Right | Composition |

**No gaps found.**

### Completeness Argument

1. **Dimensional coverage:** All meaningful combinations of {function, authority, temporality, scope} are covered. The three groups naturally correspond to codification (writing rules), process (applying rules), and sovereign structure (the framework itself).

2. **Legal theory coverage:** The distinction between constitutive rules (that create the framework) and operative rules (that work within it) maps to Hart's primary/secondary rules. Due process and rights map to constitutional constraints on the legal system itself. Treaty enables inter-system legal relations.

3. **Layer boundary:** None of these require concepts from Layer 5 (Technology). Legal formalises governance but doesn't build — creating tools, artefacts, and processes is Layer 5's gap.

---

## Primitive Specifications

Full specifications in `docs/primitives.md` (Layer 4 section).

## Product Graph

Layer 4 maps to the **Governance Graph** — formalised rule-making with transparent adjudication. See `docs/product-layers.md`.

## Reference

- `docs/derivation-method.md` — The derivation method
- `docs/primitives.md` — Full specifications
- `docs/product-layers.md` — Governance Graph
- `docs/tests/primitives/04-community-governance.md` — Related integration test scenario
