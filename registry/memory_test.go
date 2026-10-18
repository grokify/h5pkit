package registry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/grokify/h5pkit/registry"
	"github.com/grokify/h5pkit/registry/registrytest"
)

func TestMemoryStoreConformance(t *testing.T) {
	registrytest.RunStoreTests(t, func() registry.Store { return registry.NewMemoryStore() })
}

func TestMemoryStoreOptions(t *testing.T) {
	at := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	n := 0
	s := registry.NewMemoryStore(
		registry.WithClock(func() time.Time { return at }),
		registry.WithIDGenerator(func() (string, error) { n++; return fmt.Sprintf("bank-%d", n), nil }),
	)
	b, err := s.CreateBank(t.Context(), registry.Bank{Title: "Clocked"})
	if err != nil {
		t.Fatal(err)
	}
	if b.ID != "bank-1" || !b.CreatedAt.Equal(at) || !b.UpdatedAt.Equal(at) {
		t.Errorf("bank = %+v", b)
	}
}
