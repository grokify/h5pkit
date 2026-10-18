package registry

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestStatusAndRelationValid(t *testing.T) {
	for _, s := range []Status{StatusDraft, StatusPublished, StatusArchived} {
		if !s.Valid() {
			t.Errorf("%q should be valid", s)
		}
	}
	if Status("live").Valid() {
		t.Error("unknown status accepted")
	}
	for _, r := range []Relation{RelationFork, RelationRemix, RelationTranslation} {
		if !r.Valid() {
			t.Errorf("%q should be valid", r)
		}
	}
	if Relation("copy").Valid() {
		t.Error("unknown relation accepted")
	}
}

func TestBankValidate(t *testing.T) {
	tests := []struct {
		name string
		bank Bank
		ok   bool
	}{
		{"titled", Bank{Title: "Fractions"}, true},
		{"no title", Bank{Title: "  "}, false},
		{"good lineage", Bank{Title: "T", Lineage: &Lineage{ParentBankID: "p", ParentVersion: "1.0.0", Relation: RelationFork}}, true},
		{"lineage no parent", Bank{Title: "T", Lineage: &Lineage{ParentVersion: "1.0.0", Relation: RelationFork}}, false},
		{"lineage bad version", Bank{Title: "T", Lineage: &Lineage{ParentBankID: "p", ParentVersion: "v1", Relation: RelationFork}}, false},
		{"lineage bad relation", Bank{Title: "T", Lineage: &Lineage{ParentBankID: "p", ParentVersion: "1.0.0", Relation: "copy"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bank.Validate()
			if tt.ok && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.ok && !errors.Is(err, ErrInvalid) {
				t.Errorf("got %v, want ErrInvalid", err)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.1.0", -1},
		{"1.10.0", "1.9.0", 1},
		{"2.0.0", "1.99.99", 1},
		{"0.0.1", "0.0.2", -1},
	}
	for _, tt := range tests {
		got, err := CompareVersions(tt.a, tt.b)
		if err != nil || got != tt.want {
			t.Errorf("CompareVersions(%q, %q) = %d, %v; want %d", tt.a, tt.b, got, err, tt.want)
		}
	}
	for _, bad := range []string{"", "1", "1.0", "1.0.0.0", "v1.0.0", "1.0.x", "01.0.0", "1.-1.0", "1.0.0-beta"} {
		if _, err := CompareVersions(bad, "1.0.0"); !errors.Is(err, ErrInvalid) {
			t.Errorf("%q accepted: %v", bad, err)
		}
	}
}

func TestBankVersionFrozen(t *testing.T) {
	if (&BankVersion{Status: StatusDraft}).Frozen() {
		t.Error("draft reported frozen")
	}
	for _, s := range []Status{StatusPublished, StatusArchived} {
		if !(&BankVersion{Status: s}).Frozen() {
			t.Errorf("%s not frozen", s)
		}
	}
}

func TestBankJSONIsCamelCase(t *testing.T) {
	data, err := json.Marshal(Bank{ID: "b", Title: "T", Lineage: &Lineage{ParentBankID: "p", ParentVersion: "1.0.0", Relation: RelationFork}})
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"parentBankId"`, `"parentVersion"`, `"createdAt"`, `"updatedAt"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("missing %s in %s", key, data)
		}
	}
}
