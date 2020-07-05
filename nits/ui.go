package nits

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"log"
	"strings"
)

// --------------------------------------------------------------------
type Command struct {
	Aliases  []string
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
	for i := len(ui.commandContextStack) - 1; i >= 0; i-- {
		ctx := ui.commandContextStack[i]

		for _, cmd := range ctx.Commands {
			ui.println("%s: %s", strings.Join(cmd.Aliases, "|"), cmd.Help)
		}
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
	for i := len(ui.commandContextStack) - 1; i >= 0; i-- {
		for _, cmd := range ui.commandContextStack[i].Commands {
			if cmd.matches(line[0]) {
				return true, cmd.Executor(line)
			}
		}
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

func (ui *userInterface) newline() {
	ui.rl.Terminal.PrintRune('\n')
	ui.column = 0
}

func (ui *userInterface) println(s string, args ...interface{}) {
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

func (ui *userInterface) yesNo(question string) bool {
	ui.pushPrompt(question + " (Y/N)? ")
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

func (ui *userInterface) getInput() ([]string, bool) {
	for {
		line, err := ui.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			ui.println("Please use exit to leave NITS.")
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