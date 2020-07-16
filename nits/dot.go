package nits

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	shapePerson = "septagon"
	shapeEvent  = "ellipse"
	shapeDuty   = "pentagon"
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

	if err := dotEvents(f, pp.events); err != nil {
		return "", err
	}

	_, err = f.WriteString("}\n")

	return f.Name(), err
}

func dotDuties(f *os.File, duties map[*Duty]interface{}) error {
	for duty := range duties {
		if _, err := f.WriteString(fmt.Sprintf("  duty_%p [shape=%s,label=\"%s\"];\n", duty, shapeDuty, duty.Description)); err != nil {
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
		if _, err := f.WriteString(fmt.Sprintf("  person_%p [shape=%s,label=\"%s\"];\n", person, shapePerson, person.Name)); err != nil {
			return err
		}
	}
	return nil
}

func dotEvents(f *os.File, events map[*Event]interface{}) error {
	for event := range events {
		if _, err := f.WriteString(fmt.Sprintf("  event_%p [shape=%s,label=\"%s\"];\n", event, shapeEvent, event.Description)); err != nil {
			return err
		}
	}

	for event := range events {
		for _, consequence := range event.Consequences {
			if _, err := f.WriteString(fmt.Sprintf("  event_%p -> event_%p;\n", event, consequence)); err != nil {
				return err
			}
		}
	}

	return nil
}
