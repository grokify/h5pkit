// Package registrytest is the conformance suite for registry.Store
// implementations. A store passes when RunStoreTests passes against it:
//
//	func TestMyStore(t *testing.T) {
//		registrytest.RunStoreTests(t, func() registry.Store { return newMyStore(t) })
//	}
//
// newStore must return an empty store each time it is called.
package registrytest

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	h5p "github.com/grokify/h5pkit"
	"github.com/grokify/h5pkit/registry"
)

// RunStoreTests runs every conformance test against stores made by newStore.
func RunStoreTests(t *testing.T, newStore func() registry.Store) {
	t.Helper()
	tests := []struct {
		name string
		fn   func(*testing.T, registry.Store)
	}{
		{"CreateBank", testCreateBank},
		{"CreateBankRejects", testCreateBankRejects},
		{"GetBankNotFound", testGetBankNotFound},
		{"ReturnsCopies", testReturnsCopies},
		{"ListBanksFilters", testListBanksFilters},
		{"NewVersion", testNewVersion},
		{"OneDraftAtATime", testOneDraftAtATime},
		{"SaveDraft", testSaveDraft},
		{"SaveDraftRejects", testSaveDraftRejects},
		{"ListVersions", testListVersions},
		{"CanceledContext", testCanceledContext},
		{"Publish", testPublish},
		{"PublishRejects", testPublishRejects},
		{"PublishedIsImmutable", testPublishedIsImmutable},
		{"NewVersionFromPublished", testNewVersionFromPublished},
		{"ListBanksBySkillAndStandard", testListBanksBySkillAndStandard},
		{"Fork", testFork},
		{"ForkRejects", testForkRejects},
		{"Descendants", testDescendants},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn(t, newStore())
		})
	}
}

// Question returns a multiple choice question carrying an h5pkit block with
// the given skills and standards, for use in tests.
func Question(text string, skills []h5p.SkillRef, standards []h5p.StandardRef) h5p.Question {
	return h5p.Question{
		Library: "H5P.MultiChoice 1.16",
		Params: map[string]any{
			"question": text,
			"answers": []any{
				map[string]any{"text": "Yes", "correct": true},
				map[string]any{"text": "No", "correct": false},
			},
		},
		Extensions: &h5p.Extensions{H5PKit: &h5p.Extension{
			Section:   "Section",
			Skills:    skills,
			Standards: standards,
		}},
	}
}

// SameJSON reports whether a and b encode to the same JSON, which is how the
// suite compares question content regardless of its Go representation.
func SameJSON(t *testing.T, a, b any) bool {
	t.Helper()
	ja, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	jb, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	return string(ja) == string(jb)
}

func mustCreate(t *testing.T, s registry.Store, b registry.Bank) registry.Bank {
	t.Helper()
	got, err := s.CreateBank(context.Background(), b)
	if err != nil {
		t.Fatalf("CreateBank(%q): %v", b.Title, err)
	}
	return got
}

func mustNewDraft(t *testing.T, s registry.Store, bankID string) registry.BankVersion {
	t.Helper()
	v, err := s.NewVersion(context.Background(), bankID, "")
	if err != nil {
		t.Fatalf("NewVersion(%q): %v", bankID, err)
	}
	return v
}

func mustSave(t *testing.T, s registry.Store, v registry.BankVersion) {
	t.Helper()
	if err := s.SaveDraft(context.Background(), v); err != nil {
		t.Fatalf("SaveDraft(%s %s): %v", v.BankID, v.Version, err)
	}
}

func wantErr(t *testing.T, what string, err, want error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Errorf("%s: got %v, want %v", what, err, want)
	}
}

func testCreateBank(t *testing.T, s registry.Store) {
	ctx := context.Background()
	in := registry.Bank{
		Slug:     "fractions",
		Title:    "Fractions",
		Subject:  "math",
		Language: "en",
		License:  registry.License{SPDX: "CC-BY-4.0", Attribution: "Example Author"},
		Tags:     []string{"grade-4"},
	}
	b := mustCreate(t, s, in)
	if b.ID == "" {
		t.Fatal("store did not assign an ID")
	}
	if b.CreatedAt.IsZero() || b.UpdatedAt.IsZero() {
		t.Error("store did not set timestamps")
	}
	got, err := s.GetBank(ctx, b.ID)
	if err != nil {
		t.Fatalf("GetBank: %v", err)
	}
	if got.Title != in.Title || got.Slug != in.Slug || got.License != in.License || !SameJSON(t, got.Tags, in.Tags) {
		t.Errorf("GetBank = %+v, want fields of %+v", got, in)
	}

	withID := mustCreate(t, s, registry.Bank{ID: "chosen-id", Title: "Chosen"})
	if withID.ID != "chosen-id" {
		t.Errorf("caller-chosen ID replaced: %q", withID.ID)
	}
}

func testCreateBankRejects(t *testing.T, s registry.Store) {
	ctx := context.Background()
	_, err := s.CreateBank(ctx, registry.Bank{Title: " "})
	wantErr(t, "untitled bank", err, registry.ErrInvalid)

	_, err = s.CreateBank(ctx, registry.Bank{Title: "Derived", Lineage: &registry.Lineage{
		ParentBankID: "p", ParentVersion: "1.0.0", Relation: registry.RelationFork,
	}})
	wantErr(t, "bank with lineage", err, registry.ErrInvalid)

	mustCreate(t, s, registry.Bank{ID: "dup", Slug: "dup-slug", Title: "First"})
	_, err = s.CreateBank(ctx, registry.Bank{ID: "dup", Title: "Second"})
	wantErr(t, "duplicate ID", err, registry.ErrExists)
	_, err = s.CreateBank(ctx, registry.Bank{Slug: "dup-slug", Title: "Third"})
	wantErr(t, "duplicate slug", err, registry.ErrExists)
}

func testGetBankNotFound(t *testing.T, s registry.Store) {
	_, err := s.GetBank(context.Background(), "missing")
	wantErr(t, "GetBank", err, registry.ErrNotFound)
}

func testReturnsCopies(t *testing.T, s registry.Store) {
	ctx := context.Background()
	b := mustCreate(t, s, registry.Bank{Title: "Copies", Tags: []string{"a"}})
	b.Tags[0] = "changed"
	got, err := s.GetBank(ctx, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tags[0] != "a" {
		t.Error("changing a returned bank changed the store")
	}

	v := mustNewDraft(t, s, b.ID)
	v.Questions = []h5p.Question{Question("Q1", nil, nil)}
	mustSave(t, s, v)
	got1, err := s.GetVersion(ctx, b.ID, v.Version)
	if err != nil {
		t.Fatal(err)
	}
	got1.Questions[0].Library = "changed"
	got2, err := s.GetVersion(ctx, b.ID, v.Version)
	if err != nil {
		t.Fatal(err)
	}
	if got2.Questions[0].Library != "H5P.MultiChoice 1.16" {
		t.Error("changing a returned version changed the store")
	}
}

func testListBanksFilters(t *testing.T, s registry.Store) {
	ctx := context.Background()
	a := mustCreate(t, s, registry.Bank{Title: "A", Subject: "math", Language: "en", License: registry.License{SPDX: "CC-BY-4.0"}, Tags: []string{"x"}})
	b := mustCreate(t, s, registry.Bank{Title: "B", Subject: "science", Language: "pt-BR", License: registry.License{SPDX: "CC-BY-SA-4.0"}, Tags: []string{"x", "y"}})
	c := mustCreate(t, s, registry.Bank{Title: "C", Subject: "math", Language: "EN", License: registry.License{SPDX: "proprietary"}})

	tests := []struct {
		name string
		f    registry.ListFilter
		want []string
	}{
		{"all", registry.ListFilter{}, []string{a.ID, b.ID, c.ID}},
		{"tag", registry.ListFilter{Tag: "x"}, []string{a.ID, b.ID}},
		{"tag y", registry.ListFilter{Tag: "y"}, []string{b.ID}},
		{"subject", registry.ListFilter{Subject: "math"}, []string{a.ID, c.ID}},
		{"license", registry.ListFilter{License: "CC-BY-SA-4.0"}, []string{b.ID}},
		{"language ignores case", registry.ListFilter{Language: "en"}, []string{a.ID, c.ID}},
		{"combined", registry.ListFilter{Subject: "math", Tag: "x"}, []string{a.ID}},
		{"none", registry.ListFilter{Tag: "absent"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListBanks(ctx, tt.f)
			if err != nil {
				t.Fatalf("ListBanks: %v", err)
			}
			ids := make([]string, 0, len(got))
			for _, g := range got {
				ids = append(ids, g.ID)
			}
			if !sameSet(ids, tt.want) {
				t.Errorf("got %v, want %v", ids, tt.want)
			}
		})
	}
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[string]int, len(a))
	for _, x := range a {
		seen[x]++
	}
	for _, x := range b {
		seen[x]--
		if seen[x] < 0 {
			return false
		}
	}
	return true
}

func testNewVersion(t *testing.T, s registry.Store) {
	ctx := context.Background()
	_, err := s.NewVersion(ctx, "missing", "")
	wantErr(t, "NewVersion on a missing bank", err, registry.ErrNotFound)

	b := mustCreate(t, s, registry.Bank{Title: "Versions"})
	v := mustNewDraft(t, s, b.ID)
	if v.Version != "1.0.0" || v.Status != registry.StatusDraft || v.BankID != b.ID || len(v.Questions) != 0 {
		t.Errorf("first version = %+v, want an empty 1.0.0 draft", v)
	}
	_, err = s.NewVersion(ctx, b.ID, "9.9.9")
	wantErr(t, "NewVersion from a missing version", err, registry.ErrNotFound)
}

func testOneDraftAtATime(t *testing.T, s registry.Store) {
	b := mustCreate(t, s, registry.Bank{Title: "Drafts"})
	mustNewDraft(t, s, b.ID)
	_, err := s.NewVersion(context.Background(), b.ID, "")
	wantErr(t, "second draft", err, registry.ErrDraftExists)
}

func testSaveDraft(t *testing.T, s registry.Store) {
	ctx := context.Background()
	b := mustCreate(t, s, registry.Bank{Title: "Save"})
	v := mustNewDraft(t, s, b.ID)
	v.Questions = []h5p.Question{Question("Q1", nil, nil), Question("Q2", nil, nil)}
	v.Notes = "first pass"
	mustSave(t, s, v)

	got, err := s.GetVersion(ctx, b.ID, v.Version)
	if err != nil {
		t.Fatalf("GetVersion: %v", err)
	}
	if !SameJSON(t, got.Questions, v.Questions) || got.Notes != "first pass" || got.Status != registry.StatusDraft {
		t.Errorf("saved draft = %+v", got)
	}

	v.Questions = v.Questions[:1]
	mustSave(t, s, v)
	got, err = s.GetVersion(ctx, b.ID, v.Version)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Questions) != 1 {
		t.Errorf("draft not replaced: %d questions", len(got.Questions))
	}
}

func testSaveDraftRejects(t *testing.T, s registry.Store) {
	ctx := context.Background()
	err := s.SaveDraft(ctx, registry.BankVersion{BankID: "missing", Version: "1.0.0"})
	wantErr(t, "missing bank", err, registry.ErrNotFound)

	b := mustCreate(t, s, registry.Bank{Title: "Rejects"})
	err = s.SaveDraft(ctx, registry.BankVersion{BankID: b.ID, Version: "1.0.0"})
	wantErr(t, "missing version", err, registry.ErrNotFound)

	v := mustNewDraft(t, s, b.ID)
	v.Status = registry.StatusPublished
	err = s.SaveDraft(ctx, v)
	wantErr(t, "status change through SaveDraft", err, registry.ErrInvalid)
}

func testListVersions(t *testing.T, s registry.Store) {
	ctx := context.Background()
	_, err := s.ListVersions(ctx, "missing")
	wantErr(t, "missing bank", err, registry.ErrNotFound)

	b := mustCreate(t, s, registry.Bank{Title: "List"})
	vs, err := s.ListVersions(ctx, b.ID)
	if err != nil || len(vs) != 0 {
		t.Fatalf("new bank versions = %v, %v", vs, err)
	}
	mustNewDraft(t, s, b.ID)
	vs, err = s.ListVersions(ctx, b.ID)
	if err != nil || len(vs) != 1 || vs[0].Version != "1.0.0" {
		t.Fatalf("versions = %+v, %v", vs, err)
	}
}

func testCanceledContext(t *testing.T, s registry.Store) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := s.CreateBank(ctx, registry.Bank{Title: "Canceled"})
	wantErr(t, "CreateBank", err, context.Canceled)
	_, err = s.ListBanks(ctx, registry.ListFilter{})
	wantErr(t, "ListBanks", err, context.Canceled)
	_, err = s.GetVersion(ctx, "b", "1.0.0")
	wantErr(t, "GetVersion", err, context.Canceled)
}

var (
	skillA = h5p.SkillRef{Taxonomy: "tax", ID: "a"}
	skillB = h5p.SkillRef{Taxonomy: "tax", ID: "b"}
	stdA   = h5p.StandardRef{Framework: "fw", Code: "1.A"}
	stdB   = h5p.StandardRef{Framework: "fw", Code: "2.B"}
)

// mustPublished creates a bank with one published version holding the
// given questions and returns the bank and that version.
func mustPublished(t *testing.T, s registry.Store, b registry.Bank, qs ...h5p.Question) (registry.Bank, registry.BankVersion) {
	t.Helper()
	b = mustCreate(t, s, b)
	v := mustNewDraft(t, s, b.ID)
	v.Questions = qs
	mustSave(t, s, v)
	pub, err := s.Publish(context.Background(), b.ID, v.Version)
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}
	return b, pub
}

func testPublish(t *testing.T, s registry.Store) {
	ctx := context.Background()
	b, pub := mustPublished(t, s, registry.Bank{Title: "Pub"},
		Question("Q1", []h5p.SkillRef{skillB, skillA}, []h5p.StandardRef{stdB}),
		Question("Q2", []h5p.SkillRef{skillA}, []h5p.StandardRef{stdA, stdB}),
		h5p.Question{Library: "H5P.MultiChoice 1.16", Params: map[string]any{"question": "no block"}},
	)
	if pub.Status != registry.StatusPublished || pub.PublishedAt == nil || pub.PublishedAt.IsZero() {
		t.Errorf("published version = %+v", pub)
	}
	if !SameJSON(t, pub.Skills, []h5p.SkillRef{skillA, skillB}) {
		t.Errorf("skills union = %v", pub.Skills)
	}
	if !SameJSON(t, pub.Standards, []h5p.StandardRef{stdA, stdB}) {
		t.Errorf("standards union = %v", pub.Standards)
	}
	got, err := s.GetVersion(ctx, b.ID, pub.Version)
	if err != nil {
		t.Fatal(err)
	}
	if !SameJSON(t, got, pub) {
		t.Errorf("GetVersion after publish = %+v, want %+v", got, pub)
	}
}

func testPublishRejects(t *testing.T, s registry.Store) {
	ctx := context.Background()
	_, err := s.Publish(ctx, "missing", "1.0.0")
	wantErr(t, "missing bank", err, registry.ErrNotFound)

	b := mustCreate(t, s, registry.Bank{Title: "Rejects"})
	v := mustNewDraft(t, s, b.ID)
	_, err = s.Publish(ctx, b.ID, v.Version)
	wantErr(t, "empty version", err, registry.ErrInvalid)

	bad := Question("Q", []h5p.SkillRef{{Taxonomy: "tax"}}, nil)
	v.Questions = []h5p.Question{bad}
	mustSave(t, s, v)
	_, err = s.Publish(ctx, b.ID, v.Version)
	wantErr(t, "invalid extension block", err, registry.ErrInvalid)

	got, err := s.GetVersion(ctx, b.ID, v.Version)
	if err != nil || got.Status != registry.StatusDraft {
		t.Errorf("a rejected publish changed the version: %+v, %v", got, err)
	}
}

func testPublishedIsImmutable(t *testing.T, s registry.Store) {
	ctx := context.Background()
	b, pub := mustPublished(t, s, registry.Bank{Title: "Frozen"}, Question("Q1", nil, nil))

	edit := pub
	edit.Status = ""
	edit.Questions = []h5p.Question{Question("changed", nil, nil)}
	wantErr(t, "SaveDraft on published", s.SaveDraft(ctx, edit), registry.ErrImmutable)

	_, err := s.Publish(ctx, b.ID, pub.Version)
	wantErr(t, "publish twice", err, registry.ErrImmutable)

	got, err := s.GetVersion(ctx, b.ID, pub.Version)
	if err != nil {
		t.Fatal(err)
	}
	if !SameJSON(t, got.Questions, pub.Questions) {
		t.Error("published content changed")
	}
}

func testNewVersionFromPublished(t *testing.T, s registry.Store) {
	ctx := context.Background()
	b, pub := mustPublished(t, s, registry.Bank{Title: "Next"}, Question("Q1", []h5p.SkillRef{skillA}, nil))

	next := mustNewVersion(t, s, b.ID, pub.Version)
	if next.Version != "1.1.0" || next.Status != registry.StatusDraft {
		t.Errorf("next version = %s %s, want a 1.1.0 draft", next.Version, next.Status)
	}
	if !SameJSON(t, next.Questions, pub.Questions) {
		t.Error("new version did not copy the questions")
	}
	if len(next.Skills) != 0 || next.PublishedAt != nil {
		t.Errorf("draft carried publish-time fields: %+v", next)
	}

	next.Questions = append(next.Questions, Question("Q2", nil, nil))
	mustSave(t, s, next)
	orig, err := s.GetVersion(ctx, b.ID, pub.Version)
	if err != nil {
		t.Fatal(err)
	}
	if len(orig.Questions) != 1 {
		t.Error("editing the new draft changed the version it came from")
	}

	_, err = s.NewVersion(ctx, b.ID, pub.Version)
	wantErr(t, "second draft", err, registry.ErrDraftExists)

	if _, err := s.Publish(ctx, b.ID, next.Version); err != nil {
		t.Fatal(err)
	}
	empty := mustNewDraft(t, s, b.ID)
	if empty.Version != "1.2.0" || len(empty.Questions) != 0 {
		t.Errorf("empty new version = %s with %d questions, want an empty 1.2.0", empty.Version, len(empty.Questions))
	}

	vs, err := s.ListVersions(ctx, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	order := make([]string, len(vs))
	for i, v := range vs {
		order[i] = v.Version
	}
	if !SameJSON(t, order, []string{"1.0.0", "1.1.0", "1.2.0"}) {
		t.Errorf("version order = %v", order)
	}
}

func testListBanksBySkillAndStandard(t *testing.T, s registry.Store) {
	ctx := context.Background()
	a, _ := mustPublished(t, s, registry.Bank{Title: "A"}, Question("Q", []h5p.SkillRef{skillA}, []h5p.StandardRef{stdA}))
	b, _ := mustPublished(t, s, registry.Bank{Title: "B"}, Question("Q", []h5p.SkillRef{skillB}, []h5p.StandardRef{stdA}))
	draftOnly := mustCreate(t, s, registry.Bank{Title: "Draft only"})
	d := mustNewDraft(t, s, draftOnly.ID)
	d.Questions = []h5p.Question{Question("Q", []h5p.SkillRef{skillA}, nil)}
	mustSave(t, s, d)

	ids := func(f registry.ListFilter) []string {
		got, err := s.ListBanks(ctx, f)
		if err != nil {
			t.Fatal(err)
		}
		out := []string{}
		for _, g := range got {
			out = append(out, g.ID)
		}
		return out
	}
	if got := ids(registry.ListFilter{Skill: &skillA}); !sameSet(got, []string{a.ID}) {
		t.Errorf("skill a: %v, want only %s (drafts do not count)", got, a.ID)
	}
	if got := ids(registry.ListFilter{Standard: &stdA}); !sameSet(got, []string{a.ID, b.ID}) {
		t.Errorf("standard: %v", got)
	}
	if got := ids(registry.ListFilter{Skill: &skillB, Standard: &stdA}); !sameSet(got, []string{b.ID}) {
		t.Errorf("skill and standard: %v", got)
	}

	// The latest published version decides.
	next := mustNewVersion(t, s, a.ID, "1.0.0")
	next.Questions = []h5p.Question{Question("Q", []h5p.SkillRef{skillB}, nil)}
	mustSave(t, s, next)
	if _, err := s.Publish(ctx, a.ID, next.Version); err != nil {
		t.Fatal(err)
	}
	if got := ids(registry.ListFilter{Skill: &skillA}); len(got) != 0 {
		t.Errorf("skill a after a newer version dropped it: %v", got)
	}
}

func testFork(t *testing.T, s registry.Store) {
	ctx := context.Background()
	parent, pub := mustPublished(t, s,
		registry.Bank{Title: "Parent", License: registry.License{SPDX: "CC-BY-SA-4.0", Attribution: "Example Author"}},
		Question("Q1", []h5p.SkillRef{skillA}, nil))

	child, draft, err := s.Fork(ctx, parent.ID, pub.Version, registry.Bank{Title: "Child", Language: "es"}, registry.RelationTranslation)
	if err != nil {
		t.Fatalf("Fork: %v", err)
	}
	want := registry.Lineage{ParentBankID: parent.ID, ParentVersion: pub.Version, Relation: registry.RelationTranslation}
	if child.Lineage == nil || *child.Lineage != want {
		t.Errorf("lineage = %+v, want %+v", child.Lineage, want)
	}
	if child.License != parent.License {
		t.Errorf("license = %+v, want the parent's %+v", child.License, parent.License)
	}
	if draft.BankID != child.ID || draft.Version != "1.0.0" || draft.Status != registry.StatusDraft {
		t.Errorf("fork draft = %+v", draft)
	}
	if !SameJSON(t, draft.Questions, pub.Questions) {
		t.Error("fork did not copy the parent's questions")
	}

	own, _, err := s.Fork(ctx, parent.ID, pub.Version, registry.Bank{Title: "Own license", License: registry.License{SPDX: "CC-BY-4.0"}}, registry.RelationRemix)
	if err != nil {
		t.Fatal(err)
	}
	if own.License.SPDX != "CC-BY-4.0" {
		t.Errorf("explicit license replaced: %+v", own.License)
	}

	got, err := s.ListBanks(ctx, registry.ListFilter{ParentBankID: parent.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("banks with parent %s: %d, want 2", parent.ID, len(got))
	}
}

func testForkRejects(t *testing.T, s registry.Store) {
	ctx := context.Background()
	_, _, err := s.Fork(ctx, "missing", "1.0.0", registry.Bank{Title: "X"}, registry.RelationFork)
	wantErr(t, "missing parent", err, registry.ErrNotFound)

	parent := mustCreate(t, s, registry.Bank{Title: "Parent"})
	d := mustNewDraft(t, s, parent.ID)
	_, _, err = s.Fork(ctx, parent.ID, d.Version, registry.Bank{Title: "X"}, registry.RelationFork)
	wantErr(t, "draft parent", err, registry.ErrNotPublished)

	_, pub := mustPublished(t, s, registry.Bank{Title: "Published"}, Question("Q", nil, nil))
	_, _, err = s.Fork(ctx, pub.BankID, pub.Version, registry.Bank{Title: "X"}, "copy")
	wantErr(t, "unknown relation", err, registry.ErrInvalid)
	_, _, err = s.Fork(ctx, pub.BankID, pub.Version, registry.Bank{Title: " "}, registry.RelationFork)
	wantErr(t, "untitled fork", err, registry.ErrInvalid)
	_, _, err = s.Fork(ctx, pub.BankID, pub.Version, registry.Bank{Title: "X", Lineage: &registry.Lineage{
		ParentBankID: "other", ParentVersion: "1.0.0", Relation: registry.RelationFork,
	}}, registry.RelationFork)
	wantErr(t, "caller-set lineage", err, registry.ErrInvalid)
}

func testDescendants(t *testing.T, s registry.Store) {
	ctx := context.Background()
	_, err := s.Descendants(ctx, "missing")
	wantErr(t, "missing bank", err, registry.ErrNotFound)

	root, pub := mustPublished(t, s, registry.Bank{Title: "Root"}, Question("Q", nil, nil))
	child, childDraft, err := s.Fork(ctx, root.ID, pub.Version, registry.Bank{Title: "Child"}, registry.RelationFork)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Publish(ctx, child.ID, childDraft.Version); err != nil {
		t.Fatal(err)
	}
	grandchild, _, err := s.Fork(ctx, child.ID, childDraft.Version, registry.Bank{Title: "Grandchild"}, registry.RelationRemix)
	if err != nil {
		t.Fatal(err)
	}
	mustCreate(t, s, registry.Bank{Title: "Unrelated"})

	ids := func(bankID string) []string {
		got, err := s.Descendants(ctx, bankID)
		if err != nil {
			t.Fatal(err)
		}
		out := []string{}
		for _, g := range got {
			out = append(out, g.ID)
		}
		return out
	}
	if got := ids(root.ID); !sameSet(got, []string{child.ID, grandchild.ID}) {
		t.Errorf("descendants of root = %v", got)
	}
	if got := ids(child.ID); !sameSet(got, []string{grandchild.ID}) {
		t.Errorf("descendants of child = %v", got)
	}
	if got := ids(grandchild.ID); len(got) != 0 {
		t.Errorf("descendants of a leaf = %v", got)
	}
}

func mustNewVersion(t *testing.T, s registry.Store, bankID, from string) registry.BankVersion {
	t.Helper()
	v, err := s.NewVersion(context.Background(), bankID, from)
	if err != nil {
		t.Fatalf("NewVersion(%q, %q): %v", bankID, from, err)
	}
	return v
}
