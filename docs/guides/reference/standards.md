# Standards Compliance

This page covers two kinds of standards: the **learning standards** a
question can be aligned to, and the **H5P specification** points h5pkit
follows so its content works in any H5P player.

## Learning standards

A question's `h5pkit` block can name the learning standards it addresses,
each as a `StandardRef`:

```go
ext := h5p.NewExtension("Fractions", 1).
    WithStandards(h5p.StandardRef{Framework: "CCSS.MATH", Code: "3.NF.A.1"})
```

| Field | Meaning |
|---|---|
| `Framework` | The standards framework, such as `CCSS.MATH` or `NGSS` |
| `Code` | The standard's code within that framework |

h5pkit does not ship a catalog of frameworks or codes and does not check
that a code exists: it stores and compares the pair. Validation only
requires both fields to be non-empty. Use the framework's own published
identifiers, so banks from different authors align on the same values.

When a bank version is published, the registry computes the sorted,
de-duplicated union of its questions' standards, and `ListFilter.Standard`
finds banks whose latest published version covers one. Skills
(`SkillRef{Taxonomy, ID}`) work the same way. See
[Extensions](../user-guide/extensions.md) and
[Question Banks](../user-guide/question-banks.md).

## H5P specification points

h5pkit keeps its content readable by standard H5P tools:

- **Library versions are pinned.** A question names its library with a
  version, such as `H5P.MultiChoice 1.16`, and a question set exports as
  `H5P.QuestionSet 1.21`.
- **Metadata stays out of `params`.** The `h5pkit` block is a sibling of
  `library` and `params`, which H5P players ignore, so it never reaches a
  player's parameters.
- **Packages use the standard layout.** An exported `.h5p` holds `h5p.json`,
  with the main library and every question library as preloaded
  dependencies, and `content/content.json`.
- **Library code is not bundled.** The H5P platform that opens a package
  supplies the libraries.

See [H5P Specification](h5p-spec.md) and
[File Structure](file-structure.md) for the details.
