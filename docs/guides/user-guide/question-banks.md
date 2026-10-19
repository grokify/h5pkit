# Question Banks

The `registry` package manages **question banks**: named, licensed
collections of questions with immutable published versions and recorded
lineage. It depends only on the root `h5p` package and the standard library,
so any application can use it without a database dependency.

```go
import "github.com/grokify/h5pkit/registry"
```

## Concepts

| Type | Meaning |
|---|---|
| `Bank` | A named collection: title, subject, language, license, tags, and lineage. Its content lives in versions |
| `BankVersion` | One version of a bank's questions, with a `Status`, the skills and standards it covers, and notes |
| `License` | The terms the bank is offered under: an SPDX identifier, the attribution it requires, and a URL |
| `Lineage` | The published version a derived bank came from, and how it relates (`RelationFork`) |
| `Store` | Persists banks and versions. `MemoryStore` is the in-memory implementation |

A version is a `draft`, `published`, or `archived`. Only a draft can change;
`published` and `archived` are frozen.

## A bank's life

```go
ctx := context.Background()
store := registry.NewMemoryStore()

bank, err := store.CreateBank(ctx, registry.Bank{
    Title:    "Fractions",
    Subject:  "math",
    Language: "en",
    License:  registry.License{SPDX: "CC-BY-4.0", Attribution: "Example School"},
})

draft, err := store.NewVersion(ctx, bank.ID, "") // empty draft, numbered 1.0.0
draft.Questions = questions                      // []h5p.Question
err = store.SaveDraft(ctx, draft)

published, err := store.Publish(ctx, bank.ID, draft.Version)
// published.Skills and published.Standards are the questions' unions.
```

The compile-checked version is the package
[`Example`](https://pkg.go.dev/github.com/grokify/h5pkit/registry#example-package).

## Versioning rules

- **The store numbers versions.** A bank's first version is `1.0.0`; each
  later one is the next minor version after the highest existing one.
- **One draft at a time.** `NewVersion` on a bank that already has a draft
  returns `ErrDraftExists`. With a `fromVersion`, the draft starts as a copy
  of that version's questions and notes.
- **Publishing freezes.** `Publish` needs at least one question, validates
  the set and every question's `h5pkit` block, and computes the sorted,
  de-duplicated skills and standards. Saving or publishing a frozen version
  returns `ErrImmutable`.
- **Archived versions stay frozen.** The status exists for stores and
  applications that retire versions; the `Store` interface has no operation
  that sets it yet.

## Forks and lineage

```go
fork, forkDraft, err := store.Fork(ctx, bank.ID, "1.0.0",
    registry.Bank{Title: "Fractions, adapted"}, registry.RelationFork)
```

- A bank gets lineage only through `Fork`; `CreateBank` rejects a bank with
  `Lineage` set, so lineage always names real content.
- The parent must be a published or archived version. A draft returns
  `ErrNotPublished`.
- The fork starts with a `1.0.0` draft holding a copy of the parent's
  questions.
- A fork with no `License.SPDX` of its own inherits the parent's license.
- `Descendants` returns every bank derived from a bank, directly or through
  other derived banks.

## Finding banks

`ListBanks` takes a `ListFilter`. Every set field must match; a zero filter
matches every bank.

| Field | Matches |
|---|---|
| `Tag` | Banks carrying the tag |
| `Skill`, `Standard` | Banks whose latest published version covers it |
| `License` | Banks whose `License.SPDX` equals it |
| `Subject`, `Language` | Banks with that subject, or that language tag (ignoring case) |
| `ParentBankID` | Banks derived directly from that bank |

## Import and export

| Function | Does |
|---|---|
| `ExportQuestionSet(v, opts)` | Builds a playable `h5p.QuestionSet` from a version |
| `ExportPackage(path, v, opts)`, `WritePackage(w, v, opts)` | Write a version as an `.h5p` package |
| `ImportPackage(path)` | Reads an `.h5p` question set into a draft, preserving `h5pkit` blocks |

An exported package holds `h5p.json`, naming `H5P.QuestionSet 1.21` as the
main library and every question library as a dependency, and
`content/content.json`. It does not bundle library code: the H5P platform
that opens it supplies the libraries. `ImportPackage` returns a version with
an empty `BankID` and `Version`; set them, or save its questions into a draft.

## Writing your own store

A persistent store, for example over a SQL database, implements `Store` in
your own module and proves itself with the conformance suite, which defines
the rules above precisely:

```go
import "github.com/grokify/h5pkit/registry/registrytest"

func TestMyStore(t *testing.T) {
    registrytest.RunStoreTests(t, func() registry.Store { return newMyStore(t) })
}
```

## Errors

| Error | Returned when |
|---|---|
| `ErrNotFound` | A bank or version does not exist |
| `ErrExists` | A bank ID or slug is already taken |
| `ErrDraftExists` | A bank already has a draft |
| `ErrImmutable` | A published or archived version would change |
| `ErrNotPublished` | A fork names a draft version |
| `ErrInvalid` | Input fails validation |

Test them with `errors.Is`.
