package nits

import (
	"fmt"
	"strconv"
	"strings"
)

var active bool

type ConceptBucket interface {
	getConcepts() []*Concept
}

func findConcept(concepts []*Concept, name string) *Concept {
	for _, c := range concepts {
		if c.name == name || c.shortName == name {
			return c
		}
	}

	return nil
}

func (ui *userInterface) explainConcept(words []string, concepts []*Concept) {
	name := strings.Join(words[1:], " ")
	concept := findConcept(concepts, name)
	if concept == nil {
		ui.error("Concept '%s' not found.", name)
		ui.newline()
		return
	}

	ui.explain(concept.explanation)
}

func ExploreConcepts(ui *userInterface, b ConceptBucket) {
	if active {
		return
	}
	active = true
	defer func() { active = false }()

	concepts := b.getConcepts()

	showConcepts := func([]string) bool {
		for i, c := range concepts {
			ui.println("%d: %s", i+1, c.name)
		}

		ui.newline()
		return false
	}

	ui.pushCommandContext(&CommandContext{
		Description: "Exploring a set of concepts",
		Commands: []*Command{
			{
				Aliases:  []string{"show"},
				Help:     "Shows all concepts we are exploring (again)",
				Executor: showConcepts,
			},
			{
				Aliases: []string{"explain"},
				Help:    "Explains a concept.",
				Executor: func(words []string) bool {
					ui.explainConcept(words, concepts)
					return false
				},
			},
			{
				Aliases:  []string{"done"},
				Help:     "Signals that you are done exploring concepts.",
				Executor: func([]string) bool { return true },
			},
		},
	})
	defer ui.popCommandContext()
	ui.pushPrompt("Concept explorer> ")
	defer ui.popPrompt()

	ui.newline()
	ui.println("Here are the concepts that play a role in this question:")
	ui.newline()

	showConcepts(nil)
	possibleAnswers := make(answerMap)

	for i:=1; i<=len(concepts); i++ {
		r := fmt.Sprint("%d", i)
		possibleAnswers[r] = []string{r}
	}

	for {
		answer, ret := ui.getAnswer(possibleAnswers)
		if ret {
			return
		}
		n, _ := strconv.Atoi(answer)
		ui.explain(concepts[n-1].explanation)
		ui.newline()
	}
}
