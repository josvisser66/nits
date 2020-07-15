package nits

// --------------------------------------------------------------------
type Case struct {
	ShortName  string
	RootEvents []*Event
}

func (c *Case) getShortName() string {
	return c.ShortName
}

func (c *Case) getConcepts() []*Concept {
	return nil
}

func (c *Case) getTrainingConcepts() []*Concept {
	return nil
}

func (c *Case) check() {
}

func (c *Case) ask(ui *userInterface, state *studentState) {
}
