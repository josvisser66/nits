package nits

import (
	"math/rand"
)

// --------------------------------------------------------------------
type Case struct {
	Text []string
	ShortName  string
	RootEvents []Event
	preproc *preprocessedCase
}

type subQuestion struct {
	tag string
	concepts []*Concept
}

var sqMap = make(map[string]*subQuestion)

func (sq *subQuestion) add() *subQuestion {
	sqMap[sq.tag] = sq
	return sq
}

func (c *Case) getShortName() string {
	return c.ShortName
}

func (c *Case) getConcepts() []*Concept {
	return nil
}

func (c *Case) getTrainingConcepts(sq *subQuestion) []*Concept {
	if sq != nil {
		return sq.concepts
	}

	m := make(map[*Concept]interface{})

	for _, sq := range sqMap {
		for _, c := range sq.concepts {
			m[c] = nil
		}
	}

	result := make([]*Concept, 0, len(m))

	for c, _ := range m {
		result = append(result, c)
	}

	return result
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
		case CauseInFact1:
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

