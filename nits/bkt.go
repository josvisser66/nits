package nits

import (
	"fmt"
	"os"
)

var trainhmmPath string = "/Users/josv/standard-bkt/trainhmm"

func initBKT() {
	fi, err := os.Stat(trainhmmPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot stat: %s ", trainhmmPath))
	}
	if fi.Mode() & os.ModeType != 0 {
		panic(fmt.Sprintf("%s: illegal file type", trainhmmPath))
	}
	if fi.Mode() & 0555 != 0555 {
		panic(fmt.Sprintf("%s: illegal access mode (executable?)", trainhmmPath))
	}
}
