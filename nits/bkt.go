package nits

import (
	"fmt"
	"os"
)

const (
	correct   = 1
	incorrect = 2
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
	q       Question
	correct bool
}

var answers = make([]*answer, 0)

func registerAnswer(q Question, correct bool) {
	answers = append(answers, &answer{q, correct})
}
