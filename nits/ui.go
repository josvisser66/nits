package nits

import (
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

// --------------------------------------------------------------------
type Command struct {
	Aliases  []string
	Global   bool
	Help     string
	Executor func([]string) bool
}

func (c *Command) matches(cmd string) bool {
	for _, alias := range c.Aliases {
		if cmd == alias {
			return true
		}
	}

	return false
}

type CommandContext struct {
	Description string
	Commands    []*Command
}

type userInterface struct {
	rl                  *readline.Instance
	column              int
	promptStack         []string
	commandContextStack []*CommandContext
}

func (ui *userInterface) pushCommandContext(ctx *CommandContext) {
	ui.commandContextStack = append(ui.commandContextStack, ctx)
}

func (ui *userInterface) popCommandContext() {
	ui.commandContextStack = ui.commandContextStack[:len(ui.commandContextStack)-1]
}

func (ui *userInterface) giveHelp() {
	top := true
	for i := len(ui.commandContextStack) - 1; i >= 0; i-- {
		ctx := ui.commandContextStack[i]
		ui.println("Context: %s", ctx.Description)

		for _, cmd := range ctx.Commands {
			if top || cmd.Global {
				ui.println("  %s: %s", strings.Join(cmd.Aliases, "|"), cmd.Help)
			}
		}

		top = false
		ui.newline()
	}
}

func (ui *userInterface) maybeExecuteCommand(line []string) (bool, bool) {
	if len(line) == 0 {
		return false, false
	}
	if line[0] == "?" {
		ui.giveHelp()
		return true, false
	}
	top := true
	for i := len(ui.commandContextStack) - 1; i >= 0; i-- {
		for _, cmd := range ui.commandContextStack[i].Commands {
			if (top || cmd.Global) && cmd.matches(line[0]) {
				return true, cmd.Executor(line)
			}
		}
		top = false
	}

	return false, false
}

func (ui *userInterface) pushPrompt(t string) {
	ui.promptStack = append(ui.promptStack, t)
	ui.rl.SetPrompt(t)
}

func (ui *userInterface) popPrompt() {
	ui.promptStack = ui.promptStack[:len(ui.promptStack)-1]
	ui.rl.SetPrompt(ui.promptStack[len(ui.promptStack)-1])
}

func getUserHomeDir() string {
	if s, err := os.UserHomeDir(); err != nil {
		panic(err)
	} else {
		return s
	}
}

func newUserInterface() *userInterface {
	var err error

	ui := &userInterface{
		promptStack: make([]string, 0, 10),
	}
	ui.rl, err = readline.NewEx(&readline.Config{
		HistoryFile: path.Join(getUserHomeDir(), ".nits_readline"),
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

func (ui *userInterface) newline() {
	ui.rl.Terminal.PrintRune('\n')
	ui.column = 0
}

func (ui *userInterface) println(s string, args ...interface{}) {
	ui.print(s, args...)
	ui.newline()
}

func (ui *userInterface) error(s string, args ...interface{}) {
	ui.print("*** ")
	ui.print(s, args...)
	ui.newline()
}

func (ui *userInterface) print(s string, args ...interface{}) {
	t := fmt.Sprintf(s, args...)
	term := ui.rl.Terminal
	width := term.GetConfig().FuncGetWidth()
	words := strings.Split(t, " ")

	for i, word := range words {
		l := len(word)
		if l+ui.column+1 > width {
			term.PrintRune('\n')
			ui.column = 0
		}
		if ui.column > 0 && i > 0 {
			term.PrintRune(' ')
		}
		ui.rl.Terminal.Print(word)
		ui.column += l + 1
	}
}

func (ui *userInterface) printParagraphs(p []string) {
	for i, t := range p {
		ui.println(t)

		if i < len(p)-1 {
			ui.newline()
		}
	}
}

func (ui *userInterface) printAnswers(answers []*Answer) {
	r := 'A'

	for _, a := range answers {
		ui.println("%c) %s", r, a.Text)
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

func isYesNo(line string) (bool, error) {
	switch strings.TrimSpace(strings.ToLower(line)) {
	case "y":
		fallthrough
	case "yes":
		return true, nil
	case "n":
		fallthrough
	case "no":
		return false, nil
	}
	return false, errors.New("not a yes or no answer")
}

func (ui *userInterface) yesNo(question string) (bool, bool) {
	ui.pushPrompt(question + " (Y/N)? ")
	defer ui.popPrompt()

	answer, ret := ui.getAnswer(
		answerMap{
			"yes": []string{"y(es)?"},
			"no":  []string{"no?"},
		})
	return answer == "yes", ret
}

func (ui *userInterface) getInput() ([]string, bool) {
	for {
		line, err := ui.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			ui.error("Please use exit to leave NITS.")
		}
		words := strings.Split(strings.ToLower(strings.TrimSpace(line)), " ")
		if len(words) == 0 {
			continue
		}

		didExec, ret := ui.maybeExecuteCommand(words)
		if ret {
			return words, ret
		}

		if didExec {
			continue
		}

		return words, false
	}
}

type answerMap map[string][]string

func (ui *userInterface) getAnswer(answers answerMap) (string, bool) {
	for {
		words, ret := ui.getInput()
		if ret {
			return "", ret
		}
		if len(words) != 1 {
			ui.error("Please provide a one-word answer.")
			continue
		}
		a := strings.ToLower(words[0])
		for key, vec := range answers {
			for _, re := range vec {
				r, err := regexp.Compile(re)
				CHECK(err == nil, "regexp compilation error '%s': %s", re, err)
				if r.MatchString(a) {
					return key, false
				}
			}
		}
		ui.error("Invalid answer. Please try again.")
	}
}

func (ui *userInterface) explain(e *Explanation) {
	if e == nil {
		ui.println("Unfortunately there is no explanation available for this topic. :-(")
		return
	}

	ui.printParagraphs(e.Text)
	ui.newline()

	if e.References != nil {
		ui.println("References:")
		ui.newline()

		for _, r := range e.References {
			ui.print("- %s", r.GetReferenceText())
		}

		ui.newline()
	}
}
