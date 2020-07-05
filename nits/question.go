package nits

import (
	"math/rand"
)

// --------------------------------------------------------------------
type Question interface {
	ask(ui *userInterface)
	getConcepts() []*Concept
	getAllConcepts() []*Concept
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

func (q *MultipleChoiceQuestion) getConcepts() []*Concept {
	return q.Concepts
}

// getAllConcepts returns all the concepts that are directly implicated
// in the question or its answers.
func (q *MultipleChoiceQuestion) getAllConcepts() []*Concept {
	m := make(map[*Concept]interface{})

	for _, c := range q.Concepts {
		m[c] = nil
	}

	for _, a := range q.Answers {
		for _, c := range a.Concepts {
			m[c] = nil
		}
	}

	return conceptMapToSlice(m)
}

func pushCommandContext(ui *userInterface, q Question, displayQuestion func([]string) bool) {
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
				func([]string) bool {
					ExploreConcepts(ui, q)
					return false
				},
			},
		},
	})
}

func (q *MultipleChoiceQuestion) ask(ui *userInterface) {
	answers := make([]*Answer, len(q.Answers))
	copy(answers, q.Answers)
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	displayQuestion := func([]string) bool {
		ui.printParagraphs(q.Question)
		ui.newline()
		ui.printAnswers(answers)
		ui.newline()
		return false
	}

	displayQuestion(nil)
	ui.pushPrompt("Your answer? ")
	pushCommandContext(ui, q, displayQuestion)
	defer ui.popPrompt()
	defer ui.popCommandContext()

	attempts := 0

	for {
		words, ret := ui.getInput()
		if ret {
			return
		}
		if len(words) > 1 || len(words[0]) > 1 {
			ui.println("Please provide a one letter answer.")
			continue
		}
		answer := words[0][0] - 'a'
		if int(answer) >= len(answers) {
			ui.println("Please enter an answer from A to %c.", rune('A'+len(answers)-1))
			continue
		}
		if answers[answer].Correct {
			ui.println("Correct :-)")
			registerAnswer(q, attempts == 0)
			return
		}
		ui.println("Incorrect :-(")
		attempts++
	}
}

// --------------------------------------------------------------------
type Proposition struct {
	Proposition string
	Concepts    []*Concept
	True        bool
}

type PropsQuestion struct {
	Propositions []*Proposition
}

func (q *PropsQuestion) getConcepts() []*Concept {
	m := make(map[*Concept]interface{}, 0)

	for _, prop := range q.Propositions {
		for _, c := range prop.Concepts {
			m[c] = nil
		}
	}

	return conceptMapToSlice(m)
}

// https://codereview.stackexchange.com/questions/202352/int-to-roman-numerals-in-go-golang
func romanNumeral(number int) string {
	conversions := []struct {
		value int
		digit string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	roman := ""
	for _, conversion := range conversions {
		for number >= conversion.value {
			roman += conversion.digit
			number -= conversion.value
		}
	}
	return roman
}

func (q *PropsQuestion) ask(ui *userInterface) {
	displayQuestion := func([]string) bool {
		ui.println("Consider the following propositions:")
		ui.newline()

		for i, prop := range q.Propositions {
			ui.println("%4s. %s", romanNumeral(i+1), prop.Proposition)
		}

		ui.newline()
		n := len(q.Propositions) << 1
		r := 'A'

		for i := 0; i < n; i++ {
			ui.print("%c) ", r)
			r++
			k := i

			for j := range q.Propositions {
				if j > 0 {
					ui.print(", ")
				}
				var s string
				if k%2 == 1 {
					s = "true"
				} else {
					s = "false"
				}
				ui.print("%s is %s", romanNumeral(j+1), s)
				k >>= 1
			}

			ui.println(".")
		}

		ui.newline()
		return false
	}

	displayQuestion(nil)
	ui.pushPrompt("Your answer? ")
	pushCommandContext(ui, q, displayQuestion)
	defer ui.popPrompt()
	defer ui.popCommandContext()
	attempts := 0

outer:
	for {
		words, ret := ui.getInput()
		if ret {
			return
		}
		if len(words) > 1 || len(words[0]) > 1 {
			ui.println("Please provide a one letter answer.")
			continue
		}
		answer := words[0][0] - 'a'
		if int(answer) >= len(q.Propositions)<<1 {
			ui.println("Please enter an answer from A to %c", rune('A'+len(q.Propositions)<<1-1))
			continue
		}
		for _, prop := range q.Propositions {
			if (answer%2 == 0 && prop.True) || (answer%2 == 1 && !prop.True) {
				ui.println("Incorrect :-(")
				attempts++
				continue outer
			}
			answer >>= 1
		}
		ui.println("Correct :-)")
		registerAnswer(q, attempts == 0)
		return
	}
}

// getAllConcepts returns all the concepts that are directly implicated
// in the question or its answers.
func (q *PropsQuestion) getAllConcepts() []*Concept {
	return q.getConcepts()
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
