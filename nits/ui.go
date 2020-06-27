package nits

import (
	"github.com/chzyer/readline"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

// --------------------------------------------------------------------
type userInterface struct {
	rl *readline.Instance
}

func newUserInterface() *userInterface {
	var err error

	ui := &userInterface{}
	ui.rl, err = readline.NewEx(&readline.Config{
		Prompt:      "\033[31m»\033[0m ",
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

	log.SetOutput(ui.rl.Stderr())
	return ui
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

func (ui *userInterface) yesno(question string) bool {
	ui.print(question+ " (Y/N)? ", false)

	for {
		r := unicode.ToUpper(ui.rl.Terminal.ReadRune())
		ui.print("Answer is " + string(r), true)
		switch r {
		case 'Y':
			return true
		case 'N':
			return false
		}
	}
}

func (ui *userInterface) exit() {
	if ui.yesno("Are you sure you want to quit") {
		os.Exit(0)
	}
}

func (ui *userInterface) getAnswer() string {
	for {
		line, err := ui.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			ui.exit()
		}
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				ui.rl.SetVimMode(true)
			case "emacs":
				ui.rl.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case line == "mode":
			if ui.rl.IsVimMode() {
				println("current mode: vim")
			} else {
				println("current mode: emacs")
			}
		case line == "exit":
			ui.exit()
		case line == "":
		default:
			return line
		}
	}
}
