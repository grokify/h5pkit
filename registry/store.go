package registry

import (
	"context"
	"errors"
)

// Errors returned by a Store. Implementations wrap them, so callers test
// with errors.Is.
var (
	// ErrNotFound means the bank or version does not exist.
	ErrNotFound = errors.New("registry: not found")
	// ErrExists means a bank with the same ID or slug already exists.
	ErrExists = errors.New("registry: already exists")
	// ErrDraftExists means the bank already has a draft version; publish or
	// finish it before starting another.
	ErrDraftExists = errors.New("registry: bank already has a draft version")
)

// Store persists banks and their versions. Every implementation must pass
// the registrytest conformance suite, which defines the rules below
// precisely.
//
// Rules common to every method:
//
//   - Returned values are copies; changing them never changes stored state.
//   - A canceled context returns its error.
//   - A bank has at most one draft version at a time.
type Store interface {
	// CreateBank stores a new bank. An empty ID is assigned by the store;
	// CreatedAt and UpdatedAt are set by the store. The bank must pass
	// Bank.Validate and have no Lineage: derived banks are created with
	// Fork, so lineage always names real content. A duplicate ID or
	// non-empty slug returns ErrExists.
	CreateBank(ctx context.Context, b Bank) (Bank, error)

	// GetBank returns a bank by ID, or ErrNotFound.
	GetBank(ctx context.Context, id string) (Bank, error)

	// ListBanks returns the banks matching the filter, ordered by
	// CreatedAt and then ID.
	ListBanks(ctx context.Context, f ListFilter) ([]Bank, error)

	// NewVersion starts a draft. With an empty fromVersion the draft is
	// empty; otherwise it copies the questions and notes of fromVersion.
	// The draft is numbered by the store: "1.0.0" for a bank's first
	// version, else the next minor version after the highest existing one.
	// A bank that already has a draft returns ErrDraftExists.
	NewVersion(ctx context.Context, bankID, fromVersion string) (BankVersion, error)

	// SaveDraft replaces a draft's questions and notes. The version is
	// identified by BankID and Version, and its Status must be empty or
	// draft. Skills, Standards, and PublishedAt are ignored.
	SaveDraft(ctx context.Context, v BankVersion) error

	// GetVersion returns one version of a bank, or ErrNotFound.
	GetVersion(ctx context.Context, bankID, version string) (BankVersion, error)

	// ListVersions returns a bank's versions in ascending version order,
	// or ErrNotFound if the bank does not exist.
	ListVersions(ctx context.Context, bankID string) ([]BankVersion, error)
}
