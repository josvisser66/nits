package nits

// This file implements the Negligence Per Se sub question.

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

// ask asks the negligence per se sub question.
func (n *negligencePerSeSubQuestion) ask(c *Case, ui *userInterface, state *studentState) bool {
	// First we collect all the broken legal requirements into a slice and
	// shuffle it.
  pp := c.preprocess()
	blrs := make([]*BrokenLegalRequirement, 0, len(pp.brokenLegalRequirement))
	for b, _ := range pp.brokenLegalRequirement {
		blrs = append(blrs, b)
	}
	rand.Shuffle(len(blrs), func(i, j int) {
		blrs[i], blrs[j] = blrs[j], blrs[i]
	})
	// Then we find a damage that is in the direct and indirect consequences
	// of a broken legal requirement.
	var dams []InjuryOrDamage
	var blr *BrokenLegalRequirement
	for _, b := range blrs {
		dams = findDamages(b.Consequences)
		if len(dams) > 0 {
			blr = b
			break
		}
	}
	// At this point blr contains the broken legal requirement and dams the
	// damages that are in the downward chain.
	if len(dams) == 0 || blr == nil {
		return false
	}
	// Shuffle the damages.
	rand.Shuffle(len(dams), func(i, j int) {
		dams[i], dams[j] = dams[j], dams[i]
	})
	// This is the damage we are going to ask about.
	dam := dams[0]
	// Finds the breached duties that are causes of these damages.
	// This will find defendants.
	duties := findDuties(dam)
	if len(duties) == 0 {
		return false
	}
	// Finds the persons who owed these duties.
	persons := make([]*Person, 0)
	for p := range collectPersonsFromDuties(duties) {
		persons = append(persons, p)
	}
	rand.Shuffle(len(persons), func(i, j int) {
		persons[i], persons[j] = persons[j], persons[i]
	})
	// Now find a person who has committed an act that has this damage
	// in the indirect consequences.
	var defendant *Person
	outer:
	for _, p := range persons {
		for event, _ := range pp.events {
			if act , ok := event.(*Act); ok && act.Person == p {
				for _, d := range findDamages(act.Consequences) {
					if d == dam {
						defendant = p
						break outer
					}
				}
			}
		}
	}
	if defendant == nil {
		return false
	}
	attempts := 0
	for {
		ui.newline()
		ui.println("Looking at the following damage:")
		ui.println(dam.GetDescription())
		ui.println("Which legal principle can defendant %s call in?", defendant.Name)
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
