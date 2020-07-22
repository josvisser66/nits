package nits

// This file contains NITS simple text based UI implementation.

import (
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

// Command is a command that can control the NITS the UI. At every input
// opportunity the user can input a command and have it executed, after
// which the UI returns to obtaining input.
type Command struct {
	aliases  []string
	global   bool // If false this command is only available in the toplevel context.
	help     string
	executor func([]string) bool
}

// matches checks if a string matches on of the command's aliases.
func (c *Command) matches(cmd string) bool {
	for _, alias := range c.aliases {
		if cmd == alias {
			return true
		}
	}

	return false
}

// CommandContext is a set of commands that is valid in a given input
// context.
type CommandContext struct {
	description string
	commands    []*Command
}

// userInterface is the abstract representation of the text based ui.
type userInterface struct {
	rl                  *readline.Instance
	column              int
	promptStack         []string
	commandContextStack []*CommandContext
}

// pushCommandContext pushes a command context onto the stack.
func (ui *userInterface) pushCommandContext(ctx *CommandContext) {
	ui.commandContextStack = append(ui.commandContextStack, ctx)
}

// popCommandContext pops a command context from the stack.
func (ui *userInterface) popCommandContext() {
	ui.commandContextStack = ui.commandContextStack[:len(ui.commandContextStack)-1]
}

// giveHelp shows all the commands that are available at this point.
func (ui *userInterface) giveHelp() {
	top := true
	for i := len(ui.commandContextStack) - 1; i >= 0; i-- {
		ctx := ui.commandContextStack[i]
		ui.println("Context: %s", ctx.description)

		for _, cmd := range ctx.commands {
			if top || cmd.global {
				ui.println("  %s: %s", strings.Join(cmd.aliases, "|"), cmd.help)
			}
		}

		top = false
		ui.newline()
	}
}

// maybeExecute command takes a line (split by spaces into words) and tries
// to execute it as a command. If it does the first return value is true.
// The second return value is true if the command wants the top level context
// to terminate.
func (ui *userInterface) maybeExecuteCommand(words []string) (bool, bool) {
	if len(words) == 0 {
		return false, false
	}
	if words[0] == "?" {
		ui.giveHelp()
		return true, false
	}
	top := true
	for i := len(ui.commandContextStack) - 1; i >= 0; i-- {
		for _, cmd := range ui.commandContextStack[i].commands {
			if (top || cmd.global) && cmd.matches(words[0]) {
				return true, cmd.executor(words)
			}
		}
		top = false
	}

	return false, false
}

// pushPrompt pushes a prompt string onto the prompt stack.
func (ui *userInterface) pushPrompt(t string) {
	ui.promptStack = append(ui.promptStack, t)
	ui.rl.SetPrompt(t)
}

// popPrompt pops a prompt string from the prompt stack.
func (ui *userInterface) popPrompt() {
	ui.promptStack = ui.promptStack[:len(ui.promptStack)-1]
	ui.rl.SetPrompt(ui.promptStack[len(ui.promptStack)-1])
}

// mustUserHomeDir returns the location of the user's home directory
// or panics.
func mustUserHomeDir() string {
	if s, err := os.UserHomeDir(); err != nil {
		panic(err)
	} else {
		return s
	}
}

// newUserInterface creates a new UI driver object.
func newUserInterface() *userInterface {
	var err error

	ui := &userInterface{
		promptStack: make([]string, 0, 10),
	}
	ui.rl, err = readline.NewEx(&readline.Config{
		HistoryFile: path.Join(mustUserHomeDir(), ".nits_readline"),
		// AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}

	ui.pushPrompt("Your wish is my command! ")
	log.SetOutput(ui.rl.Stderr())
	return ui
}

// newline generates a newline onto the output stream.
func (ui *userInterface) newline() {
	ui.rl.Terminal.PrintRune('\n')
	ui.column = 0
}

// println prints a string and then generates a newline.
func (ui *userInterface) println(s string, args ...interface{}) {
	ui.print(s, args...)
	ui.newline()
}

// error formats an error message onto the output stream. An error
// message is followed by a newline.
func (ui *userInterface) error(s string, args ...interface{}) {
	ui.print("*** ")
	ui.print(s, args...)
	ui.newline()
}

// print prints a string to the output stream, taking the terminal
// width into account and ensuring that we are not breaking words
// in the middle.
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

// printParagraphs prints a vector of strings as paragraphs. Each
// string is a paragraph. Paragraphs are separated by a blank line.
func (ui *userInterface) printParagraphs(p []string) {
	for i, t := range p {
		ui.println(t)

		if i < len(p)-1 {
			ui.newline()
		}
	}
}

// printAnswers prints the answer set of a multiple choice
// question.
func (ui *userInterface) printAnswers(answers []*Answer) {
	r := 'A'

	for _, a := range answers {
		ui.println("%c) %s", r, a.Text)
		r += 1
	}
}

// filterInput allows us to remove input characters from
// consideration.
func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// yesNo gets a yes/no answer from the user.
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

// getInput returns a line of input (split into lower case words on
// spaces), while executing any commands that are valid in the context.
func (ui *userInterface) getInput() ([]string, bool) {
	for {
		line, err := ui.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			ui.error("Please use exit to leave NITS.")
		}
		if line == "" {
			continue
		}
		words := strings.Split(strings.ToLower(strings.TrimSpace(line)), " ")
		didExec, ret := ui.maybeExecuteCommand(words)
		if ret {
			return words, ret
		}

		if !didExec {
			return words, false
		}
	}
}

// answerMap is a map of allowed answers.
type answerMap map[string][]string

// getAnswer takes input from the user until that input is one of
// the answers in the answer map (while executing any commands that
// the user enters). The return value is the key of the entry in the
// answer map that the user entered. The second return value is a boolean
// to indicate that a command wants the calling context to terminate.
func (ui *userInterface) getAnswer(answers answerMap) (string, bool) {
	for {
		words, ret := ui.getInput()
		if ret {
			return "", ret
		}
		if len(words) == 0 {
			continue
		}
		if len(words) > 1 {
			ui.error("Please provide a one-word answer.")
			continue
		}
		a := strings.ToLower(words[0])
		for key, vec := range answers {
			for _, re := range vec {
				if regexp.MustCompile(fmt.Sprintf("^%s$", re)).MatchString(a) {
					return key, false
				}
			}
		}
		ui.error("Invalid answer. Please try again.")
	}
}

// explain print an Explanation.
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
