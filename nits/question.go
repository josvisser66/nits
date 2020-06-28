package nits

import (
	"math/rand"
	"strings"
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
	Question string
	Concepts []*Concept
	Answers  []*Answer
}

func (q *MultipleChoiceQuestion) ask(ui *userInterface) {
	answers := make([]*Answer, len(q.Answers))
	copy(answers, q.Answers)
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	displayQuestion := func() {
		ui.print(q.Question, true)
		ui.newline()
		ui.printAnswers(answers)
		ui.newline()
	}
	displayQuestion()
	ui.pushPrompt("Your answer? ")
	defer ui.popPrompt()

	for {
		words := strings.Split(strings.ToLower(strings.TrimSpace(ui.getAnswer(displayQuestion))), " ")
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
