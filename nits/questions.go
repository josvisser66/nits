package nits

import "math/rand"

// --------------------------------------------------------------------
type Question interface {
	ask(ui *userInterface)
}

// --------------------------------------------------------------------
type MultipleChoiceQuestion struct {
	Concepts []*Concept
	Question string
	Answers []string
	CorrectAnswer int
}

func (q *MultipleChoiceQuestion) ask(ui *userInterface){
	ui.print(q.Question, true)
	ui.newline()

	answers := make([]string, len(q.Answers))
	correctAnswer := q.CorrectAnswer
	copy(answers, q.Answers)
	rand.Shuffle(len(answers), func(i,j int) {
		answers[i], answers[j] = answers[j], answers[i]
		if correctAnswer == i {
			correctAnswer = j
		} else if correctAnswer == j {
			correctAnswer = i
		}
	})

	ui.printAnswers(answers)
	ui.newline()
	ui.getAnswer()
}

// --------------------------------------------------------------------
type PropsQuestion struct {
	Concepts []*Concept
	Propositions []string
	TrueOrFalse []bool
}

// --------------------------------------------------------------------
type AnswerModel struct {
}

// --------------------------------------------------------------------
type OpenQuestion struct {
	Concepts []*Concept
	Question string
	AnswerModel AnswerModel
}

// --------------------------------------------------------------------
type TrueOrFalseQuestion struct {
	Concepts []*Concept
	Question string
	Answer bool
}
