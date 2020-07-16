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
	t.ui.print(s, args...)
}

func (t *tracer) println(s string, args ...interface{}) {
	t.ui.println(s, args...)
}

func displayAnswers(ui *userInterface, state *studentState) {
	ui.println("Registered answers:")
	ui.newline()

	for _, answer := range state.answers {
		ui.println("%5t: %s", answer.correct, answer.questionShortName)
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

func debug(ui *userInterface, state *studentState) {
	ui.pushCommandContext(&CommandContext{
		Description: "NITS debugger",
		Commands: []*Command{
			{
				Aliases: []string{"answers"},
				Help:    "Displays all registered answers.",
				Executor: func(strings []string) bool {
					displayAnswers(ui, state)
					return false
				},
			},
			{
				Aliases: []string{"done"},
				Help:    "Exits the debugger.",
				Executor: func(strings []string) bool {
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

	for {
		_, ret := ui.getInput()
		if ret {
			return
		}
		ui.println("Please enter one of the debugger's commands. Use ? for help.")
	}
}
