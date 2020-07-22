package nits

// This file contains code related to answering cases.

import "math/rand"

// --------------------------------------------------------------------

// Case is the structure for a case.
type Case struct {
	Text []string
	ShortName  string
	RootEvents []Event
	preproc *preprocessedCase
}

// subQuestion is the interface that a sub question type needs
// to implement.
type subQuestion interface {
	getTag() string
	getConcepts() []*Concept
	ask(*Case, *userInterface, *studentState) bool
}

// sqMap is a global map of sub question types.
var sqMap = make(map[string]subQuestion)

// addSubQuestion is a method for registering a subQuestion in the
// global map.
func addSubQuestion(sq subQuestion) subQuestion {
	sqMap[sq.getTag()] = sq
	return sq
}

func (c *Case) getShortName() string {
	return c.ShortName
}

// getConcepts returns all the concepts in a case. Right now it returns all
// the concepts involved in all the sub questions. This is not entirely true
// because some sub questions might not apply to all cases. However at this
// time we do not yet have the machinery to make this distinction.
func (c *Case) getConcepts() []*Concept {
	result := make(map[*Concept]interface{}, 0)
	for _, sq := range sqMap {
		for _, c := range sq.getConcepts() {
			result[c] = nil
		}
	}

	ret := make([]*Concept, 0, len(result))
	for c, _ := range result {
		ret = append(ret, c)
	}

	return ret
}

// Gets all the concepts involved in a sub question. If the sub question is
// not specified (nil) you will get all the concepts in this case.
func (c *Case) getTrainingConcepts(sq subQuestion) []*Concept {
	if sq == nil {
		return c.getConcepts()
	}

	return sqMap[sq.getTag()].getConcepts()
}

// check checks the validity of a case. Not implemented at this time.
func (c *Case) check() {
	c.preprocess()
}

// pushSubQuestionCommandContext pushes a command context on the stack
// that adds ui commands relevant while answering sub questions.
func pushSubQuestionCommandContext(ui *userInterface, displaySubQuestion func([]string) bool) {
	ui.pushCommandContext(&CommandContext{
		description: "Answering a sub question in a case",
		commands: []*Command{
			{
				aliases:  []string{"subagain"},
				global:   true,
				help:     "Displays the sub question again.",
				executor: displaySubQuestion,
			},
		},
	})
}

// selectSubQuestion selects a sub question to answer. It will select a sub
// question whose concepts have not yet been mastered and that has not been
// asked two times already.
func (c *Case) selectSubQuestion(state *studentState, done map[subQuestion]int) subQuestion {
	//REMOVEME
	return &negligencePerSeSubQuestion{}
	// Let's check what the new scores are...
	state.train()

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

	sq := possibles[rand.Int() % len(possibles)]
	if trace != nil {
		trace.println("Returning sub question: %s", sq.getTag())
	}
	return sq
}

// ask asks a case question. It will ask sub questions until they are exhausted.
func (c *Case) ask(ui *userInterface, state *studentState) {
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

	// Map that keeps track of how often we have already asked a sub question.
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

