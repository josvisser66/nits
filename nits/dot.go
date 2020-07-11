package nits

import (
	"fmt"
	"io/ioutil"
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