# ROADMAP — h5pkit Question Bank Registry and Documented Extensions

**Initiative:** `INIT-H5PKIT-001`
**Repository:** `github.com/grokify/h5pkit`

Stay inside H5P where H5P has an answer; extend only where it does not, and
document each extension so other tools can use it. Phase 3 is independent
of Phases 1 and 2. Commits implementing an item carry
`Refs: RMI-H5PKIT-<NNN>`.

## Phase 1 — Extensions

**Theme:** A documented, namespaced metadata block

- [ ] `RMI-H5PKIT-001` Rename the extension namespace to `h5pkit` with read-compat for `h5pGo`, a migration helper, and a `MarshalJSON` that never writes the legacy key
- [ ] `RMI-H5PKIT-002` Add `difficultyScore`, `skills` (`SkillRef{Taxonomy, ID}`), and `standards` (`StandardRef{Framework, Code}`) with validation wired into `QuestionSet.Validate`
  - Depends on: `RMI-H5PKIT-001`
- [ ] `RMI-H5PKIT-003` Extensions guide: every field, the relative-difficulty rule, the rename, and guidance for other tools
  - Depends on: `RMI-H5PKIT-002`

## Phase 2 — Registry

**Theme:** Banks, immutable versions, lineage, and a store boundary

- [ ] `RMI-H5PKIT-004` `registry` types: `Bank`, `BankVersion`, `License`, `Lineage`, `Status`, `Relation`, `ListFilter`
  - Depends on: `RMI-H5PKIT-002`
- [ ] `RMI-H5PKIT-005` `Store` interface, `MemoryStore`, and the `registrytest` conformance suite
  - Depends on: `RMI-H5PKIT-004`
- [ ] `RMI-H5PKIT-006` Versioning rules: draft editing, publish freezes and computes skill and standards unions, `ErrImmutable`, new version from existing, fork with lineage and descendants lookup
  - Depends on: `RMI-H5PKIT-005`
- [ ] `RMI-H5PKIT-007` Interop: export a version to `QuestionSet` and `.h5p`; import an `.h5p` question set to a draft version preserving extensions
  - Depends on: `RMI-H5PKIT-006`

## Phase 3 — Content Types

**Theme:** Native H5P types instead of invented formats

- [ ] `RMI-H5PKIT-008` `H5P.TrueFalse` and `H5P.Blanks`: typed params, semantics, validation, builder methods
- [ ] `RMI-H5PKIT-009` `H5P.SingleChoiceSet` and `H5P.Essay`: typed params, semantics, validation, builder methods
  - Depends on: `RMI-H5PKIT-008`
- [ ] `RMI-H5PKIT-010` `H5P.DragText` and `H5P.DragQuestion` for ordering and matching: typed params, semantics, validation, builder methods
  - Depends on: `RMI-H5PKIT-008`

## Phase 4 — Editor, Documentation, Release

**Theme:** Make it usable and ship it

- [ ] `RMI-H5PKIT-011` TypeScript editor metadata section for difficulty score, skills, and standards on the `h5pkit` block
  - Depends on: `RMI-H5PKIT-002`
- [ ] `RMI-H5PKIT-012` Question-banks guide, standards reference, published roadmap update, Go directive bump, changelog, tag v0.6.0, record the release
  - Depends on: `RMI-H5PKIT-003`
  - Depends on: `RMI-H5PKIT-007`
  - Depends on: `RMI-H5PKIT-009`
  - Depends on: `RMI-H5PKIT-010`
  - Depends on: `RMI-H5PKIT-011`
