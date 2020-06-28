package nits

// --------------------------------------------------------------------
type Concept struct {
	Name  string
	Level int
	Explanation *Explanation
	Related []*Concept
}
