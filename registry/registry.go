// Package registry manages question banks: named collections of H5P
// questions with licenses, immutable published versions, and lineage between
// a bank and the banks derived from it.
//
// The package holds types, rules, and a [Store] boundary. It has no
// database dependency: [MemoryStore] is the reference implementation, and a
// persistent store lives in the consuming application and proves itself
// with the registrytest conformance suite. Consumer-specific fields belong
// in the consumer's own persistence, keyed by [Bank.ID].
package registry

import (
	"errors"
	"fmt"
	"strings"
	"time"

	h5p "github.com/grokify/h5pkit"
)

// Status is the lifecycle state of a bank version.
type Status string

// Version statuses. A draft can be edited; a published or archived version
// is frozen.
const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

// Valid reports whether s is a known status.
func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusPublished, StatusArchived:
		return true
	}
	return false
}

// Relation is how a derived bank relates to its parent.
type Relation string

// Lineage relations.
const (
	// RelationFork is an independent copy that may diverge freely.
	RelationFork Relation = "fork"
	// RelationRemix combines or reworks the parent's content.
	RelationRemix Relation = "remix"
	// RelationTranslation renders the parent's content in another language.
	RelationTranslation Relation = "translation"
)

// Valid reports whether r is a known relation.
func (r Relation) Valid() bool {
	switch r {
	case RelationFork, RelationRemix, RelationTranslation:
		return true
	}
	return false
}

// License states the terms a bank is offered under.
type License struct {
	// SPDX is an SPDX license identifier such as "CC-BY-4.0" or
	// "CC-BY-SA-4.0", or "proprietary".
	SPDX string `json:"spdx,omitempty"`
	// Attribution is the credit line the license requires.
	Attribution string `json:"attribution,omitempty"`
	// URL links to the license text or the bank's terms.
	URL string `json:"url,omitempty"`
}

// Lineage records the published version a bank was derived from.
type Lineage struct {
	ParentBankID  string   `json:"parentBankId"`
	ParentVersion string   `json:"parentVersion"`
	Relation      Relation `json:"relation"`
}

// Validate checks that the parent and the relation are set.
func (l *Lineage) Validate() error {
	if l == nil {
		return nil
	}
	if strings.TrimSpace(l.ParentBankID) == "" {
		return fmt.Errorf("%w: lineage needs a parent bank", ErrInvalid)
	}
	if _, err := parseVersion(l.ParentVersion); err != nil {
		return fmt.Errorf("%w: lineage parent version: %w", ErrInvalid, err)
	}
	if !l.Relation.Valid() {
		return fmt.Errorf("%w: unknown lineage relation %q", ErrInvalid, l.Relation)
	}
	return nil
}

// Bank is a named collection of questions. Its content lives in versions.
type Bank struct {
	ID          string `json:"id"`
	Slug        string `json:"slug,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Subject     string `json:"subject,omitempty"`
	// Language is a BCP 47 language tag, such as "en" or "pt-BR".
	Language string   `json:"language,omitempty"`
	License  License  `json:"license"`
	Tags     []string `json:"tags,omitempty"`
	// Lineage is nil for an original bank.
	Lineage   *Lineage  `json:"lineage,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate checks the fields a store needs: a title, and well-formed lineage
// when present. The ID may be empty; a store assigns one.
func (b *Bank) Validate() error {
	if strings.TrimSpace(b.Title) == "" {
		return fmt.Errorf("%w: bank needs a title", ErrInvalid)
	}
	return b.Lineage.Validate()
}

// BankVersion is one version of a bank's content.
type BankVersion struct {
	BankID string `json:"bankId"`
	// Version is a semantic version, MAJOR.MINOR.PATCH.
	Version   string         `json:"version"`
	Status    Status         `json:"status"`
	Questions []h5p.Question `json:"questions"`
	// Skills is the union of the questions' skills, computed on publish.
	Skills []h5p.SkillRef `json:"skills,omitempty"`
	// Standards is the union of the questions' standards, computed on
	// publish.
	Standards   []h5p.StandardRef `json:"standards,omitempty"`
	Notes       string            `json:"notes,omitempty"`
	PublishedAt *time.Time        `json:"publishedAt,omitempty"`
}

// Frozen reports whether the version can no longer change.
func (v *BankVersion) Frozen() bool {
	return v.Status == StatusPublished || v.Status == StatusArchived
}

// ListFilter selects banks. Every set field must match; a zero filter
// matches every bank.
type ListFilter struct {
	// Tag matches banks carrying the tag.
	Tag string
	// Skill matches banks whose latest published version covers the skill.
	Skill *h5p.SkillRef
	// Standard matches banks whose latest published version covers the
	// standard.
	Standard *h5p.StandardRef
	// License matches banks whose License.SPDX equals it.
	License string
	// Subject matches banks with this subject.
	Subject string
	// Language matches banks with this language tag, ignoring case.
	Language string
	// ParentBankID matches banks derived directly from this bank.
	ParentBankID string
}

// ErrInvalid reports input that breaks a registry rule. Errors wrapping it
// say which rule.
var ErrInvalid = errors.New("registry: invalid input")
