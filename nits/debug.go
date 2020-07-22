package nits

// This file contains the NITS debugger. It implements UI commands that
// can be used to get an insight into the internals of NITS.

import (
	"fmt"
	"os/exec"
)

// A global tracer. When non-nil this can be used to get some debugging
// output during the operation of NITS (inconsistently implemented).
var trace *tracer

type tracer struct {
	ui *userInterface
}

func (t *tracer) print(s string, args ...interface{}) {
	if t == nil {
		return
	}
	t.ui.print(s, args...)
}

func (t *tracer) println(s string, args ...interface{}) {
	if t == nil {
		return
	}
	t.ui.println(s, args...)
}

// displayAnswers is a UI command that displays all registered answers in the
// student state.
func displayAnswers(ui *userInterface, state *studentState) {
	ui.println("Registered answers:")
	ui.newline()

	for _, answer := range state.answers {
		ui.println("%5t: %s#%s", answer.correct, answer.questionShortName, answer.subQuestion)
	}
}

// displayQuestions is a UI command that displays all questions in the content.
func displayQuestions(ui *userInterface, state *studentState) {
	ui.println("All questions:")
	ui.newline()

	for _, q := range state.content.Questions {
		var burnt string
		if _, ok := state.burnt[q]; ok {
			burnt = "[burnt]"
		}
		ui.println("%32s %s", q.getShortName(), burnt)
	}
}

// runShell is a helper that executes a shell command.
func runShell(shellCommand string) error {
	cmd := exec.Command("/bin/bash", "-c", shellCommand)
	_, err := cmd.Output()
	return err
}

// showDot is a helper that runs the GraphViz dot command. When run on a Mac
// it will actually show the generated file in Preview.
func showDot(ui *userInterface, fname string) error {
	ui.println("Dot file: %s", fname)
	if err := runShell(fmt.Sprintf("dot -Tpdf %s >/tmp/nitsdot.pdf && test `uname` = Darwin && open /tmp/nitsdot.pdf", fname)); err != nil {
		return err
	}
	return nil
}

// nextQuetion is a UI commmand that allows the user to manually set the next
// question to ask.
func nextQuestion(ui *userInterface, state *studentState, words []string) {
	// No argument specified? Reset the next question field.
	if len(words) <2 {
		ui.println("Next question reset to nil.")
		state.nextQuestion = nil
	} else if q := state.content.findQuestion(words[1]); q == nil {
		ui.error("Question not found.")
	} else {
		state.nextQuestion = q
	}
}

// debug is the NITS debugger UI command.
func debug(ui *userInterface, state *studentState, words []string) bool {
	ui.pushCommandContext(&CommandContext{
		description: "NITS debugger",
		commands: []*Command{
			{
				aliases: []string{"questions"},
				help:    "Displays all known questions.",
				executor: func([]string) bool {
					displayQuestions(ui, state)
					return false
				},
			},
			{
				aliases: []string{"next"},
				help:    "Selects the next question (by short name).",
				executor: func(words []string) bool {
					nextQuestion(ui, state, words)
					return false
				},
			},
			{
				aliases: []string{"answers"},
				help:    "Displays all registered answers.",
				executor: func([]string) bool {
					displayAnswers(ui, state)
					return false
				},
			},
			{
				aliases: []string{"done"},
				help:    "Exits the debugger.",
				executor: func([]string) bool {
					return true
				},
			},
			{
				aliases: []string{"train"},
				help:    "Runs trainhmm.",
				executor: func([]string) bool {
					if len(state.answers) == 0 {
						ui.println("Nothing to train.")
						return false
					}
					td, err := state.writeTrainhmmInput()
					ui.println("Temporary directory: td=%s; err=%v", td, err)
					if err != nil {
						return false
					}
					err = state.runTrainhmm(td)
					ui.println("Training: err=%v", err)
					err = state.readPrediction(td)
					ui.println("Reading predictions: err=%v", err)
					for concept, skillLevel := range state.scores {
						ui.println("%-20s: %f", concept.shortName, skillLevel)
					}
					return false
				},
			},
			{
				aliases: []string{"scores"},
				help:    "Shows all concepts and the student's skill levels.",
				executor: func([]string) bool {
					if len(state.scores) == 0 {
						ui.println("This student is from Barcelona.")
						return false
					}
					for concept, skillLevel := range state.scores {
						ui.println("%-20s: %f", concept.shortName, skillLevel)
					}
					return false
				},
			},
			{
				aliases: []string{"dot"},
				help:    "Generates dot file for content and concepts.",
				executor: func([]string) bool {
					annotate, ret := ui.yesNo("Augment concepts with skill level ratios")
					if ret {
						return ret
					}
					fname, err := makeDot(state, annotate)
					if err != nil {
						ui.error("error: %s", err)
						return false
					}
					if err := showDot(ui, fname); err != nil {
						ui.error("error: %s", err)
					}
					return false
				},
			},
			{
				aliases: []string{"casedot"},
				help:    "Generates dot file for a case.",
				executor: func(words []string) bool {
					fname, err := makeCaseDot(state, words)
					if err != nil {
						ui.error("error: %s", err)
						return false
					}
					if err := showDot(ui, fname); err != nil {
						ui.error("error: %s", err)
					}
					return false
				},
			},
			{
				aliases: []string{"trace"},
				help:    "Switching detail tracing on or off",
				executor: func(words []string) bool {
					state := "on"
					if len(words) == 1 {
						if trace == nil {
							state = "off"
						}
					} else if words[1] == "on" {
						trace = &tracer{ui}
					} else if words[1] == "off" {
						trace = nil
						state = "off"
					} else {
						ui.error("on or off, please")
						return false
					}
					ui.println("tracing is %s", state)
					return false
				},
			},
			{
				aliases: []string{"select"},
				help:    "Run the question selection algorithm",
				executor: func([]string) bool {
					q := state.selectQuestion()
					if q == nil {
						ui.error("No question selected!")
					}
					return false
				},
			},
		},
	})
	ui.pushPrompt("debug> ")
	defer ui.popPrompt()
	defer ui.popCommandContext()

	// Was a debugging command given on the debug command itself?
	// E.g: debug train.
	// If so, execute it.
	if len(words) > 1 {
		var ret, didExec bool
		if didExec, ret = ui.maybeExecuteCommand(words[1:]); !didExec {
			ui.error("Unknown debugging command.")
		}
		return ret
	}

	for {
		_, ret := ui.getInput()
		if ret {
			return ret
		}
		ui.error("Please enter one of the debugger's commands. Use ? for help.")
	}
}
