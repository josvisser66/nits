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
	questionShortName string
	correct           bool
}

var answers = make([]*answer, 0)

func registerAnswer(q Question, correct bool) {
	answers = append(answers, &answer{q.getShortName(), correct})
	burn(q.getShortName())
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
		burn(v)
	} else {
		return errors.New("data format error (shortName)")
	}

	return nil
}

// --------------------------------------------------------------------
var burnt = make(map[string]interface{})

func burn(shortName string) {
	burnt[shortName] = nil
}

func selectQuestion(content *Content) Question {
	for _, q := range content.Questions {
		if _, ok := burnt[q.getShortName()]; !ok {
			return q
		}
	}

	return nil
}

// --------------------------------------------------------------------
func saveUserData() error {
	data, err := json.Marshal(answers)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(getUserHomeDir(), ".nits_data"), data, 0644)
}

func loadUserData() error {
	data, err := ioutil.ReadFile(path.Join(getUserHomeDir(), ".nits_data"))
	if err != nil {
		return err
	}
	answers = make([]*answer, 0)
	burnt = make(map[string]interface{})
	return json.Unmarshal(data, &answers)
}

// --------------------------------------------------------------------
func writeTrainhmmInput(content *Content) (string, error) {
	td, err := ioutil.TempDir("", "nits*")
	if err != nil {
		return td, err
	}
	var buffer bytes.Buffer

	for _, a := range answers {
		q := content.findQuestion(a.questionShortName)
		if q == nil {
			continue
		}
		if a.correct {
			buffer.WriteString(correct)
		} else {
			buffer.WriteString(incorrect)
		}
		buffer.WriteString(fmt.Sprintf("\tstudent\t%s\t", a.questionShortName))
		first := true
		for _, c := range q.getConcepts() {
			if !first {
				buffer.WriteString(separator)
			} else {
				first = false
			}
			buffer.WriteString(c.shortName)
		}
	}

	err = ioutil.WriteFile(path.Join(td, "input"), buffer.Bytes(), 0644)
	return td, err
}

func runTrainhmm(td string) error {
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

func readPrediction(td string, content *Content, ui *userInterface) error {
	file, err := os.Open(path.Join(td, "predict.txt"))
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for _, a := range answers {
		q := content.findQuestion(a.questionShortName)
		if q == nil {
			continue
		}
		if !scanner.Scan() {
			return errors.New("unexpected end of predict.txt")
		}
		words := strings.Split(scanner.Text(), "\t")
		i := 2
		ui.println("%s:", q.getShortName())
		for _, c := range q.getConcepts() {
			d, err := strconv.ParseFloat(words[i], 64)
			if err != nil {
				return err
			}
			ui.println("  %s: %l", c.shortName, d)
		}
	}

	return scanner.Err()
}
