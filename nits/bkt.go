package nits

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
	"strconv"
	"strings"
)

const (
	correct   = "1"
	incorrect = "2"
	separator = "~"
)

var trainhmmPath = "/Users/josv/standard-bkt/trainhmm"

// --------------------------------------------------------------------
func initBKT() {
	fi, err := os.Stat(trainhmmPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot stat: %s ", trainhmmPath))
	}
	if fi.Mode()&os.ModeType != 0 {
		panic(fmt.Sprintf("%s: illegal file type", trainhmmPath))
	}
	if fi.Mode()&0555 != 0555 {
		panic(fmt.Sprintf("%s: illegal access mode (executable?)", trainhmmPath))
	}
}

// --------------------------------------------------------------------
type answer struct {
	questionShortName string // Here because of JSON unmarshaling.
	question          Question
	correct           bool
}

type studentState struct {
	answers []*answer
	burnt   map[Question]interface{} // A set, basically.
	scores  map[*Concept]float64
	content *Content
}

func newStudentState(content *Content) *studentState {
	return &studentState{
		answers: make([]*answer, 0),
		burnt:   make(map[Question]interface{}),
		scores:  make(map[*Concept]float64),
		content: content,
	}
}

func (s *studentState) registerAnswer(q Question, correct bool) {
	s.answers = append(s.answers, &answer{
		questionShortName: q.getShortName(),
		question:          q,
		correct:           correct})
	s.burnt[q] = nil
}

func (a *answer) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["shortName"] = a.questionShortName
	m["correct"] = a.correct

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

	return nil
}

func (s *studentState) selectQuestion() Question {
	if err := s.train(); err != nil {
		panic(fmt.Sprintf("training error: %v", err))
	}
	for _, q := range s.content.Questions {
		if _, ok := s.burnt[q]; !ok {
			return q
		}
	}

	return nil
}

// --------------------------------------------------------------------
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
		q := s.content.findQuestion(a.questionShortName)
		if q == nil {
			// There is a question in the answer list that is not in the
			// content.
			continue
		}
		if a.correct {
			columns = append(columns, correct)
		} else {
			columns = append(columns, incorrect)
		}
		columns = append(columns, "student", a.questionShortName)
		names := make([]string, 0)
		for _, c := range q.getTrainingConcepts() {
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
		for _, c := range a.question.getTrainingConcepts() {
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
