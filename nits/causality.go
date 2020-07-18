package nits

const causeInFactSubQuestion = "causeInFact"

func isParentOf(suspectedParent, child Event) bool {
	if suspectedParent == child {
		return true
	}
	for _, cause := range child.getCauses() {
		if isParentOf(suspectedParent, cause) {
			return true
		}
	}
	return false
}

func (p *preprocessedCase) isCauseInFact(dam InjuryOrDamage, event Event) bool {
	for _, e := range dam.getCauses() {
		if isParentOf(e, event) {
			return true
		}
	}
	return false
}

func (c *Case) askCauseInFact(ui *userInterface, state *studentState) bool {
	dam := c.preproc.randomInjuryOrDamage()
	event := c.preproc.randomEvent()
	rightAnswer := c.preproc.isCauseInFact(dam, event)

	displayQuestion := func([]string) bool {
		ui.println("In this case, is the event:")
		ui.println(event.getDescription())
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
		words, ret := ui.getInput()
		if ret {
			return ret
		}
		if len(words) == 0 {
			continue
		}
		if answer, err := isYesNo(words[0]); err != nil {
			ui.println("Please enter a yes or no answer")
			continue
		} else if answer != rightAnswer {
			ui.println("Please try again :-(")
			attempts++
			continue
		}
		ui.println("Correct :-)")
		state.registerAnswer(c, causeInFactSubQuestion, attempts == 0)
		return false
	}
}
