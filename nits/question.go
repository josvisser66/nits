package nits

// This file contains the core abstraction for asking questions and implements
// multiple choice questions and proposition questions.

import (
	"math/rand"
)

// --------------------------------------------------------------------

// Question is a question. It can be a case, a multiple choice question, or
// a proposition question.
type Question interface {
	getShortName() string
	getConcepts() []*Concept
	getTrainingConcepts(sq subQuestion) []*Concept
	check()
	ask(ui *userInterface, state *studentState)
}

// Answer is an answer for a multiple choice question.
type Answer struct {
	Text           string
	Concepts       []*Concept
	Explanation    Explanation
	Correct        bool
	NoneOfTheAbove bool
}

// MultipleChoiceQuestion is exactly what you think it is.
type MultipleChoiceQuestion struct {
	ShortName string
	Question  []string
	Concepts  []*Concept
	Answers   []*Answer
}

// check checks for the correctness of a multiple choice question.
func (q *MultipleChoiceQuestion) check() {
	CHECK(len(q.Answers) >= 2, "Question %s does not have at least two answers", q.ShortName)
	CHECK(len(q.getConcepts()) > 0, "Question %s does not have any concepts!", q.ShortName)

	n := 0

	for _, a := range q.Answers {
		if a.Correct {
			n++
		}
	}

	CHECK(n > 0, "Question %s does not have any correct answers!", q.ShortName)
}

func (q *MultipleChoiceQuestion) getShortName() string {
	return q.ShortName
}

// getConcepts collects all the concepts involved in this question and all of
// the answers.
func (q *MultipleChoiceQuestion) getConcepts() []*Concept {
	m := make(map[*Concept]interface{})

	for _, c := range q.Concepts {
		m[c] = nil
	}

	for _, a := range q.Answers {
		for _, c := range a.Concepts {
			m[c] = nil
		}
	}

	return conceptSetToSlice(m)
}

// getTrainingConcepts collects all the training concepts for this
// question. These are the concepts related to the question *and* the
// correct answer (only). Concepts involved in incorrect answers are not
// returned.
func (q *MultipleChoiceQuestion) getTrainingConcepts(sq subQuestion) []*Concept {
	CHECK(sq == nil, "unexpected subQuestion for MultipleChoiceQuestion")
	m := make(map[*Concept]interface{})

	for _, c := range q.Concepts {
		m[c] = nil
	}

	for _, a := range q.Answers {
		if a.Correct {
			for _, c := range a.Concepts {
				m[c] = nil
			}
			break
		}
	}

	return conceptSetToSlice(m)
}

// pushCommandContext pushes a command context for general use when
// answering questions.
func pushCommandContext(name string, state *studentState, ui *userInterface, q Question, displayQuestion func([]string) bool) {
	ui.pushCommandContext(&CommandContext{
		description: name,
		commands: []*Command{
			{
				aliases:  []string{"again"},
				global:   true,
				help:     "Displays the question again.",
				executor: displayQuestion,
			},
			{
				aliases:  []string{"abandon"},
				global:   true,
				help:     "Abandons this question.",
				executor: func([]string) bool {
					return true
				},
			},
			{
				aliases: []string{"explore"},
				help:    "Explore the concepts involved in this question.",
				executor: func([]string) bool {
					exploreConcepts(ui, q)
					return false
				},
			},
			{
				aliases: []string{"burn"},
				help:    "Burn this question.",
				executor: func([]string) bool {
					state.burn(q)
					return true
				},
			},
		},
	})
}

// makeAnswerMap makes an answerMap with valid answers a .. (n-1 letters
// further).
func makeAnswerMap(n int) answerMap {
	m := make(answerMap)

	for i:=1; i <= n; i++ {
		r := string('a' + i - 1)
		m[r] = []string{r}
	}

	return m
}

// ask asks a multiple choice question.
func (q *MultipleChoiceQuestion) ask(ui *userInterface, state *studentState) {
	// Shuffles the answers, making sure that None Of The Above is always last.
	answers := make([]*Answer, 0, len(q.Answers))
	var noneOfTheAbove *Answer
	for _, a := range q.Answers {
		if a.NoneOfTheAbove {
			noneOfTheAbove = a
		} else {
			answers = append(answers, a)
		}
	}
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	if noneOfTheAbove != nil {
		answers = append(answers, noneOfTheAbove)
	}
	displayQuestion := func([]string) bool {
		ui.newline()
		ui.printParagraphs(q.Question)
		ui.newline()
		ui.printAnswers(answers)
		ui.newline()
		return false
	}

	displayQuestion(nil)
	ui.pushPrompt("Your answer? ")
	pushCommandContext("Answering a multiple choice question", state, ui, q, displayQuestion)
	defer ui.popPrompt()
	defer ui.popCommandContext()

	attempts := 0
	possibleAnswers := makeAnswerMap(len(answers))

	for {
		s, ret := ui.getAnswer(possibleAnswers)
		if ret {
			return
		}
		answer := s[0] - 'a'
		if answers[answer].Correct {
			ui.println("Correct :-)")
			state.registerAnswer(q, nil, attempts == 0)
			return
		}
		ui.println("Incorrect :-(")
		attempts++
	}
}

// --------------------------------------------------------------------

// Proposition is a proposition that can be true or false.
type Proposition struct {
	Proposition string
	Concepts    []*Concept
	True        bool
}

// PropsQuestion is a question that asks a bunch of propositions.
type PropsQuestion struct {
	ShortName    string
	Propositions []*Proposition
}

func (q *PropsQuestion) getShortName() string {
	return q.ShortName
}

// check checks the validity of a proposition question.
func (q *PropsQuestion) check() {
	CHECK(q.ShortName != "", "Proposition question does not have a short name (%s)", q.Propositions[0].Proposition)
	CHECK(len(q.getConcepts()) > 0, "Question %s does not have any concepts!", q.ShortName)
	CHECK(len(q.Propositions) >= 2, "Question %s does not have at least 2 propositions", q.ShortName)
}

// getConcepts returns all the concepts involved in all the propositions
// in this proposition question.
func (q *PropsQuestion) getConcepts() []*Concept {
	m := make(map[*Concept]interface{}, 0)

	for _, prop := range q.Propositions {
		for _, c := range prop.Concepts {
			m[c] = nil
		}
	}

	return conceptSetToSlice(m)
}

// getTrainingConcepts is the same as getConcepts (for a proposition question).
func (q *PropsQuestion) getTrainingConcepts(sq subQuestion) []*Concept {
	CHECK(sq == nil, "unexpected subQuestion for PropsQuestion")
	return q.getConcepts()
}

// Convert a number to a roman numeral. Courtesy of:
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

// ask asks a proposition question.
func (q *PropsQuestion) ask(ui *userInterface, state *studentState) {
	displayQuestion := func([]string) bool {
		ui.newline()
		ui.println("Consider the following propositions:")
		ui.newline()

		for i, prop := range q.Propositions {
			ui.println("%4s. %s", romanNumeral(i+1), prop.Proposition)
		}

		ui.newline()

		// Now we generate the answers of the type: I is true, II is false.
		// Since each proposition can be either true or false this is
		// processed as a bitmap, going through all the permutations of a
		// binary number 2^(number of propositions).
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
	pushCommandContext("Answering a proposition question", state, ui, q, displayQuestion)
	defer ui.popPrompt()
	defer ui.popCommandContext()

	attempts := 0
	possibleAnswers := makeAnswerMap(1 << len(q.Propositions))

outer:
	for {
		s, ret := ui.getAnswer(possibleAnswers)
		if ret {
			return
		}
		answer := s[0] - 'a'
		// Check the bitmap implied in the answer and see if the student
		// got each proposition right.
		for _, prop := range q.Propositions {
			if (answer%2 == 0 && prop.True) || (answer%2 == 1 && !prop.True) {
				ui.println("Incorrect :-(")
				attempts++
				continue outer
			}
			answer >>= 1
		}
		ui.println("Correct :-)")
		state.registerAnswer(q, nil, attempts == 0)
		return
	}
}

