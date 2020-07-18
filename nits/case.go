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

type subQuestion interface {
	getTag() string
	getConcepts() []*Concept
	ask(*Case, *userInterface, *studentState) bool
}

var sqMap = make(map[string]subQuestion)

func addSubQuestion(sq subQuestion) subQuestion {
	sqMap[sq.getTag()] = sq
	return sq
}

func (c *Case) getShortName() string {
	return c.ShortName
}

func (c *Case) getConcepts() []*Concept {
	return nil
}

func (c *Case) getTrainingConcepts(sq subQuestion) []*Concept {
	if sq != nil {
		return sq.getConcepts()
	}

	m := make(map[*Concept]interface{})

	for _, sq := range sqMap {
		for _, c := range sq.getConcepts() {
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

func (c *Case) selectSubQuestion(state *studentState, done map[subQuestion]int) subQuestion {
	possibles := make([]subQuestion, 0)

	for _, sq := range sqMap {
		nm := state.conceptsNotMastered(sq.getConcepts())
		if len(nm) > 0 && done[sq] < 2 {
			possibles = append(possibles, sq)
		}
	}

	if len(possibles) == 0 {
		return nil
	}

	return possibles[rand.Int() % len(possibles)]
}

func (c *Case) ask(ui *userInterface, state *studentState) {
	preprocess(c)
	displayCase := func([]string) bool {
		ui.printParagraphs(c.Text)
		ui.newline()
		return false
	}

	displayCase(nil)
	pushCommandContext("Answering a case question", state, ui, c, displayCase)
	ui.pushPrompt("Your answer? ")
	defer ui.popCommandContext()
	defer ui.popPrompt()

	done := make(map[subQuestion]int)

	for {
		sq := c.selectSubQuestion(state, done)
		if sq == nil {
			ui.println("Nothing left to ask in this case.")
			return
		}
		ret := sq.ask(c, ui, state)
		if ret {
			return
		}
		done[sq]++
	}
}

