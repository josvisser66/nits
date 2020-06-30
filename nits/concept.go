package nits

import "sort"

// --------------------------------------------------------------------
type Concept struct {
	Name  string
	Level int
	Explanation *Explanation
	Related []*Concept
}

func (c *Concept) SortRelatedConcepts() {
	sort.Slice(c.Related, func(i, j int) bool {
		return c.Related[i].Name < c.Related[j].Name
	})
}
