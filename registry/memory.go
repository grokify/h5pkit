package registry

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	h5p "github.com/grokify/h5pkit"
)

// MemoryStore is an in-memory Store, safe for concurrent use. It is the
// reference implementation of the Store rules and suits tests, tools, and
// single-process use; its contents are lost when the process exits.
//
// Questions are held as their JSON encoding, so a returned question's
// Params is the generic decoded form (map[string]any), as it would be from
// any persistent store.
type MemoryStore struct {
	mu       sync.RWMutex
	now      func() time.Time
	newID    func() (string, error)
	banks    map[string]Bank
	slugs    map[string]string
	versions map[string]map[string]*storedVersion
}

type storedVersion struct {
	meta      BankVersion // Questions is nil; see questions
	questions []byte
}

// MemoryOption configures a MemoryStore.
type MemoryOption func(*MemoryStore)

// WithClock sets the time source. The default is time.Now.
func WithClock(now func() time.Time) MemoryOption {
	return func(s *MemoryStore) { s.now = now }
}

// WithIDGenerator sets how bank IDs are assigned. The default is 16 random
// bytes, hex encoded.
func WithIDGenerator(newID func() (string, error)) MemoryOption {
	return func(s *MemoryStore) { s.newID = newID }
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore(opts ...MemoryOption) *MemoryStore {
	s := &MemoryStore{
		now:      time.Now,
		newID:    randomID,
		banks:    make(map[string]Bank),
		slugs:    make(map[string]string),
		versions: make(map[string]map[string]*storedVersion),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

var _ Store = (*MemoryStore)(nil)

func randomID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("registry: generate bank id: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}

// CreateBank implements Store.
func (s *MemoryStore) CreateBank(ctx context.Context, b Bank) (Bank, error) {
	if err := ctx.Err(); err != nil {
		return Bank{}, err
	}
	if b.Lineage != nil {
		return Bank{}, fmt.Errorf("%w: lineage is set by Fork, not CreateBank", ErrInvalid)
	}
	if err := b.Validate(); err != nil {
		return Bank{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.insertBank(b)
}

// insertBank stores b, assigning an ID and timestamps. The caller holds mu
// and has validated b.
func (s *MemoryStore) insertBank(b Bank) (Bank, error) {
	if b.ID == "" {
		id, err := s.newID()
		if err != nil {
			return Bank{}, err
		}
		b.ID = id
	}
	if _, ok := s.banks[b.ID]; ok {
		return Bank{}, fmt.Errorf("%w: bank %q", ErrExists, b.ID)
	}
	if b.Slug != "" {
		if _, ok := s.slugs[b.Slug]; ok {
			return Bank{}, fmt.Errorf("%w: slug %q", ErrExists, b.Slug)
		}
		s.slugs[b.Slug] = b.ID
	}
	now := s.now()
	b.CreatedAt, b.UpdatedAt = now, now
	b = copyBank(b)
	s.banks[b.ID] = b
	s.versions[b.ID] = make(map[string]*storedVersion)
	return copyBank(b), nil
}

// GetBank implements Store.
func (s *MemoryStore) GetBank(ctx context.Context, id string) (Bank, error) {
	if err := ctx.Err(); err != nil {
		return Bank{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.banks[id]
	if !ok {
		return Bank{}, fmt.Errorf("%w: bank %q", ErrNotFound, id)
	}
	return copyBank(b), nil
}

// ListBanks implements Store.
func (s *MemoryStore) ListBanks(ctx context.Context, f ListFilter) ([]Bank, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Bank{}
	for _, b := range s.banks {
		if s.matches(b, f) {
			out = append(out, copyBank(b))
		}
	}
	slices.SortFunc(out, func(a, b Bank) int {
		if c := a.CreatedAt.Compare(b.CreatedAt); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return out, nil
}

// matches reports whether b passes f. The caller holds mu.
func (s *MemoryStore) matches(b Bank, f ListFilter) bool {
	switch {
	case f.Tag != "" && !slices.Contains(b.Tags, f.Tag):
		return false
	case f.License != "" && b.License.SPDX != f.License:
		return false
	case f.Subject != "" && b.Subject != f.Subject:
		return false
	case f.Language != "" && !strings.EqualFold(b.Language, f.Language):
		return false
	case f.ParentBankID != "" && (b.Lineage == nil || b.Lineage.ParentBankID != f.ParentBankID):
		return false
	}
	return true
}

// NewVersion implements Store.
func (s *MemoryStore) NewVersion(ctx context.Context, bankID, fromVersion string) (BankVersion, error) {
	if err := ctx.Err(); err != nil {
		return BankVersion{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	versions, ok := s.versions[bankID]
	if !ok {
		return BankVersion{}, fmt.Errorf("%w: bank %q", ErrNotFound, bankID)
	}
	var questions []byte
	var notes string
	if fromVersion != "" {
		from, ok := versions[fromVersion]
		if !ok {
			return BankVersion{}, fmt.Errorf("%w: version %q of bank %q", ErrNotFound, fromVersion, bankID)
		}
		questions, notes = from.questions, from.meta.Notes
	}
	next := semver{1, 0, 0}
	for v, sv := range versions {
		if sv.meta.Status == StatusDraft {
			return BankVersion{}, fmt.Errorf("%w: %s of bank %q", ErrDraftExists, v, bankID)
		}
		if p, err := parseVersion(v); err == nil && p.compare(next) >= 0 {
			next = semver{p.major, p.minor + 1, 0}
		}
	}
	if questions == nil {
		questions = []byte("[]")
	}
	sv := &storedVersion{
		meta:      BankVersion{BankID: bankID, Version: next.String(), Status: StatusDraft, Notes: notes},
		questions: questions,
	}
	versions[sv.meta.Version] = sv
	s.touch(bankID)
	return sv.decode()
}

// SaveDraft implements Store.
func (s *MemoryStore) SaveDraft(ctx context.Context, v BankVersion) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if v.Status != "" && v.Status != StatusDraft {
		return fmt.Errorf("%w: SaveDraft cannot set status %q", ErrInvalid, v.Status)
	}
	questions, err := encodeQuestions(v.Questions)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	sv, err := s.version(v.BankID, v.Version)
	if err != nil {
		return err
	}
	sv.questions = questions
	sv.meta.Notes = v.Notes
	s.touch(v.BankID)
	return nil
}

// GetVersion implements Store.
func (s *MemoryStore) GetVersion(ctx context.Context, bankID, version string) (BankVersion, error) {
	if err := ctx.Err(); err != nil {
		return BankVersion{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	sv, err := s.version(bankID, version)
	if err != nil {
		return BankVersion{}, err
	}
	return sv.decode()
}

// ListVersions implements Store.
func (s *MemoryStore) ListVersions(ctx context.Context, bankID string) ([]BankVersion, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	versions, ok := s.versions[bankID]
	if !ok {
		return nil, fmt.Errorf("%w: bank %q", ErrNotFound, bankID)
	}
	out := make([]BankVersion, 0, len(versions))
	for _, sv := range versions {
		v, err := sv.decode()
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	sortVersions(out)
	return out, nil
}

// version returns a stored version. The caller holds mu.
func (s *MemoryStore) version(bankID, version string) (*storedVersion, error) {
	versions, ok := s.versions[bankID]
	if !ok {
		return nil, fmt.Errorf("%w: bank %q", ErrNotFound, bankID)
	}
	sv, ok := versions[version]
	if !ok {
		return nil, fmt.Errorf("%w: version %q of bank %q", ErrNotFound, version, bankID)
	}
	return sv, nil
}

// touch updates a bank's UpdatedAt. The caller holds mu for writing.
func (s *MemoryStore) touch(bankID string) {
	b := s.banks[bankID]
	b.UpdatedAt = s.now()
	s.banks[bankID] = b
}

func (sv *storedVersion) decode() (BankVersion, error) {
	v := sv.meta
	v.Skills = slices.Clone(sv.meta.Skills)
	v.Standards = slices.Clone(sv.meta.Standards)
	if sv.meta.PublishedAt != nil {
		t := *sv.meta.PublishedAt
		v.PublishedAt = &t
	}
	if err := json.Unmarshal(sv.questions, &v.Questions); err != nil {
		return BankVersion{}, fmt.Errorf("registry: decode questions of %s %s: %w", v.BankID, v.Version, err)
	}
	return v, nil
}

func encodeQuestions(qs []h5p.Question) ([]byte, error) {
	if qs == nil {
		qs = []h5p.Question{}
	}
	data, err := json.Marshal(qs)
	if err != nil {
		return nil, fmt.Errorf("%w: questions cannot be encoded: %w", ErrInvalid, err)
	}
	return data, nil
}

func sortVersions(vs []BankVersion) {
	slices.SortFunc(vs, func(a, b BankVersion) int {
		pa, errA := parseVersion(a.Version)
		pb, errB := parseVersion(b.Version)
		if errA != nil || errB != nil {
			return strings.Compare(a.Version, b.Version)
		}
		return pa.compare(pb)
	})
}

func copyBank(b Bank) Bank {
	b.Tags = slices.Clone(b.Tags)
	if b.Lineage != nil {
		l := *b.Lineage
		b.Lineage = &l
	}
	return b
}
