package nits

import (
	"fmt"
	"strings"
)

type ConceptBucket interface {
	GetAllConcepts() []*Concept
}

func findConcept(concepts []*Concept, name string) *Concept {
	for _, c := range concepts {
		if c.Name == name {
			return c
		}
	}

	return nil
}

func (ui *userInterface) explainConcept(words[] string, concepts []*Concept) {
	name := strings.Join(words[1:], " ")
	concept := findConcept(concepts, name)
	if concept == nil {
		ui.print(fmt.Sprintf("Concept '%s' not found."), true)
		ui.newline()
		return
	}

	ui.explain(concept.Explanation)
}

func ExploreConcepts(ui *userInterface, b ConceptBucket) {
	concepts := b.GetAllConcepts()

	showConcepts := func([]string) {
		for i, c := range concepts {
			ui.print(fmt.Sprintf("%d: %s", i+1, c.Name), true)
		}

		ui.newline()
	}

	ui.pushCommandContext(&CommandContext{
		Description: "Exploring a set of concepts",
		Commands: []*Command {
			{
				Aliases:  []string{"show"},
				Help:     "Shows all concepts we are exploring (again)",
				Executor: showConcepts,
			},
			{
				Aliases:  []string{"explain"},
				Help:     "Explains a concept.",
				Executor: func(words []string) {
					ui.explainConcept(words, concepts)
				},
			},
			{
				Aliases:  []string{"done"},
				Help:     "Signals that you are done exploring concepts.",
				Executor: func([]string){},
			},
		},
	})
	defer ui.popCommandContext()
	ui.pushPrompt("Concept explorer> ")
	defer ui.popPrompt()

	ui.newline()
	ui.print("Here are the concepts that play a role in this question:", true)
	ui.newline()

	showConcepts(nil)

	for {
		ui.getInput()
	}
}

