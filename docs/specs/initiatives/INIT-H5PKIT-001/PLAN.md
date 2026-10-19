# PLAN — h5pkit Question Bank Registry and Documented Extensions

**Initiative:** `INIT-H5PKIT-001`

## Entry Gate

- Specs reviewed; the namespace rename (`h5pGo` to `h5pkit`) agreed.
- Upstream library versions for the six new content types confirmed against
  their current releases.

## Build Order

### Phase 1 — Extensions

Rename the namespace, add `DifficultyScore`, `SkillRef`,
and `StandardRef`, validate them, and write the extensions guide. Exit:
content round-trips under the new key; validation
rejects malformed metadata; the guide documents every field.

### Phase 2 — Registry

Define `Bank`, `BankVersion`, `License`, `Lineage`, and `Store`; implement
`MemoryStore`; ship the `registrytest` conformance suite; add export and
import to and from `QuestionSet` and `.h5p`. Exit: the suite passes against
`MemoryStore`; a published version is immutable; a fork records lineage; an
exported package opens in a standard H5P player.

### Phase 3 — Documentation and First Release

Write the question-banks guide and the standards reference; update the
published roadmap; raise the Go directive; changelog and tag v0.6.0 with the
extensions and the registry; record the release. Exit: the registry and
extensions are usable from a tagged module.

### Phase 4 — Content Types

Add typed params, validation, and builder methods for true/false, blanks,
single choice set, essay, drag text, and drag question. This phase is
independent of Phases 1 to 3 and may run in parallel. Exit: each type
validates against its pinned semantics and round-trips through JSON and
`.h5p`.

### Phase 5 — Editor and Next Release

Expose the metadata fields in the TypeScript editor; document the new
content types; changelog and tag the next minor release; record the
release. Exit: PRD success criteria hold.

## Risks and Mitigations

| Risk | Mitigation |
|---|---|
| Rename surprises an external user | No stored content depends on the old key; the release notes name the rename |
| Content-type work expands to the whole H5P catalog | The six types are the list; others are new roadmap items |
| Registry absorbs consumer concerns | `Store` interface is the boundary; consumer fields live in the consumer |

## Follow-On

- A persistent `Store` implementation in a consuming application, validated
  with `registrytest`.
- Further content types from the published roadmap (image hotspots,
  interactive video) as demand appears.
- A CLI for bank import, export, and validation, built on the registry.
