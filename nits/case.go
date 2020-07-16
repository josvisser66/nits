package nits

// --------------------------------------------------------------------
type Case struct {
	ShortName  string
	RootEvents []*Event
	preproc preprocessedCase
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
	displayQuestion := func([]string) bool {
		ui.println("this is case %s", c.ShortName)
		return false
	}

	displayQuestion(nil)
	ui.pushPrompt("Your answer? ")
	pushCommandContext(state, ui, c, displayQuestion)
	defer ui.popPrompt()
	defer ui.popCommandContext()

	for {
		_, ret := ui.getInput()
		if ret {
			return
		}
	}
}

