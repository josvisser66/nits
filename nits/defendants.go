package nits

import (
	"math/rand"
	"sort"
	"strings"
)

type defendantsSubQuestion struct {}

func (p *defendantsSubQuestion) getTag() string {
	return "defendants"
}

func (p *defendantsSubQuestion) getConcepts() []*Concept {
	return []*Concept{Defendant0}
}

var _ = addSubQuestion(&defendantsSubQuestion{})


func (p *defendantsSubQuestion) ask(c *Case, ui *userInterface, state *studentState) bool {
	dams := make([]InjuryOrDamage, 0, len(c.preproc.injuriesOrDamages))
	for dam, _ := range c.preproc.injuriesOrDamages {
		dams = append(dams, dam)
	}
	rand.Shuffle(len(dams), func(i,j int) {
		dams[i], dams[j] = dams[j], dams[i]
	})
	for _, dam := range dams {
		duties := findDuties(dam, dam.getDirectCauses())
		if len(duties) == 0 {
			continue
		}
		attempts := 0
		for {
			ui.newline()
			ui.println("Consider the following damage:")
			ui.println(dam.GetDescription())
			ui.println("Please enter the names of all people who could be held responsible for this:")
			ui.println("(Enter one name per line, finish with a . on a line of its own)")
			names := make([]string, 0)

			for {
				words, ret := ui.getInput()
				if ret {
					return ret
				}
				if len(words) == 0 {
					continue
				}
				if len(words) > 1 {
					ui.println("Please enter one word names only, finish with a . on a line of its own")
					continue
				}
				if words[0] == "." {
					break
				}
				names = append(names, words[0])
			}

			ui.println("You entered:")
			for _, name := range names {
				ui.println("- %s", name)
			}

			yes, ret := ui.yesNo("Is this correct")
			if ret {
				return ret
			}
			if !yes {
				ui.println("Ok, try again")
				continue
			}

			sort.Slice(names, func(i, j int) bool {
				return names[i] < names[j]
			})

			names2 := make([]string, 0)
			for p, _ := range collectPersonsFromDuties(duties) {
				names2 = append(names2, strings.ToLower(p.Name))
			}

			sort.Slice(names2, func(i, j int) bool {
				return names2[i] < names2[j]
			})

			if func() bool {
				if len(names) != len(names2) {
					return false
				}

				for i := range names {
					if names[i] != names2[i] {
						return false
					}
				}

				return true
			}() {
				ui.println("Correct!")
				state.registerAnswer(c, p, attempts == 0)
				return false
			}

			ui.println("Incorrect :-(")
			attempts++
		}
	}

	return false
}
