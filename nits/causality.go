package nits

type causeInFactSubQuestion struct{}

func (cif *causeInFactSubQuestion) getTag() string {
	return "causeInFact"
}

func (cif *causeInFactSubQuestion) getConcepts() []*Concept {
	return []*Concept{CauseInFact1}
}

var _ = addSubQuestion(&causeInFactSubQuestion{})

func (cif *causeInFactSubQuestion) ask(c *Case, ui *userInterface, state *studentState) bool {
	dam := c.preproc.randomInjuryOrDamage()
	act := c.preproc.randomAct()
	rightAnswer := isCauseInFact(dam, act)

	displayQuestion := func([]string) bool {
		ui.println("In this case, is the act:")
		ui.println(act.Description)
		ui.println("a cause-in-fact of this injury or property damage:")
		ui.println(dam.GetDescription())
		ui.newline()
		return false
	}

	displayQuestion(nil)
	pushSubQuestionCommandContext(state, ui, displayQuestion)
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
