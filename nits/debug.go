package nits

import (
	"fmt"
	"os/exec"
)

var (
	trace *tracer
)

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

func displayAnswers(ui *userInterface, state *studentState) {
	ui.println("Registered answers:")
	ui.newline()

	for _, answer := range state.answers {
		ui.println("%5t: %s#%s", answer.correct, answer.questionShortName, answer.subQuestion)
	}
}

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

func runShell(shellCommand string) error {
	cmd := exec.Command("/bin/bash", "-c", shellCommand)
	_, err := cmd.Output()
	return err
}

func showDot(ui *userInterface, fname string) error {
	ui.println("Dot file: %s", fname)
	if err := runShell(fmt.Sprintf("dot -Tpdf %s >/tmp/aap.pdf && open /tmp/aap.pdf", fname)); err != nil {
		return err
	}
	return nil
}

func nextQuestion(ui *userInterface, state *studentState, words []string) {
	if len(words) <2 {
		state.nextQuestion = nil
	} else if q := state.content.findQuestion(words[1]); q == nil {
		ui.println("Question not found.")
	} else {
		state.nextQuestion = q
	}
}

func debug(ui *userInterface, state *studentState, words []string) bool {
	ui.pushCommandContext(&CommandContext{
		Description: "NITS debugger",
		Commands: []*Command{
			{
				Aliases: []string{"questions"},
				Help:    "Displays all known questions.",
				Executor: func([]string) bool {
					displayQuestions(ui, state)
					return false
				},
			},
			{
				Aliases: []string{"next"},
				Help:    "Selects the next question (by short name).",
				Executor: func(words []string) bool {
					nextQuestion(ui, state, words)
					return false
				},
			},
			{
				Aliases: []string{"answers"},
				Help:    "Displays all registered answers.",
				Executor: func([]string) bool {
					displayAnswers(ui, state)
					return false
				},
			},
			{
				Aliases: []string{"done"},
				Help:    "Exits the debugger.",
				Executor: func([]string) bool {
					return true
				},
			},
			{
				Aliases: []string{"train"},
				Help:    "Runs trainhmm.",
				Executor: func([]string) bool {
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
				Aliases: []string{"scores"},
				Help:    "Shows all concepts and the student's skill levels.",
				Executor: func([]string) bool {
					for concept, skillLevel := range state.scores {
						ui.println("%-20s: %f", concept.shortName, skillLevel)
					}
					return false
				},
			},
			{
				Aliases: []string{"dot"},
				Help:    "Generates dot file for content and concepts.",
				Executor: func([]string) bool {
					fname, err := makeDot(state, ui.yesNo("Augment concepts with skill level ratios"))
					if err != nil {
						ui.println("error: %s", err)
						return false
					}
					if err := showDot(ui, fname); err != nil {
						ui.println("error: %s", err)
					}
					return false
				},
			},
			{
				Aliases: []string{"casedot"},
				Help:    "Generates dot file for a case.",
				Executor: func(words []string) bool {
					fname, err := makeCaseDot(state, words)
					if err != nil {
						ui.println("error: %s", err)
						return false
					}
					if err := showDot(ui, fname); err != nil {
						ui.println("error: %s", err)
					}
					return false
				},
			},
			{
				Aliases: []string{"trace"},
				Help:    "Switching detail tracing on or off",
				Executor: func(words []string) bool {
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
						ui.println("on or off, please")
						return false
					}
					ui.println("tracing is %s", state)
					return false
				},
			},
			{
				Aliases: []string{"select"},
				Help:    "Run the question selection algorithm",
				Executor: func([]string) bool {
					q := state.selectQuestion()
					if q == nil {
						ui.println("No question selected!")
					}
					return false
				},
			},
		},
	})
	ui.pushPrompt("debug> ")
	defer ui.popPrompt()
	defer ui.popCommandContext()

	if len(words) > 1 {
		var ret, didExec bool
		if didExec, ret = ui.maybeExecuteCommand(words[1:]); !didExec {
			ui.println("Unknown debugging command.")
		}
		return ret
	}

	for {
		_, ret := ui.getInput()
		if ret {
			return ret
		}
		ui.println("Please enter one of the debugger's commands. Use ? for help.")
	}
}
