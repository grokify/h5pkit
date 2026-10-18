package registry

import (
	"archive/zip"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	h5p "github.com/grokify/h5pkit"
)

// QuestionSetLibrary is the H5P library an exported package's content runs
// in, pinned to a released H5P.QuestionSet version.
const QuestionSetLibrary = "H5P.QuestionSet 1.21"

// maxEntryBytes bounds how much of one package entry is read on import.
const maxEntryBytes = 64 << 20

// ExportOptions describes the playable set built from a bank version.
type ExportOptions struct {
	// Title is the set's title. It defaults to the bank ID and version.
	Title string
	// Introduction is shown before the first question when set.
	Introduction string
	// PassPercentage is the score needed to pass, 0 to 100.
	PassPercentage int
	// Language is the package language for h5p.json. It defaults to "und".
	Language string
	// License is the h5p.json license code, such as "CC BY". It defaults
	// to "U" (undisclosed).
	License string
	// Author is the h5p.json author.
	Author string
}

// ExportQuestionSet builds a playable question set from a version's
// questions. The questions are copied, so the result does not share state
// with v. The set must pass h5p.QuestionSet validation.
func ExportQuestionSet(v BankVersion, opts ExportOptions) (*h5p.QuestionSet, error) {
	questions, err := copyQuestions(v.Questions)
	if err != nil {
		return nil, err
	}
	title := opts.Title
	if title == "" {
		title = strings.TrimSpace(v.BankID + " " + v.Version)
	}
	qs := &h5p.QuestionSet{
		Title:          title,
		Introduction:   opts.Introduction,
		ShowIntroPage:  opts.Introduction != "",
		PassPercentage: opts.PassPercentage,
		Questions:      questions,
	}
	if err := qs.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	return qs, nil
}

// ExportPackage writes a version as an .h5p package at path. The package
// holds h5p.json, naming H5P.QuestionSet as the main library and every
// question library as a preloaded dependency, and content/content.json,
// holding the question set. It does not bundle library code: the H5P
// platform that opens it supplies the libraries.
func ExportPackage(path string, v BankVersion, opts ExportOptions) (err error) {
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("registry: create package: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("registry: close package: %w", cerr)
		}
	}()
	return WritePackage(f, v, opts)
}

// WritePackage writes a version as an .h5p package to w. See ExportPackage.
func WritePackage(w io.Writer, v BankVersion, opts ExportOptions) error {
	qs, err := ExportQuestionSet(v, opts)
	if err != nil {
		return err
	}
	deps, err := dependencies(qs.Questions)
	if err != nil {
		return err
	}
	def := h5p.PackageDefinition{
		Title:                 qs.Title,
		Language:              cmp.Or(opts.Language, "und"),
		MainLibrary:           "H5P.QuestionSet",
		EmbedTypes:            []string{"iframe"},
		License:               cmp.Or(opts.License, "U"),
		Author:                opts.Author,
		PreloadedDependencies: deps,
	}
	h5pJSON, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return fmt.Errorf("registry: encode h5p.json: %w", err)
	}
	contentJSON, err := json.MarshalIndent(qs, "", "  ")
	if err != nil {
		return fmt.Errorf("registry: encode content.json: %w", err)
	}
	zw := zip.NewWriter(w)
	for _, entry := range []struct {
		name string
		data []byte
	}{{"h5p.json", h5pJSON}, {"content/content.json", contentJSON}} {
		ew, err := zw.Create(entry.name)
		if err != nil {
			return fmt.Errorf("registry: add %s: %w", entry.name, err)
		}
		if _, err := ew.Write(entry.data); err != nil {
			return fmt.Errorf("registry: write %s: %w", entry.name, err)
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("registry: finish package: %w", err)
	}
	return nil
}

// ImportPackage reads an .h5p question set into a draft version. The
// version's BankID and Version are empty; set them before saving it.
// Extension blocks are preserved. Both the standard content.json layout, with the question set at
// the top level, and the layout written by h5p.H5PPackage, with the set
// under "questionSet", are read.
func ImportPackage(path string) (BankVersion, error) {
	r, err := zip.OpenReader(filepath.Clean(path))
	if err != nil {
		return BankVersion{}, fmt.Errorf("registry: open package: %w", err)
	}
	v, err := readPackage(&r.Reader)
	if cerr := r.Close(); cerr != nil && err == nil {
		err = fmt.Errorf("registry: close package: %w", cerr)
	}
	return v, err
}

func readPackage(r *zip.Reader) (BankVersion, error) {
	h5pJSON, err := readEntry(r, "h5p.json")
	if err != nil {
		return BankVersion{}, err
	}
	var def h5p.PackageDefinition
	if err := json.Unmarshal(h5pJSON, &def); err != nil {
		return BankVersion{}, fmt.Errorf("%w: h5p.json: %w", ErrInvalid, err)
	}
	if def.MainLibrary != "H5P.QuestionSet" {
		return BankVersion{}, fmt.Errorf("%w: main library is %q, not H5P.QuestionSet", ErrInvalid, def.MainLibrary)
	}
	contentJSON, err := readEntry(r, "content/content.json")
	if err != nil {
		return BankVersion{}, err
	}
	qs, err := decodeContent(contentJSON)
	if err != nil {
		return BankVersion{}, err
	}
	if qs.Questions == nil {
		qs.Questions = []h5p.Question{}
	}
	return BankVersion{Status: StatusDraft, Questions: qs.Questions}, nil
}

func decodeContent(data []byte) (*h5p.QuestionSet, error) {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("%w: content.json: %w", ErrInvalid, err)
	}
	if _, ok := probe["questions"]; !ok {
		wrapped, ok := probe["questionSet"]
		if !ok {
			return nil, fmt.Errorf("%w: content.json holds no question set", ErrInvalid)
		}
		data = wrapped
	}
	var qs h5p.QuestionSet
	if err := json.Unmarshal(data, &qs); err != nil {
		return nil, fmt.Errorf("%w: content.json: %w", ErrInvalid, err)
	}
	return &qs, nil
}

func readEntry(r *zip.Reader, name string) ([]byte, error) {
	for _, f := range r.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("registry: open %s: %w", name, err)
		}
		data, err := io.ReadAll(io.LimitReader(rc, maxEntryBytes+1))
		if cerr := rc.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if err != nil {
			return nil, fmt.Errorf("registry: read %s: %w", name, err)
		}
		if len(data) > maxEntryBytes {
			return nil, fmt.Errorf("%w: %s is larger than %d bytes", ErrInvalid, name, maxEntryBytes)
		}
		return data, nil
	}
	return nil, fmt.Errorf("%w: package has no %s", ErrInvalid, name)
}

// dependencies lists H5P.QuestionSet and every question library, parsed
// from "Machine.Name MAJOR.MINOR", without duplicates.
func dependencies(questions []h5p.Question) ([]h5p.LibraryDependency, error) {
	main, err := parseLibrary(QuestionSetLibrary)
	if err != nil {
		panic(err) // QuestionSetLibrary is a well-formed constant
	}
	deps := []h5p.LibraryDependency{main}
	for i, q := range questions {
		d, err := parseLibrary(q.Library)
		if err != nil {
			return nil, fmt.Errorf("%w: question %d: %w", ErrInvalid, i, err)
		}
		if !slices.Contains(deps, d) {
			deps = append(deps, d)
		}
	}
	return deps, nil
}

func parseLibrary(s string) (h5p.LibraryDependency, error) {
	name, version, ok := strings.Cut(strings.TrimSpace(s), " ")
	major, minor, ok2 := strings.Cut(version, ".")
	if !ok || !ok2 || name == "" {
		return h5p.LibraryDependency{}, fmt.Errorf("library %q is not \"Machine.Name MAJOR.MINOR\"", s)
	}
	maj, err := strconv.Atoi(major)
	if err != nil || maj < 0 {
		return h5p.LibraryDependency{}, fmt.Errorf("library %q has a bad major version", s)
	}
	mnr, err := strconv.Atoi(minor)
	if err != nil || mnr < 0 {
		return h5p.LibraryDependency{}, fmt.Errorf("library %q has a bad minor version", s)
	}
	return h5p.LibraryDependency{MachineName: name, MajorVersion: maj, MinorVersion: mnr}, nil
}

func copyQuestions(qs []h5p.Question) ([]h5p.Question, error) {
	data, err := encodeQuestions(qs)
	if err != nil {
		return nil, err
	}
	var out []h5p.Question
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("registry: copy questions: %w", err)
	}
	return out, nil
}
