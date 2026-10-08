# TRD — h5pkit Question Bank Registry and Documented Extensions

**Initiative:** `INIT-H5PKIT-001`

## Current State

- Package `h5p` (module root): `QuestionSet`, `Question{Library, Params,
  Extensions}`, `MultiChoiceQuestion` with typed `schemas.MultiChoiceParams`,
  a `QuestionSetBuilder` whose only question method is
  `AddMultipleChoiceQuestion`, JSON round-trip, `Validate`, and `.h5p`
  package create and load.
- `Extensions{H5PGo *H5PGoExtension}` serialized under the JSON key `h5pGo`,
  with `Section`, `Topic`, `Tags`, `Difficulty string`, `QuestionNumber`,
  `LearningObjective`, `Source`, and `Custom map[string]any`.
- `schemas/`: typed params and semantics for multiple choice; semantics JSON
  only for true/false and essay.
- `semantics/`: generic semantics field model.
- `ts/packages/editor`: React editor for multiple choice, true/false, and
  fill in the blanks.
- Published roadmap lists essay, true/false, and single choice set as
  short-term content types and a content marketplace as long-term.

## Design

### 1. Extension namespace and fields

Rename the block to the `h5pkit` key and give it a documented shape:

```go
type Extensions struct {
    H5PKit *Extension `json:"h5pkit,omitempty"`
    // Deprecated: read-only alias for content written before the rename.
    H5PGo  *Extension `json:"h5pGo,omitempty"`
}

type Extension struct {
    Section           string         `json:"section,omitempty"`
    Topic             string         `json:"topic,omitempty"`
    Tags              []string       `json:"tags,omitempty"`
    Difficulty        string         `json:"difficulty,omitempty"`        // display label
    DifficultyScore   *float64       `json:"difficultyScore,omitempty"`   // 0..1, relative within a bank
    Skills            []SkillRef     `json:"skills,omitempty"`
    Standards         []StandardRef  `json:"standards,omitempty"`
    QuestionNumber    int            `json:"questionNumber,omitempty"`
    LearningObjective string         `json:"learningObjective,omitempty"`
    Source            string         `json:"source,omitempty"`
    Custom            map[string]any `json:"custom,omitempty"`
}

type SkillRef    struct { Taxonomy, ID string }
type StandardRef struct { Framework, Code string }
```

- `UnmarshalJSON` on `Extensions` promotes `h5pGo` into `H5PKit` when
  `h5pkit` is absent; `MarshalJSON` never writes `h5pGo`. A
  `MigrateLegacyExtensions(*QuestionSet)` helper rewrites in place.
- Validation: `DifficultyScore` in `[0,1]`; each `SkillRef` has non-empty
  `Taxonomy` and `ID`; each `StandardRef` has non-empty `Framework` and
  `Code`. `Validate` on `QuestionSet` calls it.
- The block stays a sibling of `library` and `params` on the question (as
  today), so H5P players never see it inside library parameters.

### 2. `registry` package

`github.com/grokify/h5pkit/registry`, no dependency beyond the root package
and the standard library.

```go
type Bank struct {
    ID          string
    Slug        string
    Title       string
    Description string
    Subject     string
    Language    string          // BCP 47
    License     License
    Tags        []string
    Lineage     *Lineage        // nil for an original bank
    CreatedAt, UpdatedAt time.Time
}

type BankVersion struct {
    BankID      string
    Version     string          // semantic version
    Status      Status          // draft | published | archived
    Questions   []h5p.Question
    Skills      []h5p.SkillRef     // union of question skills, computed on publish
    Standards   []h5p.StandardRef  // union, computed on publish
    Notes       string
    PublishedAt *time.Time
}

type License struct {
    SPDX        string          // e.g. CC-BY-4.0, CC-BY-SA-4.0, proprietary
    Attribution string
    URL         string
}

type Lineage struct {
    ParentBankID  string
    ParentVersion string
    Relation      Relation      // fork | remix | translation
}
```

**Immutability.** `Publish` sets `Status = published`, computes the skill
and standards unions, and freezes the version; any mutation on a published
version returns `ErrImmutable`. New content means `NewVersion(bankID, from
version)` which copies questions into a draft.

**Store.**

```go
type Store interface {
    CreateBank(ctx, Bank) (Bank, error)
    GetBank(ctx, id) (Bank, error)
    ListBanks(ctx, ListFilter) ([]Bank, error)        // tag, skill, standard, license, subject, language, lineage
    NewVersion(ctx, bankID, fromVersion string) (BankVersion, error)
    SaveDraft(ctx, BankVersion) error                 // ErrImmutable if published
    Publish(ctx, bankID, version string) (BankVersion, error)
    GetVersion(ctx, bankID, version string) (BankVersion, error)
    ListVersions(ctx, bankID) ([]BankVersion, error)
    Fork(ctx, parentBankID, parentVersion string, into Bank, rel Relation) (Bank, BankVersion, error)
    Descendants(ctx, bankID) ([]Bank, error)
}
```

`MemoryStore` implements it. `registry/registrytest` exports
`RunStoreTests(t, func() Store)` so any implementation (an Ent store in a
consumer, for example) proves conformance.

### 3. Interop

- `registry.ExportQuestionSet(v BankVersion, opts) (*h5p.QuestionSet, error)`
  builds a playable set; `ExportPackage` wraps it in an `.h5p` with the
  required library dependencies.
- `registry.ImportPackage(path) (BankVersion, error)` loads an `.h5p`
  question set into a draft version, preserving extension blocks and
  promoting legacy `h5pGo` data.

### 4. Content types

Follow the existing pattern (`schemas.MultiChoiceParams` + `Validate`,
`NewMultiChoiceQuestion`, builder method), one typed params struct per
library, library version pinned in a constant:

| Library | Covers | Builder method |
|---|---|---|
| `H5P.TrueFalse` | true/false | `AddTrueFalseQuestion` |
| `H5P.Blanks` | fill in the blanks | `AddBlanksQuestion` |
| `H5P.SingleChoiceSet` | rapid single-choice sets | `AddSingleChoiceSet` |
| `H5P.Essay` | free text with keyword scoring | `AddEssayQuestion` |
| `H5P.DragText` | ordering and word matching | `AddDragTextQuestion` |
| `H5P.DragQuestion` | drag and drop, image matching | `AddDragQuestion` |

Semantics JSON for each is added to `schemas/` from the upstream library at
the pinned version, with a test that the typed params validate against it.
Library version strings are confirmed against the upstream releases at
implementation time, not assumed.

### 5. TypeScript editor

Add an "Metadata" section to the existing editor for difficulty score,
skills, and standards, reading and writing the `h5pkit` block. No new
question types in the editor in this initiative; the Go SDK leads.

### 6. Documentation

- `docs/guide/extensions.md`: the block, every field, the relative-difficulty
  rule, the rename and read-compat, and how another tool should treat the
  block.
- `docs/guide/question-banks.md`: registry concepts, immutability, lineage,
  conformance suite, and interop.
- `docs/reference/standards.md`: replace the placeholder with how
  `StandardRef` is used and the H5P specification points h5pkit follows.
- `docs/development/roadmap.md`: mark the content types delivered and link
  the registry.

### 7. Versioning

Release as v0.6.0. The namespace rename is noted as a behavior change with
read-compat; everything else is additive. The `go` directive is raised to
the organization's current minimum as part of the release item, subject to
CI toolchain availability.

## Sequencing Constraints

- §1 before §2 (the registry computes unions from `SkillRef` and
  `StandardRef`).
- §2 before §3 (interop operates on `BankVersion`).
- §4 is independent of §1–§3 and can run in parallel.
- §5 depends on §1.
- §6 and §7 last.

## Risks

- **Upstream semantics drift.** Pinning per library and testing typed params
  against the pinned semantics JSON catches it at upgrade time.
- **`Custom` as an escape hatch.** Document that `Custom` is for
  consumer-private data and is never interpreted by h5pkit or the registry;
  resist adding typed fields for one consumer's needs.
- **Over-modeling the registry.** Keep `Bank` to what every consumer needs;
  consumer-specific fields go in the consumer's persistence layer, keyed by
  `Bank.ID`.

## Testing Strategy

- Extension: JSON round-trip for new and legacy keys; validation table tests.
- Registry: `registrytest` suite run against `MemoryStore`; immutability,
  version copy, fork lineage, filter semantics.
- Interop: export then import yields an equal `BankVersion` (ignoring
  timestamps and status); exported `.h5p` passes the existing package
  validation.
- Content types: params validation against pinned semantics; builder output
  round-trips through JSON.
- Editor: component tests for the metadata section.
