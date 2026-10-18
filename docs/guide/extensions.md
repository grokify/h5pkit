# Extensions

H5P defines how a question is played, not how it is organized. h5pkit adds a
small, documented metadata block to each question for organization,
difficulty, skills, and standards. Any tool can read and write it.

## Where the block lives

The block sits under `extensions.h5pkit`, a sibling of a question's `library`
and `params`:

```json
{
  "library": "H5P.MultiChoice 1.16",
  "params": { "question": "What is 3/4 + 1/8?", "answers": [] },
  "extensions": {
    "h5pkit": {
      "section": "Fractions",
      "topic": "Adding unlike denominators",
      "tags": ["fractions", "addition"],
      "difficulty": "medium",
      "difficultyScore": 0.6,
      "skills": [{ "taxonomy": "example-skills", "id": "fractions.add.unlike" }],
      "standards": [{ "framework": "example-framework", "code": "5.NF.A.1" }],
      "questionNumber": 4,
      "learningObjective": "Add fractions with unlike denominators",
      "source": "Unit 3"
    }
  }
}
```

Because the block is outside `params`, an H5P player never sees it. Players
ignore unknown fields, so content with the block plays everywhere content
without it does.

## Fields

| Field | Type | Meaning |
|-------|------|---------|
| `section` | string | The section or category the question belongs to |
| `topic` | string | A more specific topic within the section |
| `tags` | string array | Free-form labels |
| `difficulty` | string | A display label, such as `easy`, `medium`, or `hard` |
| `difficultyScore` | number, 0 to 1 | Difficulty relative to the other questions in the same bank |
| `skills` | array of `{taxonomy, id}` | Skills the question exercises |
| `standards` | array of `{framework, code}` | Standards the question is aligned to |
| `questionNumber` | integer | Explicit ordering within a question set |
| `learningObjective` | string | What the question tests |
| `source` | string | Where the content came from |
| `custom` | object | Consumer-private data, never interpreted by h5pkit |

Every field is optional.

### The relative-difficulty rule

`difficultyScore` ranks a question against the **other questions in the same
bank**: 0 is the easiest question in the bank and 1 the hardest. It is not an
absolute scale, and two banks' scores are not comparable. A score of 0.8 in a
beginner bank and 0.8 in an advanced bank say nothing about each other.

Use `difficulty` for a label a person reads, and `difficultyScore` for
ordering and selection within a bank. A question without a score is
unscored, not easy.

### Skills and standards

A skill or standard is always a pair: the scheme it comes from and its
identifier in that scheme. h5pkit does not ship or validate any taxonomy or
framework; it only requires both parts to be present. Pick stable
identifiers, because banks compute and filter on these pairs.

### Custom

`custom` is an escape hatch for data that only one consumer understands.
h5pkit and its registry carry it unchanged and never read it. If a field in
`custom` turns out to be generally useful, propose it as a documented field
instead.

## Validation

`Extension.Validate` checks that `difficultyScore`, when present, is between
0 and 1, and that every skill has a `taxonomy` and an `id` and every standard
a `framework` and a `code`. `QuestionSet.Validate` runs it for every
question.

## Building the block in Go

```go
ext := h5p.NewExtension("Fractions", 4).
    WithTopic("Adding unlike denominators").
    WithTags("fractions", "addition").
    WithDifficulty("medium").
    WithDifficultyScore(0.6).
    WithSkills(h5p.SkillRef{Taxonomy: "example-skills", ID: "fractions.add.unlike"}).
    WithStandards(h5p.StandardRef{Framework: "example-framework", Code: "5.NF.A.1"}).
    WithLearningObjective("Add fractions with unlike denominators")

question := h5p.NewMultiChoiceQuestionWithExtensions(params, ext)

if err := ext.Validate(); err != nil {
    return err
}
```

`Extensions.Extension()` returns the block from a question, or nil when it
has none.

## The `h5pkit` key

The block is read and written only under `extensions.h5pkit`. Earlier
releases used `extensions.h5pGo`; that key is no longer read, so rename it in
any content that still carries it.

## Guidance for other tools

A tool that reads or writes H5P content should:

- Preserve `extensions` when it rewrites a question, even if it does not
  understand the block.
- Read and write the block only under `h5pkit`.
- Treat a missing field as unknown, not as zero or empty.
- Never move the block into `params`, where it would reach the player.
- Leave `custom` untouched.
