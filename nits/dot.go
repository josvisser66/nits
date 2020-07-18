package nits

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	shapePerson         = "diamond"
	shapeClaim          = "trapezium"
	shapeEvent          = "ellipse"
	shapeAct            = "box"
	shapeDuty           = "hexagon"
	shapeInjuryOrDamage = "house"
	shapeNegPerSe       = "parallelogram"
)

func makeDot(state *studentState, withSkills bool) (string, error) {
	f, err := ioutil.TempFile("", "nits*.dot")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString("digraph nits {\n")
	if err != nil {
		return f.Name(), err
	}

	for _, c := range allConcepts {
		if withSkills {
			_, err = f.WriteString(fmt.Sprintf("\t%s [label=\"%s (%f)\"];\n", c.shortName, c.shortName, state.scores[c]))
		} else {
			_, err = f.WriteString(fmt.Sprintf("\t%s;\n", c.shortName))
		}
		for _, rc := range c.related {
			_, err = f.WriteString(fmt.Sprintf("\t%s -> %s;\n", c.shortName, rc.shortName))
		}
	}

	for _, q := range state.content.Questions {
		_, err = f.WriteString(fmt.Sprintf("\t%s [shape=box];\n", q.getShortName()))
		for _, rc := range q.getConcepts() {
			_, err = f.WriteString(fmt.Sprintf("\t%s -> %s;\n", q.getShortName(), rc.shortName))
		}
	}

	_, err = f.WriteString("}\n")

	return f.Name(), err
}

// --------------------------------------------------------------------
func makeCaseDot(state *studentState, words []string) (string, error) {
	if len(words) < 2 {
		return "", errors.New("please provide the short name of a case as an argument")
	}
	q := state.content.findQuestion(words[1])
	if q == nil {
		return "", errors.New("no case with that name")
	}
	c, ok := q.(*Case)
	if !ok {
		return "", errors.New("that question is not a case")
	}
	pp := preprocess(c)
	f, err := ioutil.TempFile("", "nits*.dot")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString("digraph nits {\n")
	if err != nil {
		return f.Name(), err
	}

	if err := dotPersons(f, pp.persons); err != nil {
		return "", err
	}

	if err := dotDuties(f, pp.duties); err != nil {
		return "", err
	}

	if err := dotInjuryOrDamages(f, pp.injuriesOrDamages); err != nil {
		return "", err
	}

	if err := dotEvents(f, pp.events); err != nil {
		return "", err
	}

	_, err = f.WriteString("}\n")

	return f.Name(), err
}

type hasLabel interface {
	getLabel() string
}

func draw(f *os.File, t string, shape string, obj hasLabel) error {
	_, err := f.WriteString(fmt.Sprintf("  %s_%p [shape=%s,label=\"%s\"];\n", t, obj, shape, obj.getLabel()))
	return err
}

func dotInjuryOrDamages(f *os.File, id map[InjuryOrDamage]interface{}) error {
	for dam := range id {
		if err := draw(f, "dam", shapeInjuryOrDamage, dam); err != nil {
			return err
		}
		for _, person := range dam.GetPersons() {
			if _, err := f.WriteString(fmt.Sprintf("  person_%p -> dam_%p [style=dotted];\n", person, dam)); err != nil {
				return nil
			}
		}
	}
	return nil
}

func dotDuties(f *os.File, duties map[*Duty]interface{}) error {
	for duty := range duties {
		if err := draw(f, "duty", shapeDuty, duty); err != nil {
			return err
		}
		for _, person := range duty.OwedFrom {
			if _, err := f.WriteString(fmt.Sprintf("  person_%p -> duty_%p;\n", person, duty)); err != nil {
				return nil
			}
		}
		for _, person := range duty.OwedTo {
			if _, err := f.WriteString(fmt.Sprintf("  duty_%p -> person_%p;\n", duty, person)); err != nil {
				return nil
			}
		}
	}
	return nil
}

func dotPersons(f *os.File, persons map[*Person]interface{}) error {
	for person := range persons {
		if err := draw(f, "person", shapePerson, person); err != nil {
			return err
		}
	}
	return nil
}

func dotClaims(f *os.File, claims map[*Claim]interface{}) error {
	for claim := range claims {
		if err := draw(f, "claim", shapeClaim, claim); err != nil {
			return err
		}
		if _, err := f.WriteString(fmt.Sprintf("  person_%p -> claim_%p [style=dotted];\n", claim.Person, claim)); err != nil {
			return err
		}
	}
	return nil
}

func getEventShape(e Event) string {
	if _, ok := e.(*Act); ok {
		return shapeAct
	}
	return shapeEvent
}

func dotEvents(f *os.File, events map[Event]interface{}) error {
	for event := range events {
		if err := draw(f, "event", getEventShape(event), event); err != nil {
			return err
		}
	}

	for event := range events {
		for _, consequence := range event.getConsequences() {
			if _, err := f.WriteString(fmt.Sprintf("  event_%p -> event_%p;\n", event, consequence)); err != nil {
				return err
			}
		}
		if event.getDuty() != nil {
			if _, err := f.WriteString(fmt.Sprintf("  duty_%p -> event_%p [style=dotted];\n", event.getDuty(), event)); err != nil {
				return err
			}
		}
		if event.getClaims() != nil {
			for _, claim := range event.getClaims() {
				if _, err := f.WriteString(fmt.Sprintf("  claim%p -> event_%p [style=dotted];\n", claim, event)); err != nil {
					return err
				}
			}
		}
		if event.getNegPerSe() != nil {
			if err := draw(f, "negperse", shapeNegPerSe, event.getNegPerSe()); err != nil {
				return err
			}
			if _, err := f.WriteString(fmt.Sprintf("  negperse_%p -> event_%p [style=dotted];\n", event.getNegPerSe(), event)); err != nil {
				return err
			}
			for _, person := range event.getNegPerSe().Persons {
				if _, err := f.WriteString(fmt.Sprintf("  person_%p -> negperse_%p [style=dotted];\n", person, event.getNegPerSe())); err != nil {
					return err
				}
			}
		}
		for _, dam := range event.getInjuriesOrDamages() {
			if _, err := f.WriteString(fmt.Sprintf("  event_%p -> dam_%p;\n", event, dam)); err != nil {
				return err
			}
		}
	}

	return nil
}
