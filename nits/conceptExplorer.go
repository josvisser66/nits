package nits

import (
	"strconv"
	"strings"
)

var active bool

type ConceptBucket interface {
	getAllConcepts() []*Concept
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
		ui.println("Concept '%s' not found.", name)
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

	concepts := b.getAllConcepts()

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
	help := func() {
		ui.println("Please enter the number of the concept you want to explore or 'done'.")
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
		ui.newline()
	}
}
