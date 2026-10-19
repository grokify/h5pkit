package registry_test

import (
	"archive/zip"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	h5p "github.com/grokify/h5pkit"
	"github.com/grokify/h5pkit/registry"
	"github.com/grokify/h5pkit/registry/registrytest"
)

func publishedVersion(t *testing.T) registry.BankVersion {
	t.Helper()
	ctx := t.Context()
	s := registry.NewMemoryStore()
	b, err := s.CreateBank(ctx, registry.Bank{Title: "Interop"})
	if err != nil {
		t.Fatal(err)
	}
	v, err := s.NewVersion(ctx, b.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	skill := h5p.SkillRef{Taxonomy: "tax", ID: "a"}
	std := h5p.StandardRef{Framework: "fw", Code: "1.A"}
	q2 := registrytest.Question("Q2", nil, []h5p.StandardRef{std})
	q2.Library = "H5P.TrueFalse 1.8"
	v.Questions = []h5p.Question{
		registrytest.Question("Q1", []h5p.SkillRef{skill}, nil),
		q2,
		registrytest.Question("Q3", nil, nil),
	}
	if err := s.SaveDraft(ctx, v); err != nil {
		t.Fatal(err)
	}
	pub, err := s.Publish(ctx, b.ID, v.Version)
	if err != nil {
		t.Fatal(err)
	}
	return pub
}

func TestExportQuestionSet(t *testing.T) {
	v := publishedVersion(t)
	qs, err := registry.ExportQuestionSet(v, registry.ExportOptions{Title: "Quiz", Introduction: "Hi", PassPercentage: 70})
	if err != nil {
		t.Fatalf("ExportQuestionSet: %v", err)
	}
	if qs.Title != "Quiz" || !qs.ShowIntroPage || qs.PassPercentage != 70 || len(qs.Questions) != 3 {
		t.Errorf("set = %+v", qs)
	}
	qs.Questions[0].Library = "changed"
	if v.Questions[0].Library == "changed" {
		t.Error("exported set shares questions with the version")
	}

	def, err := registry.ExportQuestionSet(v, registry.ExportOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if def.Title != v.BankID+" "+v.Version {
		t.Errorf("default title = %q", def.Title)
	}

	if _, err := registry.ExportQuestionSet(registry.BankVersion{}, registry.ExportOptions{}); !errors.Is(err, registry.ErrInvalid) {
		t.Errorf("empty version exported: %v", err)
	}
}

func TestPackageRoundTrip(t *testing.T) {
	v := publishedVersion(t)
	path := filepath.Join(t.TempDir(), "bank.h5p")
	if err := registry.ExportPackage(path, v, registry.ExportOptions{Title: "Quiz", Language: "en"}); err != nil {
		t.Fatalf("ExportPackage: %v", err)
	}

	// The package passes the root package's loader.
	pkg, err := h5p.LoadH5PPackage(path)
	if err != nil {
		t.Fatalf("LoadH5PPackage: %v", err)
	}
	def := pkg.PackageDefinition
	if def == nil || def.MainLibrary != "H5P.QuestionSet" || def.Title != "Quiz" || def.Language != "en" {
		t.Fatalf("h5p.json = %+v", def)
	}
	want := []h5p.LibraryDependency{
		{MachineName: "H5P.QuestionSet", MajorVersion: 1, MinorVersion: 21},
		{MachineName: "H5P.MultiChoice", MajorVersion: 1, MinorVersion: 16},
		{MachineName: "H5P.TrueFalse", MajorVersion: 1, MinorVersion: 8},
	}
	if !registrytest.SameJSON(t, def.PreloadedDependencies, want) {
		t.Errorf("dependencies = %+v", def.PreloadedDependencies)
	}

	got, err := registry.ImportPackage(path)
	if err != nil {
		t.Fatalf("ImportPackage: %v", err)
	}
	if got.Status != registry.StatusDraft || got.BankID != "" || got.Version != "" {
		t.Errorf("imported version = %+v", got)
	}
	if !registrytest.SameJSON(t, got.Questions, v.Questions) {
		t.Error("export then import changed the questions")
	}
	skills, standards := registry.Unions(got.Questions)
	if !registrytest.SameJSON(t, skills, v.Skills) || !registrytest.SameJSON(t, standards, v.Standards) {
		t.Error("extension blocks were not preserved")
	}
}

func writeZip(t *testing.T, files map[string]string) string {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "in.h5p")
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

const questionSetH5P = `{"title":"T","language":"en","mainLibrary":"H5P.QuestionSet","embedTypes":["iframe"],"preloadedDependencies":[]}`

func TestImportPreservesExtensionBlocks(t *testing.T) {
	path := writeZip(t, map[string]string{
		"h5p.json":             questionSetH5P,
		"content/content.json": `{"questions":[{"library":"H5P.MultiChoice 1.16","params":{"question":"Q"},"extensions":{"h5pkit":{"section":"Kept","skills":[{"taxonomy":"tax","id":"a"}]}}}]}`,
	})
	v, err := registry.ImportPackage(path)
	if err != nil {
		t.Fatalf("ImportPackage: %v", err)
	}
	ext := v.Questions[0].Extensions
	if ext.H5PKit == nil || ext.H5PKit.Section != "Kept" || len(ext.H5PKit.Skills) != 1 {
		t.Errorf("extension block not preserved: %+v", ext)
	}
}

func TestImportReadsWrappedContent(t *testing.T) {
	path := writeZip(t, map[string]string{
		"h5p.json":             questionSetH5P,
		"content/content.json": `{"questionSet":{"questions":[{"library":"H5P.MultiChoice 1.16","params":{"question":"Q"}}]}}`,
	})
	v, err := registry.ImportPackage(path)
	if err != nil {
		t.Fatalf("ImportPackage: %v", err)
	}
	if len(v.Questions) != 1 {
		t.Errorf("questions = %d", len(v.Questions))
	}
}

func TestImportRejects(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  string
	}{
		{"no h5p.json", map[string]string{"content/content.json": `{"questions":[]}`}, "h5p.json"},
		{"not a question set", map[string]string{"h5p.json": `{"mainLibrary":"H5P.MultiChoice"}`, "content/content.json": `{}`}, "main library"},
		{"no content", map[string]string{"h5p.json": questionSetH5P}, "content.json"},
		{"no questions", map[string]string{"h5p.json": questionSetH5P, "content/content.json": `{"title":"x"}`}, "no question set"},
		{"bad json", map[string]string{"h5p.json": questionSetH5P, "content/content.json": `{`}, "content.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := registry.ImportPackage(writeZip(t, tt.files))
			if !errors.Is(err, registry.ErrInvalid) || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("got %v, want ErrInvalid mentioning %q", err, tt.want)
			}
		})
	}
	if _, err := registry.ImportPackage(filepath.Join(t.TempDir(), "missing.h5p")); err == nil {
		t.Error("missing file imported")
	}
}

func TestExportRejectsMalformedLibrary(t *testing.T) {
	v := registry.BankVersion{Questions: []h5p.Question{{Library: "H5P.MultiChoice", Params: map[string]any{}}}}
	var buf bytes.Buffer
	if err := registry.WritePackage(&buf, v, registry.ExportOptions{}); !errors.Is(err, registry.ErrInvalid) {
		t.Errorf("got %v, want ErrInvalid", err)
	}
}
