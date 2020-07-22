package nits

// This file contains the implementation of the causality in fact
// sub question.

type causeInFactSubQuestion struct{}

func (cif *causeInFactSubQuestion) getTag() string {
	return "causeInFact"
}

func (cif *causeInFactSubQuestion) getConcepts() []*Concept {
	return []*Concept{CauseInFact1}
}

var _ = addSubQuestion(&causeInFactSubQuestion{})

// ask asks the sub question
func (cif *causeInFactSubQuestion) ask(c *Case, ui *userInterface, state *studentState) bool {
	// Finds a random piece of damage and a random act in the case and then
	// figures out if they are connected by a link of causality.
  pp := c.preprocess()
	dam := pp.randomInjuryOrDamage()
	act := pp.randomAct()
	rightAnswer := isCauseInFact(dam, act)

	displayQuestion := func([]string) bool {
		ui.newline()
		ui.println("In this case, is the act:")
		ui.println(act.Description)
		ui.println("a cause-in-fact of this injury or property damage:")
		ui.println(dam.GetDescription())
		ui.newline()
		return false
	}

	displayQuestion(nil)
	pushSubQuestionCommandContext(ui, displayQuestion)
	defer ui.popCommandContext()

	attempts := 0

	for {
		answer, ret := ui.yesNo("Your answer")
		if ret {
			return ret
		}
		if answer != rightAnswer {
			ui.println("Please try again :-(")
			attempts++
			continue
		}
		ui.println("Correct :-)")
		state.registerAnswer(c, cif, attempts == 0)
		return false
	}
}
