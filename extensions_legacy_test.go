package h5p

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExtensionsReadsLegacyKey(t *testing.T) {
	var ext Extensions
	if err := json.Unmarshal([]byte(`{"h5pGo":{"section":"Legacy","questionNumber":3}}`), &ext); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ext.H5PKit == nil || ext.H5PKit.Section != "Legacy" || ext.H5PKit.QuestionNumber != 3 {
		t.Fatalf("legacy block not promoted: %+v", ext.H5PKit)
	}
}

func TestExtensionsPrefersNewKey(t *testing.T) {
	var ext Extensions
	data := `{"h5pkit":{"section":"New"},"h5pGo":{"section":"Old"}}`
	if err := json.Unmarshal([]byte(data), &ext); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ext.H5PKit == nil || ext.H5PKit.Section != "New" {
		t.Fatalf("h5pkit block not preferred: %+v", ext.H5PKit)
	}
}

func TestExtensionsNeverWritesLegacyKey(t *testing.T) {
	tests := []struct {
		name string
		ext  Extensions
	}{
		{"new field", Extensions{H5PKit: &Extension{Section: "S"}}},
		{"legacy field only", Extensions{H5PGo: &Extension{Section: "S"}}},
		{"both fields", Extensions{H5PKit: &Extension{Section: "S"}, H5PGo: &Extension{Section: "old"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.ext)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			got := string(data)
			if strings.Contains(got, "h5pGo") {
				t.Errorf("legacy key written: %s", got)
			}
			if got != `{"h5pkit":{"section":"S"}}` {
				t.Errorf("got %s", got)
			}
		})
	}
}

func TestLegacyContentReserializesUnderNewKey(t *testing.T) {
	legacy := `{"questions":[{"library":"H5P.MultiChoice 1.16","params":{"question":"Q?"},"extensions":{"h5pGo":{"section":"S","tags":["a"]}}}]}`
	qs, err := FromJSON([]byte(legacy))
	if err != nil {
		t.Fatalf("FromJSON: %v", err)
	}
	data, err := json.Marshal(qs)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(data), "h5pGo") || !strings.Contains(string(data), `"h5pkit":{"section":"S","tags":["a"]}`) {
		t.Errorf("not re-serialized under h5pkit: %s", data)
	}
}

func TestMigrateLegacyExtensions(t *testing.T) {
	qs := &QuestionSet{Questions: []Question{
		{Library: "H5P.MultiChoice 1.16", Extensions: &Extensions{H5PGo: &Extension{Section: "A"}}},
		{Library: "H5P.MultiChoice 1.16", Extensions: &Extensions{H5PKit: &Extension{Section: "B"}}},
		{Library: "H5P.MultiChoice 1.16", Extensions: &Extensions{H5PKit: &Extension{Section: "C"}, H5PGo: &Extension{Section: "old"}}},
		{Library: "H5P.MultiChoice 1.16"},
	}}
	if n := MigrateLegacyExtensions(qs); n != 2 {
		t.Errorf("migrated %d questions, want 2", n)
	}
	want := []string{"A", "B", "C"}
	for i, w := range want {
		ext := qs.Questions[i].Extensions
		if ext.H5PGo != nil {
			t.Errorf("question %d still has a legacy block", i)
		}
		if ext.H5PKit == nil || ext.H5PKit.Section != w {
			t.Errorf("question %d: got %+v, want section %q", i, ext.H5PKit, w)
		}
	}
	if MigrateLegacyExtensions(nil) != 0 {
		t.Error("nil question set should migrate nothing")
	}
}

func TestExtensionAccessor(t *testing.T) {
	var nilExt *Extensions
	if nilExt.Extension() != nil {
		t.Error("nil Extensions should return nil")
	}
	legacy := &Extensions{H5PGo: &Extension{Section: "L"}}
	if got := legacy.Extension(); got == nil || got.Section != "L" {
		t.Errorf("legacy block not returned: %+v", got)
	}
}
