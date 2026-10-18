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
