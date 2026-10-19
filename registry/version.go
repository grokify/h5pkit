package registry

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// semver is a parsed MAJOR.MINOR.PATCH version.
type semver struct{ major, minor, patch int }

var errVersionFormat = errors.New("version must be MAJOR.MINOR.PATCH with non-negative integers")

func parseVersion(s string) (semver, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return semver{}, fmt.Errorf("%w: %q", errVersionFormat, s)
	}
	var n [3]int
	for i, p := range parts {
		if p == "" || (len(p) > 1 && p[0] == '0') {
			return semver{}, fmt.Errorf("%w: %q", errVersionFormat, s)
		}
		v, err := strconv.Atoi(p)
		if err != nil || v < 0 {
			return semver{}, fmt.Errorf("%w: %q", errVersionFormat, s)
		}
		n[i] = v
	}
	return semver{n[0], n[1], n[2]}, nil
}

func (v semver) String() string { return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch) }

func (v semver) compare(o semver) int {
	switch {
	case v.major != o.major:
		return cmpInt(v.major, o.major)
	case v.minor != o.minor:
		return cmpInt(v.minor, o.minor)
	default:
		return cmpInt(v.patch, o.patch)
	}
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// CompareVersions orders two semantic versions: -1 if a < b, 0 if equal, 1
// if a > b. It returns an error wrapping ErrInvalid if either is malformed.
func CompareVersions(a, b string) (int, error) {
	va, err := parseVersion(a)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	vb, err := parseVersion(b)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	return va.compare(vb), nil
}
