package nits

import (
	"math/rand"
	"strings"
)

type negligencePerSeSubQuestion struct{}

func (n *negligencePerSeSubQuestion) getTag() string {
	return "negligencePerSe"
}

func (n *negligencePerSeSubQuestion) getConcepts() []*Concept {
	return []*Concept{NegligencePerSe1}
}

var _ = addSubQuestion(&negligencePerSeSubQuestion{})

func (n *negligencePerSeSubQuestion) ask(c *Case, ui *userInterface, state *studentState) bool {
	var dams []InjuryOrDamage
	var blr *BrokenLegalRequirement
	for b, _ := range c.preproc.brokenLegalRequirement {
		dams = findDamages(b.Consequences)
		if len(dams) > 0 {
			blr = b
			break
		}
	}
	if len(dams) == 0 || blr == nil {
		return false
	}
	rand.Shuffle(len(dams), func(i, j int) {
		dams[i], dams[j] = dams[j], dams[i]
	})
	dam := dams[0]
	duties := findDuties(dam, dam.getDirectCauses())
	if len(duties) == 0 {
		return false
	}
	persons := make([]*Person, 0)
	for p, _ := range collectPersonsFromDuties(duties) {
		persons = append(persons, p)
	}
	rand.Shuffle(len(persons), func (i, j int) {
		persons[i], persons[j] = persons[j], persons[i]
	})
	attempts := 0
	for {
		ui.newline()
		ui.println("Looking at the following damage:")
		ui.println(dam.GetDescription())
		ui.println("Which legal principle can defendant %s call in?", persons[0].Name)
		words, ret := ui.getInput()
		if ret {
			return ret
		}
		if strings.Join(words, " ") == "negligence per se" {
			ui.println("Correct!")
			state.registerAnswer(c, n, attempts == 0)
			return false
		}
		ui.println("Incorrect :-(")
		attempts++
	}
}

