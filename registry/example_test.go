package registry_test

import (
	"context"
	"fmt"
	"log"

	h5p "github.com/grokify/h5pkit"
	"github.com/grokify/h5pkit/registry"
)

// Example shows a bank's life: create it, draft and publish a version,
// then fork it.
func Example() {
	ctx := context.Background()
	store := registry.NewMemoryStore()

	bank, err := store.CreateBank(ctx, registry.Bank{
		Title:    "Fractions",
		Subject:  "math",
		Language: "en",
		License:  registry.License{SPDX: "CC-BY-4.0", Attribution: "Example School"},
	})
	if err != nil {
		log.Fatal(err)
	}

	qs, err := h5p.NewQuestionSetBuilder().
		SetTitle("Fractions").
		AddMultipleChoiceQuestion("What is 1/2 + 1/4?", []h5p.Answer{
			h5p.CreateAnswer("3/4", true),
			h5p.CreateAnswer("2/6", false),
		}).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	qs.Questions[0].Extensions = &h5p.Extensions{H5PKit: h5p.NewExtension("Adding", 1).
		WithDifficultyScore(0.3).
		WithSkills(h5p.SkillRef{Taxonomy: "example", ID: "fractions.add"})}

	draft, err := store.NewVersion(ctx, bank.ID, "")
	if err != nil {
		log.Fatal(err)
	}
	draft.Questions = qs.Questions
	if err := store.SaveDraft(ctx, draft); err != nil {
		log.Fatal(err)
	}
	published, err := store.Publish(ctx, bank.ID, draft.Version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(published.Version, published.Status, published.Skills[0].ID)

	fork, forkDraft, err := store.Fork(ctx, bank.ID, published.Version,
		registry.Bank{Title: "Fractions, adapted"}, registry.RelationFork)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fork.License.SPDX, fork.Lineage.ParentVersion, forkDraft.Status)
	// Output:
	// 1.0.0 published fractions.add
	// CC-BY-4.0 1.0.0 draft
}
