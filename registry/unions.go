package registry

import (
	"cmp"
	"slices"

	h5p "github.com/grokify/h5pkit"
)

// Unions returns the skills and standards referenced by the questions'
// h5pkit blocks, de-duplicated and sorted. Stores call it on publish; it is
// exported so an implementation computes the same unions.
func Unions(questions []h5p.Question) ([]h5p.SkillRef, []h5p.StandardRef) {
	var skills []h5p.SkillRef
	var standards []h5p.StandardRef
	for _, q := range questions {
		ext := q.Extensions.Extension()
		if ext == nil {
			continue
		}
		skills = append(skills, ext.Skills...)
		standards = append(standards, ext.Standards...)
	}
	slices.SortFunc(skills, func(a, b h5p.SkillRef) int {
		return cmp.Or(cmp.Compare(a.Taxonomy, b.Taxonomy), cmp.Compare(a.ID, b.ID))
	})
	slices.SortFunc(standards, func(a, b h5p.StandardRef) int {
		return cmp.Or(cmp.Compare(a.Framework, b.Framework), cmp.Compare(a.Code, b.Code))
	})
	return slices.Compact(skills), slices.Compact(standards)
}
