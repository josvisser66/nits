package nits

// This file contains the concept explorer.

import (
	"fmt"
	"strconv"
	"strings"
)

// ConceptBucket is a thing that has concepts.
type ConceptBucket interface {
	getConcepts() []*Concept
}

// findConcept finds a particular concept by name from a slice of concepts.
// The match can be either on name or short name.
func findConcept(concepts []*Concept, name string) *Concept {
	for _, c := range concepts {
		if c.name == name || c.shortName == name {
			return c
		}
	}

	return nil
}

// explainConcept implements the explain command of the concept explorer.
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

// exploreConcepts is the concept explorer.
func exploreConcepts(ui *userInterface, b ConceptBucket) {
	concepts := b.getConcepts()

	printConceptList := func([]string) bool {
		for i, c := range concepts {
			ui.println("%d: %s", i+1, c.name)
		}

		ui.newline()
		return false
	}

	ui.pushCommandContext(&CommandContext{
		description: "Exploring a set of concepts",
		commands: []*Command{
			{
				aliases:  []string{"show", "print", "display"},
				help:     "Shows all concepts we are exploring (again).",
				executor: printConceptList,
			},
			{
				aliases: []string{"explain"},
				help:    "Explains a concept by name.",
				executor: func(words []string) bool {
					ui.explainConcept(words, concepts)
					return false
				},
			},
			{
				aliases:  []string{"done"},
				help:     "Signals that you are done exploring concepts.",
				executor: func([]string) bool { return true },
			},
		},
	})
	defer ui.popCommandContext()
	ui.pushPrompt("Concept explorer> ")
	defer ui.popPrompt()

	ui.newline()
	ui.println("Here are the concepts that play a role in this question:")
	ui.newline()

	printConceptList(nil)
	// Creates an answerMap that allows selecting the concepts to be explained
	// by number.
	possibleAnswers := make(answerMap)

	for i:=1; i<=len(concepts); i++ {
		r := fmt.Sprintf("%d", i)
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
