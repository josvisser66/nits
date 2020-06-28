package nits

import (
	"math/rand"
)

// --------------------------------------------------------------------
type Question interface {
	ask(ui *userInterface)
}

// --------------------------------------------------------------------
type Answer struct {
	Text        string
	Concepts    []*Concept
	Explanation Explanation
	Correct     bool
}

type MultipleChoiceQuestion struct {
	Question []string
	Concepts []*Concept
	Answers  []*Answer
}

func (q *MultipleChoiceQuestion) GetAllConcepts() []*Concept {
	m := make(map[*Concept]interface{})

	for _, c := range q.Concepts {
		m[c] = nil
	}

	for _, a := range q.Answers {
		for _, c := range a.Concepts {
			m[c] = nil
		}
	}

	keys := make([]*Concept, len(m))
	i := 0

	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func (q *MultipleChoiceQuestion) ask(ui *userInterface) {
	answers := make([]*Answer, len(q.Answers))
	copy(answers, q.Answers)
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	displayQuestion := func([]string) {
		ui.printParagraphs(q.Question)
		ui.newline()
		ui.printAnswers(answers)
		ui.newline()
	}
	displayQuestion(nil)
	ui.pushPrompt("Your answer? ")
	defer ui.popPrompt()

	ui.pushCommandContext(&CommandContext{
		"Answering a multiple choice question",
		[]*Command{
			{
				[]string{"again"},
				"Displays the question again.",
				displayQuestion,
			},
			{
				[]string{"explore"},
				"Explore the concepts involved in this question.",
				func([]string) {
					ExploreConcepts(ui, q)
				},
			},
		},
	})
	defer ui.popCommandContext()

	for {
		words := ui.getInput()
		if len(words) > 1 || len(words[0]) > 1 {
			ui.print("Please provide a one letter answer.", true)
			continue
		}
		answer := words[0][0] - 'a'
		if answers[answer].Correct {
			ui.print("Correct!", true)
			return
		}
		ui.print("Incorrect", true)
	}
}

// --------------------------------------------------------------------
type PropsQuestion struct {
	Concepts     []*Concept
	Propositions []string
	TrueOrFalse  []bool
}

// --------------------------------------------------------------------
type AnswerModel struct {
}

// --------------------------------------------------------------------
type OpenQuestion struct {
	Concepts    []*Concept
	Question    string
	AnswerModel AnswerModel
}

// --------------------------------------------------------------------
type TrueOrFalseQuestion struct {
	Concepts []*Concept
	Question string
	Answer   bool
}
