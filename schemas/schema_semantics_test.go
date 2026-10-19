package schemas

import (
	"encoding/json"
	"testing"

	"github.com/grokify/h5pkit/semantics"
)

func TestSchemaSemantics(t *testing.T) {
	var schemaSemenaticsTests = []struct {
		v []byte
	}{
		{EssaySemanticsBytes},
		{MultiChoiceSemanticsBytes},
		{TrueFalseSemanticsBytes}}

	for _, tt := range schemaSemenaticsTests {
		try := []semantics.Field{}
		err := json.Unmarshal(tt.v, &try)
		if err != nil {
			t.Errorf("]semantics.Field unmarshal error %v", err)
		}
	}
}

func TestBehaviourExplicitFalse(t *testing.T) {
	b := Behaviour{RandomAnswers: Bool(false), EnableCheckButton: Bool(false)}
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(data)
	want := `{"enableCheckButton":false,"randomAnswers":false}`
	if got != want {
		t.Errorf("Marshal() = %s, want %s", got, want)
	}

	data, err = json.Marshal(Behaviour{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `{}` {
		t.Errorf("Marshal(empty) = %s, want {}", data)
	}
}
