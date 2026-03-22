# Layer 5: Technology

## Derivation

### The Gap

Layer 4 formalises rules and adjudication, but it governs — it doesn't build. There's no concept of creating tools, artefacts, or systems. An actor can follow rules but cannot create something new that extends the system's capabilities.

**Test:** Can you express "Alice built a tool that automates the review process, version 2 improves on version 1, and it should be deprecated when version 3 ships" in Layer 4? You can enact laws about tools and adjudicate disputes about them. But "built" (creation), "improves on" (iteration), "automates" (process transformation), and "deprecated" (artefact lifecycle) have no Layer 4 representation. Governing is not building.

### The Transition

**Governing → Building**

An actor that can only follow rules becomes one that can create new things. The fundamental new capacity: investigating the world, making artefacts, defining standards, and building systems.

### Base Operations

What can a builder DO that a rule-follower cannot?

1. **Investigate** — apply method to discover knowledge
2. **Create** — make a new artefact or tool
3. **Standardise** — establish repeatable, reliable processes
4. **Automate** — make processes self-executing

### Semantic Dimensions

| Dimension | Values | What it distinguishes |
|-----------|--------|-----------------------|
| **Object** | Artefact (thing made) / Process (way of making) | Is this about what's built or how it's built? |
| **Phase** | Genesis (creation) / Operation (use) / Retirement (end-of-life) | Where in the lifecycle? |
| **Direction** | Forward (making new) / Backward (assessing existing) | Building or evaluating? |
| **Agency** | Manual (human-driven) / Mechanical (automated) | Who does the work? |

### Decomposition

**Group 0 — Investigation** (discovering how things work)

| Primitive | Object | Phase | Direction | Agency | What it does |
|-----------|--------|-------|-----------|--------|--------------|
| **Method** | Process | Genesis | Forward | Manual | A systematic approach to investigation |
| **Measurement** | Process | Operation | Backward | Both | Quantifying properties of the world |
| **Knowledge** | Artefact | Operation | Backward | Both | Verified understanding of how things work |
| **Model** | Artefact | Genesis | Forward | Both | A simplified representation that enables prediction |

Investigation lifecycle: methods are established (Method) → measurements are taken (Measurement) → understanding is verified (Knowledge) → and simplified representations enable prediction (Model). This is how the world becomes tractable.

**Group 1 — Creation** (making things)

| Primitive | Object | Phase | Direction | Agency | What it does |
|-----------|--------|-------|-----------|--------|--------------|
| **Tool** | Artefact | Operation | Forward | Both | Artefacts that extend capabilities |
| **Technique** | Process | Operation | Forward | Manual | Skilled methods for using tools |
| **Invention** | Artefact | Genesis | Forward | Both | Creating something genuinely new |
| **Abstraction** | Artefact | Genesis | Forward | Both | Generalising from specifics to reusable patterns |

Creation lifecycle: tools are made (Tool) → techniques for using them develop (Technique) → occasionally something genuinely new emerges (Invention) → and patterns are generalised for reuse (Abstraction). Abstraction is key to scaling — concrete solutions become general principles.

**Group 2 — Systems** (making things work together)

| Primitive | Object | Phase | Direction | Agency | What it does |
|-----------|--------|-------|-----------|--------|--------------|
| **Infrastructure** | Artefact | Operation | Forward | Both | The substrate that other things run on |
| **Standard** | Process | Operation | Forward | Both | Agreed-upon conventions that enable interoperability |
| **Efficiency** | Process | Operation | Backward | Both | Getting more output from less input |
| **Automation** | Process | Operation | Forward | Mechanical | Converting manual processes to self-executing ones |

Systems lifecycle: infrastructure is established (Infrastructure) → standards enable interoperability (Standard) → efficiency is optimised (Efficiency) → and processes become self-executing (Automation). Automation is key to SELF-EVOLVE — identifying patterns that can migrate from intelligent to mechanical.

### Gap Analysis

| Behavior | Maps to | Notes |
|----------|---------|-------|
| "We use the scientific method for testing" | Method | |
| "Response time is 200ms at p99" | Measurement | |
| "We know how the caching system works" | Knowledge | |
| "This model predicts load patterns" | Model | |
| "This library provides authentication" | Tool | |
| "Here's how to use the deployment pipeline" | Technique | |
| "This is a novel approach to caching" | Invention | |
| "The pattern generalises to all queue processors" | Abstraction | |
| "The CI pipeline runs on Kubernetes" | Infrastructure | |
| "All APIs use REST with JSON" | Standard | |
| "We reduced build time by 40%" | Efficiency | |
| "Auto-format on save" | Automation | |
| Continuous integration pipeline | Infrastructure + Standard + Automation | Composition |
| Technical debt tracking | Measurement + Efficiency + Knowledge | Composition |
| Open source contribution | Invention + Tool + L2.Offer | Cross-layer composition |

**No gaps found.**

### Completeness Argument

1. **Dimensional coverage:** The {object, phase, direction, agency} space is covered. The three groups correspond to understanding (Investigation), making (Creation), and scaling (Systems).

2. **Engineering coverage:** The build-measure-learn cycle from lean methodology maps directly: Method/Invention (build) → Measurement/Knowledge (measure) → Abstraction/Efficiency (learn). Automation captures the mechanical-to-intelligent transition that's central to SELF-EVOLVE.

3. **Layer boundary:** None of these require concepts from Layer 6 (Information). Technology builds concrete artefacts and processes — symbolic representation, language, and meaning independent of physical events are Layer 6's gap.

---

## Primitive Specifications

Full specifications in `docs/primitives.md` (Layer 5 section).

## Product Graph

Layer 5 maps to the **Build Graph** — development and CI/CD with provenance. See `docs/product-layers.md`.

## Reference

- `docs/derivation-method.md` — The derivation method
- `docs/primitives.md` — Full specifications
- `docs/product-layers.md` — Build Graph
- `docs/tests/primitives/01-agent-audit-trail.md` — Related integration test scenario
