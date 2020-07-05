package nits

import (
	"fmt"
	"strconv"
	"strings"
)

type ConceptBucket interface {
	GetAllConcepts() []*Concept
}

func findConcept(concepts []*Concept, name string) *Concept {
	for _, c := range concepts {
		if c.name == name {
			return c
		}
	}

	return nil
}

func (ui *userInterface) explainConcept(words []string, concepts []*Concept) {
	name := strings.Join(words[1:], " ")
	concept := findConcept(concepts, name)
	if concept == nil {
		ui.print(fmt.Sprintf("Concept '%s' not found."), true)
		ui.newline()
		return
	}

	ui.explain(concept.explanation)
}

func ExploreConcepts(ui *userInterface, b ConceptBucket) {
	concepts := b.GetAllConcepts()

	showConcepts := func([]string) bool {
		for i, c := range concepts {
			ui.print(fmt.Sprintf("%d: %s", i+1, c.name), true)
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
	ui.print("Here are the concepts that play a role in this question:", true)
	ui.newline()

	showConcepts(nil)
	help := func() {
		ui.print("Please enter the number of the concept you want to explore or 'done'.", true)
	}

	for {
		words, ret := ui.getInput()
		if ret {
			return
		}

		if len(words) > 1 {
			help()
			continue
		}

		n, err := strconv.Atoi(words[0])
		if err != nil {
			help()
			continue
		}

		if n < 1 || n > len(concepts) {
			help()
			continue
		}

		ui.explain(concepts[n-1].explanation)

		if len(concepts[n-1].related) > 0 {
			ui.print("Here are the related concepts:", true)
			ui.newline()
		}
	}
}
