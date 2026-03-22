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
