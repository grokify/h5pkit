# PRD — h5pkit Question Bank Registry and Documented Extensions

**Initiative:** `INIT-H5PKIT-001`
**Status:** Draft
**Home repo:** `github.com/grokify/h5pkit`

## Problem

h5pkit can build, validate, serialize, and package H5P question content, and
it already carries a small metadata extension (`h5pGo`) for section, topic,
tags, difficulty, and learning objective. Two things are missing for anyone
who wants to *manage* questions rather than only author them:

1. **There is no concept of a question bank.** A `QuestionSet` is a playable
   H5P activity, not a curated, versioned, reusable collection with a
   license and provenance. Consumers that need banks (a learning game, a
   course platform, an assessment tool) each have to invent their own
   wrapper, and those wrappers do not interoperate.
2. **The extension is informal.** Its namespace reflects a former project
   name, `difficulty` is a free string, there is no place for skill or
   standards alignment, and nothing documents the block as something other
   tools may read or write. Because H5P players ignore unknown data, the
   block is safe to carry, but it is only useful across tools if its shape is
   stable and documented.

The preference is to stay inside H5P's own approach wherever H5P has one, and
to extend only where it does not, documenting each extension so that other
projects can use it independently of any one consumer.

## Goals

1. **A storage-agnostic question bank registry** in h5pkit: banks, immutable
   published versions, licenses, fork and remix lineage, and skill and
   standards tags, with a `Store` interface and an in-memory implementation.
   No database dependency in h5pkit.
2. **Documented, namespaced extensions** for the metadata H5P does not model:
   relative difficulty, skills, and standards alignment, written so a
   standard H5P player still imports the content unchanged.
3. **More native question types** in the Go SDK: true/false, fill in the
   blanks, single choice set, essay, and the drag-based types that cover
   matching and ordering, so consumers do not invent non-H5P question
   formats.
4. **Round-trip interoperability**: a bank version exports to a standard
   `.h5p` question set, and an imported `.h5p` becomes a bank version.

## Non-Goals

- A persistent (SQL or Ent) store. That belongs to the consuming application
  or a separate module, implementing the registry's `Store` interface.
- Mastery, scoring policy, rewards, certificates, or any learner state. The
  registry describes content, not people.
- A hosted registry service or marketplace. This initiative delivers the
  library; services are built on it.
- Replacing H5P's per-answer feedback with a custom explanation field. H5P
  already models feedback natively.

## Users and Stakeholders

| Who | Needs |
|---|---|
| Application developers building on h5pkit | Bank, version, license, and lineage types they can persist and expose without redesigning them |
| Content authors and educators (through those applications) | Difficulty, skill, and standards metadata that survives export and import |
| Other H5P tooling | A documented extension block they can read or ignore safely |
| h5pkit maintainers | Content types that follow the existing typed-params pattern and the roadmap already published |

## Requirements

### Functional

- **F1.** A `registry` package defines `Bank`, `BankVersion`, `License`,
  `Lineage` (parent bank and version for forks and remixes), and tag types
  for skills and standards. Versions hold `h5p.Question` values.
- **F2.** A published `BankVersion` is immutable; changing content creates a
  new version. Draft versions may be edited until published.
- **F3.** A `Store` interface covers create, get, list (with filters by tag,
  skill, standard, license), new version, publish, and lineage lookup; an
  in-memory implementation and a reusable conformance test suite ship with
  it so other stores can prove compatibility.
- **F4.** The metadata extension is renamed to the `h5pkit` namespace. The
  former `h5pGo` key is no longer read or written; no stored content depends
  on it.
- **F5.** The extension gains `difficultyScore` (a relative score within the
  bank, 0 to 1), `skills` (identifiers in a named taxonomy), and
  `standards` (framework plus code, for example a Common Core or NGSS
  reference). The existing free-string `difficulty` label remains for
  display.
- **F6.** Extension content is validated: score range, non-empty skill
  identifiers, standards with both framework and code.
- **F7.** Typed parameters and builder methods for `H5P.TrueFalse`,
  `H5P.Blanks`, `H5P.SingleChoiceSet`, `H5P.Essay`, `H5P.DragText`, and
  `H5P.DragQuestion`, each validated against its semantics.
- **F8.** A `BankVersion` exports to a `QuestionSet` and to a `.h5p`
  package; a `.h5p` question set imports to a draft `BankVersion`, with
  extension metadata preserved both ways.
- **F9.** The TypeScript editor exposes the extension fields (difficulty
  score, skills, standards) for the types it already edits.

### Non-Functional

- **N1.** Content with the extension block loads in a standard H5P player
  unchanged; the extension lives beside, not inside, library parameters.
- **N2.** Every extension field is documented in a guide that stands on its
  own, without reference to any particular consumer.
- **N3.** Existing public API is preserved except for the namespace rename
  in F4, which is called out in release notes.
- **N4.** `golangci-lint` passes; unit tests cover the registry, the
  extension validation, each new content type, and round-trip export and
  import.

## Success Criteria

- A consumer can create a bank, add questions with difficulty, skill, and
  standards metadata, publish a version, fork it into a new bank with
  lineage recorded, and export the version as a `.h5p` that a standard H5P
  player opens.
- A second `Store` implementation (in a consumer) passes the conformance
  suite without changes to h5pkit.
- The six content types in F7 round-trip through JSON and `.h5p` with
  validation.

## Risks

| Risk | Mitigation |
|---|---|
| Namespace rename breaks stored content | No content depends on the old key; the release notes name the rename |
| Extension drifts into consumer-specific fields | Keep the block to content metadata; anything about learners or rewards is refused in review |
| Content-type semantics versions move upstream | Pin the library version string per type, as `H5P.MultiChoice 1.16` is pinned today, and document the policy |
| Registry grows into a service | Non-goal stated; the `Store` interface is the boundary |

## Decisions

- **Difficulty is relative within a bank, not absolute.** A score of 0.8 in
  an arithmetic bank is not comparable to 0.8 in a calculus bank; the
  documentation says so, and no cross-bank scale is implied.
- **Standards are references, not a taxonomy.** h5pkit records
  `{framework, code}` pairs; it does not ship or validate the frameworks
  themselves.
- **Skills use caller-supplied taxonomy identifiers.** h5pkit records
  `{taxonomy, id}`; the taxonomy is the consumer's.
