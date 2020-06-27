package nits

import (
	"github.com/chzyer/readline"
	"io"
	"log"
	"os"
	"strings"
)

// --------------------------------------------------------------------
type userInterface struct {
	rl *readline.Instance
	promptStack []string
}

func newUserInterface() *userInterface {
	var err error

	ui := &userInterface{
		promptStack: make([]string, 0, 10),
	}
	ui.rl, err = readline.NewEx(&readline.Config{
		HistoryFile: "/tmp/readline.tmp",
		// AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}

	ui.pushPrompt("\033[31m»\033[0m ")
	log.SetOutput(ui.rl.Stderr())
	return ui
}

func (ui* userInterface) pushPrompt(t string) {
	ui.promptStack = append(ui.promptStack, t)
	ui.rl.SetPrompt(t)
}

func (ui* userInterface) popPrompt() {
	ui.promptStack = ui.promptStack[:len(ui.promptStack)-1]
	ui.rl.SetPrompt(ui.promptStack[len(ui.promptStack) - 1])
}

func (ui *userInterface) newline() {
	ui.rl.Terminal.PrintRune('\n')
}

func (ui *userInterface) print(t string, newline bool) {
	ui.printFromStartingWidth(t, 0, newline)
}

func (ui *userInterface) printFromStartingWidth(t string, w int, newline bool) {
	term := ui.rl.Terminal
	width := term.GetConfig().FuncGetWidth()
	words := strings.Split(t, " ")

	for _, word := range words {
		l := len(word)
		if l+w+1 > width {
			term.PrintRune('\n')
			w = 0
		}
		if w > 0 {
			term.PrintRune(' ')
		}
		ui.rl.Terminal.Print(word)
		w += l + 1
	}

	if newline {
		term.PrintRune('\n')
	}
}

func (ui *userInterface) printAnswers(answers []string) {
	r := 'A'

	for _, a := range answers {
		ui.print(string(r)+". ", false)
		ui.printFromStartingWidth(a, 3, true)
		r += 1
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (ui *userInterface) yesNo(question string) bool {
	ui.pushPrompt(question+" (Y/N)? ")
	defer ui.popPrompt()

	for {
		line, err := ui.rl.Readline()
		if err != nil {
			continue
		}
		switch strings.TrimSpace(strings.ToLower(line)) {
		case "y":
			fallthrough
		case "yes":
			return true
		case "n":
			fallthrough
		case "no":
			return false
		}
	}
}

func (ui *userInterface) exit() {
	if ui.yesNo("Are you sure you want to quit") {
		os.Exit(0)
	}
}

func (ui *userInterface) help() {
	ui.print("Help! I need somebody!", true)
}

func (ui *userInterface) getAnswer(displayQuestion func()) string {
	for {
		line, err := ui.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			ui.exit()
		}
		words := strings.Split(strings.ToLower(strings.TrimSpace(line)), " ")
		switch words[0] {
		case "quit":
			fallthrough
		case "exit":
			ui.exit()
		case "?":
			ui.help()
		case "again":
			displayQuestion()
		case "":
			continue
		default:
			return line
		}
	}
}
