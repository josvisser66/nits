package nits

import (
	"github.com/chzyer/readline"
	"io"
	"log"
	"strings"
)

// --------------------------------------------------------------------
type Command struct {
	Aliases  []string
	Help     string
	Executor func([]string)
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
	Commands []*Command
}

type userInterface struct {
	rl                  *readline.Instance
	column int
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
			ui.print(strings.Join(cmd.Aliases, "|"), false)
			ui.print(": ", false)
			ui.print(cmd.Help, true)
		}
	}
}

func (ui *userInterface) maybeExecuteCommand(line []string) bool {
	if len(line) == 0 { return false }
	if line[0] == "?" {
		ui.giveHelp()
		return true
	}
	for i := len(ui.commandContextStack) - 1; i >= 0; i-- {
		for _, cmd := range ui.commandContextStack[i].Commands {
			if cmd.matches(line[0]) {
				cmd.Executor(line)
				return true
			}
		}
	}

	return false
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
}

func (ui *userInterface) print(t string, newline bool) {
	term := ui.rl.Terminal
	width := term.GetConfig().FuncGetWidth()
	words := strings.Split(t, " ")

	for _, word := range words {
		l := len(word)
		if l+ui.column+1 > width {
			term.PrintRune('\n')
			ui.column = 0
		}
		if ui.column > 0 {
			term.PrintRune(' ')
		}
		ui.rl.Terminal.Print(word)
		ui.column += l + 1
	}

	if newline {
		term.PrintRune('\n')
		ui.column = 0
	}
}

func (ui *userInterface) printAnswers(answers []*Answer) {
	r := 'A'

	for _, a := range answers {
		ui.print(string(r)+". ", false)
		ui.print(a.Text, true)
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

func (ui *userInterface) getAnswer(displayQuestion func()) string {
	ui.pushCommandContext(&CommandContext{
		"Answering a question",
		[]*Command{
			{
				[]string{"again"},
				"Displays the question again.",
				func(i []string) {
					displayQuestion()
				},
			},
		},
	})
	defer ui.popCommandContext()

	for {
		line, err := ui.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			ui.print("Please use exit to leave NITS.", true)
		}
		words := strings.Split(strings.ToLower(strings.TrimSpace(line)), " ")
		if len(words) == 0 || ui.maybeExecuteCommand(words) {
			continue
		}
		return line
	}
}
