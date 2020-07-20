package nits

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
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
	correct   = "1"
	incorrect = "2"
	separator = "~"
	threshold = 0.95 // When do we consider a concept mastered?
	pInit     = 0.1  // Probability that concept was known a-priory.
	pLearn    = 0.2  // Probability that concept will transfer to mastered after a practice attempt.
	pSlip     = 0.1  // Probability that a mastered skill is applied incorrectly.
	pGuess    = 0.5  // Probability that an unmastered skill is applied incorrectly.
)

var trainhmmPath string

// --------------------------------------------------------------------
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
	if err = try(path.Join(getUserHomeDir(), "standard-bkt", "trainhmm")); err == nil {
		return
	}
	errors = append(errors, err)
	for _, err := range errors {
		println(err.Error())
	}
	panic("Cannot find working trainhmm binary")
}

// --------------------------------------------------------------------
type answer struct {
	questionShortName string // Here because of JSON unmarshaling.
	question          Question
	subQuestion       subQuestion
	correct           bool
}

type studentState struct {
	answers      []*answer
	burnt        map[Question]interface{} // A set, basically.
	scores       map[*Concept]float64
	content      *Content
	nextQuestion Question // for debugger, not used otherwise.
}

func newStudentState(content *Content) *studentState {
	return &studentState{
		answers: make([]*answer, 0),
		burnt:   make(map[Question]interface{}),
		scores:  make(map[*Concept]float64),
		content: content,
	}
}

func (s *studentState) registerAnswer(q Question, sq subQuestion, correct bool) {
	s.answers = append(s.answers, &answer{
		questionShortName: q.getShortName(),
		question:          q,
		subQuestion:       sq,
		correct:           correct})
	s.burn(q)
}

func (s *studentState) burn(q Question) {
	s.burnt[q] = nil
}

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
		a.questionShortName = v
	} else {
		return errors.New("data format error (shortName)")
	}
	if v, ok := m["subQuestion"].(string); ok {
		if v == "" {
			a.subQuestion = nil
		} else {
			sq, ok := sqMap[v]
			CHECK(ok, "subquestion %s not found", v)
			a.subQuestion = sq
		}
	} else {
		return errors.New("data format error (subQuestion)")
	}

	return nil
}

func (s *studentState) saveUserData() error {
	data, err := json.MarshalIndent(s.answers, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(getUserHomeDir(), ".nits_data"), data, 0644)
}

func (s *studentState) loadUserData() error {
	data, err := ioutil.ReadFile(path.Join(getUserHomeDir(), ".nits_data"))
	if err != nil {
		return err
	}
	s.answers = make([]*answer, 0)
	s.burnt = make(map[Question]interface{})
	if err := json.Unmarshal(data, &s.answers); err != nil {
		return err
	}

	for _, a := range s.answers {
		// There is a chance that we are not finding the question if either
		// the student data has been manipulated (manual testing) or if the
		// question database has changed and a question has been removed.
		// Since we have already loaded the record we are going to keep it,
		// but when saving the student state we might drop it since there
		// is nothing we can do with it (the concepts have been lost).
		a.question = s.content.findQuestion(a.questionShortName)
		if a.question != nil {
			s.burnt[a.question] = nil
		}
	}

	return nil
}

// --------------------------------------------------------------------
func (s *studentState) writeTrainhmmInput() (string, error) {
	td, err := ioutil.TempDir("", "nits*")
	if err != nil {
		return td, err
	}
	var buffer bytes.Buffer

	for _, a := range s.answers {
		columns := make([]string, 0)
		if a.question == nil {
			// Question was removed from the question database.
			continue
		}
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
func (s *studentState) avg(q Question) float64 {
	total := 0.0
	concepts := q.getTrainingConcepts(nil)
	for _, c := range concepts {
		total += s.scores[c]
	}
	return total / float64(len(concepts))
}

func (s *studentState) conceptsNotMastered(concepts []*Concept) []*Concept {
	result := make([]*Concept, 0)
	for _, c := range concepts {
		if s.scores[c] < threshold {
			result = append(result, c)
		}
	}
	return result
}

func (s *studentState) selectQuestion() Question {
	if s.nextQuestion != nil {
		q := s.nextQuestion
		s.nextQuestion = nil
		return q
	}
	if err := s.train(); err != nil {
		panic(fmt.Sprintf("training error: %v", err))
	}
	possibles := make([]Question, 0)
	for _, q := range s.content.Questions {
		if _, ok := s.burnt[q]; ok {
			continue
		}
		concepts := q.getTrainingConcepts(nil)
		if len(concepts) == 0 {
			possibles = append(possibles, q)
			continue
		}
		for _, c := range concepts {
			if s.scores[c] < threshold {
				possibles = append(possibles, q)
				break
			}
		}
	}
	if len(possibles) == 0 {
		return nil
	}
	rand.Shuffle(len(possibles), func(i, j int) {
		possibles[i], possibles[j] = possibles[j], possibles[i]
	})
	sort.Slice(possibles, func(i, j int) bool {
		return s.avg(possibles[i]) > s.avg(possibles[j])
	})
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
