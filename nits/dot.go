package nits

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
type caseDotState struct {
	events map[*Event]interface{}
}

func (s *caseDotState) dotEvent(f *os.File, e *Event) error {
	if _, ok := s.events[e]; !ok {
		_, err := f.WriteString(fmt.Sprintf("event_%p [label=\"%s\"];\n", e, e.Description))
		if err != nil {
			return err
		}
		s.events[e] = nil
		s.dotEvents(f, e.Consequences)

		for _, e2 := range e.Consequences {
			_, err := f.WriteString(fmt.Sprintf("  event_%p -> event_%p;\n", e, e2))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *caseDotState) dotEvents(f *os.File, events []*Event) error {
	for _, e := range events {
		if err := s.dotEvent(f, e); err != nil {
			return err
		}
	}
	return nil
}

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
	f, err := ioutil.TempFile("", "nits*.dot")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString("digraph nits {\n")
	if err != nil {
		return f.Name(), err
	}

	s := &caseDotState{map[*Event]interface{}{}}
	s.dotEvents(f, c.RootEvents)

	_, err = f.WriteString("}\n")

	return f.Name(), err
}
