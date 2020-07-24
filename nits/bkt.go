package nits

// This file contains all the logic related to Bayesian Knowledge Tracing.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const (
	correct   = "1"  // Indicator for trainhmm.
	incorrect = "2"  // Indicator for trainhmm.
	separator = "~"  // Separator for the concepts in the input file for trainhmm.
	threshold = 0.95 // When do we consider a concept mastered?
	pInit     = 0.1  // Probability that concept was known a-priory.
	pLearn    = 0.2  // Probability that concept will transfer to mastered after a practice attempt.
	pSlip     = 0.3  // Probability that a mastered skill is applied incorrectly.
	pGuess    = 0.5  // Probability that an unmastered skill is applied incorrectly.
)

var trainhmmPath string // Path to the trainhmm binary

// --------------------------------------------------------------------

// try tries a possible hmmPath and returns an error if this is certainly not
// the trainhmm binary we want.
func try(hmmPath string) error {
	fi, err := os.Stat(hmmPath)
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeType != 0 {
		return errors.New(fmt.Sprintf("%s: illegal file type", hmmPath))
	}
	if fi.Mode()&0555 != 0555 {
		return errors.New(fmt.Sprintf("%s: illegal access mode (executable?)", hmmPath))
	}
	println("HMM training binary is in", hmmPath)
	trainhmmPath = hmmPath
	return nil
}

// initBKT initialized the Bayesian Knowledge Training module.
// Its most important job is to find the trainhmm binary.
func initBKT() {
	var err error
	errors := make([]error, 0)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("Can't determine the directory where the executable lives")
	}
	if err = try(path.Join(dir, fmt.Sprintf("trainhmm-%s", runtime.GOOS))); err == nil {
		return
	}
	errors = append(errors, err)
	if err = try(path.Join(mustUserHomeDir(), "standard-bkt", "trainhmm")); err == nil {
		return
	}
	errors = append(errors, err)
	for _, err := range errors {
		println(err.Error())
	}
	panic("Cannot find working trainhmm binary")
}

// --------------------------------------------------------------------

// answer is a struct that contains the information of an answered question.
type answer struct {
	questionShortName string // Here because of JSON unmarshaling.
	question          Question
	subQuestion       subQuestion
	correct           bool
}

// studentState contains, guess what!
type studentState struct {
	answers      []*answer
	burnt        map[Question]interface{} // Set of burnt questions.
	scores       map[*Concept]float64     // Knowledge scores per concept.
	content      *Content                 // Link to NITS content.
	nextQuestion Question                 // Allows the user to manually specify the next question.
}

// newStudentState creates a new student state object.
func newStudentState(content *Content) *studentState {
	state := &studentState{content: content}
	state.reset()
	return state
}

func (s *studentState) reset() {
	s.answers = make([]*answer, 0)
	s.burnt = make(map[Question]interface{})
	s.scores = make(map[*Concept]float64)
}

// registerAnswer registers a new answer in the student state.
func (s *studentState) registerAnswer(q Question, sq subQuestion, correct bool) {
	s.answers = append(s.answers, &answer{
		questionShortName: q.getShortName(),
		question:          q,
		subQuestion:       sq,
		correct:           correct})
	s.burn(q)
}

// burn burns a question. It will not be asked again.
func (s *studentState) burn(q Question) {
	s.burnt[q] = nil
}

// MarshalJson marshals an answer object to a JSON object.
func (a *answer) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["shortName"] = a.questionShortName
	m["correct"] = a.correct
	if a.subQuestion == nil {
		m["subQuestion"] = ""
	} else {
		m["subQuestion"] = a.subQuestion.getTag()
	}

	return json.Marshal(m)
}

// UnmarshalJSON unmarshals a JSON object back to an answer.
func (a *answer) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	if v, ok := m["correct"].(bool); ok {
		a.correct = v
	} else {
		return errors.New("data format error (correct)")
	}
	if v, ok := m["shortName"].(string); ok {
		// Since we do not have a link to the content here we store the
		// question shortname. The caller will have to resolve that back
		// to the interface.
		a.questionShortName = v
	} else {
		return errors.New("data format error (shortName)")
	}
	if v, ok := m["subQuestion"].(string); ok {
		if v == "" {
			a.subQuestion = nil
		} else if sq, ok := sqMap[v]; ok {
			a.subQuestion = sq
		}
	} else {
		return errors.New("data format error (subQuestion)")
	}

	return nil
}

// saveUserData saves the student state to ~/.nits_data. Only the
// registered answers are saved.
func (s *studentState) saveUserData() error {
	data, err := json.MarshalIndent(s.answers, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(mustUserHomeDir(), ".nits_data"), data, 0644)
}

// loadUserData loads the student state from ~/.nits_data. Only the
// registered answers are loaded. If a question/sub-question can not
// be found the question is discarded.
func (s *studentState) loadUserData() error {
	data, err := ioutil.ReadFile(path.Join(mustUserHomeDir(), ".nits_data"))
	if err != nil {
		return err
	}
	answers := make([]*answer, 0)
	s.burnt = make(map[Question]interface{})
	if err := json.Unmarshal(data, &answers); err != nil {
		return err
	}

	// We now need to find the questions in the content by the short name
	// that unMarshalJSON has put there.
	for _, a := range answers {
		// There is a chance that we are not finding the question if either
		// the student data has been manipulated (manual testing) or if the
		// question database has changed and a question has been removed.
		// Since we have already loaded the record we are going to keep it,
		// but when saving the student state we will drop it since there
		// is nothing we can do with it (the concepts have been lost).
		a.question = s.content.findQuestion(a.questionShortName)

		if a.question != nil {
			// If the question is a case and the sub question is nil then
			// drop this question because the sub question has apparently
			// been removed from the code.
			if c, ok := a.question.(*Case); ok && a.subQuestion == nil {
				if trace != nil {
					trace.println("Discarding %s because sq not found", c.ShortName)
				}
				a.question = nil
			}
		}
	}

	// Copy all the successfully loaded questions to the state.
	s.answers = make([]*answer, 0, len(answers))
	for _, a := range answers {
		if a.question != nil {
			s.answers = append(s.answers, a)
			s.burnt[a.question] = nil
		} else if trace != nil {
			trace.println("Discarding a question.")
		}
	}

	return nil
}

// --------------------------------------------------------------------

// writeTrainhmmInput writes the input file for the trainhmm binary.
// It returns the name of the temporary directory where the file was
// written. The format of the input file is described here:
// https://iedms.github.io/standard-bkt/.
func (s *studentState) writeTrainhmmInput() (string, error) {
	td, err := ioutil.TempDir("", "nits*")
	if err != nil {
		return td, err
	}
	var buffer bytes.Buffer

	// One line per registered answer.
	for _, a := range s.answers {
		columns := make([]string, 0)
		if a.correct {
			columns = append(columns, correct)
		} else {
			columns = append(columns, incorrect)
		}
		var tag string
		if a.subQuestion == nil {
			tag = a.questionShortName
		} else {
			tag = fmt.Sprintf("%s#%s", a.questionShortName, a.subQuestion.getTag())
		}
		columns = append(columns, "student", tag)
		names := make([]string, 0)
		for _, c := range a.question.getTrainingConcepts(a.subQuestion) {
			names = append(names, c.shortName)
		}
		columns = append(columns, strings.Join(names, separator))
		_, err := buffer.WriteString(strings.Join(columns, "\t"))
		if err != nil {
			return td, err
		}
		_, err = buffer.WriteRune('\n')
		if err != nil {
			return td, err
		}
	}

	err = ioutil.WriteFile(path.Join(td, "input"), buffer.Bytes(), 0644)
	return td, err
}

// runTrainhmm runs the trainhmm binary.
func (s *studentState) runTrainhmm(td string) error {
	cmd := exec.Command(
		trainhmmPath,
		"-p", "2",
		"-0", fmt.Sprintf("%f,1.0,%f,%f,%f", pInit, pLearn, 1-pSlip, pGuess),
		"-d", separator,
		path.Join(td, "input"),
		path.Join(td, "model.txt"),
		path.Join(td, "predict.txt"),
	)
	_, err := cmd.Output()
	return err
}

// readPrediction reads the skill scores from the prediction file
// generated by trainhmm.
func (s *studentState) readPrediction(td string) error {
	file, err := os.Open(path.Join(td, "predict.txt"))
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for _, a := range s.answers {
		if a.question == nil {
			continue
		}
		if !scanner.Scan() {
			return errors.New("unexpected end of predict.txt")
		}
		words := strings.Split(scanner.Text(), "\t")
		i := 2
		d, err := strconv.ParseFloat(words[0], 64)
		if err != nil {
			return err
		}
		// If this column contains 1.0 there is no next column, otherwise
		// the next column contains 1.0-d.
		if d == 1.0 {
			i = 1
		}
		for _, c := range a.question.getTrainingConcepts(a.subQuestion) {
			d, err := strconv.ParseFloat(words[i], 64)
			if err != nil {
				return err
			}
			s.scores[c] = d
		}
	}

	return scanner.Err()
}

// train runs the trainhmm binary on the questions answered by the
// student and reads back the skill scores into the state.
func (s *studentState) train() error {
	if len(s.answers) == 0 {
		return nil
	}
	td, err := s.writeTrainhmmInput()
	if err != nil {
		return err
	}
	err = s.runTrainhmm(td)
	if err != nil {
		return err
	}
	err = s.readPrediction(td)
	if err != nil {
		return err
	}
	return os.RemoveAll(td)
}

// --------------------------------------------------------------------

// avg returns the average skill for a particular question.
func (s *studentState) avg(q Question) float64 {
	total := 0.0
	concepts := q.getTrainingConcepts(nil)
	for _, c := range concepts {
		total += s.scores[c]
	}
	return total / float64(len(concepts))
}

// conceptsNotMastered filters all the concepts from a list that the
// student has not yet mastered and puts them in the output list.
func (s *studentState) conceptsNotMastered(concepts []*Concept) []*Concept {
	result := make([]*Concept, 0)
	for _, c := range concepts {
		if s.scores[c] < threshold {
			result = append(result, c)
		}
	}
	return result
}

// selectQuestion selects the next question that the student is going to
// answer. This method implements an algorithm that I call "race to
// mastery", finding the question that leads fastest to skills being
// mastered.
func (s *studentState) selectQuestion() Question {
	// If the debugger set a specific next question, then that is the
	// question we are going to return.
	if s.nextQuestion != nil {
		q := s.nextQuestion
		s.nextQuestion = nil
		return q
	}
	// Runs the training module.
	if err := s.train(); err != nil {
		panic(fmt.Sprintf("training error: %v", err))
	}
	// These are the possible questions.
	possibles := make([]Question, 0)
	// Goes through the set of questions in the content.
	for _, q := range s.content.Questions {
		// If we already answered the question (or it got burnt by the
		// debugger) then ignore this question.
		if _, ok := s.burnt[q]; ok {
			continue
		}
		concepts := q.getTrainingConcepts(nil)
		// If this question has no registered concepts that'd be weird, but
		// it is a possible question.
		if len(concepts) == 0 {
			possibles = append(possibles, q)
			continue
		}
		// If this question has any concepts in it that are not mastered
		// yet, it is a possible question. This step skips any questions
		// that has concepts that are all mastered.
		for _, c := range concepts {
			if s.scores[c] < threshold {
				possibles = append(possibles, q)
				break
			}
		}
	}
	// Did the student exhaust the content?
	if len(possibles) == 0 {
		return nil
	}
	// Sorts the possible questions so that the question with the highest
	// average skill score is first. This is the "race to mastery" step.
	sort.Slice(possibles, func(i, j int) bool {
		return s.avg(possibles[i]) > s.avg(possibles[j])
	})
	// If tracing is enabled, writes some tracing output.
	if trace != nil {
		for _, q := range possibles {
			trace.print("%32s", q.getShortName())
			for _, c := range q.getTrainingConcepts(nil) {
				trace.print("%20s(%f)", c.shortName, s.scores[c])
			}
			trace.println("  avg=%f", s.avg(q))
		}
	}
	return possibles[0]
}
