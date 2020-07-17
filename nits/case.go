package nits

import (
	"fmt"
	"math/rand"
)

// --------------------------------------------------------------------
type Case struct {
	Text []string
	ShortName  string
	RootEvents []*Event
	preproc *preprocessedCase
}


func (c *Case) getShortName() string {
	return c.ShortName
}

func (c *Case) getConcepts() []*Concept {
	return nil
}

func (c *Case) getTrainingConcepts(subQuestion string) []*Concept {
	switch subQuestion {
	case "":
		return []*Concept{CauseInFact}
	case causeInFactSubQuestion:
		return []*Concept{CauseInFact}
	}
	panic(fmt.Sprintf("Unknown case subQuestion %s", subQuestion))
}

func (c *Case) check() {
}

func pushSubQuestionCommandContext(state *studentState, ui *userInterface, displayQuestion func([]string) bool) {
	ui.pushCommandContext(&CommandContext{
		Description: "Answering a question about a case",
		Commands: []*Command{
			{
				Aliases:  []string{"again2"},
				Global:   true,
				Help:     "Displays the subquestion again.",
				Executor: displayQuestion,
			},
		},
	})
}

func (c *Case) ask(ui *userInterface, state *studentState) {
	preprocess(c)
	displayQuestion := func([]string) bool {
		ui.printParagraphs(c.Text)
		ui.newline()
		return false
	}

	displayQuestion(nil)
	pushCommandContext("Answering a case question", state, ui, c, displayQuestion)
	ui.pushPrompt("Your answer? ")
	defer ui.popCommandContext()
	defer ui.popPrompt()

	for {
		var ret bool
		concepts := state.conceptsNotMastered(c)
		concept := concepts[rand.Int() % len(concepts)]
		switch concept {
		case CauseInFact:
			ret = c.askCauseInFact(ui, state)
		default:
			ui.println("Nothing left to ask about this case.")
			ret = true
		}
		if ret {
			return
		}
	}
}

